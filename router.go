package main

import (
	"fmt"
	"log"
	"net/http"
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

	sql := `CREATE TABLE IF NOT EXISTS routes (
            id INTEGER NOT NULL PRIMARY KEY,
            path TEXT NOT NULL UNIQUE,
            dest TEXT NOT NULL UNIQUE,
            comment TEXT,
            counter INTEGER DEFAULT 0
          );`
	_, err := db.Exec(sql)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sql)
		return nil
	}

	rows, err := db.Query("SELECT id, path, dest FROM routes")
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
			r.makeRedirect(rw, req, route)
			fmt.Printf("action=matched path=%s dest=%s\n", route.path, route.dest)
			return
		}
	}
	// no path matched; send 404 response
	http.NotFound(rw, req)
}

func (r *Router) makeRedirect(rw http.ResponseWriter, req *http.Request, route *Route) {
	http.Redirect(rw, req, route.dest, 302)
	go r.incCounter(route)
}

func (r *Router) incCounter(route *Route) {
	stmt, err := db.Prepare("UPDATE routes SET counter = counter + 1 WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(route.id)
	if err != nil {
		log.Fatal(err)
	}
}
