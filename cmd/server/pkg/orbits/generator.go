// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package orbits

import (
	"github.com/google/uuid"
	"github.com/playbymail/empyr/cmd/server/pkg/planets"
	"github.com/playbymail/empyr/cmd/server/pkg/prng"
)

type Generator func(generator prng.Generator) (*Orbit, error)

// DefaultGenerator returns a generator with the following rules:
//
//	75% chance orbit contains a planet
func DefaultGenerator() Generator {
	generatePlanet := planets.DefaultGenerator()
	return func(ts prng.Generator) (*Orbit, error) {
		var o Orbit
		o.ID = uuid.New().String()
		if ts.Intn(4) == 0 {
			// we're creating an empty orbit to give ships and colonies something to park in.
		} else {
			// we're creating a planet or asteroid belt in this orbit
			p, err := generatePlanet(ts)
			if err != nil {
				return nil, err
			}
			o.Planet = p
		}
		return &o, nil
	}
}
