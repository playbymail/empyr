// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package app implements the application. This is the web front end.
// It is not the game engine or the web server or the repository.
package app

// https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
// https://dev.to/omg-a-bear/how-i-write-http-services-in-go-2k34

import (
	"context"
	"fmt"
	"github.com/playbymail/empyr/internal/actions"
	"github.com/playbymail/empyr/internal/commands"
	"github.com/playbymail/empyr/internal/controllers"
	"github.com/playbymail/empyr/internal/encryption"
	"github.com/playbymail/empyr/internal/jot"
	"github.com/playbymail/empyr/internal/ratelimiter"
	"github.com/playbymail/empyr/internal/services"
	"github.com/playbymail/empyr/internal/views"
	"github.com/playbymail/empyr/pkg/stdlib"
	"github.com/playbymail/empyr/store"
	"path/filepath"
)

type App struct {
	Assets struct {
		Files     string
		Templates string
	}
	Database struct {
		Store   *store.Store
		Context context.Context
	}
	Encrypter   *encryption.Encrypter
	RateLimiter *ratelimiter.Limiter
	Markdown    *services.Markdown
	Paddle      *services.Paddle

	Facades struct {
		Articles     *actions.ArticlesFacade
		Products     *actions.ProductsFacade
		Transactions *actions.TransactionsFacade
	}

	Controllers struct {
		Admin         *controllers.Admin
		Articles      *controllers.Articles
		Auth          *controllers.Auth
		Blogs         *controllers.Blogs
		Home          *controllers.Home
		Lqia          *controllers.Lqia
		PaddleWebhook *controllers.PaddleWebhook
		Ptg           *controllers.Ptg
		Purchases     *controllers.Purchases
		Reports       *controllers.Reports
	}

	Commands struct {
		PaddleMigrate *commands.PaddleMigrate
	}

	JotFactory *jot.Factory
	Views      *views.View
}

type Config struct{}

func New(repo *store.Store, jf *jot.Factory, assetsPath, templatesPath string, ctx context.Context) (*App, error) {
	// verify the asset paths
	if !stdlib.IsDirExists(assetsPath) {
		return nil, fmt.Errorf("%s: invalid path", assetsPath)
	} else if !stdlib.IsDirExists(templatesPath) {
		return nil, fmt.Errorf("%s: invalid path", templatesPath)
	}

	a := &App{}
	a.Assets.Files = assetsPath
	a.Assets.Templates = templatesPath
	a.Database.Store = repo
	a.Database.Context = ctx
	a.JotFactory = jf

	// wire up the controllers for the application
	// should we be creating views for the controllers here?
	if blogsView, err := views.NewView("blogs.gohtml", filepath.Join(a.Assets.Templates, "blogs.gohtml")); err != nil {
		return nil, err
	} else if a.Controllers.Blogs, err = controllers.NewBlogsController(a.Database.Store, blogsView); err != nil {
		return nil, err
	}
	if homeView, err := views.NewView("home.gohtml", filepath.Join(a.Assets.Templates, "home.gohtml")); err != nil {
		return nil, err
	} else if a.Controllers.Home, err = controllers.NewHomeController(a.Database.Store, homeView); err != nil {
		return nil, err
	}
	if reportsView, err := views.NewView("reports.gohtml", filepath.Join(a.Assets.Templates, "reports.gohtml")); err != nil {
		return nil, err
	} else if a.Controllers.Reports, err = controllers.NewReportsController(a.Database.Store, reportsView); err != nil {
		return nil, err
	}

	return a, nil
}
