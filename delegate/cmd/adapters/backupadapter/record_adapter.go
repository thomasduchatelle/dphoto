package backupadapter

import (
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"duchatelle.io/dphoto/dphoto/catalog"
	"duchatelle.io/dphoto/dphoto/cmd/ui"
)

func NewSuggestionRepository(folders []*backupmodel.ScannedFolder) ui.SuggestionRecordRepositoryPort {
	records := make([]*ui.SuggestionRecord, len(folders))

	for i, folder := range folders {

		simplifiedDistribution := make(map[string]uint)
		for day, dayCounter := range folder.Distribution {
			simplifiedDistribution[day] = dayCounter.Count
		}

		records[i] = &ui.SuggestionRecord{
			FolderName:   "." + folder.RelativePath,
			Name:         folder.Name,
			Start:        folder.Start,
			End:          folder.End,
			Distribution: simplifiedDistribution,
			Original:     folder,
		}
	}

	return &staticRecordRepository{
		Records: records,
	}
}

func NewAlbumRepository() ui.ExistingRecordRepositoryPort {
	return new(dynamicAlbumRepository)
}

type staticRecordRepository struct {
	Records []*ui.SuggestionRecord
}

func (r *staticRecordRepository) FindSuggestionRecords() ([]*ui.SuggestionRecord, error) {
	return r.Records, nil
}

type dynamicAlbumRepository struct{}

func (r *dynamicAlbumRepository) FindExistingRecords() ([]*ui.ExistingRecord, error) {
	albums, err := catalog.FindAllAlbumsWithStats()
	if err != nil {
		return nil, err
	}

	albumsWithoutStats := make([]*catalog.Album, len(albums))
	for i, a := range albums {
		albumsWithoutStats[i] = &a.Album
	}

	timeline, err := catalog.NewTimeline(albumsWithoutStats)
	if err != nil {
		return nil, err
	}

	records := make([]*ui.ExistingRecord, len(albums))
	for i, album := range albums {

		records[i] = &ui.ExistingRecord{
			FolderName:    album.Album.FolderName,
			Name:          album.Album.Name,
			Start:         album.Album.Start,
			End:           album.Album.End,
			Count:         uint(album.TotalCount()),
			ActivePeriods: r.activePeriods(timeline, album),
		}
	}
	return records, err
}

func (r *dynamicAlbumRepository) activePeriods(timeline *catalog.Timeline, album *catalog.AlbumStat) []ui.Period {
	actives, _ := timeline.FindBetween(album.Album.Start, album.Album.End)
	var periods []ui.Period
	for _, active := range actives {
		if len(active.Albums) > 0 && active.Albums[0].FolderName == album.Album.FolderName {
			periods = append(periods, ui.Period{
				Start: active.Start,
				End:   active.End,
			})
		}
	}

	return periods
}
