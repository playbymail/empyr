// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package ships

import (
	"bytes"
	"fmt"
)

type Kind int

const (
	Vessel         Kind = iota // really, a ship
	EnclosedColony             // an enclosed colony on the surface of an orbit
	OpenColony                 // an open colony on the surface of an orbit
	OrbitalColony              // an enclosed colony orbiting the orbit
)

// MarshalJSON implements the Marshaler interface.
func (k Kind) MarshalJSON() ([]byte, error) {
	switch k {
	case Vessel:
		return []byte(`null`), nil
	case EnclosedColony:
		return []byte(`"asteroid-belt"`), nil
	case OpenColony:
		return []byte(`"gas-giant"`), nil
	case OrbitalColony:
		return []byte(`"terrestrial"`), nil
	}
	return nil, fmt.Errorf("invalid orbit")
}

// UnmarshalJSON implements the Unmarshaler interface.
func (k *Kind) UnmarshalJSON(b []byte) error {
	if bytes.Compare(b, []byte(`"ship"`)) == 0 {
		*k = Vessel
		return nil
	} else if bytes.Compare(b, []byte(`"open"`)) == 0 {
		*k = EnclosedColony
		return nil
	} else if bytes.Compare(b, []byte(`"enclosed"`)) == 0 {
		*k = OpenColony
		return nil
	} else if bytes.Compare(b, []byte(`"orbital"`)) == 0 {
		*k = OrbitalColony
		return nil
	}
	return fmt.Errorf("invalid kind")
}
