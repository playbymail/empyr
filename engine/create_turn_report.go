// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/playbymail/empyr/pkg/stdlib"
	"github.com/playbymail/empyr/store/sqlc"
	"html/template"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"
)

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
	if row, err := e.Store.Queries.ReadGameInfoByCode(e.Store.Context, cfg.Code); err != nil {
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
	gameRow, err := e.Store.Queries.ReadGameInfoByCode(e.Store.Context, cfg.Code)
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, err
	}
	clusterRow, err := e.Store.Queries.ReadClusterMetaByGameID(e.Store.Context, gameRow.ID)
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, err
	}
	_ = clusterRow
	empireRow, err := e.Store.Queries.ReadEmpireByGameIDByID(e.Store.Context, sqlc.ReadEmpireByGameIDByIDParams{
		GameID:   gameRow.ID,
		EmpireNo: cfg.EmpireNo,
	})
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, err
	}
	log.Printf("game %d: empire %d: turn %d\n", empireRow.GameID, empireRow.EmpireNo, cfg.TurnNo)

	ts, err := template.New("turn-report").Parse(turnReportTmpl)
	if err != nil {
		return nil, err
	}

	payload := TurnReport_t{
		Heading: &ReportHeading_t{
			Game:       cfg.Code,
			TurnNo:     cfg.TurnNo,
			TurnCode:   fmt.Sprintf("T%05d", cfg.TurnNo),
			EmpireNo:   cfg.EmpireNo,
			EmpireCode: fmt.Sprintf("E%03d", cfg.EmpireNo),
		},
		CreatedDate:     time.Now().UTC().Format("2006-01-02"),
		CreatedDateTime: time.Now().UTC().Format(time.RFC3339),
	}

	colonyRows, err := e.Store.Queries.ReadAllColoniesByEmpire(e.Store.Context, empireRow.EmpireID)
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, err
	}
	for _, colonyRow := range colonyRows {
		log.Printf("colony: %d\n", colonyRow.SorcsID)
		colonyReport := &ColonyReport_t{
			Id:          colonyRow.SorcsID,
			Coordinates: fmt.Sprintf("%02d/%02d/%02d%s", colonyRow.X, colonyRow.Y, colonyRow.Z, colonyRow.Suffix),
			OrbitNo:     colonyRow.OrbitNo,
			Kind:        colonyRow.SorcKind,
			Name:        colonyRow.Name,
			VitalStatistics: &ColonyStatisticsReport_t{
				TechLevel:        colonyRow.TechLevel,
				StandardOfLiving: fmt.Sprintf("%6.4f", colonyRow.Sol),
				Rations:          fmt.Sprintf("%6.2f %%", colonyRow.Rations*100),
				BirthRate:        fmt.Sprintf("%6.4f %%", colonyRow.BirthRate),
				DeathRate:        fmt.Sprintf("%6.4f %%", colonyRow.DeathRate),
			},
		}
		if colonyReport.Name == "" {
			colonyReport.Name = "Not Named"
		}
		if popRows, err := e.Store.Queries.ReadSorCPopulation(e.Store.Context, colonyRow.SorcsID); err != nil {
			log.Printf("error: %v\n", err)
			return nil, err
		} else {
			colonyReport.Census = &ColonyCensusReport_t{}
			totalPopulation, totalPay := int64(0), int64(0)
			for _, popRow := range popRows {
				totPay := int64(math.Ceil(float64(popRow.Qty) * popRow.PayRate))
				colonyReport.Census.Population = append(colonyReport.Census.Population, &PopulationReport_t{
					Group:      popRow.PopulationCd,
					Population: commas(popRow.Qty),
					PayRate:    fmt.Sprintf("%6.4f", popRow.PayRate),
					TotalPay:   commas(totPay),
					qty:        popRow.Qty,
				})
				totalPopulation += popRow.Qty
				totalPay += totPay
			}
			colonyReport.Census.TotalPopulation = commas(totalPopulation)
			colonyReport.Census.TotalPay = commas(totalPay)
			for _, group := range colonyReport.Census.Population {
				group.PctTotalPop = fmt.Sprintf("%6.2f %%", 100*float64(group.qty)/float64(totalPopulation))
			}
		}
		// lie about the other statistics for now
		colonyReport.Other = &ColonyOtherReport_t{
			TotalMass:       "89,494,472",
			TotalVolume:     "60,000,000",
			AvailableVolume: "4,952,743",
		}
		// lie about transports, too
		colonyReport.Transports = &ColonyTransportReport_t{
			Capacity:  "400,000",
			Used:      "0",
			Available: "400,000",
		}
		if inventoryRows, err := e.Store.Queries.ReadSorCInventory(e.Store.Context, colonyRow.SorcsID); err != nil {
			log.Printf("error: %v\n", err)
			return nil, err
		} else {
			type inventoryLine_t struct {
				id              string
				code            string
				techLevel       int64
				nonAssemblyQty  int64
				assembledQty    int64
				disassembledQty int64
			}
			inventoryMap := map[string]*inventoryLine_t{}
			for _, item := range inventoryRows {
				var code string
				if item.TechLevel == 0 {
					code = item.UnitCd
				} else {
					code = fmt.Sprintf("%s-%d", item.UnitCd, item.TechLevel)
				}
				line, ok := inventoryMap[code]
				if !ok {
					line = &inventoryLine_t{id: code, code: item.UnitCd, techLevel: item.TechLevel}
					inventoryMap[code] = line
				}
				if IsOperational(line.code) {
					if item.IsAssembled == 1 {
						line.assembledQty += item.Qty
					} else { // assumes disassembled are in storage
						line.disassembledQty += item.Qty
					}
				} else { // assumes non-operational are always in storage
					line.nonAssemblyQty += item.Qty
				}
			}
			// sort the inventory lines by code and tech level
			var inventoryLines []*inventoryLine_t
			for _, line := range inventoryMap {
				inventoryLines = append(inventoryLines, line)
			}
			sort.Slice(inventoryLines, func(i, j int) bool {
				if inventoryLines[i].code == inventoryLines[j].code {
					return inventoryLines[i].techLevel < inventoryLines[j].techLevel
				}
				return inventoryLines[i].code < inventoryLines[j].code
			})
			for _, line := range inventoryLines {
				colonyReport.Inventory = append(colonyReport.Inventory, &ColonyInventoryLine_t{
					Code:            line.id,
					NonAssemblyQty:  commas(line.nonAssemblyQty),
					DisassembledQty: commas(line.disassembledQty),
					AssembledQty:    commas(line.assembledQty),
					IsOPU:           IsOperational(line.code),
				})
			}
		}
		if reportRows, err := e.Store.Queries.ReadReportProductionInputs(e.Store.Context, sqlc.ReadReportProductionInputsParams{
			SorcID: colonyRow.SorcsID,
			TurnNo: cfg.TurnNo,
		}); err != nil {
			log.Printf("error: %v\n", err)
			return nil, err
		} else {
			var totalFuel, totalGold, totalMetals, totalNonMetals int64
			for _, reportRow := range reportRows {
				totalFuel += reportRow.Fuel
				totalGold += reportRow.Gold
				totalMetals += reportRow.Metals
				totalNonMetals += reportRow.NonMetals
				colonyReport.ProductionConsumed = append(colonyReport.ProductionConsumed, &ProductionConsumedLine_t{
					Category:  reportRow.Category,
					Fuel:      commas(reportRow.Fuel),
					Gold:      commas(reportRow.Gold),
					Metals:    commas(reportRow.Metals),
					NonMetals: commas(reportRow.NonMetals),
				})
			}
			if len(reportRows) != 0 {
				colonyReport.ProductionConsumed = append(colonyReport.ProductionConsumed, &ProductionConsumedLine_t{
					Category:  "Total",
					Fuel:      commas(totalFuel),
					Gold:      commas(totalGold),
					Metals:    commas(totalMetals),
					NonMetals: commas(totalNonMetals),
				})
			}
		}
		if reportRows, err := e.Store.Queries.ReadReportProductionOutputs(e.Store.Context, sqlc.ReadReportProductionOutputsParams{
			SorcID: colonyRow.SorcsID,
			TurnNo: cfg.TurnNo,
		}); err != nil {
			log.Printf("error: %v\n", err)
			return nil, err
		} else {
			for _, reportRow := range reportRows {
				colonyReport.ProductionCreated = append(colonyReport.ProductionCreated, &ProductionCreatedLine_t{
					Category:     reportRow.Category,
					Farmed:       commas(reportRow.Farmed),
					Mined:        commas(reportRow.Mined),
					Manufactured: commas(reportRow.Manufactured),
				})
			}
		}
		if fgRows, err := e.Store.Queries.ReadSorCFactoryGroups(e.Store.Context, colonyRow.SorcsID); err != nil {
			log.Printf("error: %v\n", err)
		} else {
			for _, fgRow := range fgRows {
				var code string
				if fgRow.OrdersTechLevel == 0 {
					code = fgRow.OrdersCd
				} else {
					code = fmt.Sprintf("%s-%d", fgRow.OrdersCd, fgRow.OrdersTechLevel)
				}
				rpt := &ColonyFactoryGroupsReport_t{
					GroupNo: fmt.Sprintf("%02d", fgRow.GroupNo),
					Orders:  code,
				}
				if fgRow.RetoolTurnNo.Valid {
					rpt.Orders += " *"
					rpt.RetoolTurn = fmt.Sprintf("%d", fgRow.RetoolTurnNo.Int64)
				}
				if grpRows, err := e.Store.Queries.ReadSorCFactoryGroup(e.Store.Context, fgRow.GroupID); err != nil {
					log.Printf("error: %v\n", err)
				} else {
					for _, grpRow := range grpRows {
						if grpRow.OrdersTechLevel == 0 {
							code = grpRow.OrdersCd
						} else {
							code = fmt.Sprintf("%s-%d", grpRow.OrdersCd, grpRow.OrdersTechLevel)
						}
						rptLine := &ColonyFactoryGroupReport_t{
							TechLevel:  grpRow.FactoryTechLevel,
							NbrOfUnits: commas(grpRow.NbrOfUnits),
						}
						rptLine.Pipeline[0] = &ColonyFactoryPipelineReport_t{
							Percentage: "25%",
							Unit:       code,
							Qty:        commas(grpRow.Wip25pctQty),
						}
						rptLine.Pipeline[1] = &ColonyFactoryPipelineReport_t{
							Percentage: "50%",
							Unit:       code,
							Qty:        commas(grpRow.Wip50pctQty),
						}
						rptLine.Pipeline[2] = &ColonyFactoryPipelineReport_t{
							Percentage: "75%",
							Unit:       code,
							Qty:        commas(grpRow.Wip75pctQty),
						}
						rpt.Units = append(rpt.Units, rptLine)
					}
				}
				colonyReport.FactoryGroups = append(colonyReport.FactoryGroups, rpt)
			}
		}
		if fgRows, err := e.Store.Queries.ReadSorCFarmGroups(e.Store.Context, colonyRow.SorcsID); err != nil {
			log.Printf("error: %v\n", err)
		} else {
			for _, fgRow := range fgRows {
				colonyReport.FarmGroups = append(colonyReport.FarmGroups, &ColonyFarmGroupsReport_t{
					GroupNo:    fmt.Sprintf("%02d", fgRow.GroupNo),
					TechLevel:  fgRow.TechLevel,
					NbrOfUnits: commas(int64(fgRow.NbrOfUnits.Float64)),
				})
			}
		}
		if mgRows, err := e.Store.Queries.ReadSorCMiningGroups(e.Store.Context, colonyRow.SorcsID); err != nil {
			log.Printf("error: %v\n", err)
		} else {
			for _, mgRow := range mgRows {
				rpt := &ColonyMiningGroupsReport_t{
					GroupNo:      fmt.Sprintf("%02d", mgRow.GroupNo),
					DepositNo:    fmt.Sprintf("%02d", mgRow.DepositNo),
					DepositQty:   commas(mgRow.RemainingQty),
					DepositKind:  mgRow.DepositKind,
					DepositYield: fmt.Sprintf("%d %%", mgRow.YieldPct),
				}
				if grpRows, err := e.Store.Queries.ReadSorCMiningGroup(e.Store.Context, mgRow.GroupID); err != nil {
					log.Printf("error: %v\n", err)
				} else {
					for _, grpRow := range grpRows {
						rpt.Units = append(rpt.Units, &MiningGroupUnitReport_t{
							TechLevel:  grpRow.TechLevel,
							NbrOfUnits: commas(int64(grpRow.NbrOfUnits.Float64)),
						})
					}
				}
				colonyReport.MiningGroups = append(colonyReport.MiningGroups, rpt)
			}
		}

		payload.Colonies = append(payload.Colonies, colonyReport)
	}

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

	// buffer will hold the rendered turn report
	buffer := &bytes.Buffer{}

	// execute the template, writing the result to the buffer
	if err = ts.Execute(buffer, payload); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
