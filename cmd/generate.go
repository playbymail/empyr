// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cmd

import "github.com/spf13/cobra"

var argsGenerate struct {
	path string // path containing game folders
}

// cmdGenerate runs the generate command
var cmdGenerate = &cobra.Command{
	Use:   "generate",
	Short: "generate things",
}
