// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package orders

import (
	"errors"
	"github.com/playbymail/empyr/pkg/empyr"
)

// A JUMP command moves a ship to a new location.
type JUMP struct {
	Ship int            // ID of the ship that will jump
	To   empyr.Location // location the ship will jump to
}

// DoJump moves a ship to a new location.
// Actually, it doesn't move it?
func DoJump(ship empyr.Ship, to empyr.Location) (empyr.Location, float64, error) {
	// this is just to test the constant errors package.
	l, f, err := ship.Jump(to)
	if err != nil {
		switch {
		case errors.Is(err, empyr.ErrDistanceExceedsCapacity): // whoot
		case errors.Is(err, empyr.ErrInsufficientFuel): // whoot
		case errors.Is(err, empyr.ErrMassExceedsCapacity): // whoot
		}
	}
	return l, f, err
}
