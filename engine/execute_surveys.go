// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"log"
)

// ExecuteSurveys executes all the survey orders for the current turn.
// It updates the survey results table with the results of the survey.
// The survey includes the location, habitability, population of the
// orbit, and data on all deposits in the orbit.
func (e *Engine_t) ExecuteSurveys(gameCode string, turnNo int64) error {
	// get a list of all the survey orders. these are the orders that need to be executed.
	surveyOrderRows, err := e.Store.Queries.ReadAllSurveyOrdersForGameTurn(e.Store.Context, turnNo)
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

	for _, surveyOrder := range surveyOrderRows {
		log.Printf("sorc %d: turn %d: survey %d\n", surveyOrder.ScID, surveyOrder.TurnNo, surveyOrder.TargetID)

		// scID and orbitID are the id of the SC executing the survey and the orbit being surveyed.
		scID, turnNo, orbitID := surveyOrder.ScID, surveyOrder.TurnNo, surveyOrder.TargetID
		log.Printf("sorc %d: turn %d: survey %d\n", scID, turnNo, orbitID)

		// think ReadOrbitSurvey

		//// we may need to create a report for the SC that is executing the survey
		//scReportID, err := q.ReadReport(e.Store.Context, sqlc.ReadReportParams{ScID: scID, TurnNo: turnNo})
		//if errors.Is(err, sql.ErrNoRows) {
		//	sorcReportID, err = q.CreateReport(e.Store.Context, sqlc.CreateReportParams{ScID: scID, TurnNo: turnNo})
		//}
		//if err != nil {
		//	log.Printf("sorc %d: survey: %d: %d: %v\n", scID, orbitID, row.TechLevel, err)
		//	return err
		//}
		////log.Printf("report: %d\n", reportID)
		//
		//// create a record for this survey report
		//surveyReportID, err := q.CreateReportSurvey(e.Store.Context, sqlc.CreateReportSurveyParams{ReportID: sorcReportID, OrbitID: orbitID})
		//if err != nil {
		//	log.Printf("sorc %d: survey: %d: %d: %v\n", scID, orbitID, row.TechLevel, err)
		//	return err
		//}
		////log.Printf("survey: %d\n", surveyID)
		//
		//// get the survey data for this orbit and add it to the report.
		//// each deposit in the list gets a separate row in the table.
		//surveyRows, err := e.Store.Queries.ReadOrbitSurvey(e.Store.Context, row.OrbitID)
		//if err != nil {
		//	return err
		//}
		//for _, surveyRow := range surveyRows {
		//	// log.Printf("survey: %d %q %d %d\n", surveyRow.DepositNo, surveyRow.DepositKind, surveyRow.DepositQty, surveyRow.YieldPct)
		//	parms := sqlc.CreateReportSurveyDepositParams{
		//		ReportID:        surveyReportID,
		//		DepositNo:       surveyRow.DepositNo,
		//		DepositKind:     surveyRow.DepositKind,
		//		DepositQty:      surveyRow.DepositQty,
		//		DepositYieldPct: surveyRow.YieldPct,
		//	}
		//	err = q.CreateReportSurveyDeposit(e.Store.Context, parms)
		//	if err != nil {
		//		log.Printf("survey: %d %q %d %d: %v\n", surveyRow.DepositNo, surveyRow.DepositKind, surveyRow.DepositQty, surveyRow.YieldPct, err)
		//		return err
		//	}
		//}
	}

	return tx.Commit()
}
