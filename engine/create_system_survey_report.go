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
	"os"
	"path/filepath"
	"time"
)

type CreateSystemSurveyReportsParams_t struct {
	Code   string
	TurnNo int64
	Path   string // path to the output directory
}

// CreateSystemSurveyReportsCommand creates system survey reports for all empires in the given game.
func CreateSystemSurveyReportsCommand(e *Engine_t, cfg *CreateSystemSurveyReportsParams_t) error {
	log.Printf("create: system survey: code %q\n", cfg.Code)
	log.Printf("create: system survey: path %q\n", cfg.Path)
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
		empireSurveysPath := filepath.Join(cfg.Path, fmt.Sprintf("e%03d", empireNo), "surveys")
		if !stdlib.IsDirExists(empireSurveysPath) {
			log.Printf("error: empire survey path does not exist\n")
			log.Printf("error: %q\n", empireSurveysPath)
			errorCount++
		}
	}
	if errorCount > 0 {
		return ErrInvalidPath
	}

	// try to build out the reports
	for _, row := range rows {
		empireId, empireNo := row.EmpireID, row.EmpireNo
		empireSurveysPath := filepath.Join(cfg.Path, fmt.Sprintf("e%03d", empireNo), "surveys")
		log.Printf("game: %d: turn: %d: empire %d (%d)\n", gameId, turnNo, empireId, empireNo)
		data, err := CreateSystemSurveyReportCommand(e, &CreateSystemSurveyReportParams_t{Code: cfg.Code, TurnNo: cfg.TurnNo, EmpireNo: empireNo})
		if err != nil {
			log.Printf("error: turn report: %v\n", err)
			errorCount++
			continue
		}
		reportName := filepath.Join(empireSurveysPath, fmt.Sprintf("e%03d-turn-%04d.html", empireNo, turnNo))
		if err := os.WriteFile(reportName, data, 0644); err != nil {
			log.Printf("error: %q\n", reportName)
			log.Printf("error: os.WriteFile: %v\n", err)
			errorCount++
			continue
		}
		log.Printf("game: %d: turn: %d: empire %d (%d): created system survey report\n", gameId, turnNo, empireId, empireNo)
	}
	if errorCount > 0 {
		return ErrWritingReport
	}

	return nil
}

var (
	//go:embed templates/system-survey-report.gohtml
	surveySystemReportTmpl string
)

type CreateSystemSurveyReportParams_t struct {
	Code     string // code of the game to create the turn report for
	TurnNo   int64  // turn number to create the turn report for
	EmpireNo int64  // empire number to create the turn report for
}

// CreateSystemSurveyReportCommand creates a turn report for a game.
// It returns a byte array containing the turn report as HTML.
func CreateSystemSurveyReportCommand(e *Engine_t, cfg *CreateSystemSurveyReportParams_t) ([]byte, error) {
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

	ts, err := template.New("system-survey-report").Parse(surveySystemReportTmpl)
	if err != nil {
		return nil, err
	}

	payload := SystemSurveyReport_t{
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

	// buffer will hold the rendered turn report
	buffer := &bytes.Buffer{}

	// execute the template, writing the result to the buffer
	if err = ts.Execute(buffer, payload); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
