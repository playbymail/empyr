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

package orders

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

// Loader loads all the order records from the buffer.
func Loader(buffer []byte) ([]Command, []error) {
	csvReader := csv.NewReader(bytes.NewReader(buffer))
	// we don't want the spaces
	csvReader.TrimLeadingSpace = true
	// special value means variable number of fields per record
	csvReader.FieldsPerRecord = -1

	raw, err := csvReader.ReadAll()
	if err != nil {
		// parse error on line 41, column 19: bare " in non-quoted-field
		return nil, []error{fmt.Errorf("reader: %w", err)}
	}

	// good luck here
	//var orders Orders
	var commands []Command
	var errs []error

	for row := range raw {
		line := row + 1
		cmd := parseCommand(row+1, raw[row]...)
		commands = append(commands)
		if cmd.Command == "" {
			cmd.Errors = append(cmd.Errors, fmt.Errorf("%d: unknown command", line))
		} else {

		}
		commands = append(commands, cmd)
		//for col := range raw[row] {
		//	raw[row][col] = strings.TrimSpace(raw[row][col])
		//}
		//if raw[row][0] == "game" {
		//	if orders.Game != nil {
		//		errs = append(errs, fmt.Errorf("%d: game: already defined", line))
		//		continue
		//	}
		//	// expect game , game-name , game-turn
		//	if len(raw[row]) != 3 {
		//		errs = append(errs, fmt.Errorf("%d: game: requires 3 values", line))
		//		continue
		//	}
		//	turn, err := strconv.Atoi(raw[row][2])
		//	if err != nil {
		//		errs = append(errs, fmt.Errorf("%d: game: turn: %w", line, err))
		//	}
		//	orders.Game = &Game{Name: raw[row][1], Turn: turn}
		//	continue
		//} else if raw[row][0] == "auth" {
		//	if orders.Auth != nil {
		//		errs = append(errs, fmt.Errorf("%d: auth: already defined", line))
		//		continue
		//	}
		//	// expect auth , auth-kind , value
		//	if len(raw[row]) != 3 {
		//		errs = append(errs, fmt.Errorf("%d: auth: requires 3 values", line))
		//		continue
		//	}
		//	if raw[row][1] == "token" {
		//		orders.Auth = &Auth{Kind: "token"}
		//	} else {
		//		errs = append(errs, fmt.Errorf("%d: auth: unknown kind", line))
		//		continue
		//	}
		//	orders.Auth.Value = raw[row][1]
		//	continue
		//} else if raw[row][1] == "assemble" {
		//} else if raw[row][1] == "bombard" {
		//} else if raw[row][1] == "build change" {
		//} else if raw[row][1] == "buy" {
		//} else if raw[row][1] == "control" {
		//} else if raw[row][1] == "draft" {
		//} else if raw[row][1] == "invade" {
		//} else if raw[row][1] == "mining" {
		//} else if raw[row][1] == "move" {
		//} else if raw[row][1] == "pay" {
		//} else if raw[row][1] == "permission to colonize" {
		//} else if raw[row][1] == "probe" {
		//} else if raw[row][1] == "raid" {
		//} else if raw[row][1] == "ration" {
		//} else if raw[row][1] == "sell" {
		//} else if raw[row][1] == "support" {
		//} else if raw[row][1] == "survey" {
		//} else if raw[row][1] == "transfer" {
		//} else {
		//	errs = append(errs, fmt.Errorf("unknown order on line %d", line))
		//	continue
		//}
	}

	return commands, errs
}
