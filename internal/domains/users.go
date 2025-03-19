// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

type UserID int64

type User struct {
	ID       UserID
	Username string
	Email    string
	IsAdmin  bool
	IsUser   bool
	Roles    map[string]bool
}
