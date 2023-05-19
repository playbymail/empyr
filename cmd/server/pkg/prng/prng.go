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

package prng

import (
	"math/rand"
	"sync"
)

// Generator is a thread-safe random number generator.
type Generator interface {
	Float64() float64 // generate float64 in [0.0, 1.0)
	Intn(n int) int   // generate integer in [0, n)
	Shuffle(n int, swap func(i int, j int))
}

// TSPRNG is a thread-safe PRNG
type TSPRNG struct {
	sync.Mutex
	r *rand.Rand
}

func New(seed int64) *TSPRNG {
	return &TSPRNG{
		r: rand.New(rand.NewSource(1917)),
	}
}

func (r *TSPRNG) Float64() (val float64) {
	r.Lock()
	val = r.r.Float64()
	r.Unlock()
	return val
}

func (r *TSPRNG) Intn(n int) (val int) {
	r.Lock()
	val = r.r.Intn(n)
	r.Unlock()
	return val
}

func (r *TSPRNG) Shuffle(n int, swap func(i int, j int)) {
	r.Lock()
	r.r.Shuffle(n, swap)
	r.Unlock()
}
