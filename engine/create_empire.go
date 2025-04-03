// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/playbymail/empyr/internal/domains"
	"github.com/playbymail/empyr/repos/sqlite"
	"log"
	"math"
	"strconv"
	"strings"
)

const (
	ErrEmpireNotAvailable = Error("empire not available")
	ErrMissingDeposit     = Error("missing deposit")
	ErrNoEmpireAvailable  = Error("no empire available")
)

type CreateEmpireParams_t struct {
	EmpireID int64
	Username string
	Email    string
}

func CreateEmpireCommand(e *Engine_t, cfg *CreateEmpireParams_t) (int64, error) {
	panic("not implemented")
}

// createEmpire creates a new empire in the game.
// It does not create a home colony or fiddle with mining operations.
func createEmpire(g *Game_t) (*Empire_t, error) {
	e := &Empire_t{
		Id:         g.nextAvailableEmpireID(),
		HomeSystem: g.Home.System,
		HomeStar:   g.Home.Star,
		HomeOrbit:  g.Home.Orbit,
		HomePlanet: g.Home.Planet,
	}
	if e.Id == 0 {
		return nil, ErrNoEmpireAvailable
	}
	e.EmpireNo = e.Id
	e.Name = fmt.Sprintf("Empire %d", e.Id)
	return e, nil
}

var (
	//go:embed setup/home-colony.json
	homeColonyJSON []byte
)

