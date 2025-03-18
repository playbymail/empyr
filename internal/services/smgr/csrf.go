// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package smgr

import (
	"github.com/google/uuid"
)

func generateCSRFToken() string {
	return uuid.New().String()
}
