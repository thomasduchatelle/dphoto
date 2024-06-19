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
	getCurrentAlbumSizesPort GetCurrentAlbumSizesPort,
	listUserWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort,
	mediaCounterPort MediaCounterPort,
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
		FindAlbumByOwnerPort:     findAlbumByOwnerPort,
		GetCurrentAlbumSizesPort: getCurrentAlbumSizesPort,
		AlbumReCounter: AlbumReCounter{
			ListUserWhoCanAccessAlbumPort: listUserWhoCanAccessAlbumPort,
			MediaCounterPort:              mediaCounterPort,
		},
		DriftDetector: &DriftDetector{
			GetCurrentAlbumSizesPort: getCurrentAlbumSizesPort,
			DriftObservers:           observers,
		},
	}
}

type OwnerDriftReconciler struct {
	FindAlbumByOwnerPort     FindAlbumByOwnerPort
	GetCurrentAlbumSizesPort GetCurrentAlbumSizesPort
	AlbumReCounter           AlbumReCounter
	DriftDetector            *DriftDetector
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

type GetCurrentAlbumSizesPort interface {
	GetAlbumSizes(ctx context.Context, userId usermodel.UserId, owner ...ownermodel.Owner) ([]UserAlbumSize, error)
}

type DriftDetector struct {
	GetCurrentAlbumSizesPort GetCurrentAlbumSizesPort
	DriftObservers           []DriftObserver
}

func (d *DriftDetector) InsertAlbumSize(ctx context.Context, sizes []MultiUserAlbumSize) error {
	expected := make(map[usermodel.UserId]map[catalog.AlbumId]UserAlbumSize)
	var owners []ownermodel.Owner

	for _, size := range sizes {
		for _, user := range size.Users {
			userSizes, ok := expected[user.UserId]
			if !ok {
				userSizes = make(map[catalog.AlbumId]UserAlbumSize)
			}

			userSizes[size.AlbumId] = UserAlbumSize{
				AlbumSize:    size.AlbumSize,
				Availability: user,
			}
			expected[user.UserId] = userSizes

			if !slices.Contains(owners, size.AlbumId.Owner) {
				owners = append(owners, size.AlbumId.Owner)
			}
		}
	}

	var drifts []Drift
	for userId, expectedForUser := range expected {
		currentAvailabilities, err := d.GetCurrentAlbumSizesPort.GetAlbumSizes(ctx, userId, owners...)
		if err != nil {
			return err
		}

		processed := make(map[catalog.AlbumId]any)
		for _, currentSize := range currentAvailabilities {
			processed[currentSize.AlbumSize.AlbumId] = nil

			if expectedSize, present := expectedForUser[currentSize.AlbumSize.AlbumId]; !present {
				drifts = append(drifts, NewNotExpectedDrift(currentSize.Availability, currentSize.AlbumSize.AlbumId))

			} else if currentSize.Availability != expectedSize.Availability {
				drifts = append(
					drifts,
					NewNotExpectedDrift(currentSize.Availability, currentSize.AlbumSize.AlbumId),
					NewMissingDrift(expectedSize),
				)
			} else if currentSize.AlbumSize.MediaCount != expectedSize.AlbumSize.MediaCount {
				drifts = append(drifts, NewOverrideDrift(expectedSize))
			}
		}

		for albumId, expectedSize := range expectedForUser {
			if _, present := processed[albumId]; !present {
				drifts = append(drifts, NewMissingDrift(expectedSize))
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
			size := drift.Expected.AvailableAlbumSize
			availability := size.Availability.String()

			if drift.Expected.Missing {
				log.Infof("drift: %-20s | %-30s | %-10s | %-5d", availability, size.AlbumSize.AlbumId, "MISSING", size.AlbumSize.MediaCount)
			} else {
				log.Infof("drift: %-20s | %-30s | %-10s | %-5d", availability, size.AlbumSize.AlbumId, "OVERRIDE", size.AlbumSize.MediaCount)

			}

		} else if drift.NotExpected != nil {
			log.Infof("drift: %-20s | %-30s | %-10s", drift.NotExpected.Availability, drift.NotExpected.AlbumId, "UNEXPECTED")
		}

	}
	return nil
}

type DriftSynchronizerPort interface {
	InsertAlbumSizePort
	DeleteAlbumSizePort
}

type DriftSynchronizerObserver struct {
	DriftSynchronizerPort DriftSynchronizerPort
}

func (d *DriftSynchronizerObserver) OnDetectedDrifts(ctx context.Context, drifts []Drift) error {
	for _, drift := range drifts {
		switch {
		case drift.Expected != nil:
			err := d.DriftSynchronizerPort.InsertAlbumSize(ctx, []MultiUserAlbumSize{drift.Expected.AvailableAlbumSize.ToMultiUser()})
			if err != nil {
				return err
			}

		case drift.NotExpected != nil:
			err := d.DriftSynchronizerPort.DeleteAlbumSize(ctx, drift.NotExpected.Availability, drift.NotExpected.AlbumId)
			if err != nil {
				return err
			}

		default:
			return errors.Errorf("Drift not supported: %+v", drift)
		}
	}

	return nil
}
