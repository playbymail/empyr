// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package games

import (
	"github.com/playbymail/empyr/internal/cerr"
	"github.com/playbymail/empyr/internal/domains"
)

const (
	// ErrNotFound is used when a session could not be found.
	ErrNotFound = cerr.Error("session not found")
)

// Repository defines the operations required to fetch and store authentication data.
type Repository interface {
	GetAllGames() ([]domains.Game, error)
	GetUsersGames(domains.UserID) ([]domains.Game, error)
}

// Service defines the operations this service will perform.
// It's the contract with the outside world and is defined so we can mock it for testing.
type Service interface {
	GetAllGames() ([]domains.Game, error)
	GetUsersGames(domains.UserID) ([]domains.Game, error)
}

// service defines the service we are implementing
type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r: r}
}

func (s *service) GetAllGames() ([]domains.Game, error) {
	return s.r.GetAllGames()
}

func (s *service) GetUsersGames(userID domains.UserID) ([]domains.Game, error) {
	return s.r.GetUsersGames(userID)
}
