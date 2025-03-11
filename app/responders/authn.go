// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package responders

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type ShowLoginResponder struct {
	Tmpl *template.Template
}

type ShowLoginResponse struct {
	IsHTMX bool
	User   string
}

func (r *ShowLoginResponder) Respond(w http.ResponseWriter, data ShowLoginResponse, err error) {
	w.Header().Set("Content-Type", "text/html")
	_ = r.Tmpl.ExecuteTemplate(w, "login.gohtml", nil)
}

type LogoutResponder struct {
	Templates []string // path to templates directory
	Tmpl      *template.Template
}

type LogoutResponse struct {
	IsHTMX bool
}

func NewLogoutResponder(templates string) *LogoutResponder {
	// in development, always reload the templates. in production, we should cache them.
	return &LogoutResponder{
		Templates: []string{filepath.Join(templates, "logout.gohtml")},
	}
}

func (r *LogoutResponder) Respond(w http.ResponseWriter, data LogoutResponse, err error) {
	log.Printf("responders: LogoutResponder: entered\n")

	// delete the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "empyr-session",
		Value:  "",
		MaxAge: -1,
	})

	var t *template.Template
	if r.Tmpl != nil {
		log.Printf("responders: LogoutResponder: using cached template\n")
		t = r.Tmpl
	} else {
		log.Printf("responders: LogoutResponder: loading %s\n", r.Templates)
		if t, err = template.ParseFiles(r.Templates...); err != nil {
			log.Printf("responders: LogoutResponder: template.ParseFiles: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
	_ = t.ExecuteTemplate(w, "logout.gohtml", nil)
}
