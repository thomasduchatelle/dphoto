package cmd

import (
	"context"
	"os"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
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
			drifts, err := reconciler.Reconcile(ctx, owner)
			if err != nil {
				log.WithError(err).Errorf("drift: %s FAILED after %s", owner, time.Since(start))
				printer.ErrorText("owner %s: %s", owner, err.Error())
				failed++
				continue
			}
			printer.Success("owner %s reconciled in %s (%d drift(s))", owner, time.Since(start), len(drifts))
			printDrifts(drifts)
		}

		if failed > 0 {
			printer.ErrorText("%d owner(s) failed to reconcile", failed)
			os.Exit(2)
		}
	},
}

// printDrifts groups drifts per viewer availability then prints one line per album so the operator
// can spot exactly which album needed which action.
func printDrifts(drifts []catalogviews.Drift) {
	if len(drifts) == 0 {
		return
	}

	byViewer := groupDriftsByViewer(drifts)
	viewers := make([]string, 0, len(byViewer))
	for viewer := range byViewer {
		viewers = append(viewers, viewer)
	}
	sort.Strings(viewers)

	for _, viewer := range viewers {
		group := byViewer[viewer]
		printer.Info("  %s (missing=%d overridden=%d deleted=%d)", viewer,
			countByReason(group, catalogviews.DriftReasonMissing),
			countByReason(group, catalogviews.DriftReasonOverridden),
			countByReason(group, catalogviews.DriftReasonDeleted))
		sort.Slice(group, func(i, j int) bool {
			return albumIdString(group[i].AlbumId) < albumIdString(group[j].AlbumId)
		})
		for _, drift := range group {
			printer.Info("    %-11s %s", drift.Reason, drift.AlbumId)
		}
	}
}

func groupDriftsByViewer(drifts []catalogviews.Drift) map[string][]catalogviews.Drift {
	groups := make(map[string][]catalogviews.Drift)
	for _, drift := range drifts {
		viewer := driftViewer(drift).String()
		groups[viewer] = append(groups[viewer], drift)
	}
	return groups
}

func driftViewer(drift catalogviews.Drift) catalogviews.Availability {
	if drift.Expected != nil {
		return drift.Expected.Availability
	}
	return drift.NotExpected.Availability
}

func countByReason(drifts []catalogviews.Drift, reason catalogviews.DriftReason) int {
	count := 0
	for _, drift := range drifts {
		if drift.Reason == reason {
			count++
		}
	}
	return count
}

func albumIdString(id catalog.AlbumId) string {
	return id.Owner.Value() + "/" + string(id.FolderName)
}

func init() {
	rootCmd.AddCommand(driftCmd)
	driftCmd.Flags().BoolVar(&driftArgs.apply, "apply", false, "write corrections back to the view (default: dry-run)")
}
