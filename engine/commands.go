// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/playbymail/empyr/pkg/stdlib"
	"github.com/playbymail/empyr/store/sqlc"
	// "github.com/xuri/excelize/v2"
	"html/template"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"time"
)

// commands are the commands that can be issued to the engine.
// they should be implemented elsewhere, but this is convenient for now.

var (
	//go:embed templates/cluster-map.gohtml
	clusterMapTmpl string
)

const (
	ErrGameInProgress = Error("game in progress")
	ErrInvalidPath    = Error("invalid path")
	ErrWritingReport  = Error("error writing report")
)

type CreateClusterMapParams_t struct {
	Code string // code of the game to create the cluster map for
}

// CreateClusterMapCommand creates a cluster map.
// It returns a byte array containing the map as HTML.
func CreateClusterMapCommand(e *Engine_t, cfg *CreateClusterMapParams_t) ([]byte, error) {
	ts, err := template.New("cluster-map").Parse(clusterMapTmpl)
	if err != nil {
		return nil, err
	}

	type system_t struct {
		Id      int64
		X, Y, Z int64
		Color   template.JS
	}

	payload := struct {
		Game    string
		Systems []system_t
	}{
		Game: cfg.Code,
	}

	rows, err := e.Store.Queries.ReadClusterMapByGameCode(e.Store.Context, cfg.Code)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		var color template.JS
		switch row.NumberOfStars {
		case 1:
			color = "Blue"
		case 2:
			color = "Yellow"
		case 3:
			color = "White"
		case 4:
			color = "Red"
		default:
			return nil, fmt.Errorf("assert(s.NumberOfStars != %d)", row.NumberOfStars)
		}
		payload.Systems = append(payload.Systems, system_t{
			Id:    row.ID,
			X:     row.X - 15, // shift the origin back to 0,0,0
			Y:     row.Y - 15, // shift the origin back to 0,0,0
			Z:     row.Z - 15, // shift the origin back to 0,0,0
			Color: color,
		})
	}

	// buffer will hold the cluster map
	buffer := &bytes.Buffer{}

	// execute the template, writing the result to the buffer
	if err = ts.Execute(buffer, payload); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

var (
	//go:embed templates/cluster-star-list.gohtml
	clusterStarListTmpl string
)

type CreateClusterStarListParams_t struct {
	Code string // code of the game to create the cluster star list for
}

// CreateClusterStarListCommand creates a cluster star list for a game.
// It returns a byte array containing the star list as HTML and another
// byte array containing the star list as JSON.
func CreateClusterStarListCommand(e *Engine_t, cfg *CreateClusterStarListParams_t) ([]byte, []byte, error) {
	ts, err := template.New("cluster-star-list").Parse(clusterStarListTmpl)
	if err != nil {
		return nil, nil, err
	}

	// System ID</th><th>Coordinates</th><th>Number of Stars</th><th>Distance From Center
	type system_t struct {
		Id                 int64  `json:"id,omitempty"`
		X                  int64  `json:"x,omitempty"`
		Y                  int64  `json:"y,omitempty"`
		Z                  int64  `json:"z,omitempty"`
		Coordinates        string `json:"coordinates,omitempty"`
		NumberOfStars      int64  `json:"number-of-stars,omitempty"`
		DistanceFromCenter int64  `json:"distance-from-center,omitempty"`
	}

	payload := struct {
		Game        string
		UpdatedDate string
		Systems     []system_t
	}{
		Game:        cfg.Code,
		UpdatedDate: time.Now().UTC().Format("2006-01-02"),
	}

	rows, err := e.Store.Queries.ReadClusterMapByGameCode(e.Store.Context, cfg.Code)
	if err != nil {
		return nil, nil, err
	}
	for _, row := range rows {
		coordinates := fmt.Sprintf("(%02d, %02d, %02d)", row.X, row.Y, row.Z)
		dx, dy, dz := row.X-15, row.Y-15, row.Z-15
		distance := int64(math.Ceil(math.Sqrt(float64(dx*dx + dy*dy + dz*dz))))
		payload.Systems = append(payload.Systems, system_t{
			Id:                 row.ID,
			X:                  row.X,
			Y:                  row.Y,
			Z:                  row.Z,
			Coordinates:        coordinates,
			DistanceFromCenter: distance,
			NumberOfStars:      row.NumberOfStars,
		})
	}

	// buffer will hold the cluster star list
	buffer := &bytes.Buffer{}

	// execute the template, writing the result to the buffer
	if err = ts.Execute(buffer, payload); err != nil {
		return nil, nil, err
	}

	data, err := json.Marshal(payload.Systems)
	if err != nil {
		return nil, nil, err
	}

	return buffer.Bytes(), data, nil
}

