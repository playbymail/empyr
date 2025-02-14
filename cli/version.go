// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"github.com/mdhender/semver"
	"github.com/spf13/cobra"
	"log"
)

var (
	argVersion struct {
		version semver.Version
	}

	cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of the engine",
		Long:  `Print the version number of the engine.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("empyr: version %s\n", argVersion.version.String())
		},
	}
)
