// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package smgr

import (
	"net/http"
	"time"
)

type sessionResponseWriter struct {
	http.ResponseWriter
	sessionManager *SessionManager
	request        *http.Request
	done           bool
}

func (w *sessionResponseWriter) Write(b []byte) (int, error) {
	writeCookieIfNecessary(w)

	return w.ResponseWriter.Write(b)
}

func (w *sessionResponseWriter) WriteHeader(code int) {
	writeCookieIfNecessary(w)

	w.ResponseWriter.WriteHeader(code)
}

func (w *sessionResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func writeCookieIfNecessary(w *sessionResponseWriter) {
	if w.done {
		return
	}
	session, ok := w.request.Context().Value(sessionContextKey{}).(*Session)
	if !ok {
		// todo: panic can happen after database commits. prolly not great.
		panic("session not found in request context")
	}
	cookie := &http.Cookie{
		Name:  w.sessionManager.cookieName,
		Value: session.id,
		// restricts the cookie to this specific domain, preventing it
		// from being sent to unauthorized sites.
		Domain: "mywebsite.com",
		// protects the cookie from being accessed by JavaScript, reducing
		// the risk of cross-site scripting (XSS) attacks.
		HttpOnly: true,
		// allows the cookie to be sent with all requests within the domain,
		// making it accessible across different parts of the website.
		Path: "/",
		// ensures the cookie is only transmitted over HTTPS, preventing
		// man-in-the-middle (MITM) attacks.
		Secure: true,
		// prevents the cookie from being sent with most cross-site requests,
		// reducing cross-site request forgery (CSRF) risks while still allowing
		// some navigation-based requests (e.g., clicking a link from another site).
		SameSite: http.SameSiteLaxMode,
		// sets the expiration date based on the idle session timeout.
		Expires: time.Now().Add(w.sessionManager.idleExpiration),
		// sets the expiration date based on the idle session timeout.
		MaxAge: int(w.sessionManager.idleExpiration / time.Second),
	}
	http.SetCookie(w.ResponseWriter, cookie)
	w.done = true
}
