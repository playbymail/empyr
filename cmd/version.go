// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	cmdRoot.AddCommand(cmdVersion)
}

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the engine",
	Long:  `Print the version number of the engine.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.version.String())
	},
}
