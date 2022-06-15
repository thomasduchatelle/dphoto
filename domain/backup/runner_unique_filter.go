package backup

func newUniqueFilter() runnerUniqueFilter {
	uniqueIndexes := make(map[string]interface{})

	return func(media *BackingUpMediaRequest, progressChannel chan *ProgressEvent) bool {
		if _, filterOut := uniqueIndexes[media.Id]; filterOut {
			progressChannel <- &ProgressEvent{Type: ProgressEventDuplicate, Count: 1, Size: media.AnalysedMedia.FoundMedia.Size()}
			return false
		}

		uniqueIndexes[media.Id] = nil
		progressChannel <- &ProgressEvent{Type: ProgressEventReadyForUpload, Count: 1, Size: media.AnalysedMedia.FoundMedia.Size()}
		return true
	}
}
