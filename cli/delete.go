// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"context"
	"github.com/playbymail/empyr/engine"
	"github.com/playbymail/empyr/store"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// this file implements the commands to delete assets such as databases, games, and assets

// cmdDelete represents the base command when called without any subcommands
var cmdDelete = &cobra.Command{
	Use:   "delete",
	Short: "delete all the things",
	Long:  `delete is the root of the removal commands.`,
}

// cmdDeleteGame implements the delete game command
var cmdDeleteGame = &cobra.Command{
	Use:   "game --code code",
	Short: "delete a game",
	Long:  `Delete a game.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if _, err := engine.IsValidCode(cmd.Flag("code").Value.String()); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("delete: game: elapsed time: %v\n", time.Now().Sub(started))
		}()
		code := cmd.Flag("code").Value.String()
		log.Printf("delete: game: code  %q\n", code)

		repo, err := store.Open(env.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}

		err = engine.DeleteGameCommand(e, &engine.DeleteGameParams_t{
			Code: env.Game.Code,
		})
		if err != nil {
			log.Fatalf("error: engine.DeleteGameCommand: %v\n", err)
		}

		log.Printf("delete: game: deleted game %q in %v\n", code, time.Since(started))
	},
}
