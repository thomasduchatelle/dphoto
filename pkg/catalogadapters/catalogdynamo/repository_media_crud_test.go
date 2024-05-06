package catalogdynamo

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"path"
	"testing"
	"time"
)

const IsoDate = "2006-01-02"

func mustParseDate(date string) time.Time {
	parse, err := time.Parse(IsoDate, date)
	if err != nil {
		panic(err)
	}

	return parse
}

type MediaCrudTestSuite struct {
	suite.Suite
	owner  catalog.Owner
	repo   *Repository
	medias []catalog.CreateMediaRequest
	jan21  catalog.FolderName
	feb21  catalog.FolderName
	mar21  catalog.FolderName
}

func TestRepositoryMediaCrud(t *testing.T) {
	suite.Run(t, new(MediaCrudTestSuite))
}

func (a *MediaCrudTestSuite) SetupSuite() {
	dyn := dynamotestutils.NewTestContext(context.Background(), a.T())
	err := appdynamodb.CreateTableIfNecessary(context.Background(), dyn.Table, dyn.Client, true)
	if !assert.NoError(a.T(), err) {
		assert.FailNow(a.T(), err.Error())
	}

	a.owner = "UNITTEST#2"
	a.repo = &Repository{
		client: dyn.Client,
		table:  dyn.Table,
	}

	err = a.preload()
	if err != nil {
		panic(err)
	}
}