// createHomeColony creates a home colony for the empire. The colony includes
// population, inventory, and farm+factory+mine operations.
func createHomeColony(g *Game_t, e *Empire_t) (*Entity_t, error) {
	// load the default colony setup from the JSON file
	defaultColony, err := loadDefaultHomeColony()
	if err != nil {
		return nil, err
	}

	// create a new home colony
	sc := &Entity_t{
		Id:          g.nextAvailableEntityID(),
		IsColony:    true,
		IsOnSurface: true,
		TechLevel:   defaultColony.TechLevel,
		Population: struct {
			UEM PopulationGroup_t
			USK PopulationGroup_t
			PRO PopulationGroup_t
			SLD PopulationGroup_t
			CNW PopulationGroup_t
			SPY PopulationGroup_t
		}{
			UEM: PopulationGroup_t{Qty: defaultColony.Census["UEM"], Pay: defaultColony.PayRates["UEM"], RebelQty: defaultColony.Rebels["UEM"]},
			USK: PopulationGroup_t{Qty: defaultColony.Census["USK"], Pay: defaultColony.PayRates["USK"], RebelQty: defaultColony.Rebels["USK"]},
			PRO: PopulationGroup_t{Qty: defaultColony.Census["PRO"], Pay: defaultColony.PayRates["PRO"], RebelQty: defaultColony.Rebels["PRO"]},
			SLD: PopulationGroup_t{Qty: defaultColony.Census["SLD"], Pay: defaultColony.PayRates["SLD"], RebelQty: defaultColony.Rebels["SLD"]},
			CNW: PopulationGroup_t{Qty: defaultColony.Census["CNW"], Pay: defaultColony.PayRates["CNW"], RebelQty: defaultColony.Rebels["CNW"]},
			SPY: PopulationGroup_t{Qty: defaultColony.Census["SPY"], Pay: defaultColony.PayRates["SPY"], RebelQty: defaultColony.Rebels["SPY"]},
		},
	}

	// load the default values for inventory
	sc.Inventory, err = loadHomeColonyDefaultInventory(defaultColony.Inventory)
	if err != nil {
		return nil, err
	}

	// load the default values for farms
	sc.FarmGroups, err = loadHomeColonyDefaultFarmGroups(sc)
	if err != nil {
		return nil, err
	}
	// allocate resources to the farm group.
	for _, grp := range sc.FarmGroups {
		// fake the constraints since this is the setup
		_ = grp.Allocate(FarmGroupConstraints_t{
			Fuel: 99_999_999_999,
			Pro:  99_999_999_999,
			Usk:  99_999_999_999,
			Aut:  99_999_999_999,
		})
	}

	// load the default values for factories
	sc.FactoryGroups, err = loadHomeColonyDefaultFactoryGroups(sc)
	if err != nil {
		return nil, err
	}
	// populate wanted and actual resources for the factory group.
	// assume 100% allocation of wanted resources.
	for _, grp := range sc.FactoryGroups {
		for _, unit := range grp.Units {
			wantMets, wantNmets := unitRequirements(grp.Tooling.Item.Code, grp.Tooling.TechLevel, unit.NbrOfUnits)
			unit.Wanted = &FactoryGroupResources_t{
				Pro:  float64(unit.NbrOfUnits),
				Usk:  float64(unit.NbrOfUnits) * 3,
				Fuel: factoryFuel(sc, unit.TechLevel, unit.NbrOfUnits),
				Mets: wantMets,
				Nmts: wantNmets,
			}
			unit.Actual = &FactoryGroupResources_t{
				Pro:  unit.Wanted.Pro,
				Usk:  unit.Wanted.Usk,
				Fuel: unit.Wanted.Fuel,
				Mets: unit.Wanted.Mets,
				Nmts: unit.Wanted.Nmts,
			}
			unit.Consumed = &FactoryGroupResources_t{
				Pro:  unit.Actual.Pro,
				Usk:  unit.Actual.Usk,
				Fuel: unit.Actual.Fuel,
				Mets: unit.Actual.Mets,
				Nmts: unit.Actual.Nmts,
			}
			// scale the production to units per turn
			qtyProducedPerYear := 0.0
			qtyProducedPerTurn := qtyProducedPerYear / 4
			unit.Produced = &FactoryGroupResources_t{
				WIP: [3]float64{qtyProducedPerTurn, qtyProducedPerTurn, qtyProducedPerTurn},
			}
		}
	}

	// load the default values for mining groups
	sc.MiningGroups, err = loadHomeColonyDefaultMiningGroups(e.HomeOrbit, g.Deposits)
	if err != nil {
		return nil, err
	}
	// populate wanted and actual resources for the mining group.
	// assume 100% allocation of wanted resources.
	for _, grp := range sc.MiningGroups {
		for _, unit := range grp.Units {
			unit.Wanted = &MineGroupResources_t{
				Pro:  float64(unit.NbrOfUnits),
				Usk:  float64(unit.NbrOfUnits) * 3,
				Fuel: mineFuel(sc, unit.TechLevel, unit.NbrOfUnits),
			}
			unit.Actual = &MineGroupResources_t{
				Pro:  unit.Wanted.Pro,
				Usk:  unit.Wanted.Usk,
				Fuel: unit.Wanted.Fuel,
			}
			unit.Consumed = &MineGroupResources_t{
				Pro:  unit.Actual.Pro,
				Usk:  unit.Actual.Usk,
				Fuel: unit.Actual.Fuel,
			}
			// scale the production to units per turn.
			// we are going to ignore the amount of resources available in the deposit.
			switch grp.Deposit.Resource {
			case NONE: // ignore
			case GOLD:
				unit.Produced = &MineGroupResources_t{
					Gold: float64(unit.NbrOfUnits) * 100 / 4,
				}
			case FUEL:
				unit.Produced = &MineGroupResources_t{
					Fuel: float64(unit.NbrOfUnits) * 100 / 4,
				}
			case METALLICS:
				unit.Produced = &MineGroupResources_t{
					Mets: float64(unit.NbrOfUnits) * 100 / 4,
				}
			case NON_METALLICS:
				unit.Produced = &MineGroupResources_t{
					Nmts: float64(unit.NbrOfUnits) * 100 / 4,
				}
			}
		}
		// roll the units up to the group summary
		grp.Summary.Wanted = &MineGroupResources_t{}
		grp.Summary.Actual = &MineGroupResources_t{}
		grp.Summary.Consumed = &MineGroupResources_t{}
		grp.Summary.Produced = &MineGroupResources_t{}
		for _, unit := range grp.Units {
			grp.Summary.Wanted.Fuel += unit.Wanted.Fuel
			grp.Summary.Wanted.Pro += unit.Wanted.Pro
			grp.Summary.Wanted.Usk += unit.Wanted.Usk
			grp.Summary.Actual.Fuel += unit.Actual.Fuel
			grp.Summary.Actual.Pro += unit.Actual.Pro
			grp.Summary.Actual.Usk += unit.Actual.Usk
			grp.Summary.Consumed.Fuel += unit.Consumed.Fuel
			grp.Summary.Consumed.Pro += unit.Consumed.Pro
			grp.Summary.Consumed.Usk += unit.Consumed.Usk
			grp.Summary.Produced.Fuel += unit.Produced.Fuel
			grp.Summary.Produced.Gold += unit.Produced.Gold
			grp.Summary.Produced.Mets += unit.Produced.Mets
			grp.Summary.Produced.Nmts += unit.Produced.Nmts
		}
	}
	// allocate resources to the mine group.
	for _, grp := range sc.MiningGroups {
		// fake the constraints since this is the setup
		_ = grp.Allocate(MineGroupConstraints_t{
			Fuel: 99_999_999_999,
			Pro:  99_999_999_999,
			Usk:  99_999_999_999,
			Aut:  99_999_999_999,
		})
	}

	// run some of the steps to simulate a turn?
	for _, grp := range sc.FarmGroups {
		for _, unit := range grp.Units {
			unit.Consume()
		}
	}
	for _, grp := range sc.FarmGroups {
		for _, unit := range grp.Units {
			unit.Produce()
		}
		grp.Summarize()
	}

	return sc, nil
}

var (
	//go:embed setup/factory-groups.json
	homeFactoryGroupsJSON []byte
)

type defaultFactoryGroup_t struct {
	No         int64    `json:"no"`
	NbrOfUnits int64    `json:"number-of-units"`
	Tooling    string   `json:"tooling"`
	WIP        [3]int64 `json:"wip"`
}

func loadHomeColonyDefaultFactoryGroups(sc *Entity_t) ([]*FactoryGroup_t, error) {
	var defaultFactoryGroups []defaultFactoryGroup_t
	err := json.Unmarshal(homeColonyJSON, &defaultFactoryGroups)
	if err != nil {
		return nil, err
	}

	var factoryGroups []*FactoryGroup_t
	for k, v := range defaultFactoryGroups {
		factoryGroup := &FactoryGroup_t{
			Id:     int64(k) + 1,
			Entity: sc,
			No:     v.No,
		}
		cd, techLevel, err := codeTechLevelFromString(v.Tooling)
		if err != nil {
			log.Printf("error: loadHomeColonyDefaultFactoryGroups %q: %v", v.Tooling, err)
			return nil, err
		}
		unit, ok := unitTable[cd]
		if !ok {
			log.Printf("error: loadHomeColonyDefaultFactoryGroups %q: %v", v.Tooling, ErrInvalidCode)
			return nil, ErrInvalidCode
		}
		factoryGroup.Tooling = &FactoryGroupTooling_t{
			Group:     factoryGroup,
			Item:      unit,
			TechLevel: techLevel,
		}
		factoryGroup.Units = append(factoryGroup.Units, &FactoryGroupUnit_t{
			Group:       factoryGroup,
			TechLevel:   1,
			NbrOfUnits:  v.NbrOfUnits,
			QtyProduced: v.WIP,
		})
		factoryGroups = append(factoryGroups, factoryGroup)
	}

	return factoryGroups, nil
}

