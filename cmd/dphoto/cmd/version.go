package cmd

import (
	"fmt"
	version "github.com/thomasduchatelle/dphoto/pkg/meta"

	"github.com/spf13/cobra"
)

var StaticVersion = "foo"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version",
	Long:    `Print the version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(StaticVersion)
		fmt.Println(version.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
