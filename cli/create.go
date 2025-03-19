// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"context"
	"errors"
	"fmt"
	"github.com/mdhender/phrases/v2"
	"github.com/playbymail/empyr/engine"
	"github.com/playbymail/empyr/internal/domains"
	"github.com/playbymail/empyr/pkg/stdlib"
	"github.com/playbymail/empyr/store"
	"github.com/spf13/cobra"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// this file implements the commands to create assets such as databases, games, and assets

// cmdCreate represents the base command when called without any subcommands
var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "create all the things",
	Long:  `create is the root of the generator commands.`,
}

// cmdCreateDatabase implements the create database command
var cmdCreateDatabase = &cobra.Command{
	Use:   "database",
	Short: "create a new database",
	Long:  `Create a new database.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		log.Printf("create: database: %q\n", flags.Database.Path)
		if stdlib.IsExists(flags.Database.Path) {
			if !flags.Database.ForceCreate {
				log.Fatalf("error: %v\n", ErrFileExists)
			}
			log.Printf("create: database: deleting existing database\n")
			if err := stdlib.Remove(flags.Database.Path); err != nil {
				log.Fatalf("error: stdlib.remove: %v\n", errors.Join(ErrDeleteFailed, err))
			}
		}
		log.Printf("create: database: %q\n", flags.Database.Path)
		if err := store.Create(flags.Database.Path); err != nil {
			log.Fatalf("error: store.create: %v\n", err)
		}
		log.Printf("create: database: completed in %v\n", time.Since(started))
	},
}

// cmdCreateEmpire creates a new empire
var cmdCreateEmpire = &cobra.Command{
	Use:   "empire --user --empire",
	Short: "create a new empire",
	Long:  `Create a new empire.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("create: empire: elapsed time: %v\n", time.Now().Sub(started))
		}()

		userHandle := cmd.Flag("user").Value.String()
		if userHandle == "" {
			log.Fatalf("create: empire: user: handle is required\n")
		} else if _, err := engine.IsValidHandle(userHandle); err != nil {
			log.Fatalf("create: empire: user: %v\n", err)
		}

		repo, err := store.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}

		empireId, empireNo, err := engine.CreateEmpireCommand(e, &engine.CreateEmpireParams_t{
			Code:     flags.Game.Code,
			Username: userHandle,
		})
		if err != nil {
			log.Fatalf("error: engine.CreateEmpireCommand: %v\n", err)
		}

		log.Printf("create: empire: created %d (%d)\n", empireId, empireNo)
	},
}

// cmdCreateGame implements the create game command
var cmdCreateGame = &cobra.Command{
	Use:   "game --code code --name name --descr description",
	Short: "create a new game",
	Long:  `Create a new game.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if _, err := engine.IsValidCode(cmd.Flag("code").Value.String()); err != nil {
			return err
		} else if _, err = engine.IsValidName(cmd.Flag("name").Value.String()); err != nil {
			return err
		} else if _, err = engine.IsValidDescription(cmd.Flag("descr").Value.String()); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("create: game: elapsed time: %v\n", time.Now().Sub(started))
		}()
		code := cmd.Flag("code").Value.String()
		name := cmd.Flag("name").Value.String()
		descr := cmd.Flag("descr").Value.String()
		if descr == "" {
			descr = fmt.Sprintf("A game of %s", name)
		}
		log.Printf("create: game: code  %q\n", code)
		log.Printf("create: game: name  %q\n", name)
		log.Printf("create: game: descr %q\n", descr)

		repo, err := store.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}

		gameId, err := engine.CreateGameCommand(e, &engine.CreateGameParams_t{
			Code:        flags.Game.Code,
			Name:        flags.Game.Name,
			DisplayName: fmt.Sprintf("EC-%s", flags.Game.Code),
			Rand:        rand.New(rand.NewPCG(0xdeadbeef, 0xcafedeed)),
			ForceCreate: flags.Game.ForceCreate,
		})
		if err != nil {
			log.Fatalf("error: engine.CreateGameCommand: %v\n", err)
		}

		log.Printf("create: game: created game %d\n", gameId)
	},
}

// cmdCreateStarList creates a new star list
var cmdCreateStarList = &cobra.Command{
	Use:   "star-list",
	Short: "create a star list for the cluster",
	Long:  `Create a star list for the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("create: star-list: elapsed time: %v\n", time.Now().Sub(started))
		}()

		repo, err := store.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}

		dataHtml, dataJson, err := engine.CreateClusterStarListCommand(e, &engine.CreateClusterStarListParams_t{Code: flags.Game.Code})
		if err != nil {
			log.Fatalf("error: engine.CreateClusterStarListCommand: %v\n", err)
		} else if err = os.WriteFile("cluster-star-list.html", dataHtml, 0644); err != nil {
			log.Fatalf("error: os.WriteFile: %v\n", err)
		} else if err = os.WriteFile("cluster-star-list.json", dataJson, 0644); err != nil {
			log.Fatalf("error: os.WriteFile: %v\n", err)
		}

		log.Printf("create: cluster-star-list: created %q\n", "cluster-star-list.html")
		log.Printf("create: cluster-star-list: created %q\n", "cluster-star-list.json")
	},
}

