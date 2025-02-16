// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"errors"
	"fmt"
	"github.com/mdhender/semver"
	"github.com/playbymail/empyr/config"
	"github.com/playbymail/empyr/engine"
	"github.com/playbymail/empyr/pkg/dotenv"
	"github.com/playbymail/empyr/pkg/empyr"
	"github.com/playbymail/empyr/pkg/stdlib"
	"github.com/playbymail/empyr/store"
	"log"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

var (
	version = semver.Version{Minor: 1}
)

func main() {
	log.SetFlags(log.Lshortfile)

	started := time.Now()
	defer func() {
		log.Printf("elapsed time: %v\n", time.Now().Sub(started))
	}()

	// options for the configuration will be pulled from the global command line flags
	options := []config.Option{config.WithVersion(version)}

	// build the arguments list and deal with global command line flags
	var args []string
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if flag, ok := argOptString(arg, "debug"); ok {
			options = append(options, config.WithDebugFlag(flag))
		} else if flag, ok := argOptBool(arg, "verbose"); ok {
			options = append(options, config.WithVerbose(flag))
		} else {
			args = append(args, arg)
		}
	}

	// get default configuration. uses the .env file if it exists, and
	// then pulls values from the environment.
	if err := dotenv.Load("EMPYR"); err != nil {
		log.Fatalf("main: %+v\n", err)
	}
	for _, v := range dotenv.EnvVariables("EMPYR_") {
		log.Printf("main: %-22s == %q\n", v.Key, v.Value)
	}
	env := config.Default(options...)

	if err := runRoot(env, args); err != nil {
		log.Fatalf("main: %+v\n", err)
	}

	//if command, err := cli.Initialize(
	//	cli.WithVersion(version),
	//); err != nil {
	//	log.Fatalf("main: %+v\n", err)
	//} else if err = command.Execute(); err != nil {
	//	log.Fatalf("\n%+v\n", err)
	//}

	log.Printf("\n")
	log.Printf("completed in %v\n", time.Now().Sub(started))
}

func runRoot(env *config.Environment, args []string) error {
	var arg string
	for len(args) > 0 && args[0] != "--" {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			break
		} else if arg == "create" {
			if err := runCreate(env, args); err != nil {
				return fmt.Errorf("create: %w", err)
			}
			return nil
		} else if arg == "start" {
			if err := runStart(env, args); err != nil {
				return errors.Join(fmt.Errorf("start"), err)
			}
			return nil
		} else if arg == "version" {
			if err := runVersion(env, args); err != nil {
				return errors.Join(fmt.Errorf("version"), err)
			}
			return nil
		} else if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("unknown option: %q", arg)
		} else {
			return fmt.Errorf("unknown command: %q", arg)
		}
	}
	log.Printf("usage: empyr [command] [options] [arguments]\n")
	log.Printf("  cmd: help            show help for the application\n")
	log.Printf("     : create          create things\n")
	log.Printf("     : server          start the web server\n")
	log.Printf("     : version         show the application version\n")
	log.Printf("  opt: --help          show help for the command   [false]\n")
	log.Printf("     : --debug=flag    enable various debug flags  [off]\n")
	log.Printf("     : --verbose       enhance logging             [false]\n")
	return nil
}

func runCreate(env *config.Environment, args []string) error {
	var arg string
	for len(args) > 0 && args[0] != "-- " {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			break
		} else if arg == "database" {
			if err := runCreateDatabase(env, args); err != nil {
				return fmt.Errorf("database: %w", err)
			}
			return nil
		} else if arg == "game" {
			if err := runCreateGame(env, args); err != nil {
				return fmt.Errorf("game: %w", err)
			}
			return nil
		} else if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("unknown option: %q", arg)
		} else {
			return fmt.Errorf("unknown command: %q", arg)
		}
	}
	log.Printf("usage: empyr create [command] [options] [arguments]\n")
	log.Printf("  cmd: database        create a new database\n")
	log.Printf("     : game            create a new game\n")
	log.Printf("  opt: --help          show help for the command   [false]\n")
	log.Printf("     : --debug=flag    enable various debug flags  [off]\n")
	log.Printf("     : --verbose       enhance logging             [false]\n")
	return nil
}

func runCreateDatabase(env *config.Environment, args []string) error {
	var arg string
	for len(args) > 0 && args[0] != "-- " {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			log.Printf("usage: empyr create database [options] path_to_database\n")
			log.Printf("  opt: --help          show help for the command   [false]\n")
			log.Printf("     : --debug=flag    enable various debug flags  [off]\n")
			log.Printf("     : --verbose       enhance logging             [false]\n")
			log.Printf("     : --force         force creation              [false]\n")
			return nil
		} else if flag, ok := argOptBool(arg, "force"); ok {
			env.Database.ForceCreate = flag
		} else if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("unknown option: %q", arg)
		} else if env.Database.Path == "" {
			env.Database.Path = arg
		} else {
			return fmt.Errorf("unknown argument: %q", arg)
		}
	}
	if env.Database.Path == "" {
		return fmt.Errorf("missing path to database")
	} else if stdlib.IsExists(env.Database.Path) {
		if !env.Database.ForceCreate {
			return fmt.Errorf("%q: already exists", env.Database.Path)
		}
		log.Printf("%q: deleting existing database\n", env.Database.Path)
		if err := stdlib.Remove(env.Database.Path); err != nil {
			return fmt.Errorf("%q: %w", env.Database.Path, err)
		}
	}
	log.Printf("%q: creating database\n", env.Database.Path)
	if err := store.Create(env.Database.Path); err != nil {
		return fmt.Errorf("%q: %w", env.Database.Path, err)
	}
	log.Printf("%s: created database\n", env.Database.Path)
	return nil
}

