// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package store

import (
	"github.com/playbymail/empyr/internal/domains"
	"github.com/playbymail/empyr/store/sqlc"
	"log"
)

// this file implements functions for the Authentication and Session services.

func (s *Store) GetMagicKeyUser(handle, key string) (domains.User, error) {
	log.Printf("store: get magic key: %q %q\n", handle, key)
	row, err := s.Queries.ReadUserByMagicKey(s.Context, sqlc.ReadUserByMagicKeyParams{
		Handle:   handle,
		MagicKey: key,
	})
	if err != nil {
		log.Printf("store: get magic key: %q %q: %v\n", handle, key, err)
		return domains.User{}, ErrNotFound
	}
	return domains.User{
		ID:      domains.UserID(row.ID),
		Handle:  handle,
		IsUser:  row.IsActive == 1,
		IsAdmin: row.IsActive == 1 && row.IsAdmin == 1,
	}, nil
}

func (s *Store) GetSession(sessionID domains.SessionID) (domains.Session, error) {
	return domains.Session{}, ErrNotImplemented
}

func (s *Store) GetUser(userID domains.UserID) (domains.User, error) {
	//log.Printf("store: get user: %d\n", userID)
	if userID == 0 {
		return domains.User{}, ErrNotFound
	}
	row, err := s.Queries.ReadUserByID(s.Context, int64(userID))
	if err != nil {
		log.Printf("store: get user: %d: %v\n", userID, err)
		return domains.User{}, err
	}
	return domains.User{
		ID:      domains.UserID(row.ID),
		Handle:  row.Handle,
		IsUser:  row.IsActive == 1,
		IsAdmin: row.IsActive == 1 && row.IsAdmin == 1,
	}, nil
}
