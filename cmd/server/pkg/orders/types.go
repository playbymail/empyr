// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package orders

// this file defines all of the valid order types.

// Accept order...
type Accept struct {
	PlayerID string `json:"player_id"`
	TargetID string `json:"target_id"`
}

// AddOn Order...
type AddOn struct {
	SourceID      string `json:"source_id"` // id of ship or colony being ordered
	TargetID      string `json:"target_id"` // id of ship or colony being targeted
	Item          string `json:"item"`
	TechLevel     int    `json:"tech_level"`
	Quantity      int    `json:"quantity"`
	DoNotAssemble bool   `json:"do_not_assemble,omitempty"` // if true, leave unassembled in storage
}

// AfterManeuverEnergyWeaponFire order...
type AfterManeuverEnergyWeaponFire struct {
	SourceID                string  `json:"source_id"`                           // id of ship or colony being ordered
	TargetID                string  `json:"target_id"`                           // id of ship or colony being targeted
	Percentage              float64 `json:"percentage"`                          // percentage of weapons to allocate to order
	TargetCategory          string  `json:"target_category,omitempty"`           // only if target is a colony
	MaximumTacticalDistance int     `json:"maximum_tactical_distance,omitempty"` // do not activate if maximum distance to target is exceeded
}

// AfterManeuverMissileFire order
type AfterManeuverMissileFire struct {
	SourceID                string  `json:"source_id"`                           // id of ship or colony being ordered
	TargetID                string  `json:"target_id"`                           // id of ship or colony being targeted
	Percentage              float64 `json:"percentage"`                          // percentage of missile launchers to allocate to order
	TargetCategory          string  `json:"target_category,omitempty"`           // only if target is a colony
	MaximumTacticalDistance int     `json:"maximum_tactical_distance,omitempty"` // do not activate if maximum distance to target is exceeded
}

// AssembleFactory order...
type AssembleFactory struct {
	SourceID  string `json:"source_id"`  // id of ship or colony being ordered
	Quantity  int    `json:"quantity"`   // number of factory units to assemble
	Item      string `json:"item"`       // id of items to manufacture
	TechLevel int    `json:"tech_level"` // tech level of items to manufacture
}

// AssembleFactoryGroup order...
type AssembleFactoryGroup struct {
	SourceID string `json:"source_id"` // id of ship or colony being ordered
	Quantity int    `json:"quantity"`  // number of factory units to assemble
	GroupID  string `json:"group_id"`  // id of group to add factories to
}

// AssembleItem order...
type AssembleItem struct {
	SourceID  string `json:"source_id"` // id of ship or colony being ordered
	Quantity  int    `json:"quantity"`
	Item      string `json:"item"`
	TechLevel int    `json:"tech_level"`
}

// AssembleMine order...
type AssembleMine struct {
	SourceID  string `json:"source_id"`  // id of ship or colony being ordered
	Quantity  int    `json:"quantity"`   // number of mine units to assemble
	TechLevel int    `json:"tech_level"` // tech level of mines to assemble
}

// AssembleMineGroup order...
type AssembleMineGroup struct {
	SourceID  string `json:"source_id"`  // id of ship or colony being ordered
	Quantity  int    `json:"quantity"`   // number of mine units to assemble
	DepositID string `json:"deposit_id"` // id of deposit to add mines to
}

// AutoReturnFire order...
type AutoReturnFire struct {
	SourceID   string  `json:"source_id"` // id of ship or colony being ordered
	Percentage float64 `json:"percentage"`
}

// BuildChange order...
type BuildChange struct {
	SourceID  string `json:"source_id"`  // id of ship or colony being ordered
	GroupID   string `json:"group_id"`   // id of factory group to change
	Item      string `json:"item"`       // id of item to manufacture
	TechLevel int    `json:"tech_level"` // tech level of item to manufacture
}

// Close order...
type Close struct {
	ShipID           string `json:"ship_id"`                     // id of ship being ordered
	TargetID         string `json:"target_id"`                   // id of ship or colony to close upon
	StandoffDistance int    `json:"standoff_distance,omitempty"` // optional distinance to stand off from target
}

