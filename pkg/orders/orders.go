// Copyright (c) 2025 Michael D Henderson. All rights reserved.

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
