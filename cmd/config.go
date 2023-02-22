// empyr - a game engine for Empyrean Challenge
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

package cmd

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/playbymail/empyr/internal/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"strings"
)

var cfgRoot = struct {
	cfgFile    string
	envPrefix  string
	homeFolder string
	version    semver.Version
}{
	envPrefix: "EMPYR",
	version:   semver.Version{PreRelease: "alpha"},
}

// bindConfig reads in config file and ENV variables if set.
// logic for binding viper and cobra taken from
// https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/
func bindConfig(cmd *cobra.Command) error {
	var err error

	// Find home directory.
	cfgRoot.homeFolder, err = homedir.Dir()
	if err != nil {
		return err
	}

	if cfgRoot.cfgFile == "" { // use default location of ~/.empyr.json
		viper.AddConfigPath(cfgRoot.homeFolder)
		viper.SetConfigType("json")
		viper.SetConfigName(".empyr")
	} else { // Use config file from the flag.
		viper.SetConfigFile(cfgRoot.cfgFile)
	}

	// Try to read the config file. Ignore file-not-found errors.
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	} else {
		log.Printf("viper: using   config file: %q\n", viper.ConfigFileUsed())
		viperDebugConfig := filepath.Join(viper.Get("files.path").(string), "viper.json")
		if viperDebugConfig == "viper.json" {

		}
		log.Printf("viper: writing config file: %q\n", viperDebugConfig)
		if err = viper.WriteConfigAs(viperDebugConfig); err != nil {
			return err
		}
	}

	// read in environment variables that match
	viper.SetEnvPrefix(cfgRoot.envPrefix)
	viper.AutomaticEnv()

	// bind the current command's flags to viper
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to EMPYR_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			_ = viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", cfgRoot.envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})

	return nil
}

func init() {
	// find home directory
	homeFolder, err := homedir.Dir()
	if err != nil {
		panic(fmt.Errorf("cmd: config: init: %w", err))
	}
	cfgRoot.homeFolder = homeFolder
}
