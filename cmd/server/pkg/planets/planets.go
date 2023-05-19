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

package planets

import (
	"encoding/json"
	"fmt"
	"github.com/playbymail/empyr/cmd/server/pkg/resources"
)

type Planet struct {
	ID           string
	Type         PlanetType
	Habitability int // range from 0 to 25
	Deposits     []*resources.Resource
}

type PlanetType int

const (
	ASTEROIDBELT PlanetType = iota
	GASGIANT
	TERRESTRIAL
)

func (t PlanetType) String() string {
	switch t {
	case ASTEROIDBELT:
		return "asteroid-belt"
	case GASGIANT:
		return "gas-giant"
	case TERRESTRIAL:
		return "terrestrial"
	}
	panic(fmt.Sprintf("assert(planetType != %d)", t))
}

func (p Planet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID           string                `json:"planet_id"`
		Type         string                `json:"type"`
		Habitability int                   `json:"habitability"`
		Deposits     []*resources.Resource `json:"deposits"`
	}{
		ID:           p.ID,
		Type:         p.Type.String(),
		Habitability: p.Habitability,
		Deposits:     p.Deposits,
	})
}
