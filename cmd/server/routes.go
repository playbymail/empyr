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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/playbymail/empyr/cmd/server/pkg/games"
	"github.com/playbymail/empyr/cmd/server/pkg/jsonapi"
	"github.com/playbymail/empyr/cmd/server/pkg/orders"
	"github.com/playbymail/empyr/cmd/server/pkg/systems"
	"github.com/playbymail/empyr/cmd/server/pkg/users"
	"github.com/playbymail/empyr/cmd/server/pkg/way"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

func routes(s *server, spa http.Handler, fileSavePath string) http.Handler {
	router := way.NewRouter()

	router.Handle("GET", "/api/games", s.getAllGames())
	router.Handle("GET", "/api/games/:id", s.getGame())
	router.Handle("GET", "/api/games/:id/players", s.getGamePlayers())
	router.Handle("GET", "/api/games/:id/players/:playerName", s.getGamePlayer())
	router.Handle("GET", "/api/games/:id/players/:playerName/printout", s.getGamePlayerPrintout())
	router.Handle("GET", "/api/games/:id/systems", s.getGameSystems())
	router.Handle("GET", "/api/games/:id/systems/:systemId", s.getGameSystem())
	router.Handle("GET", "/api/users", s.getAllUsers())
	router.Handle("GET", "/api/users/:id", s.getUser())
	router.Handle("GET", "/api/version", s.getVersion())

	router.Handle("POST", "/api/engine/restart", s.restart())
	router.Handle("POST", "/api/games", s.addGame())
	router.Handle("POST", "/api/games/:id/orders", s.postGameOrders())
	router.Handle("POST", "/api/games/:id/save", s.postGameSave(fileSavePath))
	router.Handle("POST", "/api/users", s.addUser())

	router.NotFound = spa

	return router
}

// addGame creates a new game and registers it with the engine.
func (s *server) addGame() http.HandlerFunc {
	type okResult struct {
		ID string `json:"id"`
	}
	type formData struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		SeedString string `json:"seed"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var input formData // should pull from form!
		meta, err := s.createGame(input.ID, input.Name, input.SeedString)
		if err != nil {
			jsonapi.Error(w, r, http.StatusBadRequest, ErrBadRequest)
			return
		}
		jsonapi.Ok(w, r, http.StatusOK, okResult{meta.ID})
	}
}

// addUser adds a new user.
func (s *server) addUser() http.HandlerFunc {
	type okResult struct {
		ID string `json:"id"`
	}
	type formData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var input formData // should pull from form!
		u, err := s.createUser("", input.Name, input.Email)
		if err != nil {
			if errors.Is(err, ErrDuplicateAddress) {
				jsonapi.Error(w, r, http.StatusBadRequest, err)
				return
			} else if errors.Is(err, ErrDuplicateUserName) {
				jsonapi.Error(w, r, http.StatusBadRequest, err)
				return
			}
			jsonapi.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		jsonapi.Ok(w, r, http.StatusOK, okResult{u.ID})
	}
}

// getAllGames returns all games
func (s *server) getAllGames() http.HandlerFunc {
	type detail struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`    // short ID for this game (should be unique)
		Players []string `json:"players"` // names of all players that have ever been in the game
	}
	type response []detail

	return func(w http.ResponseWriter, r *http.Request) {
		var list response = []detail{} // create an empty list since we never return nil
		for _, meta := range s.games {
			item := detail{ID: meta.ID, Name: meta.Name, Players: []string{}}
			for name := range meta.Players {
				item.Players = append(item.Players, name)
			}
			list = append(list, item)
		}
		jsonapi.Ok(w, r, http.StatusOK, list)
	}
}

// getAllUsers returns all users
func (s *server) getAllUsers() http.HandlerFunc {
	type detail struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	type okResult []detail
	return func(w http.ResponseWriter, r *http.Request) {
		var list okResult = []detail{} // create an empty list since we never return nil
		for _, user := range s.users.Filter(func(user *users.User) bool { return true }) {
			list = append(list, detail{
				ID:    user.ID,
				Email: user.Email,
				Name:  user.Name,
			})
		}
		jsonapi.Ok(w, r, http.StatusOK, list)
	}
}

