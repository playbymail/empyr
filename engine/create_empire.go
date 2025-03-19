// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"fmt"
	"github.com/playbymail/empyr/store/sqlc"
	"log"
)

type CreateEmpireParams_t struct {
	Code     string
	Username string
}

func CreateEmpireCommand(e *Engine_t, cfg *CreateEmpireParams_t) (int64, int64, error) {
	log.Printf("create: empire: code %q: handle %q\n", cfg.Code, cfg.Username)

	if cfg.Username == "" {
		return 0, 0, ErrMissingHandle
	} else if _, err := IsValidHandle(cfg.Username); err != nil {
		return 0, 0, err
	}

	q, tx, err := e.Store.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	var userID int64
	if row, err := q.ReadUserByUsername(e.Store.Context, cfg.Username); err != nil {
		return 0, 0, err
	} else {
		userID = row.ID
	}

	gameRow, err := q.ReadGameInfoByCode(e.Store.Context, cfg.Code)
	if err != nil {
		return 0, 0, err
	}
	if gameRow.CurrentTurn != 0 {
		return 0, 0, ErrGameInProgress
	}
	turnNo := gameRow.CurrentTurn
	clusterRow, err := q.ReadClusterMetaByGameID(e.Store.Context, gameRow.ID)
	if err != nil {
		return 0, 0, err
	}

	parms := sqlc.CreateEmpireParams{
		GameID:       gameRow.ID,
		UserID:       userID,
		EmpireNo:     gameRow.LastEmpireNo + 1,
		Name:         fmt.Sprintf("Empire %03d", gameRow.LastEmpireNo+1),
		HomeSystemID: clusterRow.HomeSystemID,
		HomeStarID:   clusterRow.HomeStarID,
		HomeOrbitID:  clusterRow.HomeOrbitID,
		HomePlanetID: clusterRow.HomePlanetID,
	}
	empireId, err := q.CreateEmpire(e.Store.Context, parms)
	if err != nil {
		return 0, 0, err
	}

	// create a home open surface colony
	sorcParams := sqlc.CreateSorCParams{
		EmpireID:    empireId,
		SorcCd:      "COPN",
		TechLevel:   1,
		Name:        "Not Named",
		OrbitID:     clusterRow.HomeOrbitID,
		IsOnSurface: 1,
		Rations:     1.0,
		Sol:         0.4881,
		BirthRate:   0,
		DeathRate:   0.0625,
	}
	sorcId, err := q.CreateSorC(e.Store.Context, sorcParams)
	if err != nil {
		log.Printf("create: empire: sorc %+v\n", sorcParams)
		return 0, 0, err
	}
	log.Printf("create: empire: id %d: no %d: colony %d\n", empireId, parms.EmpireNo, sorcId)

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
		sorcPopParms := sqlc.CreateSorCPopulationParams{
			SorcID:       sorcId,
			PopulationCd: pop.code,
			Qty:          pop.qty,
			PayRate:      pop.pay,
			RebelQty:     0,
		}
		err = q.CreateSorCPopulation(e.Store.Context, sorcPopParms)
		if err != nil {
			log.Printf("create: empire: sorc pop %+v\n", sorcPopParms)
			return 0, 0, err
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
		sorcInvParams := sqlc.CreateSorCInventoryParams{
			SorcID:    sorcId,
			UnitCd:    unit.code,
			TechLevel: unit.techLevel,
			Qty:       unit.qty,
			Mass:      Mass(unit.code, unit.techLevel, unit.qty),
		}
		if IsOperational(unit.code) {
			if unit.isAssembled {
				sorcInvParams.IsAssembled = 1
				sorcInvParams.Volume = VolumeAssembled(unit.code, unit.techLevel, unit.qty)
			} else {
				sorcInvParams.Volume = VolumeDisassembled(unit.code, unit.techLevel, unit.qty)
			}
		}
		if unit.isStored {
			sorcInvParams.IsStored = 1
			sorcInvParams.Volume = VolumeStored(unit.code, unit.techLevel, unit.qty)
		}
		err = q.CreateSorCInventory(e.Store.Context, sorcInvParams)
		if err != nil {
			log.Printf("create: empire: sorc inv %+v\n", sorcInvParams)
			return 0, 0, err
		}
	}

	// insert survey and probe orders
	err = q.CreateSorCSurveyOrder(e.Store.Context, sqlc.CreateSorCSurveyOrderParams{SorcID: sorcId, TurnNo: turnNo, TechLevel: 1, OrbitID: clusterRow.HomeOrbitID})
	if err != nil {
		log.Printf("create: empire: survey %v\n", err)
		return 0, 0, err
	}
	err = q.CreateSorCProbeOrder(e.Store.Context, sqlc.CreateSorCProbeOrderParams{SorcID: sorcId, TurnNo: turnNo, TechLevel: 1, Kind: "system", TargetID: clusterRow.HomeSystemID})
	if err != nil {
		log.Printf("create: empire: probe system %v\n", err)
		return 0, 0, err
	}
	err = q.CreateSorCProbeOrder(e.Store.Context, sqlc.CreateSorCProbeOrderParams{SorcID: sorcId, TurnNo: turnNo, TechLevel: 1, Kind: "star", TargetID: clusterRow.HomeStarID})
	if err != nil {
		log.Printf("create: empire: probe star %v\n", err)
		return 0, 0, err
	}

	// populate reports
	reportId, err := q.CreateReport(e.Store.Context, sqlc.CreateReportParams{
		SorcID: sorcId,
		TurnNo: turnNo,
	})
	if err != nil {
		log.Printf("create: empire: sorc prod %+v\n", err)
		return 0, 0, err
	}
	for _, category := range []struct {
		name      string
		fuel      int64
		gold      int64
		metals    int64
		nonMetals int64
	}{
		{"Farming", 65_000, 0, 0, 0},
		{"Mining", 150_000, 0, 0, 0},
		{"Manufacturing", 425_000, 0, 1_395_833, 1_354_167},
	} {
		_, err = q.CreateReportProductionInput(e.Store.Context, sqlc.CreateReportProductionInputParams{
			ReportID:  reportId,
			Category:  category.name,
			Fuel:      category.fuel,
			Gold:      category.gold,
			Metals:    category.metals,
			NonMetals: category.nonMetals,
		})
		if err != nil {
			log.Printf("create: empire: sorc prod %+v\n", err)
			return 0, 0, err
		}
	}
	for _, category := range []struct {
		name         string
		unitCode     string
		tech_level   int64
		farmed       int64
		mined        int64
		manufactured int64
	}{
		{name: "CNGD", unitCode: "CNGD", manufactured: 2_083_333},
		{name: "FOOD", unitCode: "FOOD", farmed: 3_250_000},
		{name: "FUEL", unitCode: "FUEL", mined: 500_000},
		{name: "GOLD", unitCode: "GOLD", mined: 0},
		{name: "METS", unitCode: "METS", mined: 2_750_000},
		{name: "NMTS", unitCode: "NMTS", mined: 0},
	} {
		_, err = q.CreateReportProductionOutput(e.Store.Context, sqlc.CreateReportProductionOutputParams{
			ReportID:     reportId,
			Category:     category.name,
			UnitCd:       category.unitCode,
			TechLevel:    category.tech_level,
			Farmed:       category.farmed,
			Mined:        category.mined,
			Manufactured: category.manufactured,
		})
		if err != nil {
			log.Printf("create: empire: sorc prod %+v\n", err)
			return 0, 0, err
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
		fgParms := sqlc.CreateSorCFarmGroupParams{
			SorcID:  sorcId,
			GroupNo: fg.groupNo,
		}
		fgID, err := q.CreateSorCFarmGroup(e.Store.Context, fgParms)
		if err != nil {
			log.Printf("create: empire: sorc farm group %+v\n", fgParms)
			return 0, 0, err
		}
		for _, unit := range fg.units {
			fgUnitParms := sqlc.CreateSorCFarmGroupUnitParams{
				FarmGroupID: fgID,
				UnitCd:      unit.code,
				TechLevel:   unit.techLevel,
				NbrOfUnits:  unit.qty,
			}
			_, err = q.CreateSorCFarmGroupUnit(e.Store.Context, fgUnitParms)
			if err != nil {
				log.Printf("create: empire: sorc farm group unit %+v\n", fgUnitParms)
				return 0, 0, err
			}
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
		fgParms := sqlc.CreateSorCFactoryGroupParams{
			SorcID:          sorcId,
			GroupNo:         fg.groupNo,
			OrdersCd:        fg.ordersCode,
			OrdersTechLevel: fg.ordersTechLevel,
		}
		fgID, err := q.CreateSorCFactoryGroup(e.Store.Context, fgParms)
		if err != nil {
			log.Printf("create: empire: sorc factory group %+v\n", fgParms)
			return 0, 0, err
		}
		for _, unit := range fg.units {
			fgUnitParms := sqlc.CreateSorCFactoryGroupUnitParams{
				FactoryGroupID:  fgID,
				UnitCd:          unit.code,
				TechLevel:       unit.techLevel,
				NbrOfUnits:      unit.nbrOfUnits,
				OrdersCd:        unit.ordersCode,
				OrdersTechLevel: unit.ordersTechLevel,
				Wip25pctQty:     unit.wip[0],
				Wip50pctQty:     unit.wip[1],
				Wip75pctQty:     unit.wip[2],
			}
			_, err = q.CreateSorCFactoryGroupUnit(e.Store.Context, fgUnitParms)
			if err != nil {
				log.Printf("create: empire: sorc factory group unit %+v\n", fgUnitParms)
				return 0, 0, err
			}
		}
	}

	// deposit maps the deposit no to the deposit id
	depositMap := map[int64]int64{}
	if rows, err := q.ReadDepositsByPlanet(e.Store.Context, clusterRow.HomePlanetID); err != nil {
		log.Printf("create: empire: sorc planet %d: %+v\n", clusterRow.HomePlanetID, err)
		return 0, 0, err
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
	if cfg.Username == "jimw" {
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
	}
	for _, mg := range miningGroups {
		depositID, ok := depositMap[mg.depositNo]
		if !ok {
			log.Printf("create: empire: sorc mining group %+v\n", mg)
			return 0, 0, fmt.Errorf("missing deposit")
		}
		mgParms := sqlc.CreateSorCMiningGroupParams{
			SorcID:    sorcId,
			GroupNo:   mg.groupNo,
			DepositID: depositID,
		}
		mgID, err := q.CreateSorCMiningGroup(e.Store.Context, mgParms)
		if err != nil {
			log.Printf("create: empire: sorc mining group %+v\n", mgParms)
			return 0, 0, err
		}
		for _, unit := range mg.units {
			mgUnitParms := sqlc.CreateSorCMiningGroupUnitParams{
				MiningGroupID: mgID,
				UnitCd:        unit.code,
				TechLevel:     unit.techLevel,
				NbrOfUnits:    unit.nbrOfUnits,
			}
			_, err = q.CreateSorCMiningGroupUnit(e.Store.Context, mgUnitParms)
			if err != nil {
				log.Printf("create: empire: sorc mining group unit %+v\n", mgUnitParms)
				return 0, 0, err
			}
		}
	}

	err = q.UpdateEmpireCounterByGameID(e.Store.Context, sqlc.UpdateEmpireCounterByGameIDParams{GameID: gameRow.ID, EmpireNo: parms.EmpireNo})
	if err != nil {
		return 0, 0, err
	}

	return empireId, parms.EmpireNo, tx.Commit()
}
