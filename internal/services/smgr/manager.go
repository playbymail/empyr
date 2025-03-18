// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package smgr implements a session manager.
package smgr

import (
	"context"
	"log"
	"net/http"
	"time"
)

func NewSessionManager(
	store SessionStore,
	gcInterval,
	idleExpiration,
	absoluteExpiration time.Duration,
	cookieName string) *SessionManager {

	m := &SessionManager{
		store:              store,
		idleExpiration:     idleExpiration,
		absoluteExpiration: absoluteExpiration,
		cookieName:         cookieName,
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
}

func (m *SessionManager) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start the session
		session, rws := m.start(r)

		// create a new response writer
		sw := &sessionResponseWriter{
			ResponseWriter: w,
			sessionManager: m,
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
			if !m.verifyCSRFToken(rws, session) {
				http.Error(sw, "CSRF token mismatch", http.StatusForbidden)
				return
			}
		}

		// call the next handler and pass the new response writer and new request
		next.ServeHTTP(sw, rws)

		// save the session
		m.save(session)

		// write the session cookie to the response if not already written
		writeCookieIfNecessary(sw)
	})
}

// gc calls the store to perform garbage collection (reap expired sessions).
func (m *SessionManager) gc(d time.Duration) {
	ticker := time.NewTicker(d)
	for range ticker.C {
		_ = m.store.gc(m.idleExpiration, m.absoluteExpiration)
	}
}

// migrate deletes an existing session and creates a new one with a fresh ID
func (m *SessionManager) migrate(session *Session) error {
	session.mu.Lock()
	defer session.mu.Unlock()
	err := m.store.destroy(session.id)
	if err != nil {
		return err
	}
	session.id = generateSessionId()
	return nil
}

// save updates the session's lastActivityAt field and writes it to the store.
func (m *SessionManager) save(session *Session) error {
	session.lastActivityAt = time.Now()
	err := m.store.write(session)
	if err != nil {
		log.Printf("error: session: save: failed to write session to store: %v", err)
		return err
	}
	return nil
}

type sessionContextKey struct{}

// start retrieves the session by reading the session cookie or generates a new one if needed.
// It then attaches the session to the request using context values.
func (m *SessionManager) start(r *http.Request) (*Session, *http.Request) {
	var session *Session

	// read session from cookie
	cookie, err := r.Cookie(m.cookieName)
	if err == nil {
		session, err = m.store.read(cookie.Value)
		if err != nil {
			log.Printf("error: session: cookie: failed to read session from store: %v", err)
		}
	}

	// generate a new session
	if session == nil || !m.validate(session) {
		session = newSession()
	}

	// attach session to context
	ctx := context.WithValue(r.Context(), sessionContextKey{}, session)
	r = r.WithContext(ctx)

	return session, r
}

// validate ensures the session is valid for use.
// It is invalid if it is expired or has been idle for too long.
func (m *SessionManager) validate(session *Session) bool {
	if time.Since(session.createdAt) > m.absoluteExpiration ||
		time.Since(session.lastActivityAt) > m.idleExpiration {
		// delete the session from the store
		err := m.store.destroy(session.id)
		if err != nil {
			log.Printf("error: %v\n", err)
		}
		return false
	}

	return true
}

// verifyCSRFToken extracts the CSRF token from a given session and validates it
// against the csrf_token form value or the X-CSRF-Token header in the request.
func (m *SessionManager) verifyCSRFToken(r *http.Request, session *Session) bool {
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
