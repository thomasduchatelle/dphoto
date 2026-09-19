package catalogviews

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"slices"
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
	findAlbumsByIdsPort FindAlbumsByIdsPort,
	DriftObservers ...DriftOption,
) *OwnerDriftReconciler {
	observers := []DriftObserver{
		new(LoggerDriftObserver),
	}
	for _, option := range DriftObservers {
		observer := option.Observer()
		if observer != nil {
			observers = append(observers, observer)
		}
	}

	return &OwnerDriftReconciler{
		FindAlbumByOwnerPort:         findAlbumByOwnerPort,
		GetCurrentAlbumSummariesPort: getCurrentAlbumSummariesPort,
		AlbumReCounter: AlbumReCounter{
			ListUserWhoCanAccessAlbumPort: listUserWhoCanAccessAlbumPort,
			MediaCounterPort:              mediaCounterPort,
			FindAlbumsByIdsPort:           findAlbumsByIdsPort,
		},
		DriftDetector: &DriftDetector{
			GetCurrentAlbumSummariesPort: getCurrentAlbumSummariesPort,
			DriftObservers:               observers,
		},
	}
}

type OwnerDriftReconciler struct {
	FindAlbumByOwnerPort         FindAlbumByOwnerPort
	GetCurrentAlbumSummariesPort GetCurrentAlbumSummariesPort
	AlbumReCounter               AlbumReCounter
	DriftDetector                *DriftDetector
}

// Reconcile is re-computing counts for each album
func (d *OwnerDriftReconciler) Reconcile(ctx context.Context, owner ownermodel.Owner) error {
	albums, err := d.FindAlbumByOwnerPort.FindAlbumsByOwner(ctx, owner)
	if err != nil {
		return err
	}

	albumIds := make([]catalog.AlbumId, len(albums))
	for i, albumId := range albums {
		albumIds[i] = albumId.AlbumId
	}

	return d.AlbumReCounter.ReCountMedias(ctx, albumIds, d.DriftDetector)
}

type GetCurrentAlbumSummariesPort interface {
	ListSummariesForUserAndOwners(ctx context.Context, userId usermodel.UserId, owner ...ownermodel.Owner) ([]UserAlbumSummary, error)
}

type DriftDetector struct {
	GetCurrentAlbumSummariesPort GetCurrentAlbumSummariesPort
	DriftObservers               []DriftObserver
}

func (d *DriftDetector) PutSummaries(ctx context.Context, summaries []AlbumSummaryForUsers) error {
	expected := make(map[usermodel.UserId]map[catalog.AlbumId]UserAlbumSummary)
	var owners []ownermodel.Owner

	for _, summary := range summaries {
		for _, user := range summary.Users {
			userSummaries, ok := expected[user.UserId]
			if !ok {
				userSummaries = make(map[catalog.AlbumId]UserAlbumSummary)
			}

			userSummaries[summary.AlbumId] = UserAlbumSummary{
				AlbumSummary: summary.AlbumSummary,
				Availability: user,
			}
			expected[user.UserId] = userSummaries

			if !slices.Contains(owners, summary.AlbumId.Owner) {
				owners = append(owners, summary.AlbumId.Owner)
			}
		}
	}

	var drifts []Drift
	for userId, expectedForUser := range expected {
		currentAvailabilities, err := d.GetCurrentAlbumSummariesPort.ListSummariesForUserAndOwners(ctx, userId, owners...)
		if err != nil {
			return err
		}

		processed := make(map[catalog.AlbumId]any)
		for _, currentSummary := range currentAvailabilities {
			processed[currentSummary.AlbumSummary.AlbumId] = nil

			if expectedSummary, present := expectedForUser[currentSummary.AlbumSummary.AlbumId]; !present {
				drifts = append(drifts, NewNotExpectedDrift(currentSummary.Availability, currentSummary.AlbumSummary.AlbumId))

			} else if currentSummary.Availability != expectedSummary.Availability {
				drifts = append(
					drifts,
					NewNotExpectedDrift(currentSummary.Availability, currentSummary.AlbumSummary.AlbumId),
					NewMissingDrift(expectedSummary),
				)
			} else if currentSummary.AlbumSummary != expectedSummary.AlbumSummary {
				drifts = append(drifts, NewOverrideDrift(expectedSummary))
			}
		}

		for albumId, expectedSummary := range expectedForUser {
			if _, present := processed[albumId]; !present {
				drifts = append(drifts, NewMissingDrift(expectedSummary))
			}
		}
	}

	if len(drifts) > 0 {
		for _, observer := range d.DriftObservers {
			err := observer.OnDetectedDrifts(ctx, drifts)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type LoggerDriftObserver struct{}

func (l LoggerDriftObserver) OnDetectedDrifts(ctx context.Context, drifts []Drift) error {
	for _, drift := range drifts {
		if drift.Expected != nil {
			summary := drift.Expected.AvailableAlbumSummary
			availability := summary.Availability.String()

			if drift.Expected.Missing {
				log.Infof("drift: %-20s | %-30s | %-10s | %-5d", availability, summary.AlbumSummary.AlbumId, "MISSING", summary.AlbumSummary.MediaCount)
			} else {
				log.Infof("drift: %-20s | %-30s | %-10s | %-5d", availability, summary.AlbumSummary.AlbumId, "OVERRIDE", summary.AlbumSummary.MediaCount)

			}

		} else if drift.NotExpected != nil {
			log.Infof("drift: %-20s | %-30s | %-10s", drift.NotExpected.Availability, drift.NotExpected.AlbumId, "UNEXPECTED")
		}

	}
	return nil
}

type DriftSynchronizerPort interface {
	PutSummariesPort
	DeleteRowPort
}

type DriftSynchronizerObserver struct {
	DriftSynchronizerPort DriftSynchronizerPort
}

func (d *DriftSynchronizerObserver) OnDetectedDrifts(ctx context.Context, drifts []Drift) error {
	for _, drift := range drifts {
		switch {
		case drift.Expected != nil:
			err := d.DriftSynchronizerPort.PutSummaries(ctx, []AlbumSummaryForUsers{drift.Expected.AvailableAlbumSummary.ToSummaryForUsers()})
			if err != nil {
				return err
			}

		case drift.NotExpected != nil:
			err := d.DriftSynchronizerPort.DeleteRow(ctx, drift.NotExpected.Availability, drift.NotExpected.AlbumId)
			if err != nil {
				return err
			}

		default:
			return errors.Errorf("Drift not supported: %+v", drift)
		}
	}

	return nil
}
