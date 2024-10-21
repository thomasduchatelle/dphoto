package backup

import "context"

// NewInterrupterObserver creates a context that will be cancelled when an error occurs
func NewInterrupterObserver(ctx context.Context) (*InterrupterObserver, context.Context) {
	cancellableCtx, cancel := context.WithCancel(ctx)

	return &InterrupterObserver{
		ctx:    ctx,
		cancel: cancel,
	}, cancellableCtx

}

type InterrupterObserver struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *InterrupterObserver) OnRejectedMedia(found FoundMedia, err error) {
	c.cancel()
}
