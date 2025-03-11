// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	http.Server
	scheme string // must be http since we run behind a proxy
	host   string // should this be blank so that we're not bound to localhost?
	port   string
	ctx    context.Context
}

// Config defines the server configuration
type Config struct {
	Scheme         string        `json:"scheme,omitempty"`
	Host           string        `json:"host,omitempty"`
	Port           string        `json:"port,omitempty"`
	ReadTimeout    time.Duration `json:"read-timeout,omitempty"`
	WriteTimeout   time.Duration `json:"write-timeout,omitempty"`
	IdleTimeout    time.Duration `json:"idle-timeout,omitempty"`
	MaxHeaderBytes int           `json:"max-header-bytes,omitempty"`
}

// New returns an initialized server.
func New(cfg Config) *Server {
	return &Server{
		scheme: "http",
		host:   cfg.Host,
		port:   cfg.Port,
		ctx:    context.Background(),
		Server: http.Server{
			Addr:           net.JoinHostPort(cfg.Host, cfg.Port),
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			IdleTimeout:    cfg.IdleTimeout,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
		},
	}
}

// BaseURL returns the base URL for the server.
// I don't know why we need this.
func (s *Server) BaseURL() string {
	return fmt.Sprintf("%s://%s", s.scheme, s.Addr)
}

// Start the server and block until the server is stopped.
// implements a graceful shutdown of the server.
//
// TODO: the context should be used to cancel background tasks when shutting down the server.
func (s *Server) Start(routes http.Handler) error {
	s.Server.Handler = routes

	started := time.Now()

	// create a channel to listen for OS signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	// start the server in a goroutine so that the server doesn't block.
	// note that WE will block and wait for signals to stop the server.
	go func() {
		log.Printf("server: listening on %s\n", fmt.Sprintf("%s://%s", s.scheme, s.Addr))
		if err := http.ListenAndServe(s.Addr, s.Handler); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server: %v\n", err)
		}
		log.Printf("server: shutdown\n")
	}()

	// server is running; block until we receive a signal.
	sig := <-stop
	log.Printf("server: signal %v: received after %v\n", sig, time.Since(started))

	// graceful shutdown with a timeout.
	started = time.Now()
	timeout := time.Second * 5
	log.Printf("server: timeout %v: creating context (%v)\n", timeout, time.Since(started))
	ctxWithTimeout, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()

	// cancel any idle connections.
	log.Printf("server: canceling idle connections (%v)\n", time.Since(started))
	s.SetKeepAlivesEnabled(false)

	log.Printf("server: shutting down the server (%v)\n", time.Since(started))
	if err := s.Shutdown(ctxWithTimeout); err != nil {
		return fmt.Errorf("server: shutdown: %w", err)
	}

	log.Printf("server: ¡stopped gracefully! (%v)\n", time.Since(started))
	return nil
}
