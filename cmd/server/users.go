// empyr - a reimagining of Vern Holford's Empyrean Challenge
// Copyright (C) 2023 Michael D Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

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
