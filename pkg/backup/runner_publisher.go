package backup

func newPublisher(volume SourceVolume) (runnerPublisher, int, error) {
	medias, err := volume.FindMedias()

	return func(channel chan FoundMedia, eventChannel chan *progressEvent) error {
		size := 0
		for _, media := range medias {
			size += media.Size()
			channel <- media
		}

		eventChannel <- &progressEvent{Type: trackScanComplete, Count: len(medias), Size: size}

		return nil
	}, len(medias), err
}
