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

import "errors"

// ErrBadRequest is used when creating an object and something is wrong with the request.
var ErrBadRequest = errors.New("bad request")

// ErrDuplicate is used when the object is not unique.
var ErrDuplicate = errors.New("duplicate")

// ErrDuplicateAddress is used when the e-mail address is not unique.
var ErrDuplicateAddress = errors.New("duplicate e-mail address")

var ErrDuplicatePlayer = errors.New("duplicate player")

// ErrDuplicateUserName is used when the user name is not unique.
var ErrDuplicateUserName = errors.New("duplicate user name")

// ErrNoData is used when a game could not be found.
var ErrNoData = errors.New("no data found")

// ErrUnknown is used when the error source is unknown.
var ErrUnknown = errors.New("internal server error")
