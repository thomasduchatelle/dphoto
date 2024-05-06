package catalog_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mocks "github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
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
	layout               = "2006-01-02T15"
	owner  catalog.Owner = "ironman"
)

var (
	myAlbumId = catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/MyAlbum")}
)

func TestFind_Found(t *testing.T) {
	a := assert.New(t)
	mockRepository, _ := mockAdapters(t)
	const owner = catalog.Owner("stark")
	albumId := catalog.AlbumId{Owner: owner, FolderName: "/MyAlbum"}

	album := catalog.Album{
		AlbumId: albumId,
		Name:    "My Album",
		Start:   time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:     time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	mockRepository.On("FindAlbums", mock.Anything, albumId).Return([]*catalog.Album{&album}, nil)

	got, err := catalog.FindAlbum(albumId)
	if a.NoError(err) {
		a.Equal(&album, got)
	}
}

func TestFind_NotFound(t *testing.T) {
	a := assert.New(t)
	mockRepository, _ := mockAdapters(t)
	const owner = "stark"

	mockRepository.On("FindAlbums", mock.Anything, myAlbumId).Return(nil, nil)

	_, err := catalog.FindAlbum(myAlbumId)
	a.ErrorIs(err, catalog.NotFoundError)
}

func TestFindAll(t *testing.T) {
	a := assert.New(t)
	mockRepository, _ := mockAdapters(t)

	album := &catalog.Album{
		AlbumId: myAlbumId,
		Name:    "My Album",
		Start:   time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:     time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	mockRepository.On("FindAlbumsByOwner", mock.Anything, owner).Return([]*catalog.Album{album}, nil)

	got, err := catalog.FindAllAlbums(owner)
	if a.NoError(err) {
		a.Equal([]*catalog.Album{album}, got)
	}
}

func TestRename_sameFolderName(t *testing.T) {
	a := assert.New(t)
	mockRepository, _ := mockAdapters(t)

	album := catalog.Album{
		AlbumId: myAlbumId,
		Name:    "My Album",
		Start:   time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:     time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	mockRepository.On("FindAlbums", mock.Anything, catalog.AlbumId{Owner: owner, FolderName: "/MyAlbum"}).Return([]*catalog.Album{&album}, nil)
	mockRepository.On("UpdateAlbum", mock.Anything, catalog.Album{
		AlbumId: myAlbumId,
		Name:    "My_Other_Album",
		Start:   album.Start,
		End:     album.End,
	}).Return(nil)

	err := catalog.RenameAlbum(myAlbumId, "My_Other_Album", false)
	a.NoError(err)
	mockRepository.AssertExpectations(t)
}

func TestShouldTransferMediasToNewAlbumWhenRenamingItsFolder(t *testing.T) {
	a := assert.New(t)
	mockRepository, mockArchive := mockAdapters(t)

	album := catalog.Album{
		AlbumId: myAlbumId,
		Name:    "/Christmas_Holidays",
		Start:   time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:     time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	mockRepository.On("FindAlbums", mock.Anything, catalog.AlbumId{Owner: owner, FolderName: "/Christmas_Holidays"}).Return([]*catalog.Album{&album}, nil)
	mockRepository.On("DeleteEmptyAlbum", mock.Anything, catalog.AlbumId{Owner: owner, FolderName: "/Christmas_Holidays"}).Return(nil)
	mockRepository.On("InsertAlbum", mock.Anything, catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: "/2020-12_Covid_Lockdown_3",
		},
		Name:  "Covid_Lockdown_3",
		Start: album.Start,
		End:   album.End,
	}).Return(nil)

	expectTransferredMedias(mockRepository, mockArchive,
		catalog.NewFindMediaRequest(owner).
			WithAlbum("/Christmas_Holidays"),
		"/2020-12_Covid_Lockdown_3")

	err := catalog.RenameAlbum(catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/Christmas_Holidays")}, "Covid_Lockdown_3", true)
	a.NoError(err)
	mockRepository.AssertExpectations(t)
}

func TestShouldTransferAppropriatelyMediasBetweenAlbumsWhenDatesAreChanged(t *testing.T) {
	a := assert.New(t)
	mockRepository, mockArchive := mockAdapters(t)

	mockRepository.On("FindAlbumsByOwner", mock.Anything, owner).Maybe().Return(catalog.AlbumCollection(), nil)

	updatedFolder := catalog.NewFolderName("/Christmas_First_Week")
	updatedStart := catalog.MustParse(layout, "2020-12-21T00")
	updatedEnd := catalog.MustParse(layout, "2020-12-27T00")

	christmas := catalog.NewFolderName("/Christmas_Holidays")
	q4 := catalog.NewFolderName("/2020-Q4")
	expectTransferredMedias(mockRepository, mockArchive,
		catalog.NewFindMediaRequest(owner).
			WithAlbum(updatedFolder, q4).
			WithinRange(catalog.MustParse(layout, "2020-12-18T00"), catalog.MustParse(layout, "2020-12-21T00")), christmas)
	expectTransferredMedias(
		mockRepository,
		mockArchive,
		catalog.NewFindMediaRequest(owner).WithAlbum(christmas, q4).WithinRange(catalog.MustParse(layout, "2020-12-21T00"), catalog.MustParse(layout, "2020-12-24T00")),
		updatedFolder,
	)
	expectTransferredMedias(mockRepository, mockArchive,
		catalog.NewFindMediaRequest(owner).
			WithAlbum(christmas, updatedFolder, q4).
			WithinRange(catalog.MustParse(layout, "2020-12-24T00"), catalog.MustParse(layout, "2020-12-26T00")), "/Christmas_Day")
	expectTransferredMedias(mockRepository, mockArchive,
		catalog.NewFindMediaRequest(owner).
			WithAlbum(christmas, q4).
			WithinRange(catalog.MustParse(layout, "2020-12-26T00"), catalog.MustParse(layout, "2020-12-27T00")), updatedFolder)

	mockRepository.On("UpdateAlbum", mock.Anything, catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: updatedFolder,
		},
		Name:  "",
		Start: updatedStart,
		End:   updatedEnd,
	}).Return(nil)

	err := catalog.UpdateAlbum(catalog.AlbumId{Owner: owner, FolderName: updatedFolder}, updatedStart, updatedEnd)
	if a.NoError(err) {
		mockRepository.AssertExpectations(t)
	}
}

func expectTransferredMedias(mockRepository *mocks.RepositoryAdapter, mockArchive *mocks.CArchiveAdapter, filter *catalog.FindMediaRequest, target catalog.FolderName) {
	ids := []catalog.MediaId{catalog.MediaId(fmt.Sprintf("to_%s", target))}
	mockRepository.On("FindMediaIds", mock.Anything, filter).Once().Return(ids, nil)
	mockRepository.On("TransferMedias", mock.Anything, owner, ids, target).Once().Return(nil)
	mockArchive.On("MoveMedias", owner, ids, target).Once().Return(nil)
}
