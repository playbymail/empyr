// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package store

import (
	"github.com/playbymail/empyr/internal/cerr"
	"github.com/playbymail/empyr/internal/domains"
	"github.com/playbymail/empyr/store/sqlc"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

// this file implements functions for the Authentication and Session services.

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
		ID:       domains.UserID(row.ID),
		Username: row.Username,
		Email:    row.Email,
		IsUser:   row.IsActive == 1,
		IsAdmin:  row.IsActive == 1 && row.IsAdmin == 1,
	}, nil
}

const (
	ErrInvalidUsername = cerr.Error("invalid username")
	ErrInvalidPassword = cerr.Error("invalid password")
	ErrInvalidEmail    = cerr.Error("invalid email")
)

func (s *Store) RegisterAdminUser(username, email, password string) (domains.UserID, error) {
	if username != strings.ToLower(username) || username != strings.TrimSpace(username) {
		return 0, ErrInvalidUsername
	} else if email != strings.ToLower(email) || email != strings.TrimSpace(email) {
		return 0, ErrInvalidEmail
	} else if len(password) < 8 {
		return 0, ErrInvalidPassword
	}
	hashedPassword, err := HashPassword(password)
	if err != nil {
		log.Printf("store: register admin: %q: %v\n", username, err)
		return 0, err
	}
	userID, err := s.Queries.CreateUser(s.Context, sqlc.CreateUserParams{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
		IsActive:       1,
		IsAdmin:        1,
	})
	if err != nil {
		log.Printf("store: register admin: %q: %v\n", username, err)
		return 0, err
	}
	return domains.UserID(userID), nil
}

func (s *Store) RegisterUser(username, email, password string) (domains.UserID, error) {
	if username != strings.ToLower(username) || username != strings.TrimSpace(username) {
		return 0, ErrInvalidUsername
	} else if email != strings.ToLower(email) || email != strings.TrimSpace(email) {
		return 0, ErrInvalidEmail
	} else if len(password) < 8 {
		return 0, ErrInvalidPassword
	}
	hashedPassword, err := HashPassword(password)
	if err != nil {
		log.Printf("store: register user: %q: %v\n", username, err)
		return 0, err
	}
	userID, err := s.Queries.CreateUser(s.Context, sqlc.CreateUserParams{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
		IsActive:       1,
		IsAdmin:        0,
	})
	if err != nil {
		log.Printf("store: register user: %q: %v\n", username, err)
		return 0, err
	}
	return domains.UserID(userID), nil
}

func (s *Store) ValidateCredentials(email, password string) (domains.User, error) {
	log.Printf("store: validate user: %q\n", email)
	row, err := s.Queries.ReadUserByEmail(s.Context, email)
	if err != nil {
		log.Printf("store: validate user: %q: %v\n", email, err)
		return domains.User{}, ErrInvalidEmail
	}
	if !CheckPassword(password, row.HashedPassword) {
		return domains.User{}, ErrInvalidPassword
	}
	return domains.User{
		ID:       domains.UserID(row.ID),
		Username: row.Username,
		Email:    row.Email,
		IsUser:   row.IsActive == 1,
		IsAdmin:  row.IsActive == 1 && row.IsAdmin == 1,
	}, nil
}

// simple password functions inspired by https://www.gregorygaines.com/blog/how-to-properly-hash-and-salt-passwords-in-golang-bcrypt/

// CheckPassword returns true if the plain text password matches the hashed password.
func CheckPassword(plainTextPassword, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword)) == nil
}

// HashPassword uses the cheapest bcrypt cost to hash the password because we are not going to use
// it for anything other than authentication in non-production environments.
func HashPassword(plainTextPassword string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashedPasswordBytes), err
}
