// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/playbymail/empyr/store/sqlc"
	"html/template"
	"log"
	"math"
	"math/rand/v2"
	"strings"
	"time"
)

// commands are the commands that can be issued to the engine.
// they should be implemented elsewhere, but this is convenient for now.

var (
	//go:embed templates/cluster-map.gohtml
	clusterMapTmpl string
)

type CreateClusterMapParams_t struct {
	Code string // code of the game to create the cluster map for
}

// CreateClusterMapCommand creates a cluster map.
// It returns a byte array containing the map as HTML.
func CreateClusterMapCommand(e *Engine_t, cfg *CreateClusterMapParams_t) ([]byte, error) {
	ts, err := template.New("cluster-map").Parse(clusterMapTmpl)
	if err != nil {
		return nil, err
	}

	type system_t struct {
		Id      int64
		X, Y, Z int64
		Color   template.JS
	}

	payload := struct {
		Game    string
		Systems []system_t
	}{
		Game: cfg.Code,
	}

	rows, err := e.Store.Queries.ReadClusterMap(e.Store.Context, cfg.Code)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		var color template.JS
		switch row.NumberOfStars {
		case 1:
			color = "Blue"
		case 2:
			color = "Yellow"
		case 3:
			color = "White"
		case 4:
			color = "Red"
		default:
			return nil, fmt.Errorf("assert(s.NumberOfStars != %d)", row.NumberOfStars)
		}
		payload.Systems = append(payload.Systems, system_t{
			Id:    row.ID.Int64,
			X:     row.X.Int64 - 15, // shift the origin back to 0,0,0
			Y:     row.Y.Int64 - 15, // shift the origin back to 0,0,0
			Z:     row.Z.Int64 - 15, // shift the origin back to 0,0,0
			Color: color,
		})
	}

	// buffer will hold the cluster map
	buffer := &bytes.Buffer{}

	// execute the template, writing the result to the buffer
	if err = ts.Execute(buffer, payload); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

var (
	//go:embed templates/cluster-star-list.gohtml
	clusterStarListTmpl string
)

type CreateClusterStarListParams_t struct {
	Code string // code of the game to create the cluster star list for
}

// CreateClusterStarListCommand creates a cluster star list for a game.
// It returns a byte array containing the star list as HTML and another
// byte array containing the star list as JSON.
func CreateClusterStarListCommand(e *Engine_t, cfg *CreateClusterStarListParams_t) ([]byte, []byte, error) {
	ts, err := template.New("cluster-star-list").Parse(clusterStarListTmpl)
	if err != nil {
		return nil, nil, err
	}

	// System ID</th><th>Coordinates</th><th>Number of Stars</th><th>Distance From Center
	type system_t struct {
		Id                 int64  `json:"id,omitempty"`
		X                  int64  `json:"x,omitempty"`
		Y                  int64  `json:"y,omitempty"`
		Z                  int64  `json:"z,omitempty"`
		Coordinates        string `json:"coordinates,omitempty"`
		NumberOfStars      int64  `json:"number-of-stars,omitempty"`
		DistanceFromCenter int64  `json:"distance-from-center,omitempty"`
	}

	payload := struct {
		Game        string
		UpdatedDate string
		Systems     []system_t
	}{
		Game:        cfg.Code,
		UpdatedDate: time.Now().UTC().Format("2006-01-02"),
	}

	rows, err := e.Store.Queries.ReadClusterMap(e.Store.Context, cfg.Code)
	if err != nil {
		return nil, nil, err
	}
	for _, row := range rows {
		coordinates := fmt.Sprintf("(%02d, %02d, %02d)", row.X.Int64, row.Y.Int64, row.Z.Int64)
		dx, dy, dz := row.X.Int64-15, row.Y.Int64-15, row.Z.Int64-15
		distance := int64(math.Ceil(math.Sqrt(float64(dx*dx + dy*dy + dz*dz))))
		payload.Systems = append(payload.Systems, system_t{
			Id:                 row.ID.Int64,
			X:                  row.X.Int64,
			Y:                  row.Y.Int64,
			Z:                  row.Z.Int64,
			Coordinates:        coordinates,
			DistanceFromCenter: distance,
			NumberOfStars:      row.NumberOfStars,
		})
	}

	// buffer will hold the cluster star list
	buffer := &bytes.Buffer{}

	// execute the template, writing the result to the buffer
	if err = ts.Execute(buffer, payload); err != nil {
		return nil, nil, err
	}

	data, err := json.Marshal(payload.Systems)
	if err != nil {
		return nil, nil, err
	}

	return buffer.Bytes(), data, nil
}

