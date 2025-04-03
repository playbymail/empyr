// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import "fmt"

// Code   Operational Requirements                                     Output and Notes
// ANM    Missile Launcher of same TL                                  Destroys Missiles; see combat
// ASC    1 soldier or military robot equivalent + 0.1 fuel in combat  Provides 10 x TL combat factors
//                                                                     does not require transports to attack
// ASW    1 soldier or military robot equivalent                       Provides 2 x TL^2 combat factors
// AUT    Must be assembled                                            Replaces 1 x TL^2 Unskilled
//                                                                     see Automation in Production chapter
// CNGD                                                                Consumption determines SOL
// ESH    1 soldier / 100, uses 10 x TL fuel                           Deflects 10 x TL^2 energy units per use
// EWP    1 soldier / 100, uses 4 x TL fuel                            Destroys 10 x TL^2 mass per hit
// FCT    1 professional + 3 unskilled, uses 0.5 fuel or power         Produces 20 x TL mass per turn
//                                                                     see Manufacturing
// FOOD                                                                Feeds 4 to 16 population each turn
//                                                                     see Basic Needs
// FRM    1 professional + 3 unskilled, fuel varies                    Production varies via TL
//                                                                     see Farming
// FUEL                                                                Raw material used by many units
// GOLD                                                                Raw material with no use in the game
// HEN    1 professional / 100, uses 40 fuel per light year            Lift capacity 1045 x TL, range is 3 x √TL
// LAB    3 professional + 1 unskilled, 0.5 fuel/power                 Produces 0.25 research per turn
// LFS    1 x TL fuel or power                                         Supports 1 x TL^2 population
// METS                                                                Raw material used in production
// MIN    1 professional + 3 unskilled, 0.5 fuel or power              Mines 25 x TL per turn in raw ore
//                                                                     Actual net depends on yield of deposit
//                                                                     see Mining
// MSL    1 soldier / 100                                              Launches 1 missile per attack
//                                                                     see Combat
// MSS    Missile Launcher of same TL                                  Destroys 100 x TL^2 Mass
// MTBT   2 x TL military supplies                                     Same as 2 x TL soldiers
// MTSP                                                                Required by soldiers in combat
// NMTS                                                                Raw material used in production
// PWP                                                                 Produces TL Power per turn (think hydro electric)
// RPV                                                                 Obtains probe data, expended when used
// RSCH                                                                Expended to increase TLs
// SEN    uses 0.05 x TL fuel                                          Used to obtain probe information
// SLS    May only be built in Orbiting Colonies                       Encloses (1 x TL^2) divided by type factor
// SPD    1 professional / 100, uses 1 x TL fuel                       Produces 3000 x TL^2 thrust
// STU                                                                 Encloses (1 x TL^2) divided by type factor
// TPT    1 professional / 10, uses 0.1 x TL^2 fuel                    Transports 20 x TL^2 Mass

type Unit_t struct {
	Code         string
	Name         string
	Category     string
	IsAssembly   bool // is assembly and is operational are synonymous
	IsConsumable bool
	IsResource   bool
}