// getGame returns a specific game
func (s *server) getGame() http.HandlerFunc {
	type response struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`    // short ID for this game (should be unique)
		Owner   string   `json:"owner"`   // user id of the game's owner (creator)
		Players []string `json:"players"` // names of all players that have ever been in the game
		Systems []*systems.System
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := way.Param(r.Context(), "id")
		meta, ok := s.games[id]
		if !ok {
			jsonapi.Error(w, r, http.StatusNotFound, ErrNoData)
			return
		}
		var game response
		game.ID = meta.ID
		game.Name = meta.Name
		game.Players = []string{}
		for name := range meta.Players {
			game.Players = append(game.Players, name)
		}
		game.Systems = meta.Game.Systems[:10]
		jsonapi.Ok(w, r, http.StatusOK, game)
	}
}

// getGamePlayer returns a specific player in a game
func (s *server) getGamePlayer() http.HandlerFunc {
	type response struct {
		ID         string `json:"id"`
		PlayerName string `json:"name"`
		UserName   string `json:"user"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := way.Param(ctx, "id")
		playerName := way.Param(ctx, "playerName")
		meta, ok := s.games[id]
		if !ok {
			jsonapi.Error(w, r, http.StatusNotFound, ErrNoData)
			return
		}
		var user *users.User
		for name, player := range meta.Players {
			if name == playerName {
				user = player
				break
			}
		}
		if user == nil {
			jsonapi.Error(w, r, http.StatusNotFound, ErrNoData)
			return
		}
		jsonapi.Ok(w, r, http.StatusOK, response{ID: user.ID, PlayerName: playerName, UserName: user.Name})
	}
}

