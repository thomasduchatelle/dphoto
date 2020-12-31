/*
Copyright © 2020 Thomas Duchatelle <duchatelle.thomas@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	debug = false
	info  = false
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dphoto",
	Short: "Backup photos and videos to your personal AWS Cloud",
	Long:  `Backup photos and videos to your personal AWS Cloud.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		log.SetOutput(os.Stdout)
		formatter := new(log.TextFormatter)
		formatter.FullTimestamp = true
		log.SetFormatter(formatter)

		if debug {
			log.SetLevel(log.DebugLevel)
		} else if info {
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}

		log.WithFields(log.Fields{
			"debug": debug,
			"info":  info,
		}).Traceln("Start logging")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&info, "info", false, "enable info logging")
}
