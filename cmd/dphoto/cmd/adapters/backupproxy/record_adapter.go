package backupproxy

import (
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/ui"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

func NewSuggestionRepository(owner string, folders []*backup.ScannedFolder, rejectCount int) ui.SuggestionRecordRepositoryPort {
	records := make([]*ui.SuggestionRecord, len(folders))

	for i, folder := range folders {

		simplifiedDistribution := make(map[string]int)
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
		Owner:       owner,
		Records:     records,
		RejectCount: rejectCount,
	}
}

func NewAlbumRepository(owner string) ui.ExistingRecordRepositoryPort {
	return &dynamicAlbumRepository{
		owner: owner,
	}
}

type staticRecordRepository struct {
	Owner       string
	Records     []*ui.SuggestionRecord
	RejectCount int
}

func (r *staticRecordRepository) FindSuggestionRecords() []*ui.SuggestionRecord {
	return r.Records
}

func (r *staticRecordRepository) Count() int {
	return len(r.Records)
}

func (r *staticRecordRepository) Rejects() int {
	return r.RejectCount
}

type dynamicAlbumRepository struct {
	owner string
}

func (r *dynamicAlbumRepository) FindExistingRecords() ([]*ui.ExistingRecord, error) {
	albums, err := catalog.FindAllAlbums(r.owner)
	if err != nil {
		return nil, err
	}

	timeline, err := catalog.NewTimeline(albums)
	if err != nil {
		return nil, err
	}

	records := make([]*ui.ExistingRecord, len(albums))
	for i, album := range albums {

		records[i] = &ui.ExistingRecord{
			FolderName:    album.FolderName,
			Name:          album.Name,
			Start:         album.Start,
			End:           album.End,
			Count:         album.TotalCount,
			ActivePeriods: r.activePeriods(timeline, album),
		}
	}
	return records, err
}

func (r *dynamicAlbumRepository) activePeriods(timeline *catalog.Timeline, album *catalog.Album) []ui.Period {
	actives, _ := timeline.FindBetween(album.Start, album.End)
	var periods []ui.Period
	for _, active := range actives {
		if len(active.Albums) > 0 && active.Albums[0].FolderName == album.FolderName {
			periods = append(periods, ui.Period{
				Start: active.Start,
				End:   active.End,
			})
		}
	}

	return periods
}
