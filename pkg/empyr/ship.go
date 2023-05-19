// empyr - a reimagining of Vern Holford's Empyrean Challenge
// Copyright (C) 2023 Michael D Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package empyr

import (
	"fmt"
)

type Ship struct {
	Mass   int // total mass of the ship and contents, in tonnes
	Volume int // total volume of the ship, in cubic meters

	// Fuel holds the amount of fuel available to power the ship and engines.
	Fuel float64

	// JumpDrives holds the number and type of all hyper-engines on the ship.
	JumpDrives []HyperEngine
	// OrbitalDrives holds the number and type of all space-drives on the ship.
	OrbitalDrives []SpaceDrive

	Location struct {
		Current  Location  // current location of the ship
		Previous *Location // if set, where the ship jumped from
	}
}

// estimateJumpCost does
func (s Ship) estimateJumpCost(to Location) (fuelConsumed float64, err error) {
	distance := s.Location.Current.DistanceFrom(to)

	// get current drive setup
	drives := s.getJumpDrives()

	// fuel needed is 40 units per light year jumped per jump drive
	fuelConsumed = 40 * distance * float64(drives.Quantity)
	if fuelConsumed > s.Fuel {
		return 0, fmt.Errorf("jump-drive: %w", ErrInsufficientFuel)
	}

	// each jump drive can move at most TL * 1000 metric tonnes
	maxMass := drives.TechLevel * 1_000 * drives.Quantity
	if s.Mass > maxMass {
		return 0, fmt.Errorf("jump-drive: %w", ErrMassExceedsCapacity)
	}

	// maximum distance per jump is 1 light year per tech level of the jump drive
	maxDistance := float64(drives.TechLevel)
	if distance > maxDistance {
		return 0, fmt.Errorf("jump-drive: %w", ErrDistanceExceedsCapacity)
	}

	return fuelConsumed, nil
}

// getJumpDrives returns the current jump drive configuration.
// The game penalizes the player for mixing the tech level of drives.
// The combined tech level is the lowest of all the installed drives.
func (s Ship) getJumpDrives() (e struct{ TechLevel, Quantity int }) {
	for n, engine := range s.JumpDrives {
		e.Quantity += engine.Quantity
		if n == 0 || e.TechLevel > engine.TechLevel {
			e.TechLevel = engine.TechLevel
		}
	}
	return e
}

// Jump moves a ship to a new location.
// Executing this command does not update the ship - it only returns the new location, fuel consumed, and any errors.
func (s Ship) Jump(to Location) (Location, float64, error) {
	fuelConsumed, err := s.estimateJumpCost(to)
	if err != nil {
		return s.Location.Current, fuelConsumed, err
	}
	return to, fuelConsumed, nil
}
