// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mdhender/semver"
	"github.com/playbymail/empyr/cli"
	"github.com/playbymail/empyr/config"
	"github.com/playbymail/empyr/engine"
	"github.com/playbymail/empyr/pkg/dotenv"
	"github.com/playbymail/empyr/pkg/stdlib"
	"github.com/playbymail/empyr/store"
	"html/template"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	version = semver.Version{Minor: 1}
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
	if version.String() != "0.0.0" {
		log.Printf("empyr version %s\n", version.String())
		os.Exit(0)
	}

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
	//for _, v := range dotenv.EnvVariables("EMPYR_") {
	//	log.Printf("main: %-22s == %q\n", v.Key, v.Value)
	//}
	env := config.Default(options...)

	if err := runRoot(env, args); err != nil {
		log.Fatalf("main: %+v\n", err)
	}

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
		} else if arg == "show" {
			if err := runShow(env, args); err != nil {
				return errors.Join(fmt.Errorf("show"), err)
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
		} else if arg == "cluster" {
			if err := runCreateCluster(env, args); err != nil {
				return fmt.Errorf("cluster: %w", err)
			}
			return nil
		} else if arg == "turn-reports" {
			if err := runCreateTurnReports(env, args); err != nil {
				return fmt.Errorf("turn-reports: %w", err)
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
	log.Printf("     : cluster         create cluster assets\n")
	log.Printf("     : turn-reports    create turn reports for all empires\n")
	log.Printf("  opt: --help          show help for the command   [false]\n")
	log.Printf("     : --debug=flag    enable various debug flags  [off]\n")
	log.Printf("     : --verbose       enhance logging             [false]\n")
	return nil
}

func runCreateCluster(env *config.Environment, args []string) error {
	var arg string
	for len(args) > 0 && args[0] != "-- " {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			log.Printf("usage: empyr create cluster [command] [options]\n")
			log.Printf("  opt: --help          show help for the command   [false]\n")
			log.Printf("     : --debug=flag    enable various debug flags  [off]\n")
			log.Printf("     : --verbose       enhance logging             [false]\n")
			log.Printf("     : --path          path to database            [required]\n")
			return nil
		} else if arg == "empire" {
			if err := runCreateClusterEmpire(env, args); err != nil {
				return fmt.Errorf("empire: %w", err)
			}
			return nil
		} else if arg == "star-list" {
			if err := runCreateClusterStarList(env, args); err != nil {
				return fmt.Errorf("star-list: %w", err)
			}
			return nil
		} else if arg == "system-map" {
			if err := runCreateClusterSystemMap(env, args); err != nil {
				return fmt.Errorf("system-map: %w", err)
			}
			return nil
		} else if flag, ok := argOptString(arg, "code"); ok {
			env.Game.Code = flag
		} else if flag, ok := argOptString(arg, "path"); ok {
			env.Database.Path = flag
		} else if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("unknown option: %q", arg)
		} else {
			return fmt.Errorf("unknown argument: %q", arg)
		}
	}
	return fmt.Errorf("missing command")
}

func runCreateClusterEmpire(env *config.Environment, args []string) error {
	var playerHandle string
	var arg string
	for len(args) > 0 && args[0] != "-- " {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			log.Printf("usage: empyr create empire [options]\n")
			log.Printf("  opt: --help          show help for the command   [false]\n")
			log.Printf("     : --debug=flag    enable various debug flags  [off]\n")
			log.Printf("     : --verbose       enhance logging             [false]\n")
			log.Printf("     : --path          path to database            [required]\n")
			log.Printf("     : --handle        player handle.              [required]\n")
			return nil
		} else if flag, ok := argOptString(arg, "code"); ok {
			env.Game.Code = flag
		} else if flag, ok := argOptString(arg, "handle"); ok {
			playerHandle = flag
		} else if flag, ok := argOptString(arg, "path"); ok {
			env.Database.Path = flag
		} else if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("unknown option: %q", arg)
		} else {
			return fmt.Errorf("unknown argument: %q", arg)
		}
	}
	if playerHandle == "" {
		return fmt.Errorf("missing player handle")
	} else if strings.TrimSpace(playerHandle) != playerHandle {
		return fmt.Errorf("player handle: must not have whitespace at head or tail")
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

	var err error
	if env.Store, err = store.Open(env.Database.Path, context.Background()); err != nil {
		return fmt.Errorf("%q: %w", env.Database.Path, err)
	}
	defer env.Store.Close()
	e, err := engine.Open(env.Store)
	if err != nil {
		return err
	}

	empireId, empireNo, err := engine.CreateEmpireCommand(e, &engine.CreateEmpireParams_t{
		Code:   env.Game.Code,
		Handle: playerHandle,
	})
	if err != nil {
		return fmt.Errorf("create empire: %w", err)
	}

	log.Printf("create: cluster: empire: created: id %d: no %d\n", empireId, empireNo)
	return nil
}

