// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package repos

import (
	"github.com/playbymail/empyr/repos/sqlite"
)

func (s *Store) CreateSCProbeOrder(scID, turnNo, techLevel int64, kind string, targetID int64) (int64, error) {
	parms := sqlite.CreateSCProbeOrderParams{ScID: scID, Kind: kind, TargetID: targetID}
	return s.Queries.CreateSCProbeOrder(s.Context, parms)
}

func (s *Store) CreateSCSurveyOrder(scID, targetID int64) (int64, error) {
	parms := sqlite.CreateSCSurveyOrderParams{ScID: scID, TargetID: targetID}
	return s.Queries.CreateSCSurveyOrder(s.Context, parms)
}
