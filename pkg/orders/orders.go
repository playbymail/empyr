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

package orders

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// Load loads all the order records from the buffer.
func LoadOrdersCSV(filename string) (commands []Command, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(fmt.Errorf("loadOrdersCSV: %w", err))
	}

	// create a CSV reader that accepts a variable number of fields per record
	// and ignores leading spaces in fields.
	csvReader := csv.NewReader(bytes.NewReader(data))
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true

	// parse the data into rows and columns
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("loadOrdersCSV: %w", err)
	}

	// return the parsed orders
	return parseCSV(rows), nil
}
