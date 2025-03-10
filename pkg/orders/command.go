// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package orders

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
)

type Command struct {
	Line      int
	Command   string
	Arguments []string
	Errors    []error
}

// loadCommands loads all the command records from the buffer.
func loadCommands(buffer []byte) ([]Command, error) {
	csvReader := csv.NewReader(bytes.NewReader(buffer))
	// we don't want the spaces
	csvReader.TrimLeadingSpace = true
	// special value means variable number of fields per record
	csvReader.FieldsPerRecord = -1

	raw, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reader: %w", err)
	}

	var commands []Command
	for row := range raw {
		commands = append(commands, parseCommand(row+1, raw[row]...))
	}

	return commands, nil
}

func parseCSV(rows [][]string) (commands []Command) {
	for row := range rows {
		cmd := parseCommand(row+1, rows[row]...)
		commands = append(commands, cmd)
	}
	return commands
}

// parseCommand assumes that the command has been split into tokens, one token per argument.
func parseCommand(line int, tokens ...string) Command {
	c := Command{Line: line}
	for _, token := range tokens {
		arg := strings.TrimSpace(token)
		if c.Command != "" {
			c.Arguments = append(c.Arguments, arg)
			continue
		}

		switch arg {
		case "attack spies":
			c.Command = "attack-spies"
		case "game":
			c.Command = arg
		case "auth":
			c.Command = arg
		case "assemble":
			c.Command = arg
		case "bombard":
			c.Command = arg
		case "build change":
			c.Command = arg
		case "buy":
			c.Command = arg
		case "check for spies":
			c.Command = "check-for-spies"
		case "check rebels":
			c.Command = "check-rebels"
		case "control":
			c.Command = arg
		case "convert rebels":
			c.Command = "convert-rebels"
		case "draft":
			c.Command = arg
		case "incite rebels":
			c.Command = "incite-rebels"
		case "information":
			c.Command = "information"
		case "invade":
			c.Command = arg
		case "mining":
			c.Command = arg
		case "move":
			c.Command = arg
		case "name":
			c.Command = arg
		case "News":
			c.Command = "news"
		case "pay":
			c.Command = arg
		case "permission denied":
			c.Command = "permission-denied"
		case "permission granted":
			c.Command = "permission-granted"
		case "permission to colonize":
			c.Command = "permission-to-colonize"
		case "probe":
			c.Command = arg
		case "raid":
			c.Command = arg
		case "ration":
			c.Command = arg
		case "sell":
			c.Command = arg
		case "Set Up":
			c.Command = "set-up"
		case "support":
			c.Command = arg
		case "survey":
			c.Command = arg
		case "transfer":
			c.Command = arg
		case "un-control":
			c.Command = arg
		default:
			c.Arguments = append(c.Arguments, arg)
		}
	}
	if c.Command == "" {
		c.Errors = append(c.Errors, fmt.Errorf("%d: unknown command", line))
	}
	return c
}
