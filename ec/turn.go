// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package ec

import (
	"fmt"
	"github.com/playbymail/empyr/models/orders"
	"log"
	"sort"
)

func (e *Engine) AddOrders(orders []orders.Order) error {
	eo := &Orders{Orders: orders}
	// gather secrets
	for _, order := range eo.Orders {
		if secret, ok := order.(*Secret); ok {
			if eo.Secret != nil {
				return fmt.Errorf("multiple secrets")
			}
			eo.Secret = secret
			eo.Handle = secret.Handle
			eo.Game = secret.Game
			eo.Turn = secret.Turn
		}
	}
	if eo.Secret == nil {
		return fmt.Errorf("missing secret")
	}
	e.Orders = append(e.Orders, eo)
	return nil
}

func (e *Engine) Process() error {
	// process secrets phase
	for _, po := range e.Orders { // for each player orders
		err := e.SecretsPhase(po)
		if err != nil {
			// any error with secrets means the order file should be skipped
			po.Validated = false
		}
		if po.Validated {
			log.Printf("secrets: validated %s\n", po.Handle)
		} else if po.Error == nil {
			log.Printf("secrets: failed    %s\n", po.Handle)
		} else {
			log.Printf("secrets: failed    %s %v\n", po.Handle, po.Error)
		}
	}

	// sort orders by handle for consistent processing in future phases
	sort.Slice(e.Orders, func(i, j int) bool {
		if !e.Orders[i].Validated {
			return false
		}
		return e.Orders[i].Handle < e.Orders[j].Handle
	})

	return nil
}

func (e *Engine) SecretsPhase(orders *Orders) error {
	if orders.Secret == nil {
		orders.Error = fmt.Errorf("missing secret")
		return nil
	}
	secret := orders.Secret
	orders.Handle = secret.Handle
	orders.Game = secret.Game
	orders.Turn = secret.Turn

	// todo: unbake the secrets!
	switch secret.Token {
	case "003d626a-27c9-4f92-80f3-880384f22d08":
		if orders.Handle != "mdhender" {
			orders.Error = fmt.Errorf("invalid secret")
			return nil
		}
		if orders.Game != "G1" {
			orders.Error = fmt.Errorf("invalid game")
			return nil
		}
		if orders.Turn != 5 {
			orders.Error = fmt.Errorf("invalid turn")
			return nil
		}
		orders.Validated = true
	default:
		orders.Error = fmt.Errorf("invalid secret")
		return nil
	}
	return nil
}