var unitTable = map[string]*Unit_t{
	"ANM":  {Code: "ANM", Name: "Anti-Missiles", Category: "Vehicles", IsAssembly: false, IsConsumable: true, IsResource: false},
	"ASC":  {Code: "ASC", Name: "Assault Craft", Category: "Vehicles", IsAssembly: false, IsConsumable: false, IsResource: false},
	"ASW":  {Code: "ASW", Name: "Assault Weapons", Category: "Vehicles", IsAssembly: false, IsConsumable: false, IsResource: false},
	"AUT":  {Code: "AUT", Name: "Automation", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"CNGD": {Code: "CNGD", Name: "Consumer Goods", Category: "Consumables", IsAssembly: false, IsConsumable: true, IsResource: false},
	"ESH":  {Code: "ESH", Name: "Energy Shields", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"EWP":  {Code: "EWP", Name: "Energy Weapons", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"FCT":  {Code: "FCT", Name: "Factories", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"FOOD": {Code: "FOOD", Name: "Food", Category: "Consumables", IsAssembly: false, IsConsumable: true, IsResource: true},
	"FRM":  {Code: "FRM", Name: "Farms", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"FUEL": {Code: "FUEL", Name: "Fuel", Category: "Consumables", IsAssembly: false, IsConsumable: true, IsResource: true},
	"GOLD": {Code: "GOLD", Name: "Gold", Category: "Consumables", IsAssembly: false, IsConsumable: true, IsResource: true},
	"HEN":  {Code: "HEN", Name: "Hyper Engines", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"LAB":  {Code: "LAB", Name: "Laboratories", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"LFS":  {Code: "LFS", Name: "Life Supports", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"METS": {Code: "METS", Name: "Metals", Category: "Consumables", IsAssembly: false, IsConsumable: true, IsResource: true},
	"MIN":  {Code: "MIN", Name: "Mines", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"MSL":  {Code: "MSL", Name: "Missile Launchers", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"MSS":  {Code: "MSS", Name: "Missiles", Category: "Vehicles", IsAssembly: false, IsConsumable: false, IsResource: false},
	"MTBT": {Code: "MTBT", Name: "Military Robots", Category: "Bots", IsAssembly: false, IsConsumable: false, IsResource: false},
	"MTSP": {Code: "MTSP", Name: "Military Supplies", Category: "Consumables", IsAssembly: false, IsConsumable: true, IsResource: false},
	"NMTS": {Code: "NMTS", Name: "Non-Metals", Category: "Consumables", IsAssembly: false, IsConsumable: true, IsResource: true},
	"PWP":  {Code: "PWP", Name: "Power Plants", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"RPV":  {Code: "RPV", Name: "Robot Probe Vehicles", Category: "Bots", IsAssembly: false, IsConsumable: true, IsResource: false},
	"RSCH": {Code: "RSCH", Name: "Research", Category: "Consumables", IsAssembly: false, IsConsumable: true, IsResource: false},
	"SEN":  {Code: "SEN", Name: "Sensors", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"SLS":  {Code: "SLS", Name: "Light Structure", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"SPD":  {Code: "SPD", Name: "Space Drives", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"STU":  {Code: "STU", Name: "Structure", Category: "Assembly", IsAssembly: true, IsConsumable: false, IsResource: false},
	"TPT":  {Code: "TPT", Name: "Transports", Category: "Vehicles", IsAssembly: false, IsConsumable: false, IsResource: false},
}

// farmFuel returns the amount of fuel required to operate a group of farm units.
// We have to know the tech level of the units and which orbit the colony is in.
//
// The rules allow FRM-1 only in open-air surface colonies on habitable planets.
// The rules also imply that FRM-2 through FRM-5 can't be used on a ship.
// We don't enforce that here.
func farmFuel(sc *Entity_t, techLevel, quantity int64) (fuel float64) {
	// Farm units use no fuel if they are on orbiting colonies within the fifth orbit.
	if sc.IsColony && !sc.IsOnSurface && sc.Location.OrbitNo <= 5 {
		return 0
	}
	tl, qty := float64(techLevel), float64(quantity)
	if techLevel <= 5 {
		// FRM-1 through FRM-5 use 0.5 fuel per unit per tech level
		return tl * 0.5 * qty
	}
	// FRM-6 and above use 1 fuel per unit per tech level
	return tl * 1 * qty
}

func factoryFuel(sc *Entity_t, techLevel, quantity int64) (fuel float64) {
	// Factories use no fuel if they are on orbiting colonies within the fifth orbit.
	if sc.IsColony && !sc.IsOnSurface && sc.Location.OrbitNo <= 5 {
		return 0
	}
	tl, qty := float64(techLevel), float64(quantity)
	// FCT use 0.5 fuel per unit per tech level
	return tl * 0.5 * qty
}

func mineFuel(sc *Entity_t, techLevel, quantity int64) (fuel float64) {
	tl, qty := float64(techLevel), float64(quantity)
	// MIN use 0.5 fuel per unit per tech level
	return tl * 0.5 * qty
}

// unitFuel returns the amount of fuel required to operate a group of units.
// There are some exceptions that this function does not handle:
// 1. farms on orbiting colonies within the fifth orbit require no fuel; they use solar power.
// 2. factories on orbiting colonies within the fifth orbit require no fuel; they use solar power.
// 3. hyper engines use 40 fuel per light year jumped; only the engines needed for the jump require fuel.
func unitFuel(code string, techLevel, quantity int64) (fuel float64) {
	tl, qty := float64(techLevel), float64(quantity)
	switch code {
	case "ESH":
		fuel = tl * 10
	case "EWP":
		fuel = tl * 4
	case "FCT": // some factories can use solar power instead of fuel
		fuel = tl * 0.5
	case "FRM": // some farms can use solar power instead of fuel
		if techLevel <= 5 {
			fuel = tl * 0.5
		} else {
			fuel = tl
		}
	case "HEN": // fuel is actually 40 fuel per light year jumped
		fuel = 40
	case "LAB":
		fuel = tl * 0.5
	case "LFS":
		fuel = tl
	case "MIN":
		fuel = tl * 0.5
	case "SEN":
		fuel = tl * 0.05
	case "SPD":
		fuel = tl
	case "TPT":
		fuel = tl * tl * 0.1
	}
	return fuel * qty
}

// unitMassAndVolume returns the mass and volume of a group of units.
func unitMassAndVolume(code string, techLevel, quantity int64, isAssembled bool) (mass, volume float64) {
	tl, qty := float64(techLevel), float64(quantity)
	switch code {
	case "ANM":
		mass, volume = tl*4, tl*4
	case "ASC":
		mass, volume = tl*5, tl*5
	case "ASW":
		mass, volume = tl*2, tl*2
	case "AUT":
		mass, volume = tl*4, tl*2
	case "CNGD":
		mass, volume = 0.6, 0.3
	case "ESH":
		mass, volume = tl*20, tl*10
	case "EWP":
		mass, volume = tl*10, tl*5
	case "FCT":
		mass, volume = tl*2+12, tl+6
	case "FOOD":
		mass, volume = 6, 3
	case "FRM":
		mass, volume = tl*2+6, tl+3
	case "FUEL":
		mass, volume = 1, 0.5
	case "GOLD":
		mass, volume = 1, 0.5
	case "HEN":
		mass, volume = tl*45, tl*22.5
	case "LAB":
		mass, volume = tl*2+8, tl+4
	case "LFS":
		mass, volume = tl*8, tl*4
	case "METS":
		mass, volume = 1, 0.5
	case "MIN":
		mass, volume = tl*2+10, tl+5
	case "MSL":
		mass, volume = tl*25, tl*12.5
	case "MSS":
		mass, volume = tl*4, tl*4
	case "MTBT":
		mass, volume = tl*2+20, tl+10
	case "MTSP":
		mass, volume = 0.04, 0.02
	case "NMTS":
		mass, volume = 1, 0.5
	case "PWP":
		mass, volume = tl*2+10, tl+5
	case "RPV":
		mass, volume = 500/tl, 500/tl
	case "RSCH":
		mass, volume = 0, 0
	case "SEN":
		mass, volume = tl*3000, tl*1500
	case "SLS":
		mass, volume = tl*0.01, tl*0.005
	case "SPD":
		mass, volume = tl*25, tl*12.5
	case "STU":
		mass, volume = tl*0.1, tl*0.05
	case "TPT":
		mass, volume = tl*4, tl*4
	default:
		panic(fmt.Sprintf("assert(code != %q)", code))
	}
	if !isAssembled {
		volume = volume / 2
	}
	return mass * qty, volume * qty
}

// unitRequirements returns the metals and non-metals required to build a group of units.
func unitRequirements(code string, techLevel, quantity int64) (mets, nmts float64) {
	tl, qty := float64(techLevel), float64(quantity)
	switch code {
	case "ANM":
		mets, nmts = tl*2, tl*2
	case "ASC":
		mets, nmts = tl*3, tl*2
	case "ASW":
		mets, nmts = tl*1, tl*1
	case "AUT":
		mets, nmts = tl*2, tl*2
	case "CNGD":
		mets, nmts = 0.2, 0.4
	case "ESH":
		mets, nmts = tl*10, tl*10
	case "EWP":
		mets, nmts = tl*5, tl*5
	case "FCT":
		mets, nmts = tl+8, tl+4
	case "FOOD":
		mets, nmts = 0, 0
	case "FRM":
		mets, nmts = tl+4, tl+2
	case "FUEL":
		mets, nmts = 0, 0
	case "GOLD":
		mets, nmts = 0, 0
	case "HEN":
		mets, nmts = tl*25, tl*20
	case "LAB":
		mets, nmts = tl+5, tl+3
	case "LFS":
		mets, nmts = tl*3, tl*5
	case "METS":
		mets, nmts = 0, 0
	case "MIN":
		mets, nmts = tl+5, tl+5
	case "MSL":
		mets, nmts = tl+15, tl+10
	case "MSS":
		mets, nmts = tl*2, tl*2
	case "MTBT":
		mets, nmts = tl+10, tl+10
	case "MTSP":
		mets, nmts = 0.02, 0.02
	case "NMTS":
		mets, nmts = 0, 0
	case "PWP":
		mets, nmts = tl+5, tl+5
	case "RPV":
		mets, nmts = 200/tl, 300/tl
	case "RSCH":
		mets, nmts = 0, 0
	case "SEN":
		mets, nmts = tl*1000, tl*2000
	case "SLS":
		mets, nmts = tl*0.005, tl*0.005
	case "SPD":
		mets, nmts = tl*15, tl*10
	case "STU":
		mets, nmts = tl*0.07, tl*0.03
	case "TPT":
		mets, nmts = tl*2, tl*2
	default:
		panic(fmt.Sprintf("assert(code != %q)", code))
	}
	return mets * qty, nmts * qty
}

type InventoryItem_t struct {
	Code      string
	TechLevel int64
}

var assemblyUnits = map[string]bool{
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
