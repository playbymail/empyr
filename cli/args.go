// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"encoding/json"
	"fmt"
	"github.com/mdhender/semver"
	"github.com/mdhender/xii"
	"log"
)

// this file defines the command line arguments structure

var env struct {
	Env struct {
		Prefix string
	}
	Database struct {
		Path        string
		DryRun      bool
		ForceCreate bool
	}
	Debug struct {
		Database struct {
			Open  bool
			Close bool
		}
		DumpEnv bool
	}
	Game struct {
		Code        string
		Name        string
		Description string
		TurnNo      int64
		ForceCreate bool
	}
	Verbose bool
	Version semver.Version
}

// painfully apply the environment variables to the arguments
func applyEnvironmentVariables() {
	xiistr(&env.Database.Path, "_DATABASE_PATH")
	xiibool(&env.Database.DryRun, "_DATABASE_DRYRUN")
	xiibool(&env.Database.ForceCreate, "_DATABASE_FORCECREATE")

	xiibool(&env.Debug.Database.Open, "_DEBUG_DATABASE_OPEN")
	xiibool(&env.Debug.Database.Close, "_DEBUG_DATABASE_CLOSE")
	xiibool(&env.Debug.DumpEnv, "_DEBUG_DUMPARGS")

	xiistr(&env.Game.Code, "_GAME_CODE")
	xiistr(&env.Game.Name, "_GAME_NAME")
	xiistr(&env.Game.Description, "_GAME_DESCRIPTION")
	xiiint(&env.Game.TurnNo, "_GAME_TURNNO")
	xiibool(&env.Game.ForceCreate, "_GAME_FORCECREATE")

	xiibool(&env.Verbose, "_VERBOSE")
}

func dumpEnv(toLog bool) {
	data, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		panic(err)
	}
	if toLog {
		log.Printf("env: %s\n", data)
		return
	}
	fmt.Printf("%s\n", data)
}

// xiibool is a helper function to apply the environment variable to a boolean
func xiibool(b *bool, varname string) {
	if val, err := xii.AsBool(env.Env.Prefix+varname, xii.BoolOpts{DefaultValue: *b}); err == nil {
		*b = val
	}
}

// xiiint is a helper function to apply the environment variable to an integer
func xiiint(i *int64, varname string) {
	if val, err := xii.AsInt(env.Env.Prefix+varname, xii.IntOpts{DefaultValue: int(*i)}); err == nil {
		*i = int64(val)
	}
}

// xiistr is a helper function to apply the environment variable to a string
func xiistr(s *string, varname string) {
	if val, err := xii.AsString(env.Env.Prefix+varname, xii.StringOpts{DefaultValue: *s}); err == nil {
		*s = val
	}
}
