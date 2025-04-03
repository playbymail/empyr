// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import "math"

type FarmGroup_t struct {
	Id      int64
	Entity  *Entity_t
	No      int64
	Units   []*FarmGroupUnit_t
	Summary struct {
		Wanted    *FarmGroupInputs_t
		Allocated *FarmGroupInputs_t
		Consumed  *FarmGroupInputs_t
		Produced  *FarmGroupOutputs_t
	}
}

type FarmGroupUnit_t struct {
	Group      *FarmGroup_t
	TechLevel  int64
	NbrOfUnits int64
	Wanted     *FarmGroupInputs_t
	Allocated  *FarmGroupInputs_t
	Consumed   *FarmGroupInputs_t
	Produced   *FarmGroupOutputs_t
}

type FarmGroupInputs_t struct {
	NbrOfUnits float64
	Pro        float64
	Usk        float64
	Aut        float64
	Fuel       float64
}

type FarmGroupOutputs_t struct {
	Food float64
}

func (grp *FarmGroup_t) Allocate(constraints FarmGroupConstraints_t) FarmGroupConstraints_t {
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
	}
	return constraints
}

func (grp *FarmGroup_t) Consume() {
	for _, unit := range grp.Units {
		unit.Consumed = unit.Consume(unit.Allocated)
	}
}

func (grp *FarmGroup_t) Produce() {
	for _, unit := range grp.Units {
		unit.Produced = unit.Produce(unit.Allocated)
	}
}

func (grp *FarmGroup_t) Want() {
	for _, unit := range grp.Units {
		unit.Wanted = unit.Want()
	}
}

func (grp *FarmGroup_t) Summarize() {
	// roll the units up to the group summary
	grp.Summary.Wanted = &FarmGroupInputs_t{}
	grp.Summary.Allocated = &FarmGroupInputs_t{}
	grp.Summary.Consumed = &FarmGroupInputs_t{}
	grp.Summary.Produced = &FarmGroupOutputs_t{}
	for _, unit := range grp.Units {
		grp.Summary.Wanted.Fuel += unit.Wanted.Fuel
		grp.Summary.Wanted.Pro += unit.Wanted.Pro
		grp.Summary.Wanted.Usk += unit.Wanted.Usk
		grp.Summary.Allocated.Fuel += unit.Allocated.Fuel
		grp.Summary.Allocated.Pro += unit.Allocated.Pro
		grp.Summary.Allocated.Usk += unit.Allocated.Usk
		grp.Summary.Consumed.Fuel += unit.Consumed.Fuel
		grp.Summary.Consumed.Pro += unit.Consumed.Pro
		grp.Summary.Consumed.Usk += unit.Consumed.Usk
		grp.Summary.Produced.Food += unit.Produced.Food
	}
}

type FarmGroupConstraints_t struct {
	NbrOfUnits float64
	Pro        float64
	Usk        float64
	Fuel       float64
	Aut        float64
}

func (unit *FarmGroupUnit_t) Allocate(constraints FarmGroupConstraints_t) *FarmGroupInputs_t {
	// fetch the resources required per unit in this group
	fuelPerUnit := unit.fuelRequiredPerUnit()
	proPerUnit, uskPerUnit := unit.laborRequiredPerUnit()

	// calculate the maximum number of units that can be allocated based
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

	allocated := &FarmGroupInputs_t{}

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
func (unit *FarmGroupUnit_t) Consume(allocated *FarmGroupInputs_t) *FarmGroupInputs_t {
	// consume what we have been allocated
	return &FarmGroupInputs_t{
		NbrOfUnits: allocated.NbrOfUnits,
		Pro:        allocated.Pro,
		Usk:        allocated.Usk,
		Aut:        allocated.Aut,
		Fuel:       allocated.Fuel,
	}
}

// Produce calculates the outputs produced from the inputs per turn.
func (unit *FarmGroupUnit_t) Produce(inputs *FarmGroupInputs_t) *FarmGroupOutputs_t {
	// scale the production to units per turn
	var qtyPerUnitPerYear float64
	if unit.TechLevel == 1 {
		qtyPerUnitPerYear = 100
	} else {
		qtyPerUnitPerYear = float64(unit.TechLevel) * 20
	}
	qtyPerYear := qtyPerUnitPerYear * inputs.NbrOfUnits
	qtyPerTurn := qtyPerYear / 4

	return &FarmGroupOutputs_t{
		Food: math.Floor(qtyPerTurn),
	}
}

// Want calculates the resources required to operate the group unit at 100% capacity.
func (unit *FarmGroupUnit_t) Want() *FarmGroupInputs_t {
	// fetch the resources required per unit in this group
	fuelPerUnit := unit.fuelRequiredPerUnit()
	proPerUnit, uskPerUnit := unit.laborRequiredPerUnit()

	nbrOfUnits := float64(unit.NbrOfUnits)
	return &FarmGroupInputs_t{
		NbrOfUnits: nbrOfUnits,
		Pro:        nbrOfUnits * proPerUnit,
		Usk:        nbrOfUnits * uskPerUnit,
		Fuel:       nbrOfUnits * fuelPerUnit,
	}
}

// farms require a variable amount of fuel to operate
func (unit *FarmGroupUnit_t) fuelRequiredPerUnit() (fuel float64) {
	return farmFuel(unit.Group.Entity, unit.TechLevel, 1)
}

// farms require 1 PRO and 3 USK per unit per turn
func (unit *FarmGroupUnit_t) laborRequiredPerUnit() (pro, usk float64) {
	return float64(unit.NbrOfUnits), float64(unit.NbrOfUnits) * 3
}
