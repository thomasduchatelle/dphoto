package cmd

import (
	"github.com/spf13/cobra"
)

var opsCmd = &cobra.Command{
	Use:   "ops",
	Short: "Operational subcommands - requires direct access to AWS",
}

func init() {
	rootCmd.AddCommand(opsCmd)
}
