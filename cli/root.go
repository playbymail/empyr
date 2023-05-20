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

package cli

import (
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
		viper.SetConfigFile(argsRoot.cfgFile)

		return bindConfig(cmd)
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
	return cmdRoot.Execute()
}

var argsRoot struct {
	cfgFile  string
	timeSelf bool
}

func init() {
	cmdRoot.Version = config.version.String()

	cmdRoot.PersistentFlags().StringVar(&argsRoot.cfgFile, "config", "", "configuration file (optional)")
	cmdRoot.PersistentFlags().BoolVar(&argsRoot.timeSelf, "time", false, "time commands")
}
