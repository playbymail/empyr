// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package engine

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
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
