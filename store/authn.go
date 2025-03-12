// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package store

import (
	"github.com/playbymail/empyr/internal/domains"
	"log"
)

// this file implements functions for the Authentication and Session services.

func (s *Store) GetMagicKeyUser(key string) (domains.User, error) {
	log.Printf("store: get magic key user: %s\n", key)
	row, err := s.Queries.ReadUserByMagicKey(s.Context, key)
	if err != nil {
		log.Printf("store: get magic key user: %v\n", err)
		return domains.User{}, ErrNotFound
	}
	return domains.User{
		ID:     domains.UserID(row.ID),
		Handle: row.Handle,
		IsUser: true,
	}, nil
}

func (s *Store) GetSession(sessionID domains.SessionID) (domains.Session, error) {
	return domains.Session{}, ErrNotImplemented
}

func (s *Store) GetUser(userID domains.UserID) (domains.User, error) {
	return domains.User{}, ErrNotImplemented
}
