package catalog_test

import (
	"duchatelle.io/dphoto/dphoto/catalog"
	"duchatelle.io/dphoto/dphoto/mocks"
	"github.com/stretchr/testify/assert"
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

	RepositoryMock.On("InsertAlbum", catalog.Album{
		Name:       "Christm@s  2nd-week !\"£$%^&*",
		FolderName: "2020-12_Christm_s_2nd_week",
		Start:      start,
		End:        end,
	}).Return(nil)

	RepositoryMock.On("UpdateMedias", catalog.NewUpdateFilter().WithAlbum("2020-Q4", "2021-Q1", "Christmas Holidays").WithinRange(start, catalog.MustParse(layout, "2020-12-31T18")).WithinRange(catalog.MustParse(layout, "2021-01-01T18"), end), MoveTo("2020-12_Christm_s_2nd_week")).Return("", 0, nil)

	err := catalog.Create(catalog.CreateAlbum{
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

	const deletedFolder = "Christmas Holidays"
	const q4 = "2020-Q4"
	const q1 = "2021-Q1"

	RepositoryMock.On("FindAllAlbums").Maybe().Return(catalog.AlbumCollection(), nil)
	RepositoryMock.On("DeleteEmptyAlbum", deletedFolder).Return(nil)

	RepositoryMock.On("UpdateMedias", catalog.NewUpdateFilter().WithAlbum(deletedFolder).WithinRange(catalog.MustParse(layout, "2020-12-26T00"), catalog.MustParse(layout, "2020-12-31T18")), MoveTo(q4)).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalog.NewUpdateFilter().WithAlbum(deletedFolder).WithinRange(catalog.MustParse(layout, "2021-01-01T18"), catalog.MustParse(layout, "2021-01-04T00")), MoveTo(q1)).Return("", 0, nil)

	// side effect - medias has never been assigned to these albums
	RepositoryMock.On("UpdateMedias", catalog.NewUpdateFilter().WithAlbum(deletedFolder, q4).WithinRange(catalog.MustParse(layout, "2020-12-18T00"), catalog.MustParse(layout, "2020-12-24T00")), MoveTo("Christmas First Week")).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalog.NewUpdateFilter().WithAlbum(deletedFolder, q4, "Christmas First Week").WithinRange(catalog.MustParse(layout, "2020-12-24T00"), catalog.MustParse(layout, "2020-12-26T00")), MoveTo("Christmas Day")).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalog.NewUpdateFilter().WithAlbum(deletedFolder, q1, q4).WithinRange(catalog.MustParse(layout, "2020-12-31T18"), catalog.MustParse(layout, "2021-01-01T18")), MoveTo("New Year")).Return("", 0, nil)

	// when
	err := catalog.DeleteAlbum(deletedFolder, false)

	a.NoError(err)
	RepositoryMock.AssertExpectations(t)
}

func TestFind(t *testing.T) {
	a := assert.New(t)
	mockAdapters()

	album := catalog.Album{
		Name:       "My Album",
		FolderName: "MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("FindAlbum", "MyAlbum").Return(&album, nil)

	got, err := catalog.FindAlbum("MyAlbum")
	if a.NoError(err) {
		a.Equal(&album, got)
	}
}

func TestFindAll(t *testing.T) {
	a := assert.New(t)
	mockAdapters()

	album := &catalog.Album{
		Name:       "My Album",
		FolderName: "MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("FindAllAlbums").Return([]*catalog.Album{album}, nil)

	got, err := catalog.FindAllAlbums()
	if a.NoError(err) {
		a.Equal([]*catalog.Album{album}, got)
	}
}

func TestRename_sameFolderName(t *testing.T) {
	a := assert.New(t)
	mockAdapters()

	album := catalog.Album{
		Name:       "My Album",
		FolderName: "MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("FindAlbum", "MyAlbum").Return(&album, nil)
	RepositoryMock.On("UpdateAlbum", catalog.Album{
		Name:       "My Other Album",
		FolderName: "MyAlbum",
		Start:      album.Start,
		End:        album.End,
	}).Return(nil)

	err := catalog.RenameAlbum("MyAlbum", "My Other Album", false)
	a.NoError(err)
	RepositoryMock.AssertExpectations(t)
}

func TestRename_updateFolderName(t *testing.T) {
	a := assert.New(t)
	mockAdapters()

	album := catalog.Album{
		Name:       "Christmas Holidays",
		FolderName: "Christmas Holidays",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("FindAlbum", "Christmas Holidays").Return(&album, nil)
	RepositoryMock.On("DeleteEmptyAlbum", "Christmas Holidays").Return(nil)
	RepositoryMock.On("InsertAlbum", catalog.Album{
		Name:       "Covid Lockdown 3",
		FolderName: "2020-12_Covid_Lockdown_3",
		Start:      album.Start,
		End:        album.End,
	}).Return(nil)
	RepositoryMock.On("UpdateMedias", catalog.NewUpdateFilter().WithAlbum("Christmas Holidays"), "2020-12_Covid_Lockdown_3").Return("", 0, nil)

	err := catalog.RenameAlbum("Christmas Holidays", "Covid Lockdown 3", true)
	a.NoError(err)
	RepositoryMock.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	a := assert.New(t)
	mockAdapters()
	layout := "2006-01-02T15"

	RepositoryMock.On("FindAllAlbums").Maybe().Return(catalog.AlbumCollection(), nil)

	updatedFolder := "Christmas First Week"
	updatedStart := catalog.MustParse(layout, "2020-12-21T00")
	updatedEnd := catalog.MustParse(layout, "2020-12-27T00")

	christmas := "Christmas Holidays"
	q4 := "2020-Q4"
	RepositoryMock.On("UpdateMedias", catalog.NewUpdateFilter().WithAlbum(updatedFolder, q4).WithinRange(catalog.MustParse(layout, "2020-12-18T00"), catalog.MustParse(layout, "2020-12-21T00")), MoveTo(christmas)).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalog.NewUpdateFilter().WithAlbum(christmas, q4).WithinRange(catalog.MustParse(layout, "2020-12-21T00"), catalog.MustParse(layout, "2020-12-24T00")), MoveTo(updatedFolder)).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalog.NewUpdateFilter().WithAlbum(christmas, updatedFolder, q4).WithinRange(catalog.MustParse(layout, "2020-12-24T00"), catalog.MustParse(layout, "2020-12-26T00")), MoveTo("Christmas Day")).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalog.NewUpdateFilter().WithAlbum(christmas, q4).WithinRange(catalog.MustParse(layout, "2020-12-26T00"), catalog.MustParse(layout, "2020-12-27T00")), MoveTo(updatedFolder)).Return("", 0, nil)

	RepositoryMock.On("UpdateAlbum", catalog.Album{
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
