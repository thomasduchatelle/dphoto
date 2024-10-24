package backup

type IChannelPublisher interface {
	AnalysedMediaObserver
	CatalogReferencerObserver
}

func ChannelPublisherImplementsObserverInterface() IChannelPublisher {
	return &ChannelPublisher{}
}
