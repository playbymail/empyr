// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package middlewares

import (
	"github.com/playbymail/empyr/internal/router"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Static middleware serves static files from the filesystem.
func Static(root string) router.Middleware {
	log.Printf("static: %q: registered as middleware\n", root)
	notFound := false
	if sb, err := os.Stat(root); err != nil {
		log.Printf("static: %q: %v\n", root, err)
		notFound = true
	} else if !sb.IsDir() {
		log.Printf("static: %q: not a directory\n", root)
		notFound = true
	}
	if notFound {
		return func(_ http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			})
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s: static: %q\n", r.Method, r.URL, root)
			if path := filepath.Clean(r.URL.Path); path != "/" {
				path = filepath.Join(root, path)
				if sb, err := os.Stat(path); err == nil {
					if sb.IsDir() {
						// never serve directories or other non-regular files
						log.Printf("%s %s: static: path is directory\n", r.Method, r.URL)
						http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
						return
					} else if !sb.Mode().IsRegular() {
						// never serve directories or other non-regular files
						log.Printf("%s %s: static: path is special file\n", r.Method, r.URL)
						http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
						return
					}
					log.Printf("%s %s: static: %s\n", r.Method, r.URL, path)
					http.ServeFile(w, r, path)
					return
				}
			}
			// path is not an asset, so pass through to the next handler
			log.Printf("%s %s: static: %q: path is not an asset\n", r.Method, r.URL, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}