// getGamePlayerPrintout returns the turn printout for a specific player
func (s *server) getGamePlayerPrintout() http.HandlerFunc {
	type response struct {
		GameNo     string `json:"gameNo"`
		PlayerName string `json:"name"`
		TurnNo     string `json:"turnNo"`
		Factories  []struct {
			Number int
			Levels int
			Types  []string
		} `json:"factories"`
		Population []struct {
			Type   string
			Number int
		}
		Weapons []struct {
			Type   string
			Number int
		}
		Transportation []struct {
			Type   string
			Number int
			Level  string
		}
		Deposits []struct {
			Type   string
			Number int
		}
		ID         string `json:"idNo"`
		HomeNation struct {
			ID              string `json:"id"`
			Location        string
			System          string
			NumberOfPlanets int
			Type            string
		}
		Race struct {
			ID   string
			Name string
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := way.Param(ctx, "id")
		playerName := way.Param(ctx, "playerName")
		meta, ok := s.games[id]
		if !ok {
			jsonapi.Error(w, r, http.StatusNotFound, ErrNoData)
			return
		}
		var user *users.User
		for name, player := range meta.Players {
			if name == playerName {
				user = player
				break
			}
		}
		if user == nil {
			jsonapi.Error(w, r, http.StatusNotFound, ErrNoData)
			return
		}
		// todo: finish turn printout
		jsonapi.Ok(w, r, http.StatusOK, response{PlayerName: playerName})
	}
}

// getGamePlayers returns all players in a game
func (s *server) getGamePlayers() http.HandlerFunc {
	type response struct {
		Names []string `json:"names"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := way.Param(r.Context(), "id")
		meta, ok := s.games[id]
		if !ok {
			jsonapi.Error(w, r, http.StatusNotFound, ErrNoData)
			return
		}
		var list response
		for name := range meta.Players {
			list.Names = append(list.Names, name)
		}
		jsonapi.Ok(w, r, http.StatusOK, list)
	}
}

// getGameSystem returns a specific system in a game
func (s *server) getGameSystem() http.HandlerFunc {
	type response struct {
		System *systems.System `json:"system"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := way.Param(r.Context(), "id")
		systemId := way.Param(r.Context(), "systemId")
		meta, ok := s.games[id]
		if !ok {
			jsonapi.Error(w, r, http.StatusNotFound, ErrNoData)
			return
		}
		var data response
		for _, system := range meta.Game.Systems {
			if system.ID == systemId || system.Name == systemId {
				data.System = system
				break
			}
		}
		if data.System == nil {
			jsonapi.Error(w, r, http.StatusNotFound, ErrNoData)
			return
		}
		jsonapi.Ok(w, r, http.StatusOK, data)
	}
}

// getGameSystems returns all systems in a game
func (s *server) getGameSystems() http.HandlerFunc {
	type response struct {
		Names []string `json:"systems"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := way.Param(r.Context(), "id")
		meta, ok := s.games[id]
		if !ok {
			jsonapi.Error(w, r, http.StatusNotFound, ErrNoData)
			return
		}
		var list response
		for _, system := range meta.Game.Systems {
			list.Names = append(list.Names, system.Name)
		}
		jsonapi.Ok(w, r, http.StatusOK, list)
	}
}

// getPlayer returns a specific player
func (s *server) getPlayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonapi.Error(w, r, http.StatusInternalServerError, fmt.Errorf("assert(server.getPlayer implemented)"))
	}
}

// getUser returns a specific user
func (s *server) getUser() http.HandlerFunc {
	type okResult struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := way.Param(r.Context(), "id")
		user := s.users.ID(id)
		if user == nil {
			jsonapi.Error(w, r, http.StatusNotFound, ErrNoData)
			return
		}
		jsonapi.Ok(w, r, http.StatusOK, okResult{user.ID, user.Name, user.Email})
	}
}

// getVersion returns the application version.
func (s *server) getVersion() http.HandlerFunc {
	type okResult struct {
		Major string `json:"major"`
		Minor string `json:"minor"`
		Patch string `json:"patch"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		jsonapi.Ok(w, r, http.StatusOK, okResult{versionMajor, versionMinor, versionPatch})
	}
}

// postGameOrders .
func (s *server) postGameOrders() http.HandlerFunc {
	type response struct {
		Message string `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := way.Param(r.Context(), "id")
		log.Printf("[orders] %s %s: id %q\n", r.Method, r.URL.Path, id)

		// Enforce a maximum read of 1MB from the request body.
		orders, errs := orders.Decode(http.MaxBytesReader(w, r.Body, 1<<20))
		if errs != nil {
			for _, err := range errs {
				log.Printf("[orders] %+v\n", err)
			}
			jsonapi.Error(w, r, http.StatusBadRequest, errs[0])
			return
		}
		log.Printf("[orders] %s %v\n", id, orders)

		// do something with the orders
		var debug bool
		for i, o := range orders {
			switch {
			case o.Debug != nil:
				debug = o.Debug.On
				if debug {
					log.Printf("[orders] %s %4d debug %v\n", id, i, *o.Debug)
				}
			case o.DefensiveSupport != nil:
				if debug {
					log.Printf("[orders] %s %4d defensiveSupport %v\n", id, i, *o.DefensiveSupport)
				}
			case o.Dock != nil:
				if debug {
					log.Printf("[orders] %s %4d dock %v\n", id, i, *o.Dock)
				}
			case o.Run != nil:
				if debug {
					log.Printf("[orders] %s %4d run %v\n", id, i, *o.Run)
				}
			case o.Undock != nil:
				if debug {
					log.Printf("[orders] %s %4d undock %v\n", id, i, *o.Undock)
				}
			}
		}

		jsonapi.Ok(w, r, http.StatusOK, response{"order accepted"})
	}
}

// postGameSave .
func (s *server) postGameSave(fileSavePath string) http.HandlerFunc {
	log.Printf("[game] save file %q\n", fileSavePath)
	type response struct {
		Message string `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := way.Param(r.Context(), "id")
		log.Printf("[orders] %s %s: id %q\n", r.Method, r.URL.Path, id)
		meta, ok := s.games[id]
		if !ok {
			jsonapi.Error(w, r, http.StatusNotFound, ErrNoData)
			return
		}
		gameSaveFile := path.Join(fileSavePath, id+".json")
		log.Printf("[game] save game %q %q\n", meta.Name, gameSaveFile)
		data, err := json.MarshalIndent(meta, "", "  ")
		if err != nil {
			jsonapi.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		if err = ioutil.WriteFile(gameSaveFile, data, 0644); err != nil {
			jsonapi.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		jsonapi.Ok(w, r, http.StatusOK, response{fmt.Sprintf("save file %q", gameSaveFile)})
	}
}

// restart the engine if the requestor is an admin.
func (s *server) restart() http.HandlerFunc {
	type okResult struct {
		Message string `json:"msg"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		e, err := NewEngine()
		if err != nil {
			// no choice but to halt the service here
			log.Printf("[server] restart engine failed %+v\n", err)
			os.Exit(2)
		}
		s.e = e
		jsonapi.Ok(w, r, http.StatusOK, okResult{"engine restarted"})
	}
}

// filterGames ...
func (s *server) filterGames(fn func(*games.Game) bool) []*games.Game {
	list := []*games.Game{} // create an empty list since we never want to return a nil list
	for _, meta := range s.games {
		if fn(meta.Game) {
			list = append(list, meta.Game)
		}
	}
	return list
}
