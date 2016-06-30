package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var page string
	if len(r.URL.Query()["page"]) > 0 {
		page = r.URL.Query()["page"][0]
	} else {
		page = "1"
	}

	businesses, err := AllBusinesses(page)
	if err != nil {
		panic(err)
	}

	if err := json.NewEncoder(w).Encode(businesses); err != nil {
		panic(err)
	}
}

func IdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	businessId := vars["businessId"]

	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		panic(err)
	}

	row, err := db.Query("SELECT * FROM businesses WHERE id=" + businessId)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	business := new(Business)
	for row.Next() {
		biz := new(Business)
		err := row.Scan(&biz.Id, &biz.Uuid, &biz.Name, &biz.Address, &biz.Address2, &biz.City, &biz.State, &biz.Zip, &biz.Country, &biz.Phone, &biz.Website, &biz.CreatedAt, &biz.UpdatedAt)
		if err != nil {
			panic(err)
		}
		business = biz
	}

	if err = row.Err(); err != nil {
		panic(err)
	}

	if err := json.NewEncoder(w).Encode(business); err != nil {
		panic(err)
	}
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
func AllBusinesses(page string) ([]*Business, error) {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		panic(err)
	}

	pageNum, err := strconv.Atoi(page)
	idStart := pageNum * 50
	idEnd := (pageNum * 50) + 50
	sqlQuery := "SELECT * FROM businesses WHERE id >= " + strconv.Itoa(idStart) + " AND id < " + strconv.Itoa(idEnd)

	rows, err := db.Query(sqlQuery)
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
