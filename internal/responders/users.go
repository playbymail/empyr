// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package responders

import (
	"encoding/json"
	"github.com/playbymail/empyr/internal/domains"
	"html/template"
	"log"
	"net/http"
)

type CreateUserResponder struct {
	IsHTMX bool
	Tmpl   *template.Template // Injected template for rendering
}

func (r *CreateUserResponder) Respond(w http.ResponseWriter, user domains.User, err error) {
	log.Printf("%s %s: %s\n", "?", "?", "current responder")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%s %s: htmx %q\n", "?", "?", w.Header().Get("HX-Request"))

	// Detect if the request is from HTMX
	// if htmx := w.Header().Get("HX-Request"); htmx != "" {
	if r.IsHTMX {
		// Render partial HTML
		w.Header().Set("Content-Type", "text/html")
		r.Tmpl.ExecuteTemplate(w, "user-row.gohtml", user)
		return
	}

	// Default to JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
