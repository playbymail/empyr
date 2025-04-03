// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import "math"

// A FactoryGroup_t is a group of factories working together to produce a
// single type of unit.
//
// Factory groups have a cache, which is the raw materials used in the
// production process. The cache allows us to move materials from the colony
// to the factory and still have them count as part of the colony's mass.
// If we didn't, the colony would lose and gain mass in sync with the
// production cycle.
type FactoryGroup_t struct {
	Id      int64
	Entity  *Entity_t
	No      int64
	Tooling struct {
		Current *FactoryGroupTooling_t
		Retool  *FactoryGroupTooling_t
	}
	Units   []*FactoryGroupUnit_t
	Summary struct {
		Wanted    *FactoryGroupInputs_t
		Allocated *FactoryGroupInputs_t
		Consumed  *FactoryGroupInputs_t
		Produced  [3]float64
	}
}

// A FactoryGroupTooling_t tells the factory group which item to produce.
// If the group has work in progress, it will stop producing the current
// item, push the remaining work through the pipeline, and then start
// producing the new item.
type FactoryGroupTooling_t struct {
	Group     *FactoryGroup_t
	TurnNo    int64 // turn number when the tooling was changed
	Item      *Unit_t
	TechLevel int64
}

// A FactoryGroupUnit_t contains all the units in a group that
// have the same technology level.
type FactoryGroupUnit_t struct {
	Group      *FactoryGroup_t
	TechLevel  int64
	NbrOfUnits int64
	Summary    struct {
		Wanted    *FactoryGroupInputs_t
		Allocated *FactoryGroupInputs_t
		Consumed  *FactoryGroupInputs_t
		Produced  [3]float64
	}
	Pipeline [3]*FactoryGroupPipeline_t
}

// A FactoryGroupPipeline_t is a pipeline of factories that have the same tooling
// and technology level. At the beginning of each production cycle, the pipeline
// has queued up work to be done. These units are the work in progress remaining
// from the previous production cycle and are backlogged because of resource
// constraints.
type FactoryGroupPipeline_t struct {
	GroupUnit  *FactoryGroupUnit_t
	TechLevel  int64
	NbrOfUnits int64
	Wanted     *FactoryGroupInputs_t
	Allocated  *FactoryGroupInputs_t
	Consumed   *FactoryGroupInputs_t
	Backlog    float64 // work not completed in the previous production cycle
	Produced   float64
}

type FactoryGroupInputs_t struct {
	NbrOfUnits float64
	Pro        float64
	Usk        float64
	Aut        float64
	Gold       float64
	Fuel       float64
	Mets       float64
	Nmts       float64
}

type FactoryGroupOutputs_t struct {
	WIP [3]float64
}

type FactoryGroupConstraints_t struct {
	NbrOfUnits float64
	Pro        float64
	Usk        float64
	Aut        float64
	Gold       float64
	Fuel       float64
	Mets       float64
	Nmts       float64
}

