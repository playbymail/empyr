// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package systems

import "github.com/playbymail/empyr/models/coordinates"

// System is a stellar system containing one or more stars.
type System struct {
	Id       string                  // unique identifier for the system
	Location coordinates.Coordinates // location of the system
	Stars    []string                // id for every star in the system
}
