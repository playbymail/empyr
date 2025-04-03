// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"fmt"
	"github.com/playbymail/empyr/repos/sqlite"
	"github.com/xuri/excelize/v2"
	"log"
	"math"
	"path/filepath"
)

type CreateExcelReportParams_t struct {
	Code     string // code of the game to create the turn report for
	EmpireID int64  // empire number to create the turn report for
}

// CreateExcelReportCommand creates a turn report for a game.
// It returns a byte array containing the turn report as HTML.
func CreateExcelReportCommand(e *Engine_t, cfg *CreateExcelReportParams_t, path string) error {
	gameRow, err := e.Store.Queries.ReadAllGameInfo(e.Store.Context)
	if err != nil {
		log.Printf("error: %v\n", err)
		return err
	}
	empireRow, err := e.Store.Queries.ReadEmpireByID(e.Store.Context, cfg.EmpireID)
	if err != nil {
		log.Printf("error: %v\n", err)
		return err
	}
	log.Printf("game %q: empire %d: turn %d\n", gameRow.Code, empireRow.EmpireID, gameRow.CurrentTurn)

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	pathXls := filepath.Join(path, fmt.Sprintf("%s.t%05d.e%03d.xlsx", cfg.Code, gameRow.CurrentTurn, empireRow.EmpireID))

	// create the turn report cover sheet
	sheet := "Cover"
	if index, err := f.NewSheet(sheet); err != nil {
		log.Printf("error: %v\n", err)
		return err
	} else {
		f.SetActiveSheet(index)
		_ = f.SetCellValue(sheet, "A1", "Game")
		_ = f.SetCellValue(sheet, "A2", gameRow.Code)
		_ = f.SetCellValue(sheet, "B1", "Player")
		_ = f.SetCellValue(sheet, "B2", empireRow.Username)
		_ = f.SetCellValue(sheet, "C1", "Empire")
		_ = f.SetCellValue(sheet, "C2", empireRow.EmpireID)
		_ = f.SetCellValue(sheet, "D1", "Current Turn")
		_ = f.SetCellValue(sheet, "D2", gameRow.CurrentTurn)
		_ = f.SetCellValue(sheet, "E1", "Home System")
		_ = f.SetCellValue(sheet, "E2", empireRow.HomeSystemID)
		_ = f.SetCellValue(sheet, "F1", "Home Star")
		_ = f.SetCellValue(sheet, "F2", empireRow.HomeStarID)
		_ = f.SetCellValue(sheet, "G1", "Home Orbit")
		_ = f.SetCellValue(sheet, "G2", empireRow.HomeOrbitID)
	}

	// create the systems sheet
	sheet = "Systems"
	if index, err := f.NewSheet(sheet); err != nil {
		log.Printf("error: %v\n", err)
		return err
	} else {
		f.SetActiveSheet(index)
		rows, err := e.Store.Queries.ReadAllSystems(e.Store.Context)
		if err != nil {
			log.Printf("error: %v\n", err)
			return err
		}
		_ = f.SetCellValue(sheet, "A1", "System ID")
		_ = f.SetCellValue(sheet, "B1", "X")
		_ = f.SetCellValue(sheet, "C1", "Y")
		_ = f.SetCellValue(sheet, "D1", "Z")
		_ = f.SetCellValue(sheet, "E1", "Nbr of Stars")
		rowNo := 2
		for _, row := range rows {
			_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNo), row.SystemID)
			_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNo), row.X)
			_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNo), row.Y)
			_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNo), row.Z)
			_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNo), row.NbrOfStars)
			rowNo++
		}
	}

	// create the stars sheet
	sheet = "Stars"
	if index, err := f.NewSheet(sheet); err != nil {
		log.Printf("error: %v\n", err)
		return err
	} else {
		f.SetActiveSheet(index)
		_ = f.SetCellValue(sheet, "A1", "Star ID")
		_ = f.SetCellValue(sheet, "B1", "Coordinates")
		_ = f.SetCellValue(sheet, "C1", "Orbit")
		_ = f.SetCellValue(sheet, "D1", "Planet Type")
		_ = f.SetCellValue(sheet, "E1", "Deposit Type")
		_ = f.SetCellValue(sheet, "F1", "Deposit Qty")

		// get all stars in the home system
		starRows, err := e.Store.Queries.ReadAllStarsInSystem(e.Store.Context, empireRow.HomeSystemID)
		if err != nil {
			log.Printf("error: %v\n", err)
			return err
		}

		rowNo := 2
		for _, starRow := range starRows {
			coords := fmt.Sprintf("%02d/%02d/%02d%s", starRow.X, starRow.Y, starRow.Z, starRow.Sequence)
			rows, err := e.Store.Queries.ReadStarSurvey(e.Store.Context, starRow.ID)
			if err != nil {
				log.Printf("error: %v\n", err)
				return err
			}
			for _, row := range rows {
				depositQty := int(math.Ceil(math.Log10(row.Quantity.Float64)))
				var orbitKind string
				switch row.OrbitKind {
				case "ASTR":
					orbitKind = "Asteroid Belt"
				case "ERTH", "RCKY":
					orbitKind = "Rocky Planet"
				case "GASG", "ICEG":
					orbitKind = "Gas Giant"
				default:
					orbitKind = "Unknown"
				}
				_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNo), starRow.ID)
				_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNo), coords)
				_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNo), row.OrbitNo)
				_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNo), orbitKind)
				_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNo), row.DepositKind)
				_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNo), depositQty)
				rowNo++
			}
		}
	}

	// create the planet survey sheet
	sheet = "Planet Survey"
	if index, err := f.NewSheet(sheet); err != nil {
		log.Printf("error: %v\n", err)
		return err
	} else {
		f.SetActiveSheet(index)
		_ = f.SetCellValue(sheet, "A1", "Coordinates")
		_ = f.SetCellValue(sheet, "B1", "Orbit")
		_ = f.SetCellValue(sheet, "C1", "Planet Type")
		_ = f.SetCellValue(sheet, "D1", "Deposit No")
		_ = f.SetCellValue(sheet, "E1", "Deposit Type")
		_ = f.SetCellValue(sheet, "F1", "Deposit Qty")
		_ = f.SetCellValue(sheet, "G1", "Yield Pct")
		rowNo := 2
		rows, err := e.Store.Queries.ReadOrbitSurvey(e.Store.Context, sqlc.ReadOrbitSurveyParams{
			OrbitID: empireRow.HomeOrbitID,
			TurnNo:  empireRow.GameCurrentTurn,
		})
		if err != nil {
			log.Printf("error: %v\n", err)
			return err
		}
		for _, row := range rows {
			_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNo), row.OrbitNo)
			_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNo), row.OrbitKind)
			_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNo), row.DepositNo)
			_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNo), row.DepositKind)
			_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNo), row.DepositQty)
			_ = f.SetCellValue(sheet, fmt.Sprintf("G%d", rowNo), float64(row.YieldPct)/100.0)
			rowNo++
		}
	}

	// Save spreadsheet by the given path.
	if err := f.SaveAs(pathXls); err != nil {
		log.Printf("create excel report: error: %v\n", err)
		return err
	}
	log.Printf("create excel report: saved to %s\n", pathXls)

	return nil
}