type CreateEmpireParams_t struct {
	Code   string
	Handle string
}

func CreateEmpireCommand(e *Engine_t, cfg *CreateEmpireParams_t) (int64, int64, error) {
	log.Printf("create: empire: code %q\n", cfg.Code)

	if cfg.Handle == "" {
		return 0, 0, fmt.Errorf("handle: missing")
	} else if strings.TrimSpace(cfg.Handle) != cfg.Handle {
		return 0, 0, fmt.Errorf("handle: invalid")
	}

	q, tx, err := e.Store.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	gameRow, err := q.ReadGameByCode(e.Store.Context, cfg.Code)
	if err != nil {
		return 0, 0, err
	}
	parms := sqlc.CreateEmpireParams{
		GameID:       gameRow.ID,
		EmpireNo:     gameRow.LastEmpireNo + 1,
		Name:         fmt.Sprintf("Empire %03d", gameRow.LastEmpireNo+1),
		Handle:       cfg.Handle,
		HomeSystemID: gameRow.HomeSystemID,
		HomeStarID:   gameRow.HomeStarID,
		HomeOrbitID:  gameRow.HomeOrbitID,
		HomePlanetID: gameRow.HomePlanetID,
	}
	if parms.Handle == "player1" {
		parms.Handle = fmt.Sprintf("player-%03d", parms.EmpireNo)
	}
	empireId, err := q.CreateEmpire(e.Store.Context, parms)
	if err != nil {
		return 0, 0, err
	}

	// create a home colony
	sorcId, err := q.CreateSorC(e.Store.Context, sqlc.CreateSorCParams{
		EmpireID: empireId,
		Kind:     "open-colony",
	})
	if err != nil {
		return 0, 0, err
	}
	log.Printf("create: empire: id %d: no %d: colony %d\n", empireId, parms.EmpireNo, sorcId)

	detailsId, err := q.CreateSorCDetails(e.Store.Context, sqlc.CreateSorCDetailsParams{
		SorcID:      sorcId,
		TurnNo:      0,
		TechLevel:   1,
		Name:        "Not Named",
		UemQty:      59_000_000,
		UskQty:      60_000_000,
		UskPay:      0.125,
		ProQty:      15_000_000,
		ProPay:      0.3750,
		SldQty:      2_500_000,
		SldPay:      0.25,
		CnwQty:      10_000,
		SpyQty:      20,
		BirthRate:   0.0,
		DeathRate:   0.00625,
		Sol:         0.481,
		OrbitID:     parms.HomeOrbitID,
		IsOnSurface: 1,
	})
	if err != nil {
		return 0, 0, err
	}
	log.Printf("create: empire: id %d: no %d: colony %d: details %d\n", empireId, parms.EmpireNo, sorcId, detailsId)

	for _, unit := range []struct {
		kind      string
		techLevel int64
		qty       int64
	}{
		{"STUN", 0, 60_000_000},
	} {
		err = q.CreateSorCInfrastructure(e.Store.Context, sqlc.CreateSorCInfrastructureParams{
			SorcID:    sorcId,
			Kind:      unit.kind,
			TechLevel: unit.techLevel,
			Qty:       unit.qty,
		})
		if err != nil {
			return 0, 0, err
		}
	}

	for _, unit := range []struct {
		kind      string
		techLevel int64
		qty       int64
	}{
		{"FRM", 1, 130_000},
		{"MSL", 1, 50_000},
		{"SEN", 1, 20},
	} {
		err = q.CreateSorCSuperstructure(e.Store.Context, sqlc.CreateSorCSuperstructureParams{
			SorcID:    sorcId,
			Kind:      unit.kind,
			TechLevel: unit.techLevel,
			Qty:       unit.qty,
		})
		if err != nil {
			return 0, 0, err
		}
	}

	err = q.UpdateGameEmpireCounter(e.Store.Context, sqlc.UpdateGameEmpireCounterParams{GameID: gameRow.ID, EmpireNo: parms.EmpireNo})
	if err != nil {
		return 0, 0, err
	}

	return empireId, parms.EmpireNo, tx.Commit()
}

