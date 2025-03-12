// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package responders

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func NewView(name, path string, files ...string) *ResponderTemplate {
	rt := &ResponderTemplate{Name: name + ".gohtml"}
	for _, file := range files {
		rt.Files = append(rt.Files, filepath.Join(path, file))
	}
	return rt
}

type ResponderTemplate struct {
	Name   string             // name of the base template
	Files  []string           // names of the templates to load
	Cached *template.Template // cached templates
}

func (rt *ResponderTemplate) GetTemplate() (*template.Template, error) {
	// in production, we should cache them.
	if rt.Cached != nil {
		return rt.Cached, nil
	}
	// in development, always reload the templates.
	log.Printf("responders: parsing %v\n", rt.Files)
	return template.ParseFiles(rt.Files...)
}

func (rt *ResponderTemplate) Render(w http.ResponseWriter, route string, data any) {
	log.Printf("responders: %s: rendering %s\n", route, rt.Name)
	t, err := rt.GetTemplate()
	if err != nil {
		log.Printf("responders: %s: %v\n", route, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	buffer := &bytes.Buffer{}
	err = t.ExecuteTemplate(buffer, rt.Name, data)
	if err != nil {
		log.Printf("responders: %s: %v\n", route, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = buffer.WriteTo(w)
}
