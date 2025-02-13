// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package resources

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrBadRequest = errors.New("bad request")
var ErrInvalidGenerator = errors.New("invalid generator")

// Resource is any resource that can be mined.
type Resource struct {
	ID              string
	Type            Type
	YieldPct        float64
	Unlimited       bool
	InitialAmount   int64
	AmountRemaining int64
}

func (nr Resource) MarshalJSON() ([]byte, error) {
	amtRemaining := nr.AmountRemaining
	if nr.Unlimited {
		amtRemaining = 99_999_999
	}
	return json.Marshal(&struct {
		ID              string `json:"deposit_id"`
		Type            string `json:"type"`
		Yield           int    `json:"yield"`
		AmountRemaining int64  `json:"remaining"`
	}{
		ID:              nr.ID,
		Type:            nr.Type.String(),
		Yield:           int(nr.YieldPct * 100),
		AmountRemaining: amtRemaining,
	})
}

type Type int

const (
	FUEL Type = iota
	GOLD
	METAL
	NONMETAL
)

func (t Type) String() string {
	switch t {
	case FUEL:
		return "fuel"
	case GOLD:
		return "gold"
	case METAL:
		return "metal"
	case NONMETAL:
		return "non-metal"
	}
	panic(fmt.Sprintf("assert(type != %d", t))
}
