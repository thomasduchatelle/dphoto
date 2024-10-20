package backup

// NewAsyncPublisher observes step outputs and publishes them on a channel for the next one to consume
func NewAsyncPublisher(sizeHint int) *ChannelPublisher {
	analysedMediaChannelPublisher := &ChannelPublisher{
		AnalysedMediaChannel: make(chan *AnalysedMedia, sizeHint),
	}
	return analysedMediaChannelPublisher
}

type ChannelPublisher struct {
	AnalysedMediaChannel chan *AnalysedMedia
}

func (a *ChannelPublisher) OnAnalysedMedia(media *AnalysedMedia) {
	a.AnalysedMediaChannel <- media
}

func (a *ChannelPublisher) Close() {
	close(a.AnalysedMediaChannel)
}
