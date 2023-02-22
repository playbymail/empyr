// empyr - a game engine for Empyrean Challenge
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

package scanner

import (
	"bytes"
	"encoding/csv"
)

// Orders loads all the order records from the buffer.
func Orders(buffer []byte) ([][]string, error) {
	csvReader := csv.NewReader(bytes.NewReader(buffer))
	// we don't want the spaces
	csvReader.TrimLeadingSpace = true
	// special value means variable number of fields per record
	csvReader.FieldsPerRecord = -1

	return csvReader.ReadAll()
}
