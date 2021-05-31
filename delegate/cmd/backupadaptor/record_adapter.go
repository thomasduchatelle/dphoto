package backupadaptor

import (
	"duchatelle.io/dphoto/dphoto/backup"
	"duchatelle.io/dphoto/dphoto/catalog"
	"duchatelle.io/dphoto/dphoto/cmd/ui"
)

func NewSuggestionRepository(suggestions []*backup.FoundAlbum) ui.RecordRepositoryPort {
	records := make([]*ui.Record, len(suggestions))

	for i, suggestion := range suggestions {
		records[i] = &ui.Record{
			Suggestion: true,
			FolderName: "",
			Name:       suggestion.Name,
			Start:      suggestion.Start,
			End:        suggestion.End,
			Count:      0,
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
