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

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cmdRoot.AddCommand(cmdEnv)
}

var cmdEnv = &cobra.Command{
	Use:   "env",
	Short: "dump the environment",
	Long:  `Print out the environment settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%-30s == %q\n", "version", cfgRoot.version.String())
		fmt.Printf("%-30s == %q\n", "homeFolder", cfgRoot.homeFolder)
		fmt.Printf("%-30s == %q\n", "viperConfigFile", viper.ConfigFileUsed())
	},
}
