// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package responders

import (
	"encoding/json"
	"github.com/playbymail/empyr/app/domains"
	"html/template"
	"log"
	"net/http"
)

type CreateUserResponder struct {
	Tmpl *template.Template // Injected template for rendering
}

// CreateUserResponse encapsulates the data used to respond to CreateUser request.
type CreateUserResponse struct {
	IsHTMX bool
	User   domains.User
}

func (r *CreateUserResponder) Respond(w http.ResponseWriter, data CreateUserResponse, err error) {
	log.Printf("%s %s: %s\n", "?", "?", "current responder")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if data.IsHTMX { // render partial HTML for HTMX requests
		w.Header().Set("Content-Type", "text/html")
		_ = r.Tmpl.ExecuteTemplate(w, "user-row.gohtml", data.User)
		return
	}

	// Default to JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(data.User)
}
