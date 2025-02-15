// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package config provides the configuration for the application.
package config

import (
	"fmt"
	"github.com/mdhender/semver"
	"os"
	"strings"
)

type Environment struct {
	Database struct {
		Path        string
		ForceCreate bool
	}
	Debug struct {
		Database struct {
			Open  bool
			Close bool
		}
	}
	Game struct {
		Code string
		Name string
	}
	Verbose bool
	Version semver.Version
}

func Default(options ...Option) *Environment {
	cfg := &Environment{}

	// get default values from the environment
	if value, ok := os.LookupEnv("EMPYR_DATABASE_PATH"); ok {
		cfg.Database.Path = value
	}
	if value, ok := os.LookupEnv("EMPYR_GAME_CODE"); ok {
		cfg.Game.Code = value
	}
	if value, ok := os.LookupEnv("EMPYR_GAME_NAME"); ok {
		cfg.Game.Name = value
	}
	if value, ok := os.LookupEnv("EMPYR_VERBOSE"); ok {
		switch strings.ToLower(value) {
		case "1", "true", "yes":
			cfg.Verbose = true
		}
	}

	// apply options
	for _, option := range options {
		if err := option(cfg); err != nil {
			return nil
		}
	}

	// return the updated configuration
	return cfg
}

type Option func(*Environment) error

func WithDatabaseFlag(flag string) Option {
	return func(c *Environment) error {
		switch flag {
		case "database.force-create":
			c.Database.ForceCreate = true
		case "database.no-force-create":
			c.Database.ForceCreate = true
		default:
			return fmt.Errorf("database: %q: unknown flag", flag)
		}
		return nil
	}
}

func WithDatabasePath(path string) Option {
	return func(c *Environment) error {
		c.Database.Path = path
		return nil
	}
}

func WithDebugFlag(flag string) Option {
	return func(c *Environment) error {
		switch flag {
		case "database.open":
			c.Debug.Database.Open = true
		case "database.close":
			c.Debug.Database.Close = true
		default:
			return fmt.Errorf("debug: %q: unknown flag", flag)
		}
		return nil
	}
}

func WithVerbose(verbose bool) Option {
	return func(c *Environment) error {
		c.Verbose = verbose
		return nil
	}
}

// WithVerboseFlag sets the verbose variable from the given flag.
// If the flag (converted to lowercase) is "1", "true", or "yes", then
// verbose is set to true. Otherwise, verbose is set to false.
func WithVerboseFlag(flag string) Option {
	return func(c *Environment) error {
		switch strings.ToLower(flag) {
		case "1", "true", "yes":
			c.Verbose = true
		}
		return nil
	}
}

func WithVersion(version semver.Version) Option {
	return func(c *Environment) error {
		c.Version = version
		return nil
	}
}
