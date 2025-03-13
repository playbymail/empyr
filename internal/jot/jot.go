// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package jot implements a simple JSON Web Token (JWT) that can only be used
// to authenticate a user session. It should not be used in any production
// environment.
package jot

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/mdhender/semver"
	"github.com/playbymail/empyr/internal/domains"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	cookieName = "empyr-jot"
	cookiePath = "/"
	cookieTTL  = 4 * 7 * 24 * time.Hour
)

var (
	version = semver.Version{Major: 0, Minor: 0, Patch: 1}
)

func Version() string {
	return version.String()
}

type Factory struct {
	key []byte
}

// NewFactory returns a Factory that can create and verify JOT tokens.
func NewFactory(secret string) *Factory {
	// hashedKey is the SHA256 hash of the secret.
	hashedKey := sha256.Sum256([]byte(secret))
	// make a copy of the hashed key for use by the Factory.
	keyCopy := make([]byte, len(hashedKey))
	copy(keyCopy, hashedKey[:])
	return &Factory{key: keyCopy}
}

// DeleteCookie deletes the session cookie.
func (f *Factory) DeleteCookie(w http.ResponseWriter) {
	deleteCookie(w)
}

// SetCookie creates a session cookie.
func (f *Factory) SetCookie(w http.ResponseWriter, t *Token) {
	setCookie(w, t)
}

// Version returns the version of the JOT library.
func (f *Factory) Version() string {
	return version.String()
}

// GetToken returns the Token from the request. It checks the request for
// either the bearer token or the cookie. If neither is found, it returns nil.
func (f *Factory) GetToken(r *http.Request) *Token {
	text, ok := f.fromBearerToken(r)
	if !ok {
		if text, ok = f.fromCookie(r); !ok {
			return nil
		}
	}
	// token should be in the form: id.expiresAt.signature
	fields := strings.SplitN(text, ".", 3)
	if len(fields) != 3 {
		return nil
	}
	// log.Printf("jot: getToken: fields: %+v\n", fields)

	var t Token
	if id, err := strconv.Atoi(fields[0]); err != nil {
		// log.Printf("jot: getToken: id: strconv: %v\n", err)
		return nil
	} else {
		t.id = domains.UserID(id)
		// log.Printf("jot: getToken: id: %d\n", t.id)
	}
	if expiresAt, err := strconv.Atoi(fields[1]); err != nil {
		// log.Printf("jot: getToken: expiresAt: strconv: %v\n", err)
		return nil
	} else {
		t.expiresAt = time.Unix(int64(expiresAt), 0)
		// log.Printf("jot: getToken: expiresAt: %s\n", t.expiresAt.Format(time.RFC3339))
	}
	if t.signature = fields[2]; t.signature == "" {
		// log.Printf("jot: getToken: signature: empty\n")
		return nil
	}

	// verify that the token is valid. we check for expiration first because
	// it is cheaper to check than to verify the signature.
	if t.expiresAt.IsZero() {
		// log.Printf("jot: getToken: token has no expiration timestamp\n")
	} else if now := time.Now().UTC(); !now.Before(t.expiresAt) {
		// log.Printf("jot: getToken: token has expired\n")
	} else if t.isValid = f.verify(t.id, t.expiresAt, t.signature); !t.isValid {
		// log.Printf("jot: getToken: token not signed\n")
	} else {
		// log.Printf("jot: getToken: token is valid\n")
	}

	return &t
}

func (f *Factory) NewToken(id domains.UserID) *Token {
	expiresAt := time.Now().Add(cookieTTL).UTC()
	return &Token{
		id:        id,
		expiresAt: expiresAt,
		signature: f.sign(id, expiresAt),
	}
}

// fromBearerToken returns the text of the token from the Authorization header.
// If there is no bearer token, it returns false.
func (f *Factory) fromBearerToken(r *http.Request) (string, bool) {
	// log.Printf("jot: bearer: entered\n")
	headerAuthText := r.Header.Get("Authorization")
	if headerAuthText == "" {
		return "", false
	}
	// log.Printf("jon: bearer: found authorization header\n")
	authTokens := strings.SplitN(headerAuthText, " ", 2)
	if len(authTokens) != 2 {
		return "", false
	}
	// log.Printf("jot: bearer: found authorization token\n")
	authType, authToken := authTokens[0], strings.TrimSpace(authTokens[1])
	if authType != "Bearer" {
		return "", false
	}
	// log.Printf("jot: bearer: found bearer token\n")
	return authToken, true
}

// fromCookie returns the text of the token from a session cookie.
// If there is no session cookie, it returns false.
func (f *Factory) fromCookie(r *http.Request) (string, bool) {
	// log.Printf("jot: cookie: entered\n")
	c, err := r.Cookie(cookieName)
	if err != nil {
		// log.Printf("jot: cookie: %+v\n", err)
		return "", false
	} else if c.Value == "" {
		// log.Printf("jot: cookie: missing value\n")
		return "", false
	}
	// log.Printf("jot: cookie: %q\n", c.Value)
	return c.Value, true
}

// sign returns the hex-encoded signature of the message using HMAC-SHA256 and the factory key.
func (f *Factory) sign(id domains.UserID, expiresAt time.Time) string {
	msg := fmt.Sprintf("%d.%d", id, expiresAt.Unix())
	// sign the message using HMAC-SHA256 and the factory key
	hm := hmac.New(sha256.New, f.key)
	hm.Write([]byte(msg))
	// return the hex-encoded signature
	return fmt.Sprintf("%x", hm.Sum(nil))
}

// verify returns true if the signature is valid for the message.
func (f *Factory) verify(id domains.UserID, expiresAt time.Time, signature string) bool {
	return signature == f.sign(id, expiresAt)
}

// Token implements my version of JOT.
type Token struct {
	id        domains.UserID
	expiresAt time.Time
	signature string
	isValid   bool
}

// IsValid returns true only if the Token is signed, active, and not expired.
func (t *Token) IsValid() bool {
	return t != nil && t.isValid
}

// SetCookie creates session cookie.
func (t *Token) SetCookie(w http.ResponseWriter) {
	setCookie(w, t)
}

// String implements the Stringer interface.
// Please don't call this before signing the token.
func (t *Token) String() string {
	return fmt.Sprintf("%d.%d.%s", t.id, t.expiresAt.Unix(), t.signature)
}

// UserID returns the user ID associated with the Token.
// If the Token is not valid, this function returns 0.
func (t *Token) UserID() domains.UserID {
	if !t.IsValid() {
		return 0
	}
	return t.id
}

func (t *Token) secondsUntilExpiration() int {
	timeLeft := int(t.expiresAt.Sub(time.Now().UTC()).Seconds())
	if timeLeft <= 0 {
		return -1
	}
	return timeLeft
}

// deleteCookie is a helper function to delete the session cookie.
func deleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Path:     cookiePath,
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// setCookie is a helper function to create a session cookie.
func setCookie(w http.ResponseWriter, t *Token) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Path:     cookiePath,
		Value:    t.String(),
		MaxAge:   t.secondsUntilExpiration(),
		HttpOnly: true,
		Secure:   true,
	})
}
