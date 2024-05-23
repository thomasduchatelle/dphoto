package catalog

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// MediaFutureReference is the response of a simulation of inserting the media: the unique ID (if the media already exists) or a unique ID the media can use.
type MediaFutureReference struct {
	Signature          MediaSignature
	ProvisionalMediaId MediaId
	AlreadyExists      bool
}

type FindExistingSignaturePort interface {
	FindSignatures(ctx context.Context, owner ownermodel.Owner, signatures []MediaSignature) (map[MediaSignature]MediaId, error)
}

type MediasInsertSimulator struct {
	FindExistingSignaturePort FindExistingSignaturePort
}

func (m *MediasInsertSimulator) SimulateInsertingMedia(ctx context.Context, owner ownermodel.Owner, signatures []MediaSignature) ([]MediaFutureReference, error) {
	if len(signatures) == 0 {
		return nil, nil
	}

	existing, err := m.FindExistingSignaturePort.FindSignatures(ctx, owner, signatures)
	if err != nil {
		return nil, err
	}

	references := make([]MediaFutureReference, len(signatures))
	for index, signature := range signatures {
		mediaId, exists := existing[signature]
		if !exists {
			mediaId, err = GenerateMediaId(signature)
			if err != nil {
				return nil, errors.Wrapf(err, "SimulateInsertingMedia failed for %s owner", owner)
			}
		}

		references[index] = MediaFutureReference{
			Signature:          signature,
			ProvisionalMediaId: mediaId,
			AlreadyExists:      exists,
		}
	}

	return references, nil
}
