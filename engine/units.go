// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import "fmt"

var operationalUnits = map[string]bool{
	"AUT": true,
	"ESH": true,
	"EWP": true,
	"FCT": true,
	"FRM": true,
	"HEN": true,
	"LAB": true,
	"LFS": true,
	"MIN": true,
	"MSL": true,
	"PWP": true,
	"SEN": true,
	"SLS": true,
	"SPD": true,
	"STU": true,
}

func IsOperational(code string) bool {
	return operationalUnits[code]
}

func Mass(code string, techLevel int64, quantity int64) float64 {
	tl, qty, massPerUnit := float64(techLevel), float64(quantity), 0.0
	switch code {
	case "ANM":
		massPerUnit = 4 * tl
	case "ASC":
		massPerUnit = 5 * tl
	case "ASW":
		massPerUnit = 2 * tl
	case "AUT":
		massPerUnit = 4 * tl
	case "CNGD":
		massPerUnit = 0.6
	case "ESH":
		massPerUnit = 20 * tl
	case "EWP":
		massPerUnit = 10 * tl
	case "FCT":
		massPerUnit = 2*tl + 12
	case "FOOD":
		massPerUnit = 6
	case "FRM":
		massPerUnit = 2*tl + 6
	case "FUEL":
		massPerUnit = 1
	case "GOLD":
		massPerUnit = 1
	case "HEN":
		massPerUnit = 45 * tl
	case "LAB":
		massPerUnit = 2*tl + 8
	case "LFS":
		massPerUnit = 8 * tl
	case "METS":
		massPerUnit = 1
	case "MIN":
		massPerUnit = 2*tl + 10
	case "MSL":
		massPerUnit = 25 * tl
	case "MSS":
		massPerUnit = 4 * tl
	case "MTBT":
		massPerUnit = 2*tl + 20
	case "MTSP":
		massPerUnit = 0.04
	case "NMTS":
		massPerUnit = 1
	case "PWP":
		massPerUnit = 2*tl + 10
	case "RPV":
		massPerUnit = 500 / tl
	case "RSCH":
		massPerUnit = 0
	case "SEN":
		massPerUnit = 3000 * tl
	case "SLS":
		massPerUnit = 0.01 * tl
	case "SPD":
		massPerUnit = 25 * tl
	case "STU":
		massPerUnit = 0.1 * tl
	case "TPT":
		massPerUnit = 4 * tl
	default:
		panic(fmt.Sprintf("assert(code != %q)", code))
	}
	return massPerUnit * qty
}

// Volume returns the volume of a unit, assuming that it is not operational has is not in storage.
func Volume(code string, techLevel int64, quantity int64) float64 {
	tl, qty, volumePerUnit := float64(techLevel), float64(quantity), 0.0
	switch code {
	case "ANM":
		volumePerUnit = 4 * tl
	case "ASC":
		volumePerUnit = 5 * tl
	case "ASW":
		volumePerUnit = 2 * tl
	case "AUT":
		volumePerUnit = 2 * tl
	case "CNGD":
		volumePerUnit = 0.3
	case "ESH":
		volumePerUnit = 10 * tl
	case "EWP":
		volumePerUnit = 5 * tl
	case "FCT":
		volumePerUnit = tl + 6
	case "FOOD":
		volumePerUnit = 3
	case "FRM":
		volumePerUnit = tl + 3
	case "FUEL":
		volumePerUnit = 0.5
	case "GOLD":
		volumePerUnit = 0.5
	case "HEN":
		volumePerUnit = 22.5 * tl
	case "LAB":
		volumePerUnit = tl + 4
	case "LFS":
		volumePerUnit = 4 * tl
	case "METS":
		volumePerUnit = 0.5
	case "MIN":
		volumePerUnit = tl + 5
	case "MSL":
		volumePerUnit = 12.5 * tl
	case "MSS":
		volumePerUnit = 4 * tl
	case "MTBT":
		volumePerUnit = tl + 10
	case "MTSP":
		volumePerUnit = 0.02
	case "NMTS":
		volumePerUnit = 0.5
	case "PWP":
		volumePerUnit = tl + 5
	case "RPV":
		volumePerUnit = 500 / tl
	case "RSCH":
		volumePerUnit = 0
	case "SEN":
		volumePerUnit = 1500 * tl
	case "SLS":
		volumePerUnit = 0.005 * tl
	case "SPD":
		volumePerUnit = 12.5 * tl
	case "STU":
		volumePerUnit = 0.05 * tl
	case "TPT":
		volumePerUnit = 4 * tl
	default:
		panic(fmt.Sprintf("assert(code != %q)", code))
	}
	return volumePerUnit * qty
}

