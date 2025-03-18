// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package smgr

import (
	"github.com/google/uuid"
	"net/http"
	"sync"
	"time"
)

type Session struct {
	createdAt      time.Time
	lastActivityAt time.Time
	id             string
	mu             sync.Mutex
	data           map[string]any
}

// generateSessionId generates a random session ID with at least 64 bits of entropy.
func generateSessionId() string {
	return uuid.New().String()
}

// newSession creates and initializes a new session.
func newSession() *Session {
	return &Session{
		id:             generateSessionId(),
		data:           map[string]any{"csrf_token": generateCSRFToken()},
		createdAt:      time.Now(),
		lastActivityAt: time.Now(),
	}
}

func (s *Session) Get(key string) any {
	return s.data[key]
}

func (s *Session) Put(key string, value any) {
	s.data[key] = value
}

func (s *Session) Delete(key string) {
	delete(s.data, key)
}

// GetSession retrieves the session from the request context.
func GetSession(r *http.Request) *Session {
	session, ok := r.Context().Value(sessionContextKey{}).(*Session)
	if !ok {
		// todo: panic is not great.
		panic("session not found in request context")
	}
	return session
}
