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

package resources

import (
	"math/rand"

	"github.com/google/uuid"
)

// Config options for the flexible generator.
// The amounts for the resource types are used when generating deposits.
type Config struct {
	Fuel         Parms
	Gold         Parms
	Metal        Parms
	NonMetal     Parms
	sumPortions  int
	fuelPortion  int
	goldPortion  int
	metalPortion int
}

// Parms defines the parameters for generating the deposit.
//
// Portion is used when randomly selecting the type of deposit.
// The value must be greater than zero and the chance of selecting
// the type is (Portion / sum of all portions). Ideally, the sum
// of all Portions will be 100, but that is not required.
//
// Yield is the maximum yield percentage, 1 <= Yield <= 99. If less than 1,
// it will be set to 1. If greater than 99, it will be set to 99.
//
// Units is the number of units (in millions of units) to allocate to the
// deposit. If the minimum is less than 1, it will be set to 10. If the
// maximum is not greater than the minimum, it will be set to 9 times the
// minimum value.
//
// Unlimited is used to set the Unlimited flag on the generated resource.
type Parms struct {
	Portion            int
	MinYield, MaxYield int // as a percentage
	MinUnits, MaxUnits int // in millions of units
	Unlimited          bool
}

// normalize enforces constraints on the parameters.
// We trust the caller mostly, but not completely.
func normalize(p Parms) Parms {
	if p.Portion < 1 {
		p.Portion = 1
	} else if p.Portion > 9999 {
		p.Portion = 9999
	}
	if p.MinYield < 0 {
		p.MinYield = 1
	}
	if p.MaxYield > 99 {
		p.MaxYield = 99
	}
	if !(p.MinYield < p.MaxYield) {
		p.MinYield, p.MaxYield = p.MaxYield, p.MinYield
	}
	if p.MinUnits < 1 {
		p.MinUnits = 10
	}
	if !(p.MinUnits < p.MaxUnits) {
		p.MaxUnits = 9 * p.MinUnits
	}
	return p
}

// Generator returns a fairly customizable generator.
func (c *Config) Generator(prng *rand.Rand) func() (*Resource, error) {
	var cfg Config
	cfg.Fuel.Portion = 33
	cfg.Fuel.MinYield, cfg.Fuel.MaxYield = 10, 90
	cfg.Fuel.MinUnits, cfg.Fuel.MaxUnits = 10, 90
	cfg.Gold.Portion = 1
	cfg.Gold.MinYield, cfg.Gold.MaxYield = 1, 9
	cfg.Gold.MinUnits, cfg.Gold.MaxUnits = 10, 90
	cfg.Metal.Portion = 33
	cfg.Metal.MinYield, cfg.Metal.MaxYield = 10, 90
	cfg.Metal.MinUnits, cfg.Metal.MaxUnits = 10, 90
	cfg.NonMetal.Portion = 33
	cfg.NonMetal.MinYield, cfg.NonMetal.MaxYield = 10, 90
	cfg.NonMetal.MinUnits, cfg.NonMetal.MaxUnits = 10, 90

	// if provided, the user's config will override the defaults.
	if c != nil {
		cfg = *c
	}

	// enforce our constraints on the allocations
	cfg.Fuel = normalize(cfg.Fuel)
	cfg.Gold = normalize(cfg.Gold)
	cfg.Metal = normalize(cfg.Metal)
	cfg.NonMetal = normalize(cfg.NonMetal)

	// derive the final ranges for allocation
	cfg.sumPortions = cfg.Fuel.Portion + cfg.Gold.Portion + cfg.Metal.Portion + cfg.NonMetal.Portion
	cfg.fuelPortion = cfg.Fuel.Portion
	cfg.goldPortion = cfg.fuelPortion + cfg.Gold.Portion
	cfg.metalPortion = cfg.goldPortion + cfg.Metal.Portion

	return func() (*Resource, error) {
		var nr Resource
		nr.ID = uuid.New().String()
		switch roll := prng.Intn(cfg.sumPortions); {
		case roll < cfg.fuelPortion:
			nr.Type = FUEL
			nr.Unlimited = cfg.Fuel.Unlimited
			nr.YieldPct = float64(prng.Intn(cfg.Fuel.MaxYield-cfg.Fuel.MinYield)+cfg.Fuel.MinYield) / 100
			nr.InitialAmount = 1_000_000 * int64(prng.Intn(cfg.Fuel.MaxUnits-cfg.Fuel.MinUnits)+cfg.Fuel.MinUnits)
		case roll < cfg.goldPortion:
			nr.Type = GOLD
			nr.Unlimited = cfg.Gold.Unlimited
			nr.YieldPct = float64(prng.Intn(cfg.Gold.MaxYield-cfg.Gold.MinYield)+cfg.Gold.MinYield) / 100
			nr.InitialAmount = 1_000_000 * int64(prng.Intn(cfg.Gold.MaxUnits-cfg.Gold.MinUnits)+cfg.Gold.MinUnits)
		case roll < cfg.metalPortion:
			nr.Type = METAL
			nr.Unlimited = cfg.Metal.Unlimited
			nr.YieldPct = float64(prng.Intn(cfg.Metal.MaxYield-cfg.Metal.MinYield)+cfg.Metal.MinYield) / 100
			nr.InitialAmount = 1_000_000 * int64(prng.Intn(cfg.Metal.MaxUnits-cfg.Metal.MinUnits)+cfg.Metal.MinUnits)
		default:
			nr.Type = NONMETAL
			nr.Unlimited = cfg.NonMetal.Unlimited
			nr.YieldPct = float64(prng.Intn(cfg.NonMetal.MaxYield-cfg.NonMetal.MinYield)+cfg.NonMetal.MinYield) / 100
			nr.InitialAmount = 1_000_000 * int64(prng.Intn(cfg.NonMetal.MaxUnits-cfg.NonMetal.MinUnits)+cfg.NonMetal.MinUnits)
		}
		nr.AmountRemaining = nr.InitialAmount
		return &nr, nil
	}
}
