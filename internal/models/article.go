// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package models

import "time"

type Article struct {
	ID            int
	Title         string
	Slug          string
	Published     bool
	DatePublished time.Time
	DateUpdated   time.Time
}
