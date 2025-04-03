// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import "math"

type MineGroup_t struct {
	Id      int64
	Entity  *Entity_t
	No      int64
	Deposit *Deposit_t
	Units   []*MineGroupUnit_t
	Summary struct {
		Wanted    *MineGroupInputs_t
		Allocated *MineGroupInputs_t
		Consumed  *MineGroupInputs_t
		Produced  *MineGroupOutputs_t
	}
}

type MineGroupUnit_t struct {
	Group      *MineGroup_t
	Deposit    *Deposit_t
	TechLevel  int64
	NbrOfUnits int64
	Wanted     *MineGroupInputs_t
	Allocated  *MineGroupInputs_t
	Consumed   *MineGroupInputs_t
	Produced   *MineGroupOutputs_t
}

type MineGroupInputs_t struct {
	NbrOfUnits float64
	Pro        float64
	Usk        float64
	Aut        float64
	Fuel       float64
	Ore        float64
}

type MineGroupOutputs_t struct {
	Refined float64
}

type MineGroupConstraints_t struct {
	NbrOfUnits float64
	Pro        float64
	Usk        float64
	Aut        float64
	Fuel       float64
	Ore        float64
}

func (grp *MineGroup_t) Allocate(constraints MineGroupConstraints_t) MineGroupConstraints_t {
	// allocate resources to each group unit
	for _, unit := range grp.Units {
		// constrain the number of units to the actual number of units in the group unit
		constraints.NbrOfUnits = float64(unit.NbrOfUnits)
		// allocate the resources to the group unit
		unit.Allocated = unit.Allocate(constraints)
		// consume the resources allocated to the group unit
		constraints.Pro -= unit.Allocated.Pro
		constraints.Usk -= unit.Allocated.Usk
		constraints.Aut -= unit.Allocated.Aut
		constraints.Fuel -= unit.Allocated.Fuel
		constraints.Ore -= unit.Allocated.Ore
	}
	return constraints
}

func (grp *MineGroup_t) Consume() {
	for _, unit := range grp.Units {
		unit.Consumed = unit.Consume(unit.Allocated)
	}
}

func (grp *MineGroup_t) Want() {
	for _, unit := range grp.Units {
		unit.Wanted = unit.Want()
	}
}

func (unit *MineGroupUnit_t) Allocate(constraints MineGroupConstraints_t) *MineGroupInputs_t {
	// fetch the resources required per unit in this group
	fuelPerUnit := unit.fuelRequiredPerUnit()
	proPerUnit, uskPerUnit := unit.laborRequiredPerUnit()
	orePerUnit := unit.massUnitsRequiredPerUnit()

	// calculate the maximum number of units that can be allocated based
	// on the constraints (the available amount of fuel and labor and ore)
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
		// limit the maximum number of units to the amount of ore available
		if orePerUnit*constraints.NbrOfUnits > constraints.Ore {
			constraints.NbrOfUnits = math.Floor(constraints.Fuel / orePerUnit)
			changed = true
		}
	}

	allocated := &MineGroupInputs_t{}

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

// Consume is called by the engine to consume the resources used by the group unit.
// TODO: check that the resources allocated are still available!
func (unit *MineGroupUnit_t) Consume(allocated *MineGroupInputs_t) *MineGroupInputs_t {
	// consume what we have been allocated
	return &MineGroupInputs_t{
		NbrOfUnits: allocated.NbrOfUnits,
		Pro:        allocated.Pro,
		Usk:        allocated.Usk,
		Aut:        allocated.Aut,
		Fuel:       allocated.Fuel,
		Ore:        allocated.Ore,
	}
}

// Produce calculates the outputs produced from the inputs per turn.
func (unit *MineGroupUnit_t) Produce(inputs *MineGroupInputs_t) *MineGroupOutputs_t {
	// scale the production to units per turn
	muConsumedPerUnitPerYear := float64(unit.TechLevel) * 100
	muConsumedPerYear := muConsumedPerUnitPerYear * inputs.NbrOfUnits
	muConsumedPerTurn := muConsumedPerYear / 4
	yieldPerTurn := muConsumedPerTurn * float64(unit.Deposit.Yield) / 100

	return &MineGroupOutputs_t{
		Refined: yieldPerTurn,
	}
}

// Want calculates the resources required to operate the group unit at 100% capacity.
func (unit *MineGroupUnit_t) Want() *MineGroupInputs_t {
	// fetch the resources required per unit in this group
	fuelPerUnit := unit.fuelRequiredPerUnit()
	proPerUnit, uskPerUnit := unit.laborRequiredPerUnit()
	orePerUnit := unit.massUnitsRequiredPerUnit()

	nbrOfUnits := float64(unit.NbrOfUnits)
	return &MineGroupInputs_t{
		NbrOfUnits: nbrOfUnits,
		Pro:        nbrOfUnits * proPerUnit,
		Usk:        nbrOfUnits * uskPerUnit,
		Fuel:       nbrOfUnits * fuelPerUnit,
		Ore:        nbrOfUnits * orePerUnit,
	}
}

// mines require 0.5 units of fuel per tech level
func (unit *MineGroupUnit_t) fuelRequiredPerUnit() (fuel float64) {
	return float64(unit.TechLevel) * 0.5
}

// mines require 1 PRO and 3 USK per unit per turn
func (unit *MineGroupUnit_t) laborRequiredPerUnit() (pro, usk float64) {
	return float64(unit.NbrOfUnits), float64(unit.NbrOfUnits) * 3
}

// mines consume 100 mass units per technology level per unit per year.
// We'll scale that to mass units per turn.
func (unit *MineGroupUnit_t) massUnitsRequiredPerUnit() (mu float64) {
	return float64(unit.TechLevel) * 100 / 4
}
