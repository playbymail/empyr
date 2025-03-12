// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domains

type PlayerID int64

// Player represents a single player in the game.
// A player is a human or AI or just an NPC.
// Each player controls a single empire in any given game.
type Player struct {
	ID PlayerID // unique identifier for the player
}
