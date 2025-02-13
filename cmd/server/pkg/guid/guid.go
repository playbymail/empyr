// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package guid implements a thread-safe ID generator.
// The returned ID will be unique for a given process,
// but will not be globally unique.
// This is because it uses the system clock and a sequence.
//
// The ID returned is YYYYMMDD-HH24MI-SSSS-SSSS-NNNNNNNNNNNN,
// where the date is UTC and N is a zero-padded sequence
// that starts at zero.
package guid

import (
	"fmt"
	"sync"
	"time"
)

// New returns a new generator starting at the sequence.
// It's useful if you're interested in restoring a generator.
func New(seq int) *Generator {
	if seq < 0 {
		seq = 0
	}
	return &Generator{sequence: seq}
}

// CurrVal returns the current value of the sequence.
func (g *Generator) CurrVal() int {
	return g.sequence
}

// Next returns a string that looks like a UUID.
func (g *Generator) Next() string {
	g.private.Lock()
	g.sequence++
	now := time.Now().UTC()
	yyyymmdd := now.Format("20060102")
	hh24mi := now.Format("1504")
	seconds := now.Format("05.000000")
	// 8-4-4-4-12
	uuid := fmt.Sprintf("%s-%s-%s%s-%s-%012d",
		yyyymmdd,
		hh24mi,
		seconds[:2], seconds[3:5],
		seconds[5:],
		g.sequence)
	g.private.Unlock()
	return uuid
}

// NextVal increments the sequence, then returns the new value.
func (g *Generator) NextVal() int {
	g.private.Lock()
	g.sequence++
	val := g.sequence
	g.private.Unlock()
	return val
}

// Generator holds our lock and current value.
type Generator struct {
	sequence int
	private  struct {
		sync.Mutex
	}
}