// cmdCreateSystemMap creates a new system map
var cmdCreateSystemMap = &cobra.Command{
	Use:   "system-map",
	Short: "create a system map for the cluster",
	Long:  `Create a system map for the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("create: system-map: elapsed time: %v\n", time.Now().Sub(started))
		}()

		repo, err := store.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}

		data, err := engine.CreateClusterMapCommand(e, &engine.CreateClusterMapParams_t{Code: flags.Game.Code})
		if err != nil {
			log.Fatalf("error: engine.CreateClusterMapCommand: %v\n", err)
		} else if err = os.WriteFile("cluster-system-map.html", data, 0644); err != nil {
			log.Fatalf("error: os.WriteFile: %v\n", err)
		}

		log.Printf("create: system-map: created %q\n", "cluster-system-map.html")
	},
}

// cmdCreateSystemSurveyReport creates system survey reports for one empire in a game
var cmdCreateSystemSurveyReport = &cobra.Command{
	Use:   "system-survey-report --empire-no --turn-no",
	Short: "create system survey reports for one empire in a game",
	Long:  `Create system survey reports for one empire in a game.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("create: system-survey-report: elapsed time: %v\n", time.Now().Sub(started))
		}()

		repo, err := store.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}

		var empireNo int64
		if n, err := strconv.Atoi(cmd.Flag("empire-no").Value.String()); err != nil || n < 1 || n > 256 {
			log.Fatalf("error: empire-no must be between 1 and 255")
		} else {
			empireNo = int64(n)
		}
		if flags.Game.TurnNo < 0 || flags.Game.TurnNo > 9999 {
			log.Fatalf("error: turn-no must be between 0 and 9999")
		}

		empireSurveysPath := filepath.Join(flags.Surveys.Path, fmt.Sprintf("e%03d", empireNo), "surveys")
		reportName := filepath.Join(empireSurveysPath, fmt.Sprintf("e%03d-turn-%04d.html", empireNo, flags.Game.TurnNo))

		data, err := engine.CreateSystemSurveyReportCommand(e, &engine.CreateSystemSurveyReportParams_t{Code: flags.Game.Code, TurnNo: flags.Game.TurnNo, EmpireNo: empireNo})
		if err != nil {
			log.Fatalf("error: system survey report: %v\n", err)
		}

		if err := os.WriteFile(reportName, data, 0644); err != nil {
			log.Fatalf("error: os.WriteFile: %v\n", err)
		}
		log.Printf("created system survey report empire %3d as %s\n", empireNo, reportName)

		log.Printf("create: system-survey-report: completed\n")
	},
}

