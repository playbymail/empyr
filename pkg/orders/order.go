// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package orders

type Orders struct {
	Game *Game
	Auth *Auth
}

type Auth struct {
	Kind  string
	Value string
}

type Game struct {
	Name string
	Turn int
}