type CreateGameParams_t struct {
	Code                        string
	Name                        string
	DisplayName                 string
	NumberOfEmpires             int64
	PopulateSystemDistanceTable bool
	Rand                        *rand.Rand
}

// CreateGameCommand creates a new game.
func CreateGameCommand(e *Engine_t, cfg *CreateGameParams_t) (int64, error) {
	log.Printf("create: game: code %q: name %q: display %q\n", cfg.Code, cfg.Name, cfg.DisplayName)
	return e.CreateGame(cfg.Code, cfg.Name, cfg.DisplayName, cfg.NumberOfEmpires, cfg.PopulateSystemDistanceTable, cfg.Rand)
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
func CreateTurnReportCommand(e *Engine_t, cfg *CreateTurnReportParams_t) ([]byte, error) {
	ts, err := template.New("turn-report").Parse(turnReportTmpl)
	if err != nil {
		return nil, err
	}

	var gameId int64
	var empireId int64
	if row, err := e.Store.Queries.ReadGameEmpire(e.Store.Context, sqlc.ReadGameEmpireParams{GameCode: cfg.Code, EmpireNo: cfg.EmpireNo}); err != nil {
		log.Printf("error: %v\n", err)
		return nil, err
	} else {
		gameId = row.GameID
		empireId = row.EmpireID
	}
	log.Printf("game: %d: empire: %d\n", gameId, empireId)

	type turn_report_inventory_t struct {
		Unit      string
		Qty       int
		Assembled bool
		Storage   bool
	}

	type turn_report_system_t struct {
		Id          int
		Coordinates struct {
			X int
			Y int
			Z int
		}
	}

	type turn_report_star_t struct {
		Id          int
		System      *turn_report_system_t
		Sequence    string
		Coordinates string
	}

	type turn_report_orbit_t struct {
		Id        int
		Star      *turn_report_star_t
		OrbitNo   int
		Kind      string
		Habitable bool
	}

	type turn_report_planet_t struct {
		Id           int64
		Orbit        *turn_report_orbit_t
		Habitability int64
	}

	type turn_report_colony_t struct {
		Id              int64
		StarCoordinates string
		OrbitNo         int64
		Name            string
		Kind            string
		TL              string
		SOL             string
		Vitals          struct {
			BirthRate string
			DeathRate string
			Rations   string
			PayRates  struct {
				USK string
				PRO string
				SLD string
			}
			Census []int64
		}
		Inventory []*turn_report_inventory_t
		Factories []int64
		Mines     []int64
	}

	payload := struct {
		Site struct {
			CSS string
		}
		Game            string
		CreatedDate     string
		CreatedDateTime string
		TurnNo          int64
		TurnCode        string
		EmpireNo        int64
		EmpireCode      string
		Colonies        []*turn_report_colony_t
	}{
		Game:            cfg.Code,
		CreatedDate:     time.Now().UTC().Format("2006-01-02"),
		CreatedDateTime: time.Now().UTC().Format(time.RFC3339),
		TurnNo:          cfg.TurnNo,
		TurnCode:        fmt.Sprintf("T%05d", cfg.TurnNo),
		EmpireNo:        cfg.EmpireNo,
		EmpireCode:      fmt.Sprintf("E%03d", cfg.EmpireNo),
	}
	payload.Site.CSS = "a02/css/monospace.css"

	payload.Colonies = append(payload.Colonies, &turn_report_colony_t{})
	// buffer will hold the rendered turn report
	buffer := &bytes.Buffer{}

	// execute the template, writing the result to the buffer
	if err = ts.Execute(buffer, payload); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
