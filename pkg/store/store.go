// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package store

//go:generate sqlc generate

import (
	"context"
	"database/sql"
	_ "embed"
	"github.com/playbymail/empyr/pkg/store/sqlc"
)

type Store struct {
	path string
	db   *sql.DB
	ctx  context.Context
	q    *sqlc.Queries
}
