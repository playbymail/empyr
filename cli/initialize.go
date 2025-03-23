// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package cli implements the command line interface for empyr.
package cli

import (
	"github.com/playbymail/empyr/pkg/dotenv"
	"github.com/spf13/cobra"
	"log"
	"os"
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
	flags.Env.Prefix = "EMPYR"
	// apply the options
	for _, option := range options {
		if err := option(); err != nil {
			return nil, err
		}
	}

	// the environment flag must be set
	if env, ok := os.LookupEnv(flags.Env.Prefix + "_ENV"); !ok {
		return nil, ErrEnvFlagNotSet
	} else if !(env == "development" || env == "test" || env == "production") {
		return nil, ErrEnvFlagInvalid
	} else {
		flags.Environment = env
	}

	// load the env files, and then pull values from the environment.
	if err := dotenv.Load(flags.Env.Prefix); err != nil {
		log.Fatalf("empyr: %+v\n", err)
	}
	applyEnvironmentVariables()

	cmdRoot.PersistentFlags().BoolVar(&flags.Debug.DumpEnv, "dump-env", flags.Debug.DumpEnv, "dump environment variables")

	cmdRoot.AddCommand(cmdCreate, cmdDB, cmdDelete, cmdExecute, cmdShow, cmdStart, cmdVersion)

	cmdCreate.AddCommand(cmdCreateDatabase, cmdCreateEmpire, cmdCreateGame, cmdCreateStarList, cmdCreateSystemMap)

	cmdCreateDatabase.Flags().BoolVar(&flags.Database.ForceCreate, "force-create", flags.Database.ForceCreate, "force creation of the database")
	cmdCreateDatabase.Flags().StringVar(&flags.Database.Path, "path", flags.Database.Path, "path to the database")
	cmdCreateDatabase.Flags().StringVar(&flags.Game.Code, "code", flags.Game.Code, "code for the game")
	if err := cmdCreateDatabase.MarkFlagRequired("code"); err != nil {
		return nil, err
	}
	cmdCreateDatabase.Flags().StringVar(&flags.Game.Name, "name", flags.Game.Name, "name of the game")
	if err := cmdCreateDatabase.MarkFlagRequired("name"); err != nil {
		return nil, err
	}
	cmdCreateDatabase.Flags().StringVar(&flags.Game.Description, "descr", flags.Game.Description, "description of the game")

	cmdCreateEmpire.Flags().Int64("id", 0, "empire id for user")
	cmdCreateEmpire.Flags().StringVar(&flags.Empire.UserHandle, "user", flags.Empire.UserHandle, "user handle for empire")
	if err := cmdCreateEmpire.MarkFlagRequired("user"); err != nil {
		return nil, err
	}

	cmdDB.PersistentFlags().String("path", "", "path to the database")
	if err := cmdDB.MarkPersistentFlagRequired("path"); err != nil {
		return nil, err
	}
	cmdDB.AddCommand(cmdDBCreate, cmdDBOpen)

	cmdExecute.AddCommand(cmdExecuteProbes, cmdExecuteReset, cmdExecuteSurveys)

	cmdShow.AddCommand(cmdShowEnv)

	return cmdRoot, nil
}