// Allocate allocates resources to the factory group unit.
//
// The function accepts constraints, which are the maximum resources that can be allocated to the factory group unit. It returns the resources that are allocated.
//
// Create and initialize the allocated resources to zero.
//
// If the factory group unit has no current tooling, return immediately.
//
// If the factory group is retooling, it must return all cached METS and NMTS resources.
//
// If the factory group is retooling, it will only allocate FUEL and labor.
//
// If the factory group unit has current tooling, it claims METS and NMTS resources by moving them from the constraint to cached. The maximum METS and NMTS that can be cached depend on the item being built and the number of factory units in the group unit.
//
// Determine the number of items the factory group unit can produce in a year, assuming 100% allocation of resources. Divide that by 4 to get the maximum number of items the group unit can produce this turn.
//
// Allocate FUEL and labor to the work in the pipeline, starting with the 75% complete group, then the 50% complete group, and then the 25% complete group. All FUEL and labor that is allocated must be removed from the constraint.
//
// Using the remaining FUEL and labor in the constraint along with the cached METS and NMETS to determine the number of new items that will be started this turn. Remove any allocated FUEL and labor from the constraint.
func (unit *FactoryGroupUnit_t) Allocate(constraints FactoryGroupConstraints_t) *FactoryGroupInputs_t {
	// Create and initialize allocated resources to zero
	allocated := &FactoryGroupInputs_t{}
	isRetooling := unit.Group.Tooling.Retool != nil
	isCaching := !isRetooling && unit.Group.Tooling.Current != nil

	// Initialize cached resources if not already done
	if unit.Cached == nil {
		unit.Cached = &FactoryGroupInputs_t{}
	}

	// If the factory group unit has no current tooling,
	// clear the cache and return without allocating any resources
	if unit.Group.Tooling.Current == nil {
		constraints.Mets += unit.Cached.Mets
		constraints.Nmts += unit.Cached.Nmts
		unit.Cached.Mets = 0
		unit.Cached.Nmts = 0
		return allocated
	}

	// Calculate the resources required per unit in this group.
	fuelPerUnit := unit.fuelRequiredPerUnit()
	proPerUnit, uskPerUnit := unit.laborRequiredPerUnit()

	// Get the number of units in this group
	nbrOfUnits := float64(unit.NbrOfUnits)

	// If the factory group is retooling, return all cached METS and NMTS resources
	if isRetooling {
		// Return cached resources to constraints
		constraints.Mets += unit.Cached.Mets
		constraints.Nmts += unit.Cached.Nmts
		unit.Cached.Mets = 0
		unit.Cached.Nmts = 0
	} else if isCaching { // the factory group will try to cache a year's worth of input materials
		// calculate the resources required to produce a single item
		item, techLevel := unit.Group.Tooling.Current.Item, unit.Group.Tooling.Current.TechLevel
		metsPerItem, nmtsPerItem := unitRequirements(item.Code, techLevel, 1)
		muPerUnit := metsPerItem + nmtsPerItem

		// calculate resources needed for full year's production.
		maxMUPerYear := float64(unit.TechLevel) * 20 * nbrOfUnits
		maxItemsPerYear := maxMUPerYear / muPerUnit
		maxMetsPerYear := metsPerItem * maxItemsPerYear
		maxNmtsPerYear := nmtsPerItem * maxItemsPerYear

		// cache up to a year's METS
		if deltaMets := maxMetsPerYear - unit.Cached.Mets; deltaMets > 0 {
			deltaMets = math.Min(deltaMets, constraints.Mets)
			constraints.Mets -= deltaMets
			unit.Cached.Mets += deltaMets
			allocated.Mets += deltaMets
		}
		// cache up to a year's NMTS
		if deltaNmts := maxNmtsPerYear - unit.Cached.Nmts; deltaNmts > 0 {
			deltaNmts = math.Min(deltaNmts, constraints.Nmts)
			constraints.Nmts -= deltaNmts
			unit.Cached.Nmts += deltaNmts
			allocated.Nmts += deltaNmts
		}
	}

	// Calculate the maximum number of units that can be allocated based
	// on the constraints (the available amount of fuel and labor)
	for changed := true; changed && constraints.NbrOfUnits > 0; changed = false {
		// limit the maximum number of units to the amount of fuel available
		if fuelPerUnit*constraints.NbrOfUnits > constraints.Fuel {
			constraints.NbrOfUnits = math.Floor(constraints.Fuel / fuelPerUnit)
			changed = true
		}
		// limit the maximum number of units to the amount of PRO available
		if proPerUnit*constraints.NbrOfUnits > constraints.Pro {
			constraints.NbrOfUnits = math.Floor(constraints.Pro / proPerUnit)
			changed = true
		}
		// limit the maximum number of units to the amount of USK and AUT available
		if uskPerUnit*constraints.NbrOfUnits > constraints.Usk+constraints.Aut {
			constraints.NbrOfUnits = math.Floor((constraints.Usk + constraints.Aut) / uskPerUnit)
			changed = true
		}
	}

	// sanity check
	if constraints.NbrOfUnits < 1 {
		return allocated
	}

	// allocate resources based on the maximum number of units that satisfy the constraints
	allocated.NbrOfUnits = constraints.NbrOfUnits
	allocated.Pro = constraints.NbrOfUnits * proPerUnit
	// always allocate AUT before USK
	allocated.Aut = math.Min(constraints.Aut, constraints.NbrOfUnits*uskPerUnit)
	if allocated.Aut < constraints.NbrOfUnits*uskPerUnit {
		allocated.Usk = constraints.NbrOfUnits*uskPerUnit - allocated.Aut
	}
	allocated.Fuel = constraints.NbrOfUnits * fuelPerUnit

	return allocated
}

