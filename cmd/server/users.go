// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/playbymail/empyr/cmd/server/pkg/users"
)

// createUser
func (s *server) createUser(id, name, email string) (*users.User, error) {
	for _, user := range s.users.All() {
		if email == user.Email {
			return user, ErrDuplicateAddress
		} else if name == user.Name {
			return user, ErrDuplicateUserName
		}
	}
	return s.users.Create(id, name, email)
}

// filterUser
func (s *server) filterUsers(fn func(*users.User) bool) []*users.User {
	return s.users.Filter(fn)
}
