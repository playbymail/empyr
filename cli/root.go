// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"github.com/spf13/cobra"
)

// cmdRoot represents the base command when called without any subcommands
var cmdRoot = &cobra.Command{
	Short: "empyr: a game engine",
	Long:  `empyr is an engine inspired by better games.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		const toLog = true
		if env.Debug.DumpEnv {
			dumpEnv(toLog)
		}
	},
}
