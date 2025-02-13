// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

// Engine implements the rules of the game.
// Input is always game state and an order.
// Output is updated game state.
type Engine struct {
}

// NewEngine returns an initialized engine.
func NewEngine() (*Engine, error) {
	e := &Engine{}
	e.initialize()
	return e, nil
}

func (e *Engine) initialize() {
}

func (e *Engine) reset() {
	e.initialize()
}
