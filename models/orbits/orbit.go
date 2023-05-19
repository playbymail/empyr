// empyr - a game engine for Empyrean Challenge
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

package orbits

import (
	"bytes"
	"fmt"
	"github.com/playbymail/empyr/models/coordinates"
)

type Orbit struct {
	Id           string // unique identifier for the orbit
	Location     coordinates.Coordinates
	Kind         OrbitKind // kind of orbit
	Habitability int       // range 0..25
	ControlledBy string    // id of nation controlling this orbit
	Colonies     struct {
		Open    []string // id of open surface colonies
		Closed  []string // id of closed surface colonies
		Orbital []string // id of orbital colonies
	}
	Deposits []Deposit // deposits of resources
}

type OrbitKind int

const (
	Empty OrbitKind = iota
	AsteroidBelt
	GasGiant
	Terrestrial
)

// MarshalJSON implements the Marshaler interface.
func (k OrbitKind) MarshalJSON() ([]byte, error) {
	switch k {
	case Empty:
		return []byte(`null`), nil
	case AsteroidBelt:
		return []byte(`"asteroid-belt"`), nil
	case GasGiant:
		return []byte(`"gas-giant"`), nil
	case Terrestrial:
		return []byte(`"terrestrial"`), nil
	}
	return nil, fmt.Errorf("invalid orbit")
}

// UnmarshalJSON implements the Unmarshaler interface.
func (k *OrbitKind) UnmarshalJSON(b []byte) error {
	if b == nil {
		*k = Empty
		return nil
	} else if bytes.Compare(b, []byte(`"asteroid-belt"`)) == 0 {
		*k = AsteroidBelt
		return nil
	} else if bytes.Compare(b, []byte(`"gas-giant"`)) == 0 {
		*k = GasGiant
		return nil
	} else if bytes.Compare(b, []byte(`"terrestrial"`)) == 0 {
		*k = Terrestrial
		return nil
	}
	return fmt.Errorf("invalid orbit")
}

type Deposit struct {
	Id           string // unique identifier for deposit
	Resource     Resource
	ControlledBy string // id of nation controlling this deposit
	QtyInitial   int
	QtyRemaining int
}
type Resource int
