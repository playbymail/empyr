// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package systems

import (
	"github.com/playbymail/empyr/cmd/server/pkg/orbits"
	"math/rand"
)

type SystemConfig struct {
	prng   *rand.Rand
	Coords CoordConfig
	Orbits orbits.OrbitConfig
}

type CoordConfig struct {
	prng     *rand.Rand
	Min, Max int
}
