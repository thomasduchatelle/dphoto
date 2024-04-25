package cmd

import (
	"github.com/spf13/cobra"
)

// dynamodbCmd represents the dynamodb command
var dynamodbCmd = &cobra.Command{
	Use:   "dynamodb",
	Short: "Operate on the DynamoDB database",
}

func init() {
	rootCmd.AddCommand(dynamodbCmd)
}
