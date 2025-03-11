// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package controllers

import (
	"github.com/playbymail/empyr/internal/flash"
	"github.com/playbymail/empyr/internal/views"
	"github.com/playbymail/empyr/store"
	"log"
	"net/http"
)

type Blogs struct {
	db   *store.Store
	view *views.View
}

// NewBlogsController creates a new instance of the Blogs controller
func NewBlogsController(db *store.Store, view *views.View) (*Blogs, error) {
	c := &Blogs{
		db:   db,
		view: view,
	}
	// add any initialization logic here if needed
	return c, nil
}

func (c Blogs) Show(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: %s\n", r.Method, r.URL.Path, r.RemoteAddr)

	store := flash.GetStore(r)

	// - Render the template
	c.view.Render(w, r, "blogs.gohtml", struct {
		Error string
	}{
		Error: store.Get("error"),
	})
}
