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
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Lexeme struct {
	Line    int
	Kind    Kind
	Float   float64
	Integer int
	Text    string
}

func (l *Lexeme) String() string {
	switch l.Kind {
	case EOF:
		return fmt.Sprintf("{%d EOF}", l.Line)
	case EOL:
		return fmt.Sprintf("{%d EOL}", l.Line)
	case COMMA:
		return fmt.Sprintf("{%d ','}", l.Line)
	case FLOAT:
		return fmt.Sprintf("{%d %g}", l.Line, l.Float)
	case INTEGER:
		return fmt.Sprintf("{%d %d}", l.Line, l.Integer)
	case PARENCL:
		return fmt.Sprintf("{%d ')'}", l.Line)
	case PARENOP:
		return fmt.Sprintf("{%d '('}", l.Line)
	case PERCENTAGE:
		return fmt.Sprintf("{%d %d%%}", l.Line, l.Integer)
	case POPULATION:
		return fmt.Sprintf("{%d %s}", l.Line, l.Text)
	case PRODUCT:
		if l.Integer == 0 {
			return fmt.Sprintf("{%d %s}", l.Line, l.Text)
		}
		return fmt.Sprintf("{%d %s-%d}", l.Line, l.Text, l.Integer)
	case RESOURCE:
		return fmt.Sprintf("{%d %s}", l.Line, l.Text)
	case TECHLEVEL:
		return fmt.Sprintf("{%d TL-%d}", l.Line, l.Integer)
	}
	return fmt.Sprintf("{%d %q}", l.Line, l.Text)
}

// Kind is the type of lexeme.
type Kind int

// enums for Kind
const (
	EOF Kind = iota
	EOL
	COMMA
	DEPOSITID
	FACTGRP
	FLOAT
	INTEGER
	MINEGRP
	PARENCL
	PARENOP
	PERCENTAGE
	POPULATION
	PRODUCT
	QTEXT
	RESEARCH
	RESOURCE
	TECHLEVEL
	TEXT
	UUID
)

