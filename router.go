package main

import (
	"fmt"
	"net/http"
)

type Route struct {
	path string
	dest string
}

type Router struct {
	routes []*Route
}

func NewRouter() *Router {
	router := &Router{}
	router.AddRoute("/google", "http://google.com")
	return router
}

func (r *Router) AddRoute(path string, dest string) {
	r.routes = append(r.routes, &Route{path, dest})
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
