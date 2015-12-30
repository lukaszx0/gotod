package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Route struct {
	id   int64
	path string
	dest string
}

type Router struct {
	routes []*Route
}

func NewRouter() *Router {
	router := &Router{}

	db, err := sql.Open("sqlite3", "./goto.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := `CREATE TABLE IF NOT EXISTS redirects (
            id INTEGER NOT NULL PRIMARY KEY,
            path TEXT,
            dest TEXT,
            comment TEXT
          );`
	_, err = db.Exec(sql)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sql)
		return nil
	}

	rows, err := db.Query("SELECT id, path, dest FROM redirects")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var path string
		var dest string

		rows.Scan(&id, &path, &dest)
		router.AddRoute(id, path, dest)
	}

	fmt.Printf("Loaded %d paths from database\n", len(router.routes))

	return router
}

func (r *Router) AddRoute(id int64, path string, dest string) {
	r.routes = append(r.routes, &Route{id, path, dest})
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if route.path == req.URL.Path {
			r.makeRedirect(rw, req, route.dest)
			fmt.Printf("action=matched path=%s dest=%s\n", route.path, route.dest)
			return
		}
	}
	// no path matched; send 404 response
	http.NotFound(rw, req)
}

func (r *Router) makeRedirect(rw http.ResponseWriter, req *http.Request, dest string) {
	http.Redirect(rw, req, dest, 302)
}
