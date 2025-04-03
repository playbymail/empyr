// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

// Entity_t is either a ship or a colony.
type Entity_t struct {
	Id          int64
	IsColony    bool
	IsEnclosed  bool
	IsOnSurface bool
	TechLevel   int64
	Location    *Orbit_t
	Population  struct {
		UEM PopulationGroup_t
		USK PopulationGroup_t
		PRO PopulationGroup_t
		SLD PopulationGroup_t
		CNW PopulationGroup_t
		SPY PopulationGroup_t
	}
	Inventory     []*Inventory_t
	FactoryGroups []*FactoryGroup_t
	FarmGroups    []*FarmGroup_t
	MiningGroups  []*MineGroup_t
}

type Inventory_t struct {
	Unit        *Unit_t
	TechLevel   int64
	Qty         int64
	Mass        float64
	Volume      float64
	IsAssembled bool
}

type PopulationGroup_t struct {
	Qty      int64
	Pay      float64
	RebelQty int64
}
