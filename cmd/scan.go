// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cmd

import "github.com/spf13/cobra"

// cmdScan runs the scan command
var cmdScan = &cobra.Command{
	Use:   "scan",
	Short: "Scan things",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
