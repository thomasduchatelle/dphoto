package scanui

import (
	"context"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/ui"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func NewSuggestionRepository(owner string, folders []*backup.ScannedFolder) ui.SuggestionRecordRepositoryPort {
	records := make([]*ui.SuggestionRecord, len(folders))
	rejectCount := 0

	for i, folder := range folders {
		rejectCount += folder.RejectsCount

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
	ctx := context.TODO()

	// TODO It's incorrect to assume userId = owner
	owner := ownermodel.Owner(r.owner)
	albums, err := pkgfactory.AlbumView(ctx).ListAlbums(ctx, usermodel.CurrentUser{UserId: usermodel.UserId(r.owner), Owner: &owner}, catalogviews.ListAlbumsFilter{OnlyDirectlyOwned: true})

	timeline, err := newTimeline(albums)
	if err != nil {
		return nil, err
	}

	records := make([]*ui.ExistingRecord, len(albums))
	for i, album := range albums {

		records[i] = &ui.ExistingRecord{
			FolderName:    album.FolderName.String(),
			Name:          album.Name,
			Start:         album.Start,
			End:           album.End,
			Count:         album.MediaCount,
			ActivePeriods: r.activePeriods(timeline, &album.Album),
		}
	}
	return records, err
}

func newTimeline(visibleAlbums []*catalogviews.VisibleAlbum) (*catalog.Timeline, error) {
	albums := make([]*catalog.Album, len(visibleAlbums))
	for i, visible := range visibleAlbums {
		albums[i] = &visible.Album
	}

	return catalog.NewTimeline(albums)
}

func (r *dynamicAlbumRepository) activePeriods(timeline *catalog.Timeline, album *catalog.Album) []ui.Period {
	segments := timeline.FindSegmentsBetweenAndFilter(album.Start, album.End, album.AlbumId)

	var periods []ui.Period
	for _, segment := range segments {
		periods = append(periods, ui.Period{
			Start: segment.Start,
			End:   segment.End,
		})
	}

	return periods
}
