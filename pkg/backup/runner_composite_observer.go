package backup

import "context"

// CompositeRunnerObserver dispatches events to multiple observers of different types
type CompositeRunnerObserver struct {
	Observers []interface{}
}

func (c *CompositeRunnerObserver) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	for _, observer := range c.Observers {
		if typed, ok := observer.(AnalysedMediaObserver); ok {
			err := typed.OnAnalysedMedia(ctx, media)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *CompositeRunnerObserver) OnRejectedMedia(ctx context.Context, found FoundMedia, err error) {
	for _, observer := range c.Observers {
		if typed, ok := observer.(RejectedMediaObserver); ok {
			typed.OnRejectedMedia(ctx, found, err)
		}
	}
}

func (c *CompositeRunnerObserver) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	for _, observer := range c.Observers {
		if typed, ok := observer.(CatalogReferencerObserver); ok {
			err := typed.OnMediaCatalogued(ctx, requests)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *CompositeRunnerObserver) OnFilteredOut(ctx context.Context, media AnalysedMedia, reference CatalogReference, cause error) error {
	for _, observer := range c.Observers {
		if typed, ok := observer.(CataloguerFilterObserver); ok {
			err := typed.OnFilteredOut(ctx, media, reference, cause)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
