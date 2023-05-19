// empyr - a reimagining of Vern Holford's Empyrean Challenge
// Copyright (C) 2023 Michael D Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

package main

import (
	"log"
	"net/http"

	"github.com/playbymail/empyr/cmd/server/pkg/http/handlers"
)

func run(cfg *config) error {
	var options []func(*server) error
	options = append(options, setSalt(cfg.Server.Salt))

	e, err := NewEngine()
	if err != nil {
		return err
	}
	options = append(options, addEngine(e))

	// inject some users
	options = append(options, addUser("bf4c8168-6aab-409d-80cf-a4ee901904ef", "mdhender", "mdhender@example.com"))
	options = append(options, addUser("236bb1a5-1ae8-411a-a71f-791f4f03aa99", "yojimbo", "yojimbo@example.com"))

	// inject some games
	options = append(options, addGame("6b91f8d4-42ed-4148-bb20-eb9b31c91eb0", "1917", "sample-mdhender", "mdhender", "yojimbo"))
	options = append(options, addGame("5f03f14b-41e1-46d2-b273-f42c2a7d466e", "1812", "sample-yojimbo", "yojimbo", "mdhender"))

	srv, err := newServer(cfg, options...)
	if err != nil {
		return err
	}
	srv.Handler = routes(srv, http.StripPrefix("/", handlers.SPA(cfg.Server.PublicRoot)), cfg.Games.FileSavePath)

	log.Printf("[server] listening on %s\n", srv.Addr)
	return srv.ListenAndServe()
}
