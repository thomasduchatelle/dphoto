package dynamo

import (
	"duchatelle.io/dphoto/dphoto/album"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type AlbumCrudTestSuite struct {
	suite.Suite
	suffix string
	repo   *rep
}

func TestRepositoryAlbum(t *testing.T) {
	suite.Run(t, new(AlbumCrudTestSuite))
}

func (a *AlbumCrudTestSuite) SetupSuite() {
	a.suffix = time.Now().Format("20060102150405")

	a.repo = &rep{
		db:            dynamodb.New(session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-west-1")})), &aws.Config{Endpoint: aws.String("http://localhost:8000")}),
		table:         "test-albums-" + a.suffix,
		RootOwner:     "UNITTEST#1",
		localDynamodb: true,
	}

	err := a.repo.CreateTableIfNecessary()
	if err != nil {
		panic(err)
	}
}

func (a *AlbumCrudTestSuite) TestInsertAndFind() {
	folderName := "Christmas"

	err := a.repo.InsertAlbum(album.Album{
		Name:       "Christmas",
		FolderName: folderName,
		Start:      mustParseDate("2020-12-24"),
		End:        mustParseDate("2020-12-26"),
	})
	if !a.NoError(err, "it should insert a new album in DB") {
		return
	}

	name := "it should find previously saved album"
	found, err := a.repo.FindAlbum(folderName)
	if a.NoError(err, name) {
		a.Equal(&album.Album{
			Name:       "Christmas",
			FolderName: folderName,
			Start:      mustParseDate("2020-12-24"),
			End:        mustParseDate("2020-12-26"),
		}, found, name)
	}
}

func (a *AlbumCrudTestSuite) TestInsertTwiceFails() {
	folderName := "New Year"

	err := a.repo.InsertAlbum(album.Album{
		Name:       "New Year",
		FolderName: folderName,
		Start:      mustParseDate("2020-12-31"),
		End:        mustParseDate("2021-01-01"),
	})
	if !a.NoError(err, "it should insert a new album in DB") {
		return
	}

	err = a.repo.InsertAlbum(album.Album{
		Name:       "New Year Again",
		FolderName: folderName,
		Start:      mustParseDate("2020-12-31"),
		End:        mustParseDate("2021-01-01"),
	})
	log.WithField("Error", err).Infoln("insert twice fails")
	a.Error(err, "it should fail to override an existing album")
}

func (a *AlbumCrudTestSuite) TestFindNotFound() {
	ttName := "it should return [?, NotFoundError] when searched album do not exists"
	_, err := a.repo.FindAlbum("_donotexist")
	if a.Error(err, ttName) {
		a.Equal(album.NotFoundError, err)
	}
}

func (a *AlbumCrudTestSuite) TestDeleteEmpty() {
	folderName := "ToBeDeleted"
	err := a.repo.InsertAlbum(album.Album{
		Name:       folderName,
		FolderName: folderName,
		Start:      mustParseDate("2020-12-24"),
		End:        mustParseDate("2020-12-26"),
	})
	if !a.NoError(err, "it should insert an album to delete") {
		return
	}

	err = a.repo.DeleteEmptyAlbum(folderName)
	a.NoError(err, "it should delete an album that do not have any medias")
}

func (a *AlbumCrudTestSuite) TestUpdate() {
	folderName := "Update1"

	err := a.repo.InsertAlbum(album.Album{
		Name:       folderName,
		FolderName: folderName,
		Start:      mustParseDate("2020-12-01"),
		End:        mustParseDate("2021-01-31"),
	})
	if !a.NoError(err, "it should insert an album to update") {
		return
	}

	update := album.Album{
		Name:       "Another Name",
		FolderName: folderName,
		Start:      mustParseDate("2021-01-01"),
		End:        mustParseDate("2021-02-01"),
	}
	err = a.repo.UpdateAlbum(update)
	name := "it should update an exiting album"
	if a.NoError(err, name) {
		updated, err := a.repo.FindAlbum(folderName)
		if a.NoError(err, name) {
			a.Equal(&update, updated, name)
		}
	}
}

func (a *AlbumCrudTestSuite) TestUpdateNotExisting() {
	folderName := "_do_not_exist"
	update := album.Album{
		Name:       "Another Name",
		FolderName: folderName,
		Start:      mustParseDate("2021-01-01"),
		End:        mustParseDate("2021-02-01"),
	}
	err := a.repo.UpdateAlbum(update)
	a.Error(err, "it should fail to update an album that do not exist.")
}
