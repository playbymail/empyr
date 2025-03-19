// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package qxb

import (
	"fmt"
	"github.com/playbymail/empyr/internal/services/smgr"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// really this is just middleware

func NewQXB(root string, mgr *smgr.SessionManager) (*QXB, error) {
	log.Printf("qxb: root %q\n", root)
	// we must have an absolute path to detect directory traversal attacks
	if path, err := filepath.Abs(root); err != nil {
		log.Printf("error: %q: %v\n", root, err)
		panic(fmt.Sprintf("%s: invalid directory: %v", root, err))
	} else {
		root = path
	}
	// the root must exist and be a directory
	if sb, err := os.Stat(root); err != nil {
		log.Printf("qxb: %q: %v\n", root, err)
		return nil, err
	} else if !sb.IsDir() {
		log.Printf("qxb: %q: not a directory\n", root)
		return nil, fmt.Errorf("%s: not a directory", root)
	}

	return &QXB{
		root: root,
		mgr:  mgr,
	}, nil
}

type QXB struct {
	root string
	mgr  *smgr.SessionManager
}

var (
	// create a regular expression that will match any hidden file in any
	// directory in the path. a hidden file starts with a dot.
	hiddenFiles = regexp.MustCompile(`/\.`)
)

func (q *QXB) Handle(next http.Handler) http.Handler {
	log.Printf("qxb: registered as middleware\n")
	next = q.mgr.Handle(next)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s: entered\n", r.Method, r.URL.Path)

		// clean the request path
		path := filepath.Clean(r.URL.Path)
		if path == "." || path == "/" {
			// assume that we're being asked to serve the landing page.
			// let the next handler handle the request
			next.ServeHTTP(w, r)
			return
		}
		// log.Printf("%s %s: clean %q\n", r.Method, r.URL, path)

		// search for hidden files after cleaning the path and before building the full file path.
		if hiddenFiles.MatchString(path) {
			// log.Printf("%s %s: path contains hidden file or directory\n", r.Method, r.URL)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		// build the full file path
		fullPath := filepath.Join(q.root, path)
		// log.Printf("%s %s: full %q\n", r.Method, r.URL, fullPath)

		// prevent directory traversal attacks. (note: the filepath.Clean
		// function already does this, but we do it again to be safe.)
		if !strings.HasPrefix(fullPath, q.root) {
			// log.Printf("%s %s: path is outside of root directory\n", r.Method, r.URL)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		// check if the file exists and is a regular file
		sb, err := os.Stat(fullPath)
		if err != nil {
			// assume that any error with the stat means that
			// we have an API route, not an asset to serve, so
			// let the next handler handle the request
			next.ServeHTTP(w, r)
			return
		} else if sb.IsDir() { // never serve directories
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if !sb.Mode().IsRegular() { // never serve special files
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		// file exists and is a regular file, so serve it
		// log.Printf("%s %s: serving asset\n", r.Method, r.URL)

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
		file, err := os.Open(fullPath)
		if err != nil {
			// log.Printf("%s %s: static: open: %v\n", r.Method, r.URL, err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		defer file.Close()

		// serve the file
		http.ServeContent(w, r, path, modTime, file)
	})
}
