// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package actions

import "database/sql"

type TransactionsFacade struct {
	db *sql.DB
}
