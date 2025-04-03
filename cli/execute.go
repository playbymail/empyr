// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"context"
	"github.com/playbymail/empyr/engine"
	"github.com/playbymail/empyr/repos"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// this file implements the commands to execute orders for a game turn

// cmdExecute represents the base command when called without any subcommands
var cmdExecute = &cobra.Command{
	Use:   "execute",
	Short: "execute orders",
	Long:  `execute is the root of the execution commands.`,
}

var cmdExecuteProbes = &cobra.Command{
	Use:   "probes",
	Short: "execute probe orders",
	Long:  `execute probe orders.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("execute: probes: elapsed time: %v\n", time.Now().Sub(started))
		}()
		log.Printf("execute: probes: game %q\n", flags.Game.Code)
		repo, err := repos.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}
		_ = e
	},
}

var cmdExecuteReset = &cobra.Command{
	Use:   "reset",
	Short: "execute reset turn results",
	Long:  `Deletes all results for the current turn in the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("execute: reset: elapsed time: %v\n", time.Now().Sub(started))
		}()
		log.Printf("execute: reset: game %q\n", flags.Game.Code)
		repo, err := repos.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()

		// reset the turn results deletes all reports, probes, and surveys.
		err = repo.ResetTurnResults(flags.Game.Code)
		if err != nil {
			log.Fatalf("error: reset: %v\n", err)
		}
	},
}

var cmdExecuteSurveys = &cobra.Command{
	Use:   "surveys",
	Short: "execute survey orders",
	Long:  `execute survey orders.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("execute: surveys: elapsed time: %v\n", time.Now().Sub(started))
		}()
		log.Printf("execute: surveys: game %q\n", flags.Game.Code)
		repo, err := repos.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}
		err = e.ExecuteSurveys(flags.Game.Code, 0)
	},
}
