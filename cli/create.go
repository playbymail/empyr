// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"context"
	"errors"
	"fmt"
	"github.com/playbymail/empyr/engine"
	"github.com/playbymail/empyr/pkg/stdlib"
	"github.com/playbymail/empyr/repos"
	"github.com/playbymail/empyr/repos/empires"
	"github.com/spf13/cobra"
	"log"
	"math/rand/v2"
	"os"
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
		if err := repos.Create(flags.Database.Path); err != nil {
			log.Fatalf("error: repos.create: %v\n", err)
		}

		code := cmd.Flag("code").Value.String()
		log.Printf("create: game: code  %q\n", code)

		name := cmd.Flag("name").Value.String()
		log.Printf("create: game: name  %q\n", name)

		descr := cmd.Flag("descr").Value.String()
		if descr == "" {
			descr = fmt.Sprintf("A game of %s", name)
		}
		log.Printf("create: game: descr %q\n", descr)

		repo, err := repos.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: repos.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}

		err = engine.CreateGameCommand(e, &engine.CreateGameParams_t{
			Code:        flags.Game.Code,
			Name:        flags.Game.Name,
			DisplayName: fmt.Sprintf("EC-%s", flags.Game.Code),
			Rand:        rand.New(rand.NewPCG(0xdeadbeef, 0xcafedeed)),
			ForceCreate: flags.Game.ForceCreate,
		})
		if err != nil {
			log.Fatalf("error: engine.CreateGameCommand: %v\n", err)
		}

		log.Printf("create: game: created game %q\n", flags.Game.Code)

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

		repo, err := repos.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: repos.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}
		empireRepo := empires.NewRepo(repo)

		var empireID int64
		if n, err := strconv.Atoi(cmd.Flag("id").Value.String()); err != nil {
			log.Fatalf("create: empire: %v\n", err)
		} else if n < 1 || n > 250 {
			log.Fatalf("create: empire: id must be between 1 and 250\n")
		} else {
			empireID = int64(n)
		}
		empireID, err = empireRepo.CreateEmpireWithID(empireID)
		if err != nil {
			log.Fatalf("error: createEmpire: %v\n", err)
		}
		// update the username and email if we have them
		// create the default colony
		// some players may change their mining orders
		// some players may retool their factories
		empireID, err = engine.CreateEmpireCommand(e, &engine.CreateEmpireParams_t{
			EmpireID: int64(empireID),
			Username: userHandle,
		})
		if err != nil {
			log.Fatalf("error: engine.CreateEmpireCommand: %v\n", err)
		}

		log.Printf("create: empire: created %d\n", empireID)
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

		repo, err := repos.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: repos.open: %v\n", err)
		}
		defer repo.Close()
		e, err := engine.Open(repo)
		if err != nil {
			log.Fatalf("error: engine.open: %v\n", err)
		}

		err = engine.CreateGameCommand(e, &engine.CreateGameParams_t{
			Code:        flags.Game.Code,
			Name:        flags.Game.Name,
			DisplayName: fmt.Sprintf("EC-%s", flags.Game.Code),
			Rand:        rand.New(rand.NewPCG(0xdeadbeef, 0xcafedeed)),
			ForceCreate: flags.Game.ForceCreate,
		})
		if err != nil {
			log.Fatalf("error: engine.CreateGameCommand: %v\n", err)
		}

		log.Printf("create: game: created game %s\n", flags.Game.Code)
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

		repo, err := repos.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: repos.open: %v\n", err)
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

		repo, err := repos.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: repos.open: %v\n", err)
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
