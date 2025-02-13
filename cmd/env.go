// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdEnv = &cobra.Command{
	Use:   "env",
	Short: "dump the environment",
	Long:  `Print out the environment settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%-30s == %q\n", "version", cmdRoot.Version)
		fmt.Printf("%-30s == %q\n", "homeFolder", config.homeFolder)
		fmt.Printf("%-30s == %q\n", "configFile", viper.ConfigFileUsed())
	},
}
