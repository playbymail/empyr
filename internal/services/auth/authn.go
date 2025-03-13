// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package auth implements an authentication service.
package auth

import (
	"github.com/playbymail/empyr/internal/cerr"
	"github.com/playbymail/empyr/internal/domains"
	"log"
)

const (
	// ErrNotFound is used when a beer could not be found.
	ErrNotFound = cerr.Error("beer not found")
)

// Repository defines the operations required to fetch and store authentication data.
type Repository interface {
	GetMagicKeyUser(handle, key string) (domains.User, error)
	GetUser(domains.UserID) (domains.User, error)
}

// Service defines the operations this service will perform.
// It's the contract with the outside world and is defined so we can mock it for testing.
type Service interface {
	AuthenticateMagicKey(handle, key string) (domains.User, error)
}

// service defines the service we are implementing
type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r: r}
}

func (s *service) AuthenticateMagicKey(handle, key string) (domains.User, error) {
	log.Printf("services: authn: authenticate magic key: %s\n", key)
	return s.r.GetMagicKeyUser(handle, key)
}

func (s *service) GetUser(userID domains.UserID) (domains.User, error) {
	return s.r.GetUser(userID)
}
