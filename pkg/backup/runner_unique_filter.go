package backup

func newUniqueFilter() runnerUniqueFilter {
	uniqueIndexes := make(map[string]interface{})

	return func(media *BackingUpMediaRequest, progressChannel chan *progressEvent) bool {
		uniqueId := media.CatalogReference.UniqueIdentifier()
		if _, filterOut := uniqueIndexes[uniqueId]; filterOut {
			progressChannel <- &progressEvent{Type: trackDuplicatedInVolume, Count: 1, Size: media.AnalysedMedia.FoundMedia.Size()}
			return false
		}

		uniqueIndexes[uniqueId] = nil
		progressChannel <- &progressEvent{Type: trackCatalogued, Count: 1, Size: media.AnalysedMedia.FoundMedia.Size()}
		return true
	}
}
