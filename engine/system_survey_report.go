// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

// SystemSurveyReport_t is the report for a sytem survey.
// It contains data on the game, the current turn, the empire,
// and all the surveys that were completed during the turn.
type SystemSurveyReport_t struct {
	Heading *ReportHeading_t

	Surveys []*SurveyReport_t // list of surveys sorted by ID

	CreatedDate     string // date the report was created
	CreatedDateTime string // date and time the report was created
}

type SurveyReport_t struct {
	ID       int64
	SorCID   int64
	Name     string // name of the system, eg "02/13/28A"
	StarID   int64  // star ID, eg 1
	OrbitID  int64
	OrbitNo  int64
	Deposits []*SurveyReportLine_t
}

type SurveyReportLine_t struct {
	DepositNo string
	Resource  string
	Quantity  string
	YieldPct  string
}
