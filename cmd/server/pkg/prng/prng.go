// Copyright (c) 2025 Michael D Henderson. All rights reserved.

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
