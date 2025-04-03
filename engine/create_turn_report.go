// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/playbymail/empyr/pkg/stdlib"
	"github.com/playbymail/empyr/repos/sqlite"
	"html/template"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type CreateTurnReportsParams_t struct {
	Path string // path to the output directory
}

// CreateTurnReportsCommand creates turn reports for all empires in the given game.
func CreateTurnReportsCommand(e *Engine_t, cfg *CreateTurnReportsParams_t) error {
	log.Printf("create: turn-reports: path %q\n", cfg.Path)
	if !stdlib.IsDirExists(cfg.Path) {
		log.Printf("error: reports path does not exist\n")
		return ErrInvalidPath
	}

	var gameCode string
	var turnNo int64
	if row, err := e.Store.Queries.ReadAllGameInfo(e.Store.Context); err != nil {
		return err
	} else {
		gameCode, turnNo = row.Code, row.CurrentTurn
	}
	log.Printf("game: %q: turn: %d\n", gameCode, turnNo)

	listOfEmpireID, err := e.Store.Queries.ReadActiveEmpires(e.Store.Context)
	if err != nil {
		log.Printf("error: %v\n", err)
		return err
	}

	// before we start, make sure the output directory exists for each empire
	errorCount := 0
	for _, empireID := range listOfEmpireID {
		empireReportPath := filepath.Join(cfg.Path, fmt.Sprintf("e%03d", empireID), "reports")
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
	for _, empireID := range listOfEmpireID {
		empireReportPath := filepath.Join(cfg.Path, fmt.Sprintf("e%03d", empireID), "reports")
		log.Printf("game: %d: turn: %d: empire %d (%d)\n", gameCode, turnNo, empireID, empireID)
		data, err := CreateTurnReportCommand(e, &CreateTurnReportParams_t{EmpireID: empireID})
		if err != nil {
			log.Printf("error: turn report: %v\n", err)
			errorCount++
			continue
		}
		reportName := filepath.Join(empireReportPath, fmt.Sprintf("e%03d-turn-%04d.html", empireID, turnNo))
		if err := os.WriteFile(reportName, data, 0644); err != nil {
			log.Printf("error: %q\n", reportName)
			log.Printf("error: os.WriteFile: %v\n", err)
			errorCount++
			continue
		}
		log.Printf("game: %d: turn: %d: empire %d (%d): created turn report\n", gameCode, turnNo, empireID, empireID)
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
	EmpireID int64 // empire number to create the turn report for
}

// CreateTurnReportCommand creates a turn report for a game.
// It returns a byte array containing the turn report as HTML.
func CreateTurnReportCommand(e *Engine_t, cfg *CreateTurnReportParams_t) ([]byte, error) {
	gameRow, err := e.Store.Queries.ReadAllGameInfo(e.Store.Context)
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, err
	}
	gameCode, turnNo := gameRow.Code, gameRow.CurrentTurn
	empireRow, err := e.Store.Queries.ReadEmpireByID(e.Store.Context, cfg.EmpireID)
	if err != nil {
		log.Printf("error: %v\n", err)
		return nil, err
	}
	log.Printf("game %d: empire %d: turn %d\n", gameCode, empireRow.EmpireID, turnNo)

	ts, err := template.New("turn-report").Parse(turnReportTmpl)
	if err != nil {
		return nil, err
	}

	payload := TurnReport_t{
		Heading: &ReportHeading_t{
			Game:       gameCode,
			TurnNo:     turnNo,
			TurnCode:   fmt.Sprintf("T%05d", turnNo),
			EmpireNo:   empireRow.EmpireID,
			EmpireCode: fmt.Sprintf("E%03d", empireRow.EmpireID),
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
		log.Printf("colony: %d\n", colonyRow.ScID)
		colonyReport := &ColonyReport_t{
			Id:          colonyRow.ScID,
			Coordinates: colonyRow.StarName,
			OrbitNo:     colonyRow.OrbitNo,
			Kind:        colonyRow.ScKind,
			Name:        colonyRow.Name,
			VitalStatistics: &ColonyStatisticsReport_t{
				TechLevel:        colonyRow.ScTechLevel,
				StandardOfLiving: fmt.Sprintf("%6.4f", colonyRow.Sol),
				Rations:          fmt.Sprintf("%6.2f %%", colonyRow.Rations*100),
				BirthRate:        fmt.Sprintf("%6.4f %%", colonyRow.BirthRate),
				DeathRate:        fmt.Sprintf("%6.4f %%", colonyRow.DeathRate),
			},
		}
		if colonyReport.Name == "" {
			colonyReport.Name = "Not Named"
		}
		if popRows, err := e.Store.Queries.ReadSCPopulation(e.Store.Context, colonyRow.ScID); err != nil {
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
		if inventoryRows, err := e.Store.Queries.ReadSCInventory(e.Store.Context, colonyRow.ScID); err != nil {
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
		//if reportRows, err := e.Store.Queries.ReadReportProductionInputs(e.Store.Context, sqlc.ReadReportProductionInputsParams{
		//	SorcID: colonyRow.ScID,
		//	TurnNo: cfg.TurnNo,
		//}); err != nil {
		//	log.Printf("error: %v\n", err)
		//	return nil, err
		//} else {
		//	var totalFuel, totalGold, totalMetals, totalNonMetals int64
		//	for _, reportRow := range reportRows {
		//		totalFuel += reportRow.Fuel
		//		totalGold += reportRow.Gold
		//		totalMetals += reportRow.Metals
		//		totalNonMetals += reportRow.NonMetals
		//		colonyReport.ProductionConsumed = append(colonyReport.ProductionConsumed, &ProductionConsumedLine_t{
		//			Category:  reportRow.Category,
		//			Fuel:      commas(reportRow.Fuel),
		//			Gold:      commas(reportRow.Gold),
		//			Metals:    commas(reportRow.Metals),
		//			NonMetals: commas(reportRow.NonMetals),
		//		})
		//	}
		//	if len(reportRows) != 0 {
		//		colonyReport.ProductionConsumed = append(colonyReport.ProductionConsumed, &ProductionConsumedLine_t{
		//			Category:  "Total",
		//			Fuel:      commas(totalFuel),
		//			Gold:      commas(totalGold),
		//			Metals:    commas(totalMetals),
		//			NonMetals: commas(totalNonMetals),
		//		})
		//	}
		//}
		//if reportRows, err := e.Store.Queries.ReadReportProductionOutputs(e.Store.Context, sqlc.ReadReportProductionOutputsParams{
		//	SorcID: colonyRow.ScID,
		//	TurnNo: cfg.TurnNo,
		//}); err != nil {
		//	log.Printf("error: %v\n", err)
		//	return nil, err
		//} else {
		//	for _, reportRow := range reportRows {
		//		colonyReport.ProductionCreated = append(colonyReport.ProductionCreated, &ProductionCreatedLine_t{
		//			Category:     reportRow.Category,
		//			Farmed:       commas(reportRow.Farmed),
		//			Mined:        commas(reportRow.Mined),
		//			Manufactured: commas(reportRow.Manufactured),
		//		})
		//	}
		//}
		if fgRows, err := e.Store.Queries.ReadSCFactoryGroups(e.Store.Context, colonyRow.ScID); err != nil {
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
				if fgRetool, err := e.Store.Queries.ReadSCFactoryGroupRetoolOrder(e.Store.Context, sqlc.ReadSCFactoryGroupRetoolOrderParams{
					ScID:    colonyRow.ScID,
					GroupNo: fgRow.GroupNo,
				}); err == nil {
					rpt.Orders += " *"
					rpt.RetoolTurn = fmt.Sprintf("%d", fgRetool.TurnNo)
				}
				if grpRows, err := e.Store.Queries.ReadSCFactoryGroup(e.Store.Context, sqlc.ReadSCFactoryGroupParams{
					ScID:    colonyRow.ScID,
					GroupNo: fgRow.GroupNo,
				}); err != nil {
					log.Printf("error: %v\n", err)
				} else {
					for _, grpRow := range grpRows {
						if grpRow.OrdersTechLevel == 0 {
							code = grpRow.OrdersCd
						} else {
							code = fmt.Sprintf("%s-%d", grpRow.OrdersCd, grpRow.OrdersTechLevel)
						}
						rptLine := &ColonyFactoryGroupReport_t{
							TechLevel:  grpRow.GroupTechLevel,
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
		if fgRows, err := e.Store.Queries.ReadSCFarmGroups(e.Store.Context, colonyRow.ScID); err != nil {
			log.Printf("error: %v\n", err)
		} else {
			for _, fgRow := range fgRows {
				colonyReport.FarmGroups = append(colonyReport.FarmGroups, &ColonyFarmGroupsReport_t{
					GroupNo:    fmt.Sprintf("%02d", fgRow.GroupNo),
					TechLevel:  fgRow.GroupTechLevel,
					NbrOfUnits: commas(int64(fgRow.NbrOfUnits.Float64)),
				})
			}
		}
		if mgRows, err := e.Store.Queries.ReadSCMiningGroups(e.Store.Context, sqlc.ReadSCMiningGroupsParams{
			ScID:   colonyRow.ScID,
			TurnNo: turnNo,
		}); err != nil {
			log.Printf("error: %v\n", err)
		} else {
			for _, mgRow := range mgRows {
				rpt := &ColonyMiningGroupsReport_t{
					GroupNo:      fmt.Sprintf("%02d", mgRow.GroupNo),
					DepositNo:    fmt.Sprintf("%02d", mgRow.DepositNo),
					DepositQty:   commas(mgRow.DepositQty),
					DepositKind:  mgRow.DepositKind,
					DepositYield: fmt.Sprintf("%d %%", mgRow.YieldPct),
				}
				if grpRows, err := e.Store.Queries.ReadSCMiningGroup(e.Store.Context, sqlc.ReadSCMiningGroupParams{
					ScID:    colonyRow.ScID,
					GroupNo: mgRow.GroupNo,
				}); err != nil {
					log.Printf("error: %v\n", err)
				} else {
					for _, grpRow := range grpRows {
						rpt.Units = append(rpt.Units, &MiningGroupUnitReport_t{
							TechLevel:  grpRow.GroupTechLevel,
							NbrOfUnits: commas(int64(grpRow.NbrOfUnits.Float64)),
						})
					}
				}
				colonyReport.MiningGroups = append(colonyReport.MiningGroups, rpt)
			}
		}

		payload.Colonies = append(payload.Colonies, colonyReport)
	}

	// buffer will hold the rendered turn report
	buffer := &bytes.Buffer{}

	// execute the template, writing the result to the buffer
	if err = ts.Execute(buffer, payload); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