func loadHomeColonyDefaultFarmGroups(sc *Entity_t) ([]*FarmGroup_t, error) {
	var farmGroups []*FarmGroup_t

	farmGroup := &FarmGroup_t{
		Id:     1, // todo: set this later
		Entity: sc,
		No:     1,
		Units:  []*FarmGroupUnit_t{{TechLevel: 1, NbrOfUnits: 130_000}},
	}
	farmGroup.Units[0].Group = farmGroup
	farmGroups = append(farmGroups, farmGroup)

	return farmGroups, nil
}

// apply the inventory
func loadHomeColonyDefaultInventory(defaultInventory defaultInventory_t) ([]*Inventory_t, error) {
	var inventory []*Inventory_t

	for cd, qty := range defaultInventory.Storage.NonAssembled {
		code, techLevel, err := codeTechLevelFromString(cd)
		if err != nil {
			log.Printf("error: createHomeColony %q: %v", cd, err)
			return nil, err
		}
		mass, volume := unitMassAndVolume(code, techLevel, qty, false)
		unit, ok := unitTable[code]
		if !ok {
			log.Printf("error: createHomeColony %q: %v", cd, ErrInvalidCode)
			return nil, ErrInvalidCode
		}
		inventory = append(inventory, &Inventory_t{
			Unit:      unit,
			TechLevel: techLevel,
			Qty:       qty,
			Mass:      mass,
			Volume:    volume,
		})
	}
	for cd, qty := range defaultInventory.Storage.Disassembled {
		code, techLevel, err := codeTechLevelFromString(cd)
		if err != nil {
			log.Printf("error: createHomeColony %q: %v", cd, err)
			return nil, err
		}
		mass, volume := unitMassAndVolume(code, techLevel, qty, false)
		unit, ok := unitTable[code]
		if !ok {
			log.Printf("error: createHomeColony %q: %v", cd, ErrInvalidCode)
			return nil, ErrInvalidCode
		}
		inventory = append(inventory, &Inventory_t{
			Unit:      unit,
			TechLevel: techLevel,
			Qty:       qty,
			Mass:      mass,
			Volume:    volume,
		})
	}
	for cd, qty := range defaultInventory.Assembled {
		code, techLevel, err := codeTechLevelFromString(cd)
		if err != nil {
			log.Printf("error: createHomeColony %q: %v", cd, err)
			return nil, err
		}
		mass, volume := unitMassAndVolume(code, techLevel, qty, true)
		unit, ok := unitTable[code]
		if !ok {
			log.Printf("error: createHomeColony %q: %v", cd, ErrInvalidCode)
			return nil, ErrInvalidCode
		}
		inventory = append(inventory, &Inventory_t{
			Unit:        unit,
			TechLevel:   techLevel,
			Qty:         qty,
			Mass:        mass,
			Volume:      volume,
			IsAssembled: true,
		})
	}
	return inventory, nil
}

const (
	ErrNoFuelDeposits = Error("no fuel deposits on home world")
	ErrNoGoldDeposits = Error("no gold deposits on home world")
	ErrNoMetsDeposits = Error("no mets deposits on home world")
	ErrNoNmtsDeposits = Error("no nmts deposits on home world")
)

func loadHomeColonyDefaultMiningGroups(orbit *Orbit_t, deposits map[int64]*Deposit_t) ([]*MineGroup_t, error) {
	var miningGroups []*MineGroup_t

	// get the best deposits for the home planet
	gold, fuel, mets, nmts := getInitialMiningSites(orbit, deposits)
	if gold == nil {
		return nil, ErrNoGoldDeposits
	} else if fuel == nil {
		return nil, ErrNoFuelDeposits
	} else if mets == nil {
		return nil, ErrNoMetsDeposits
	} else if nmts == nil {
		return nil, ErrNoNmtsDeposits
	}

	// create mining groups for each production goal
	for no, goal := range []struct {
		Deposit        *Deposit_t
		ProductionGoal float64
	}{
		{Deposit: gold, ProductionGoal: 1_000},
		{Deposit: fuel, ProductionGoal: 500_000},
		{Deposit: mets, ProductionGoal: 2_750_000 / 2},
		{Deposit: nmts, ProductionGoal: 2_750_000 / 2},
	} {
		yieldPct := float64(goal.Deposit.Yield) / 100 // convert percentage to float
		rawOreNeeded := int64(math.Ceil(goal.ProductionGoal / yieldPct))
		nbrUnitsNeeded := (rawOreNeeded + 24) / 25 // Ceiling division by 25 (each unit mines 25 units of ore)
		miningGroups = append(miningGroups, &MineGroup_t{
			Id:      0, // must set this later
			No:      int64(no + 1),
			Deposit: goal.Deposit,
			Units: []*MineGroupUnit_t{{
				Deposit:    goal.Deposit,
				TechLevel:  1,
				NbrOfUnits: nbrUnitsNeeded,
			}},
		})
	}

	return miningGroups, nil
}

