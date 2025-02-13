// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package planets

import (
	"github.com/playbymail/empyr/cmd/server/pkg/prng"
	"github.com/playbymail/empyr/cmd/server/pkg/resources"
	"log"
)

type Generator func(prng.Generator) (*Planet, error)

// Default generator returns a generator with the following distribution:
//
//	    Habitability: 0..24 million
//	NaturalResources: 1..40 deposits
func DefaultGenerator() Generator {
	minHabitability, maxHabitability := 0, 25
	log.Printf("[planet] habitability: min %d max %d\n", minHabitability, maxHabitability)
	minDeposits, maxDeposits := 1, 4 // s/b 40
	log.Printf("[planet] deposits: min %d max %d\n", minDeposits, maxDeposits)
	generateResource := resources.DefaultGenerator()
	return func(ts prng.Generator) (*Planet, error) {
		var p Planet
		switch ts.Intn(3) {
		case 0:
			p.Type = ASTEROIDBELT
		case 1:
			p.Type = GASGIANT
		default:
			p.Type = TERRESTRIAL
		}
		p.Habitability = (minHabitability + ts.Intn(maxHabitability-minHabitability)) * 1_000_000
		numberOfDeposits := minDeposits + ts.Intn(maxDeposits-minDeposits)
		p.Deposits = make([]*resources.Resource, numberOfDeposits, numberOfDeposits)
		for i := range p.Deposits {
			resource, err := generateResource(ts)
			if err != nil {
				return nil, err
			}
			p.Deposits[i] = resource
		}
		return &p, nil
	}
}

// GenerateHomeworld returns a nice planet for a home world
func GenerateHomeworld() Generator {
	minHabitability, maxHabitability := 0, 25
	log.Printf("[home-world] habitability: min %d max %d\n", minHabitability, maxHabitability)
	minDeposits, maxDeposits := 1, 4 // s/b 40
	log.Printf("[home-world] deposits: min %d max %d\n", minDeposits, maxDeposits)
	generateResource := resources.DefaultGenerator()
	return func(ts prng.Generator) (*Planet, error) {
		var p Planet
		p.Type = TERRESTRIAL
		p.Habitability = maxHabitability * 1_000_000
		numberOfDeposits := maxDeposits
		p.Deposits = make([]*resources.Resource, numberOfDeposits, numberOfDeposits)
		for i := range p.Deposits {
			switch i {
			case 0:
				p.Deposits[i] = &resources.Resource{
					Type:      resources.FUEL,
					YieldPct:  9,
					Unlimited: true,
				}
			case 1:
				p.Deposits[i] = &resources.Resource{
					Type:      resources.METAL,
					YieldPct:  9,
					Unlimited: true,
				}
			case 2:
				p.Deposits[i] = &resources.Resource{
					Type:      resources.NONMETAL,
					YieldPct:  9,
					Unlimited: true,
				}
			default:
				resource, err := generateResource(ts)
				if err != nil {
					return nil, err
				}
				p.Deposits[i] = resource
			}
		}
		return &p, nil
	}
}