// CloseProximityTargeting order...
type CloseProximityTargeting struct {
	SourceID   string  `json:"source_id"`  // id of ship or colony being ordered
	Percentage float64 `jons:"percentage"` // percentage of weapons to allocate to order
}

// CombineFactoryGroup order...
type CombineFactoryGroup struct {
	SourceID    string `json:"source_id"` // id of ship or colony being ordered
	FromGroupID string `json:"from_group_id"`
	ToGroupID   string `json:"to_group_id"`
	WIPOnly     bool   `json:"wip_only,omitempty"`
	WIPQuarters []int  `json:"wip_quarters,omitempty"` // quarters of WIP to combine
}

// ControlPlanet order...
type ControlPlanet struct {
	ColonyID string `json:"colony_id"` // id of surface colony being ordered
}

// Coordinates for Tactical or Stellar locations
type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

// Debug order sets the debug flag to on or off.
type Debug struct {
	On bool `json:"on"`
}

// DefineCargoHold order...
type DefineCargoHold struct {
	ShipID   string `json:"ship_id"` // id of ship being ordered
	Quantity int    `json:"quantity"`
}

// DefensiveSupport order...
type DefensiveSupport struct {
	SourceID string `json:"source_id"` // id of ship or colony being ordered
	TargetID string `json:"target_id"` // id of ship or colony to provide support to
	Items    []struct {
		Item      string `json:"item"`
		TechLevel int    `json:"tech_level"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
}

// Disassemble order...
type Disassemble struct {
	SourceID  string `json:"source_id"` // id of ship or colony being ordered
	Item      string `json:"item"`
	TechLevel int    `json:"tech_level"`
	GroupID   string `json:"group_id,omitempty"` // required for factories and mines
	Quantity  int    `json:"quantity"`
}

// Disband order...
type Disband struct {
	SourceID       string `json:"source_id"`       // id of ship or colony being ordered
	RaceID         string `json:"race_id"`         // id of race to draft from
	PopulationType string `json:"population_type"` // type of population to draft
	Quantity       int    `json:"quantity"`
}

// Dock order...
type Dock struct {
	ShipID   string `json:"ship_id"`   // id of ship being ordered
	TargetID string `json:"target_id"` // id of S/C being docked with
}

// Dodge order...
type Dodge struct {
	ShipID     string  `json:"ship_id"`    // id of ship being ordered
	Percentage float64 `json:"percentage"` // percentage of speed to allocate to maneuver
}

// Draft order...
type Draft struct {
	SourceID       string `json:"source_id"`       // id of ship or colony being ordered
	RaceID         string `json:"race_id"`         // id of race to draft from
	PopulationType string `json:"population_type"` // type of population to draft
	Quantity       int    `json:"quantity"`
}

// ExpendCommittedBufferResearchPoints order...
type ExpendCommittedBufferResearchPoints struct {
	ColonyID string `json:"colony_id"` // id of colony being ordered
	Quantity int    `json:"quantity"`  // amount of research points to expend
	Item     string `json:"item"`      // item to apply the research points to
}

// ExpendPrototype order...
type ExpendPrototype struct {
	ColonyID  string `json:"colony_id"`  // id of colony being ordered
	Quantity  int    `json:"quantity"`   // number of prototypes to expend
	Item      string `json:"item"`       // id of prototype item to consume
	TechLevel string `json:"tech_level"` // tech level of prototype item to consume
}

// ExpendResearchPointsOnly order...
type ExpendResearchPointsOnly struct {
	ColonyID string `json:"colony_id"` // id of colony being ordered
	Quantity int    `json:"quantity"`  // amount of research points to expend
	Item     string `json:"item"`      // item to apply the research points to
}

// FactoryGroupChange order...
type FactoryGroupChange struct {
	ColonyID string `json:"colony_id"` // id of colony being ordered
	FromID   string `json:"from_id"`   // id of factory group to move from
	ToID     string `json:"to_id"`     // id of factory group to move to
	Quantity int    `json:"quantity"`  // number of factory units to move
}

// Give order...
type Give struct {
	SourceID string `json:"source_id"` // id of ship or colony being ordered
	Type     string `json:"type"`      // must be one of ALLY, RULER, VICEROY
	TargetID string `json:"target_id"` // id of the entity being given to
}

// HomePortChange order...
type HomePortChange struct {
	ShipID   string `json:"ship_id"`   // id of ship being ordered
	ColonyID string `json:"colony_id"` // id of colony being targeted
}

// Invade order...
type Invade struct {
	SourceID string `json:"source_id"` // id of ship or colony being ordered
	TargetID string `json:"target_id"` // id of ship or colony being targeted
	Items    []struct {
		Item      string `json:"item"`
		TechLevel int    `json:"tech_level"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
}

// Jump order...
type Jump struct {
	ShipID string `json:"ship_id"` // id of ship being ordered
	Coords Coords `json:"coords"`  // coordinates of destination
	Offset int    `json:"offset"`  // distance (in tactical distance units) to arrive from destination
}

// Junk order...
type Junk struct {
	SourceID string `json:"source_id"` // id of ship or colony being ordered
	TargetID string `json:"target_id"` // id of ship or colony being targeted
}

// LaunchRobotProbe order...
// Type must be one of "SURVEY", "ORBIT", "SYSTEM", "SHIP", or "COLONY"
type LaunchRobotProbe struct {
	SourceID   string `json:"source_id"`   // id of ship or colony being ordered
	Type       string `json:"type"`        // type of probe to launch
	Coords     Coords `json:"coords"`      // coordinates (only for orbit and system probes)
	StarLetter string `json:"star_letter"` // only for multiple star systems
	Orbit      int    `json:"orbit"`       // number of orbit to probe
}

// LoadCargo order...
type LoadCargo struct {
	ColonyID  string `json:"colony_id"` // id of colony being ordered
	ToID      string `json:"to_id"`     // id of ship to load
	Item      string `json:"item"`
	TechLevel int    `json:"tech_level"`
	Quantity  int    `json:"quantity"`
}

// Merge order...
type Merge struct {
	SourceID string `json:"source_id"` // id of ship or colony being ordered
	TargetID string `json:"target_id"` // id of ship or colony being targeted
}

// Message order...
// Text must be UTF-8 and is truncated at 200 runes.
type Message struct {
	SourceID string `json:"source_id"` // id of ship or colony being ordered
	TargetID string `json:"target_id"` // id of the ship or colony to deliver the message to
	Text     string `json:"text"`      // text to be sent
}

// MineChange order...
type MineChange struct {
	SourceID  string `json:"source_id"`  // id of ship or colony being ordered
	GroupID   string `json:"group_id"`   // id of mine group to modify
	DepositID string `json:"deposit_id"` // id of deposit to move mine units to
	Quantity  int    `json:"quantity"`   // number of mine unites to move to deposit
}

// MineShutDown order...
type MineShutDown struct {
	SourceID string `json:"source_id"` // id of ship or colony being ordered
	GroupID  string `json:"group_id"`  // id of mine group to shut down
	Quantity int    `json:"quantity"`  // number of mine unites to shut down
}

// MineStartUp order...
type MineStartUp struct {
	SourceID string `json:"source_id"` // id of ship or colony being ordered
	GroupID  string `json:"group_id"`  // id of mine group to start up
	Quantity int    `json:"quantity"`  // number of mine unites to start up
}

// Move order...
// Movement within a system.
type Move struct {
	ShipID string `json:"ship_id"` // id of ship being ordered
	Orbit  int    `json:"orbit"`   // number of orbit to move to (must be 1..10)
	Offset int    `json:"offset"`  // distance (in tactical distance units) to arrive from destination
}

// Name order...
// Note: original order used "Entity" instead of "Type".
// Entity must be one of SHIP, COLONY, SYSTEM, STAR, PLANET, or PLAYER
type Name struct {
	EntityID string `json:"entity_id"` // id of entity being ordered
	Type     string `json:"type"`      // type of entity to assign name to
	Name     string `json:"name"`      // name to assign to the entity
}

// Note order...
// Text must be UTF-8 and is truncated at 200 runes.
type Note struct {
	SourceID string `json:"source_id"` // id of ship or colony being ordered
	Text     string `json:"text"`      // text to be displayed on Owner's report for the ship or colony
}

// OffensiveSupport order...
type OffensiveSupport struct {
	SourceID string `json:"source_id"` // id of ship or colony being ordered
	TargetID string `json:"target_id"` // id of ship or colony to provide support against
	Items    []struct {
		Item      string `json:"item"`
		TechLevel int    `json:"tech_level"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
}

// Pay order...
// If PopulationType is ALL, then Amount is a percentage of current pay rates to pay out.
type Pay struct {
	ColonyID       string  `json:"colony_id"`       // id of colony being ordered
	Amount         float64 `json:"amount"`          // amount (in Consumer Goods) to pay per unit
	PopulationType string  `json:"population_type"` // type of population to pay
}

// PermissionToColonize order...
type PermissionToColonize struct {
	ColonyID string `json:"colony_id"` // id of colony being ordered
	ShipID   string `json:"ship_id"`   // id of ship being granted permission
}

// PickUpItem order...
type PickUpItem struct {
	SourceID  string `json:"source_id"` // id of ship or colony being ordered
	ToID      string `json:"to_id"`     // id of ship to transfer items to
	Item      string `json:"item"`
	TechLevel int    `json:"tech_level"`
	Quantity  int    `json:"quantity"`
}

// PickUpPopulation order...
type PickUpPopulation struct {
	SourceID       string `json:"source_id"` // id of ship or colony being ordered
	ToID           string `json:"to_id"`     // id of ship to transfer population to
	PopulationType string `json:"population_type"`
	RaceID         string `json:"race_id"`
	Quantity       int    `json:"quantity"`
}

// PreManeuverEnergyWeaponFire order...
type PreManeuverEnergyWeaponFire struct {
	SourceID                string  `json:"source_id"`                         // id of ship or colony being ordered
	TargetID                string  `json:"target_id"`                         // id of ship or colony being targeted
	Percentage              float64 `json:"percentage"`                        // percentage of weapons to allocate to order
	TargetCategory          string  `json:"targetCategory,omitempty"`          // only if target is a colony
	MaximumTacticalDistance int     `json:"maximumTacticalDistance,omitempty"` // do not activate if maximum distance to target is exceeded
}

// PreManeuverMissileFire order
type PreManeuverMissileFire struct {
	SourceID                string  `json:"source_id"`                           // id of ship or colony being ordered
	TargetID                string  `json:"target_id"`                           // id of ship or colony being targeted
	Percentage              float64 `json:"percentage"`                          // percentage of missile launchers to allocate to order
	TargetCategory          string  `json:"target_category,omitempty"`           // only if target is a colony
	MaximumTacticalDistance int     `json:"maximum_tactical_distance,omitempty"` // do not activate if maximum distance to target is exceeded
}

// Probe order...
// Also called "Intensive Probe."
// TODO: split to ProbeColony and ProbeShip
type Probe struct {
	SourceID string `json:"source_id"` // id of ship or colony being ordered
	TargetID string `json:"target_id"` // id of ship or colony to probe
}

// ProbeOrbit order...
// If Orbit is zero (0) then probe all orbits.
type ProbeOrbit struct {
	SourceID string `json:"source_id"`       // id of ship or colony being ordered
	TargetID string `json:"target_id"`       // id of system to probe
	Orbit    int    `json:"orbit,omitempty"` // number of orbit to probe (must be 0..10)
}

// ProbeSystem order...
// TODO: Consider a ProbeCoordinates order to provide some symmetry
type ProbeSystem struct {
	SourceID  string `json:"source_id"` // id of ship or colony being ordered
	TargetID  string `json:"target_id"` // id of system to center probe on
	Magnitude int    `json:"magnitude"` // radius (in light years) to probe from system's central coordinates
}

// Ration order...
type Ration struct {
	SourceID string  `json:"source_id"` // id of ship or colony being ordered
	Amount   float64 `json:"amount"`    // percentage of standard rations
}

// Run order...
type Run struct {
	ShipID   string `json:"ship_id"`   // id of ship being ordered
	TargetID string `json:"target_id"` // id of S/C to run from
}

// Scrap order...
type Scrap struct {
	SourceID  string `json:"source_id"` // id of ship or colony being ordered
	Item      string `json:"item"`
	TechLevel int    `json:"tech_level"`
	Quantity  int    `json:"quantity"`
}

// SetUp order...
type SetUp struct {
	SourceID     string `json:"source_id"` // id of ship or colony being ordered
	TypeOfColony string `json:"type_of_colony"`
	Quantity     int    `json:"quantity"`
	Items        []struct {
		Factory *struct {
			Quantity      int    `json:"quantity"`
			ItemToBuild   string `json:"item_to_build"`
			ItemTechLevel int    `json:"item_tech_level"`
		} `json:"factory,omitempty"`
		Item *struct {
			Quantity      int    `json:"quantity"`
			Item          string `json:"item"`
			TechLevel     int    `json:"tech_level"`
			DoNotAssemble bool   `json:"do_not_assemble,omitempty"`
		} `json:"item,omitempty"`
		Mine *struct {
			Quantity    int    `json:"quantity"`
			DepositID   string `json:"deposit_id,omitempty"`
			DepositType string `json:"deposit_type,omitempty"`
		} `json:"mine,omitempty"`
	} `json:"items"`
}

// ShutDown order...
// The item to shut down must be FRM or LAB. Any other value will be rejected.
// TODO: consider renaming to "FarmShutDown" and "LabShutDown"
type ShutDown struct {
	SourceID  string `json:"source_id"`  // id of ship or colony being ordered
	ItemID    string `json:"item_id"`    // id of farm or lab to shut down (must be FRM or LAB!)
	TechLevel int    `json:"tech_level"` // tech level of farm or lab to shut down
	Quantity  int    `json:"quantity"`   // number of mine unites to shut down
}

// StartUp order...
// The item to shut down must be FRM or LAB. Any other value will be rejected.
// TODO: consider renaming to "FarmStartUp" and "LabStartUp"
type StartUp struct {
	SourceID  string `json:"source_id"`  // id of ship or colony being ordered
	ItemID    string `json:"item_id"`    // id of farm or lab to start up (must be FRM or LAB!)
	TechLevel int    `json:"tech_level"` // tech level of item to start up
	Quantity  int    `json:"quantity"`   // number of mine unites to shut down
}

// Survey order...
type Survey struct {
	SourceID string `json:"source_id"` // id of ship or colony being ordered
	PlanetID string `json:"planet_id"` // id of planet to survey
}

// TacticalManeuver order...
type TacticalManeuver struct {
	ShipID string `json:"ship_id"` // id of ship being ordered
	To     Coords `json:"to"`      // targeted coordinates
}

// Transfer order...
type Transfer struct {
	SourceID  string `json:"source_id"` // id of ship or colony being ordered
	ToID      string `json:"ship_id"`   // id of ship to transfer items to
	Item      string `json:"item"`
	TechLevel int    `json:"tech_level"`
	Quantity  int    `json:"quantity"`
}

// UncontrolPlanet order...
type UncontrolPlanet struct {
	ColonyID string `json:"colony_id"` // id of surface colony being ordered
}

// Undock order...
type Undock struct {
	ShipID string `json:"ship_id"` // id of ship being ordered
}

// UnloadCargo order...
type UnloadCargo struct {
	ColonyID  string `json:"colony_id"` // id of colony being ordered
	ShipID    string `json:"ship_id"`   // id of ship to unload
	Item      string `json:"item"`
	TechLevel int    `json:"tech_level"`
	Quantity  int    `json:"quantity"`
}

// Withdraw order...
// Note: manual calls the target ID the "defending ID".
type Withdraw struct {
	SourceID string `json:"source_id"` // id of ship or colony being ordered
	TargetID string `json:"target_id"` // id of ship or colony being targeted
}
