package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	Version = "1.0.0"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Aliases: []string{"v"},
	Short: "Print the version",
	Long:  `Print the version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
