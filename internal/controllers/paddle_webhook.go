// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package controllers

import (
	"github.com/playbymail/empyr/internal/config"
	"github.com/playbymail/empyr/internal/services"
)

type PaddleWebhook struct {
	Config *config.Config
	Paddle *services.Paddle
}