type CreateEmpireParams_t struct {
	Code       string
	UserHandle string
}

func CreateEmpireCommand(e *Engine_t, cfg *CreateEmpireParams_t) (int64, int64, error) {
	log.Printf("create: empire: code %q\n", cfg.Code)

	if cfg.UserHandle == "" {
		return 0, 0, ErrMissingHandle
	} else if _, err := IsValidHandle(cfg.UserHandle); err != nil {
		return 0, 0, err
	}

	q, tx, err := e.Store.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	var userID int64
	if row, err := q.ReadUserByHandle(e.Store.Context, cfg.UserHandle); err != nil {
		return 0, 0, err
	} else {
		userID = row.ID
	}

	gameRow, err := q.ReadGameByCode(e.Store.Context, cfg.Code)
	if err != nil {
		return 0, 0, err
	}
	if gameRow.CurrentTurn != 0 {
		return 0, 0, ErrGameInProgress
	}
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
		{code: "UEM", qty: 59_000_000, pay: 0.0},
		{code: "USK", qty: 60_000_000, pay: 0.125},
		{code: "PRO", qty: 15_000_000, pay: 0.375},
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
		code string
		qty  int64
	}{
		{code: "STUN", qty: 60_000_000},
		{code: "FRM-1", qty: 130_000},
		{code: "MSL-1", qty: 50_000},
		{code: "SEN-1", qty: 20},
	} {
		sorcInvParams := sqlc.CreateSorCInventoryParams{
			SorcID: sorcId,
			UnitCd: unit.code,
			Qty:    unit.qty,
		}
		err = q.CreateSorCInventory(e.Store.Context, sorcInvParams)
		if err != nil {
			log.Printf("create: empire: sorc inv %+v\n", sorcInvParams)
			return 0, 0, err
		}
	}

	err = q.UpdateEmpireCounterByGameID(e.Store.Context, sqlc.UpdateEmpireCounterByGameIDParams{GameID: gameRow.ID, EmpireNo: parms.EmpireNo})
	if err != nil {
		return 0, 0, err
	}

	return empireId, parms.EmpireNo, tx.Commit()
}

type CreateGameParams_t struct {
	Code                        string
	Name                        string
	DisplayName                 string
	NumberOfEmpires             int64
	IncludeEmptyResources       bool
	PopulateSystemDistanceTable bool
	Rand                        *rand.Rand
	ForceCreate                 bool
}

// CreateGameCommand creates a new game.
func CreateGameCommand(e *Engine_t, cfg *CreateGameParams_t) (int64, error) {
	log.Printf("create: game: code %q: name %q: display %q\n", cfg.Code, cfg.Name, cfg.DisplayName)

	g, err := e.CreateGame(cfg.Code, cfg.Name, cfg.DisplayName, cfg.IncludeEmptyResources, cfg.PopulateSystemDistanceTable, cfg.Rand, cfg.ForceCreate)
	if err != nil {
		return 0, err
	}
	return g.Id, err
}

var (
	//go:embed templates/turn-report.gohtml
	turnReportTmpl string
)

type CreateTurnReportParams_t struct {
	Code     string // code of the game to create the turn report for
	TurnNo   int64  // turn number to create the turn report for
	EmpireNo int64  // empire number to create the turn report for
}

