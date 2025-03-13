// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package app

import (
	"github.com/playbymail/empyr/app/actions"
	"github.com/playbymail/empyr/app/domains"
	"github.com/playbymail/empyr/app/responders"
	"github.com/playbymail/empyr/internal/services/auth"
	"github.com/playbymail/empyr/internal/services/games"
	"github.com/playbymail/empyr/internal/services/sessions"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func (a *App) Router(
	authService auth.Service,
	gamesService games.Service,
	sessionsService sessions.Service,
) http.Handler {
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

	showLoginResponder := responders.NewShowLoginResponder(responders.NewView("login", a.Assets.Templates, "login.gohtml"))
	showLoginAction := actions.ShowLoginAction{
		Sessions:  sessionsService,
		Responder: showLoginResponder,
	}
	mux.HandleFunc("GET /login", showLoginAction.ServeHTTP)

	loginUserAction := actions.LoginUserAction{
		Sessions:       sessionsService,
		Authentication: authService,
		Responder:      showLoginResponder,
	}
	mux.HandleFunc("GET /login/{handle}/{magicKey}", loginUserAction.ServeHTTP)
	mux.HandleFunc("POST /login/{handle}/{magicKey}", loginUserAction.ServeHTTP)

	logoutUserAction := actions.LogoutUserAction{
		Sessions:  sessionsService,
		Responder: responders.NewLogoutUserResponder(responders.NewView("logout", a.Assets.Templates, "logout.gohtml")),
	}
	mux.HandleFunc("GET /logout", logoutUserAction.ServeHTTP)
	mux.HandleFunc("POST /logout", logoutUserAction.ServeHTTP)

	showGamesResponder := responders.NewShowGamesResponder(responders.NewView("games", a.Assets.Templates, "games.gohtml"))
	showGamesAction := actions.ShowGamesAction{
		Sessions:  sessionsService,
		Games:     gamesService,
		Responder: showGamesResponder,
	}
	mux.HandleFunc("GET /games", showGamesAction.ServeHTTP)

	createUserResponder := &responders.CreateUserResponder{Tmpl: tmpl}
	createUserAction := &actions.CreateUserAction{Service: userService, Responder: createUserResponder}

	mux.HandleFunc("POST /users", createUserAction.ServeHTTP)

	// the "/" route is special. it serves the landing page but also serves as
	// the catch-all for all not-found routes.
	// Gotta love Go's routing.
	// Actually, I don't love this part of it.
	showLandingResponder := responders.NewShowLandingResponder(responders.NewView("landing", a.Assets.Templates, "landing.gohtml"))
	showLandingAction := actions.ShowLandingAction{
		Sessions:  sessionsService,
		Responder: showLandingResponder,
	}
	mux.HandleFunc("GET /", showLandingAction.ServeHTTP)

	return sessions.AddUserToContext(mux, sessionsService)

	//r := router.New(authn.Sessions(a.JotFactory))
	//
	//r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	//	log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	//	id := authn.User(r).ID()
	//	log.Printf("%s %s: id is %d\n", r.Method, r.URL.Path, id)
	//	fmt.Println("[the handler ran here]")
	//	_, _ = fmt.Fprintln(w, "Hello world of", r.URL.Path)
	//})
	//r.Get("/logout", logoutAction.ServeHTTP)
	//
	//return r
}

// Mock repository implementation
type InMemoryUserRepo struct {
	data map[string]domains.User
}

func (repo *InMemoryUserRepo) Save(user domains.User) error {
	repo.data[user.Email] = user
	return nil
}
