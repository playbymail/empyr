// empyr - a reimagining of Vern Holford's Empyrean Challenge
// Copyright (C) 2023 Michael D Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package handlers

import (
	"log"
	"net/http"
	"os"
	"path"
)

// SPA is a static file server that works with single page applications (Vue, React, Svelte, etc).
//
// If a file is not found, it will return index.html from the root.
// We assume that the caller stripped the path prefix if needed.
//
//	http.StripPrefix("/static", handlers.SPA("/tmp"))
//
// Code based on http.ServeFile, but updated to refuse a directory listing.
func SPA(root string) http.HandlerFunc {
	root = path.Clean(root)
	log.Printf("[spa] root %q\n", root)

	indexFile := root + "/index.html"
	if stat, err := os.Stat(indexFile); err != nil {
		log.Printf("[spa] %q %+v\n", indexFile, err)
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	} else if stat.IsDir() {
		log.Printf("[spa] %q is a directory!", indexFile)
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	} else if !stat.Mode().IsRegular() {
		log.Printf("[spa] %q is not a regular file!", indexFile)
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
	log.Printf("[spa] index %q\n", indexFile)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		file := path.Clean(r.URL.Path)
		if file == "" || file == "/" || file == "." {
			file = "index.html"
		}
		file = path.Clean(root + "/" + file)

		stat, err := os.Stat(file)
		if err != nil {
			// with spa, we must serve the index file when we can't find the requested item
			file = indexFile
			if stat, err = os.Stat(file); err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
		}

		// we never serve directory indexes or irregular files
		if stat.IsDir() || !stat.Mode().IsRegular() {
			// we never serve irregular files
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// let the stdlib serve the file.
		rdr, err := os.Open(file)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		defer rdr.Close()
		http.ServeContent(w, r, file, stat.ModTime(), rdr)
	}
}
