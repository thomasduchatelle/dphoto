package backup

type IChannelPublisher interface {
	AnalysedMediaObserver
}

func ChannelPublisherImplementsObserverInterface() IChannelPublisher {
	return &ChannelPublisher{}
}
