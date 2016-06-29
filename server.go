package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Root of everything
func main() {
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Routes and Routing
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"BusinessHandler",
		"GET",
		"/api/businesses",
		BusinessHandler,
	},

	Route{
		"IdHandler",
		"GET",
		"/api/businesses/{businessId}",
		IdHandler,
	},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

// Handlers
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func BusinessHandler(w http.ResponseWriter, r *http.Request) {
	businesses := Businesses{
		Business{Name: "Fiduciary Planners"},
		Business{Name: "Jokes on Us"},
	}
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM businesses WHERE id < 5")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id sql.NullInt64
		var uuid sql.NullString
		var name sql.NullString
		var address sql.NullString
		var address2 sql.NullString
		var city sql.NullString
		var state sql.NullString
		var zip sql.NullString
		var country sql.NullString
		var phone sql.NullString
		var website sql.NullString
		var created_at sql.NullString
		var updated_at sql.NullString

		if err := rows.Scan(&id, &uuid, &name, &address, &address2, &city, &state, &zip, &country, &phone, &website, &created_at, &updated_at); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ID %d with website %s & name %s zip code %s\n",
			id.Int64, website.String, name.String, zip.String)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(businesses); err != nil {
		panic(err)
	}
}

func IdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	businessId := vars["businessId"]
	fmt.Fprintln(w, "businessId:", businessId)
}

// Model
type Business struct {
	Id        int       `json:"id"`
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Address2  string    `json:"address2"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	Zip       string    `json:"zip"`
	Country   string    `json:"country"`
	Phone     string    `json:"phone"`
	Website   string    `json:"website"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Businesses []Business

// Database funny business

// Logger
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
