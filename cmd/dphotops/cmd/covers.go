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

type backfillCoversViewUpdater struct {
	ctx context.Context
}

func (b backfillCoversViewUpdater) UpdateCovers(ctx context.Context, albumId catalog.AlbumId, covers []catalog.Cover) error {
	repo := pkgfactory.AlbumViewRepository(b.ctx)
	if len(covers) == 0 {
		return repo.DeleteCoversForAllViewers(ctx, albumId)
	}
	return repo.PutCoversForAllViewers(ctx, albumId, covers)
}

var coversCmd = &cobra.Command{
	Use:   "covers",
	Short: "Album covers administration",
}

var coversBackfillCmd = &cobra.Command{
	Use:   "backfill",
	Short: "Reconcile cover sets on every album of every owner",
	Long: `Iterate every album of every owner and apply the cover invariant: drop every
RANDOM cover, keep every CHERRY_PICKED cover, and fill empty slots up to 4 with
RANDOM covers drawn from each album's eligible images.

Safe to re-run: CHERRY_PICKED covers are always preserved and the cover set is
capped at 4. Note that re-running MAY redraw existing RANDOM covers (they are
dropped and re-picked on every pass).`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		owners, err := pkgfactory.AclRepository(ctx).ListOwners(ctx)
		if err != nil {
			printer.FatalIfError(err, 1)
		}

		printer.Info("Backfilling covers across %d owner(s)", len(owners))

		backfill := &catalog.BackfillCovers{
			FindAlbumByOwnerPort:      pkgfactory.AlbumQueries(ctx),
			RefreshCoversPort:         pkgfactory.CoverMaintenanceCase(ctx),
			BackfillCoversViewUpdater: backfillCoversViewUpdater{ctx: ctx},
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