func (a *MediaCrudTestSuite) preload() error {
	log.Infoln("Initialising dataset in dynamodb...")
	a.jan21 = catalog.NewFolderName("/media/2021-jan")
	a.feb21 = catalog.NewFolderName("/media/2021-feb")
	a.mar21 = catalog.NewFolderName("/media/2021-mar")

	err := a.repo.InsertAlbum(context.TODO(), catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: a.jan21,
		},
		Name:  "Media Container Jan",
		Start: mustParseDate("2021-01-01"),
		End:   mustParseDate("2021-02-01"),
	})
	if !a.NoError(err, "failed album initialisation") {
		return err
	}

	err = a.repo.InsertAlbum(context.TODO(), catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: a.feb21,
		},
		Name:  "Media Container Feb",
		Start: mustParseDate("2021-02-01"),
		End:   mustParseDate("2021-03-01"),
	})
	if !a.NoError(err, "failed album initialisation") {
		return err
	}

	err = a.repo.InsertAlbum(context.TODO(), catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: a.mar21,
		},
		Name:  "Media Container Mar",
		Start: mustParseDate("2021-03-01"),
		End:   mustParseDate("2021-04-01"),
	})
	if !a.NoError(err, "failed album initialisation") {
		return err
	}

	img001Signature := catalog.MediaSignature{
		SignatureSha256: "dc58865da1228b7a187693c702905d00d6a59439a07d52f2a8e7ae43764b55b9",
		SignatureSize:   16384,
	}
	img002Signature := catalog.MediaSignature{
		SignatureSha256: "4d37f8780f5f5f14b914683b1fd36a9a567f5ea63a835b76100d9970303d6ad6",
		SignatureSize:   32000,
	}
	img003Signature := catalog.MediaSignature{
		SignatureSha256: "77f218b4deaab40c47d21799f74a5c400b413d597e3f8926ef7d00572b8bb3d2",
		SignatureSize:   16384,
	}
	a.medias = []catalog.CreateMediaRequest{
		{
			Id:         mustGenerateMediaId(catalog.GenerateMediaId(img001Signature)),
			Signature:  img001Signature,
			FolderName: a.jan21,
			Filename:   "img001.jpeg",
			Type:       "Image",
			Details: catalog.MediaDetails{
				Width:        1280,
				Height:       720,
				DateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Orientation:  "TopLeft",
				Make:         "Google",
				Model:        "Pixel",
				GPSLatitude:  0.123,
				GPSLongitude: 0.456,
			},
		},
		{
			Id:         mustGenerateMediaId(catalog.GenerateMediaId(img002Signature)),
			Signature:  img002Signature,
			FolderName: a.feb21,
			Filename:   "img002.jpeg",
			Type:       "Image",
			Details: catalog.MediaDetails{
				DateTime: time.Date(2021, 2, 20, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			Id:         mustGenerateMediaId(catalog.GenerateMediaId(img003Signature)),
			Signature:  img003Signature,
			FolderName: a.jan21,
			Filename:   "img003.jpeg",
			Type:       "Image",
			Details: catalog.MediaDetails{
				DateTime: time.Date(2021, 1, 12, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	err = a.repo.InsertMedias(context.TODO(), a.owner, a.medias)
	a.NoError(err, "failed media initialisation")

	return err
}

func mustGenerateMediaId(id string, err error) catalog.MediaId {
	if err != nil {
		panic(err)
	}
	return catalog.MediaId(id)
}

func (a *MediaCrudTestSuite) fullPathNames(medias []*catalog.CreateMediaRequest) []string {
	names := make([]string, 0, len(medias))
	for _, a := range medias {
		names = append(names, path.Join(a.FolderName.String(), a.Filename))
	}

	return names
}

func (a *MediaCrudTestSuite) TestFindAlbums() {
	albums, err := a.repo.FindAlbumsByOwner(context.Background(), a.owner)
	if a.NoError(err) {
		names := make(map[catalog.FolderName]int)
		for _, a := range albums {
			names[a.FolderName] = a.TotalCount
		}

		a.Equal(map[catalog.FolderName]int{
			a.jan21: 2,
			a.feb21: 1,
			a.mar21: 0},
			names,
			"it should list all albums no matter how many medias are also stored",
		)
	}
}

func (a *MediaCrudTestSuite) TestFindMedias() {
	allTime := catalog.TimeRange{}
	tests := []struct {
		name       string
		folderName catalog.FolderName
		size       int64
		timeRange  catalog.TimeRange
		medias     [][]string // medias is a slice of slice to represent pages (that have been removed)
	}{
		{
			"it should find no media in empty albums",
			a.mar21,
			0,
			allTime,
			[][]string{{}},
		},
		{
			"it should find 2 medias in Jan",
			a.jan21,
			0,
			allTime,
			[][]string{{"/media_2021-jan/img001.jpeg", "/media_2021-jan/img003.jpeg"}},
		},
		{
			"it should filter on the date to only get medias between 2 dates",
			a.jan21,
			42,
			newDateRange("2021-01-12", "2021-01-13"),
			[][]string{{"/media_2021-jan/img003.jpeg"}},
		},
	}

	for _, tt := range tests {
		var pages [][]string

		medias, err := a.repo.FindMedias(context.TODO(), catalog.NewFindMediaRequest(a.owner).WithAlbum(tt.folderName).WithinRange(tt.timeRange.Start, tt.timeRange.End))
		if a.NoError(err, tt.name) {
			pages = append(pages, extractFilenames(tt.folderName, medias))

			a.Equal(tt.medias, pages, tt.name)
		}
	}
}

func (a *MediaCrudTestSuite) TestFindMediaIds() {
	allTime := catalog.TimeRange{}
	tests := []struct {
		name       string
		folderName catalog.FolderName
		timeRange  catalog.TimeRange
		medias     [][]catalog.MediaId // medias is a slice of slice to represent pages (that have been removed)
	}{
		{
			"it should find no media in empty albums",
			a.mar21,
			allTime,
			[][]catalog.MediaId{nil},
		},
		{
			"it should find 2 medias in Jan",
			a.jan21,
			allTime,
			[][]catalog.MediaId{{a.medias[0].Id, a.medias[2].Id}},
		},
		{
			"it should filter on the date to only get medias between 2 dates",
			a.jan21,
			newDateRange("2021-01-12", "2021-01-13"),
			[][]catalog.MediaId{{a.medias[2].Id}},
		},
	}

	for _, tt := range tests {
		var pages [][]catalog.MediaId

		ids, err := a.repo.FindMediaIds(context.TODO(), catalog.NewFindMediaRequest(a.owner).WithAlbum(tt.folderName).WithinRange(tt.timeRange.Start, tt.timeRange.End))
		if a.NoError(err, tt.name) {
			pages = append(pages, ids)

			a.Equal(tt.medias, pages, tt.name)
		}
	}
}

func (a *MediaCrudTestSuite) TestFindMedias_AllDetails() {
	name := "it should find a media with all its details"
	medias, err := a.repo.FindMedias(context.TODO(), catalog.NewFindMediaRequest(a.owner).WithAlbum(a.jan21))
	if a.NoError(err, name) {
		a.Len(extractFilenames(a.jan21, medias), 2, name)
		a.Equal(&catalog.MediaMeta{
			Id:        a.medias[0].Id,
			Signature: a.medias[0].Signature,
			Filename:  a.medias[0].Filename,
			Type:      a.medias[0].Type,
			Details:   a.medias[0].Details,
		}, medias[0])
	}
}

func (a *MediaCrudTestSuite) TestDeleteNonEmpty() {
	err := a.repo.DeleteEmptyAlbum(context.TODO(), catalog.AlbumId{Owner: a.owner, FolderName: a.jan21})
	a.Equal(catalog.NotEmptyError, err, "it should not delete an album with images in it")
}

func (a *MediaCrudTestSuite) TestFindExistingSignatures() {
	exiting := []*catalog.MediaSignature{
		{SignatureSha256: "dc58865da1228b7a187693c702905d00d6a59439a07d52f2a8e7ae43764b55b9", SignatureSize: 16384},
		{SignatureSha256: "4d37f8780f5f5f14b914683b1fd36a9a567f5ea63a835b76100d9970303d6ad6", SignatureSize: 32000},
	}
	search := make([]*catalog.MediaSignature, 0, dynamoutils.DynamoReadBatchSize*2+20)
	for i := 0; i < dynamoutils.DynamoReadBatchSize*2+20; i++ {
		search = append(search, &catalog.MediaSignature{
			SignatureSha256: fmt.Sprintf("%064d", i),
			SignatureSize:   42,
		})
	}

	signatures, err := a.repo.FindExistingSignatures(context.TODO(), a.owner, search)
	if a.NoError(err) {
		a.Empty(signatures, "it should not find any of non-existing signature")
	} else {
		return
	}

	search[42] = exiting[0]
	search[69] = exiting[1]
	signatures, err = a.repo.FindExistingSignatures(context.TODO(), a.owner, search)
	if a.NoError(err) {
		a.Equal(exiting, signatures, "it should filter out any non exiting signature to keep the only 2 that exist")
	}

	search[1] = exiting[1]
	signatures, err = a.repo.FindExistingSignatures(context.TODO(), a.owner, search)
	if a.NoError(err, "it should ignore duplicated keys") {
		a.Equal(exiting, signatures, "it should ignore duplicated keys")
	}
}

func (a *MediaCrudTestSuite) TestFindMediaCurrentAlbum() {
	type args struct {
		owner   catalog.Owner
		mediaId catalog.MediaId
	}
	tests := []struct {
		name    string
		args    args
		want    catalog.FolderName
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "it should find the folder name of a media (sinple case)",
			args:    args{owner: a.owner, mediaId: a.medias[0].Id},
			want:    a.jan21,
			wantErr: assert.NoError,
		},
		{
			name:    "it should find the folder name even when there is a lot of # everywhere",
			args:    args{owner: "this#is#my#owner", mediaId: "a#random#id"},
			want:    catalog.NewFolderName("this#is#my#folder"),
			wantErr: assert.NoError,
		},
		{
			name: "it should return NotFound when the media doesn't ",
			args: args{owner: "this#is#my#owner", mediaId: "foobar"},
			want: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.NotFoundError)
			},
		},
	}

	err := a.repo.InsertAlbum(context.TODO(), catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      "this#is#my#owner",
			FolderName: catalog.NewFolderName("this#is#my#folder"),
		},
		Name:  "Media Container Jan",
		Start: mustParseDate("2021-01-01"),
		End:   mustParseDate("2021-02-01"),
	})
	if !assert.NoError(a.T(), err, "failed album initialisation") {
		assert.FailNow(a.T(), err.Error())
	}
	err = a.repo.InsertMedias(context.TODO(), "this#is#my#owner", []catalog.CreateMediaRequest{
		{
			Id:         "a#random#id",
			Signature:  catalog.MediaSignature{SignatureSha256: "qwertyuiop", SignatureSize: 42},
			FolderName: "this#is#my#folder",
			Filename:   "img001.jpeg",
			Type:       "Image",
			Details: catalog.MediaDetails{
				Width:        1280,
				Height:       720,
				DateTime:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Orientation:  "TopLeft",
				Make:         "Google",
				Model:        "Pixel",
				GPSLatitude:  0.123,
				GPSLongitude: 0.456,
			},
		},
	})
	if !assert.NoError(a.T(), err) {
		assert.FailNow(a.T(), err.Error())
	}

	for _, tt := range tests {
		a.T().Run(tt.name, func(t *testing.T) {
			got, err := a.repo.FindMediaCurrentAlbum(context.TODO(), tt.args.owner, tt.args.mediaId)
			if tt.wantErr(t, err) && err == nil {
				assert.Equal(t, &catalog.AlbumId{Owner: tt.args.owner, FolderName: tt.want}, got)
			}
		})
	}
}

