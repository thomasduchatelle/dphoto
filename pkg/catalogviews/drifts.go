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
	ExpectedSummary UserAlbumSummary
	Missing         bool
}

type NotExpectedDrift struct {
	Availability Availability
	AlbumId      catalog.AlbumId
}

func NewOverrideDrift(summary UserAlbumSummary) Drift {
	return Drift{
		Expected: &MissingOrInvalidDrift{
			ExpectedSummary: summary,
		},
	}
}

func NewMissingDrift(summary UserAlbumSummary) Drift {
	return Drift{
		Expected: &MissingOrInvalidDrift{
			ExpectedSummary: summary,
			Missing:         true,
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
