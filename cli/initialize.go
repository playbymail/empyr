// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package cli implements the command line interface for empyr.
package cli

import (
	"github.com/playbymail/empyr/pkg/dotenv"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
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
	flags.Application.Assets.Public = "app/assets/public"
	flags.Application.Assets.Templates = "app/assets/templates"
	flags.Reports.Path = "reports"
	flags.Surveys.Path = "surveys"
	flags.Server.Scheme = "http"
	flags.Server.Host = "localhost"
	flags.Server.Port = "8080"
	flags.Server.ReadTimeout = 5 * time.Second
	flags.Server.WriteTimeout = 10 * time.Second
	flags.Server.IdleTimeout = 120 * time.Second
	flags.Server.MaxHeaderBytes = 1 << 20

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

	cmdCreate.AddCommand(cmdCreateDatabase, cmdCreateEmpire, cmdCreateGame, cmdCreateStarList, cmdCreateSystemMap, cmdCreateSystemSurveyReport, cmdCreateSystemSurveyReports, cmdCreateTurnReport, cmdCreateTurnReports, cmdCreateUser)
	cmdCreateDatabase.Flags().BoolVar(&flags.Database.ForceCreate, "force-create", flags.Database.ForceCreate, "force creation of the database")
	cmdCreateDatabase.Flags().StringVar(&flags.Database.Path, "path", flags.Database.Path, "path to the database")
	cmdCreateEmpire.Flags().StringVar(&flags.Empire.UserHandle, "user", flags.Empire.UserHandle, "user handle for empire")
	if err := cmdCreateEmpire.MarkFlagRequired("user"); err != nil {
		return nil, err
	}
	cmdCreateGame.Flags().StringVar(&flags.Game.Code, "code", flags.Game.Code, "code for the game")
	if err := cmdCreateGame.MarkFlagRequired("code"); err != nil {
		return nil, err
	}
	cmdCreateGame.Flags().StringVar(&flags.Game.Name, "name", flags.Game.Name, "name of the game")
	if err := cmdCreateGame.MarkFlagRequired("name"); err != nil {
		return nil, err
	}
	cmdCreateGame.Flags().StringVar(&flags.Game.Description, "descr", flags.Game.Description, "description of the game")
	cmdCreateGame.Flags().BoolVar(&flags.Game.ForceCreate, "force-create", flags.Game.ForceCreate, "force creation of the game")
	cmdCreateSystemSurveyReport.Flags().Int64("empire-no", 0, "empire number for the report")
	if err := cmdCreateSystemSurveyReport.MarkFlagRequired("empire-no"); err != nil {
		return nil, err
	}
	cmdCreateSystemSurveyReport.Flags().Int64Var(&flags.Game.TurnNo, "turn-no", flags.Game.TurnNo, "turn number for the report")
	cmdCreateSystemSurveyReports.Flags().Int64Var(&flags.Game.TurnNo, "turn-no", flags.Game.TurnNo, "turn number for the report")
	cmdCreateTurnReport.Flags().Int64("empire-no", 0, "empire number for the report")
	if err := cmdCreateTurnReport.MarkFlagRequired("empire-no"); err != nil {
		return nil, err
	}
	cmdCreateTurnReport.Flags().Int64Var(&flags.Game.TurnNo, "turn-no", flags.Game.TurnNo, "turn number for the report")
	cmdCreateTurnReports.Flags().Int64Var(&flags.Game.TurnNo, "turn-no", flags.Game.TurnNo, "turn number for the report")
	cmdCreateUser.Flags().Bool("is-admin", false, "admin user")
	cmdCreateUser.Flags().String("email", "", "user email")
	if err := cmdCreateUser.MarkFlagRequired("email"); err != nil {
		return nil, err
	}
	cmdCreateUser.Flags().String("password", "", "user password")
	cmdCreateUser.Flags().String("user", "", "user name")
	if err := cmdCreateUser.MarkFlagRequired("user"); err != nil {
		return nil, err
	}

	cmdDB.PersistentFlags().String("path", "", "path to the database")
	if err := cmdDB.MarkPersistentFlagRequired("path"); err != nil {
		return nil, err
	}
	cmdDB.AddCommand(cmdDBCreate, cmdDBOpen)

	cmdDelete.AddCommand(cmdDeleteGame)
	cmdDeleteGame.Flags().StringVar(&flags.Game.Code, "code", flags.Game.Code, "code for the game")
	if err := cmdDeleteGame.MarkFlagRequired("code"); err != nil {
		return nil, err
	}

	cmdExecute.AddCommand(cmdExecuteProbes, cmdExecuteReset, cmdExecuteSurveys)

	cmdShow.AddCommand(cmdShowEnv)

	cmdStart.AddCommand(cmdStartServer)

	return cmdRoot, nil
}
