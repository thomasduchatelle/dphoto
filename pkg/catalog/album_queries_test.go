package catalog_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mocks "github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"testing"
	"time"
)

func mockAdapters(t *testing.T) *mocks.RepositoryAdapter {
	mockRepository := mocks.NewRepositoryAdapter(t)
	catalog.Init(mockRepository)

	return mockRepository
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
	mockRepository := mockAdapters(t)
	const owner = catalog.Owner("stark")
	albumId := catalog.AlbumId{Owner: owner, FolderName: "/MyAlbum"}

	album := catalog.Album{
		AlbumId: albumId,
		Name:    "My Album",
		Start:   time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
		End:     time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	mockRepository.On("FindAlbumByIds", mock.Anything, albumId).Return([]*catalog.Album{&album}, nil)

	got, err := catalog.FindAlbum(albumId)
	if a.NoError(err) {
		a.Equal(&album, got)
	}
}

func TestFind_NotFound(t *testing.T) {
	a := assert.New(t)
	mockRepository := mockAdapters(t)
	const owner = "stark"

	mockRepository.On("FindAlbumByIds", mock.Anything, myAlbumId).Return(nil, nil)

	_, err := catalog.FindAlbum(myAlbumId)
	a.ErrorIs(err, catalog.AlbumNotFoundError)
}

func TestFindAll(t *testing.T) {
	a := assert.New(t)
	mockRepository := mockAdapters(t)

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