// CreateTurnReportCommand creates a turn report for a game.
// It returns a byte array containing the turn report as HTML.
func CreateTurnReportCommand(e *Engine_t, cfg *CreateTurnReportParams_t, path string) ([]byte, error) {
	panic("not implemented")
	//var gameId, homeSystemId, homeStarId, homePlanetId int64
	//if row, err := e.Store.Queries.ReadGameByCode(e.Store.Context, cfg.Code); err != nil {
	//	log.Printf("error: %v\n", err)
	//	return nil, err
	//} else {
	//	gameId, homeSystemId, homeStarId, homePlanetId = row.ID, row.HomeSystemID, row.HomeStarID, row.HomePlanetID
	//}
	//
	//var empireId int64
	//if row, err := e.Store.Queries.ReadGameEmpire(e.Store.Context, sqlc.ReadGameEmpireParams{GameCode: cfg.Code, EmpireNo: cfg.EmpireNo}); err != nil {
	//	log.Printf("error: %v\n", err)
	//	return nil, err
	//} else {
	//	gameId = row.GameID
	//	empireId = row.EmpireID
	//}
	//log.Printf("game %d: empire %d: turn %d\n", gameId, empireId, cfg.TurnNo)
	//
	//f := excelize.NewFile()
	//defer func() {
	//	if err := f.Close(); err != nil {
	//		fmt.Println(err)
	//	}
	//}()
	//pathXls := filepath.Join(path, fmt.Sprintf("%s.t%05d.e%03d.xlsx", cfg.Code, cfg.TurnNo, cfg.EmpireNo))
	//
	//// create the turn report cover sheet
	//sheet := "Cover"
	//if index, err := f.NewSheet(sheet); err != nil {
	//	log.Printf("error: %v\n", err)
	//	return nil, err
	//} else {
	//	f.SetActiveSheet(index)
	//	_ = f.SetCellValue(sheet, "A1", "Game")
	//	_ = f.SetCellValue(sheet, "B1", "Player")
	//	_ = f.SetCellValue(sheet, "C1", "Empire")
	//	_ = f.SetCellValue(sheet, "D1", "Current Turn")
	//	_ = f.SetCellValue(sheet, "E1", "Home System")
	//	_ = f.SetCellValue(sheet, "F1", "Home Star")
	//	_ = f.SetCellValue(sheet, "G1", "Home Planet")
	//	_ = f.SetCellValue(sheet, "A2", "A01")
	//	_ = f.SetCellValue(sheet, "B2", "jimw")
	//	_ = f.SetCellValue(sheet, "C2", 1)
	//	_ = f.SetCellValue(sheet, "D2", cfg.TurnNo)
	//	_ = f.SetCellValue(sheet, "E2", homeSystemId)
	//	_ = f.SetCellValue(sheet, "F2", homeStarId)
	//	_ = f.SetCellValue(sheet, "G2", homePlanetId)
	//}
	//
	//// create the systems sheet
	//sheet = "Systems"
	//if index, err := f.NewSheet(sheet); err != nil {
	//	log.Printf("error: %v\n", err)
	//	return nil, err
	//} else {
	//	f.SetActiveSheet(index)
	//	rows, err := e.Store.Queries.ReadAllSystems(e.Store.Context, gameId)
	//	if err != nil {
	//		log.Printf("error: %v\n", err)
	//		return nil, err
	//	}
	//	_ = f.SetCellValue(sheet, "A1", "System ID")
	//	_ = f.SetCellValue(sheet, "B1", "X")
	//	_ = f.SetCellValue(sheet, "C1", "Y")
	//	_ = f.SetCellValue(sheet, "D1", "Z")
	//	_ = f.SetCellValue(sheet, "E1", "Nbr of Stars")
	//	for n, row := range rows {
	//		rowNo := n + 2
	//		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNo), row.ID.Int64)
	//		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNo), row.X.Int64)
	//		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNo), row.Y.Int64)
	//		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNo), row.Z.Int64)
	//		_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNo), row.NumberOfStars)
	//	}
	//}
	//
	//// get all stars in the home system
	//var starIds []int64
	//if rows, err := e.Store.Queries.ReadStarsInSystem(e.Store.Context, homeSystemId); err != nil {
	//	log.Printf("error: %v\n", err)
	//	return nil, err
	//} else {
	//	starIds = rows
	//}
	//
	//// create the stars sheet
	//sheet = "Stars"
	//if index, err := f.NewSheet(sheet); err != nil {
	//	log.Printf("error: %v\n", err)
	//	return nil, err
	//} else {
	//	f.SetActiveSheet(index)
	//	_ = f.SetCellValue(sheet, "A1", "Star ID")
	//	_ = f.SetCellValue(sheet, "B1", "Coordinates")
	//	_ = f.SetCellValue(sheet, "C1", "Orbit")
	//	_ = f.SetCellValue(sheet, "D1", "Planet Type")
	//	_ = f.SetCellValue(sheet, "E1", "Deposit Type")
	//	_ = f.SetCellValue(sheet, "F1", "Deposit Qty")
	//	rowOffset := 2
	//	for _, starId := range starIds {
	//		rows, err := e.Store.Queries.ReadStarSurvey(e.Store.Context, starId)
	//		if err != nil {
	//			log.Printf("error: %v\n", err)
	//			return nil, err
	//		}
	//		for _, row := range rows {
	//			rowNo := rowOffset
	//			rowOffset++
	//			depositQty := int(math.Ceil(math.Log10(row.Quantity.Float64)))
	//			_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNo), row.StarID)
	//			_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNo), fmt.Sprintf("%02d/%02d/%02d%s", row.X, row.Y, row.Z, row.Sequence))
	//			_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNo), row.OrbitNo)
	//			_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNo), Orbit_e(row.OrbitKind).Code())
	//			_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNo), Resource_e(row.DepositKind).Code())
	//			_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNo), depositQty)
	//		}
	//	}
	//}
	//
	//// create the planet survey sheet
	//sheet = "Planet Survey"
	//if index, err := f.NewSheet(sheet); err != nil {
	//	log.Printf("error: %v\n", err)
	//	return nil, err
	//} else {
	//	f.SetActiveSheet(index)
	//	_ = f.SetCellValue(sheet, "A1", "Coordinates")
	//	_ = f.SetCellValue(sheet, "B1", "Orbit")
	//	_ = f.SetCellValue(sheet, "C1", "Planet Type")
	//	_ = f.SetCellValue(sheet, "D1", "Deposit No")
	//	_ = f.SetCellValue(sheet, "E1", "Deposit Type")
	//	_ = f.SetCellValue(sheet, "F1", "Deposit Qty")
	//	rowOffset := 2
	//	rows, err := e.Store.Queries.ReadPlanetSurvey(e.Store.Context, homePlanetId)
	//	if err != nil {
	//		log.Printf("error: %v\n", err)
	//		return nil, err
	//	}
	//	for _, row := range rows {
	//		rowNo := rowOffset
	//		rowOffset++
	//		depositQty := int(row.Quantity)
	//		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNo), fmt.Sprintf("%02d/%02d/%02d%s", row.X, row.Y, row.Z, row.Sequence))
	//		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNo), row.OrbitNo)
	//		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNo), Planet_e(row.PlanetKind).Code())
	//		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNo), row.DepositNo)
	//		_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNo), Resource_e(row.DepositKind).Code())
	//		_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNo), depositQty)
	//	}
	//}
	//
	//ts, err := template.New("turn-report").Parse(turnReportTmpl)
	//if err != nil {
	//	return nil, err
	//}
	//
	//type inventory_item_data_t struct {
	//	Code          string // the code (eg. "STU")
	//	TL            int64  // the tech level
	//	Qty           int64  // the quantity
	//	Mass          int64
	//	IsOperational bool
	//	InStorage     bool
	//	IsAssembled   bool
	//}
	//inventoryItemData := []inventory_item_data_t{
	//	// census
	//	{"UEM", 0, 59_000_000, 590_000, false, false, false},
	//	{"USK", 0, 60_000_000, 600_000, false, false, false},
	//	{"PRO", 0, 15_000_000, 150_000, false, false, false},
	//	{"SLD", 0, 2_500_000, 25_000, false, false, false},
	//	{"CNW", 0, 10_000, 200, false, false, false},
	//	{"SPY", 0, 20, 1, false, false, false},
	//	// storage/non-assembly
	//	{"ASWP", 1, 2_400_000, 2_400_000, false, true, false},
	//	{"ANM", 1, 150_000, 150_000, false, true, false},
	//	{"MTSP", 1, 150_000, 150_000, false, true, false},
	//	{"TPT", 1, 20_000, 20_000, false, true, false},
	//	{"METS", 0, 5_354_167, 5_354_167, false, true, false},
	//	{"NMTS", 0, 2_645_833, 2_645_833, false, true, false},
	//	{"FUEL", 0, 1_360_000, 1_360_000, false, true, false},
	//	{"GOLD", 0, 20_000, 20_000, false, true, false},
	//	{"FOOD", 0, 4_269_990, 4_269_990, false, true, false},
	//	{"CNGD", 0, 940_821, 940_821, false, true, false},
	//	// storage/unassembled
	//	{"AUT", 1, 93_750, 93_750, true, true, false},
	//	{"EWP", 1, 37_500, 37_500, true, true, false},
	//	{"MINU", 1, 31_250, 31_250, true, true, false},
	//	{"SEN", 1, 20, 20, true, true, false},
	//	{"STU", 0, 14_500_000, 14_500_000, true, true, false},
	//	// assembled items
	//	{"FRMU", 1, 130_000, 0, true, false, true},
	//	{"MSL", 1, 50_000, 0, true, false, true},
	//	{"SEN", 1, 20, 0, true, false, true},
	//	{"STU", 0, 60_000_000, 30_000_000, true, false, true},
	//}
	//
	//type fg_wip_data_t struct {
	//	Code            string // the code being produced (eg. "FRMU")
	//	TL              int64  // optional tech level of the item being produced
	//	UnitsInProgress int64  // number of units in the pipeline
	//}
	//
	//type factory_group_data_t struct {
	//	GroupNo    int64
	//	NbrOfUnits int64
	//	TL         int64
	//	Order      *fg_wip_data_t    // current order for the factory group
	//	WIP        [3]*fg_wip_data_t // newest to oldest, next to finish is always last
	//}
	//factoryGroupData := []factory_group_data_t{
	//	{GroupNo: 1, NbrOfUnits: 250_000, TL: 1,
	//		Order: &fg_wip_data_t{"CNGD", 0, 2_083_334},
	//		WIP: [3]*fg_wip_data_t{
	//			{"CNGD", 0, 2_083_333},
	//			{"CNGD", 0, 2_083_333},
	//			{"CNGD", 0, 2_083_333}}},
	//	{GroupNo: 2, NbrOfUnits: 75_000, TL: 1,
	//		Order: &fg_wip_data_t{"MTSP", 0, 9_375_000},
	//		WIP: [3]*fg_wip_data_t{
	//			{"MTSP", 0, 9_375_000},
	//			{"MTSP", 0, 9_375_000},
	//			{"MTSP", 0, 9_375_000}}},
	//	{GroupNo: 3, NbrOfUnits: 75_000, TL: 1,
	//		Order: &fg_wip_data_t{"AUT", 1, 93_750},
	//		WIP: [3]*fg_wip_data_t{
	//			{"AUT", 1, 93_750},
	//			{"AUT", 1, 93_750},
	//			{"AUT", 1, 93_750}}},
	//	{GroupNo: 4, NbrOfUnits: 75_000, TL: 1,
	//		Order: &fg_wip_data_t{"EWP", 1, 37_500},
	//		WIP: [3]*fg_wip_data_t{
	//			{"EWP", 1, 37_500},
	//			{"EWP", 1, 37_500},
	//			{"EWP", 1, 37_500}}},
	//	{GroupNo: 5, NbrOfUnits: 75_000, TL: 1,
	//		Order: &fg_wip_data_t{"MINU", 1, 31_250},
	//		WIP: [3]*fg_wip_data_t{
	//			{"MINU", 1, 31_250},
	//			{"MINU", 1, 31_250},
	//			{"MINU", 1, 31_250}}},
	//	{GroupNo: 6, NbrOfUnits: 250_000, TL: 1,
	//		Order: &fg_wip_data_t{"STU", 0, 12_500_000},
	//		WIP: [3]*fg_wip_data_t{
	//			{"STU", 0, 12_500_000},
	//			{"STU", 0, 12_500_000},
	//			{"STU", 0, 12_500_000}}},
	//	{GroupNo: 7, NbrOfUnits: 50_000, TL: 1,
	//		Order: &fg_wip_data_t{"RSCH", 0, 12_500},
	//		WIP: [3]*fg_wip_data_t{
	//			{"RSCH", 0, 12_500},
	//			{"RSCH", 0, 12_500},
	//			{"RSCH", 0, 12_500}}},
	//}
	//_ = factoryGroupData
	//
	//type mining_group_data_t struct {
	//	GroupNo    int64
	//	NbrOfUnits int64
	//	TL         int64
	//	DepositNo  int64
	//	DepositQty int64
	//	Resource   string
	//	YieldPct   int64
	//}
	//miningGroupData := []mining_group_data_t{
	//	{1, 100_000, 1, 13, 37_500_000, "FUEL", 20},
	//	{2, 200_000, 1, 28, 35_000_000, "METS", 55},
	//}
	//_ = miningGroupData
	//
	//type spy_results_data_t struct {
	//	Group   string
	//	Qty     int64
	//	Results []string
	//}
	//spyResultsData := []spy_results_data_t{
	//	spy_results_data_t{
	//		Group: "A", Qty: 10, Results: []string{
	//			"Report on the rebel situation:",
	//			"  Rebels         = 0",
	//			"  Rebel soldiers = 0",
	//		},
	//	},
	//	spy_results_data_t{
	//		Group: "B", Qty: 10, Results: []string{
	//			"Report on foreign espionage operations as you requested:",
	//			"  Owner: #   0:  Type A Spies:    0",
	//			"                 Type B Spies:    0",
	//			"                 Type C Spies:    0",
	//			"                 Type D Spies:    0",
	//			"                 Type E Spies:    0",
	//			"                 Type F Spies:    0",
	//		},
	//	},
	//}
	//
	//type pay_rates_t struct {
	//	USK     string
	//	PRO     string
	//	SLD     string
	//	Rations string
	//}
	//
	//type census_t struct {
	//	TotalPopulation int64
	//	UemQty          string
	//	UemPct          string
	//	UskQty          string
	//	UskPct          string
	//	ProQty          string
	//	ProPct          string
	//	SldQty          string
	//	SldPct          string
	//	CnwQty          string
	//	CnwPct          string
	//	SpyQty          string
	//	SpyPct          string
	//	BirthRate       string
	//	DeathRate       string
	//}
	//
	//// mass of population is 1 per 100 people
	//// hard code the inventory items for now
	//type inventory_item_t struct {
	//	Qty  string
	//	Code string // combined code with optional tech level
	//}
	//
	//type fg_wip_result_t struct {
	//	Code            string
	//	UnitsInProgress string
	//	PctComplete     string
	//}
	//type factory_group_result_t struct {
	//	GroupNo             string
	//	NbrOfUnits          string
	//	TL                  string
	//	OrderCode           string
	//	WIP25, WIP50, WIP75 fg_wip_result_t
	//}
	//
	//type mining_group_result_t struct {
	//	MineNo     int64
	//	NbrOfUnits string
	//	TL         int64
	//	DepositNo  int64
	//	DepositQty string
	//	Type       string
	//	YieldPct   string
	//}
	//
	//type spy_results_t struct {
	//	Group   string
	//	Qty     string
	//	Results []string
	//}
	//
	//type colony_t struct {
	//	Id                      int64
	//	Coordinates             string
	//	OrbitNo                 int64
	//	Name                    string
	//	Kind                    string
	//	TL                      int64
	//	SOL                     string
	//	Census                  *census_t
	//	PayRates                *pay_rates_t
	//	StorageNonAssemblyItems []*inventory_item_t
	//	StorageUnassembledItems []*inventory_item_t
	//	AssembledItems          []*inventory_item_t
	//	FactoryGroups           []*factory_group_result_t
	//	MiningGroups            []*mining_group_result_t
	//	InternalSpies           []*spy_results_t
	//}
	//
	//payload := struct {
	//	Game            string
	//	CreatedDate     string
	//	CreatedDateTime string
	//	TurnNo          int64
	//	TurnCode        string
	//	EmpireNo        int64
	//	EmpireCode      string
	//	Colonies        []*colony_t
	//}{
	//	Game:            cfg.Code,
	//	CreatedDate:     time.Now().UTC().Format("2006-01-02"),
	//	CreatedDateTime: time.Now().UTC().Format(time.RFC3339),
	//	TurnNo:          cfg.TurnNo,
	//	TurnCode:        fmt.Sprintf("T%05d", cfg.TurnNo),
	//	EmpireNo:        cfg.EmpireNo,
	//	EmpireCode:      fmt.Sprintf("E%03d", cfg.EmpireNo),
	//}
	//
	//colonyRows, err := e.Store.Queries.ReadEmpireAllColoniesForTurn(e.Store.Context, sqlc.ReadEmpireAllColoniesForTurnParams{EmpireID: empireId, TurnNo: cfg.TurnNo})
	//if err != nil {
	//	return nil, err
	//}
	//for _, row := range colonyRows {
	//	var kind SorC_e
	//	switch SorC_e(row.Kind) {
	//	case SCShip:
	//		kind = SCShip
	//	case SCOpenSurfaceColony:
	//		kind = SCOpenSurfaceColony
	//	case SCEnclosedSurfaceColony:
	//		kind = SCEnclosedSurfaceColony
	//	case SCOrbitalColony:
	//		kind = SCOrbitalColony
	//	default:
	//		panic(fmt.Sprintf("assert(sorc.Kind != %d)", row.Kind))
	//	}
	//	colony := &colony_t{
	//		Id:      row.SorcID.Int64,
	//		Name:    row.Name.String,
	//		Kind:    kind.Code(),
	//		OrbitNo: row.OrbitNo.Int64,
	//		TL:      row.TechLevel.Int64,
	//		SOL:     fmt.Sprintf("%5.4f", row.Sol.Float64),
	//		Census: &census_t{
	//			BirthRate: fmt.Sprintf("%5.4f", row.BirthRate.Float64),
	//			DeathRate: fmt.Sprintf("%5.4f", row.DeathRate.Float64),
	//			UemQty:    commas(row.UemQty.Int64),
	//			UskQty:    commas(row.UskQty.Int64),
	//			ProQty:    commas(row.ProQty.Int64),
	//			SldQty:    commas(row.SldQty.Int64),
	//			CnwQty:    commas(row.CnwQty.Int64),
	//			SpyQty:    commas(row.SpyQty.Int64),
	//		},
	//		PayRates: &pay_rates_t{
	//			USK:     fmt.Sprintf("%5.4f", row.UskPay.Float64),
	//			PRO:     fmt.Sprintf("%5.4f", row.ProPay.Float64),
	//			SLD:     fmt.Sprintf("%5.4f", row.SldPay.Float64),
	//			Rations: fmt.Sprintf("%5.4f", row.Rations.Float64),
	//		},
	//	}
	//	colony.Coordinates = fmt.Sprintf("%02d/%02d/%02d%s", row.X.Int64, row.Y.Int64, row.Z.Int64, row.Suffix.String)
	//	switch kind {
	//	case SCOpenSurfaceColony:
	//		colony.Kind = "Open Colony"
	//	case SCEnclosedSurfaceColony:
	//		colony.Kind = "Enclosed Colony"
	//	case SCOrbitalColony:
	//		colony.Kind = "Orbital Colony"
	//	}
	//	colony.Census.TotalPopulation = row.UemQty.Int64 + row.UskQty.Int64 + row.ProQty.Int64 + row.SldQty.Int64 + 2*row.CnwQty.Int64 + 2*row.SpyQty.Int64
	//	colony.Census.UemPct = fmt.Sprintf("%7.4f%%", float64(row.UemQty.Int64)/float64(colony.Census.TotalPopulation)*100)
	//	colony.Census.UskPct = fmt.Sprintf("%7.4f%%", float64(row.UskQty.Int64)/float64(colony.Census.TotalPopulation)*100)
	//	colony.Census.ProPct = fmt.Sprintf("%7.4f%%", float64(row.ProQty.Int64)/float64(colony.Census.TotalPopulation)*100)
	//	colony.Census.SldPct = fmt.Sprintf("%7.4f%%", float64(row.SldQty.Int64)/float64(colony.Census.TotalPopulation)*100)
	//	colony.Census.CnwPct = fmt.Sprintf("%7.4f%%", float64(row.CnwQty.Int64)/float64(colony.Census.TotalPopulation)*100)
	//	colony.Census.SpyPct = fmt.Sprintf("%7.4f%%", float64(row.SpyQty.Int64)/float64(colony.Census.TotalPopulation)*100)
	//
	//	for _, item := range inventoryItemData {
	//		var code string
	//		if item.TL == 0 {
	//			code = item.Code
	//		} else {
	//			code = fmt.Sprintf("%s-%d", item.Code, item.TL)
	//		}
	//		if item.InStorage && !item.IsOperational {
	//			colony.StorageNonAssemblyItems = append(colony.StorageNonAssemblyItems, &inventory_item_t{
	//				Qty:  commas(item.Qty),
	//				Code: code,
	//			})
	//		} else if item.InStorage && item.IsOperational {
	//			colony.StorageUnassembledItems = append(colony.StorageUnassembledItems, &inventory_item_t{
	//				Qty:  commas(item.Qty),
	//				Code: code,
	//			})
	//		} else if item.IsOperational && item.IsAssembled {
	//			colony.AssembledItems = append(colony.AssembledItems, &inventory_item_t{
	//				Qty:  commas(item.Qty),
	//				Code: code,
	//			})
	//		}
	//	}
	//
	//	for _, item := range factoryGroupData {
	//		fg := &factory_group_result_t{
	//			GroupNo:    fmt.Sprintf("%5d", item.GroupNo),
	//			NbrOfUnits: fmt.Sprintf("%14s", commas(item.NbrOfUnits)),
	//			TL:         fmt.Sprintf("%2d", item.TL),
	//			OrderCode:  fmt.Sprintf("%-8s", codeTL(item.Order.Code, item.Order.TL)),
	//		}
	//		if wip := item.WIP[0]; wip != nil {
	//			fg.WIP25.Code = fmt.Sprintf("%-8s", codeTL(item.Order.Code, item.Order.TL))
	//			fg.WIP25.UnitsInProgress = fmt.Sprintf("%14s", commas(wip.UnitsInProgress))
	//		}
	//		if wip := item.WIP[1]; wip != nil {
	//			fg.WIP50.Code = fmt.Sprintf("%-8s", codeTL(item.Order.Code, item.Order.TL))
	//			fg.WIP50.UnitsInProgress = fmt.Sprintf("%14s", commas(wip.UnitsInProgress))
	//		}
	//		if wip := item.WIP[2]; wip != nil {
	//			fg.WIP75.Code = fmt.Sprintf("%-8s", codeTL(item.Order.Code, item.Order.TL))
	//			fg.WIP75.UnitsInProgress = fmt.Sprintf("%14s", commas(wip.UnitsInProgress))
	//		}
	//		colony.FactoryGroups = append(colony.FactoryGroups, fg)
	//	}
	//
	//	for _, item := range miningGroupData {
	//		mg := &mining_group_result_t{
	//			MineNo:     item.GroupNo,
	//			NbrOfUnits: commas(item.NbrOfUnits),
	//			TL:         item.TL,
	//			DepositNo:  item.DepositNo,
	//			DepositQty: commas(item.DepositQty),
	//			Type:       item.Resource,
	//			YieldPct:   fmt.Sprintf("%d %%", item.YieldPct),
	//		}
	//		colony.MiningGroups = append(colony.MiningGroups, mg)
	//	}
	//
	//	for _, item := range spyResultsData {
	//		results := &spy_results_t{
	//			Group:   item.Group,
	//			Qty:     commas(item.Qty),
	//			Results: item.Results,
	//		}
	//		colony.InternalSpies = append(colony.InternalSpies, results)
	//	}
	//
	//	payload.Colonies = append(payload.Colonies, colony)
	//}
	//
	//// Save spreadsheet by the given path.
	//if err := f.SaveAs(pathXls); err != nil {
	//	log.Printf("create turn report: error: %v\n", err)
	//	return nil, err
	//}
	//
	//// buffer will hold the rendered turn report
	//buffer := &bytes.Buffer{}
	//
	//// execute the template, writing the result to the buffer
	//if err = ts.Execute(buffer, payload); err != nil {
	//	return nil, err
	//}
	//
	//return buffer.Bytes(), nil
}

