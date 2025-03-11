// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package views

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"net/http"
)

type View struct {
	assetsFS  FS
	viewsFS   FS
	name      string
	path      string
	templates *template.Template
}

type FS struct {
	FS   embed.FS
	Path string
}

func NewView(name string, path string) (*View, error) {
	if _, err := template.ParseFiles(path); err != nil {
		return nil, err
	}
	return &View{
		name: name,
		path: path,
	}, nil
}

func New(assetsFS, viewsFS FS) *View {
	return &View{
		assetsFS:  assetsFS,
		viewsFS:   viewsFS,
		templates: template.Must(template.ParseFS(viewsFS.FS, "views/*.gohtml")),
	}
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	log.Printf("%s %s: rendering template %q\n", r.Method, r.URL.Path, name)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	t := v.templates
	if t == nil {
		// in development, we want to reload the templates on each request
		var err error
		t, err = template.ParseFiles(v.path)
		if err != nil {
			log.Printf("%s %s: %v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Printf("%s %s: parsed components\n", r.Method, r.URL.Path)
	}

	// parse into a buffer so that we can handle errors without writing to the response
	buf := &bytes.Buffer{}
	if err := t.ExecuteTemplate(buf, v.name, data); err != nil {
		log.Printf("%s %s: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// write to the response
	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("%s %s: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
