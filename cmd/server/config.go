// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"flag"
	"github.com/peterbourgon/ff/v3"
	"os"
	"time"
)

type config struct {
	Debug    bool
	FileName string
	Server   struct {
		Scheme  string
		Host    string
		Port    string
		Timeout struct {
			Idle  time.Duration
			Read  time.Duration
			Write time.Duration
		}
		PublicRoot string
		Salt       string
	}
	Cookies struct {
		HttpOnly bool
		Secure   bool
	}
	Games struct {
		FileSavePath string
	}
	SampleData *sampleData
}

type sampleData struct {
	Game struct {
		ID      string
		Name    string
		Players []string // list of names, not ID
	}
	Users struct {
		Mdhender, Yojimbo struct {
			ID, Name, Email string
		}
	}
}

// getConfig returns a configuration.
// It accepts an optional configuration file name.
// If provided, the file must contain a valid JSON object.
//
// The command line overrides environment variables overides configuration file override default values.
func getConfig() (*config, error) {
	var cfg config
	cfg.Games.FileSavePath = "../testdata/"
	cfg.SampleData = &sampleData{}
	cfg.SampleData.Game.ID = "6b91f8d4-42ed-4148-bb20-eb9b31c91eb0"
	cfg.SampleData.Game.Name = "sample"
	cfg.SampleData.Game.Players = []string{"mdhender", "yojimbo"}
	cfg.SampleData.Users.Mdhender.ID = "bf4c8168-6aab-409d-80cf-a4ee901904ef"
	cfg.SampleData.Users.Mdhender.Name = "mdhender"
	cfg.SampleData.Users.Mdhender.Email = "mdhender@example.com"
	cfg.SampleData.Users.Yojimbo.ID = "236bb1a5-1ae8-411a-a71f-791f4f03aa99"
	cfg.SampleData.Users.Yojimbo.Name = "yojimbo"
	cfg.SampleData.Users.Yojimbo.Email = "yojimbo@example.com"
	cfg.Server.Scheme = "http"
	cfg.Server.Host = "localhost"
	cfg.Server.Port = "8080"
	cfg.Server.Timeout.Idle = 10 * time.Second
	cfg.Server.Timeout.Read = 5 * time.Second
	cfg.Server.Timeout.Write = 10 * time.Second
	cfg.Server.PublicRoot = "../public"

	var (
		fs                    = flag.NewFlagSet("server", flag.ExitOnError)
		fileName              = fs.String("config", cfg.FileName, "config file (optional)")
		debug                 = fs.Bool("debug", cfg.Debug, "log debug information (optional)")
		gamesFileSavePath     = fs.String("file-save-path", cfg.Games.FileSavePath, "path to save game files to")
		serverCookiesHttpOnly = fs.Bool("cookies-http-only", cfg.Cookies.HttpOnly, "set HttpOnly flag on cookies")
		serverCookiesSecure   = fs.Bool("cookies-secure", cfg.Cookies.Secure, "set Secure flag on cookies")
		serverScheme          = fs.String("scheme", cfg.Server.Scheme, "http scheme, either 'http' or 'https'")
		serverHost            = fs.String("host", cfg.Server.Host, "host name (or IP) to listen on")
		serverPort            = fs.String("port", cfg.Server.Port, "port to listen on")
		serverPublicRoot      = fs.String("public-root", cfg.Server.PublicRoot, "path to serve static files from")
		serverSalt            = fs.String("salt", cfg.Server.Salt, "set salt for hashing")
		serverTimeoutIdle     = fs.Duration("idle-timeout", cfg.Server.Timeout.Idle, "http idle timeout")
		serverTimeoutRead     = fs.Duration("read-timeout", cfg.Server.Timeout.Read, "http read timeout")
		serverTimeoutWrite    = fs.Duration("write-timeout", cfg.Server.Timeout.Write, "http write timeout")
	)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("SERVER"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.JSONParser)); err != nil {
		return nil, err
	}

	cfg.Debug = *debug
	cfg.FileName = *fileName
	cfg.Cookies.HttpOnly = *serverCookiesHttpOnly
	cfg.Cookies.Secure = *serverCookiesSecure
	cfg.Games.FileSavePath = *gamesFileSavePath
	cfg.Server.Scheme = *serverScheme
	cfg.Server.Host = *serverHost
	cfg.Server.Port = *serverPort
	cfg.Server.PublicRoot = *serverPublicRoot
	cfg.Server.Salt = *serverSalt
	cfg.Server.Timeout.Idle = *serverTimeoutIdle
	cfg.Server.Timeout.Read = *serverTimeoutRead
	cfg.Server.Timeout.Write = *serverTimeoutWrite

	return &cfg, nil
}
