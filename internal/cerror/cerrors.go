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

// Package cerror implements constant error wrappers from Cheney.
//
// Usage is something like:
//
//	const ErrFoo = cerror.Error("foo")
//	if errors.Is(err, ErrFoo) { ... }
package cerror

// The original article by Cheney is at:
//   https://dave.cheney.net/2016/04/07/constant-errors
// Myren has a good article:
//   https://smyrman.medium.com/writing-constant-errors-with-go-1-13-10c4191617

// Error allows us to create constant error values
type Error string

// Error implements the Error interface.
func (err Error) Error() string {
	return string(err)
}
