package main

import (
	"fmt"
	"net/http"
)

type Router struct {
	routes []*Route
}

func NewRouter() *Router {
	router := &Router{}

	routes := loadRoutes()
	for _, route := range routes {
		router.AddRoute(&route)
	}

	return router
}

func (r *Router) AddRoute(route *Route) {
	r.routes = append(r.routes, route)
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if route.Path == req.URL.Path {
			r.makeRedirect(rw, req, route)
			fmt.Printf("action=matched path=%s dest=%s\n", route.Path, route.Dest)
			return
		}
	}
	// no path matched; send 404 response
	http.NotFound(rw, req)
}

func (r *Router) makeRedirect(rw http.ResponseWriter, req *http.Request, route *Route) {
	http.Redirect(rw, req, route.Dest, 302)
	go incRouteCounter(route)
}
