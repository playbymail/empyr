// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cmd

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"time"
)

// cmdRoot represents the base command when called without any subcommands
var cmdRoot = &cobra.Command{
	Short: "empyr: a game engine",
	Long:  `empyr is an engine inspired by better games.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// bind viper and cobra here since this hook runs early and always
		log.Printf("cobra: running root.PersistentPreRunE")

		// find and bind home directory
		if homeFolder, ok := os.LookupEnv(config.envPrefix + "_HOME"); ok {
			config.homeFolder = homeFolder
		} else if homeFolder, err := homedir.Dir(); err != nil {
			panic(fmt.Errorf("cobra: config: init: %w", err))
		} else {
			config.homeFolder = homeFolder
		}

		// find and bind the configuration file
		if argsRoot.cfgFile != "" {
			// use config file from the command line
		} else if cfgFile, ok := os.LookupEnv(config.envPrefix + "_CONFIG"); ok {
			// use config file from the environment
			argsRoot.cfgFile = cfgFile
		} else {
			// use default location of ~/.empyr.json
			argsRoot.cfgFile = filepath.Clean(filepath.Join(config.homeFolder, ".empyr.json"))
		}
		if sb, err := os.Stat(argsRoot.cfgFile); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				log.Fatalf("missing config %s\n", argsRoot.cfgFile)
			}
		} else if sb.IsDir() {
			log.Fatalf("config %s is not a file\n", argsRoot.cfgFile)
		} else {
			viper.SetConfigFile(argsRoot.cfgFile)
			return bindConfig(cmd)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		started := time.Now()

		if argsRoot.timeSelf {
			elapsed := time.Now().Sub(started)
			log.Printf("elapsed time: %v\n", elapsed)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the root Command.
func Execute() error {
	cmdRoot.Version = config.version.String()

	cmdRoot.PersistentFlags().StringVar(&argsRoot.cfgFile, "config", "", "configuration file (optional)")
	cmdRoot.PersistentFlags().BoolVar(&argsRoot.timeSelf, "time", false, "time commands")

	cmdRoot.AddCommand(cmdEnv)

	cmdRoot.AddCommand(cmdGenerate)
	cmdGenerate.PersistentFlags().StringVarP(&argsGenerate.path, "path", "p", argsGenerate.path, "path containing game folders")
	cmdGenerate.AddCommand(cmdGenerateCluster)
	cmdGenerateCluster.Flags().StringVarP(&argsGenerateCluster.game, "game", "g", argsGenerateCluster.game, "code of game to update")
	cmdGenerateCluster.Flags().Float64VarP(&argsGenerateCluster.radius, "radius", "r", 15.0, "cluster radius")
	cmdGenerateCluster.Flags().StringVar(&argsGenerateCluster.mapFile, "html-map", "", "name of map file to create (optional)")
	cmdGenerateCluster.Flags().BoolVar(&argsGenerateCluster.kind.cluster, "clustered", false, "clustered point distribution")
	cmdGenerateCluster.Flags().BoolVar(&argsGenerateCluster.kind.surface, "surface", false, "surface point distribution")
	cmdGenerateCluster.Flags().BoolVar(&argsGenerateCluster.kind.uniform, "uniform", true, "uniform point distribution")
	cmdGenerateCluster.MarkFlagsMutuallyExclusive("clustered", "surface", "uniform")
	cmdGenerate.AddCommand(cmdGenerateGame)
	cmdGenerateGame.Flags().StringVarP(&argsGenerateGame.code, "code", "c", argsGenerateGame.code, "short code for game")
	if err := cmdGenerateGame.MarkFlagRequired("code"); err != nil {
		panic(fmt.Errorf("generate: game: code: %w", err))
	}
	cmdGenerateGame.Flags().StringVarP(&argsGenerateGame.descr, "descr", "d", argsGenerateGame.descr, "description of game")
	cmdGenerateGame.Flags().BoolVarP(&argsGenerateGame.force, "force", "f", argsGenerateGame.force, "force creation of game")
	cmdGenerateGame.Flags().StringVarP(&argsGenerateGame.name, "name", "n", argsGenerateGame.name, "short name for game")

	cmdRoot.AddCommand(cmdScan)

	cmdScan.AddCommand(cmdScanOrders)
	cmdScanOrders.Flags().StringVarP(&argsScanOrders.ordersPath, "path", "p", "", "path to orders files")
	if err := cmdScanOrders.MarkFlagRequired("path"); err != nil {
		panic(fmt.Errorf("scan: orders: %w", err))
	}

	return cmdRoot.Execute()
}

var argsRoot struct {
	cfgFile  string
	timeSelf bool
}
