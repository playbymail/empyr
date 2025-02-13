// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/playbymail/empyr/models/games"
	"github.com/playbymail/empyr/models/player"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var argsGenerateGame struct {
	force bool   // if true, delete any existing game
	code  string // code for game
	name  string // name of game
	descr string // description of game
}

// cmdGenerateGame runs the game generator command
var cmdGenerateGame = &cobra.Command{
	Use:   "game",
	Short: "generate a new game",
	RunE: func(cmd *cobra.Command, args []string) error {
		// isweird checks if a string contains any special characters or escape sequences
		// by comparing the raw string with its quoted representation.
		// It's a clever way to detect special characters by leveraging Go's string quoting behavior.
		//
		// Returns true if the string contains special characters, false otherwise.
		// Example: isweird("hello") returns false, isweird("hello\n") returns true
		isweird := func(s string) bool {
			return `"`+s+`"` != fmt.Sprintf("%q", s)
		}

		if argsGenerate.path != "" {
			argsGenerate.path = filepath.Clean(argsGenerate.path)
			if sb, err := os.Stat(argsGenerate.path); err != nil {
				return err
			} else if !sb.IsDir() {
				return fmt.Errorf("path must be a valid directory")
			}
		}

		argsGenerateGame.code = strings.ToUpper(argsGenerateGame.code)
		if len(argsGenerateGame.code) == 0 {
			return fmt.Errorf("missing code")
		} else if !unicode.IsLetter(rune(argsGenerateGame.code[0])) {
			return fmt.Errorf("code must start with letter")
		} else if isweird(argsGenerateGame.code) {
			return fmt.Errorf("code must not contain weird characters")
		}

		if argsGenerateGame.name == "" {
			// name will default to code
			argsGenerateGame.name = argsGenerateGame.code
		} else if isweird(argsGenerateGame.name) {
			return fmt.Errorf("name must not contain weird characters")
		}

		if argsGenerateGame.descr == "" {
			// description will default to name
			argsGenerateGame.descr = argsGenerateGame.name
		} else if isweird(argsGenerateGame.descr) {
			return fmt.Errorf("descr must not contain weird characters")
		}

		gamePath := filepath.Join(argsGenerate.path, argsGenerateGame.code)
		if sb, err := os.Stat(gamePath); err == nil {
			if !argsGenerateGame.force {
				log.Fatalf("generate: game: path: %s exists\n", gamePath)
			} else if !sb.IsDir() {
				log.Fatalf("generate: game: path: %s is not a directory\n", gamePath)
			} else if err = os.RemoveAll(gamePath); err != nil { // remove it
				log.Fatalf("generate: game: %v\n", err)
			}
		}
		log.Printf("generate: game: creating %s\n", gamePath)
		if err := os.MkdirAll(gamePath, 0755); err != nil {
			log.Fatal(err)
		}
		log.Printf("generate: game: created  %s\n", gamePath)

		g := games.Game{
			Id:    uuid.New().String(),
			Code:  argsGenerateGame.code,
			Name:  argsGenerateGame.name,
			Descr: argsGenerateGame.descr,
		}

		if buffer, err := json.MarshalIndent(&g, "", "\t"); err != nil {
			log.Fatal(err)
		} else if err = os.WriteFile(filepath.Join(gamePath, "game.json"), buffer, 0664); err != nil {
			log.Fatal(err)
		}

		if buffer, err := json.MarshalIndent([]player.Player{}, "", "\t"); err != nil {
			log.Fatal(err)
		} else if err = os.WriteFile(filepath.Join(gamePath, "players.json"), buffer, 0664); err != nil {
			log.Fatal(err)
		}

		return nil
	},
}
