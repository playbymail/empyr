// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/spf13/cobra"
	"log"
	"time"
)

// cmdRoot represents the base command when called without any subcommands
var cmdRoot = &cobra.Command{
	Short: "empyr: a game engine",
	Long:  `empyr is an engine inspired by better games.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("cobra: running root.PersistentPreRunE")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		started := time.Now()
		defer func() {
			log.Printf("root: elapsed time: %v\n", time.Now().Sub(started))
		}()
		return nil
	},
}
