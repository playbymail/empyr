// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements a server for the Empyrean Challenge game.
package main

import (
	"log"
	"os"
)

var (
	versionMajor string = "0"
	versionMinor string = "0"
	versionPatch string = "0"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Printf("%+v\n", err)
		os.Exit(2)
	}

	if err := run(cfg); err != nil {
		log.Printf("%+v\n", err)
		os.Exit(2)
	}
}
