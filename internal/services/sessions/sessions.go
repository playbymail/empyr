// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package sessions implements a session management service.
package sessions

import (
	"context"
	"github.com/playbymail/empyr/internal/cerr"
	"github.com/playbymail/empyr/internal/domains"
	"github.com/playbymail/empyr/internal/jot"
	"net/http"
)

const (
	// ErrNotFound is used when a session could not be found.
	ErrNotFound = cerr.Error("session not found")
)

// Repository defines the operations required to fetch and store authentication data.
type Repository interface {
	GetSession(id domains.SessionID) (domains.Session, error)
	GetUser(domains.UserID) (domains.User, error)
}

// Service defines the operations this service will perform.
// It's the contract with the outside world and is defined so we can mock it for testing.
type Service interface {
	CreateSession(w http.ResponseWriter, userID domains.UserID) (SessionID, error)
	DeleteSession(w http.ResponseWriter, sessionID SessionID) error
	GetSession(r *http.Request) (SessionID, error)
	GetUser(r *http.Request) domains.User
}

// service defines the service we are implementing
type service struct {
	r  Repository
	jf *jot.Factory
}

func NewService(r Repository, jf *jot.Factory) Service {
	return &service{r: r, jf: jf}
}

func (s *service) CreateSession(w http.ResponseWriter, userID domains.UserID) (SessionID, error) {
	t := s.jf.NewToken(userID)
	t.SetCookie(w)
	return SessionID(t.UserID()), nil
}

func (s *service) DeleteSession(w http.ResponseWriter, id SessionID) error {
	s.jf.DeleteCookie(w)
	return nil
}

func (s *service) GetSession(r *http.Request) (SessionID, error) {
	t := s.jf.GetToken(r)
	if !t.IsValid() {
		return SessionID(0), ErrNotFound
	}
	return SessionID(t.UserID()), nil
}

func (s *service) GetUser(r *http.Request) domains.User {
	// log.Printf("%s %s: mw: sessions: GetUser entered\n", r.Method, r.URL.Path)
	// try from the context first
	user, ok := r.Context().Value(sessionsContextKey("user")).(domains.User)
	if ok {
		return user
	}
	// then from the request
	t := s.jf.GetToken(r)
	if !t.IsValid() {
		return domains.User{}
	}
	user, err := s.r.GetUser(t.UserID())
	if err != nil {
		return domains.User{}
	}
	return user
}

// sessionsContextKey is a custom type to prevent key collisions.
type sessionsContextKey string

func AddUserToContext(h http.Handler, sessionsService Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("%s %s: mw: sessions: entered\n", r.Method, r.URL.Path)
		user := sessionsService.GetUser(r)
		ctx := context.WithValue(r.Context(), sessionsContextKey("user"), user)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
