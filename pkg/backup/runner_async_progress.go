package backup

// NewProgressObserver publishes on a channel the progress of the backup
func NewProgressObserver(sizeHint int) *ProgressObserver {
	return &ProgressObserver{
		EventChannel: make(chan *ProgressEvent, sizeHint*5),
	}
}

type ProgressObserver struct {
	EventChannel chan *ProgressEvent
}

func (c *ProgressObserver) OnRejectedMedia(found FoundMedia, err error) {
	c.EventChannel <- &ProgressEvent{Type: ProgressEventAnalysed, Count: 1, Size: found.Size()}
	c.EventChannel <- &ProgressEvent{Type: ProgressEventDuplicate, Count: 1, Size: found.Size()} // Hacky - Should be SKIPPED not ProgressEventDuplicate
}

func (c *ProgressObserver) OnDecoratedAnalyser(found FoundMedia, cacheHit bool) {
	if cacheHit {
		c.EventChannel <- &ProgressEvent{Type: ProgressEventAnalysedFromCache, Count: 1, Size: found.Size()}
	}
}

func (c *ProgressObserver) OnAnalysedMedia(media *AnalysedMedia) {
	c.EventChannel <- &ProgressEvent{Type: ProgressEventAnalysed, Count: 1, Size: media.FoundMedia.Size()}
}
