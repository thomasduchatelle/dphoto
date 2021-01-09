package album

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	RepositoryMock *MockRepositoryPort
)

func mockAdapters() {
	RepositoryMock = new(MockRepositoryPort)
	Repository = RepositoryMock
}

// it should create a new album and re-assign existing medias to it
func TestCreate(t *testing.T) {
	a := assert.New(t)
	mockAdapters()
	layout := "2006-01-02T15"

	start := time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 1, 03, 0, 0, 0, 0, time.UTC)

	RepositoryMock.On("FindAll").Maybe().Return(albumCollection(), nil)

	RepositoryMock.On("Insert", Album{
		Name:       "Christm@s  2nd-week !\"£$%^&*",
		FolderName: "2020-12_Christm_s_2nd_week",
		Start:      start,
		End:        end,
	}).Return(nil)

	RepositoryMock.On("UpdateMedias", NewFilter().WithAlbum("2020-Q4", "2021-Q1", "Christmas Holidays").WithinRange(start, mustParse(layout, "2020-12-31T18")).WithinRange(mustParse(layout, "2021-01-01T18"), end), MoveTo("2020-12_Christm_s_2nd_week")).Return(nil)

	err := Create(CreateAlbum{
		Name:  "Christm@s  2nd-week !\"£$%^&*",
		Start: &start,
		End:   &end,
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

	RepositoryMock.On("FindAll").Maybe().Return(albumCollection(), nil)
	RepositoryMock.On("DeleteEmpty", deletedFolder).Return(nil)

	RepositoryMock.On("UpdateMedias", NewFilter().WithAlbum(deletedFolder).WithinRange(mustParse(layout, "2020-12-26T00"), mustParse(layout, "2020-12-31T18")), MoveTo(q4)).Return(nil)
	RepositoryMock.On("UpdateMedias", NewFilter().WithAlbum(deletedFolder).WithinRange(mustParse(layout, "2021-01-01T18"), mustParse(layout, "2021-01-04T00")), MoveTo(q1)).Return(nil)

	// side effect - medias has never been assigned to these albums
	RepositoryMock.On("UpdateMedias", NewFilter().WithAlbum(deletedFolder, q4).WithinRange(mustParse(layout, "2020-12-18T00"), mustParse(layout, "2020-12-24T00")), MoveTo("Christmas First Week")).Return(nil)
	RepositoryMock.On("UpdateMedias", NewFilter().WithAlbum(deletedFolder, q4, "Christmas First Week").WithinRange(mustParse(layout, "2020-12-24T00"), mustParse(layout, "2020-12-26T00")), MoveTo("Christmas Day")).Return(nil)
	RepositoryMock.On("UpdateMedias", NewFilter().WithAlbum(deletedFolder, q1, q4).WithinRange(mustParse(layout, "2020-12-31T18"), mustParse(layout, "2021-01-01T18")), MoveTo("New Year")).Return(nil)

	// when
	err := Delete(deletedFolder, false)

	a.NoError(err)
	RepositoryMock.AssertExpectations(t)
}

func TestFind(t *testing.T) {
	a := assert.New(t)
	mockAdapters()

	album := Album{
		Name:       "My Album",
		FolderName: "MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("Find", "MyAlbum").Return(&album, nil)

	got, err := Find("MyAlbum")
	if a.NoError(err) {
		a.Equal(&album, got)
	}
}

func TestFindAll(t *testing.T) {
	a := assert.New(t)
	mockAdapters()

	album := Album{
		Name:       "My Album",
		FolderName: "MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("FindAll").Return([]Album{album}, nil)

	got, err := FindAll()
	if a.NoError(err) {
		a.Equal([]Album{album}, got)
	}
}

func TestRename_sameFolderName(t *testing.T) {
	a := assert.New(t)
	mockAdapters()

	album := Album{
		Name:       "My Album",
		FolderName: "MyAlbum",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("Find", "MyAlbum").Return(&album, nil)
	RepositoryMock.On("Update", Album{
		Name:       "My Other Album",
		FolderName: "MyAlbum",
		Start:      album.Start,
		End:        album.End,
	}).Return(nil)

	err := Rename("MyAlbum", "My Other Album", false)
	a.NoError(err)
	RepositoryMock.AssertExpectations(t)
}

func TestRename_updateFolderName(t *testing.T) {
	a := assert.New(t)
	mockAdapters()

	album := Album{
		Name:       "Christmas Holidays",
		FolderName: "Christmas Holidays",
		Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	RepositoryMock.On("Find", "Christmas Holidays").Return(&album, nil)
	RepositoryMock.On("DeleteEmpty", "Christmas Holidays").Return(nil)
	RepositoryMock.On("Insert", Album{
		Name:       "Covid Lockdown 3",
		FolderName: "2020-12_Covid_Lockdown_3",
		Start:      album.Start,
		End:        album.End,
	}).Return(nil)
	RepositoryMock.On("UpdateMedias", NewFilter().WithAlbum("Christmas Holidays"), MediaUpdate{FolderName: "2020-12_Covid_Lockdown_3"}).Return(nil)

	err := Rename("Christmas Holidays", "Covid Lockdown 3", true)
	a.NoError(err)
	RepositoryMock.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	a := assert.New(t)
	mockAdapters()
	layout := "2006-01-02T15"

	RepositoryMock.On("FindAll").Maybe().Return(albumCollection(), nil)

	updatedFolder := "Christmas First Week"
	christmas := "Christmas Holidays"
	q4 := "2020-Q4"
	RepositoryMock.On("UpdateMedias", NewFilter().WithAlbum(updatedFolder, q4).WithinRange(mustParse(layout, "2020-12-18T00"), mustParse(layout, "2020-12-21T00")), MoveTo(christmas)).Return(nil)
	RepositoryMock.On("UpdateMedias", NewFilter().WithAlbum(christmas, q4).WithinRange(mustParse(layout, "2020-12-21T00"), mustParse(layout, "2020-12-24T00")), MoveTo(updatedFolder)).Return(nil)
	RepositoryMock.On("UpdateMedias", NewFilter().WithAlbum(christmas, updatedFolder, q4).WithinRange(mustParse(layout, "2020-12-24T00"), mustParse(layout, "2020-12-26T00")), MoveTo("Christmas Day")).Return(nil)
	RepositoryMock.On("UpdateMedias", NewFilter().WithAlbum(christmas, q4).WithinRange(mustParse(layout, "2020-12-26T00"), mustParse(layout, "2020-12-27T00")), MoveTo(updatedFolder)).Return(nil)

	err := Update(updatedFolder, mustParse(layout, "2020-12-21T00"), mustParse(layout, "2020-12-27T00"))
	if a.NoError(err) {
		RepositoryMock.AssertExpectations(t)
	}
}

func MoveTo(folderName string) MediaUpdate {
	return MediaUpdate{
		FolderName: folderName,
	}
}
