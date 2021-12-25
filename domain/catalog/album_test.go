package catalog_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogmodel"
	"github.com/thomasduchatelle/dphoto/mocks"
	"testing"
	"time"
)

var (
	RepositoryMock *mocks.RepositoryPort
)

func mockAdapters() {
	RepositoryMock = new(mocks.RepositoryPort)
	catalog.Repository = RepositoryMock
}

const layout = "2006-01-02T15"

// it should create a new album and re-assign existing medias to it
func TestCreate(t *testing.T) {
	a := assert.New(t)
	mockAdapters()
	layout := "2006-01-02T15"

	start := time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 1, 03, 0, 0, 0, 0, time.UTC)

	RepositoryMock.On("FindAllAlbums").Maybe().Return(catalog.AlbumCollection(), nil)

	RepositoryMock.On("InsertAlbum", catalogmodel.Album{
		Name:       "Christm@s  2nd-week !\"£$%^&*",
		FolderName: "/2020-12_Christm_s_2nd-week",
		Start:      start,
		End:        end,
	}).Return(nil)

	RepositoryMock.On("UpdateMedias", catalogmodel.NewUpdateFilter().WithAlbum("/2020-Q4", "/2021-Q1", "/Christmas_Holidays").WithinRange(start, catalog.MustParse(layout, "2020-12-31T18")).WithinRange(catalog.MustParse(layout, "2021-01-01T18"), end), MoveTo("/2020-12_Christm_s_2nd-week")).Return("", 0, nil)

	err := catalog.Create(catalogmodel.CreateAlbum{
		Name:  "Christm@s  2nd-week !\"£$%^&*",
		Start: start,
		End:   end,
	})

	a.NoError(err)
	RepositoryMock.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	a := assert.New(t)
	mockAdapters()

	const deletedFolder = "/Christmas_Holidays"
	const q4 = "/2020-Q4"
	const q1 = "/2021-Q1"

	RepositoryMock.On("FindAllAlbums").Maybe().Return(catalog.AlbumCollection(), nil)
	RepositoryMock.On("DeleteEmptyAlbum", deletedFolder).Return(nil)

	RepositoryMock.On("UpdateMedias", catalogmodel.NewUpdateFilter().WithAlbum(deletedFolder).WithinRange(catalog.MustParse(layout, "2020-12-26T00"), catalog.MustParse(layout, "2020-12-31T18")), MoveTo(q4)).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalogmodel.NewUpdateFilter().WithAlbum(deletedFolder).WithinRange(catalog.MustParse(layout, "2021-01-01T18"), catalog.MustParse(layout, "2021-01-04T00")), MoveTo(q1)).Return("", 0, nil)

	// side effect - medias has never been assigned to these albums
	RepositoryMock.On("UpdateMedias", catalogmodel.NewUpdateFilter().WithAlbum(deletedFolder, q4).WithinRange(catalog.MustParse(layout, "2020-12-18T00"), catalog.MustParse(layout, "2020-12-24T00")), MoveTo("/Christmas_First_Week")).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalogmodel.NewUpdateFilter().WithAlbum(deletedFolder, q4, "/Christmas_First_Week").WithinRange(catalog.MustParse(layout, "2020-12-24T00"), catalog.MustParse(layout, "2020-12-26T00")), MoveTo("/Christmas_Day")).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalogmodel.NewUpdateFilter().WithAlbum(deletedFolder, q1, q4).WithinRange(catalog.MustParse(layout, "2020-12-31T18"), catalog.MustParse(layout, "2021-01-01T18")), MoveTo("/New_Year")).Return("", 0, nil)

	// when
	err := catalog.DeleteAlbum(deletedFolder, false)

	a.NoError(err)
	RepositoryMock.AssertExpectations(t)
}

