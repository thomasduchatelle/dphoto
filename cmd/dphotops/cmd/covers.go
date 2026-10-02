package cmd

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

var coversCmd = &cobra.Command{
	Use:   "covers",
	Short: "Album covers administration",
}

var coversBackfillCmd = &cobra.Command{
	Use:   "backfill",
	Short: "Fill empty cover slots on every album of every owner",
	Long: `Iterate every album of every owner and complete empty cover sets with random
RANDOM covers drawn from each album's eligible images.

Idempotent: albums already carrying a full set of covers are left untouched, and
CHERRY_PICKED covers are never altered. Safe to re-run.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		owners, err := pkgfactory.AclRepository(ctx).ListOwners(ctx)
		if err != nil {
			printer.FatalIfError(err, 1)
		}

		printer.Info("Backfilling covers across %d owner(s)", len(owners))

		backfill := &catalog.BackfillCovers{
			FindAlbumByOwnerPort: pkgfactory.AlbumQueries(ctx),
			CompleteCoversPort:   pkgfactory.CompleteCoversCase(ctx),
		}

		var failedOwners int
		for _, owner := range owners {
			start := time.Now()
			report, err := backfill.BackfillForOwner(ctx, owner)
			if err != nil {
				log.WithError(err).Errorf("covers backfill: %s FAILED after %s", owner, time.Since(start))
				printer.ErrorText("owner %s: %s", owner, err.Error())
				failedOwners++
				continue
			}
			for _, failure := range report.Failures {
				log.WithError(failure.Err).Errorf("covers backfill: album %s FAILED", failure.AlbumId)
				printer.ErrorText("  album %s: %s", failure.AlbumId, failure.Err.Error())
			}
			if len(report.Failures) > 0 {
				printer.ErrorText("owner %s: %d album(s) failed out of %d (%s)", owner, len(report.Failures), report.Albums, time.Since(start))
				failedOwners++
				continue
			}
			printer.Success("owner %s backfilled in %s (%d album(s))", owner, time.Since(start), report.Albums)
		}

		if failedOwners > 0 {
			printer.ErrorText("%d owner(s) failed to backfill", failedOwners)
			os.Exit(2)
		}
	},
}

func init() {
	rootCmd.AddCommand(coversCmd)
	coversCmd.AddCommand(coversBackfillCmd)
}
