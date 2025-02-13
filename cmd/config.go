// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package cmd

import (
	"fmt"
	"github.com/mdhender/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"strings"
)

var config = struct {
	envPrefix  string
	homeFolder string
	version    semver.Version
}{
	envPrefix: "EMPYR",
	version:   semver.Version{Major: 0, Minor: 1, Patch: 0, PreRelease: "alpha"},
}

// bindConfig reads in config file and ENV variables if set.
func bindConfig(cmd *cobra.Command) error {
	// use the config file as determined by our command line and environment.
	viper.SetConfigFile(argsRoot.cfgFile)

	// the following code for binding viper and cobra taken from
	// https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/

	// Try to read the config file. Ignore file-not-found errors.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	} else {
		log.Printf("viper: using   config file: %q\n", viper.ConfigFileUsed())
		if debugConfig, ok := viper.Get("files.path").(string); ok && debugConfig != "" {
			viperDebugConfig := filepath.Clean(filepath.Join(debugConfig, "viper.json"))
			log.Printf("viper: writing config file: %q\n", viperDebugConfig)
			if err = viper.WriteConfigAs(viperDebugConfig); err != nil {
				return err
			}
		}
	}

	// read in environment variables that match
	viper.SetEnvPrefix(config.envPrefix)
	viper.AutomaticEnv()

	// bind the current command's flags to viper
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to EMPYR_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			_ = viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", config.envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})

	return nil
}
