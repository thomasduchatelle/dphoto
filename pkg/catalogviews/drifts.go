package catalogviews

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

type DriftObserver interface {
	OnDetectedDrifts(ctx context.Context, drifts []Drift) error
}

type Drift struct {
	Expected    *MissingOrInvalidDrift // Expected is exclusive from NotExpected, they can not be set together
	NotExpected *NotExpectedDrift      // NotExpected is exclusive from Expected, they can not be set together
}

type MissingOrInvalidDrift struct {
	AvailableAlbumSize UserAlbumSize
	Missing            bool
}

type NotExpectedDrift struct {
	Availability Availability
	AlbumId      catalog.AlbumId
}

func NewOverrideDrift(size UserAlbumSize) Drift {
	return Drift{
		Expected: &MissingOrInvalidDrift{
			AvailableAlbumSize: size,
		},
	}
}

func NewMissingDrift(size UserAlbumSize) Drift {
	return Drift{
		Expected: &MissingOrInvalidDrift{
			AvailableAlbumSize: size,
			Missing:            true,
		},
	}
}

func NewNotExpectedDrift(availability Availability, albumId catalog.AlbumId) Drift {
	return Drift{
		NotExpected: &NotExpectedDrift{
			Availability: availability,
			AlbumId:      albumId,
		},
	}
}
