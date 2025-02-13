// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package resources

import (
	"github.com/google/uuid"
	"github.com/playbymail/empyr/cmd/server/pkg/prng"
)

type Generator func(prng.Generator) (*Resource, error)

// DefaultGenerator returns a generator with the following distribution:
//
//	 Fuel: 33%, Yield 10..90%, Units: 10..99 million
//	 Gold:  1%, Yield  1.. 9%, Units: 10..99 million
//	Metal: 32%, Yield 10..90%, Units: 10..99 million
//
// NonMetal: 34%, Yield 10..90%, Units: 10..99 million
func DefaultGenerator() Generator {
	return func(ts prng.Generator) (*Resource, error) {
		var nr Resource
		nr.ID = uuid.New().String()
		switch roll := ts.Intn(100); {
		case roll < 33:
			nr.Type = FUEL
			nr.YieldPct = float64(10+ts.Intn(80)) / 100
			nr.InitialAmount = 1000000 * int64(1+ts.Intn(99))
		case roll < 34:
			nr.Type = GOLD
			nr.YieldPct = float64(1+ts.Intn(9)) / 100
			nr.InitialAmount = 1000000 * int64(1+ts.Intn(99))
		case roll < 66:
			nr.Type = METAL
			nr.YieldPct = float64(10+ts.Intn(80)) / 100
			nr.InitialAmount = 1_000_000 * int64(1+ts.Intn(99))
		default:
			nr.Type = NONMETAL
			nr.YieldPct = float64(10+ts.Intn(80)) / 100
			nr.InitialAmount = 1_000_000 * int64(1+ts.Intn(99))
		}
		nr.AmountRemaining = nr.InitialAmount
		return &nr, nil
	}
}
