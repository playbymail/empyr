// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/mdhender/semver"
	"github.com/playbymail/empyr/cli"
	"github.com/playbymail/empyr/pkg/dotenv"
	"log"
	"time"
)

var (
	version = semver.Version{Minor: 1}
)

func main() {
	//log.SetFlags(log.Lshortfile)
	log.SetFlags(log.LstdFlags | log.LUTC)

	started := time.Now()
	defer func() {
		log.Printf("elapsed time: %v\n", time.Now().Sub(started))
	}()

	if err := dotenv.Load("EMPYR"); err != nil {
		log.Fatalf("main: %+v\n", err)
	} else if command, err := cli.Initialize(
		cli.WithVersion(version),
	); err != nil {
		log.Fatalf("main: %+v\n", err)
	} else if err = command.Execute(); err != nil {
		log.Fatalf("\n%+v\n", err)
	}

	log.Printf("\n")
	log.Printf("completed in %v\n", time.Now().Sub(started))
}
