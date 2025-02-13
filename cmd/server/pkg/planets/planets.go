// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package planets

import (
	"encoding/json"
	"fmt"
	"github.com/playbymail/empyr/cmd/server/pkg/resources"
)

type Planet struct {
	ID           string
	Type         PlanetType
	Habitability int // range from 0 to 25
	Deposits     []*resources.Resource
}

type PlanetType int

const (
	ASTEROIDBELT PlanetType = iota
	GASGIANT
	TERRESTRIAL
)

func (t PlanetType) String() string {
	switch t {
	case ASTEROIDBELT:
		return "asteroid-belt"
	case GASGIANT:
		return "gas-giant"
	case TERRESTRIAL:
		return "terrestrial"
	}
	panic(fmt.Sprintf("assert(planetType != %d)", t))
}

func (p Planet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID           string                `json:"planet_id"`
		Type         string                `json:"type"`
		Habitability int                   `json:"habitability"`
		Deposits     []*resources.Resource `json:"deposits"`
	}{
		ID:           p.ID,
		Type:         p.Type.String(),
		Habitability: p.Habitability,
		Deposits:     p.Deposits,
	})
}
