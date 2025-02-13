// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"crypto/md5"
	"encoding/binary"
	"github.com/playbymail/empyr/cmd/server/pkg/users"
	"io"
	"log"
	"net"
	"net/http"
)

// server defines the server
type server struct {
	http.Server
	e     *Engine
	users *users.Users
	games map[string]*GameMeta
	seqno int
	salt  string
}

// serverContextKey is the context key type for storing parameters in context.Context.
type serverContextKey string

// newServer returns an initialized server.
// the main change from the default server is that we override the default timeouts.
// see the following sources for an explanation of why:
//
//	https://blog.cloudflare.com/exposing-go-on-the-internet/
//	https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
//	https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
func newServer(cfg *config, options ...func(*server) error) (s *server, err error) {
	s = &server{}
	s.Addr = net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)
	s.IdleTimeout = cfg.Server.Timeout.Idle
	s.ReadTimeout = cfg.Server.Timeout.Read
	s.WriteTimeout = cfg.Server.Timeout.Write
	s.MaxHeaderBytes = 1 << 20

	s.games = make(map[string]*GameMeta)
	s.users = users.New()

	// allow caller to override the default values
	for _, option := range options {
		if err := option(s); err != nil {
			return nil, err
		}
	}

	if s.e == nil {
		if s.e, err = NewEngine(); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func addEngine(e *Engine) func(*server) error {
	return func(s *server) error {
		s.e = e
		return nil
	}
}

func addGame(id, seedString, name string, userNames ...string) func(*server) error {
	return func(s *server) error {
		meta, err := s.createGame(id, name, seedString)
		if err != nil {
			return err
		}
		for _, name := range userNames {
			if matches := s.users.Filter(func(user *users.User) bool { return user.Name == name }); len(matches) == 0 {
				return ErrNoData
			} else if err := meta.addPlayer(matches[0], matches[0].Name); err != nil {
				return err
			} else if err := meta.Game.AddPlayer(matches[0].Name, matches[0].Name); err != nil {
				return err
			}
		}
		log.Printf("[server] sample.game %v\n", *meta)
		return nil
	}
}

func addUser(id, name, email string) func(*server) error {
	return func(s *server) error {
		user, err := s.createUser(id, name, email)
		if err != nil {
			return err
		}
		log.Printf("[server] sample.user %v\n", *user)
		return nil
	}
}

func setSalt(salt string) func(*server) error {
	return func(s *server) error {
		s.salt = salt
		return nil
	}
}

func (s *server) nextval() int {
	s.seqno++
	return s.seqno
}

func (s *server) seed(seedString string) int64 {
	switch seedString {
	case "1812":
		return 1812
	case "1917":
		return 1917
	}
	hasher := md5.New()
	io.WriteString(hasher, s.salt)
	io.WriteString(hasher, seedString)
	return int64(binary.BigEndian.Uint64(hasher.Sum(nil)))
}
