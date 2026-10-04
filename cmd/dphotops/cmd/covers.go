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
	Short: "Reconcile cover sets on every album of every owner",
	Long: `Iterate every album of every owner and reconcile its cover set from the album's
eligible images: covers whose media is no longer in the album are stripped and
empty slots are filled randomly, up to 4. Existing RANDOM and CHERRY_PICKED
covers are preserved whenever possible.

Safe to re-run: the cover set is always capped at 4 and valid covers are kept.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		owners, err := pkgfactory.AclRepository(ctx).ListOwners(ctx)
		if err != nil {
			printer.FatalIfError(err, 1)
		}

		printer.Info("Backfilling covers across %d owner(s)", len(owners))

		backfill := &catalog.BackfillCovers{
			FindAlbumByOwnerPort: pkgfactory.AlbumQueries(ctx),
			CoverService:         pkgfactory.CoverServiceCase(ctx),
			Observers:            []catalog.CoverBackfillObserver{pkgfactory.AlbumView(ctx)},
		}

		var failedOwners int
		for _, owner := range owners {
			start := time.Now()
			changed, err := backfill.BackfillForOwner(ctx, owner)
			if err != nil {
				log.WithError(err).Errorf("covers backfill: %s FAILED after %s", owner, time.Since(start))
				printer.ErrorText("owner %s: %s", owner, err.Error())
				failedOwners++
				continue
			}
			printer.Success("owner %s backfilled in %s (%d album(s) updated)", owner, time.Since(start), len(changed))
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
