// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import "fmt"

// empire_turn_results_t is the JSON structure for the empire turn results.
type empire_turn_results_t struct {
	Turn    int `json:"turn"`
	Empire  int `json:"empire"`
	Systems []struct {
		Id          int `json:"id"`
		Coordinates struct {
			X int `json:"x"`
			Y int `json:"y"`
			Z int `json:"z"`
		} `json:"coordinates"`
		Stars []int `json:"stars"`
	} `json:"systems"`
	Stars []struct {
		Id       int    `json:"id"`
		SystemId int    `json:"system-id"`
		Sequence string `json:"sequence"`
	} `json:"stars"`
	Orbits []struct {
		Id        int    `json:"id"`
		StarId    int    `json:"star-id"`
		OrbitNo   int    `json:"orbit-no"`
		Kind      string `json:"kind"`
		Habitable bool   `json:"habitable,omitempty"`
	} `json:"orbits"`
	Planets []struct {
		Id           int `json:"id"`
		OrbitId      int `json:"orbit-id"`
		Habitability int `json:"habitability,omitempty"`
	} `json:"planets"`
	Deposits []struct {
		Id       int     `json:"id"`
		PlanetId int     `json:"planet-id"`
		Qty      int     `json:"qty"`
		Kind     string  `json:"kind"`
		Yield    float64 `json:"yield"`
	} `json:"deposits"`
	Colonies   []*turn_results_colony_t
	Ships      []interface{} `json:"ships"`
	Population []struct {
		Id       int `json:"id"`
		Location int `json:"location"`
		UEM      int `json:"uem"`
		USK      int `json:"usk"`
		PRO      int `json:"pro"`
		SLD      int `json:"sld"`
		CNW      int `json:"cnw"`
	} `json:"population"`
	FactoryGroups []struct {
		Id     int    `json:"id"`
		Qty    int    `json:"qty"`
		TL     int    `json:"tl"`
		Orders string `json:"orders"`
		Wip    struct {
			Q1 struct {
				Unit string `json:"unit"`
				Qty  int    `json:"qty"`
			} `json:"q1"`
			Q2 struct {
				Unit string `json:"unit"`
				Qty  int    `json:"qty"`
			} `json:"q2"`
			Q3 struct {
				Unit string `json:"unit"`
				Qty  int    `json:"qty"`
			} `json:"q3"`
		} `json:"wip"`
	} `json:"factory-groups"`
	MiningGroups []struct {
		Id        int `json:"id"`
		DepositId int `json:"deposit-id"`
		Units     []struct {
			TL  int `json:"tl"`
			Qty int `json:"qty"`
		} `json:"units"`
	} `json:"mining-groups"`
}

type turn_results_colony_t struct {
	Id       int     `json:"id"`
	PlanetId int     `json:"planet-id"`
	Kind     string  `json:"kind"`
	Name     string  `json:"name"`
	TL       int     `json:"tl"`
	SOL      float64 `json:"sol"`
	Vitals   struct {
		BirthRate float64 `json:"birth-rate"`
		DeathRate float64 `json:"death-rate"`
		Rations   float64 `json:"rations"`
		PayRates  struct {
			USK float64 `json:"usk"`
			PRO float64 `json:"pro"`
			SLD float64 `json:"sld"`
		} `json:"pay-rates"`
		Census []int `json:"census"`
	} `json:"vitals"`
	Inventory []struct {
		Unit      string `json:"unit"`
		Qty       int    `json:"qty"`
		Assembled bool   `json:"assembled"`
		Storage   bool   `json:"storage"`
	} `json:"inventory"`
	Factories []int `json:"factories"`
	Mines     []int `json:"mines"`
}

type turn_report_payload_t struct {
	EmpireId string // formatted as 3-digit number
	TurnNo   int
	Systems  []*turn_report_system_t
	Stars    []*turn_report_star_t
	Orbits   []*turn_report_orbit_t
	Planets  []*turn_report_planet_t
	Colonies []*turn_report_colony_t
}

type turn_report_system_t struct {
	Id          int
	Coordinates struct {
		X int
		Y int
		Z int
	}
	Stars []*turn_report_star_t
}

type turn_report_star_t struct {
	Id          int
	System      *turn_report_system_t
	Sequence    string
	Coordinates string
}

type turn_report_orbit_t struct {
	Id        int
	Star      *turn_report_star_t
	OrbitNo   int
	Kind      string
	Habitable bool
}

