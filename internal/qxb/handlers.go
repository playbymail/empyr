// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package qxb

import (
	"fmt"
	"github.com/playbymail/empyr/internal/services/smgr"
	"github.com/playbymail/empyr/store"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// really this is just middleware

func NewQXB(assets string, mgr *smgr.SessionManager, db *store.Store) (*QXB, error) {
	log.Printf("qxb: assets %q\n", assets)
	// we must have an absolute path to detect directory traversal attacks
	if path, err := filepath.Abs(assets); err != nil {
		log.Printf("error: %q: %v\n", assets, err)
		panic(fmt.Sprintf("%s: invalid directory: %v", assets, err))
	} else {
		assets = path
	}
	// the assets path must exist and be a directory
	if sb, err := os.Stat(assets); err != nil {
		log.Printf("qxb: assets: %q: %v\n", assets, err)
		return nil, err
	} else if !sb.IsDir() {
		log.Printf("qxb: assets: %q: not a directory\n", assets)
		return nil, fmt.Errorf("%s: not a directory", assets)
	}

	return &QXB{
		assets: assets,
		mgr:    mgr,
		db:     db,
	}, nil
}

type QXB struct {
	assets string
	mgr    *smgr.SessionManager
	db     *store.Store
}

var (
	// create a regular expression that will match any hidden file in any
	// directory in the path. a hidden file starts with a dot.
	hiddenFiles = regexp.MustCompile(`/\.`)
)

func (q *QXB) AdminOnly(next http.Handler) http.Handler {
	log.Printf("qxb: admin: registered as middleware\n")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
		if session := smgr.GetSession(r); session == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else if username, ok := session.Get("username").(string); !ok || username == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else if row, err := q.db.Queries.ReadUserByUsername(q.db.Context, username); err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else if row.ID == 0 || row.IsActive != 1 {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		} else if row.IsAdmin != 1 {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (q *QXB) MustAuth(next http.Handler) http.Handler {
	log.Printf("qxb: auth: registered as middleware\n")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
		if session := smgr.GetSession(r); session == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else if username, ok := session.Get("username").(string); !ok || username == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else if row, err := q.db.Queries.ReadUserByUsername(q.db.Context, username); err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		} else if row.ID == 0 || row.IsActive != 1 {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (q *QXB) Assets(next http.Handler) http.Handler {
	log.Printf("qxb: assets: registered as middleware\n")
	next = q.mgr.Sessions(next)

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
		fullPath := filepath.Join(q.assets, path)
		// log.Printf("%s %s: full %q\n", r.Method, r.URL, fullPath)

		// prevent directory traversal attacks. (note: the filepath.Clean
		// function already does this, but we do it again to be safe.)
		if !strings.HasPrefix(fullPath, q.assets) {
			// log.Printf("%s %s: path is outside of assets directory\n", r.Method, r.URL)
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
