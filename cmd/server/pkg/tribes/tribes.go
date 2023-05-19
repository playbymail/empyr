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
