// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package stars

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/playbymail/empyr/cmd/server/pkg/orbits"
	"github.com/playbymail/empyr/cmd/server/pkg/prng"
)

type Generator func(ts prng.Generator) (*Star, error)

// DefaultGenerator returns a generator with the following rules:
//
//	11 orbits
//	Orbit[0] is treated as the "11th Orbit" in the rulebook.
func DefaultGenerator() Generator {
	generateOrbit := orbits.DefaultGenerator()
	return func(ts prng.Generator) (*Star, error) {
		var s Star
		s.ID = uuid.New().String()
		s.Name = fmt.Sprintf("%02d-%02d-%02d", 0, 0, 0)
		for i := 1; i <= 10; i++ {
			orbit, err := generateOrbit(ts)
			if err != nil {
				return nil, err
			}
			s.Orbits[i] = orbit
		}
		return &s, nil
	}
}
