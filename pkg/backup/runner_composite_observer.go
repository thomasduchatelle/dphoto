package backup

// CompositeRunnerObserver dispatches events to multiple observers of different types
type CompositeRunnerObserver struct {
	Observers []interface{}
}

func (c *CompositeRunnerObserver) OnAnalysedMedia(media *AnalysedMedia) {
	for _, observer := range c.Observers {
		if typed, ok := observer.(AnalysedMediaObserver); ok {
			typed.OnAnalysedMedia(media)
		}
	}
}

func (c *CompositeRunnerObserver) OnRejectedMedia(found FoundMedia, err error) {
	for _, observer := range c.Observers {
		if typed, ok := observer.(RejectedMediaObserver); ok {
			typed.OnRejectedMedia(found, err)
		}
	}
}
