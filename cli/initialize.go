// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package cli implements the command line interface for empyr.
package cli

import (
	"github.com/mdhender/semver"
	"github.com/spf13/cobra"
)

func Initialize(options ...Option) (*cobra.Command, error) {
	for _, option := range options {
		if err := option(); err != nil {
			return nil, err
		}
	}

	cmdRoot.AddCommand(cmdDB)

	cmdDB.PersistentFlags().String("path", "", "path to the database")
	if err := cmdDB.MarkPersistentFlagRequired("path"); err != nil {
		return nil, err
	}
	cmdDB.AddCommand(cmdDBCreate)
	cmdDB.AddCommand(cmdDBOpen)

	cmdRoot.AddCommand(cmdCreate)
	cmdCreate.PersistentFlags().String("path", "", "path to the database")
	if err := cmdCreate.MarkPersistentFlagRequired("path"); err != nil {
		return nil, err
	}
	cmdCreate.AddCommand(cmdCreateDatabase)
	cmdCreate.AddCommand(cmdCreateGame)
	cmdCreateGame.Flags().String("code", "", "code for the game")
	if err := cmdCreateGame.MarkFlagRequired("code"); err != nil {
		return nil, err
	}
	cmdCreateGame.Flags().String("name", "", "name of the game")
	if err := cmdCreateGame.MarkFlagRequired("name"); err != nil {
		return nil, err
	}

	cmdRoot.AddCommand(cmdVersion)

	return cmdRoot, nil
}

type Option func() error

func WithVersion(version semver.Version) Option {
	return func() error {
		argVersion.version = version
		return nil
	}
}