func (a *MediaCrudTestSuite) TestTransferMedias() {
	name := "it should find transferred media in the new album and not anymore on the previous album"

	err := a.repo.TransferMedias(context.TODO(), a.owner, []catalog.MediaId{a.medias[0].Id, a.medias[1].Id}, a.mar21)
	if a.NoError(err, name) {
		mediasInMar21, err := a.repo.FindMediaIds(context.TODO(), catalog.NewFindMediaRequest(a.owner).WithAlbum(a.mar21))
		if a.NoError(err, name) {
			a.Equal([]catalog.MediaId{a.medias[0].Id, a.medias[1].Id}, mediasInMar21, name)
		}

		mediasInJan21, err := a.repo.FindMediaIds(context.TODO(), catalog.NewFindMediaRequest(a.owner).WithAlbum(a.jan21))
		if a.NoError(err, name) {
			a.Equal([]catalog.MediaId{a.medias[2].Id}, mediasInJan21, name)
		}

		mediasInFeb21, err := a.repo.FindMediaIds(context.TODO(), catalog.NewFindMediaRequest(a.owner).WithAlbum(a.feb21))
		if a.NoError(err, name) {
			a.Nil(mediasInFeb21, name)
		}
	}
}

func extractFilenames(albumFolderName catalog.FolderName, medias []*catalog.MediaMeta) []string {
	filenames := make([]string, 0, len(medias))
	for _, m := range medias {
		filenames = append(filenames, path.Join(albumFolderName.String(), m.Filename))
	}

	return filenames
}

func newDateRange(start, end string) catalog.TimeRange {
	return catalog.TimeRange{
		Start: mustParseDate(start),
		End:   mustParseDate(end),
	}
}
