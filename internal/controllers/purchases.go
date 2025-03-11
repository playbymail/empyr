// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package controllers

import (
	"github.com/playbymail/empyr/internal/flash"
	"github.com/playbymail/empyr/internal/ratelimiter"
	"net/http"
)

type Purchases struct {
	views   view
	limiter *ratelimiter.Limiter
}

type view struct{}

func (v view) Render(w http.ResponseWriter, name string, data interface{}) {}

func (c Purchases) Show(w http.ResponseWriter, r *http.Request) {
	store := flash.GetStore(r)

	c.views.Render(w, "purchases.html", struct {
		Error string
	}{
		Error: store.Get("error"),
	})
}

func (c Purchases) Download(w http.ResponseWriter, r *http.Request) {
	store := flash.GetStore(r)

	if !c.limiter.Allow(r.URL.Path+r.Header.Get("CF-Connecting-IP"), 5) {
		store.Set("error", "403")
		http.Redirect(w, r, "/purchases", http.StatusFound)
		return
	}
}
