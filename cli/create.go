// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"context"
	"errors"
	"fmt"
	"github.com/playbymail/empyr/engine"
	"github.com/playbymail/empyr/pkg/clean"
	"github.com/playbymail/empyr/pkg/stdlib"
	"github.com/playbymail/empyr/store"
	"github.com/spf13/cobra"
	"log"
	"math/rand/v2"
	"time"
)

// this file implements the commands to create assets such as databases, games, and assets

// cmdCreate represents the base command when called without any subcommands
var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "create all the things",
	Long:  `create is the root of the generator commands.`,
}

// cmdCreateDatabase implements the create database command
var cmdCreateDatabase = &cobra.Command{
	Use:   "database",
	Short: "create a new database",
	Long:  `Create a new database.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		log.Printf("create: database: %q\n", env.Database.Path)
		if stdlib.IsExists(env.Database.Path) {
			if !env.Database.ForceCreate {
				log.Fatalf("error: %v\n", ErrFileExists)
			}
			log.Printf("create: database: deleting existing database\n")
			if err := stdlib.Remove(env.Database.Path); err != nil {
				log.Fatalf("error: stdlib.remove: %v\n", errors.Join(ErrDeleteFailed, err))
			}
		}
		log.Printf("create: database: %q\n", env.Database.Path)
		if err := store.Create(env.Database.Path); err != nil {
			log.Fatalf("error: store.create: %v\n", err)
		}
		log.Printf("create: database: completed in %v\n", time.Since(started))
	},
}

// cmdCreateGame implements the create game command
var cmdCreateGame = &cobra.Command{
	Use:   "game --code code --name name --descr description",
	Short: "create a new game",
	Long:  `Create a new game.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if _, err := clean.IsValidCode(cmd.Flag("code").Value.String()); err != nil {
			return err
		} else if _, err = clean.IsValidName(cmd.Flag("name").Value.String()); err != nil {
			return err
		} else if _, err = clean.IsValidDescription(cmd.Flag("descr").Value.String()); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("create: game: elapsed time: %v\n", time.Now().Sub(started))
		}()
		code := cmd.Flag("code").Value.String()
		name := cmd.Flag("name").Value.String()
		descr := cmd.Flag("descr").Value.String()
		if descr == "" {
			descr = fmt.Sprintf("A game of %s", name)
		}
		log.Printf("create: game: code  %q\n", code)
		log.Printf("create: game: name  %q\n", name)
		log.Printf("create: game: descr %q\n", descr)

		repo, err := store.Open(env.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}

		gameId, err := engine.CreateGameCommand(e, &engine.CreateGameParams_t{
			Code:        env.Game.Code,
			Name:        env.Game.Name,
			DisplayName: fmt.Sprintf("EC-%s", env.Game.Code),
			Rand:        rand.New(rand.NewPCG(0xdeadbeef, 0xcafedeed)),
			ForceCreate: env.Game.ForceCreate,
		})
		if err != nil {
			log.Fatalf("error: engine.CreateGameCommand: %v\n", err)
		}

		log.Printf("create: game: created game %d in %v\n", gameId, time.Since(started))
	},
}

// isWeird checks if a string contains any special characters or escape sequences
// by comparing the raw string with its quoted representation. It's a clever way
// to detect special characters by leveraging Go's string quoting behavior.
//
// Returns true if the string contains special characters, false otherwise.
// Example: isweird("hello") returns false, isweird("hello\n") returns true
func isWeird(s string) bool {
	return `"`+s+`"` != fmt.Sprintf("%q", s)
}

func hasBadBytes(s string) bool {
	for _, ch := range []byte(s) {
		if !goodBytes[ch] {
			return true
		}
	}
	return false
}

var (
	goodCodeBytes [256]bool
	goodBytes     [256]bool
)

func init() {
	for _, ch := range []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") {
		goodCodeBytes[ch] = true
	}
	for _, ch := range []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-.") {
		goodBytes[ch] = true
	}
}
