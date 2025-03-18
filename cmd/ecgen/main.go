// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

func main() {

}

func unitTable() {
	type unit_t struct {
		Code                    string
		Name                    string
		Category                string
		OperationalRequirements string
		Notes                   string
		Text                    string
		Mass                    func(tl int) int
		VolumeAssembled         func(tl int) int
		VolumeDisassembled      func(tl int) int
		MetalsToBuild           func(tl int) int
		NonMetalsToBuild        func(tl int) int
	}
	for _, unit := range []unit_t{
		{Code: "ANM    ", Name: "Anti-Missiles             ", Category: "Vehicles     ", OperationalRequirements: "Missile Launcher of same TL                                  ", Notes: "Destroys Missiles; see combat",
			Text: "4 x TL       4 x TL      2 x TL           2 x TL               ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "ASC    ", Name: "Assault Craft             ", Category: "Vehicles     ", OperationalRequirements: "1 soldier or military robot equivalent + 0.1 fuel in combat  ", Notes: "Provides 10 x TL combat factors; does not require transports to attack",
			Text: "5 x TL       5 x TL      3 x TL           2 x TL               ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "ASW    ", Name: "Assault Weapons           ", Category: "Vehicles     ", OperationalRequirements: "1 soldier or military robot equivalent                       ", Notes: "Provides 2 x TL^2 combat factors",
			Text: "2 x TL       2 x TL      1 x TL           1 x TL               ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "AUT    ", Name: "Automation                ", Category: "Assembly     ", OperationalRequirements: "Must be assembled                                            ", Notes: "Replaces 1 x TL^2 Unskilled; see Automation in Production chapter",
			Text: "4 x TL       2 x TL      2 x TL           2 x TL               ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "CNGD   ", Name: "Consumer Goods            ", Category: "Consumables  ", OperationalRequirements: "                                                             ", Notes: "Consumption determines SOL",
			Text: "0.6          0.3         0.2              0.4                  ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "ESH    ", Name: "Energy Shields            ", Category: "Assembly     ", OperationalRequirements: "1 soldier / 100, uses 10 x TL fuel                           ", Notes: "Deflects 10 x TL^2 energy units per use",
			Text: "20 x TL      10 x TL     10 x TL          10 x TL              ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "EWP    ", Name: "Energy Weapons            ", Category: "Assembly     ", OperationalRequirements: "1 soldier / 100, uses 4 x TL fuel                            ", Notes: "Destroys 10 x TL^2 mass per hit",
			Text: "10 x TL      5 x TL      5 x TL           5 x TL               ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "FCT    ", Name: "Factories                 ", Category: "Assembly     ", OperationalRequirements: "1 professional + 3 unskilled, uses 0.5 fuel or power         ", Notes: "Produces 20 x TL mass per turn; see Manufacturing",
			Text: "2 x TL + 12  TL + 6      8 + TL           4 + TL               ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "FOOD   ", Name: "Food                      ", Category: "Consumables  ", OperationalRequirements: "                                                             ", Notes: "Feeds 4 to 16 population each turn; see Basic Needs",
			Text: "6            3           0                0                    ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "FRM    ", Name: "Farms                     ", Category: "Assembly     ", OperationalRequirements: "1 professional + 3 unskilled, fuel varies                    ", Notes: "Production varies via TL; see Farming",
			Text: "2 x TL + 6   TL + 3      4 + TL           2 + TL               ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "FUEL   ", Name: "Fuel                      ", Category: "Consumables  ", OperationalRequirements: "                                                             ", Notes: "Raw material used by many units",
			Text: "1            0.5         0                0                    ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "GOLD   ", Name: "Gold                      ", Category: "Consumables  ", OperationalRequirements: "                                                             ", Notes: "Raw material with no use in the game",
			Text: "1            0.5         0                0                    ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "HEN    ", Name: "Hyper Engines             ", Category: "Assembly     ", OperationalRequirements: "1 professional / 100, uses 40 fuel per light year            ", Notes: "Lift capacity 1045 x TL, range is 3 x √TL",
			Text: "45 x TL      22.5 x TL   25 x TL          20 x TL              ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "LAB    ", Name: "Laboratories              ", Category: "Assembly     ", OperationalRequirements: "3 professional + 1 unskilled, 0.5 fuel/power                 ", Notes: "Produces 0.25 research per turn",
			Text: "2 x TL + 8   TL + 4      5 + TL           3 + TL               ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "LFS    ", Name: "Life Supports             ", Category: "Assembly     ", OperationalRequirements: "1 x TL fuel or power                                         ", Notes: "Supports 1 x TL^2 population",
			Text: "8 x TL       4 x TL      3 x TL           5 x TL               ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "METS   ", Name: "Metals                    ", Category: "Consumables  ", OperationalRequirements: "                                                             ", Notes: "Raw material used in production",
			Text: "1            0.5         0                0                    ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "MIN    ", Name: "Mines                     ", Category: "Assembly     ", OperationalRequirements: "1 professional + 3 unskilled, 0.5 fuel or power              ", Notes: "Mines 25 x TL per turn in raw ore, Actual net depends on yield of deposit; see Mining",
			Text: "2 x TL + 10  TL + 5      5 + TL           5 + TL               ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "MSL    ", Name: "Missile Launchers         ", Category: "Assembly     ", OperationalRequirements: "1 soldier / 100                                              ", Notes: "Launches 1 missile per attack; see Combat",
			Text: "25 x TL      12.5 x TL   15 + TL          10 + TL              ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "MSS    ", Name: "Missiles                  ", Category: "Vehicles     ", OperationalRequirements: "Missile Launcher of same TL                                  ", Notes: "Destroys 100 x TL^2 Mass",
			Text: "4 x TL       4 x TL      2 x TL           2 x TL               ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "MTBT   ", Name: "Military Robots           ", Category: "Bots         ", OperationalRequirements: "2 x TL military supplies                                     ", Notes: "Same as 2 x TL soldiers",
			Text: "2 x TL + 20  TL + 10     10 + TL          10 + TL              ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "MTSP   ", Name: "Military Supplies         ", Category: "Consumables  ", OperationalRequirements: "                                                             ", Notes: "Required by soldiers in combat",
			Text: "0.04         0.02        0.02             0.02                 ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "NMTS   ", Name: "Non-Metals                ", Category: "Consumables  ", OperationalRequirements: "                                                             ", Notes: "Raw material used in production",
			Text: "1            0.5         0                0                    ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "PWP    ", Name: "Power Plants              ", Category: "Assembly     ", OperationalRequirements: "                                                             ", Notes: "Produces TL Power per turn (think hydro electric)",
			Text: "2 x TL + 10  TL + 5      5 + TL           5 + TL               ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "RPV    ", Name: "Robot Probe Vehicles      ", Category: "Bots         ", OperationalRequirements: "                                                             ", Notes: "Obtains probe data, expended when used",
			Text: "500 / TL     500 / TL    200 / TL         300 / TL             ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "RSCH   ", Name: "Research                  ", Category: "Consumables  ", OperationalRequirements: "                                                             ", Notes: "Expended to increase TLs",
			Text: "0            0           0                0                    ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "SEN    ", Name: "Sensors                   ", Category: "Assembly     ", OperationalRequirements: "uses 0.05 x TL fuel                                          ", Notes: "Used to obtain probe information",
			Text: "3000 x TL    1500 x TL   1000 x TL        2000 x TL            ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "SLS    ", Name: "Light Structure           ", Category: "Assembly     ", OperationalRequirements: "May only be built in Orbiting Colonies                       ", Notes: "Encloses (1 x TL^2) divided by type factor",
			Text: "0.01 x TL    0.005 x TL  0.005 x TL       0.005 x TL           ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "SPD    ", Name: "Space Drives              ", Category: "Assembly     ", OperationalRequirements: "1 professional / 100, uses 1 x TL fuel                       ", Notes: "Produces 3000 x TL^2 thrust",
			Text: "25 x TL      12.5 X TL   15 X TL          10 x TL              ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "STU    ", Name: "Structure                 ", Category: "Assembly     ", OperationalRequirements: "                                                             ", Notes: "Encloses (1 x TL^2) divided by type factor",
			Text: "0.1 x TL     0.05 x TL   0.07 x TL        0.03 x TL            ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
		{Code: "TPT    ", Name: "Transports                ", Category: "Vehicles     ", OperationalRequirements: "1 professional / 10, uses 0.1 x TL^2 fuel                    ", Notes: "Transports 20 x TL^2 Mass",
			Text: "4 x TL       4 x TL      2 x TL           2 x TL               ",
			Mass: func(tl int) int { return 4 * tl }, VolumeAssembled: func(tl int) int { return 4 * tl }, VolumeDisassembled: func(tl int) int { return 4 * tl }, MetalsToBuild: func(tl int) int { return 2 * tl }, NonMetalsToBuild: func(tl int) int { return 2 * tl }},
	} {
		_ = unit
	}
}