type defaultStorage_t struct {
	NonAssembled map[string]int64 `json:"non-assembled"`
	Disassembled map[string]int64 `json:"disassembled"`
}
type defaultInventory_t struct {
	Storage   defaultStorage_t `json:"storage"`
	Assembled map[string]int64 `json:"assembled"`
}
type defaultColony_t struct {
	TechLevel        int64              `json:"tech-level"`
	StandardOfLiving float64            `json:"standard-of-living"`
	Census           map[string]int64   `json:"census"`
	BirthRate        float64            `json:"birth-rate"`
	DeathRate        float64            `json:"death-rate"`
	Rations          int                `json:"rations"`
	PayRates         map[string]float64 `json:"pay-rates"`
	Rebels           map[string]int64   `json:"rebels"`
	Inventory        defaultInventory_t `json:"inventory"`
}

// load the default values for colonies
func loadDefaultHomeColony() (*defaultColony_t, error) {
	defaultColony := defaultColony_t{
		Census:   map[string]int64{},
		PayRates: map[string]float64{},
		Rebels:   map[string]int64{},
		Inventory: defaultInventory_t{
			Storage: defaultStorage_t{
				NonAssembled: map[string]int64{},
				Disassembled: map[string]int64{},
			},
			Assembled: map[string]int64{},
		},
	}
	err := json.Unmarshal(homeColonyJSON, &defaultColony)
	if err != nil {
		return nil, err
	}
	return &defaultColony, nil
}

