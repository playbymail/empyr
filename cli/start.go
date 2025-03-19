// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cli

import (
	"context"
	"github.com/playbymail/empyr/app"
	"github.com/playbymail/empyr/internal/jot"
	"github.com/playbymail/empyr/internal/qxb"
	"github.com/playbymail/empyr/internal/services/auth"
	"github.com/playbymail/empyr/internal/services/games"
	"github.com/playbymail/empyr/internal/services/sessions"
	"github.com/playbymail/empyr/internal/services/smgr"
	"github.com/playbymail/empyr/server"
	"github.com/playbymail/empyr/store"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var cmdStart = &cobra.Command{
	Use:   "start",
	Short: "start application components",
}

var cmdStartServer = &cobra.Command{
	Use:   "server",
	Short: "start the web server",
	Long:  `Start the web server`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()

		// open the database. this will create the database if it doesn't exist.
		repo, err := store.Open(flags.Database.Path, context.Background())
		if err != nil {
			log.Fatalf("error: store.open: %v\n", err)
		}
		// keep the repository open until the application exits.
		defer func() {
			_ = repo.Close()
		}()

		mgr := smgr.NewSessionManager(
			smgr.NewInMemorySessionStore(),
			15*time.Minute,
			12*time.Hour,
			24*time.Hour,
			flags.Sessions.Domain,
			"empyr-sm")

		q, err := qxb.NewQXB(flags.Application.Assets.Public, mgr)
		if err != nil {
			log.Fatalf("error: qxb.new: %v\n", err)
		}

		// make a new jot factory
		jf := jot.NewFactory(flags.Sessions.Secret)

		authService := auth.NewService(repo)
		gamesService := games.NewService(repo)
		sessionsService := sessions.NewService(repo, jf)

		// initialize the web application
		a, err := app.New(repo, jf, flags.Application.Assets.Public, flags.Application.Assets.Templates, context.Background())
		if err != nil {
			_ = repo.Close()
			log.Fatalf("error: app.new: %v\n", err)
		}

		// configure the server
		cfg := server.Config{
			Scheme:         flags.Server.Scheme,
			Host:           flags.Server.Host,
			Port:           flags.Server.Port,
			ReadTimeout:    flags.Server.ReadTimeout,
			WriteTimeout:   flags.Server.WriteTimeout,
			IdleTimeout:    flags.Server.IdleTimeout,
			MaxHeaderBytes: int(flags.Server.MaxHeaderBytes),
		}

		srv := server.New(cfg)

		// start the server. blocks until the server receives a signal to stop.
		err = srv.Start(q.Handle(a.Router(authService, gamesService, sessionsService)))

		// force the repository to close before we exit.
		_ = repo.Close()

		log.Printf("server: shut down after %v\n", time.Since(started))
		if err != nil {
			log.Fatalf("error: server.start: %v\n", err)
		}
	},
}
