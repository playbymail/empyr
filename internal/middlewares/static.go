// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package middlewares

import (
	"fmt"
	"github.com/playbymail/empyr/internal/router"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Static middleware serves static files from the filesystem.
func Static(root string) router.Middleware {
	log.Printf("static: %q: registered as middleware\n", root)

	// ensure the staticDir has an absolute path so that we can detect directory traversal attacks
	if path, err := filepath.Abs(root); err != nil {
		log.Printf("error: static: %q: %v\n", root, err)
		panic("static: invalid directory: " + err.Error())
	} else {
		root = path
	}

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
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			})
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s: static: %q\n", r.Method, r.URL, root)
			// clean the request path
			path := filepath.Clean(r.URL.Path)
			log.Printf("%s %s: static: clean %q\n", r.Method, r.URL, path)
			//if path == "/" {
			//	// let the next handler handle the request
			//	log.Printf("%s %s: static: / is not an asset\n", r.Method, r.URL)
			//	next.ServeHTTP(w, r)
			//	return
			//} else if strings.HasSuffix(path, "/index.html") {
			//	// prevent stdlib from permanently redirecting index.html to /
			//	log.Printf("%s %s: static: caught /index.html\n", r.Method, r.URL)
			//	next.ServeHTTP(w, r)
			//	return
			//}

			// build the full file path
			filePath := filepath.Join(root, path)

			// prevent directory traversal attacks
			if !strings.HasPrefix(filePath, root) {
				log.Printf("%s %s: static: path is outside of root directory\n", r.Method, r.URL)
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			// check if the file exists and is a regular file
			sb, err := os.Stat(filePath)
			if err != nil && os.IsNotExist(err) {
				// the file does not exist, call the next handler
				log.Printf("%s %s: static: not an asset\n", r.Method, r.URL)
				next.ServeHTTP(w, r)
				return
			} else if err != nil {
				// return an internal server error on any other error
				log.Printf("%s %s: static: stat: %v\n", r.Method, r.URL, err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			} else if sb.IsDir() {
				// file exists and is a directory
				log.Printf("%s %s: static: path is directory\n", r.Method, r.URL)
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			} else if !sb.Mode().IsRegular() {
				// file exists and is not a regular file
				log.Printf("%s %s: static: path is special file\n", r.Method, r.URL)
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			// file exists and is a regular file, so serve it
			log.Printf("%s %s: static: serving asset\n", r.Method, r.URL)

			// set caching headers
			modTime := sb.ModTime()
			etag := fmt.Sprintf("\"%x-%x\"", sb.Size(), modTime.Unix())
			w.Header().Set("ETag", etag)

			// handle conditional requests
			if match := r.Header.Get("If-None-Match"); match == etag {
				w.WriteHeader(http.StatusNotModified)
				return
			} else if since := r.Header.Get("If-Modified-Since"); since != "" {
				if t, err := time.Parse(http.TimeFormat, since); err == nil && modTime.Before(t.Add(1*time.Second)) {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}

			// open the file for reading so we can serve it
			file, err := os.Open(filePath)
			if err != nil {
				log.Printf("%s %s: static: open: %v\n", r.Method, r.URL, err)
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			defer file.Close()

			// serve the file
			http.ServeContent(w, r, path, modTime, file)
		})
	}
}
