// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package cli implements the command line interface for empyr.
package cli

import (
	"github.com/playbymail/empyr/pkg/dotenv"
	"github.com/spf13/cobra"
	"log"
)

// Initialize returns a new cobra.Command that is initialized from the current environment.
// All configuration options are processed before we initialize the command flags.
//
// Flag values are initialized from the following sources:
//  1. Environment files (e.g., .env)
//  2. Environment variables (e.g., EMPYR_PATH)
//  3. Command line arguments
//
// The values are loaded in the order listed above, so that command line arguments
// will override any environment variables, which override any environment files.
func Initialize(options ...Option) (*cobra.Command, error) {
	// bootstrap the arguments
	env.Env.Prefix = "EMPYR"
	env.Reports.Path = "reports"

	// apply the options
	for _, option := range options {
		if err := option(); err != nil {
			return nil, err
		}
	}

	// load the env files, and then pull values from the environment.
	if err := dotenv.Load(env.Env.Prefix); err != nil {
		log.Fatalf("empyr: %+v\n", err)
	}
	applyEnvironmentVariables()

	cmdRoot.PersistentFlags().BoolVar(&env.Debug.DumpEnv, "dump-env", env.Debug.DumpEnv, "dump environment variables")

	cmdRoot.AddCommand(cmdCreate, cmdDB, cmdDelete, cmdShow, cmdVersion)

	cmdCreate.AddCommand(cmdCreateDatabase, cmdCreateEmpire, cmdCreateGame, cmdCreateStarList, cmdCreateSystemMap, cmdCreateTurnReport, cmdCreateTurnReports)
	cmdCreateDatabase.Flags().BoolVar(&env.Database.ForceCreate, "force-create", env.Database.ForceCreate, "force creation of the database")
	cmdCreateDatabase.Flags().StringVar(&env.Database.Path, "path", env.Database.Path, "path to the database")
	cmdCreateEmpire.Flags().StringVar(&env.Empire.Handle, "handle", env.Empire.Handle, "player handle for empire")
	cmdCreateGame.Flags().StringVar(&env.Game.Code, "code", env.Game.Code, "code for the game")
	if err := cmdCreateGame.MarkFlagRequired("code"); err != nil {
		return nil, err
	}
	cmdCreateGame.Flags().StringVar(&env.Game.Name, "name", env.Game.Name, "name of the game")
	if err := cmdCreateGame.MarkFlagRequired("name"); err != nil {
		return nil, err
	}
	cmdCreateGame.Flags().StringVar(&env.Game.Description, "descr", env.Game.Description, "description of the game")
	cmdCreateGame.Flags().BoolVar(&env.Game.ForceCreate, "force-create", env.Game.ForceCreate, "force creation of the game")
	cmdCreateTurnReport.Flags().Int64("empire-no", 0, "empire number for the report")
	if err := cmdCreateTurnReport.MarkFlagRequired("empire-no"); err != nil {
		return nil, err
	}
	cmdCreateTurnReport.Flags().Int64Var(&env.Game.TurnNo, "turn-no", env.Game.TurnNo, "turn number for the report")
	cmdCreateTurnReports.Flags().Int64Var(&env.Game.TurnNo, "turn-no", env.Game.TurnNo, "turn number for the report")

	cmdDB.PersistentFlags().String("path", "", "path to the database")
	if err := cmdDB.MarkPersistentFlagRequired("path"); err != nil {
		return nil, err
	}
	cmdDB.AddCommand(cmdDBCreate, cmdDBOpen)

	cmdDelete.AddCommand(cmdDeleteGame)
	cmdDeleteGame.Flags().StringVar(&env.Game.Code, "code", env.Game.Code, "code for the game")
	if err := cmdDeleteGame.MarkFlagRequired("code"); err != nil {
		return nil, err
	}

	cmdShow.AddCommand(cmdShowEnv)

	return cmdRoot, nil
}
