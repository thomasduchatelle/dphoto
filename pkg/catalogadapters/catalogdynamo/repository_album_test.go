package catalogdynamo

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"testing"
)

type AlbumCrudTestSuite struct {
	suite.Suite
	suffix string
	repo   *Repository
	owner  catalog.Owner
}

func TestRepositoryAlbum(t *testing.T) {
	suite.Run(t, new(AlbumCrudTestSuite))
}

func (a *AlbumCrudTestSuite) SetupSuite() {
	dyn := dynamotestutils.NewTestContext(context.Background(), a.T())
	err := appdynamodb.CreateTableIfNecessary(context.Background(), dyn.Table, dyn.Client, true)
	if !assert.NoError(a.T(), err) {
		assert.FailNow(a.T(), err.Error())
	}

	a.owner = "UNITTEST#1"
	a.repo = &Repository{
		client: dyn.Client,
		table:  dyn.Table,
	}
}

func (a *AlbumCrudTestSuite) TestInsertAndFind() {
	folderName := catalog.NewFolderName("/christmas")

	err := a.repo.InsertAlbum(context.TODO(), catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: folderName,
		},
		Name:  "Christmas",
		Start: mustParseDate("2020-12-24"),
		End:   mustParseDate("2020-12-26"),
	})
	if !a.NoError(err, "it should insert a new album in DB") {
		return
	}

	name := "it should find previously saved album"
	found, err := a.repo.FindAlbums(catalog.AlbumId{Owner: a.owner, FolderName: folderName})
	if a.NoError(err, name) && a.Len(found, 1, name) {
		a.Equal(&catalog.Album{
			AlbumId: catalog.AlbumId{
				Owner:      a.owner,
				FolderName: folderName,
			},
			Name:  "Christmas",
			Start: mustParseDate("2020-12-24"),
			End:   mustParseDate("2020-12-26"),
		}, found[0], name)
	}
}

func (a *AlbumCrudTestSuite) TestInsertTwiceFails() {
	folderName := catalog.NewFolderName("New Year")

	err := a.repo.InsertAlbum(context.TODO(), catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: folderName,
		},
		Name:  "New Year",
		Start: mustParseDate("2020-12-31"),
		End:   mustParseDate("2021-01-01"),
	})
	if !a.NoError(err, "it should insert a new album in DB") {
		return
	}

	err = a.repo.InsertAlbum(context.TODO(), catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: folderName,
		},
		Name:  "New Year Again",
		Start: mustParseDate("2020-12-31"),
		End:   mustParseDate("2021-01-01"),
	})
	log.WithField("Error", err).Infoln("insert twice fails")
	a.Error(err, "it should fail to override an existing album")
}

func (a *AlbumCrudTestSuite) TestFindNotFound() {
	ttName := "it should return [?, NotFoundError] when searched album do not exists"
	albums, err := a.repo.FindAlbums(catalog.AlbumId{Owner: a.owner, FolderName: "_donotexist"})
	if a.NoError(err, ttName) {
		a.Empty(albums)
	}
}

func (a *AlbumCrudTestSuite) TestDeleteEmpty() {
	folderName := catalog.NewFolderName("ToBeDeleted")
	albumId := catalog.AlbumId{
		Owner:      a.owner,
		FolderName: folderName,
	}

	err := a.repo.InsertAlbum(context.TODO(), catalog.Album{
		AlbumId: albumId,
		Name:    "ToBeDeleted",
		Start:   mustParseDate("2020-12-24"),
		End:     mustParseDate("2020-12-26"),
	})
	if !a.NoError(err, "it should insert an album to delete") {
		return
	}

	err = a.repo.DeleteEmptyAlbum(albumId)
	a.NoError(err, "it should delete an album that do not have any medias")
}

func (a *AlbumCrudTestSuite) TestUpdate() {
	folderName := "Update1"

	err := a.repo.InsertAlbum(context.TODO(), catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: catalog.NewFolderName("Update1"),
		},
		Name:  folderName,
		Start: mustParseDate("2020-12-01"),
		End:   mustParseDate("2021-01-31"),
	})
	if !a.NoError(err, "it should insert an album to update") {
		return
	}

	update := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: catalog.NewFolderName("Update1"),
		},
		Name:  "Another Name",
		Start: mustParseDate("2021-01-01"),
		End:   mustParseDate("2021-02-01"),
	}
	err = a.repo.UpdateAlbum(update)
	name := "it should update an exiting album"
	if a.NoError(err, name) {
		updated, err := a.repo.FindAlbums(catalog.AlbumId{Owner: a.owner, FolderName: catalog.NewFolderName(folderName)})
		if a.NoError(err, name) && a.Len(updated, 1) {
			a.Equal(&update, updated[0], name)
		}
	}
}

func (a *AlbumCrudTestSuite) TestUpdateNotExisting() {
	folderName := catalog.NewFolderName("_do_not_exist")

	update := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: folderName,
		},
		Name:  "Another Name",
		Start: mustParseDate("2021-01-01"),
		End:   mustParseDate("2021-02-01"),
	}
	err := a.repo.UpdateAlbum(update)
	a.Error(err, "it should fail to update an album that do not exist.")
}
