// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package orders

import (
	"fmt"
	"strconv"
	"strings"
)

func Parse(lexemes []*Lexeme) []any {
	var orders []any

	var cmd *Lexeme
	var order any
	for len(lexemes) != 0 && lexemes[0].Kind != EOF {
		cmd, lexemes = lexemes[0], lexemes[1:]
		switch cmd.Text {
		case "abandon":
			order, lexemes = parseAbandon(cmd, lexemes)
		case "assemble":
			order, lexemes = parseAssemble(cmd, lexemes)
		case "bombard":
			order, lexemes = parseBombard(cmd, lexemes)
		case "buy":
			order, lexemes = parseBuy(cmd, lexemes)
		case "check-rebels":
			order, lexemes = parseCheckRebels(cmd, lexemes)
		case "claim":
			order, lexemes = parseClaim(cmd, lexemes)
		case "convert-rebels":
			order, lexemes = parseConvertRebels(cmd, lexemes)
		case "counter-agents":
			order, lexemes = parseCounterAgents(cmd, lexemes)
		case "discharge":
			order, lexemes = parseDischarge(cmd, lexemes)
		case "draft":
			order, lexemes = parseDraft(cmd, lexemes)
		case "expand":
			order, lexemes = parseExpand(cmd, lexemes)
		case "grant":
			order, lexemes = parseGrant(cmd, lexemes)
		case "incite-rebels":
			order, lexemes = parseInciteRebels(cmd, lexemes)
		case "invade":
			order, lexemes = parseInvade(cmd, lexemes)
		case "jump":
			order, lexemes = parseJump(cmd, lexemes)
		case "move":
			order, lexemes = parseMove(cmd, lexemes)
		case "name":
			order, lexemes = parseName(cmd, lexemes)
		case "news":
			order, lexemes = parseNews(cmd, lexemes)
		case "pay":
			order, lexemes = parsePay(cmd, lexemes)
		case "probe":
			order, lexemes = parseProbe(cmd, lexemes)
		case "raid":
			order, lexemes = parseRaid(cmd, lexemes)
		case "ration":
			order, lexemes = parseRation(cmd, lexemes)
		case "recycle":
			order, lexemes = parseRecycle(cmd, lexemes)
		case "retool":
			order, lexemes = parseRetoolFactoryGroup(cmd, lexemes)
		case "revoke":
			order, lexemes = parseRevoke(cmd, lexemes)
		case "scrap":
			order, lexemes = parseScrap(cmd, lexemes)
		case "secret":
			order, lexemes = parseSecret(cmd, lexemes)
		case "sell":
			order, lexemes = parseSell(cmd, lexemes)
		case "setup":
			order, lexemes = parseSetup(cmd, lexemes)
		case "steal-secrets":
			order, lexemes = parseStealSecrets(cmd, lexemes)
		case "store":
			order, lexemes = parseStore(cmd, lexemes)
		case "support":
			order, lexemes = parseSupport(cmd, lexemes)
		case "suppress-agents":
			order, lexemes = parseSuppressAgents(cmd, lexemes)
		case "survey":
			order, lexemes = parseSurvey(cmd, lexemes)
		case "transfer":
			order, lexemes = parseTransfer(cmd, lexemes)
		default:
			order, lexemes = parseUnknown(cmd, lexemes)
		}
		orders = append(orders, order)
	}
	return orders
}

func eatLine(l []*Lexeme) []*Lexeme {
	for len(l) != 0 && l[0].Kind != EOF {
		if l[0].Kind == EOL {
			return l[1:]
		}
		l = l[1:]
	}
	return l
}

// expectCargo wants population or product or research or resource
func expectCargo(l []*Lexeme) (string, int, []*Lexeme, error) {
	if len(l) == 0 {
		return "", 0, l, fmt.Errorf("want material, got eof")
	}
	if population, rest, err := expectPopulation(l); err == nil {
		return population, 0, rest, nil
	}
	if product, tl, rest, err := expectProduct(l); err == nil {
		return product, tl, rest, nil
	}
	if research, rest, err := expectResearch(l); err == nil {
		return research, 0, rest, nil
	}
	if resource, rest, err := expectResource(l); err == nil {
		return resource, 0, rest, nil
	}
	return "", 0, l, fmt.Errorf("want material, got %q", l[0].Text)
}

