// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package systems

import "github.com/playbymail/empyr/cmd/server/pkg/stars"

type System struct {
	ID    string        `json:"system_id"`
	Name  string        `json:"name"`
	X     int           `json:"x"`
	Y     int           `json:"y"`
	Z     int           `json:"z"`
	Stars []*stars.Star `json:"stars,omitempty"`
}
