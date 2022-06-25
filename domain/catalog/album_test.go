package catalog_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/mocks"
	"testing"
	"time"
)

func mockAdapters(t *testing.T) (*mocks.RepositoryAdapter, *mocks.CArchiveAdapter) {
	mockRepository := mocks.NewRepositoryAdapter(t)
	mockArchive := mocks.NewCArchiveAdapter(t)
	catalog.Init(mockRepository, mockArchive)

	return mockRepository, mockArchive
}

const (
	layout = "2006-01-02T15"
	owner  = "ironman"
)

// it should create a new album and re-assign existing medias to it
func TestShouldCreateAnAlbumAndMoveMediasToIt(t *testing.T) {

	tests := []struct {
		name             string
		mediaIds         []string
		mockExpectations func(*mocks.RepositoryAdapter, *mocks.CArchiveAdapter)
	}{
		{"it should move medias to the newly created album", []string{"file_1", "file_2"}, func(mockRepository *mocks.RepositoryAdapter, mockArchive *mocks.CArchiveAdapter) {
			mockRepository.On("TransferMedias", owner, []string{"file_1", "file_2"}, "/2020-12_Christm_s_2nd-week").Once().Return(nil)
			mockArchive.On("MoveMedias", owner, []string{"file_1", "file_2"}, "/2020-12_Christm_s_2nd-week").Once().Return(nil)
		}},
		{"it should not call adapters to move medias if there is no media to move", nil, func(mockRepository *mocks.RepositoryAdapter, mockArchive *mocks.CArchiveAdapter) {}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			mockRepository, mockArchive := mockAdapters(t)

			start := time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC)
			end := time.Date(2021, 1, 03, 0, 0, 0, 0, time.UTC)

			mockRepository.On("FindAllAlbums", owner).Maybe().Return(catalog.AlbumCollection(), nil)

			mockRepository.On("InsertAlbum", catalog.Album{
				Owner:      owner,
				Name:       "Christm@s  2nd-week !\"£$%^&*",
				FolderName: "/2020-12_Christm_s_2nd-week",
				Start:      start,
				End:        end,
			}).Return(nil)

			mockRepository.On("FindMediaIds",
				catalog.NewFindMediaRequest(owner).
					WithAlbum("/2020-Q4", "/2021-Q1", "/Christmas_Holidays").
					WithinRange(start, catalog.MustParse(layout, "2020-12-31T18")).
					WithinRange(catalog.MustParse(layout, "2021-01-01T18"), end),
			).Once().Return(tt.mediaIds, nil)

			tt.mockExpectations(mockRepository, mockArchive)

			err := catalog.Create(catalog.CreateAlbum{
				Owner: owner,
				Name:  "Christm@s  2nd-week !\"£$%^&*",
				Start: start,
				End:   end,
			})

			a.NoError(err)
		})

	}
}

func TestShouldReassignMediasToOtherAlbumsWhenDeletingAnAlbum(t *testing.T) {
	a := assert.New(t)
	mockRepository, mockArchive := mockAdapters(t)

	const deletedFolder = "/Christmas_Holidays"
	const q4 = "/2020-Q4"
	const q1 = "/2021-Q1"

	mockRepository.On("FindAllAlbums", owner).Maybe().Return(catalog.AlbumCollection(), nil)
	mockRepository.On("DeleteEmptyAlbum", owner, deletedFolder).Return(nil)

	expectTransferredMedias(mockRepository, mockArchive,
		catalog.NewFindMediaRequest(owner).
			WithAlbum(deletedFolder).
			WithinRange(catalog.MustParse(layout, "2020-12-26T00"), catalog.MustParse(layout, "2020-12-31T18")),
		q4)
	expectTransferredMedias(mockRepository, mockArchive,
		catalog.NewFindMediaRequest(owner).
			WithAlbum(deletedFolder).
			WithinRange(catalog.MustParse(layout, "2021-01-01T18"), catalog.MustParse(layout, "2021-01-04T00")),
		q1)

	// side effect - clear-up other assignments that might have been missed
	expectTransferredMedias(mockRepository, mockArchive,
		catalog.NewFindMediaRequest(owner).
			WithAlbum(deletedFolder, q4).
			WithinRange(catalog.MustParse(layout, "2020-12-18T00"), catalog.MustParse(layout, "2020-12-24T00")),
		"/Christmas_First_Week")
	expectTransferredMedias(mockRepository, mockArchive,
		catalog.NewFindMediaRequest(owner).
			WithAlbum(deletedFolder, q4, "/Christmas_First_Week").
			WithinRange(catalog.MustParse(layout, "2020-12-24T00"), catalog.MustParse(layout, "2020-12-26T00")),
		"/Christmas_Day")
	expectTransferredMedias(mockRepository, mockArchive,
		catalog.NewFindMediaRequest(owner).
			WithAlbum(deletedFolder, q1, q4).
			WithinRange(catalog.MustParse(layout, "2020-12-31T18"), catalog.MustParse(layout, "2021-01-01T18")),
		"/New_Year")

	// when
	err := catalog.DeleteAlbum(owner, deletedFolder, false)

	a.NoError(err)
}

