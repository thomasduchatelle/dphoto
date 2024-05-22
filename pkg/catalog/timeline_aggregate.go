package catalog

type TimelineAggregate struct {
	timeline *Timeline
	albums   []*Album
}

func NewTimelineAggregate(albums []*Album) *TimelineAggregate {
	return &TimelineAggregate{
		albums: albums,
	}
}

func (t *TimelineAggregate) AddNew(addedAlbum Album) (MediaTransferRecords, error) {
	t.albums = append(t.albums, &addedAlbum)

	timeline, err := NewTimeline(t.albums)
	if err != nil {
		return nil, err
	}

	records := make(MediaTransferRecords)
	for _, seg := range timeline.FindForAlbum(addedAlbum.AlbumId) {
		if len(seg.Albums) > 1 {
			selector := MediaSelector{
				FromAlbums: extractAlbumIds(seg.Albums[1:]),
				Start:      seg.Start,
				End:        seg.End,
			}
			if selectors, found := records[addedAlbum.AlbumId]; found {
				records[addedAlbum.AlbumId] = append(selectors, selector)
			} else {
				records[addedAlbum.AlbumId] = []MediaSelector{selector}
			}
		}
	}

	if len(records) == 0 {
		return nil, nil
	}

	return records, nil
}