func VolumeAssembled(code string, techLevel int64, quantity int64) float64 {
	tl, qty, volumePerUnit := float64(techLevel), float64(quantity), 0.0
	switch code {
	case "ANM":
		volumePerUnit = 4 * tl
	case "ASC":
		volumePerUnit = 5 * tl
	case "ASW":
		volumePerUnit = 2 * tl
	case "AUT":
		volumePerUnit = 2 * tl
	case "CNGD":
		volumePerUnit = 0.3
	case "ESH":
		volumePerUnit = 10 * tl
	case "EWP":
		volumePerUnit = 5 * tl
	case "FCT":
		volumePerUnit = tl + 6
	case "FOOD":
		volumePerUnit = 3
	case "FRM":
		volumePerUnit = tl + 3
	case "FUEL":
		volumePerUnit = 0.5
	case "GOLD":
		volumePerUnit = 0.5
	case "HEN":
		volumePerUnit = 22.5 * tl
	case "LAB":
		volumePerUnit = tl + 4
	case "LFS":
		volumePerUnit = 4 * tl
	case "METS":
		volumePerUnit = 0.5
	case "MIN":
		volumePerUnit = tl + 5
	case "MSL":
		volumePerUnit = 12.5 * tl
	case "MSS":
		volumePerUnit = 4 * tl
	case "MTBT":
		volumePerUnit = tl + 10
	case "MTSP":
		volumePerUnit = 0.02
	case "NMTS":
		volumePerUnit = 0.5
	case "PWP":
		volumePerUnit = tl + 5
	case "RPV":
		volumePerUnit = 500 / tl
	case "RSCH":
		volumePerUnit = 0
	case "SEN":
		volumePerUnit = 1500 * tl
	case "SLS":
		volumePerUnit = 0.005 * tl
	case "SPD":
		volumePerUnit = 12.5 * tl
	case "STU":
		volumePerUnit = 0.05 * tl
	case "TPT":
		volumePerUnit = 4 * tl
	default:
		panic(fmt.Sprintf("assert(code != %q)", code))
	}
	return volumePerUnit * qty
}

func VolumeDisassembled(code string, techLevel int64, quantity int64) float64 {
	tl, qty, volumePerUnit := float64(techLevel), float64(quantity), 0.0
	switch code {
	case "ANM":
		volumePerUnit = 4 * tl
	case "ASC":
		volumePerUnit = 5 * tl
	case "ASW":
		volumePerUnit = 2 * tl
	case "AUT":
		volumePerUnit = 2 * tl
	case "CNGD":
		volumePerUnit = 0.3
	case "ESH":
		volumePerUnit = 10 * tl
	case "EWP":
		volumePerUnit = 5 * tl
	case "FCT":
		volumePerUnit = tl + 6
	case "FOOD":
		volumePerUnit = 3
	case "FRM":
		volumePerUnit = tl + 3
	case "FUEL":
		volumePerUnit = 0.5
	case "GOLD":
		volumePerUnit = 0.5
	case "HEN":
		volumePerUnit = 22.5 * tl
	case "LAB":
		volumePerUnit = tl + 4
	case "LFS":
		volumePerUnit = 4 * tl
	case "METS":
		volumePerUnit = 0.5
	case "MIN":
		volumePerUnit = tl + 5
	case "MSL":
		volumePerUnit = 12.5 * tl
	case "MSS":
		volumePerUnit = 4 * tl
	case "MTBT":
		volumePerUnit = tl + 10
	case "MTSP":
		volumePerUnit = 0.02
	case "NMTS":
		volumePerUnit = 0.5
	case "PWP":
		volumePerUnit = tl + 5
	case "RPV":
		volumePerUnit = 500 / tl
	case "RSCH":
		volumePerUnit = 0
	case "SEN":
		volumePerUnit = 1500 * tl
	case "SLS":
		volumePerUnit = 0.005 * tl
	case "SPD":
		volumePerUnit = 12.5 * tl
	case "STU":
		volumePerUnit = 0.05 * tl
	case "TPT":
		volumePerUnit = 4 * tl
	default:
		panic(fmt.Sprintf("assert(code != %q)", code))
	}
	return volumePerUnit * qty
}

