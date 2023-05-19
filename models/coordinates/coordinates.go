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

package coordinates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type Coordinates struct { // location being set up
	X, Y, Z int
	System  string // suffix for multi-star system, A...Z
	Orbit   int
	Ring    int // orbital ring 0..127 (0 means on surface)
}

func (c Coordinates) String() string {
	if c.Orbit == 0 {
		return fmt.Sprintf("(%d %d %d%s)", c.X, c.Y, c.Z, c.System)
	}
	return fmt.Sprintf("(%d %d %d%s %d)", c.X, c.Y, c.Z, c.System, c.Orbit)
}

// MarshalJSON implements the Marshaler interface.
func (c Coordinates) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON implements the Unmarshaler interface.
func (c *Coordinates) UnmarshalJSON(b []byte) error {
	if !bytes.HasPrefix(b, []byte{'('}) || !bytes.HasSuffix(b, []byte{')'}) {
		return fmt.Errorf("invalid coordinates")
	}
	b = b[1 : len(b)-2]
	f := bytes.Fields(b)
	if len(f) < 3 || len(f) > 4 {
		return fmt.Errorf("invalid coordinates")
	}
	x, y, z := string(f[0]), string(f[1]), string(f[2])
	var err error
	if c.X, err = strconv.Atoi((x)); err != nil {
		return fmt.Errorf("invalid coordinates")
	}
	if c.Y, err = strconv.Atoi((y)); err != nil {
		return fmt.Errorf("invalid coordinates")
	}
	c.System = z[len(z)-1:]
	if "A" <= c.System && c.System <= "Z" {
		z = z[:len(z)-1]
	} else {
		c.System = ""
	}
	if c.Z, err = strconv.Atoi(z); err != nil {
		return fmt.Errorf("invalid coordinates")
	}
	return nil
}
