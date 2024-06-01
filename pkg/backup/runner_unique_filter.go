package backup

func newUniqueFilter() runnerUniqueFilter {
	uniqueIndexes := make(map[string]interface{})

	return func(media *BackingUpMediaRequest, progressChannel chan *ProgressEvent) bool {
		uniqueId := media.CatalogReference.UniqueIdentifier()
		if _, filterOut := uniqueIndexes[uniqueId]; filterOut {
			progressChannel <- &ProgressEvent{Type: ProgressEventDuplicate, Count: 1, Size: media.AnalysedMedia.FoundMedia.Size()}
			return false
		}

		uniqueIndexes[uniqueId] = nil
		progressChannel <- &ProgressEvent{Type: ProgressEventReadyForUpload, Count: 1, Size: media.AnalysedMedia.FoundMedia.Size()}
		return true
	}
}
