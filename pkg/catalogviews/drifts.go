package catalogviews

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

// DriftReason categorises a Drift.
type DriftReason string

const (
	DriftReasonMissing    DriftReason = "MISSING"    // row absent from the projection but expected
	DriftReasonOverridden DriftReason = "OVERRIDDEN" // row present but with a stale content
	DriftReasonDeleted    DriftReason = "DELETED"    // row present but not expected any more (must be removed)
)

// Drift describes a single divergence between a projection and its canonical source, for one album
// and one viewer. Expected is set when the row must be written; NotExpected is set when a row must
// be removed. Both are set for a swap (availability changed).
type Drift struct {
	AlbumId     catalog.AlbumId
	Reason      DriftReason
	Expected    *UserAlbumSummary // Expected is the row that should exist ; set on MISSING and OVERRIDDEN.
	NotExpected *UserAlbumSummary // NotExpected is the row that currently exists but must be removed ; set on DELETED.
}

func NewMissingDrift(expected UserAlbumSummary) Drift {
	return Drift{
		AlbumId:  expected.AlbumSummary.AlbumId,
		Reason:   DriftReasonMissing,
		Expected: &expected,
	}
}

func NewOverrideDrift(expected UserAlbumSummary) Drift {
	return Drift{
		AlbumId:  expected.AlbumSummary.AlbumId,
		Reason:   DriftReasonOverridden,
		Expected: &expected,
	}
}

func NewDeletedDrift(current UserAlbumSummary) Drift {
	return Drift{
		AlbumId:     current.AlbumSummary.AlbumId,
		Reason:      DriftReasonDeleted,
		NotExpected: &current,
	}
}

type DriftObserver interface {
	OnDetectedDrifts(ctx context.Context, drifts []Drift) error
}
