package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	version "github.com/thomasduchatelle/dphoto/pkg/meta"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version",
	Long:    `Print the version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s [%s]", version.Version(), version.BuildVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
