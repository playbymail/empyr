// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package orders

// Order is a hot mess, but it allows our JSON input to be simpler
// to read and parse. The consumer of an Order must test all the
// properties for non-nil to determine which one to process.
type Order struct {
	priority                            int                                  // priority for sorting orders
	Accept                              *Accept                              `json:"accept,omitempty"`
	AddOn                               *AddOn                               `json:"add_on,omitempty"`
	AfterManeuverEnergyWeaponFire       *AfterManeuverEnergyWeaponFire       `json:"after_maneuver_energy_weapon_fire,omitempty"`
	AfterManeuverMissileFire            *AfterManeuverMissileFire            `json:"after_maneuver_missile_fire,omitempty"`
	AssembleFactory                     *AssembleFactory                     `json:"assemble_factory,omitempty"`
	AssembleFactoryGroup                *AssembleFactoryGroup                `json:"assemble_factory_group,omitempty"`
	AssembleItem                        *AssembleItem                        `json:"assemble_item,omitempty"`
	AssembleMine                        *AssembleMine                        `json:"assemble_mine,omitempty"`
	AssembleMineGroup                   *AssembleMineGroup                   `json:"assemble_mine_group,omitempty"`
	AutoReturnFire                      *AutoReturnFire                      `json:"auto_return_fire,omitempty"`
	BuildChange                         *BuildChange                         `json:"build_change,omitempty"`
	Close                               *Close                               `json:"close,omitempty"`
	CloseProximityTargeting             *CloseProximityTargeting             `json:"close_proximity_targeting,omitempty"`
	CombineFactoryGroup                 *CombineFactoryGroup                 `json:"combine_factory_group,omitempty"`
	ControlPlanet                       *ControlPlanet                       `json:"control_planet,omitempty"`
	Debug                               *Debug                               `json:"debug,omitempty"`
	DefensiveSupport                    *DefensiveSupport                    `json:"defensive_support,omitempty"`
	DefineCargoHold                     *DefineCargoHold                     `json:"define_cargo_hold,omitempty"`
	Disassemble                         *Disassemble                         `json:"disassemble,omitempty"`
	Disband                             *Disband                             `json:"disband,omitempty"`
	Dock                                *Dock                                `json:"dock,omitempty"`
	Dodge                               *Dodge                               `json:"dodge,omitempty"`
	Draft                               *Draft                               `json:"draft,omitempty"`
	ExpendCommittedBufferResearchPoints *ExpendCommittedBufferResearchPoints `json:"expend_committed_buffer_research_points,omitempty"`
	ExpendPrototype                     *ExpendPrototype                     `json:"expend_prototype,omitempty"`
	ExpendResearchPointsOnly            *ExpendResearchPointsOnly            `json:"expend_research_points_only,omitempty"`
	FactoryGroupChange                  *FactoryGroupChange                  `json:"factory_group_change,omitempty"`
	Give                                *Give                                `json:"give,omitempty"`
	HomePortChange                      *HomePortChange                      `json:"home_port_change,omitempty"`
	Invade                              *Invade                              `json:"invade,omitempty"`
	Jump                                *Jump                                `json:"jump,omitempty"`
	Junk                                *Junk                                `json:"junk,omitempty"`
	LaunchRobotProbe                    *LaunchRobotProbe                    `json:"launch_robot_probe,omitempty"`
	LoadCargo                           *LoadCargo                           `json:"load_cargo,omitempty"`
	Merge                               *Merge                               `json:"merge,omitempty"`
	Message                             *Message                             `json:"message,omitempty"`
	MineChange                          *MineChange                          `json:"mine_change,omitempty"`
	MineShutDown                        *MineShutDown                        `json:"mine_shut_down,omitempty"`
	MineStartUp                         *MineStartUp                         `json:"mine_start_up,omitempty"`
	Move                                *Move                                `json:"move,omitempty"`
	Name                                *Name                                `json:"name,omitempty"`
	Note                                *Note                                `json:"note,omitempty"`
	OffensiveSupport                    *OffensiveSupport                    `json:"offensive_support,omitempty"`
	Pay                                 *Pay                                 `json:"pay,omitempty"`
	PermissionToColonize                *PermissionToColonize                `json:"permission_to_colonize,omitempty"`
	PickUpItem                          *PickUpItem                          `json:"pick_up_item,omitempty"`
	PickUpPopulation                    *PickUpPopulation                    `json:"pick_up_population,omitempty"`
	PreManeuverEnergyWeaponFire         *PreManeuverEnergyWeaponFire         `json:"pre_maneuver_energy_weapon_fire,omitempty"`
	PreManeuverMissileFire              *PreManeuverMissileFire              `json:"pre_maneuver_missile_fire,omitempty"`
	Probe                               *Probe                               `json:"probe,omitempty"`
	ProbeOrbit                          *ProbeOrbit                          `json:"probe_orbit,omitempty"`
	ProbeSystem                         *ProbeSystem                         `json:"probe_system,omitempty"`
	Ration                              *Ration                              `json:"ration,omitempty"`
	Run                                 *Run                                 `json:"run,omitempty"`
	Scrap                               *Scrap                               `json:"scrap,omitempty"`
	SetUp                               *SetUp                               `json:"set_up,omitempty"`
	ShutDown                            *ShutDown                            `json:"shut_down,omitempty"`
	StartUp                             *StartUp                             `json:"start_up,omitempty"`
	Survey                              *Survey                              `json:"survey,omitempty"`
	TacticalManeuver                    *TacticalManeuver                    `json:"tactical_maneuver,omitempty"`
	Transfer                            *Transfer                            `json:"transfer,omitempty"`
	UncontrolPlanet                     *UncontrolPlanet                     `json:"uncontrol_planet,omitempty"`
	Undock                              *Undock                              `json:"undock,omitempty"`
	UnloadCargo                         *UnloadCargo                         `json:"unload_cargo,omitempty"`
	Withdraw                            *Withdraw                            `json:"withdraw,omitempty"`
}
