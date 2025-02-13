// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements the command line interface for empyr.
package main

import (
	"github.com/mdhender/semver"
	"log"
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

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cmdRoot.AddCommand(cmdDB)

	cmdDB.PersistentFlags().String("path", "", "path to the database")
	if err := cmdDB.MarkPersistentFlagRequired("path"); err != nil {
		return err
	}
	cmdDB.AddCommand(cmdDBCreate)
	cmdDB.AddCommand(cmdDBOpen)

	cmdRoot.AddCommand(cmdCreate)
	cmdCreate.PersistentFlags().String("path", "", "path to the database")
	if err := cmdCreate.MarkPersistentFlagRequired("path"); err != nil {
		return err
	}
	cmdCreate.AddCommand(cmdCreateGame)
	cmdCreateGame.Flags().String("code", "", "code for the game")
	if err := cmdCreateGame.MarkFlagRequired("code"); err != nil {
		return err
	}
	cmdCreateGame.Flags().String("name", "", "name of the game")
	if err := cmdCreateGame.MarkFlagRequired("name"); err != nil {
		return err
	}

	cmdRoot.AddCommand(cmdVersion)

	return cmdRoot.Execute()
}
