package catalogviews

import (
	"context"
	"slices"

	log "github.com/sirupsen/logrus"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type DriftOption struct {
	observer     DriftObserver
	synchronizer DriftSynchronizerPort
}

// DriftOptionObserver adds a custom observer.
func DriftOptionObserver(observer DriftObserver) DriftOption {
	return DriftOption{observer: observer}
}

func DriftOptionSynchronizer(synchronizer DriftSynchronizerPort) DriftOption {
	return DriftOption{synchronizer: synchronizer}
}

// DriftOptionDryMode enable or not the DRY mode.
func DriftOptionDryMode(dry bool, synchronizer DriftSynchronizerPort) DriftOption {
	if !dry {
		return DriftOptionSynchronizer(synchronizer)
	}
	return DriftOptionObserver(nil)
}

func (o *DriftOption) Observer() DriftObserver {
	switch {
	case o.observer != nil:
		return o.observer
	case o.synchronizer != nil:
		return &DriftSynchronizerObserver{DriftSynchronizerPort: o.synchronizer}
	default:
		return nil
	}
}

// NewDriftReconciler creates a new DriftReconciler in DRY mode ; use the option DriftOptionSynchronizer to reconcile.
func NewDriftReconciler(
	findAlbumByOwnerPort FindAlbumByOwnerPort,
	getCurrentAlbumSummariesPort GetCurrentAlbumSummariesPort,
	listUserWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort,
	mediaCounterPort MediaCounterPort,
	driftOptions ...DriftOption,
) *OwnerDriftReconciler {
	observers := []DriftObserver{new(LoggerDriftObserver)}
	for _, option := range driftOptions {
		if observer := option.Observer(); observer != nil {
			observers = append(observers, observer)
		}
	}

	return &OwnerDriftReconciler{
		FindAlbumByOwnerPort: findAlbumByOwnerPort,
		AlbumSummaryReprojector: AlbumSummaryReprojector{
			ListUserWhoCanAccessAlbumPort: listUserWhoCanAccessAlbumPort,
			MediaCounterPort:              mediaCounterPort,
		},
		DriftDetector: &DriftDetector{
			GetCurrentAlbumSummariesPort: getCurrentAlbumSummariesPort,
		},
		DriftObservers: observers,
	}
}

type OwnerDriftReconciler struct {
	FindAlbumByOwnerPort    FindAlbumByOwnerPort
	AlbumSummaryReprojector AlbumSummaryReprojector
	DriftDetector           *DriftDetector
	DriftObservers          []DriftObserver
}

// Reconcile rebuilds the album-list projection for the owner, detects drifts against the current
// projection, hands them to every observer (logger, synchronizer, ...) and returns them so callers
// can render their own reports.
func (d *OwnerDriftReconciler) Reconcile(ctx context.Context, owner ownermodel.Owner) ([]Drift, error) {
	albums, err := d.FindAlbumByOwnerPort.FindAlbumsByOwner(ctx, owner)
	if err != nil {
		return nil, err
	}

	expected, err := d.AlbumSummaryReprojector.Reproject(ctx, albums)
	if err != nil {
		return nil, err
	}

	drifts, users, err := d.DriftDetector.Detect(ctx, expected)
	if err != nil {
		return nil, err
	}

	if len(drifts) > 0 {
		for _, observer := range d.DriftObservers {
			if err := observer.OnDetectedDrifts(ctx, drifts); err != nil {
				return nil, err
			}
		}
	}

	for _, observer := range d.DriftObservers {
		cleaner, ok := observer.(LegacyRowsCleaner)
		if !ok {
			continue
		}
		for _, userId := range users {
			if err := cleaner.OnReconciledUser(ctx, userId); err != nil {
				return nil, err
			}
		}
	}

	return drifts, nil
}

type LegacyRowsCleaner interface {
	OnReconciledUser(ctx context.Context, userId usermodel.UserId) error
}

type GetCurrentAlbumSummariesPort interface {
	ListSummariesForUserAndOwners(ctx context.Context, userId usermodel.UserId, owner ...ownermodel.Owner) ([]UserAlbumSummary, error)
}

// DriftDetector compares an expected projection against the current stored projection and returns
// the per-album drifts.
type DriftDetector struct {
	GetCurrentAlbumSummariesPort GetCurrentAlbumSummariesPort
}

func (d *DriftDetector) Detect(ctx context.Context, summaries []AlbumSummaryForUsers) ([]Drift, []usermodel.UserId, error) {
	expected := make(map[usermodel.UserId]map[catalog.AlbumId]UserAlbumSummary)
	var owners []ownermodel.Owner

	for _, summary := range summaries {
		for _, user := range summary.Users {
			userSummaries, ok := expected[user.UserId]
			if !ok {
				userSummaries = make(map[catalog.AlbumId]UserAlbumSummary)
				expected[user.UserId] = userSummaries
			}
			userSummaries[summary.AlbumId] = UserAlbumSummary{
				AlbumSummary: summary.AlbumSummary,
				Availability: user,
			}
			if !slices.Contains(owners, summary.AlbumId.Owner) {
				owners = append(owners, summary.AlbumId.Owner)
			}
		}
	}

	var drifts []Drift
	users := make([]usermodel.UserId, 0, len(expected))
	for userId, expectedForUser := range expected {
		users = append(users, userId)
		current, err := d.GetCurrentAlbumSummariesPort.ListSummariesForUserAndOwners(ctx, userId, owners...)
		if err != nil {
			return nil, nil, err
		}
		drifts = append(drifts, detectDriftsForUser(expectedForUser, current)...)
	}
	return drifts, users, nil
}

func detectDriftsForUser(expected map[catalog.AlbumId]UserAlbumSummary, current []UserAlbumSummary) []Drift {
	var drifts []Drift
	processed := make(map[catalog.AlbumId]struct{})

	for _, currentSummary := range current {
		processed[currentSummary.AlbumSummary.AlbumId] = struct{}{}

		expectedSummary, present := expected[currentSummary.AlbumSummary.AlbumId]
		switch {
		case !present:
			drifts = append(drifts, NewDeletedDrift(currentSummary))
		case currentSummary.Availability != expectedSummary.Availability:
			drifts = append(drifts, NewDeletedDrift(currentSummary), NewMissingDrift(expectedSummary))
		case hasSummaryDrift(currentSummary.AlbumSummary, expectedSummary.AlbumSummary):
			drifts = append(drifts, NewOverrideDrift(expectedSummary))
		}
	}

	for albumId, expectedSummary := range expected {
		if _, present := processed[albumId]; !present {
			drifts = append(drifts, NewMissingDrift(expectedSummary))
		}
	}

	return drifts
}

func hasSummaryDrift(a, b AlbumSummary) bool {
	return a.MediaCount != b.MediaCount ||
		a.Name != b.Name ||
		!a.Start.Equal(b.Start) ||
		!a.End.Equal(b.End)
}

type LoggerDriftObserver struct{}

func (l LoggerDriftObserver) OnDetectedDrifts(_ context.Context, drifts []Drift) error {
	for _, drift := range drifts {
		switch drift.Reason {
		case DriftReasonMissing, DriftReasonOverridden:
			summary := drift.Expected
			log.Infof("drift: %-20s | %-30s | %-10s | %-5d", summary.Availability, drift.AlbumId, drift.Reason, summary.AlbumSummary.MediaCount)
		case DriftReasonDeleted:
			log.Infof("drift: %-20s | %-30s | %-10s", drift.NotExpected.Availability, drift.AlbumId, drift.Reason)
		}
	}
	return nil
}

type DriftSynchronizerPort interface {
	PutSummariesPort
	DeleteRowPort
	DeleteLegacyRowsForUserPort
}

type DriftSynchronizerObserver struct {
	DriftSynchronizerPort DriftSynchronizerPort
}

func (d *DriftSynchronizerObserver) OnDetectedDrifts(ctx context.Context, drifts []Drift) error {
	for _, drift := range drifts {
		switch drift.Reason {
		case DriftReasonMissing, DriftReasonOverridden:
			if err := d.DriftSynchronizerPort.PutSummaries(ctx, []AlbumSummaryForUsers{drift.Expected.ToSummaryForUsers()}); err != nil {
				return err
			}
		case DriftReasonDeleted:
			if err := d.DriftSynchronizerPort.DeleteRow(ctx, drift.NotExpected.Availability, drift.AlbumId); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *DriftSynchronizerObserver) OnReconciledUser(ctx context.Context, userId usermodel.UserId) error {
	return d.DriftSynchronizerPort.DeleteLegacyRowsForUser(ctx, userId)
}
