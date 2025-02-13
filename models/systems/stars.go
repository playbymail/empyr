// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package systems

import (
	"github.com/playbymail/empyr/models/coordinates"
	"github.com/playbymail/empyr/models/orbits"
)

// Star is a single star system containing one or more Orbit(s)
type Star struct {
	Id       string // unique identifier for the star system
	Location coordinates.Coordinates
	Orbits   [11]orbits.Orbit
}
