// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"fmt"
	"github.com/playbymail/empyr/store/sqlc"
	"log"
)

const (
	ErrEmpireNotAvailable = Error("empire not available")
	ErrMissingDeposit     = Error("missing deposit")
)

type CreateEmpireParams_t struct {
	EmpireID int64
	Username string
	Email    string
}

func CreateEmpireCommand(e *Engine_t, cfg *CreateEmpireParams_t) (int64, error) {
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
	err = q.UpdateEmpireStatus(e.Store.Context, sqlc.UpdateEmpireStatusParams{
		EmpireID: cfg.EmpireID,
		IsActive: 1,
	})
	if err != nil {
		return 0, err
	}
	parms := sqlc.CreateEmpireParams{
		EmpireID:     empireID,
		EmpireName:   fmt.Sprintf("Empire %03d", empireID),
		Username:     cfg.Username,
		Email:        cfg.Email,
		HomeSystemID: gameRow.HomeSystemID,
		HomeStarID:   gameRow.HomeStarID,
		HomeOrbitID:  gameRow.HomeOrbitID,
	}
	err = q.CreateEmpire(e.Store.Context, parms)
	if err != nil {
		return 0, err
	}

	// create a home open surface colony
	scParams := sqlc.CreateSCParams{
		EmpireID:    empireID,
		ScCd:        "COPN",
		ScTechLevel: 1,
		Name:        "Not Named",
		Location:    gameRow.HomeOrbitID,
		IsOnSurface: 1,
		Rations:     1.0,
		Sol:         0.4881,
		BirthRate:   0,
		DeathRate:   0.0625,
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
		scPopParms := sqlc.CreateSCPopulationParams{
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
		scInvParams := sqlc.CreateSCInventoryParams{
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
		fgParms := sqlc.CreateSCFactoryGroupParams{
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
			fgUnitParms := sqlc.CreateSCFactoryGroupUnitParams{
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
		fgParms := sqlc.CreateSCFarmGroupParams{
			ScID:    scId,
			GroupNo: fg.groupNo,
		}
		err = q.CreateSCFarmGroup(e.Store.Context, fgParms)
		if err != nil {
			log.Printf("create: empire: sc farm group %+v\n", fgParms)
			return 0, err
		}
		for _, unit := range fg.units {
			fgUnitParms := sqlc.CreateSCFarmGroupUnitParams{
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
		}
	}

	// deposit maps the deposit no to the deposit id
	depositMap := map[int64]int64{}
	if rows, err := q.ReadDepositsByOrbit(e.Store.Context, gameRow.HomeOrbitID); err != nil {
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
	}
	for _, mg := range miningGroups {
		depositID, ok := depositMap[mg.depositNo]
		if !ok {
			log.Printf("create: empire: sc mining group %+v\n", mg)
			return 0, ErrMissingDeposit
		}
		mgParms := sqlc.CreateSCMiningGroupParams{
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
			mgUnitParms := sqlc.CreateSCMiningGroupUnitParams{
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
		}
	}

	// todo: reports are going to need us to fake survey, probe, and production data from the previous turn
	//// insert survey and probe orders
	//err = q.CreateSCSurveyOrder(e.Store.Context, sqlc.CreateSCSurveyOrderParams{ScID: scId, TargetID: gameRow.HomeOrbitID})
	//if err != nil {
	//	log.Printf("create: empire: survey %v\n", err)
	//	return err
	//}
	//err = q.CreateSCProbeOrder(e.Store.Context, sqlc.CreateSCProbeOrderParams{ScID: scId, Kind: "system", TargetID: gameRow.HomeSystemID})
	//if err != nil {
	//	log.Printf("create: empire: probe system %v\n", err)
	//	return err
	//}
	//err = q.CreateSCProbeOrder(e.Store.Context, sqlc.CreateSCProbeOrderParams{ScID: scId, Kind: "star", TargetID: gameRow.HomeStarID})
	//if err != nil {
	//	log.Printf("create: empire: probe star %v\n", err)
	//	return err
	//}
	//

	return empireID, tx.Commit()
}
