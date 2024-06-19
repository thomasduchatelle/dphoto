package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

var (
	viewsCmdArgs = struct {
		dry bool
	}{}
)
var viewsCmd = &cobra.Command{
	Use:     "views <owner>",
	Short:   "Control and re-generate the views for an owner if necessary",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"view"},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		owner := ownermodel.Owner(args[0])

		counter := new(CounterDriftObserver)
		err := pkgfactory.OwnerDriftReconciler(ctx, viewsCmdArgs.dry, catalogviews.DriftOptionObserver(counter)).Reconcile(ctx, owner)
		printer.FatalIfError(err, 1)

		if counter.count == 0 {
			printer.Success("No drift detected")
		} else if viewsCmdArgs.dry {
			printer.Success("%d drift(s) detected", counter.count)
		} else {
			printer.Success("Views reconciled, %d drifts fixed", counter.count)
		}
	},
}

func init() {
	migrateCmd.AddCommand(viewsCmd)

	viewsCmd.Flags().BoolVar(&viewsCmdArgs.dry, "dry", false, "dry run, do not apply changes")
}

type CounterDriftObserver struct {
	count int
}

func (c *CounterDriftObserver) OnDetectedDrifts(ctx context.Context, drifts []catalogviews.Drift) error {
	c.count += len(drifts)
	return nil
}