// coordinates are (x, y, z(suffix?)(, orbit)?)
func expectCoordinates(l []*Lexeme) (Coordinates, []*Lexeme, error) {
	var err error
	if len(l) < 7 {
		return Coordinates{}, l, fmt.Errorf("want coordinates, got eof")
	}
	var c Coordinates
	if l[0].Kind != PARENOP {
		return Coordinates{}, l, fmt.Errorf("want coordinates, got %q", l[0].Text)
	}
	if l[1].Kind == INTEGER {
		c.X = l[1].Integer
	} else {
		return Coordinates{}, l, fmt.Errorf("want coordinates, got %q", l[1].Text)
	}
	if l[2].Kind != COMMA {
		return Coordinates{}, l, fmt.Errorf("want coordinates, got %q", l[2].Text)
	}
	if l[3].Kind == INTEGER {
		c.Y = l[3].Integer
	} else {
		return Coordinates{}, l, fmt.Errorf("want coordinates, got %q", l[3].Text)
	}
	if l[4].Kind != COMMA {
		return Coordinates{}, l, fmt.Errorf("want coordinates, got %q", l[4].Text)
	}
	if l[5].Kind == INTEGER {
		c.Z = l[5].Integer
	} else if l[5].Kind == TEXT { // may have system suffix
		prefix, suffix := l[5].Text[:len(l[5].Text)-1], strings.ToLower(l[5].Text[len(l[5].Text)-1:])
		if c.Z, err = strconv.Atoi(prefix); err != nil {
			return Coordinates{}, l, fmt.Errorf("want coordinates, got %q", l[3].Text)
		}
		if !("a" <= suffix && suffix <= "z") {
			return Coordinates{}, l, fmt.Errorf("want coordinates, got %q", l[3].Text)
		}
		c.System = suffix
	} else {
		return Coordinates{}, l, fmt.Errorf("want coordinates, got %q", l[3].Text)
	}
	if l[6].Kind == PARENCL {
		return c, l[7:], nil
	}
	if len(l) < 9 {
		return Coordinates{}, l, fmt.Errorf("want coordinates, got eof")
	}
	if l[6].Kind != COMMA {
		return Coordinates{}, l, fmt.Errorf("want coordinates, got %q", l[6].Text)
	} else if l[7].Kind != INTEGER {
		return Coordinates{}, l, fmt.Errorf("want coordinates, got %q", l[7].Text)
	} else if l[8].Kind != PARENCL {
		return Coordinates{}, l, fmt.Errorf("want coordinates, got %q", l[8].Text)
	}
	c.Orbit = l[7].Integer
	if !(0 <= c.Orbit && c.Orbit <= 10) {
		return Coordinates{}, l, fmt.Errorf("want coordinates, got %q", l[7].Text)
	}
	return c, l[9:], nil
}

func expectDepositId(l []*Lexeme) (string, []*Lexeme, error) {
	if len(l) == 0 {
		return "", l, fmt.Errorf("want depositId, got eof")
	}
	switch l[0].Kind {
	case DEPOSITID:
		return l[0].Text, l[1:], nil
	}
	return "", l, fmt.Errorf("want depositId, got %q", l[0].Text)
}

func expectEOL(l []*Lexeme) ([]*Lexeme, error) {
	if len(l) == 0 {
		return l, nil
	}
	switch l[0].Kind {
	case EOF:
		return l, nil
	case EOL:
		return l[1:], nil
	}
	return l, fmt.Errorf("want EOL, got %q", l[0].Text)
}

func expectFactoryGroup(l []*Lexeme) (string, []*Lexeme, error) {
	if len(l) == 0 {
		return "", l, fmt.Errorf("want factoryGroup, got eof")
	}
	switch l[0].Kind {
	case FACTGRP:
		return l[0].Text, l[1:], nil
	}
	return "", l, fmt.Errorf("want factoryGroup, got %q", l[0].Text)
}

func expectInteger(l []*Lexeme) (int, []*Lexeme, error) {
	if len(l) == 0 {
		return 0, l, fmt.Errorf("want integer, got eof")
	}
	switch l[0].Kind {
	case INTEGER:
		return l[0].Integer, l[1:], nil
	}
	return 0, l, fmt.Errorf("want integer, got %q", l[0].Text)
}

// expectMaterial wants population or product or research
func expectMaterial(l []*Lexeme) (string, int, []*Lexeme, error) {
	if len(l) == 0 {
		return "", 0, l, fmt.Errorf("want material, got eof")
	}
	if population, rest, err := expectPopulation(l); err == nil {
		return population, 0, rest, nil
	}
	if product, tl, rest, err := expectProduct(l); err == nil {
		return product, tl, rest, nil
	}
	if research, rest, err := expectResearch(l); err == nil {
		return research, 0, rest, nil
	}
	return "", 0, l, fmt.Errorf("want material, got %q", l[0].Text)
}

