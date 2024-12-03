package backup

import "context"

func newPublisher(volume SourceVolume) (runnerPublisher, int, error) {
	medias, err := volume.FindMedias(context.TODO())

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