type turn_report_planet_t struct {
	Id           int
	Orbit        *turn_report_orbit_t
	Habitability int
}

type turn_report_colony_t struct {
	Id     int
	Planet *turn_report_planet_t
	Name   string
	Kind   string
	TL     string
	SOL    string
	Vitals struct {
		BirthRate string
		DeathRate string
		Rations   string
		PayRates  struct {
			USK string
			PRO string
			SLD string
		}
		Census []int
	}
	Inventory []*turn_report_inventory_t
	Factories []int
	Mines     []int
}

type turn_report_inventory_t struct {
	Unit      string
	Qty       int
	Assembled bool
	Storage   bool
}

func adapt_emp_to_turn_report_payload_t(input *empire_turn_results_t) *turn_report_payload_t {
	var output turn_report_payload_t
	output.EmpireId = fmt.Sprintf("%03d", input.Empire)
	output.TurnNo = input.Turn

	// make some maps to make it easier to find things
	for _, system := range input.Systems {
		s := &turn_report_system_t{
			Id: system.Id,
			Coordinates: struct {
				X int
				Y int
				Z int
			}{
				X: system.Coordinates.X,
				Y: system.Coordinates.Y,
				Z: system.Coordinates.Z,
			},
		}
		output.Systems = append(output.Systems, s)
	}
	systems := map[int]*turn_report_system_t{}
	for _, s := range output.Systems {
		systems[s.Id] = s
	}
	for _, star := range input.Stars {
		system := systems[star.SystemId]
		s := &turn_report_star_t{
			Id:          star.Id,
			System:      system,
			Sequence:    star.Sequence,
			Coordinates: fmt.Sprintf("%02d/%02d/%02d", system.Coordinates.X, system.Coordinates.Y, system.Coordinates.Z),
		}
		system.Stars = append(system.Stars, s)
		output.Stars = append(output.Stars, s)
	}
	stars := map[int]*turn_report_star_t{}
	for _, s := range output.Stars {
		stars[s.Id] = s
		if len(s.System.Stars) > 1 {
			s.Coordinates += s.Sequence
		}
	}
	for _, orbit := range input.Orbits {
		o := &turn_report_orbit_t{
			Id:      orbit.Id,
			Star:    stars[orbit.StarId],
			OrbitNo: orbit.OrbitNo,
			Kind:    orbit.Kind,
		}
		output.Orbits = append(output.Orbits, o)
	}
	orbits := map[int]*turn_report_orbit_t{}
	for _, o := range output.Orbits {
		orbits[o.Id] = o
	}
	for _, planet := range input.Planets {
		p := &turn_report_planet_t{
			Id:           planet.Id,
			Orbit:        orbits[planet.OrbitId],
			Habitability: planet.Habitability,
		}
		output.Planets = append(output.Planets, p)
	}
	planets := map[int]*turn_report_planet_t{}
	for _, p := range output.Planets {
		planets[p.Id] = p
	}
	colonies := map[int]*turn_report_colony_t{}
	for _, colony := range input.Colonies {
		c := &turn_report_colony_t{
			Id:     colony.Id,
			Planet: planets[colony.PlanetId],
			Name:   colony.Name,
			TL:     fmt.Sprintf("%d", colony.TL),
			SOL:    fmt.Sprintf("%5.4f", colony.SOL),
		}
		if c.Name == "" {
			c.Name = "Not Named"
		}
		switch colony.Kind {
		case "open":
			c.Kind = "Open Colony"
		default:
			panic(fmt.Sprintf("assert(colony.kind != %q)", colony.Kind))
		}
		c.Vitals.BirthRate = fmt.Sprintf("%5.4f", 100*colony.Vitals.BirthRate)
		c.Vitals.DeathRate = fmt.Sprintf("%5.4f", 100*colony.Vitals.DeathRate)
		c.Vitals.Rations = fmt.Sprintf("%8.4f", 100*colony.Vitals.Rations)
		c.Vitals.PayRates.USK = fmt.Sprintf("%5.4f", colony.Vitals.PayRates.USK)
		c.Vitals.PayRates.PRO = fmt.Sprintf("%5.4f", colony.Vitals.PayRates.PRO)
		c.Vitals.PayRates.SLD = fmt.Sprintf("%5.4f", colony.Vitals.PayRates.SLD)
		colonies[colony.Id] = c
		output.Colonies = append(output.Colonies, c)
	}

	return &output
}
