// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package tribes

import (
	"encoding/json"
	"github.com/playbymail/empyr/cmd/server/pkg/planets"
	"github.com/playbymail/empyr/cmd/server/pkg/systems"
)

type Tribe struct {
	Name       string          `json:"name"`        // name of the tribe
	HomeSystem *systems.System `json:"home_system"` // tribe's home system
	HomeWorld  *planets.Planet `json:"home_world"`  // tribe's home world
}

func (t *Tribe) MarshalJSON() ([]byte, error) {
	data := struct {
		Name       string `json:"name"`
		HomeSystem string `json:"home_system"`
		HomeWorld  string `json:"home_world"`
	}{
		Name:       t.Name,
		HomeSystem: t.HomeSystem.Name,
		HomeWorld:  t.HomeWorld.ID,
	}
	return json.Marshal(&data)
}