type CreateTurnReportsParams_t struct {
	Code   string
	TurnNo int64
	Path   string // path to the output directory
}

// CreateTurnReportsCommand creates turn reports for all empires in the given game.
func CreateTurnReportsCommand(e *Engine_t, cfg *CreateTurnReportsParams_t) error {
	log.Printf("create: turn-reports: code %q\n", cfg.Code)
	log.Printf("create: turn-reports: path %q\n", cfg.Path)
	if !stdlib.IsDirExists(cfg.Path) {
		log.Printf("error: reports path does not exist\n")
		return ErrInvalidPath
	}

	var gameId, turnNo int64
	if row, err := e.Store.Queries.ReadGameByCode(e.Store.Context, cfg.Code); err != nil {
		return err
	} else {
		gameId, turnNo = row.ID, row.CurrentTurn
	}
	log.Printf("game: %d: turn: %d\n", gameId, turnNo)

	rows, err := e.Store.Queries.ReadAllEmpiresByGameID(e.Store.Context, gameId)
	if err != nil {
		log.Printf("error: %v\n", err)
		return err
	}

	// before we start, make sure the output directory exists for each empire
	errorCount := 0
	for _, row := range rows {
		empireNo := row.EmpireNo
		empireReportPath := filepath.Join(cfg.Path, fmt.Sprintf("e%03d", empireNo), "reports")
		if !stdlib.IsDirExists(empireReportPath) {
			log.Printf("error: empire report path does not exist\n")
			log.Printf("error: %q\n", empireReportPath)
			errorCount++
		}
	}
	if errorCount > 0 {
		return ErrInvalidPath
	}

	// try to build out the reports
	for _, row := range rows {
		empireId, empireNo := row.EmpireID, row.EmpireNo
		empireReportPath := filepath.Join(cfg.Path, fmt.Sprintf("e%03d", empireNo), "reports")
		log.Printf("game: %d: turn: %d: empire %d (%d)\n", gameId, turnNo, empireId, empireNo)
		data, err := CreateTurnReportCommand(e, &CreateTurnReportParams_t{
			Code:     cfg.Code,
			TurnNo:   cfg.TurnNo,
			EmpireNo: empireNo,
		}, empireReportPath)
		if err != nil {
			log.Printf("error: turn report: %v\n", err)
			errorCount++
			continue
		}
		reportName := filepath.Join(empireReportPath, fmt.Sprintf("e%03d-turn-%04d.html", empireNo, turnNo))
		if err := os.WriteFile(reportName, data, 0644); err != nil {
			log.Printf("error: %q\n", reportName)
			log.Printf("error: os.WriteFile: %v\n", err)
			errorCount++
			continue
		}
		log.Printf("game: %d: turn: %d: empire %d (%d): created turn report\n", gameId, turnNo, empireId, empireNo)
	}
	if errorCount > 0 {
		return ErrWritingReport
	}

	return nil
}

type DeleteGameParams_t struct {
	Code string
}

// DeleteGameCommand deletes an existing game..
func DeleteGameCommand(e *Engine_t, cfg *DeleteGameParams_t) error {
	log.Printf("delete: game: code %q\n", cfg.Code)
	return e.DeleteGame(cfg.Code)
}

func codeTL(code string, tl int64) string {
	if tl == 0 {
		return code
	}
	return fmt.Sprintf("%s-%d", code, tl)
}

func commas(n int64) string {
	in := fmt.Sprintf("%d", n)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}