func expectMineGroup(l []*Lexeme) (string, []*Lexeme, error) {
	if len(l) == 0 {
		return "", l, fmt.Errorf("want mineGroup, got eof")
	}
	switch l[0].Kind {
	case MINEGRP:
		return l[0].Text, l[1:], nil
	}
	return "", l, fmt.Errorf("want mineGroup, got %q", l[0].Text)
}

func expectNumber(l []*Lexeme) (float64, []*Lexeme, error) {
	if len(l) == 0 {
		return 0, l, fmt.Errorf("want number, got eof")
	}
	switch l[0].Kind {
	case FLOAT:
		return l[0].Float, l[1:], nil
	case INTEGER:
		return float64(l[0].Integer), l[1:], nil
	}
	return 0, l, fmt.Errorf("want number, got %q", l[0].Text)
}

func expectPercentage(l []*Lexeme) (int, []*Lexeme, error) {
	if len(l) == 0 {
		return 0, l, fmt.Errorf("want percentage, got eof")
	}
	switch l[0].Kind {
	case PERCENTAGE:
		return l[0].Integer, l[1:], nil
	}
	return 0, l, fmt.Errorf("want percentage, got %q", l[0].Text)
}

func expectPopulation(l []*Lexeme) (string, []*Lexeme, error) {
	if len(l) == 0 {
		return "", l, fmt.Errorf("want population, got eof")
	}
	switch l[0].Kind {
	case POPULATION:
		return l[0].Text, l[1:], nil
	}
	return "", l, fmt.Errorf("want population, got %q", l[0].Text)
}

func expectProduct(l []*Lexeme) (string, int, []*Lexeme, error) {
	if len(l) == 0 {
		return "", 0, l, fmt.Errorf("want product, got eof")
	}
	switch l[0].Kind {
	case PRODUCT:
		return l[0].Text, l[0].Integer, l[1:], nil
	}
	return "", 0, l, fmt.Errorf("want product, got %q", l[0].Text)
}

func expectQuotedText(l []*Lexeme) (string, []*Lexeme, error) {
	if len(l) == 0 {
		return "", l, fmt.Errorf("want quotedText, got eof")
	}
	switch l[0].Kind {
	case QTEXT:
		return l[0].Text, l[1:], nil
	}
	return "", l, fmt.Errorf("want quotedText, got %q", l[0].Text)
}

func expectResearch(l []*Lexeme) (string, []*Lexeme, error) {
	if len(l) == 0 {
		return "", l, fmt.Errorf("want research, got eof")
	}
	switch l[0].Kind {
	case RESEARCH:
		return l[0].Text, l[1:], nil
	}
	return "", l, fmt.Errorf("want research, got %q", l[0].Text)
}

func expectResource(l []*Lexeme) (string, []*Lexeme, error) {
	if len(l) == 0 {
		return "", l, fmt.Errorf("want resource, got eof")
	}
	switch l[0].Kind {
	case RESOURCE:
		return l[0].Text, l[1:], nil
	}
	return "", l, fmt.Errorf("want resource, got %q", l[0].Text)
}

func expectText(l []*Lexeme) (string, []*Lexeme, error) {
	if len(l) == 0 {
		return "", l, fmt.Errorf("want text, got eof")
	}
	switch l[0].Kind {
	case TEXT:
		return l[0].Text, l[1:], nil
	}
	return "", l, fmt.Errorf("want text, got %q", l[0].Text)
}

func expectUuid(l []*Lexeme) (string, []*Lexeme, error) {
	if len(l) == 0 {
		return "", l, fmt.Errorf("want uuid, got eof")
	}
	switch l[0].Kind {
	case UUID:
		return l[0].Text, l[1:], nil
	}
	return "", l, fmt.Errorf("want uuid, got %q", l[0].Text)
}

func expectUnit(l []*Lexeme) (Unit, []*Lexeme, error) {
	if len(l) == 0 {
		return Unit{}, l, fmt.Errorf("want unit, got eof")
	}
	switch l[0].Kind {
	case POPULATION:
		return Unit{Name: l[0].Text}, l[1:], nil
	case PRODUCT:
		return Unit{Name: l[0].Text, TechLevel: l[0].Integer}, l[1:], nil
	case RESEARCH:
		return Unit{Name: l[0].Text}, l[1:], nil
	case RESOURCE:
		return Unit{Name: l[0].Text}, l[1:], nil
	case TECHLEVEL:
		return Unit{Name: l[0].Text, TechLevel: l[0].Integer}, l[1:], nil
	}
	return Unit{}, l, fmt.Errorf("want unit, got %q", l[0].Text)
}