func runCreateClusterStarList(env *config.Environment, args []string) error {
	var arg string
	for len(args) > 0 && args[0] != "-- " {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			log.Printf("usage: empyr create star-lists [options]\n")
			log.Printf("  opt: --help          show help for the command   [false]\n")
			log.Printf("     : --debug=flag    enable various debug flags  [off]\n")
			log.Printf("     : --verbose       enhance logging             [false]\n")
			log.Printf("     : --path          path to database            [required]\n")
			return nil
		} else if flag, ok := argOptString(arg, "code"); ok {
			env.Game.Code = flag
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

	var err error
	if env.Store, err = store.Open(env.Database.Path, context.Background()); err != nil {
		return fmt.Errorf("%q: %w", env.Database.Path, err)
	}
	defer env.Store.Close()
	e, err := engine.Open(env.Store)
	if err != nil {
		return err
	}

	dataHtml, dataJson, err := engine.CreateClusterStarListCommand(e, &engine.CreateClusterStarListParams_t{Code: env.Game.Code})
	if err != nil {
		return err
	} else if err = os.WriteFile("cluster-star-list.html", dataHtml, 0644); err != nil {
		return err
	} else if err = os.WriteFile("cluster-star-list.json", dataJson, 0644); err != nil {
		return err
	}

	log.Printf("create: cluster-star-list: created %q\n", "cluster-star-list.html")
	log.Printf("create: cluster-star-list: created %q\n", "cluster-star-list.json")
	return nil
}

func runCreateClusterSystemMap(env *config.Environment, args []string) error {
	var arg string
	for len(args) > 0 && args[0] != "-- " {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			log.Printf("usage: empyr create cluster system-map [options]\n")
			log.Printf("  opt: --help          show help for the command   [false]\n")
			log.Printf("     : --debug=flag    enable various debug flags  [off]\n")
			log.Printf("     : --verbose       enhance logging             [false]\n")
			log.Printf("     : --path          path to database            [required]\n")
			return nil
		} else if flag, ok := argOptString(arg, "code"); ok {
			env.Game.Code = flag
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

	var err error
	if env.Store, err = store.Open(env.Database.Path, context.Background()); err != nil {
		return fmt.Errorf("%q: %w", env.Database.Path, err)
	}
	defer env.Store.Close()
	e, err := engine.Open(env.Store)
	if err != nil {
		return err
	}

	data, err := engine.CreateClusterMapCommand(e, &engine.CreateClusterMapParams_t{Code: env.Game.Code})
	if err != nil {
		return err
	} else if err = os.WriteFile("cluster-system-map.html", data, 0644); err != nil {
		return err
	}

	log.Printf("create: cluster: system-map: created %q\n", "cluster-system-map.html")
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
			log.Printf("     : --dry-run       show what would happen      [false]\n")
			return nil
		} else if flag, ok := argOptBool(arg, "dry-run"); ok {
			env.Database.DryRun = flag
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
	populateSystemDistances := false
	numberOfEmpires := int64(8)
	var arg string
	for len(args) > 0 && args[0] != "-- " {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			log.Printf("usage: empyr [options] create [options] game [options]\n")
			log.Printf("  opt: --help          show help for the command   [false]\n")
			log.Printf("     : --debug=flag    enable various debug flags  [off]\n")
			log.Printf("     : --verbose       enhance logging             [false]\n")
			log.Printf("     : --path          path to database            [required]\n")
			return nil
		} else if flag, ok := argOptBool(arg, "populate-system-distances"); ok {
			populateSystemDistances = flag
		} else if flag, ok := argOptInt(arg, "number-of-empires"); ok {
			numberOfEmpires = int64(flag)
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

	var err error
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
	if env.Game.Name == "" {
		return fmt.Errorf("missing game name")
	}

	if env.Store, err = store.Open(env.Database.Path, context.Background()); err != nil {
		return fmt.Errorf("%q: %w", env.Database.Path, err)
	}
	defer env.Store.Close()
	e, err := engine.Open(env.Store)
	if err != nil {
		return fmt.Errorf("%q: %w", env.Database.Path, err)
	}

	gameId, err := engine.CreateGameCommand(e, &engine.CreateGameParams_t{
		Code:                        env.Game.Code,
		Name:                        env.Game.Name,
		DisplayName:                 fmt.Sprintf("EC-%s", env.Game.Code),
		NumberOfEmpires:             numberOfEmpires,
		PopulateSystemDistanceTable: populateSystemDistances,
		Rand:                        rand.New(rand.NewPCG(0xdeadbeef, 0xcafedeed)),
	})
	if err != nil {
		return fmt.Errorf("game: %w", err)
	}

	log.Printf("create: game: created game %d\n", gameId)
	return nil
}

var (
	//go:embed templates/turn-report.gohtml
	turnReportTmpl string
)

func runCreateTurnReports(env *config.Environment, args []string) error {
	foundTurnId := false
	var arg string
	for len(args) > 0 && args[0] != "-- " {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			log.Printf("usage: empyr create turn-reports [options]\n")
			log.Printf("  opt: --help          show help for the command   [false]\n")
			log.Printf("     : --debug=flag    enable various debug flags  [off]\n")
			log.Printf("     : --verbose       enhance logging             [false]\n")
			log.Printf("     : --path          path to database            [required]\n")
			log.Printf("     : --turn=id       turn id to report on        [required]\n")
			return nil
		} else if flag, ok := argOptString(arg, "code"); ok {
			env.Game.Code = flag
		} else if flag, ok := argOptString(arg, "path"); ok {
			env.Database.Path = flag
		} else if no, ok := argOptInt(arg, "turn"); ok {
			if no < 0 {
				return fmt.Errorf("turn id must be >= 0")
			}
			env.Game.TurnNo, foundTurnId = no, true
		} else if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("unknown option: %q", arg)
		} else {
			return fmt.Errorf("unknown argument: %q", arg)
		}
	}
	if !foundTurnId {
		return fmt.Errorf("turn id required")
	}

	var err error
	if env.Store, err = store.Open(env.Database.Path, context.Background()); err != nil {
		return fmt.Errorf("%q: %w", env.Database.Path, err)
	}
	defer env.Store.Close()
	e, err := engine.Open(env.Store)
	if err != nil {
		return fmt.Errorf("%q: %w", env.Database.Path, err)
	}

	if data, err := engine.CreateTurnReportCommand(e, &engine.CreateTurnReportParams_t{
		Code:     env.Game.Code,
		TurnNo:   env.Game.TurnNo,
		EmpireNo: 1,
	}); err != nil {
		return fmt.Errorf("turn report: %w", err)
	} else {
		reportName := fmt.Sprintf("e%03d-turn-%04d-db.html", 1, env.Game.TurnNo)
		if err := os.WriteFile(reportName, data, 0644); err != nil {
			return err
		}
		log.Printf("created turn report empire %3d as %s\n", 1, reportName)
	}

	ts, err := template.New("turn-report").Parse(turnReportTmpl)
	if err != nil {
		return err
	}

	for empire := 1; empire < 256; empire++ {
		// attempt to read the empire's turn report
		filename := fmt.Sprintf("e%03d-turn-%04d.json", empire, env.Game.TurnNo)
		var payload struct {
			Site struct {
				CSS string
			}
			Report *turn_report_payload_t
		}
		payload.Site.CSS = "a02/css/monospace.css"
		if data, err := os.ReadFile(filename); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		} else {
			log.Printf("loaded report %s\n", filename)
			var input empire_turn_results_t
			if err := json.Unmarshal(data, &input); err != nil {
				return err
			}
			payload.Report = adapt_emp_to_turn_report_payload_t(&input)
			log.Printf("unmarshalled report %s\n", filename)
		}

		// buffer will hold the table containing the star lists
		buffer := &bytes.Buffer{}
		if err = ts.Execute(buffer, payload); err != nil {
			return err
		}
		// write the buffer to our output file
		reportName := fmt.Sprintf("e%03d-turn-%04d.html", empire, env.Game.TurnNo)
		if err := os.WriteFile(reportName, buffer.Bytes(), 0644); err != nil {
			return err
		}
		log.Printf("created turn report empire %3d as %s\n", empire, reportName)
	}

	log.Printf("%q: created turn reports\n", env.Database.Path)
	return nil
}