func runCreateGame(env *config.Environment, args []string) error {
	var arg string
	for len(args) > 0 && args[0] != "-- " {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			log.Printf("usage: empyr create game [options]\n")
			log.Printf("  opt: --help          show help for the command   [false]\n")
			log.Printf("     : --debug=flag    enable various debug flags  [off]\n")
			log.Printf("     : --verbose       enhance logging             [false]\n")
			log.Printf("     : --path          path to database            [required]\n")
			return nil
		} else if flag, ok := argOptString(arg, "code"); ok {
			env.Game.Code = flag
		} else if flag, ok := argOptString(arg, "name"); ok {
			env.Game.Name = flag
		} else if flag, ok := argOptString(arg, "path"); ok {
			env.Database.Path = flag
		} else if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("unknown option: %q", arg)
		} else {
			return fmt.Errorf("unknown argument: %q", arg)
		}
	}
	if env.Database.Path == "" {
		return fmt.Errorf("missing path to database")
	} else if !stdlib.IsFileExists(env.Database.Path) {
		return fmt.Errorf("%q: does not exist", env.Database.Path)
	}
	if env.Game.Code == "" {
		return fmt.Errorf("missing game code")
	} else if strings.ToUpper(env.Game.Code) != env.Game.Code {
		return fmt.Errorf("%q: code must be uppercase", env.Game.Code)
	} else if strings.TrimSpace(env.Game.Code) != env.Game.Code {
		return fmt.Errorf("%q: code must not contain whitespace", env.Game.Code)
	}
	log.Printf("game: code: %q\n", env.Game.Code)
	if env.Game.Name == "" {
		return fmt.Errorf("missing game name")
	}
	log.Printf("game: name: %q\n", env.Game.Name)
	log.Printf("%q: creating game\n", env.Database.Path)
	r := rand.New(rand.NewPCG(0xdeadbeef, 0xcafedeed))
	gc, err := engine.CreateCluster(r)
	if err != nil {
		return fmt.Errorf("cluster: %w", err)
	}
	log.Printf("%q: creating cluster %p\n", env.Database.Path, gc)
	_, err = empyr.NewGame(env.Game.Code, env.Game.Name)
	if err != nil {
		return fmt.Errorf("code: %q: %w", err)
	}
	log.Printf("%q: created game\n", env.Database.Path)
	return nil
}

func runStart(env *config.Environment, args []string) error {
	var arg string
	for len(args) > 0 && args[0] != "-- " {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			break
		} else if arg == "server" {
			if err := runStartServer(env, args); err != nil {
				return errors.Join(fmt.Errorf("server"), err)
			}
			return nil
		} else if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("unknown option: %q", arg)
		} else {
			return fmt.Errorf("unknown command: %q", arg)
		}
	}
	log.Printf("usage: empyr start [command] [options] [arguments]\n")
	log.Printf("  cmd: server          start the web server\n")
	log.Printf("  opt: --help          show help for the command   [false]\n")
	log.Printf("     : --debug=flag    enable various debug flags  [off]\n")
	log.Printf("     : --verbose       enhance logging             [false]\n")
	log.Printf("     : --db=path       path to database            [required]\n")
	return nil
}

func runStartServer(env *config.Environment, args []string) error {
	var arg string
	for len(args) > 0 && args[0] != "-- " {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			break
		} else if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("unknown option: %q", arg)
		} else {
			return fmt.Errorf("unknown command: %q", arg)
		}
	}
	log.Printf("usage: empyr start server [options] [arguments]\n")
	log.Printf("  opt: --help          show help for the command   [false]\n")
	log.Printf("     : --debug=flag    enable various debug flags  [off]\n")
	log.Printf("     : --verbose       enhance logging             [false]\n")
	log.Printf("     : --db=path       path to database            [required]\n")
	log.Printf("  arg: --port=port     port to listen on           [8080]\n")
	log.Printf("     : --host=host     host to bind to             [localhost]\n")
	return nil
}

func runVersion(env *config.Environment, args []string) error {
	log.Printf("version: %s\n", env.Version.String())
	return nil
}

func argOptBool(arg, flag string) (value bool, ok bool) {
	if arg == "--"+flag {
		return true, true
	} else if arg == "--no-"+flag {
		return false, true
	}
	return false, false
}

func argOptHelp(arg string) bool {
	return arg == "--help" || arg == "-h" || arg == "help" || arg == "?" || arg == "-?" || arg == "/?"
}

func argOptString(arg, flag string) (value string, ok bool) {
	flag = "--" + flag + "="
	if strings.HasPrefix(arg, flag) {
		return strings.TrimPrefix(arg, flag), true
	}
	return "", false
}