func createEmpireCommand(e *Engine_t, cfg *CreateEmpireParams_t) (int64, error) {
	if cfg.Username == "" {
		return 0, ErrMissingHandle
	} else if _, err := IsValidHandle(cfg.Username); err != nil {
		return 0, err
	}
	if cfg.Email == "" {
		cfg.Email = fmt.Sprintf("%s@epimethean.dev", cfg.Username)
	}
	log.Printf("create: empire: user %q: email %q\n", cfg.Username, cfg.Email)

	q, tx, err := e.Store.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	gameRow, err := q.ReadAllGameInfo(e.Store.Context)
	if err != nil {
		return 0, err
	}
	if gameRow.CurrentTurn != 0 {
		return 0, ErrGameInProgress
	}

	type deposit_t struct {
		no       int64
		kind     string
		qty      int64
		yieldPct float64
	}
	deposits := map[int64]deposit_t{}
	if rows, err := q.ReadDepositsByOrbit(e.Store.Context, sqlite.ReadDepositsByOrbitParams{
		OrbitID: gameRow.HomeOrbitID,
		TurnNo:  gameRow.CurrentTurn,
	}); err != nil {
		return 0, err
	} else {
		for _, row := range rows {
			deposits[row.DepositNo] = deposit_t{
				no:       row.DepositNo,
				kind:     row.Kind,
				qty:      row.Qty,
				yieldPct: float64(row.YieldPct) / 100,
			}
		}
	}

	empireID := cfg.EmpireID
	if empireID == 0 {
		empireID, err = q.ReadNextEmpireNumber(e.Store.Context)
		if err != nil {
			return 0, err
		}
	}
	if isActive, err := q.IsEmpireActive(e.Store.Context, cfg.EmpireID); err != nil {
		return 0, err
	} else if isActive != 0 {
		return 0, ErrEmpireNotAvailable
	}
	err = q.UpdateEmpireStatus(e.Store.Context, sqlite.UpdateEmpireStatusParams{
		EmpireID: cfg.EmpireID,
		IsActive: 1,
	})
	if err != nil {
		return 0, err
	}
	parms := sqlite.CreateEmpireWithIDParams{
		ID:           empireID,
		HomeSystemID: gameRow.HomeSystemID,
		HomeStarID:   gameRow.HomeStarID,
		HomeOrbitID:  gameRow.HomeOrbitID,
	}
	err = q.CreateEmpireWithID(e.Store.Context, parms)
	if err != nil {
		return 0, err
	}
	err = q.CreateEmpirePlayer(e.Store.Context, sqlite.CreateEmpirePlayerParams{
		EmpireID:   empireID,
		Effdt:      0,
		Enddt:      domains.MaxGameTurnNo,
		EmpireName: fmt.Sprintf("E%03d", empireID),
		Username:   cfg.Username,
		Email:      cfg.Email,
	})

	// create a home open surface colony
	scParams := sqlite.CreateSCParams{
		EmpireID:    empireID,
		ScCd:        "COPN",
		ScTechLevel: 1,
	}
	scId, err := q.CreateSC(e.Store.Context, scParams)
	if err != nil {
		log.Printf("create: empire: sc %+v\n", scParams)
		return 0, err
	}
	log.Printf("create: empire: id %d: colony %d\n", empireID, scId)

	for _, pop := range []struct {
		code string
		qty  int64
		pay  float64
	}{
		{code: "UEM", qty: 5_900_000, pay: 0.0},
		{code: "USK", qty: 6_000_000, pay: 0.125},
		{code: "PRO", qty: 1_500_000, pay: 0.375},
		{code: "SLD", qty: 2_500_000, pay: 0.25},
		{code: "CNW", qty: 10_000, pay: 0.5},
		{code: "SPY", qty: 20, pay: 0.625},
	} {
		scPopParms := sqlite.CreateSCPopulationParams{
			ScID:         scId,
			PopulationCd: pop.code,
			Qty:          pop.qty,
			PayRate:      pop.pay,
			RebelQty:     0,
		}
		err = q.CreateSCPopulation(e.Store.Context, scPopParms)
		if err != nil {
			log.Printf("create: empire: sc pop %+v\n", scPopParms)
			return 0, err
		}
	}

	for _, unit := range []struct {
		code        string
		techLevel   int64
		qty         int64
		isAssembled bool
		isStored    bool
	}{
		{code: "ANM", techLevel: 1, qty: 150_000, isStored: true},
		{code: "ASW", techLevel: 1, qty: 2_400_000, isStored: true},
		{code: "AUT", techLevel: 1, qty: 93_750, isStored: true},
		{code: "CNGD", techLevel: 0, qty: 940_821, isStored: true},
		{code: "EWP", techLevel: 1, qty: 37_500, isStored: true},
		{code: "FOOD", techLevel: 0, qty: 4_269_990, isStored: true},
		{code: "FRM", techLevel: 1, qty: 130_000, isAssembled: true},
		{code: "FUEL", techLevel: 0, qty: 1_360_000, isStored: true},
		{code: "GOLD", techLevel: 0, qty: 20_000, isStored: true},
		{code: "METS", techLevel: 0, qty: 5_354_167, isStored: true},
		{code: "MIN", techLevel: 1, qty: 300_000, isAssembled: true},
		{code: "MIN", techLevel: 1, qty: 31_250, isStored: true},
		{code: "MSL", techLevel: 1, qty: 50_000, isAssembled: true},
		{code: "MTSP", techLevel: 1, qty: 150_000, isStored: true},
		{code: "NMTS", techLevel: 0, qty: 2_645_833, isStored: true},
		{code: "SEN", techLevel: 1, qty: 20, isAssembled: true},
		{code: "SEN", techLevel: 1, qty: 20, isStored: true},
		{code: "STU", qty: 60_000_000, isAssembled: true},
		{code: "STU", techLevel: 0, qty: 14_500_000, isStored: true},
		{code: "TPT", techLevel: 1, qty: 20_000, isStored: true},
	} {
		scInvParams := sqlite.CreateSCInventoryParams{
			ScID:      scId,
			UnitCd:    unit.code,
			TechLevel: unit.techLevel,
			Qty:       unit.qty,
			Mass:      Mass(unit.code, unit.techLevel, unit.qty),
		}
		if IsOperational(unit.code) {
			if unit.isAssembled {
				scInvParams.IsAssembled = 1
				scInvParams.Volume = VolumeAssembled(unit.code, unit.techLevel, unit.qty)
			} else {
				scInvParams.Volume = VolumeDisassembled(unit.code, unit.techLevel, unit.qty)
			}
		}
		if unit.isStored {
			scInvParams.IsStored = 1
			scInvParams.Volume = VolumeStored(unit.code, unit.techLevel, unit.qty)
		}
		err = q.CreateSCInventory(e.Store.Context, scInvParams)
		if err != nil {
			log.Printf("create: empire: sc inv %+v\n", scInvParams)
			return 0, err
		}
	}

	type factoryGroupUnit struct {
		code            string // the code being produced (eg. "FRM")
		techLevel       int64  // optional tech level of the item being produced
		nbrOfUnits      int64  // number of units in the pipeline
		ordersCode      string
		ordersTechLevel int64
		wip             [3]int64 // newest to oldest, next to finish is always last
	}
	type factoryGroupData struct {
		groupNo         int64
		ordersCode      string
		ordersTechLevel int64
		retoolTurnNo    int64 // the turn number the factory group is on
		units           []*factoryGroupUnit
	}
	for _, fg := range []factoryGroupData{
		{groupNo: 1, ordersCode: "CNGD",
			units: []*factoryGroupUnit{
				{code: "FCT", techLevel: 1, nbrOfUnits: 250_000, ordersCode: "CNGD", wip: [3]int64{2_083_333, 2_083_333, 2_083_333}},
			},
		},
		{groupNo: 2, ordersCode: "MTSP",
			units: []*factoryGroupUnit{
				{code: "FCT", techLevel: 1, nbrOfUnits: 75_000, ordersCode: "MTSP", wip: [3]int64{9_375_000, 9_375_000, 9_375_000}},
			},
		},
		{groupNo: 3, ordersCode: "AUT", ordersTechLevel: 1,
			units: []*factoryGroupUnit{
				{code: "FCT", techLevel: 1, nbrOfUnits: 75_000, ordersCode: "AUT", ordersTechLevel: 1, wip: [3]int64{93_750, 93_750, 93_750}},
			},
		},
		{groupNo: 4, ordersCode: "EWP", ordersTechLevel: 1,
			units: []*factoryGroupUnit{
				{code: "FCT", techLevel: 1, nbrOfUnits: 75_000, ordersCode: "EWP", ordersTechLevel: 1, wip: [3]int64{37_500, 37_500, 37_500}},
			},
		},
		{groupNo: 5, ordersCode: "MIN", ordersTechLevel: 1,
			units: []*factoryGroupUnit{
				{code: "FCT", techLevel: 1, nbrOfUnits: 75_000, ordersCode: "MIN", ordersTechLevel: 1, wip: [3]int64{31_250, 31_250, 31_250}},
			},
		},
		{groupNo: 6, ordersCode: "STU",
			units: []*factoryGroupUnit{
				{code: "FCT", techLevel: 1, nbrOfUnits: 250_000, ordersCode: "STU", wip: [3]int64{12_500_000, 12_500_000, 12_500_000}},
			},
		},
		{groupNo: 7, ordersCode: "RSCH",
			units: []*factoryGroupUnit{
				{code: "FCT", techLevel: 1, nbrOfUnits: 50_000, ordersCode: "RSCH", wip: [3]int64{12_500, 12_500, 12_500}},
			},
		},
	} {
		fgParms := sqlite.CreateSCFactoryGroupParams{
			ScID:            scId,
			GroupNo:         fg.groupNo,
			OrdersCd:        fg.ordersCode,
			OrdersTechLevel: fg.ordersTechLevel,
		}
		err := q.CreateSCFactoryGroup(e.Store.Context, fgParms)
		if err != nil {
			log.Printf("create: empire: sc factory group %+v\n", fgParms)
			return 0, err
		}
		for _, unit := range fg.units {
			fgUnitParms := sqlite.CreateSCFactoryGroupUnitParams{
				ScID:            scId,
				GroupNo:         fg.groupNo,
				GroupTechLevel:  unit.techLevel,
				NbrOfUnits:      unit.nbrOfUnits,
				OrdersCd:        unit.ordersCode,
				OrdersTechLevel: unit.ordersTechLevel,
				Wip25pctQty:     unit.wip[0],
				Wip50pctQty:     unit.wip[1],
				Wip75pctQty:     unit.wip[2],
			}
			err = q.CreateSCFactoryGroupUnit(e.Store.Context, fgUnitParms)
			if err != nil {
				log.Printf("create: empire: sc factory group unit %+v\n", fgUnitParms)
				return 0, err
			}
			// assumes FCT-1 for fuel and labor costs and actual production values
			if unit.techLevel != 1 {
				panic(fmt.Sprintf("assert(FCT-1 == FCT-%d)", unit.techLevel))
			}
			ps := sqlite.CreateSCFactoryProductionSummaryParams{
				ScID:          scId,
				GroupNo:       fg.groupNo,
				TurnNo:        0,
				FuelConsumed:  int64(math.Ceil(float64(unit.nbrOfUnits) * 0.5)),
				MetsConsumed:  0,
				NmtsConsumed:  0,
				UnitCd:        fg.ordersCode,
				UnitTechLevel: fg.ordersTechLevel,
			}
			if unit.nbrOfUnits < 5 {
				ps.ProConsumed = unit.nbrOfUnits * 6
				ps.UskConsumed = unit.nbrOfUnits * 18
			} else if unit.nbrOfUnits < 50 {
				ps.ProConsumed = unit.nbrOfUnits * 5
				ps.UskConsumed = unit.nbrOfUnits * 15
			} else if unit.nbrOfUnits < 500 {
				ps.ProConsumed = unit.nbrOfUnits * 4
				ps.UskConsumed = unit.nbrOfUnits * 12
			} else if unit.nbrOfUnits < 5_000 {
				ps.ProConsumed = unit.nbrOfUnits * 3
				ps.UskConsumed = unit.nbrOfUnits * 9
			} else if unit.nbrOfUnits < 50_000 {
				ps.ProConsumed = unit.nbrOfUnits * 2
				ps.UskConsumed = unit.nbrOfUnits * 6
			} else {
				ps.ProConsumed = unit.nbrOfUnits * 1
				ps.UskConsumed = unit.nbrOfUnits * 3
			}

			massConvertedPerYear := 20.0 * float64(unit.nbrOfUnits)

			metsPerUnit, nmtsPerUnit := 0.0, 0.0
			if fg.ordersCode == "AUT" {
				metsPerUnit, nmtsPerUnit = 2, 2
			} else if fg.ordersCode == "CNGD" {
				metsPerUnit, nmtsPerUnit = 0.2, 0.4
			} else if fg.ordersCode == "EWP" {
				metsPerUnit, nmtsPerUnit = 5, 5
			} else if fg.ordersCode == "MIN" {
				metsPerUnit, nmtsPerUnit = 6, 6
			} else if fg.ordersCode == "MTSP" {
				metsPerUnit, nmtsPerUnit = 0.02, 0.02
			} else if fg.ordersCode == "RSCH" {
				metsPerUnit, nmtsPerUnit = 0, 0
			} else if fg.ordersCode == "STU" {
				metsPerUnit, nmtsPerUnit = 0.07, 0.03
			} else {
				panic(fmt.Sprintf("assert(ordersCode != %q)", fg.ordersCode))
			}
			totalMassPerUnit := metsPerUnit + nmtsPerUnit
			if fg.ordersCode == "RSCH" {
				totalMassPerUnit = 1
			}
			unitsProducedPerYear := massConvertedPerYear / totalMassPerUnit
			unitsProducedPerTurn := math.Floor(unitsProducedPerYear / 4)

			ps.MetsConsumed = int64(math.Ceil(unitsProducedPerTurn * metsPerUnit))
			ps.NmtsConsumed = int64(math.Ceil(unitsProducedPerTurn * nmtsPerUnit))
			ps.UnitsProduced = int64(unitsProducedPerTurn)

			err = q.CreateSCFactoryProductionSummary(e.Store.Context, ps)
			if err != nil {
				log.Printf("create: empire: sc factory production summary %+v\n", ps)
				return 0, err
			}
		}
	}

	type farmGroupUnitData struct {
		code      string // always "FRM"
		techLevel int64  // tech level of the farming unit
		qty       int64  // number of units in the group
	}
	type farmGroupData struct {
		groupNo int64
		units   []farmGroupUnitData
	}
	for _, fg := range []farmGroupData{
		{groupNo: 1, units: []farmGroupUnitData{
			{code: "FRM", techLevel: 1, qty: 130_000}}},
	} {
		fgParms := sqlite.CreateSCFarmGroupParams{
			ScID:    scId,
			GroupNo: fg.groupNo,
		}
		err = q.CreateSCFarmGroup(e.Store.Context, fgParms)
		if err != nil {
			log.Printf("create: empire: sc farm group %+v\n", fgParms)
			return 0, err
		}
		for _, unit := range fg.units {
			fgUnitParms := sqlite.CreateSCFarmGroupUnitParams{
				ScID:           scId,
				GroupNo:        fg.groupNo,
				GroupTechLevel: unit.techLevel,
				NbrOfUnits:     unit.qty,
			}
			err = q.CreateSCFarmGroupUnit(e.Store.Context, fgUnitParms)
			if err != nil {
				log.Printf("create: empire: sc farm group unit %+v\n", fgUnitParms)
				return 0, err
			}
			// assumes FRM-1 for fuel and labor costs and actual production values
			if unit.techLevel != 1 {
				panic(fmt.Sprintf("assert(FRM-1 == FRM-%d)", unit.techLevel))
			}
			var ps = struct {
				groupNo      int64
				fuelConsumed int64
				proConsumed  int64
				uskConsumed  int64
				autConsumed  int64
				foodProduced int64
			}{
				groupNo:      fg.groupNo,
				fuelConsumed: int64(math.Ceil(float64(unit.qty) * 0.5)),
				proConsumed:  unit.qty,
				uskConsumed:  unit.qty * 3,
				autConsumed:  0,
				foodProduced: unit.qty * 100,
			}
			err = q.CreateSCFarmProductionSummary(e.Store.Context, sqlite.CreateSCFarmProductionSummaryParams{
				ScID:         scId,
				GroupNo:      ps.groupNo,
				TurnNo:       0,
				FuelConsumed: ps.fuelConsumed,
				ProConsumed:  ps.proConsumed,
				UskConsumed:  ps.uskConsumed,
				AutConsumed:  ps.autConsumed,
				FoodProduced: ps.foodProduced,
			})
			if err != nil {
				log.Printf("create: empire: sc farm production summary %+v\n", ps)
				return 0, err
			}
		}
	}

	// deposit maps the deposit no to the deposit id
	depositMap := map[int64]int64{}
	if rows, err := q.ReadDepositsByOrbit(e.Store.Context, sqlite.ReadDepositsByOrbitParams{
		OrbitID: gameRow.HomeOrbitID,
		TurnNo:  gameRow.CurrentTurn,
	}); err != nil {
		log.Printf("create: empire: sc orbit %d: %+v\n", gameRow.HomeOrbitID, err)
		return 0, err
	} else {
		for _, row := range rows {
			depositMap[row.DepositNo] = row.ID
		}
	}

	type miningUnit struct {
		code       string // always "MIN"
		techLevel  int64  // tech level of the mining unit
		nbrOfUnits int64  // number of units in the group
	}
	type miningGroup struct {
		groupNo   int64
		depositId int64
		depositNo int64
		units     []miningUnit
	}
	miningGroups := []miningGroup{
		{groupNo: 1, depositNo: 1, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 1_000}}},
		{groupNo: 2, depositNo: 2, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 40_817}}},
		{groupNo: 3, depositNo: 6, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 159_421}}},
	}
	if cfg.Username == "cortrah" {
		miningGroups = []miningGroup{
			{groupNo: 1, depositNo: 1, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 41_000}}},
			{groupNo: 2, depositNo: 2, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 40_817}}},
			{groupNo: 3, depositNo: 4, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 18_762}}},
			{groupNo: 4, depositNo: 6, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 159_421}}},
			{groupNo: 5, depositNo: 15, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20_000}}},
			{groupNo: 6, depositNo: 20, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20_000}}},
		}
	} else if cfg.Username == "jimw" {
		miningGroups = []miningGroup{
			{groupNo: 1, depositNo: 1, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20000}}},
			{groupNo: 2, depositNo: 2, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 90000}}},
			{groupNo: 3, depositNo: 3, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 4, depositNo: 4, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 5, depositNo: 5, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 6, depositNo: 6, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 159400}}},
			{groupNo: 7, depositNo: 7, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 8, depositNo: 8, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 9, depositNo: 9, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 10, depositNo: 10, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 11, depositNo: 11, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 12, depositNo: 12, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 13, depositNo: 13, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 14, depositNo: 14, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 15, depositNo: 15, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 16, depositNo: 16, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 17, depositNo: 17, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 18, depositNo: 18, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 19, depositNo: 19, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 10000}}},
			{groupNo: 20, depositNo: 20, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20000}}},
			{groupNo: 21, depositNo: 21, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 22, depositNo: 22, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 23, depositNo: 23, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 24, depositNo: 24, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 25, depositNo: 25, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 26, depositNo: 26, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 27, depositNo: 27, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 28, depositNo: 28, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 29, depositNo: 29, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 30, depositNo: 30, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 31, depositNo: 31, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 32, depositNo: 32, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 33, depositNo: 33, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 34, depositNo: 34, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
			{groupNo: 35, depositNo: 35, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 20}}},
		}
	} else if cfg.Username == "longshadow" {
		miningGroups = []miningGroup{
			{groupNo: 1, depositNo: 1, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 1000}}},
			{groupNo: 2, depositNo: 2, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 74000}}},
			{groupNo: 3, depositNo: 6, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 165000}}},
			{groupNo: 4, depositNo: 21, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 60000}}},
		}
	} else {
		miningGroups = []miningGroup{
			{groupNo: 1, depositNo: 1, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 1000}}},
			{groupNo: 2, depositNo: 2, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 74000}}},
			{groupNo: 3, depositNo: 6, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 165000}}},
			{groupNo: 4, depositNo: 21, units: []miningUnit{{code: "MIN", techLevel: 1, nbrOfUnits: 60000}}},
		}
	}
	for _, mg := range miningGroups {
		depositID, ok := depositMap[mg.depositNo]
		if !ok {
			log.Printf("create: empire: sc mining group %+v\n", mg)
			return 0, ErrMissingDeposit
		}
		mgParms := sqlite.CreateSCMiningGroupParams{
			ScID:      scId,
			GroupNo:   mg.groupNo,
			DepositID: depositID,
		}
		err = q.CreateSCMiningGroup(e.Store.Context, mgParms)
		if err != nil {
			log.Printf("create: empire: sc mining group %+v\n", mgParms)
			return 0, err
		}
		for _, unit := range mg.units {
			mgUnitParms := sqlite.CreateSCMiningGroupUnitParams{
				ScID:           scId,
				GroupNo:        mg.groupNo,
				GroupTechLevel: unit.techLevel,
				NbrOfUnits:     unit.nbrOfUnits,
			}
			err = q.CreateSCMiningGroupUnit(e.Store.Context, mgUnitParms)
			if err != nil {
				log.Printf("create: empire: sc mining group unit %+v\n", mgUnitParms)
				return 0, err
			}
			// assumes MIN-1 for fuel and labor costs and actual production values
			if unit.techLevel != 1 {
				panic(fmt.Sprintf("assert(MIN-1 == MIN-%d)", unit.techLevel))
			}
			deposit, ok := deposits[mg.depositNo]
			if !ok {
				log.Printf("create: empire: sc mining group unit %+v\n", mgUnitParms)
				return 0, ErrMissingDeposit
			}
			amtMinedPerTurn := unit.nbrOfUnits * 25
			amtProducedPerTurn := int64(math.Ceil(float64(amtMinedPerTurn) * deposit.yieldPct))
			if amtProducedPerTurn > deposit.qty {
				amtProducedPerTurn = deposit.qty
			}
			var ps = struct {
				groupNo      int64
				fuelConsumed int64
				proConsumed  int64
				uskConsumed  int64
				autConsumed  int64
				fuelProduced int64
				goldProduced int64
				metsProduced int64
				nmtsProduced int64
			}{
				groupNo:      mg.groupNo,
				fuelConsumed: int64(math.Ceil(float64(unit.nbrOfUnits) * 0.5)),
				proConsumed:  unit.nbrOfUnits,
				uskConsumed:  unit.nbrOfUnits * 3,
				autConsumed:  0,
			}
			switch deposit.kind {
			case "FUEL":
				ps.fuelProduced = amtProducedPerTurn
			case "GOLD":
				ps.goldProduced = amtProducedPerTurn
			case "METS":
				ps.metsProduced = amtProducedPerTurn
			case "NMTS":
				ps.nmtsProduced = amtProducedPerTurn
			}
			err = q.CreateSCMiningProductionSummary(e.Store.Context, sqlite.CreateSCMiningProductionSummaryParams{
				ScID:         scId,
				GroupNo:      ps.groupNo,
				TurnNo:       0,
				FuelConsumed: ps.fuelConsumed,
				ProConsumed:  ps.proConsumed,
				UskConsumed:  ps.uskConsumed,
				AutConsumed:  ps.autConsumed,
				FuelProduced: ps.fuelProduced,
				GoldProduced: ps.goldProduced,
				MetsProduced: ps.metsProduced,
				NmtsProduced: ps.nmtsProduced,
			})
			if err != nil {
				log.Printf("create: empire: sc farm production summary %+v\n", ps)
				return 0, err
			}
		}
	}

	// insert survey and probe orders to get reports for the empire started
	err = q.CreateSCSurveyOrder(e.Store.Context, sqlite.CreateSCSurveyOrderParams{ScID: scId, TargetID: gameRow.HomeOrbitID})
	if err != nil {
		log.Printf("create: empire: survey %v\n", err)
		return 0, err
	}
	err = q.CreateSCProbeOrder(e.Store.Context, sqlite.CreateSCProbeOrderParams{ScID: scId, Kind: "system", TargetID: gameRow.HomeSystemID})
	if err != nil {
		log.Printf("create: empire: probe system %v\n", err)
		return 0, err
	}
	err = q.CreateSCProbeOrder(e.Store.Context, sqlite.CreateSCProbeOrderParams{ScID: scId, Kind: "star", TargetID: gameRow.HomeStarID})
	if err != nil {
		log.Printf("create: empire: probe star %v\n", err)
		return 0, err
	}
	//

	return empireID, tx.Commit()
}

func codeTechLevelFromString(s string) (code string, techLevel int64, err error) {
	cd, tl, ok := strings.Cut(s, "-")
	if !ok {
		return s, 0, nil
	}
	code = cd
	if n, err := strconv.Atoi(tl); err != nil {
		log.Printf("create: empire: inventory: %q: techLevel %v\n", s, err)
	} else {
		techLevel = int64(n)
	}
	return code, techLevel, err
}
