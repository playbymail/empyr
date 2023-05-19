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

package main

// Engine implements the rules of the game.
// Input is always game state and an order.
// Output is updated game state.
type Engine struct {
}

// NewEngine returns an initialized engine.
func NewEngine() (*Engine, error) {
	e := &Engine{}
	e.initialize()
	return e, nil
}

func (e *Engine) initialize() {
}

func (e *Engine) reset() {
	e.initialize()
}
