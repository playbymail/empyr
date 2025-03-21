// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package smgr implements a session manager.
package smgr

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func NewSessionManager(
	store SessionStore,
	gcInterval,
	idleExpiration,
	absoluteExpiration time.Duration,
	domain string,
	cookieName string) *SessionManager {

	m := &SessionManager{
		store:              store,
		idleExpiration:     idleExpiration,
		absoluteExpiration: absoluteExpiration,
		domain:             domain,
		cookieName:         cookieName,
		anonymous:          "d0aba137-d7a7-4e25-86f9-ab6094b33a46",
	}

	go m.gc(gcInterval)

	return m
}

// SessionManager coordinates all session operations.
type SessionManager struct {
	store              SessionStore
	idleExpiration     time.Duration
	absoluteExpiration time.Duration
	cookieName         string
	domain             string
	anonymous          string
}

func (sm *SessionManager) Sessions(next http.Handler) http.Handler {
	log.Printf("smgr: registered as middleware\n")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s: entered\n", r.Method, r.URL.Path)

		// start the session
		session, rws := sm.start(r)

		// create a new response writer
		sw := &sessionResponseWriter{
			ResponseWriter: w,
			sessionManager: sm,
			request:        rws,
		}

		// add essential headers to help maintain proper session handling and prevent
		// caching issues that could lead to incorrect user data being displayed.
		//
		// the Vary: Cookie header ensures that caches, such as CDN or browser caches,
		// differentiate responses based on the presence or value of the Cookie header.
		// this prevents serving cached responses intended for one user to another.
		w.Header().Add("Vary", "Cookie")
		// the Cache-Control: no-cache="Set-Cookie" directive instructs caches not to store
		// responses that include the Set-Cookie header, ensuring that session-related
		// headers are always processed fresh from the server rather than being retrieved
		// from cache.
		w.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)

		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch || r.Method == http.MethodDelete {
			if !sm.verifyCSRFToken(rws, session) {
				http.Error(sw, "CSRF token mismatch", http.StatusForbidden)
				return
			}
		}

		// call the next handler and pass the new response writer and new request
		next.ServeHTTP(sw, rws)

		// save the session
		sm.save(session)

		// write the session cookie to the response if not already written
		writeCookieIfNecessary(sw)
	})
}

func (sm *SessionManager) Login(r *http.Request, username string) error {
	session := GetSession(r)
	err := sm.migrate(session)
	if err != nil {
		return fmt.Errorf("failed to migrate session: %w", err)
	}
	session.Put("username", username)
	return nil
}

func (sm *SessionManager) Logout(r *http.Request) error {
	session := GetSession(r)
	err := sm.migrate(session)
	if err != nil {
		return fmt.Errorf("failed to migrate session: %w", err)
	}
	session.Put("username", "")
	return nil
}

// gc calls the store to perform garbage collection (reap expired sessions).
func (sm *SessionManager) gc(d time.Duration) {
	ticker := time.NewTicker(d)
	for range ticker.C {
		_ = sm.store.gc(sm.idleExpiration, sm.absoluteExpiration)
	}
}

// migrate deletes an existing session and creates a new one with a fresh ID
func (sm *SessionManager) migrate(session *Session) error {
	session.mu.Lock()
	defer session.mu.Unlock()
	err := sm.store.destroy(session.id)
	if err != nil {
		return err
	}
	session.id = generateSessionId()
	return nil
}

// save updates the session's lastActivityAt field and writes it to the store.
func (sm *SessionManager) save(session *Session) error {
	session.lastActivityAt = time.Now()
	err := sm.store.write(session)
	if err != nil {
		log.Printf("error: session: save: failed to write session to store: %v", err)
		return err
	}
	return nil
}

type sessionContextKey struct{}

// start retrieves the session by reading the session cookie or generates a new one if needed.
// It then attaches the session to the request using context values.
func (sm *SessionManager) start(r *http.Request) (*Session, *http.Request) {
	var session *Session

	// read session from cookie
	log.Printf("smgr: start: reading session from cookie %q\n", sm.cookieName)
	cookie, err := r.Cookie(sm.cookieName)
	if err == nil {
		log.Printf("smgr: start: cookie: session from cookie %q\n", sm.cookieName)
		log.Printf("smgr: start: cookie: session from cookie %q\n", cookie.Value)
		session, err = sm.store.read(cookie.Value)
		if err != nil {
			log.Printf("error: session: cookie: failed to read session from store: %v", err)
		}
	} else {
		log.Printf("smgr: start: cookie: no session from cookie %q\n", sm.cookieName)
	}

	// generate a new session
	if session == nil {
		log.Printf("smgr: start: generating new session: session == nil\n")
	} else if !sm.validate(session) {
		log.Printf("smgr: start: generating new session: session is invalid\n")
	}
	if session == nil || !sm.validate(session) {
		session = newSession()
	}

	// attach session to context
	ctx := context.WithValue(r.Context(), sessionContextKey{}, session)
	r = r.WithContext(ctx)

	return session, r
}

// validate ensures the session is valid for use.
// It is invalid if it is expired or has been idle for too long.
func (sm *SessionManager) validate(session *Session) bool {
	if time.Since(session.createdAt) > sm.absoluteExpiration ||
		time.Since(session.lastActivityAt) > sm.idleExpiration {
		// delete the session from the store
		err := sm.store.destroy(session.id)
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		return false
	}

	return true
}

// verifyCSRFToken extracts the CSRF token from a given session and validates it
// against the csrf_token form value or the X-CSRF-Token header in the request.
func (sm *SessionManager) verifyCSRFToken(r *http.Request, session *Session) bool {
	sToken, ok := session.Get("csrf_token").(string)
	if !ok {
		return false
	}
	token := r.FormValue("csrf_token")
	if token == "" {
		token = r.Header.Get("X-XSRF-Token")
	}
	return token == sToken
}
