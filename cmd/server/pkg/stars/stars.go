// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package stars

import (
	"github.com/playbymail/empyr/cmd/server/pkg/orbits"
)

// Star may be a member of a multiple-star system.
// By the way, I think that orbit 11 is the default jump point target.
type Star struct {
	ID     string
	Name   string
	Orbits [11]*orbits.Orbit
}