func VolumeStored(code string, techLevel int64, quantity int64) float64 {
	tl, qty, volumePerUnit := float64(techLevel), float64(quantity), 0.0
	switch code {
	case "ANM":
		volumePerUnit = 4 * tl
	case "ASC":
		volumePerUnit = 5 * tl
	case "ASW":
		volumePerUnit = 2 * tl
	case "AUT":
		volumePerUnit = 2 * tl
	case "CNGD":
		volumePerUnit = 0.3
	case "ESH":
		volumePerUnit = 10 * tl
	case "EWP":
		volumePerUnit = 5 * tl
	case "FCT":
		volumePerUnit = tl + 6
	case "FOOD":
		volumePerUnit = 3
	case "FRM":
		volumePerUnit = tl + 3
	case "FUEL":
		volumePerUnit = 0.5
	case "GOLD":
		volumePerUnit = 0.5
	case "HEN":
		volumePerUnit = 22.5 * tl
	case "LAB":
		volumePerUnit = tl + 4
	case "LFS":
		volumePerUnit = 4 * tl
	case "METS":
		volumePerUnit = 0.5
	case "MIN":
		volumePerUnit = tl + 5
	case "MSL":
		volumePerUnit = 12.5 * tl
	case "MSS":
		volumePerUnit = 4 * tl
	case "MTBT":
		volumePerUnit = tl + 10
	case "MTSP":
		volumePerUnit = 0.02
	case "NMTS":
		volumePerUnit = 0.5
	case "PWP":
		volumePerUnit = tl + 5
	case "RPV":
		volumePerUnit = 500 / tl
	case "RSCH":
		volumePerUnit = 0
	case "SEN":
		volumePerUnit = 1500 * tl
	case "SLS":
		volumePerUnit = 0.005 * tl
	case "SPD":
		volumePerUnit = 12.5 * tl
	case "STU":
		volumePerUnit = 0.05 * tl
	case "TPT":
		volumePerUnit = 4 * tl
	default:
		panic(fmt.Sprintf("assert(code != %q)", code))
	}
	return volumePerUnit * qty
}

func ResourcesToBuild(code string, techLevel int64, quantity int64) (float64, float64) {
	tl, qty, metalsPerUnit, nonMetalsPerUnit := float64(techLevel), float64(quantity), 0.0, 0.0
	switch code {
	case "ANM":
		metalsPerUnit, nonMetalsPerUnit = 2*tl, 2*tl
	case "ASC":
		metalsPerUnit, nonMetalsPerUnit = 3*tl, 2*tl
	case "ASW":
		metalsPerUnit, nonMetalsPerUnit = 1*tl, 1*tl
	case "AUT":
		metalsPerUnit, nonMetalsPerUnit = 2*tl, 2*tl
	case "CNGD":
		metalsPerUnit, nonMetalsPerUnit = 0.2, 0.4
	case "ESH":
		metalsPerUnit, nonMetalsPerUnit = 10*tl, 10*tl
	case "EWP":
		metalsPerUnit, nonMetalsPerUnit = 5*tl, 5*tl
	case "FCT":
		metalsPerUnit, nonMetalsPerUnit = 8+tl, 4+tl
	case "FOOD":
		metalsPerUnit, nonMetalsPerUnit = 0, 0
	case "FRM":
		metalsPerUnit, nonMetalsPerUnit = 4+tl, 2+tl
	case "FUEL":
		metalsPerUnit, nonMetalsPerUnit = 0, 0
	case "GOLD":
		metalsPerUnit, nonMetalsPerUnit = 0, 0
	case "HEN":
		metalsPerUnit, nonMetalsPerUnit = 25*tl, 20*tl
	case "LAB":
		metalsPerUnit, nonMetalsPerUnit = 5+tl, 3+tl
	case "LFS":
		metalsPerUnit, nonMetalsPerUnit = 3*tl, 5*tl
	case "METS":
		metalsPerUnit, nonMetalsPerUnit = 0, 0
	case "MIN":
		metalsPerUnit, nonMetalsPerUnit = 5+tl, 5+tl
	case "MSL":
		metalsPerUnit, nonMetalsPerUnit = 15+tl, 10+tl
	case "MSS":
		metalsPerUnit, nonMetalsPerUnit = 2*tl, 2*tl
	case "MTBT":
		metalsPerUnit, nonMetalsPerUnit = 10+tl, 10+tl
	case "MTSP":
		metalsPerUnit, nonMetalsPerUnit = 0.02, 0.02
	case "NMTS":
		metalsPerUnit, nonMetalsPerUnit = 0, 0
	case "PWP":
		metalsPerUnit, nonMetalsPerUnit = 5+tl, 5+tl
	case "RPV":
		metalsPerUnit, nonMetalsPerUnit = 200/tl, 300/tl
	case "RSCH":
		metalsPerUnit, nonMetalsPerUnit = 0, 0
	case "SEN":
		metalsPerUnit, nonMetalsPerUnit = 1000*tl, 2000*tl
	case "SLS":
		metalsPerUnit, nonMetalsPerUnit = 0.005*tl, 0.005*tl
	case "SPD":
		metalsPerUnit, nonMetalsPerUnit = 15*tl, 10*tl
	case "STU":
		metalsPerUnit, nonMetalsPerUnit = 0.07*tl, 0.03*tl
	case "TPT":
		metalsPerUnit, nonMetalsPerUnit = 2*tl, 2*tl
	default:
		panic(fmt.Sprintf("assert(code != %q)", code))
	}
	return metalsPerUnit * qty, nonMetalsPerUnit * qty
}
