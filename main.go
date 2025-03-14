// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/mdhender/semver"
	"github.com/playbymail/empyr/cli"
	"github.com/playbymail/empyr/config"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	version = semver.Version{Minor: 5}
)

func main() {
	// versions hack
	for _, arg := range os.Args {
		if arg == "version" || arg == "-version" || arg == "--version" {
			fmt.Printf("%s\n", version.String())
			os.Exit(0)
		}
	}

	log.SetFlags(log.Lshortfile)

	if command, err := cli.Initialize(
		cli.WithVersion(version),
		cli.WithEnvPrefix("EMPYR"),
	); err != nil {
		log.Fatalf("main: %+v\n", err)
	} else if err = command.Execute(); err != nil {
		log.Fatalf("\n%+v\n", err)
	}
}

func runRoot(env *config.Environment, args []string) error {
	var arg string
	for len(args) > 0 && args[0] != "--" {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			break
		} else if arg == "start" {
			if err := runStart(env, args); err != nil {
				return errors.Join(fmt.Errorf("start"), err)
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
	log.Printf("usage: empyr [options] start [command] [options] [arguments]\n")
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
	log.Printf("usage: empyr [options] start [options] server [options] [arguments]\n")
	log.Printf("  opt: --help          show help for the command   [false]\n")
	log.Printf("     : --debug=flag    enable various debug flags  [off]\n")
	log.Printf("     : --verbose       enhance logging             [false]\n")
	log.Printf("     : --db=path       path to database            [required]\n")
	log.Printf("  arg: --port=port     port to listen on           [8080]\n")
	log.Printf("     : --host=host     host to bind to             [localhost]\n")
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

func argOptInt(arg, flag string) (value int64, ok bool) {
	flag = "--" + flag + "="
	if strings.HasPrefix(arg, flag) {
		if n, err := strconv.Atoi(strings.TrimPrefix(arg, flag)); err == nil {
			return int64(n), true
		}
	}
	return 0, false
}

func argOptString(arg, flag string) (value string, ok bool) {
	flag = "--" + flag + "="
	if strings.HasPrefix(arg, flag) {
		return strings.TrimPrefix(arg, flag), true
	}
	return "", false
}
