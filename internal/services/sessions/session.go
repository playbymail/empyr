// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sessions

import "github.com/playbymail/empyr/internal/domains"

type SessionID int64

type Session struct {
	ID     SessionID
	UserID domains.UserID
}
