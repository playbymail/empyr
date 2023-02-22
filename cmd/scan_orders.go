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
	"github.com/playbymail/empyr/internal/scanner"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var cmdScanOrders = &cobra.Command{
	Use:   "scan-orders",
	Short: "scan an orders file",
	Long:  `Load an orders file, scan it, and report on all errors.`,
	Run: func(cmd *cobra.Command, args []string) {
		orders, err := os.ReadFile(cfgScanOrders.inputFileName)
		if err != nil {
			log.Fatal(fmt.Errorf("scan-orders: %w", err))
		}

		if r, err := scanner.Orders(orders); err != nil {
			log.Fatal(fmt.Errorf("scan-orders: %w", err))
		} else {
			fmt.Println(r)
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdScanOrders)

	cmdScanOrders.Flags().StringVarP(&cfgScanOrders.inputFileName, "input-file", "i", "", "file to load and scan")
	if err := cmdScanOrders.MarkFlagRequired("input-file"); err != nil {
		panic(fmt.Errorf("scan-orders: %w", err))
	}
}

var cfgScanOrders struct {
	inputFileName string
}