func TestFind(t *testing.T) {
	a := assert.New(t)
	mockRepository, _ := mockAdapters(t)
	const owner = "stark"

	album := catalog.Album{
		Owner:      owner,
		Name:       "My Album",
		FolderName: "/MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	mockRepository.On("FindAlbum", owner, "/MyAlbum").Return(&album, nil)

	got, err := catalog.FindAlbum(owner, "/MyAlbum")
	if a.NoError(err) {
		a.Equal(&album, got)
	}
}

func TestFindAll(t *testing.T) {
	a := assert.New(t)
	mockRepository, _ := mockAdapters(t)
	const owner = "stark"

	album := &catalog.Album{
		Owner:      owner,
		Name:       "My Album",
		FolderName: "/MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	mockRepository.On("FindAllAlbums", owner).Return([]*catalog.Album{album}, nil)

	got, err := catalog.FindAllAlbums(owner)
	if a.NoError(err) {
		a.Equal([]*catalog.Album{album}, got)
	}
}

func TestRename_sameFolderName(t *testing.T) {
	a := assert.New(t)
	mockRepository, _ := mockAdapters(t)
	const owner = "stark"

	album := catalog.Album{
		Owner:      owner,
		Name:       "My Album",
		FolderName: "/MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	mockRepository.On("FindAlbum", owner, "/MyAlbum").Return(&album, nil)
	mockRepository.On("UpdateAlbum", catalog.Album{
		Owner:      owner,
		Name:       "/My_Other_Album",
		FolderName: "/MyAlbum",
		Start:      album.Start,
		End:        album.End,
	}).Return(nil)

	err := catalog.RenameAlbum(owner, "/MyAlbum", "/My_Other_Album", false)
	a.NoError(err)
	mockRepository.AssertExpectations(t)
}

func TestShouldTransferMediasToNewAlbumWhenRenamingItsFolder(t *testing.T) {
	a := assert.New(t)
	mockRepository, mockArchive := mockAdapters(t)

	album := catalog.Album{
		Owner:      owner,
		Name:       "/Christmas_Holidays",
		FolderName: "/Christmas_Holidays",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	mockRepository.On("FindAlbum", owner, "/Christmas_Holidays").Return(&album, nil)
	mockRepository.On("DeleteEmptyAlbum", owner, "/Christmas_Holidays").Return(nil)
	mockRepository.On("InsertAlbum", catalog.Album{
		Owner:      owner,
		Name:       "/Covid_Lockdown_3",
		FolderName: "/2020-12_Covid_Lockdown_3",
		Start:      album.Start,
		End:        album.End,
	}).Return(nil)

	expectTransferredMedias(mockRepository, mockArchive,
		catalog.NewFindMediaRequest(owner).
			WithAlbum("/Christmas_Holidays"),
		"/2020-12_Covid_Lockdown_3")

	err := catalog.RenameAlbum(owner, "/Christmas_Holidays", "/Covid_Lockdown_3", true)
	a.NoError(err)
	mockRepository.AssertExpectations(t)
}

func TestShouldTransferAppropriatelyMediasBetweenAlbumsWhenDatesAreChanged(t *testing.T) {
	a := assert.New(t)
	mockRepository, mockArchive := mockAdapters(t)

	mockRepository.On("FindAllAlbums", owner).Maybe().Return(catalog.AlbumCollection(), nil)

	updatedFolder := "/Christmas_First_Week"
	updatedStart := catalog.MustParse(layout, "2020-12-21T00")
	updatedEnd := catalog.MustParse(layout, "2020-12-27T00")

	christmas := "/Christmas_Holidays"
	q4 := "/2020-Q4"
	expectTransferredMedias(mockRepository, mockArchive,
		catalog.NewFindMediaRequest(owner).
			WithAlbum(updatedFolder, q4).
			WithinRange(catalog.MustParse(layout, "2020-12-18T00"), catalog.MustParse(layout, "2020-12-21T00")), christmas)
	expectTransferredMedias(mockRepository, mockArchive,
		catalog.NewFindMediaRequest(owner).
			WithAlbum(christmas, q4).
			WithinRange(catalog.MustParse(layout, "2020-12-21T00"), catalog.MustParse(layout, "2020-12-24T00")), updatedFolder)
	expectTransferredMedias(mockRepository, mockArchive,
		catalog.NewFindMediaRequest(owner).
			WithAlbum(christmas, updatedFolder, q4).
			WithinRange(catalog.MustParse(layout, "2020-12-24T00"), catalog.MustParse(layout, "2020-12-26T00")), "/Christmas_Day")
	expectTransferredMedias(mockRepository, mockArchive,
		catalog.NewFindMediaRequest(owner).
			WithAlbum(christmas, q4).
			WithinRange(catalog.MustParse(layout, "2020-12-26T00"), catalog.MustParse(layout, "2020-12-27T00")), updatedFolder)

	mockRepository.On("UpdateAlbum", catalog.Album{
		Owner:      owner,
		Name:       "",
		FolderName: updatedFolder,
		Start:      updatedStart,
		End:        updatedEnd,
	}).Return(nil)

	err := catalog.UpdateAlbum(owner, updatedFolder, updatedStart, updatedEnd)
	if a.NoError(err) {
		mockRepository.AssertExpectations(t)
	}
}

func expectTransferredMedias(mockRepository *mocks.RepositoryAdapter, mockArchive *mocks.CArchiveAdapter, filter *catalog.FindMediaRequest, target string) {
	ids := []string{"to_" + target}
	mockRepository.On("FindMediaIds", filter).Once().Return(ids, nil)
	mockRepository.On("TransferMedias", owner, ids, target).Once().Return(nil)
	mockArchive.On("MoveMedias", owner, ids, target).Once().Return(nil)
}
