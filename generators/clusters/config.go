// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package clusters

import "github.com/playbymail/empyr/generators/points"

type config struct {
	initSystems   int                  // number of systems to seed cluster with
	mapFile       string               // if set, create a map
	pgen          func() *points.Point // points generator
	clustered     bool
	radius        float64
	sphereSize    float64
	templatesPath string // path to template files
}
