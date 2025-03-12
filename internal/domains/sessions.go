// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

import "time"

type SessionID int64

type Session struct {
	ID        SessionID
	User      UserID
	ExpiresAt time.Time
}
