// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package store

import (
	"github.com/playbymail/empyr/store/sqlc"
)

func (s *Store) CreateSorCProbeOrder(sorcID, turnNo, techLevel int64, kind string, targetID int64) error {
	parms := sqlc.CreateSorCProbeOrderParams{SorcID: sorcID, TurnNo: turnNo, TechLevel: techLevel, Kind: kind, TargetID: targetID}
	return s.Queries.CreateSorCProbeOrder(s.Context, parms)
}

func (s *Store) CreateSorCSurveyOrder(sorcID, turnNo, techLevel, orbitID int64) error {
	parms := sqlc.CreateSorCSurveyOrderParams{SorcID: sorcID, TurnNo: turnNo, TechLevel: techLevel, OrbitID: orbitID}
	return s.Queries.CreateSorCSurveyOrder(s.Context, parms)
}
