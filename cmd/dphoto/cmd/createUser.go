package cmd

import (
	"github.com/logrusorgru/aurora/v3"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"strings"

	"github.com/spf13/cobra"
)

var (
	createUserArg = struct {
		email string
		owner string
	}{}

	CreateUserCase func(email, ownerOptional string) error
)
var createUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Create a user capable of backing up its media to a owner of its own",
	Run: func(cmd *cobra.Command, args []string) {
		email := strings.Trim(createUserArg.email, " ")
		if email == "" {
			printer.ErrorText("--email is mandatory")
		}
		err := CreateUserCase(email, createUserArg.owner)
		printer.FatalIfError(err, 1)

		printer.Success("User %s has been created", aurora.Cyan(email))
	},
}

func init() {
	rootCmd.AddCommand(createUserCmd)

	createUserCmd.Flags().StringVarP(&createUserArg.email, "email", "e", "", "email with which the user is identified and can authenticate with google account")
	createUserCmd.Flags().StringVarP(&createUserArg.email, "owner", "o", "", "(optional) identifier of the owner (tenant) on which this email will backup its media")
}
