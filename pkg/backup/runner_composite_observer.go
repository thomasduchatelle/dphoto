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
