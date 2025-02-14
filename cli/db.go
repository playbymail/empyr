// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"context"
	"github.com/playbymail/empyr/pkg/store"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// this file implements the commands to manage the database

// cmdDB represents the base command when called without any subcommands
var cmdDB = &cobra.Command{
	Use:   "db --path database",
	Short: "database management",
	Long:  `db is the root of the database management commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("db: elapsed time: %v\n", time.Now().Sub(started))
		}()
	},
}

// cmdDBCreate creates and initializes a new database.
// It returns an error if the database already exists.
var cmdDBCreate = &cobra.Command{
	Use:   "create --path database",
	Short: "initialize a new database at the given path",
	Long:  `Create and initialize a new database. It is an error if the path is not absolute or the file already exists.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("db: create: elapsed time: %v\n", time.Now().Sub(started))
		}()
		path := cmd.Flag("path").Value.String()
		log.Printf("db: create: %s\n", path)
		if err := store.Create(path); err != nil {
			log.Fatal(err)
		}
		log.Printf("db: create: finished\n")
	},
}

// cmdDBOpens an existing database.
// It returns an error if it fails to open the database.
var cmdDBOpen = &cobra.Command{
	Use:   "open --path database",
	Short: "open an existing database",
	Long:  `Open a database. Useful for. testing.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("db: open: elapsed time: %v\n", time.Now().Sub(started))
		}()
		path := cmd.Flag("path").Value.String()
		log.Printf("db: open: %s\n", path)
		s, err := store.Open(path, context.Background())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("db: open: %s\n", path)
		_ = s.Close()
		log.Printf("db: open: finished\n")
	},
}
