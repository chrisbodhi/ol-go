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
	businesses, err := AllBusinesses()
	if err != nil {
		panic(err)
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
	Address2  *string   `json:"address2"`
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

// Database funny business
func AllBusinesses() ([]*Business, error) {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM businesses WHERE id < 5")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	businesses := make([]*Business, 0)
	for rows.Next() {
		biz := new(Business)
		err := rows.Scan(&biz.Id, &biz.Uuid, &biz.Name, &biz.Address, &biz.Address2, &biz.City, &biz.State, &biz.Zip, &biz.Country, &biz.Phone, &biz.Website, &biz.CreatedAt, &biz.UpdatedAt)
		if err != nil {
			panic(err)
		}
		businesses = append(businesses, biz)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return businesses, nil
}
