package backup

import "context"

type Interrupter interface {
	Cancel()
}

// NewInterrupterObserver creates a context that will be cancelled when an error occurs
func NewInterrupterObserver(ctx context.Context, options Options) (Interrupter, context.Context) {
	cancellableCtx, cancel := context.WithCancel(ctx)

	if options.SkipRejects {
		return &DefaultInterrupterObserver{
			cancel: cancel,
		}, cancellableCtx
	}

	return &AnalyserInterrupterObserver{
		DefaultInterrupterObserver: &DefaultInterrupterObserver{
			cancel: cancel,
		},
	}, cancellableCtx

}

// DefaultInterrupterObserver interrupts everything EXCEPT the analyser rejects
type DefaultInterrupterObserver struct {
	cancel context.CancelFunc
}

func (c *DefaultInterrupterObserver) Cancel() {
	c.cancel()
}

// AnalyserInterrupterObserver interrupts everything, including when the analyser rejected a media..
type AnalyserInterrupterObserver struct {
	*DefaultInterrupterObserver
}

func (c *AnalyserInterrupterObserver) OnRejectedMedia(found FoundMedia, err error) {
	c.cancel()
}
