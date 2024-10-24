package backup

import "context"

// NewAsyncPublisher observes step outputs and publishes them on a channel for the next one to consume
func NewAsyncPublisher(sizeHint int, batchSize int) *ChannelPublisher {
	bufferedChannelSize := 1 + sizeHint/batchSize
	if sizeHint == 0 {
		bufferedChannelSize = 0
	}

	analysedMediaChannelPublisher := &ChannelPublisher{
		AnalysedMediaChannel:      make(chan *AnalysedMedia, sizeHint),
		FoundChannel:              make(chan FoundMedia, sizeHint),
		BufferedAnalysedChannel:   make(chan []*AnalysedMedia, bufferedChannelSize),
		CataloguedChannel:         make(chan *BackingUpMediaRequest, bufferedChannelSize),
		BufferedCataloguedChannel: make(chan []*BackingUpMediaRequest, bufferedChannelSize),
		CompletionChannel:         make(chan []error, 1),
	}
	return analysedMediaChannelPublisher
}

type ChannelPublisher struct {
	FoundChannel              chan FoundMedia
	AnalysedMediaChannel      chan *AnalysedMedia
	BufferedAnalysedChannel   chan []*AnalysedMedia
	CataloguedChannel         chan *BackingUpMediaRequest
	BufferedCataloguedChannel chan []*BackingUpMediaRequest
	CompletionChannel         chan []error
}

func (a *ChannelPublisher) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	a.AnalysedMediaChannel <- media
	return nil
}

func (a *ChannelPublisher) AnalysedMediaChannelCloser() {
	close(a.AnalysedMediaChannel)
}

func (a *ChannelPublisher) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	for _, request := range requests {
		a.CataloguedChannel <- &request
	}
	return nil
}

func (a *ChannelPublisher) CataloguedChannelCloser() {
	close(a.CataloguedChannel)
}

func (a *ChannelPublisher) WaitToFinish() []error {
	return <-a.CompletionChannel
}
