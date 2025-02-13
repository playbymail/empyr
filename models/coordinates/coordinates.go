// Copyright (c) 2025 Michael D Henderson. All rights reserved.

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