func TestFind(t *testing.T) {
	a := assert.New(t)
	mockAdapters()

	album := catalogmodel.Album{
		Name:       "My Album",
		FolderName: "/MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("FindAlbum", "/MyAlbum").Return(&album, nil)

	got, err := catalog.FindAlbum("/MyAlbum")
	if a.NoError(err) {
		a.Equal(&album, got)
	}
}

func TestFindAll(t *testing.T) {
	a := assert.New(t)
	mockAdapters()

	album := &catalogmodel.Album{
		Name:       "My Album",
		FolderName: "/MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("FindAllAlbums").Return([]*catalogmodel.Album{album}, nil)

	got, err := catalog.FindAllAlbums()
	if a.NoError(err) {
		a.Equal([]*catalogmodel.Album{album}, got)
	}
}

func TestRename_sameFolderName(t *testing.T) {
	a := assert.New(t)
	mockAdapters()

	album := catalogmodel.Album{
		Name:       "My Album",
		FolderName: "/MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("FindAlbum", "/MyAlbum").Return(&album, nil)
	RepositoryMock.On("UpdateAlbum", catalogmodel.Album{
		Name:       "/My_Other_Album",
		FolderName: "/MyAlbum",
		Start:      album.Start,
		End:        album.End,
	}).Return(nil)

	err := catalog.RenameAlbum("/MyAlbum", "/My_Other_Album", false)
	a.NoError(err)
	RepositoryMock.AssertExpectations(t)
}

func TestRename_updateFolderName(t *testing.T) {
	a := assert.New(t)
	mockAdapters()

	album := catalogmodel.Album{
		Name:       "/Christmas_Holidays",
		FolderName: "/Christmas_Holidays",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("FindAlbum", "/Christmas_Holidays").Return(&album, nil)
	RepositoryMock.On("DeleteEmptyAlbum", "/Christmas_Holidays").Return(nil)
	RepositoryMock.On("InsertAlbum", catalogmodel.Album{
		Name:       "/Covid_Lockdown_3",
		FolderName: "/2020-12_Covid_Lockdown_3",
		Start:      album.Start,
		End:        album.End,
	}).Return(nil)
	RepositoryMock.On("UpdateMedias", catalogmodel.NewUpdateFilter().WithAlbum("/Christmas_Holidays"), "/2020-12_Covid_Lockdown_3").Return("", 0, nil)

	err := catalog.RenameAlbum("/Christmas_Holidays", "/Covid_Lockdown_3", true)
	a.NoError(err)
	RepositoryMock.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	a := assert.New(t)
	mockAdapters()
	layout := "2006-01-02T15"

	RepositoryMock.On("FindAllAlbums").Maybe().Return(catalog.AlbumCollection(), nil)

	updatedFolder := "/Christmas_First_Week"
	updatedStart := catalog.MustParse(layout, "2020-12-21T00")
	updatedEnd := catalog.MustParse(layout, "2020-12-27T00")

	christmas := "/Christmas_Holidays"
	q4 := "/2020-Q4"
	RepositoryMock.On("UpdateMedias", catalogmodel.NewUpdateFilter().WithAlbum(updatedFolder, q4).WithinRange(catalog.MustParse(layout, "2020-12-18T00"), catalog.MustParse(layout, "2020-12-21T00")), MoveTo(christmas)).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalogmodel.NewUpdateFilter().WithAlbum(christmas, q4).WithinRange(catalog.MustParse(layout, "2020-12-21T00"), catalog.MustParse(layout, "2020-12-24T00")), MoveTo(updatedFolder)).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalogmodel.NewUpdateFilter().WithAlbum(christmas, updatedFolder, q4).WithinRange(catalog.MustParse(layout, "2020-12-24T00"), catalog.MustParse(layout, "2020-12-26T00")), MoveTo("/Christmas_Day")).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalogmodel.NewUpdateFilter().WithAlbum(christmas, q4).WithinRange(catalog.MustParse(layout, "2020-12-26T00"), catalog.MustParse(layout, "2020-12-27T00")), MoveTo(updatedFolder)).Return("", 0, nil)

	RepositoryMock.On("UpdateAlbum", catalogmodel.Album{
		Name:       "",
		FolderName: updatedFolder,
		Start:      updatedStart,
		End:        updatedEnd,
	}).Return(nil)

	err := catalog.UpdateAlbum(updatedFolder, updatedStart, updatedEnd)
	if a.NoError(err) {
		RepositoryMock.AssertExpectations(t)
	}
}

func MoveTo(folderName string) string {
	return folderName
}