// cmdCreateSystemSurveyReports creates turn reports for all the empires in a game
var cmdCreateSystemSurveyReports = &cobra.Command{
	Use:   "system-survey-reports --turn-no",
	Short: "create turn reports for all empires in a game",
	Long:  `Create turn reports for all empires in a game.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("create: system-survey-reports: elapsed time: %v\n", time.Now().Sub(started))
		}()

		repo, err := store.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}

		err = engine.CreateSystemSurveyReportsCommand(e, &engine.CreateSystemSurveyReportsParams_t{
			Code:   flags.Game.Code,
			TurnNo: flags.Game.TurnNo,
			Path:   flags.Surveys.Path,
		})
		if err != nil {
			log.Fatalf("error: engine.CreateSystemSurveyReportsCommand: %v\n", err)
		}

		log.Printf("create: system-survey-reports: completed\n")
	},
}

// cmdCreateTurnReport creates turn report for one empire in a game
var cmdCreateTurnReport = &cobra.Command{
	Use:   "turn-report --empire-no --turn-no",
	Short: "create turn reports for one empire in a game",
	Long:  `Create turn reports for one empire in a game.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("create: turn-report: elapsed time: %v\n", time.Now().Sub(started))
		}()

		repo, err := store.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}

		var empireNo int64
		if n, err := strconv.Atoi(cmd.Flag("empire-no").Value.String()); err != nil || n < 1 || n > 256 {
			log.Fatalf("error: empire-no must be between 1 and 255")
		} else {
			empireNo = int64(n)
		}
		if flags.Game.TurnNo < 0 || flags.Game.TurnNo > 9999 {
			log.Fatalf("error: turn-no must be between 0 and 9999")
		}

		empireReportsPath := filepath.Join(flags.Reports.Path, fmt.Sprintf("e%03d", empireNo), "reports")
		reportName := filepath.Join(empireReportsPath, fmt.Sprintf("e%03d-turn-%04d.html", empireNo, flags.Game.TurnNo))

		data, err := engine.CreateTurnReportCommand(e, &engine.CreateTurnReportParams_t{Code: flags.Game.Code, TurnNo: flags.Game.TurnNo, EmpireNo: empireNo})
		if err != nil {
			log.Fatalf("error: turn report: %v\n", err)
		}

		if err := os.WriteFile(reportName, data, 0644); err != nil {
			log.Fatalf("error: os.WriteFile: %v\n", err)
		}
		log.Printf("created turn report empire %3d as %s\n", empireNo, reportName)

		log.Printf("create: turn-report: completed\n")
	},
}

// cmdCreateTurnReports creates turn reports for all the empires in a game
var cmdCreateTurnReports = &cobra.Command{
	Use:   "turn-reports --turn-no",
	Short: "create turn reports for all empires in a game",
	Long:  `Create turn reports for all empires in a game.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("create: turn-reports: elapsed time: %v\n", time.Now().Sub(started))
		}()

		repo, err := store.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}

		err = engine.CreateTurnReportsCommand(e, &engine.CreateTurnReportsParams_t{
			Code:   flags.Game.Code,
			TurnNo: flags.Game.TurnNo,
			Path:   flags.Reports.Path,
		})
		if err != nil {
			log.Fatalf("error: engine.CreateTurnReportsCommand: %v\n", err)
		}

		log.Printf("create: turn-reports: completed\n")
	},
}

// cmdCreateUser implements the create user command
var cmdCreateUser = &cobra.Command{
	Use:   "user --username username --password password --is-admin",
	Short: "create a new user",
	Long:  `Create a new user.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("create: user: elapsed time: %v\n", time.Now().Sub(started))
		}()

		repo, err := store.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		defer repo.Close()

		isAdmin := cmd.Flag("is-admin").Value.String() == "true"
		username := cmd.Flag("user").Value.String()
		email := cmd.Flag("email").Value.String()

		password := cmd.Flag("password").Value.String()
		if password == "" {
			password = phrases.Generate(4)
		}

		log.Printf("create: user %q: email %q: is-admin %v\n", username, email, isAdmin)

		var userID domains.UserID
		if isAdmin {
			userID, err = repo.RegisterAdminUser(username, email, password)
		} else {
			userID, err = repo.RegisterUser(username, email, password)
		}
		if err != nil {
			log.Fatalf("error: repo.RegisterUser: %v\n", err)
		}
		log.Printf("create: user %q: password %q: id %d\n", username, password, userID)

		log.Printf("create: user: completed\n")
	},
}
