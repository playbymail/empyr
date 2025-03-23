// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"github.com/spf13/cobra"
)

// this file implements the commands to delete assets such as databases, games, and assets

// cmdDelete represents the base command when called without any subcommands
var cmdDelete = &cobra.Command{
	Use:   "delete",
	Short: "delete all the things",
	Long:  `delete is the root of the removal commands.`,
}
