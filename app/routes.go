// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package app

import (
	"github.com/playbymail/empyr/app/actions"
	"github.com/playbymail/empyr/app/domains"
	"github.com/playbymail/empyr/app/responders"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func (a *App) Router() http.Handler {
	// Load templates
	tmpl := template.Must(template.ParseFiles(filepath.Join(a.Assets.Templates, "user-row.gohtml")))

	// Dependency injection
	userRepo := &InMemoryUserRepo{data: make(map[string]domains.User)}
	userService := &domains.UserService{Repo: userRepo}

	mux := http.NewServeMux()

	// static assets in the root directory need to be called out separately
	// serve specific root assets
	mux.Handle("GET /favicon.ico", http.FileServer(http.Dir(a.Assets.Files)))
	mux.Handle("GET /robots.txt", http.FileServer(http.Dir(a.Assets.Files)))
	// serve assets from specific directories
	log.Printf("assets: %s", a.Assets.Files)
	mux.Handle("GET /css/", http.FileServer(http.Dir(a.Assets.Files)))
	mux.Handle("GET /js/", http.FileServer(http.Dir(a.Assets.Files)))

	mux.HandleFunc("GET /home", a.Controllers.Home.Show)
	mux.HandleFunc("GET /blogs", a.Controllers.Blogs.Show)
	mux.HandleFunc("GET /reports", a.Controllers.Reports.Show)

	showLoginAction := actions.ShowLoginAction{Responder: &responders.ShowLoginResponder{Tmpl: tmpl}}
	mux.HandleFunc("GET /login/", showLoginAction.ServeHTTP)
	mux.HandleFunc("GET /login/{magicKey}", showLoginAction.ServeHTTP)

	logoutAction := actions.LogoutAction{Responder: responders.NewLogoutResponder(a.Assets.Templates)}
	mux.HandleFunc("GET /logout", logoutAction.ServeHTTP)
	mux.HandleFunc("POST /logout", logoutAction.ServeHTTP)

	createUserResponder := &responders.CreateUserResponder{Tmpl: tmpl}
	createUserAction := &actions.CreateUserAction{Service: userService, Responder: createUserResponder}

	mux.HandleFunc("POST /users", createUserAction.ServeHTTP)

	// the "/" route is special. it serves the landing page
	mux.HandleFunc("GET /", a.Controllers.Home.Show)

	return mux
}

// Mock repository implementation
type InMemoryUserRepo struct {
	data map[string]domains.User
}

func (repo *InMemoryUserRepo) Save(user domains.User) error {
	repo.data[user.Email] = user
	return nil
}
