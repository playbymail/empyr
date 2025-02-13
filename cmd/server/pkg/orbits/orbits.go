// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package orbits

import (
	"github.com/playbymail/empyr/cmd/server/pkg/planets"
)

type Orbit struct {
	ID     string          `json:"orbit_id"`
	Planet *planets.Planet `json:"planet,omitempty"`
}