func expectWord(l []*Lexeme, words ...string) (string, []*Lexeme, error) {
	if len(l) == 0 {
		return "", l, fmt.Errorf("want keyword, got eof")
	}
	switch l[0].Kind {
	case TEXT:
		upper := strings.ToUpper(l[0].Text)
		for _, word := range words {
			if word == upper {
				return word, l[1:], nil
			}
		}
	}
	return "", l, fmt.Errorf("want keyword, got %q", l[0].Text)
}

func parseAbandon(cmd *Lexeme, l []*Lexeme) (*Abandon, []*Lexeme) {
	var err error
	o := &Abandon{Line: cmd.Line}
	if o.Location, l, err = expectCoordinates(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("location: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseAssemble(cmd *Lexeme, l []*Lexeme) (any, []*Lexeme) {
	if order, rest := parseAssembleFactoryGroup(cmd, l); order.Errors == nil {
		return order, rest
	}
	if order, rest := parseAssembleMineGroup(cmd, l); order.Errors == nil {
		return order, rest
	}
	if order, rest := parseAssembleUnit(cmd, l); order.Errors == nil {
		return order, rest
	}
	return parseUnknown(cmd, l)
}

func parseAssembleFactoryGroup(cmd *Lexeme, l []*Lexeme) (*AssembleFactoryGroup, []*Lexeme) {
	var err error
	o := &AssembleFactoryGroup{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if o.Manufacture, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("manufacture: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseAssembleMineGroup(cmd *Lexeme, l []*Lexeme) (*AssembleMineGroup, []*Lexeme) {
	var err error
	o := &AssembleMineGroup{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.DepositId, l, err = expectDepositId(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("depositId: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseAssembleUnit(cmd *Lexeme, l []*Lexeme) (*AssembleUnit, []*Lexeme) {
	var err error
	o := &AssembleUnit{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseBombard(cmd *Lexeme, l []*Lexeme) (*Bombard, []*Lexeme) {
	var err error
	o := &Bombard{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.PctCommitted, l, err = expectPercentage(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("pctCommitted: %w", err))
		return o, eatLine(l)
	}
	if o.TargetId, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("targetId: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseBuy(cmd *Lexeme, l []*Lexeme) (*Buy, []*Lexeme) {
	var err error
	o := &Buy{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		// quantity is not required when buying research
		if unit, _, nerr := expectUnit(l); nerr != nil || !strings.HasPrefix(unit.Name, "TL-") {
			o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
			return o, eatLine(l)
		}
		o.Quantity = 1
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if o.Bid, l, err = expectNumber(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("bid: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseCheckRebels(cmd *Lexeme, l []*Lexeme) (*CheckRebels, []*Lexeme) {
	var err error
	o := &CheckRebels{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseClaim(cmd *Lexeme, l []*Lexeme) (*Claim, []*Lexeme) {
	var err error
	o := &Claim{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Location, l, err = expectCoordinates(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("location: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseConvertRebels(cmd *Lexeme, l []*Lexeme) (*ConvertRebels, []*Lexeme) {
	var err error
	o := &ConvertRebels{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseCounterAgents(cmd *Lexeme, l []*Lexeme) (*CounterAgents, []*Lexeme) {
	var err error
	o := &CounterAgents{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseDraft(cmd *Lexeme, l []*Lexeme) (*Draft, []*Lexeme) {
	var err error
	o := &Draft{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Profession, l, err = expectPopulation(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("profession: %w", err))
		return o, eatLine(l)
	} else if !(o.Profession == "CONS" || o.Profession == "PRO" || o.Profession == "SLD" || o.Profession == "SPY") {
		o.Errors = append(o.Errors, fmt.Errorf("profession: invalid profession %q", o.Profession))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseDischarge(cmd *Lexeme, l []*Lexeme) (*Discharge, []*Lexeme) {
	var err error
	o := &Discharge{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Profession, l, err = expectPopulation(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("profession: %w", err))
		return o, eatLine(l)
	} else if !(o.Profession == "CONS" || o.Profession == "PRO" || o.Profession == "SLD" || o.Profession == "SPY") {
		o.Errors = append(o.Errors, fmt.Errorf("profession: invalid profession %q", o.Profession))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseExpand(cmd *Lexeme, l []*Lexeme) (any, []*Lexeme) {
	if order, rest := parseExpandFactoryGroup(cmd, l); order.Errors == nil {
		return order, rest
	}
	if order, rest := parseExpandMineGroup(cmd, l); order.Errors == nil {
		return order, rest
	}
	return parseUnknown(cmd, l)
}

func parseExpandFactoryGroup(cmd *Lexeme, l []*Lexeme) (*ExpandFactoryGroup, []*Lexeme) {
	var err error
	o := &ExpandFactoryGroup{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.FactoryGroup, l, err = expectFactoryGroup(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("factoryGroup: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseExpandMineGroup(cmd *Lexeme, l []*Lexeme) (*ExpandMineGroup, []*Lexeme) {
	var err error
	o := &ExpandMineGroup{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.MineGroup, l, err = expectMineGroup(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("mineGroup: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseGrant(cmd *Lexeme, l []*Lexeme) (*Grant, []*Lexeme) {
	var err error
	o := &Grant{Line: cmd.Line}
	if o.Location, l, err = expectCoordinates(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("location: %w", err))
		return o, eatLine(l)
	}
	if o.Kind, l, err = expectWord(l, "COLONIZE", "TRADE"); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("kind: %w", err))
		return o, eatLine(l)
	}
	if o.TargetId, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("targetId: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseInciteRebels(cmd *Lexeme, l []*Lexeme) (*InciteRebels, []*Lexeme) {
	var err error
	o := &InciteRebels{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.TargetId, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("targetId: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseInvade(cmd *Lexeme, l []*Lexeme) (*Invade, []*Lexeme) {
	var err error
	o := &Invade{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.PctCommitted, l, err = expectPercentage(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("pctCommitted: %w", err))
		return o, eatLine(l)
	}
	if o.TargetId, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("targetId: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseJump(cmd *Lexeme, l []*Lexeme) (*Jump, []*Lexeme) {
	var err error
	o := &Jump{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Location, l, err = expectCoordinates(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("location: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseMove(cmd *Lexeme, l []*Lexeme) (*Move, []*Lexeme) {
	var err error
	o := &Move{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Orbit, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("orbit: %w", err))
		return o, eatLine(l)
	} else if !(0 < o.Orbit && o.Orbit <= 10) {
		o.Errors = append(o.Errors, fmt.Errorf("orbit: invalid orbit %d", o.Orbit))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseName(cmd *Lexeme, l []*Lexeme) (any, []*Lexeme) {
	var err error
	o := &NameUnit{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err == nil {
		if o.Name, l, err = expectQuotedText(l); err != nil {
			o.Errors = append(o.Errors, fmt.Errorf("name: %w", err))
			return o, eatLine(l)
		}
		if l, err = expectEOL(l); err != nil {
			o.Errors = append(o.Errors, err)
			return o, eatLine(l)
		}
		return o, l
	}
	var location Coordinates
	if location, l, err = expectCoordinates(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("location: %w", err))
		return o, eatLine(l)
	}
	var name string
	if name, l, err = expectQuotedText(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("name: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	if 0 <= location.Orbit && location.Orbit <= 10 {
		return &Name{Line: o.Line, Location: location, Name: name}, l
	}
	o.Errors = append(o.Errors, fmt.Errorf("location: invalid orbit %d", location.Orbit))
	return o, l
}

func parseNews(cmd *Lexeme, l []*Lexeme) (*News, []*Lexeme) {
	var err error
	o := &News{Line: cmd.Line}
	if o.Location, l, err = expectCoordinates(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("location: %w", err))
		return o, eatLine(l)
	}
	if o.Article, l, err = expectQuotedText(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("article: %w", err))
		return o, eatLine(l)
	}
	if o.Signature, l, err = expectQuotedText(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("signature: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parsePay(cmd *Lexeme, l []*Lexeme) (any, []*Lexeme) {
	var err error
	pl := &PayLocal{Line: cmd.Line}
	if pl.Id, l, err = expectInteger(l); err == nil {
		if pl.Profession, l, err = expectPopulation(l); err != nil {
			pl.Errors = append(pl.Errors, fmt.Errorf("location: %w", err))
			return pl, eatLine(l)
		}
		if pl.Rate, l, err = expectNumber(l); err != nil {
			pl.Errors = append(pl.Errors, fmt.Errorf("rate: %w", err))
			return pl, eatLine(l)
		}
		if l, err = expectEOL(l); err != nil {
			pl.Errors = append(pl.Errors, err)
			return pl, eatLine(l)
		}
		return pl, l
	}
	pa := &PayAll{Line: cmd.Line}
	if pa.Profession, l, err = expectPopulation(l); err != nil {
		pa.Errors = append(pl.Errors, fmt.Errorf("location: %w", err))
		return pa, eatLine(l)
	}
	if pa.Rate, l, err = expectNumber(l); err != nil {
		pa.Errors = append(pl.Errors, fmt.Errorf("rate: %w", err))
		return pa, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		pa.Errors = append(pl.Errors, err)
		return pa, eatLine(l)
	}
	return pa, l
}

func parseProbe(cmd *Lexeme, l []*Lexeme) (any, []*Lexeme) {
	var err error
	o := &Probe{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err == nil {
		return o, l
	}
	if o.Orbit, l, err = expectInteger(l); err == nil {
		if !(0 < o.Orbit && o.Orbit <= 10) {
			o.Errors = append(o.Errors, fmt.Errorf("orbit: invalid orbit %d", o.Orbit))
			return o, eatLine(l)
		}
		if l, err = expectEOL(l); err != nil {
			o.Errors = append(o.Errors, err)
			return o, eatLine(l)
		}
		return o, l
	}
	ps := &ProbeSystem{Line: o.Line, Id: o.Id}
	if ps.Location, l, err = expectCoordinates(l); err != nil {
		ps.Errors = append(ps.Errors, fmt.Errorf("location: %w", err))
		return ps, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		ps.Errors = append(ps.Errors, err)
		return ps, eatLine(l)
	}
	return ps, l
}

func parseRaid(cmd *Lexeme, l []*Lexeme) (*Raid, []*Lexeme) {
	var err error
	o := &Raid{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.PctCommitted, l, err = expectPercentage(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("pctCommitted: %w", err))
		return o, eatLine(l)
	}
	if o.TargetId, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("targetId: %w", err))
		return o, eatLine(l)
	}
	if o.TargetUnit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("material: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseRation(cmd *Lexeme, l []*Lexeme) (any, []*Lexeme) {
	var err error
	rl := &RationLocal{Line: cmd.Line}
	if rl.Id, l, err = expectInteger(l); err == nil {
		if rl.Rate, l, err = expectPercentage(l); err != nil {
			rl.Errors = append(rl.Errors, fmt.Errorf("rate: %w", err))
			return rl, eatLine(l)
		}
		if l, err = expectEOL(l); err != nil {
			rl.Errors = append(rl.Errors, err)
			return rl, eatLine(l)
		}
		return rl, l
	}
	ra := &RationAll{Line: cmd.Line}
	if ra.Rate, l, err = expectPercentage(l); err != nil {
		ra.Errors = append(rl.Errors, fmt.Errorf("rate: %w", err))
		return ra, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		ra.Errors = append(rl.Errors, err)
		return ra, eatLine(l)
	}
	return ra, l
}

func parseRecycle(cmd *Lexeme, l []*Lexeme) (any, []*Lexeme) {
	if order, rest := parseRecycleFactoryGroup(cmd, l); order.Errors == nil {
		return order, rest
	}
	if order, rest := parseRecycleMineGroup(cmd, l); order.Errors == nil {
		return order, rest
	}
	if order, rest := parseRecycleUnit(cmd, l); order.Errors == nil {
		return order, rest
	}
	return parseUnknown(cmd, l)
}

func parseRecycleFactoryGroup(cmd *Lexeme, l []*Lexeme) (*RecycleFactoryGroup, []*Lexeme) {
	var err error
	o := &RecycleFactoryGroup{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.FactoryGroup, l, err = expectFactoryGroup(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("factoryGroup: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseRecycleMineGroup(cmd *Lexeme, l []*Lexeme) (*RecycleMineGroup, []*Lexeme) {
	var err error
	o := &RecycleMineGroup{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.MineGroup, l, err = expectMineGroup(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("mineGroup: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseRecycleUnit(cmd *Lexeme, l []*Lexeme) (*RecycleUnit, []*Lexeme) {
	var err error
	o := &RecycleUnit{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseRetoolFactoryGroup(cmd *Lexeme, l []*Lexeme) (*RetoolFactoryGroup, []*Lexeme) {
	var err error
	o := &RetoolFactoryGroup{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.FactoryGroup, l, err = expectFactoryGroup(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("factoryGroup: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseRevoke(cmd *Lexeme, l []*Lexeme) (*Revoke, []*Lexeme) {
	var err error
	o := &Revoke{Line: cmd.Line}
	if o.Location, l, err = expectCoordinates(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("location: %w", err))
		return o, eatLine(l)
	}
	if o.Kind, l, err = expectWord(l, "COLONIZE", "TRADE"); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("kind: %w", err))
		return o, eatLine(l)
	}
	if o.TargetId, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("targetId: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseScrap(cmd *Lexeme, l []*Lexeme) (any, []*Lexeme) {
	if order, rest := parseScrapFactoryGroup(cmd, l); order.Errors == nil {
		return order, rest
	}
	if order, rest := parseScrapMineGroup(cmd, l); order.Errors == nil {
		return order, rest
	}
	if order, rest := parseScrapUnit(cmd, l); order.Errors == nil {
		return order, rest
	}
	return parseUnknown(cmd, l)
}

func parseScrapFactoryGroup(cmd *Lexeme, l []*Lexeme) (*ScrapFactoryGroup, []*Lexeme) {
	var err error
	o := &ScrapFactoryGroup{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.FactoryGroup, l, err = expectFactoryGroup(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("factoryGroup: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseScrapMineGroup(cmd *Lexeme, l []*Lexeme) (*ScrapMineGroup, []*Lexeme) {
	var err error
	o := &ScrapMineGroup{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.MineGroup, l, err = expectMineGroup(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("mineGroup: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseScrapUnit(cmd *Lexeme, l []*Lexeme) (*ScrapUnit, []*Lexeme) {
	var err error
	o := &ScrapUnit{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseSecret(cmd *Lexeme, l []*Lexeme) (*Secret, []*Lexeme) {
	var err error
	o := &Secret{Line: cmd.Line}
	if o.Handle, l, err = expectText(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("handle: %w", err))
		return o, eatLine(l)
	}
	if o.Game, l, err = expectText(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("game: %w", err))
		return o, eatLine(l)
	} else if o.Game = strings.ToUpper(o.Game); !strings.HasPrefix(o.Game, "G") {
		o.Errors = append(o.Errors, fmt.Errorf("game: invalid game %q", o.Game))
		return o, eatLine(l)
	}
	if o.Turn, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("turn: %w", err))
		return o, eatLine(l)
	}
	if o.Token, l, err = expectUuid(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("uuid: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseSell(cmd *Lexeme, l []*Lexeme) (*Sell, []*Lexeme) {
	var err error
	o := &Sell{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		// quantity is not required when selling research
		if unit, _, nerr := expectUnit(l); nerr != nil || !strings.HasPrefix(unit.Name, "TL-") {
			o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
			return o, eatLine(l)
		}
		o.Quantity = 1
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if o.Ask, l, err = expectNumber(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("ask: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseSetup(cmd *Lexeme, l []*Lexeme) (*Setup, []*Lexeme) {
	var err error
	o := &Setup{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Location, l, err = expectCoordinates(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("coordinates: %w", err))
		return o, eatLine(l)
	}
	if o.Kind, l, err = expectWord(l, "COLONY", "SHIP"); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("kind: %w", err))
		return o, eatLine(l)
	}
	if o.Action, l, err = expectWord(l, "TRANSFER"); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("action: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	for {
		_, rest, err := expectWord(l, "END")
		if err == nil {
			break
		}
		l = rest
		qty, rest, err := expectInteger(l)
		if err != nil {
			o.Errors = append(o.Errors, fmt.Errorf("transfer: quantity: %w", err))
			l = eatLine(l)
			continue
		}
		l = rest
		unit, rest, err := expectUnit(l)
		if err != nil {
			o.Errors = append(o.Errors, fmt.Errorf("transfer: unit: %w", err))
			l = eatLine(l)
			continue
		}
		l = rest
		if l, err = expectEOL(l); err != nil {
			o.Errors = append(o.Errors, fmt.Errorf("transfer: eol: %w", err))
			l = eatLine(l)
			continue
		}
		o.Items = append(o.Items, &TransferDetail{Quantity: qty, Unit: unit})
	}
	if _, l, err = expectWord(l, "END"); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("end: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseStealSecrets(cmd *Lexeme, l []*Lexeme) (*StealSecrets, []*Lexeme) {
	var err error
	o := &StealSecrets{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.TargetId, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("targetId: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseStore(cmd *Lexeme, l []*Lexeme) (any, []*Lexeme) {
	if order, rest := parseStoreFactoryGroup(cmd, l); order.Errors == nil {
		return order, rest
	}
	if order, rest := parseStoreMineGroup(cmd, l); order.Errors == nil {
		return order, rest
	}
	if order, rest := parseStoreUnit(cmd, l); order.Errors == nil {
		return order, rest
	}
	return parseUnknown(cmd, l)
}

func parseStoreFactoryGroup(cmd *Lexeme, l []*Lexeme) (*StoreFactoryGroup, []*Lexeme) {
	var err error
	o := &StoreFactoryGroup{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.FactoryGroup, l, err = expectFactoryGroup(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("factoryGroup: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseStoreMineGroup(cmd *Lexeme, l []*Lexeme) (*StoreMineGroup, []*Lexeme) {
	var err error
	o := &StoreMineGroup{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.MineGroup, l, err = expectMineGroup(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("mineGroup: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseStoreUnit(cmd *Lexeme, l []*Lexeme) (*StoreUnit, []*Lexeme) {
	var err error
	o := &StoreUnit{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseSupport(cmd *Lexeme, l []*Lexeme) (any, []*Lexeme) {
	var err error
	sd := &SupportDefend{Line: cmd.Line}
	if sd.Id, l, err = expectInteger(l); err != nil {
		sd.Errors = append(sd.Errors, fmt.Errorf("id: %w", err))
		return sd, eatLine(l)
	}
	if sd.PctCommitted, l, err = expectPercentage(l); err != nil {
		sd.Errors = append(sd.Errors, fmt.Errorf("pctCommitted: %w", err))
		return sd, eatLine(l)
	}
	if sd.SupportId, l, err = expectInteger(l); err != nil {
		sd.Errors = append(sd.Errors, fmt.Errorf("supportId: %w", err))
		return sd, eatLine(l)
	}
	if targetId, rest, err := expectInteger(l); err == nil {
		// this is a support attack order
		sa := &SupportAttack{Line: sd.Line, Id: sd.Id, SupportId: sd.SupportId, PctCommitted: sd.PctCommitted, TargetId: targetId}
		l = rest
		if l, err = expectEOL(l); err != nil {
			sa.Errors = append(sa.Errors, err)
			return sa, eatLine(l)
		}
		return sa, l
	}
	if l, err = expectEOL(l); err != nil {
		sd.Errors = append(sd.Errors, err)
		return sd, eatLine(l)
	}
	return sd, l
}

func parseSuppressAgents(cmd *Lexeme, l []*Lexeme) (*SuppressAgents, []*Lexeme) {
	var err error
	o := &SuppressAgents{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.TargetId, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("targetId: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseSurvey(cmd *Lexeme, l []*Lexeme) (any, []*Lexeme) {
	var err error
	o := &Survey{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err == nil {
		return o, l
	}
	if o.Orbit, l, err = expectInteger(l); err == nil {
		if !(0 < o.Orbit && o.Orbit <= 10) {
			o.Errors = append(o.Errors, fmt.Errorf("orbit: invalid orbit %d", o.Orbit))
			return o, eatLine(l)
		}
		if l, err = expectEOL(l); err != nil {
			o.Errors = append(o.Errors, err)
			return o, eatLine(l)
		}
		return o, l
	}
	ss := &SurveySystem{Line: o.Line, Id: o.Id}
	if ss.Location, l, err = expectCoordinates(l); err != nil {
		ss.Errors = append(ss.Errors, fmt.Errorf("location: %w", err))
		return ss, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		ss.Errors = append(ss.Errors, err)
		return ss, eatLine(l)
	}
	return ss, l
}

func parseTransfer(cmd *Lexeme, l []*Lexeme) (*Transfer, []*Lexeme) {
	var err error
	o := &Transfer{Line: cmd.Line}
	if o.Id, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("id: %w", err))
		return o, eatLine(l)
	}
	if o.Quantity, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("quantity: %w", err))
		return o, eatLine(l)
	}
	if o.Unit, l, err = expectUnit(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("unit: %w", err))
		return o, eatLine(l)
	}
	if o.TargetId, l, err = expectInteger(l); err != nil {
		o.Errors = append(o.Errors, fmt.Errorf("targetId: %w", err))
		return o, eatLine(l)
	}
	if l, err = expectEOL(l); err != nil {
		o.Errors = append(o.Errors, err)
		return o, eatLine(l)
	}
	return o, l
}

func parseUnknown(cmd *Lexeme, l []*Lexeme) (*Unknown, []*Lexeme) {
	o := &Unknown{
		Line:    cmd.Line,
		Command: cmd.Text,
		Errors:  []error{fmt.Errorf("unknown command %q", cmd.Text)}}
	return o, eatLine(l)
}