// Produce calculates the outputs produced from the inputs per turn.
func (unit *FactoryGroupUnit_t) Produce(inputs *FactoryGroupInputs_t) *FactoryGroupOutputs_t {
	isRetooling := unit.Group.Tooling.Retool != nil
	isCaching := !isRetooling && unit.Group.Tooling.Current != nil

	output := &FactoryGroupOutputs_t{}

	// the group unit cannot produce new items while retooling or there is no current tooling
	if isRetooling {
		// move items in the pipeline
		return output
	}
	tooling := unit.Group.Tooling.Current
	if tooling == nil { // should never happen
		return output
	}

	// requirements for the current tooling
	metsPerUnit, nmtsPerUnit := unit.metsAndNmtsPerUnit()
	muPerUnit := metsPerUnit + nmtsPerUnit

	// calculate this group unit's maximum production capacity in units
	muConsumedPerYear := float64(unit.TechLevel) * 20 * inputs.NbrOfUnits
	maxProductionCapacity := muConsumedPerYear / muPerUnit

	// does the group have any capacity to produce new units?
	totalWIP := unit.WIP[0] + unit.WIP[1] + unit.WIP[2]
	remainingCapacity := maxProductionCapacity - totalWIP
	if remainingCapacity < 1 {
		// TODO: if maximum production capacity exceeds the WIP, we should destroy units in WIP
		return output
	}

	// for each stage in the pipeline, use cached resources to move items to the next stage

	// how many units can we produce this turn?
	productionCapacity := math.Min(maxProductionCapacity, remainingCapacity) / 4

	// scale the production to units per turn
	muConsumedPerTurn := muConsumedPerYear / 4
	yieldPerTurn := muConsumedPerTurn * float64(unit.Deposit.Yield) / 100

	return &FactoryGroupOutputs_t{
		Refined: yieldPerTurn,
	}
}

// Want calculates the resources required to operate the group unit at 100% capacity.
func (unit *FactoryGroupUnit_t) Want() *FactoryGroupInputs_t {
	// fetch the resources required per unit in this group
	fuelPerUnit := unit.fuelRequiredPerUnit()
	proPerUnit, uskPerUnit := unit.laborRequiredPerUnit()

	nbrOfUnits := float64(unit.NbrOfUnits)
	return &FactoryGroupInputs_t{
		NbrOfUnits: nbrOfUnits,
		Pro:        nbrOfUnits * proPerUnit,
		Usk:        nbrOfUnits * uskPerUnit,
		Fuel:       nbrOfUnits * fuelPerUnit,
	}
}

// factories require 0.5 units of fuel per tech level
func (unit *FactoryGroupUnit_t) fuelRequiredPerUnit() (fuel float64) {
	return float64(unit.TechLevel) * 0.5
}

// factories require a variable number of PRO and USK per unit per turn
func (unit *FactoryGroupUnit_t) laborRequiredPerUnit() (pro, usk float64) {
	if unit.NbrOfUnits < 5 {
		return 6, 18
	} else if unit.NbrOfUnits < 50 {
		return 5, 15
	} else if unit.NbrOfUnits < 500 {
		return 4, 12
	} else if unit.NbrOfUnits < 5_000 {
		return 3, 9
	} else if unit.NbrOfUnits < 50_000 {
		return 2, 6
	} else {
		return 1, 3
	}
}

// factories consume 30 mass units per technology level per unit per year.
// We'll scale that to mass units per turn.
func (unit *FactoryGroupUnit_t) massUnitsRequiredPerUnit() (mu float64) {
	return float64(unit.TechLevel) * 20 / 4
}

func (unit *FactoryGroupUnit_t) metsAndNmtsPerUnit() (mets, nmts float64) {
	return unitRequirements(unit.Group.Tooling.Current.Item.Code, unit.Group.Tooling.Current.TechLevel, 1)
}
