// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"context"
	"fmt"
	"github.com/playbymail/empyr/engine"
	"github.com/playbymail/empyr/repos"
	"github.com/playbymail/empyr/repos/sqlite"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
	"log"
	"path/filepath"
	"time"
)

// cmdExport represents the base command when called without any subcommands
var cmdExport = &cobra.Command{
	Use:   "export",
	Short: "export things",
	Long:  `export is the root of the export commands.`,
}

var cmdExportEmpires = &cobra.Command{
	Use:   "empires",
	Short: "export empires",
	Long:  `export empires data.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("export: empires: elapsed time: %v\n", time.Now().Sub(started))
		}()
		log.Printf("export: empires: game %q\n", flags.Game.Code)
		repo, err := repos.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}

		outputPath := cmd.Flags().Lookup("output").Value.String()

		activeEmpires, err := e.Store.Queries.ReadActiveEmpires(e.Store.Context)
		if err != nil {
			log.Fatalf("error: store.queries.read_active_empires: %v\n", err)
		}
		for _, empireID := range activeEmpires {
			empireRow, err := e.Store.Queries.ReadEmpireByID(e.Store.Context, empireID)
			if err != nil {
				log.Fatalf("error: readEmpireByID: %v\n", err)
			}
			gameCode, turnNo := empireRow.GameCode, empireRow.GameCurrentTurn

			pathXls := filepath.Join(outputPath, fmt.Sprintf("%s.t%05d.e%03d.xlsx", gameCode, turnNo, empireID))
			log.Printf("export: empire %d: %s\n", empireID, pathXls)

			f := excelize.NewFile()
			defer func() {
				if err := f.Close(); err != nil {
					fmt.Println(err)
				}
			}()

			if _, err := exportCoverTab(empireID, f, e.Store.Context, e.Store.Queries); err != nil {
				log.Fatalf("export: empire %d: %v\n", empireID, err)
			} else if _, err = exportSystemsTab(empireID, f, e.Store.Context, e.Store.Queries); err != nil {
				log.Fatalf("export: empire %d: %v\n", empireID, err)
			} else if _, err = exportStarProbesTab(empireID, turnNo, f, e.Store.Context, e.Store.Queries); err != nil {
				log.Fatalf("export: empire %d: %v\n", empireID, err)
			}

			// write the spreadsheet to the given path
			if err := f.SaveAs(pathXls); err != nil {
				log.Fatalf("error: excel.saveAs %v\n", err)
			}
			log.Printf("export: empire %d: saved to %s\n", empireID, pathXls)
		}
	},
}

// create the turn report cover sheet
func exportCoverTab(empireID int64, f *excelize.File, ctx context.Context, q *sqlite.Queries) (index int, err error) {
	const sheet = "Cover"
	index, err = f.NewSheet(sheet)
	if err != nil {
		log.Printf("export: sheet %q: %v\n", sheet, err)
		return index, err
	}
	f.SetActiveSheet(index)

	row, err := q.ExportCoverTabByID(ctx, empireID)
	if err != nil {
		log.Printf("export: sheet %q: %v\n", sheet, err)
		return index, err
	}

	_ = f.SetCellValue(sheet, "A1", "Game")
	_ = f.SetCellValue(sheet, "A2", row.GameCode)
	_ = f.SetCellValue(sheet, "B1", "Player")
	_ = f.SetCellValue(sheet, "B2", row.Username)
	_ = f.SetCellValue(sheet, "C1", "Empire")
	_ = f.SetCellValue(sheet, "C2", empireID)
	_ = f.SetCellValue(sheet, "D1", "Current Turn")
	_ = f.SetCellValue(sheet, "D2", row.GameCurrentTurn)
	_ = f.SetCellValue(sheet, "E1", "Home System")
	_ = f.SetCellValue(sheet, "E2", row.HomeSystemName)
	_ = f.SetCellValue(sheet, "F1", "Home Star")
	_ = f.SetCellValue(sheet, "F2", row.HomeStarName)
	_ = f.SetCellValue(sheet, "G1", "Home Planet")
	_ = f.SetCellValue(sheet, "G2", row.OrbitNo)
	_ = f.SetCellValue(sheet, "H1", "Home Planet Name")
	_ = f.SetCellValue(sheet, "H2", "not named")

	return index, nil
}

// create the turn report systems sheet
func exportSystemsTab(empireID int64, f *excelize.File, ctx context.Context, q *sqlite.Queries) (index int, err error) {
	const sheet = "Systems"
	index, err = f.NewSheet(sheet)
	if err != nil {
		log.Printf("export: sheet %q: %v\n", sheet, err)
		return index, err
	}
	f.SetActiveSheet(index)

	rowNo := 1 // heading row
	_ = f.SetCellValue(sheet, "A1", "X")
	_ = f.SetCellValue(sheet, "B1", "Y")
	_ = f.SetCellValue(sheet, "C1", "Z")
	_ = f.SetCellValue(sheet, "D1", "Nbr of Stars")

	rows, err := q.ExportSystems(ctx)
	if err != nil {
		log.Printf("export: sheet %q: %v\n", sheet, err)
		return index, err
	}

	for _, row := range rows {
		rowNo++
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNo), row.X)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNo), row.Y)
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNo), row.Z)
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNo), row.NbrOfStars)
	}

	return index, nil
}

// export the results of all star probes by the empire
func exportStarProbesTab(empireID, turnNo int64, f *excelize.File, ctx context.Context, q *sqlite.Queries) (index int, err error) {
	const sheet = "Star Probes"
	index, err = f.NewSheet(sheet)
	if err != nil {
		log.Printf("export: sheet %q: %v\n", sheet, err)
		return index, err
	}
	f.SetActiveSheet(index)

	rowNo := 1 // heading row
	_ = f.SetCellValue(sheet, "A1", "Star")

	rows, err := q.ExportStarProbes(ctx, sqlite.ExportStarProbesParams{EmpireID: empireID, TurnNo: turnNo})
	if err != nil {
		log.Printf("export: sheet %q: %v\n", sheet, err)
		return index, err
	}

	for _, row := range rows {
		rowNo++
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNo), row.StarName)
	}

	return index, nil
}
