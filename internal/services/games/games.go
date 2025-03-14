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
	GetAllGameInfo() ([]domains.GameInfo, error)
	GetUserGames(domains.UserID) ([]domains.UserGame, error)
}

// Service defines the operations this service will perform.
// It's the contract with the outside world and is defined so we can mock it for testing.
type Service interface {
	GetAllGameInfo() ([]domains.GameInfo, error)
	GetUserGames(domains.UserID) ([]domains.UserGame, error)
}

// service defines the service we are implementing
type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r: r}
}

func (s *service) GetAllGameInfo() ([]domains.GameInfo, error) {
	return s.r.GetAllGameInfo()
}

func (s *service) GetUserGames(userID domains.UserID) ([]domains.UserGame, error) {
	return s.r.GetUserGames(userID)
}
