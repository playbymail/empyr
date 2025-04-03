// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package repos implements the data store via a Sqlite database.
package repos

import (
	"context"
	"database/sql"
	_ "embed"
	"github.com/playbymail/empyr/repos/sqlite"
)

type Store struct {
	Path    string
	DB      *sql.DB
	Context context.Context
	Queries *sqlite.Queries
}
