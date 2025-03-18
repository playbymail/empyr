// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package smgr

import (
	"sync"
	"time"
)

// SessionStore defines the interface for a storage engine for sessions.
type SessionStore interface {
	// read reads the session that has the given unique ID.
	read(id string) (*Session, error)
	// write writes a session to the storage engine.
	write(session *Session) error
	// destroy deletes a session with the given ID.
	destroy(id string) error
	// gc performs garbage collection.
	// It queries all expired sessions and deletes them from storage.
	gc(idleExpiration, absoluteExpiration time.Duration) error
}

// implement an in-memory session store for testing purposes.

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]*Session),
	}
}

type InMemorySessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func (s *InMemorySessionStore) read(id string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, _ := s.sessions[id]

	return session, nil
}

func (s *InMemorySessionStore) write(session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[session.id] = session

	return nil
}

func (s *InMemorySessionStore) destroy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, id)

	return nil
}

func (s *InMemorySessionStore) gc(idleExpiration, absoluteExpiration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, session := range s.sessions {
		if time.Since(session.lastActivityAt) > idleExpiration ||
			time.Since(session.createdAt) > absoluteExpiration {
			delete(s.sessions, id)
		}
	}
	return nil
}
