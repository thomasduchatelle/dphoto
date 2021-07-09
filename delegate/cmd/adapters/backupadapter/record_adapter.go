package backupadapter

import (
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"duchatelle.io/dphoto/dphoto/catalog"
	"duchatelle.io/dphoto/dphoto/cmd/ui"
)

func NewSuggestionRepository(suggestions []*backupmodel.ScannedFolder) ui.RecordRepositoryPort {
	records := make([]*ui.Record, len(suggestions))

	for i, suggestion := range suggestions {

		count := uint(0)
		for _, dayCounter := range suggestion.Distribution {
			count += dayCounter.Count
		}

		records[i] = &ui.Record{
			Suggestion: true,
			FolderName: suggestion.FolderName,
			Name:       suggestion.Name,
			Start:      suggestion.Start,
			End:        suggestion.End,
			Count:      count,
		}
	}

	return &staticRecordRepository{
		Records: records,
	}
}

func NewAlbumRepository() ui.RecordRepositoryPort {
	return new(dynamicAlbumRepository)
}

type staticRecordRepository struct {
	Records []*ui.Record
}

func (r *staticRecordRepository) FindRecords() ([]*ui.Record, error) {
	return r.Records, nil
}

type dynamicAlbumRepository struct{}

func (r *dynamicAlbumRepository) FindRecords() ([]*ui.Record, error) {
	albums, err := catalog.FindAllAlbumsWithStats()
	records := make([]*ui.Record, len(albums))

	for i, album := range albums {
		records[i] = &ui.Record{
			Suggestion: false,
			FolderName: album.Album.FolderName,
			Name:       album.Album.Name,
			Start:      album.Album.Start,
			End:        album.Album.End,
			Count:      uint(album.TotalCount()),
		}
	}
	return records, err
}
