// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package store implements the data store via a Sqlite database.
package store

//go:generate sqlc generate

import (
	"context"
	"database/sql"
	_ "embed"
	"github.com/playbymail/empyr/store/sqlc"
)

type Store struct {
	Path    string
	DB      *sql.DB
	Context context.Context
	Queries *sqlc.Queries
}