func Scan(buffer []byte) ([]*Lexeme, error) {
	// isdigit returns true if the byte is a digit
	isdigit := func(ch byte) bool {
		return '0' <= ch && ch <= '9'
	}
	// ishexdigit returns true if the byte is a hex digit
	ishexdigit := func(ch byte) bool {
		return ('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
	}
	// isspace returns true if the rune is a space or invalid rune.
	// A new-line is not considered a space.
	isspace := func(r rune) bool {
		return r != '\n' && (unicode.IsSpace(r) || r == utf8.RuneError)
	}
	// isuuid returns true if the lexeme looks like a UUID
	isuuid := func(b []byte) bool {
		if len(b) != 36 {
			return false
		}
		for i, ch := range b {
			switch i {
			case 8, 13, 18, 23:
				if ch != '-' {
					return false
				}
			default:
				if !ishexdigit(ch) {
					return false
				}
			}
		}
		return true
	}

	next := func() *Lexeme {
		var lexeme []byte
		if len(buffer) == 0 {
			return nil
		}

		for len(buffer) != 0 {
			r, w := utf8.DecodeRune(buffer)

			// is it whitespace?
			if isspace(r) {
				// skip spaces
				for len(buffer) != 0 {
					r, w = utf8.DecodeRune(buffer)
					if !isspace(r) {
						break
					}
					buffer = buffer[w:]
				}
				continue
			}

			// skip comments
			if r == ';' {
				for len(buffer) != 0 && buffer[0] != '\n' {
					r, w = utf8.DecodeRune(buffer)
					buffer = buffer[w:]
				}
				continue
			}

			lexeme, buffer = append(lexeme, buffer[:w]...), buffer[w:]

			// is it a single character lexeme (such as a new-line)?
			switch r {
			case '\n':
				return &Lexeme{Kind: EOL, Text: "\n"}
			case '(':
				return &Lexeme{Kind: PARENOP, Text: "("}
			case ')':
				return &Lexeme{Kind: PARENCL, Text: ")"}
			case ',':
				return &Lexeme{Kind: COMMA, Text: ","}
			}

			// is it quoted text?
			if r == '"' {
				r, w = utf8.DecodeRune(buffer)
				for len(buffer) != 0 && r != '\n' && r != '"' {
					if r == '\\' {
						// escaped quotes are accepted as part of the quoted string
						if len(buffer) > 0 || buffer[1] == '"' {
							// consume the escape character
							buffer = buffer[w:]
							r, w = utf8.DecodeRune(buffer)
						}
					}
					lexeme, buffer = append(lexeme, buffer[:w]...), buffer[w:]
					r, w = utf8.DecodeRune(buffer)
				}
				if r == '"' {
					// consume the closing quote
					buffer = buffer[w:]
				}
				// don't include the quote marks in the lexeme
				return &Lexeme{Kind: QTEXT, Text: string(lexeme[1:])}
			}

			// is it an integer or integer followed by a percent sign?
			if unicode.IsDigit(r) || ((r == '-' || r == '+') && (len(buffer) != 0 && isdigit(buffer[0]))) {
				kind := INTEGER
				r, w = utf8.DecodeRune(buffer)
				for unicode.IsDigit(r) {
					lexeme, buffer = append(lexeme, buffer[:w]...), buffer[w:]
					r, w = utf8.DecodeRune(buffer)
				}
				if r == '.' && len(buffer) > 0 && isdigit(buffer[1]) { // may be a floating point number
					kind = FLOAT
					lexeme, buffer = append(lexeme, buffer[:w]...), buffer[w:]
					r, w = utf8.DecodeRune(buffer)
					for len(buffer) != 0 && unicode.IsDigit(r) {
						lexeme, buffer = append(lexeme, buffer[:w]...), buffer[w:]
						r, w = utf8.DecodeRune(buffer)
					}
				} else if r == '%' { // may be a percentage number
					kind = PERCENTAGE
					lexeme, buffer = append(lexeme, buffer[:w]...), buffer[w:]
					r, w = utf8.DecodeRune(buffer)
				}
				// must be terminated by a delimiter (space, eol, eof, comma, parentheses, comment)
				if r == '\n' || isspace(r) || r == ',' || r == '(' || r == ')' || r == ';' {
					return &Lexeme{Kind: kind, Text: string(lexeme)}
				}
				// it's not an integer or percentage, so fall through
			}

			// the lexeme is everything up to the next comma, comment, new-line, or paren, or space
			for len(buffer) != 0 && !(r == ',' || r == ';' || r == '\n' || r == '(' || r == ')' || isspace(r)) {
				lexeme, buffer = append(lexeme, buffer[:w]...), buffer[w:]
				r, w = utf8.DecodeRune(buffer)
			}

			if isuuid(lexeme) {
				return &Lexeme{Kind: UUID, Text: string(bytes.ToLower(lexeme))}
			}

			return &Lexeme{Kind: TEXT, Text: string(bytes.ToLower(lexeme))}
		}

		// force an end of line at the end of the input
		return &Lexeme{Kind: EOL}
	}

	var lexemes []*Lexeme
	for no, lexeme := 1, next(); lexeme != nil; lexeme = next() {
		lexeme.Line = no
		switch lexeme.Kind {
		case EOL:
			no = no + 1
			// filter out blank lines
			if len(lexemes) == 0 || lexemes[len(lexemes)-1].Kind == EOL {
				continue
			}
		case FLOAT:
			lexeme.Float, _ = strconv.ParseFloat(lexeme.Text, 64)
		case INTEGER:
			lexeme.Integer, _ = strconv.Atoi(lexeme.Text)
		case PERCENTAGE:
			lexeme.Integer, _ = strconv.Atoi(strings.TrimRight(lexeme.Text, "%"))
		case TEXT:
			switch lexeme.Text {
			case "civilian":
				lexeme.Kind, lexeme.Text = POPULATION, "CIV"
			case "construction-crew":
				lexeme.Kind, lexeme.Text = POPULATION, "CONS"
			case "professional":
				lexeme.Kind, lexeme.Text = POPULATION, "PRO"
			case "soldier":
				lexeme.Kind, lexeme.Text = POPULATION, "SLD"
			case "spy":
				lexeme.Kind, lexeme.Text = POPULATION, "SPY"
			case "unskilled-worker":
				lexeme.Kind, lexeme.Text = POPULATION, "UNSK"
			case "unsk":
				lexeme.Kind, lexeme.Text = POPULATION, "UNSK"
			case "fuel":
				lexeme.Kind, lexeme.Text = RESOURCE, "FUEL"
			case "gold":
				lexeme.Kind, lexeme.Text = RESOURCE, "GOLD"
			case "metallics":
				lexeme.Kind, lexeme.Text = RESOURCE, "MTL"
			case "non-metallics":
				lexeme.Kind, lexeme.Text = RESOURCE, "NMTL"
			case "research":
				lexeme.Kind, lexeme.Text = RESEARCH, "RESEARCH"
			default:
				if strings.HasPrefix(lexeme.Text, "dp-") { // deposit id
					if id, err := strconv.Atoi(lexeme.Text[3:]); err == nil {
						lexeme.Kind, lexeme.Text = DEPOSITID, fmt.Sprintf("DP-%d", id)
					}
				} else if strings.HasPrefix(lexeme.Text, "fg-") { // factory group
					if id, err := strconv.Atoi(lexeme.Text[3:]); err == nil {
						lexeme.Kind, lexeme.Text = FACTGRP, fmt.Sprintf("FG-%d", id)
					}
				} else if strings.HasPrefix(lexeme.Text, "mg-") { // mining group
					if id, err := strconv.Atoi(lexeme.Text[3:]); err == nil {
						lexeme.Kind, lexeme.Text = MINEGRP, fmt.Sprintf("MG-%d", id)
					}
				} else if strings.HasPrefix(lexeme.Text, "tl-") { // tech level
					if tl, err := strconv.Atoi(lexeme.Text[3:]); err == nil && 0 < tl && tl <= 10 {
						lexeme.Kind, lexeme.Text, lexeme.Integer = TECHLEVEL, fmt.Sprintf("TL-%d", tl), tl
					}
				} else {
					// product will be xxx, xxx-yyy, or xxx-yyy-tl
					// if product includes tl, we must extract it.
					var product string
					var tl int
					if fields := strings.Split(lexeme.Text, "-"); len(fields) == 1 {
						// product is xxx
						product, tl = lexeme.Text, 0
					} else {
						// product is xxx-yyy-tl or xxx-yyy
						firstFields, lastField := fields[:len(fields)-1], fields[len(fields)-1]
						if n, err := strconv.Atoi(lastField); err == nil {
							// product is xxx-yyy-tl
							product, tl = strings.Join(firstFields, "-"), n
						} else {
							// product is xxx-yyy
							product, tl = lexeme.Text, 0
						}
					}
					switch product {
					case "amsl", "anti-missile":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "AMSL", tl
					case "ascr", "assault-craft":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "ASCR", tl
					case "aswp", "assault-weapons":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "ASWP", tl
					case "auto", "automation":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "AUTO", tl
					case "cngd", "consumer-goods":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "CNGD", tl
					case "eshd", "energy-shield":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "ESHD", tl
					case "ewpn", "energy-weapon":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "EWPN", tl
					case "fact", "factory":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "FACT", tl
					case "farm":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "FARM", tl
					case "food":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "FOOD", tl
					case "hdrv", "hyper-engine":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "HDRV", tl
					case "ls", "life-support":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "LS", tl
					case "lsu", "light-structural-unit":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "LSU", tl
					case "milr", "military-robot":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "MILR", tl
					case "mils", "military-supplies":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "MILS", tl
					case "mine":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "MINE", tl
					case "mssl", "missile":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "MSSL", tl
					case "msln", "missile-launcher":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "MSLN", tl
					case "snsr", "sensor":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "SNSR", tl
					case "sdrv", "space-drive":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "SDRV", tl
					case "su", "structural-unit":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "SU", tl
					case "slsu", "super-light-structural-unit":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "SLSU", tl
					case "trns", "transport":
						lexeme.Kind, lexeme.Text, lexeme.Integer = PRODUCT, "TRNS", tl
					}
				}
			}
		case UUID:
			lexeme.Text = strings.ToLower(lexeme.Text)
		}
		lexemes = append(lexemes, lexeme)
	}

	// force an end of file at the end of the input
	if len(lexemes) == 0 {
		lexemes = append(lexemes, &Lexeme{Kind: EOF})
	} else {
		lexemes = append(lexemes, &Lexeme{Line: lexemes[len(lexemes)-1].Line, Kind: EOF})
	}

	return lexemes, nil
}