func runShow(env *config.Environment, args []string) error {
	var arg string
	for len(args) > 0 && args[0] != "-- " {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			log.Printf("usage: empyr show [command]\n")
			log.Printf("  cmd: env             dump environment\n")
			log.Printf("     : version         print application version\n")
			return nil
		} else if arg == "env" {
			if err := runShowEnvironment(env, args); err != nil {
				return fmt.Errorf("env: %w", err)
			}
			return nil
		} else if arg == "version" {
			if err := runShowVersion(env, args); err != nil {
				return fmt.Errorf("version: %w", err)
			}
			return nil
		} else if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("unknown option: %q", arg)
		} else if env.Database.Path == "" {
			env.Database.Path = arg
		} else {
			return fmt.Errorf("unknown argument: %q", arg)
		}
	}
	return nil
}

func runShowEnvironment(env *config.Environment, args []string) error {
	for _, v := range dotenv.EnvVariables("EMPYR_") {
		log.Printf("env: %-22s == %q\n", v.Key, v.Value)
	}
	log.Printf("env: env.database.path        %q\n", env.Database.Path)
	log.Printf("env: env.database.dryrun      %v\n", env.Database.DryRun)
	log.Printf("env: env.database.forcecreate %v\n", env.Database.ForceCreate)
	log.Printf("env: env.debug.database.open  %v\n", env.Debug.Database.Open)
	log.Printf("env: env.debug.database.close %v\n", env.Debug.Database.Close)
	log.Printf("env: env.game.code            %q\n", env.Game.Code)
	log.Printf("env: env.game.name            %q\n", env.Game.Name)
	log.Printf("env: env.game.turnno          %d\n", env.Game.TurnNo)
	log.Printf("env: env.verbose              %v\n", env.Verbose)
	return nil
}

func runShowVersion(env *config.Environment, args []string) error {
	log.Printf("version: %s\n", env.Version.String())
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
