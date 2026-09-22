package cmd

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

var driftArgs = struct {
	apply bool
}{}

var driftCmd = &cobra.Command{
	Use:   "drift",
	Short: "Reconcile the album-list view against canonical records (dry-run by default)",
	Long: `Iterate over every owner and re-project their album-list view from the canonical
records. Detects and repairs missing rows, stale counts, stale display fields, and
orphan rows.

Runs in dry mode by default: differences are logged only. Pass --apply to write
back the corrections.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		owners, err := pkgfactory.AclRepository(ctx).ListOwners(ctx)
		if err != nil {
			printer.FatalIfError(err, 1)
		}

		printer.Info("Reconciling %d owner(s) - apply=%t", len(owners), driftArgs.apply)

		reconciler := pkgfactory.OwnerDriftReconciler(ctx, !driftArgs.apply)

		var failed int
		for _, owner := range owners {
			start := time.Now()
			if err := reconciler.Reconcile(ctx, owner); err != nil {
				log.WithError(err).Errorf("drift: %s FAILED after %s", owner, time.Since(start))
				printer.ErrorText("owner %s: %s", owner, err.Error())
				failed++
				continue
			}
			printer.Success("owner %s reconciled in %s", owner, time.Since(start))
		}

		if failed > 0 {
			printer.ErrorText("%d owner(s) failed to reconcile", failed)
			os.Exit(2)
		}
	},
}

func init() {
	rootCmd.AddCommand(driftCmd)
	driftCmd.Flags().BoolVar(&driftArgs.apply, "apply", false, "write corrections back to the view (default: dry-run)")
}
