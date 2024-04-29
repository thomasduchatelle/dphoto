package catalog

import "context"

type ListMediaIdsPort interface {
	ListMediaIdsFromSelector(ctx context.Context, selector []MediaSelector) ([]MediaId, error)
}

type TransferMediaListener struct {
	ListMediaIds ListMediaIdsPort
}

func (t *TransferMediaListener) TransferMedia(ctx context.Context, selector []MediaSelector) error {
	mediaIds, err := t.ListMediaIds.ListMediaIdsFromSelector(ctx, selector)
	if err != nil {
		return err
	}

	if len(mediaIds) > 0 {

	}

	return nil
}
