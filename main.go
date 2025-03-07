// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mdhender/semver"
	"github.com/playbymail/empyr/config"
	"github.com/playbymail/empyr/engine"
	"github.com/playbymail/empyr/pkg/dotenv"
	"github.com/playbymail/empyr/pkg/stdlib"
	"github.com/playbymail/empyr/store"
	"html/template"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
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
		} else if arg == "star-lists" {
			if err := runCreateStarLists(env, args); err != nil {
				return fmt.Errorf("star-lists: %w", err)
			}
			return nil
		} else if arg == "cluster-map" {
			if err := runCreateClusterMap(env, args); err != nil {
				return fmt.Errorf("cluster-map: %w", err)
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
	log.Printf("     : star-lists      create list of systems in game\n")
	log.Printf("     : cluster-map     create cluster map of systems\n")
	log.Printf("     : turn-reports    create turn reports for all empires\n")
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
			log.Printf("usage: empyr [options] create [options] database [options] path_to_database\n")
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
			log.Printf("usage: empyr [options] create [options] game [options]\n")
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
	for _, obj := range []struct {
		name string
		ptr  any
	}{
		{name: "systems", ptr: gc.Systems[1:]},
		{name: "stars", ptr: gc.Stars[1:]},
		{name: "orbits", ptr: gc.Orbits[1:]},
		{name: "planets", ptr: gc.Planets[1:]},
	} {
		if data, err := json.MarshalIndent(obj.ptr, "", "  "); err != nil {
			return err
		} else if err = os.WriteFile(obj.name+".json", data, 0644); err != nil {
			log.Fatalf("cluster: %s: %v\n", obj.name, err)
		} else {
			log.Printf("cluster: %s: wrote %s\n", obj.name, obj.name+".json")
		}
	}

	//log.Printf("%q: creating cluster %p\n", env.Database.Path, gc)
	//g, err := empyr.NewGame(env.Game.Code, env.Game.Name)
	//if err != nil {
	//	return fmt.Errorf("code: %q: %w", err)
	//}
	//if data, err := json.MarshalIndent(g, "", "  "); err != nil {
	//	return err
	//} else {
	//	log.Printf("game: %s\n", string(data))
	//}

	log.Printf("%q: created game\n", env.Database.Path)
	return nil
}

var (
	//go:embed templates/cluster-map.gohtml
	clusterMapTmpl string
	//go:embed templates/turn-report.gohtml
	turnReportTmpl string
)

func runCreateClusterMap(env *config.Environment, args []string) error {
	var arg string
	for len(args) > 0 && args[0] != "-- " {
		arg, args = args[0], args[1:]
		if argOptHelp(arg) {
			log.Printf("usage: empyr create cluster-map [options]\n")
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

	// load the systems data from JSON files
	// TODO: this should be done in the engine and from the database, not flat files
	var systems []*engine.System_t
	if data, err := os.ReadFile("systems.json"); err != nil {
		return err
	} else if err = json.Unmarshal(data, &systems); err != nil {
		return err
	}

	// buffer will hold the table containing the star lists
	buffer := &bytes.Buffer{}

	ts, err := template.New("cluster-map").Parse(clusterMapTmpl)
	if err != nil {
		return err
	}
	type cluster_t struct {
		Id      int
		X, Y, Z int
		Size    float64
		// Black, Blue, Gray, Green, Magenta, Purple, Random, Red, Teal, White, Yellow
		Color template.JS
	}
	var data []cluster_t
	for _, s := range systems {
		var color template.JS
		switch len(s.Stars) {
		case 4:
			color = "Green"
		case 3:
			color = "Gray"
		case 2:
			color = "Blue"
		case 1:
			color = "Black"
		default:
			panic(fmt.Sprintf("assert(len(s.Stars) != %d)", len(s.Stars)))
		}
		data = append(data, cluster_t{
			Id:    s.Id,
			X:     s.Coordinates.X - 15,
			Y:     s.Coordinates.Y - 15,
			Z:     s.Coordinates.Z - 15,
			Size:  0.333333,
			Color: color,
		})
	}
	if err = ts.Execute(buffer, data); err != nil {
		return err
	}
	// write the buffer to our output file
	if err := os.WriteFile("cluster-map.html", buffer.Bytes(), 0644); err != nil {
		return err
	}

	log.Printf("%q: created cluster-map: %3d systems\n", env.Database.Path, len(systems))
	return nil
}

func runCreateStarLists(env *config.Environment, args []string) error {
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

	// load the systems data from JSON files
	// TODO: this should be done in the engine and from the database, not flat files
	var systems []*engine.System_t
	if data, err := os.ReadFile("systems.json"); err != nil {
		return err
	} else if err = json.Unmarshal(data, &systems); err != nil {
		return err
	}

	// buffer will hold the table containing the star lists
	buffer := &bytes.Buffer{}

	// create a new tabwriter that writes to our buffer
	w := tabwriter.NewWriter(buffer, 0, 0, 2, ' ', tabwriter.Debug)

	// write table header
	_, _ = fmt.Fprintln(w, "System ID\tCoordinates\tNumber of Stars\tDistance From Center")

	// write system data rows
	for i, system := range systems {
		// Format coordinates as (x, y, z)
		coords := fmt.Sprintf("(%02d, %02d, %02d)", system.Coordinates.X, system.Coordinates.Y, system.Coordinates.Z)
		// distance from center uses 15,15,15 as the center
		dx := float64(system.Coordinates.X - 15)
		dy := float64(system.Coordinates.Y - 15)
		dz := float64(system.Coordinates.Z - 15)
		distance := int(math.Ceil(math.Sqrt(dx*dx + dy*dy + dz*dz)))

		numberOfStars := len(system.Stars)

		// Write the row
		_, _ = fmt.Fprintf(w, "%d\t%s\t%d\t%d\n", i+1, coords, numberOfStars, distance)
	}

	// flush ensures all data is written to the buffer
	_ = w.Flush()

	// write the buffer to our output file
	if err := os.WriteFile("star-lists.txt", buffer.Bytes(), 0644); err != nil {
		return err
	}

	log.Printf("%q: created star-lists: %3d systems\n", env.Database.Path, len(systems))
	return nil
}

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
		} else if flag, ok := argOptString(arg, "name"); ok {
			env.Game.Name = flag
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

func argOptInt(arg, flag string) (value int, ok bool) {
	flag = "--" + flag + "="
	if strings.HasPrefix(arg, flag) {
		if value, err := strconv.Atoi(strings.TrimPrefix(arg, flag)); err == nil {
			return value, true
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
