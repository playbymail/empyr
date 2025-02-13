// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package ships

import "github.com/playbymail/empyr/models/coordinates"

// Ship is either a ship or a colony(?!!?).
type Ship struct {
	Id       string // unique identifier for ship or colony
	Kind     Kind
	Location coordinates.Coordinates
	// attributes like hull, cargo, bridge, engines
}
