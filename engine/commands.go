// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"math/rand/v2"
	"time"
)

// commands are the commands that can be issued to the engine.
// they should be implemented elsewhere, but this is convenient for now.

var (
	//go:embed templates/cluster-map.gohtml
	clusterMapTmpl string
)

const (
	ErrGameInProgress = Error("game in progress")
	ErrInvalidPath    = Error("invalid path")
	ErrWritingReport  = Error("error writing report")
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

	rows, err := e.Store.Queries.ReadClusterMap(e.Store.Context)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		var color template.JS
		switch row.NbrOfStars {
		case 1:
			color = "Blue"
		case 2:
			color = "Yellow"
		case 3:
			color = "White"
		case 4:
			color = "Red"
		default:
			return nil, fmt.Errorf("assert(s.NumberOfStars != %d)", row.NbrOfStars)
		}
		payload.Systems = append(payload.Systems, system_t{
			Id:    row.SystemID,
			X:     row.X - 15, // shift the origin back to 0,0,0
			Y:     row.Y - 15, // shift the origin back to 0,0,0
			Z:     row.Z - 15, // shift the origin back to 0,0,0
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

	rows, err := e.Store.Queries.ReadClusterMap(e.Store.Context)
	if err != nil {
		return nil, nil, err
	}
	for _, row := range rows {
		coordinates := fmt.Sprintf("(%02d, %02d, %02d)", row.X, row.Y, row.Z)
		dx, dy, dz := row.X-15, row.Y-15, row.Z-15
		distance := int64(math.Ceil(math.Sqrt(float64(dx*dx + dy*dy + dz*dz))))
		payload.Systems = append(payload.Systems, system_t{
			Id:                 row.SystemID,
			X:                  row.X,
			Y:                  row.Y,
			Z:                  row.Z,
			Coordinates:        coordinates,
			DistanceFromCenter: distance,
			NumberOfStars:      row.NbrOfStars,
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

type CreateGameParams_t struct {
	Code                        string
	Name                        string
	DisplayName                 string
	NumberOfEmpires             int64
	IncludeEmptyResources       bool
	PopulateSystemDistanceTable bool
	Rand                        *rand.Rand
	ForceCreate                 bool
}

// CreateGameCommand creates a new game.
func CreateGameCommand(e *Engine_t, cfg *CreateGameParams_t) error {
	log.Printf("create: game: code %q: name %q: display %q\n", cfg.Code, cfg.Name, cfg.DisplayName)

	_, err := e.CreateGame(cfg.Code, cfg.Name, cfg.DisplayName, cfg.IncludeEmptyResources, cfg.PopulateSystemDistanceTable, cfg.Rand, cfg.ForceCreate)
	return err
}

func codeTL(code string, tl int64) string {
	if tl == 0 {
		return code
	}
	return fmt.Sprintf("%s-%d", code, tl)
}

func commas(n int64) string {
	in := fmt.Sprintf("%d", n)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}
