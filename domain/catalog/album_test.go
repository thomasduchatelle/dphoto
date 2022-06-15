package catalog_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/mocks"
	"testing"
	"time"
)

var (
	RepositoryMock *mocks.RepositoryPort
)

func mockAdapters() {
	RepositoryMock = new(mocks.RepositoryPort)
	catalog.dbPort = RepositoryMock
}

const layout = "2006-01-02T15"

// it should create a new album and re-assign existing medias to it
func TestCreate(t *testing.T) {
	a := assert.New(t)
	mockAdapters()
	layout := "2006-01-02T15"
	const owner = "stark"

	start := time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 1, 03, 0, 0, 0, 0, time.UTC)

	RepositoryMock.On("FindAllAlbums", owner).Maybe().Return(catalog.AlbumCollection(), nil)

	RepositoryMock.On("InsertAlbum", catalog.Album{
		Owner:      owner,
		Name:       "Christm@s  2nd-week !\"£$%^&*",
		FolderName: "/2020-12_Christm_s_2nd-week",
		Start:      start,
		End:        end,
	}).Return(nil)

	RepositoryMock.On("UpdateMedias", catalog.NewFindMediaRequest(owner).WithAlbum("/2020-Q4", "/2021-Q1", "/Christmas_Holidays").WithinRange(start, catalog.MustParse(layout, "2020-12-31T18")).WithinRange(catalog.MustParse(layout, "2021-01-01T18"), end), MoveTo("/2020-12_Christm_s_2nd-week")).Return("", 0, nil)

	err := catalog.Create(catalog.CreateAlbum{
		Owner: owner,
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
	const owner = "stark"

	const deletedFolder = "/Christmas_Holidays"
	const q4 = "/2020-Q4"
	const q1 = "/2021-Q1"

	RepositoryMock.On("FindAllAlbums", owner).Maybe().Return(catalog.AlbumCollection(), nil)
	RepositoryMock.On("DeleteEmptyAlbum", owner, deletedFolder).Return(nil)

	RepositoryMock.On("UpdateMedias", catalog.NewFindMediaRequest(owner).WithAlbum(deletedFolder).WithinRange(catalog.MustParse(layout, "2020-12-26T00"), catalog.MustParse(layout, "2020-12-31T18")), MoveTo(q4)).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalog.NewFindMediaRequest(owner).WithAlbum(deletedFolder).WithinRange(catalog.MustParse(layout, "2021-01-01T18"), catalog.MustParse(layout, "2021-01-04T00")), MoveTo(q1)).Return("", 0, nil)

	// side effect - medias has never been assigned to these albums
	RepositoryMock.On("UpdateMedias", catalog.NewFindMediaRequest(owner).WithAlbum(deletedFolder, q4).WithinRange(catalog.MustParse(layout, "2020-12-18T00"), catalog.MustParse(layout, "2020-12-24T00")), MoveTo("/Christmas_First_Week")).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalog.NewFindMediaRequest(owner).WithAlbum(deletedFolder, q4, "/Christmas_First_Week").WithinRange(catalog.MustParse(layout, "2020-12-24T00"), catalog.MustParse(layout, "2020-12-26T00")), MoveTo("/Christmas_Day")).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalog.NewFindMediaRequest(owner).WithAlbum(deletedFolder, q1, q4).WithinRange(catalog.MustParse(layout, "2020-12-31T18"), catalog.MustParse(layout, "2021-01-01T18")), MoveTo("/New_Year")).Return("", 0, nil)

	// when
	err := catalog.DeleteAlbum(owner, deletedFolder, false)

	a.NoError(err)
	RepositoryMock.AssertExpectations(t)
}

