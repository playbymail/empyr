// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"database/sql"
	"errors"
	"github.com/playbymail/empyr/store/sqlc"
	"log"
)

func (e *Engine_t) ExecuteSurveys(gameCode string) error {
	// get the game id and turn number
	var gameID, turnNo int64
	if row, err := e.Store.Queries.ReadGameInfoByCode(e.Store.Context, gameCode); err != nil {
		return err
	} else {
		gameID, turnNo = row.ID, row.CurrentTurn
	}
	//log.Printf("game %q: id %d: turn: %d\n", gameCode, gameID, turnNo)

	// get a list of all the survey orders for this turn. these are the
	// orders that need to be executed.
	rows, err := e.Store.Queries.ReadAllSurveyOrdersForGameTurn(e.Store.Context, sqlc.ReadAllSurveyOrdersForGameTurnParams{GameID: gameID, TurnNo: turnNo})
	if err != nil {
		return err
	}

	// start a transaction
	q, tx, err := e.Store.Begin()
	if err != nil {
		return err
	}
	_ = q
	defer tx.Rollback()

	for _, row := range rows {
		//log.Printf("sorc %d: survey: %d: %d\n", row.SorcID, row.OrbitID, row.TechLevel)

		// sorcID and orbitID are the id of the sorc executing the survey and the orbit being surveyed.
		sorcID, orbitID := row.SorcID, row.OrbitID

		// we may need to create a report for the sorc that is executing the survey
		sorcReportID, err := q.ReadReport(e.Store.Context, sqlc.ReadReportParams{SorcID: sorcID, TurnNo: turnNo})
		if errors.Is(err, sql.ErrNoRows) {
			sorcReportID, err = q.CreateReport(e.Store.Context, sqlc.CreateReportParams{SorcID: sorcID, TurnNo: turnNo})
		}
		if err != nil {
			log.Printf("sorc %d: survey: %d: %d: %v\n", sorcID, orbitID, row.TechLevel, err)
			return err
		}
		//log.Printf("report: %d\n", reportID)

		// create a record for this survey report
		surveyReportID, err := q.CreateReportSurvey(e.Store.Context, sqlc.CreateReportSurveyParams{ReportID: sorcReportID, OrbitID: orbitID})
		if err != nil {
			log.Printf("sorc %d: survey: %d: %d: %v\n", sorcID, orbitID, row.TechLevel, err)
			return err
		}
		//log.Printf("survey: %d\n", surveyID)

		// get the survey data for this orbit and add it to the report.
		// each deposit in the list gets a separate row in the table.
		surveyRows, err := e.Store.Queries.ReadOrbitSurvey(e.Store.Context, row.OrbitID)
		if err != nil {
			return err
		}
		for _, surveyRow := range surveyRows {
			// log.Printf("survey: %d %q %d %d\n", surveyRow.DepositNo, surveyRow.DepositKind, surveyRow.DepositQty, surveyRow.YieldPct)
			parms := sqlc.CreateReportSurveyDepositParams{
				ReportID:        surveyReportID,
				DepositNo:       surveyRow.DepositNo,
				DepositKind:     surveyRow.DepositKind,
				DepositQty:      surveyRow.DepositQty,
				DepositYieldPct: surveyRow.YieldPct,
			}
			err = q.CreateReportSurveyDeposit(e.Store.Context, parms)
			if err != nil {
				log.Printf("survey: %d %q %d %d: %v\n", surveyRow.DepositNo, surveyRow.DepositKind, surveyRow.DepositQty, surveyRow.YieldPct, err)
				return err
			}
		}
	}

	return tx.Commit()
}
