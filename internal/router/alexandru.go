// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package router implements a chi-like router based on a gist from alexandru.
//
// Chi-like syntactic sugar layer on top of stdlib http.ServeMux.
// https://gist.github.com/alexaandru/747f9d7bdfb1fa35140b359bf23fa820
package router

import (
	"net/http"
	"slices"
)

// Router implements a Chi-like interface for routing over the stdlib.
type Router struct {
	*http.ServeMux
	chain []Middleware
}

type Middleware func(http.Handler) http.Handler

// New returns a Router that uses the given middleware.
func New(mx ...Middleware) *Router {
	return &Router{ServeMux: &http.ServeMux{}, chain: mx}
}

// Use adds middleware to the router
func (r *Router) Use(mx ...Middleware) {
	r.chain = append(r.chain, mx...)
}

// Group provides a way to group routes together so they can use a common set of middleware.
//
// It creates a new router, setting that router's chain to a clone of the current router's chain.
// It then calls the provided function, passing in the new router as a parameter. It's assumed
// this function will add middleware and handlers to the new router.
//
// The mux for the new router is the same as for the current router. All routes that are added
// to the new router will be served by the current router's mux.
func (r *Router) Group(fn func(gr *Router)) {
	fn(&Router{ServeMux: r.ServeMux, chain: slices.Clone(r.chain)})
}

// Delete is a helper function for the HTTP DELETE method that wraps calls to the handler with the given middleware.
func (r *Router) Delete(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodDelete, path, fn, mx)
}

// Get is a helper function for the HTTP GET method that wraps calls to the handler with the given middleware.
func (r *Router) Get(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodGet, path, fn, mx)
}

// Head is a helper function for the HTTP HEAD method that wraps calls to the handler with the given middleware.
func (r *Router) Head(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodHead, path, fn, mx)
}

// Options is a helper function for the HTTP OPTIONS method that wraps calls to the handler with the given middleware.
func (r *Router) Options(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodOptions, path, fn, mx)
}

// Post is a helper function for the HTTP POST method that wraps calls to the handler with the given middleware.
func (r *Router) Post(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodPost, path, fn, mx)
}

// Put is a helper function for the HTTP PUT method that wraps calls to the handler with the given middleware.
func (r *Router) Put(path string, fn http.HandlerFunc, mx ...Middleware) {
	r.handle(http.MethodPut, path, fn, mx)
}

// handle injects the method into the route for the stdlib's use and wraps the handler
// with the given middleware.
func (r *Router) handle(method, path string, fn http.HandlerFunc, mx []Middleware) {
	r.Handle(method+" "+path, r.wrap(fn, mx))
}

// wrap reverses the order of the middleware so that they'll be called right to left
// (or inside-out, if you prefer to view it that way).
//
// It does more than that. It wraps the middleware around the handler and includes
// the current chain in there, too.
func (r *Router) wrap(fn http.HandlerFunc, mx []Middleware) (out http.Handler) {
	out, mx = http.Handler(fn), append(slices.Clone(r.chain), mx...)
	slices.Reverse(mx)
	for _, m := range mx {
		out = m(out)
	}
	return out
}