func TestFind(t *testing.T) {
	a := assert.New(t)
	mockAdapters()
	const owner = "stark"

	album := catalog.Album{
		Owner:      owner,
		Name:       "My Album",
		FolderName: "/MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("FindAlbum", owner, "/MyAlbum").Return(&album, nil)

	got, err := catalog.FindAlbum(owner, "/MyAlbum")
	if a.NoError(err) {
		a.Equal(&album, got)
	}
}

func TestFindAll(t *testing.T) {
	a := assert.New(t)
	mockAdapters()
	const owner = "stark"

	album := &catalog.Album{
		Owner:      owner,
		Name:       "My Album",
		FolderName: "/MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("FindAllAlbums", owner).Return([]*catalog.Album{album}, nil)

	got, err := catalog.FindAllAlbums(owner)
	if a.NoError(err) {
		a.Equal([]*catalog.Album{album}, got)
	}
}

func TestRename_sameFolderName(t *testing.T) {
	a := assert.New(t)
	mockAdapters()
	const owner = "stark"

	album := catalog.Album{
		Owner:      owner,
		Name:       "My Album",
		FolderName: "/MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("FindAlbum", owner, "/MyAlbum").Return(&album, nil)
	RepositoryMock.On("UpdateAlbum", catalog.Album{
		Owner:      owner,
		Name:       "/My_Other_Album",
		FolderName: "/MyAlbum",
		Start:      album.Start,
		End:        album.End,
	}).Return(nil)

	err := catalog.RenameAlbum(owner, "/MyAlbum", "/My_Other_Album", false)
	a.NoError(err)
	RepositoryMock.AssertExpectations(t)
}

func TestRename_updateFolderName(t *testing.T) {
	a := assert.New(t)
	mockAdapters()
	const owner = "stark"

	album := catalog.Album{
		Owner:      owner,
		Name:       "/Christmas_Holidays",
		FolderName: "/Christmas_Holidays",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("FindAlbum", owner, "/Christmas_Holidays").Return(&album, nil)
	RepositoryMock.On("DeleteEmptyAlbum", owner, "/Christmas_Holidays").Return(nil)
	RepositoryMock.On("InsertAlbum", catalog.Album{
		Owner:      owner,
		Name:       "/Covid_Lockdown_3",
		FolderName: "/2020-12_Covid_Lockdown_3",
		Start:      album.Start,
		End:        album.End,
	}).Return(nil)
	RepositoryMock.On("UpdateMedias", catalog.NewFindMediaRequest(owner).WithAlbum("/Christmas_Holidays"), "/2020-12_Covid_Lockdown_3").Return("", 0, nil)

	err := catalog.RenameAlbum(owner, "/Christmas_Holidays", "/Covid_Lockdown_3", true)
	a.NoError(err)
	RepositoryMock.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	a := assert.New(t)
	mockAdapters()
	layout := "2006-01-02T15"
	const owner = "ironman"

	RepositoryMock.On("FindAllAlbums", owner).Maybe().Return(catalog.AlbumCollection(), nil)

	updatedFolder := "/Christmas_First_Week"
	updatedStart := catalog.MustParse(layout, "2020-12-21T00")
	updatedEnd := catalog.MustParse(layout, "2020-12-27T00")

	christmas := "/Christmas_Holidays"
	q4 := "/2020-Q4"
	RepositoryMock.On("UpdateMedias", catalog.NewFindMediaRequest(owner).WithAlbum(updatedFolder, q4).WithinRange(catalog.MustParse(layout, "2020-12-18T00"), catalog.MustParse(layout, "2020-12-21T00")), MoveTo(christmas)).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalog.NewFindMediaRequest(owner).WithAlbum(christmas, q4).WithinRange(catalog.MustParse(layout, "2020-12-21T00"), catalog.MustParse(layout, "2020-12-24T00")), MoveTo(updatedFolder)).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalog.NewFindMediaRequest(owner).WithAlbum(christmas, updatedFolder, q4).WithinRange(catalog.MustParse(layout, "2020-12-24T00"), catalog.MustParse(layout, "2020-12-26T00")), MoveTo("/Christmas_Day")).Return("", 0, nil)
	RepositoryMock.On("UpdateMedias", catalog.NewFindMediaRequest(owner).WithAlbum(christmas, q4).WithinRange(catalog.MustParse(layout, "2020-12-26T00"), catalog.MustParse(layout, "2020-12-27T00")), MoveTo(updatedFolder)).Return("", 0, nil)

	RepositoryMock.On("UpdateAlbum", catalog.Album{
		Owner:      owner,
		Name:       "",
		FolderName: updatedFolder,
		Start:      updatedStart,
		End:        updatedEnd,
	}).Return(nil)

	err := catalog.UpdateAlbum(owner, updatedFolder, updatedStart, updatedEnd)
	if a.NoError(err) {
		RepositoryMock.AssertExpectations(t)
	}
}

func MoveTo(folderName string) string {
	return folderName
}
