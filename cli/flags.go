// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"encoding/json"
	"fmt"
	"github.com/mdhender/semver"
	"github.com/mdhender/xii"
	"log"
	"time"
)

// this file defines the command line argument flags structure

var flags struct {
	Environment string // development, test, production
	Env         struct {
		Prefix string
	}
	Application struct {
		Assets struct {
			Public    string
			Templates string
		}
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
	Empire struct {
		Id         int64  // database key
		No         int64  // empire number in game
		UserHandle string // player handle
	}
	Game struct {
		Code        string
		Name        string
		Description string
		TurnNo      int64
		ForceCreate bool
	}
	Reports struct {
		Path string
	}
	Server struct {
		Scheme         string
		Host           string
		Port           string
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
		IdleTimeout    time.Duration
		MaxHeaderBytes int64
	}
	Sessions struct {
		Secret string
	}
	Verbose bool
	Version semver.Version
}

// painfully apply the environment variables to the arguments
func applyEnvironmentVariables() {
	xiistr(&flags.Application.Assets.Public, "_APP_ASSETS_PUBLIC")
	xiistr(&flags.Application.Assets.Templates, "_APP_ASSETS_TEMPLATES")

	xiistr(&flags.Database.Path, "_DATABASE_PATH")
	xiibool(&flags.Database.DryRun, "_DATABASE_DRYRUN")
	xiibool(&flags.Database.ForceCreate, "_DATABASE_FORCECREATE")

	xiibool(&flags.Debug.Database.Open, "_DEBUG_DATABASE_OPEN")
	xiibool(&flags.Debug.Database.Close, "_DEBUG_DATABASE_CLOSE")
	xiibool(&flags.Debug.DumpEnv, "_DEBUG_DUMPARGS")

	xiiint(&flags.Empire.Id, "_EMPIRE_ID")
	xiiint(&flags.Empire.No, "_EMPIRE_NO")
	xiistr(&flags.Empire.UserHandle, "_EMPIRE_USERHANDLE")

	xiistr(&flags.Game.Code, "_GAME_CODE")
	xiistr(&flags.Game.Name, "_GAME_NAME")
	xiistr(&flags.Game.Description, "_GAME_DESCRIPTION")
	xiiint(&flags.Game.TurnNo, "_GAME_TURNNO")
	xiibool(&flags.Game.ForceCreate, "_GAME_FORCECREATE")

	xiistr(&flags.Reports.Path, "_REPORTS_PATH")

	xiistr(&flags.Server.Scheme, "_SERVER_SCHEME")
	xiistr(&flags.Server.Host, "_SERVER_HOST")
	xiistr(&flags.Server.Port, "_SERVER_PORT")
	xiiint(&flags.Server.MaxHeaderBytes, "_SERVER_MAXHEADERBYTES")

	xiistr(&flags.Sessions.Secret, "_SESSIONS_SECRET")

	xiibool(&flags.Verbose, "_VERBOSE")
}

func dumpEnv(toLog bool) {
	data, err := json.MarshalIndent(flags, "", "  ")
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
	if val, err := xii.AsBool(flags.Env.Prefix+varname, xii.BoolOpts{DefaultValue: *b}); err == nil {
		*b = val
	}
}

// xiiint is a helper function to apply the environment variable to an integer
func xiiint(i *int64, varname string) {
	if val, err := xii.AsInt(flags.Env.Prefix+varname, xii.IntOpts{DefaultValue: int(*i)}); err == nil {
		*i = int64(val)
	}
}

// xiistr is a helper function to apply the environment variable to a string
func xiistr(s *string, varname string) {
	if val, err := xii.AsString(flags.Env.Prefix+varname, xii.StringOpts{DefaultValue: *s}); err == nil {
		*s = val
	}
}
