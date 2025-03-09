// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import "github.com/spf13/cobra"

var cmdShow = &cobra.Command{
	Use:   "show",
	Short: "Show data",
}

var cmdShowEnv = &cobra.Command{
	Use:   "env",
	Short: "Show the environment",
	Run: func(cmd *cobra.Command, args []string) {
		const toLog = false
		dumpEnv(toLog)
	},
}
