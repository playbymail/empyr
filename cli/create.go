// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"fmt"
	"github.com/playbymail/empyr/pkg/clean"
	"github.com/playbymail/empyr/pkg/empyr"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

// this file implements the commands to create assets such as games and assets

// cmdCreate represents the base command when called without any subcommands
var cmdCreate = &cobra.Command{
	Use:   "create --path database",
	Short: "create games and assets",
	Long:  `create is the root of the generator commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()
		defer func() {
			log.Printf("create: elapsed time: %v\n", time.Now().Sub(started))
		}()
	},
}

// cmdCreateDatabase implements the create database command
var cmdCreateDatabase = &cobra.Command{
	Use:   "database --path database",
	Short: "create a new database",
	Long:  `Create a new database.`,
}

// cmdCreateGame implements the create game command
var cmdCreateGame = &cobra.Command{
	Use:   "game --path database --code code --name name --descr description",
	Short: "create a new game",
	Long:  `Create a new game (includes the database).`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if _, err := clean.IsValidCode(cmd.Flag("code").Value.String()); err != nil {
			return err
		} else if _, err = clean.IsValidName(cmd.Flag("name").Value.String()); err != nil {
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
		g, err := empyr.NewGame(code, name)
		if err != nil {
			log.Fatalf("create: game: %v", err)
		}
		// save the map as an HTML file
		if buffer, err := g.ClusterHTML(); err != nil {
			log.Fatalf("create: game: html: %v", err)
		} else if err = os.WriteFile("cluster-map.html", buffer.Bytes(), 0644); err != nil {
			log.Fatalf("create: game: html: %v", err)
		} else {
			log.Printf("create: game: html: %q: %d bytes\n", "cluster-map.html", buffer.Len())
		}
		log.Printf("create: game: created game %s (%q)\n", code, name)
	},
}

// isWeird checks if a string contains any special characters or escape sequences
// by comparing the raw string with its quoted representation. It's a clever way
// to detect special characters by leveraging Go's string quoting behavior.
//
// Returns true if the string contains special characters, false otherwise.
// Example: isweird("hello") returns false, isweird("hello\n") returns true
func isWeird(s string) bool {
	return `"`+s+`"` != fmt.Sprintf("%q", s)
}

func hasBadBytes(s string) bool {
	for _, ch := range []byte(s) {
		if !goodBytes[ch] {
			return true
		}
	}
	return false
}

var (
	goodCodeBytes [256]bool
	goodBytes     [256]bool
)

func init() {
	for _, ch := range []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") {
		goodCodeBytes[ch] = true
	}
	for _, ch := range []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-.") {
		goodBytes[ch] = true
	}
}
