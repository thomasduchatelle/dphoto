package dynamo

import (
	"github.com/thomasduchatelle/dphoto/dphoto/catalog"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
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
	repo   *Rep
	medias []catalog.CreateMediaRequest
	jan21  string
	feb21  string
	mar21  string
}

func TestRepositoryMediaCrud(t *testing.T) {
	suite.Run(t, new(MediaCrudTestSuite))
}

func (a *MediaCrudTestSuite) SetupSuite() {
	suffix := time.Now().Format("20060102150405")

	a.repo = &Rep{
		db:            dynamodb.New(session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-west-1")})), &aws.Config{Endpoint: aws.String("http://localhost:8000")}),
		table:         "test-medias-crud-" + suffix,
		RootOwner:     "UNITTEST#2",
		localDynamodb: true,
	}

	err := a.repo.CreateTableIfNecessary()
	if err != nil {
		panic(err)
	}

	err = a.preload()
	if err != nil {
		panic(err)
	}
}

func (a *MediaCrudTestSuite) preload() error {
	log.Infoln("Initialising dataset in dynamodb...")
	a.jan21 = "/media/2021-jan"
	a.feb21 = "/media/2021-feb"
	a.mar21 = "/media/2021-mar"

	err := a.repo.InsertAlbum(catalog.Album{
		Name:       "Media Container Jan",
		FolderName: a.jan21,
		Start:      mustParseDate("2021-01-01"),
		End:        mustParseDate("2021-02-01"),
	})
	if !a.NoError(err, "failed album initialisation") {
		return err
	}

	err = a.repo.InsertAlbum(catalog.Album{
		Name:       "Media Container Feb",
		FolderName: a.feb21,
		Start:      mustParseDate("2021-02-01"),
		End:        mustParseDate("2021-03-01"),
	})
	if !a.NoError(err, "failed album initialisation") {
		return err
	}

	err = a.repo.InsertAlbum(catalog.Album{
		Name:       "Media Container Mar",
		FolderName: a.mar21,
		Start:      mustParseDate("2021-03-01"),
		End:        mustParseDate("2021-04-01"),
	})
	if !a.NoError(err, "failed album initialisation") {
		return err
	}

	a.medias = []catalog.CreateMediaRequest{
		{
			Location: catalog.MediaLocation{
				FolderName: a.jan21,
				Filename:   "img001.jpeg",
			},
			Type: "Image",
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
			Signature: catalog.MediaSignature{
				SignatureSha256: "dc58865da1228b7a187693c702905d00d6a59439a07d52f2a8e7ae43764b55b9",
				SignatureSize:   16384,
			},
		},
		{
			Location: catalog.MediaLocation{
				FolderName: a.feb21,
				Filename:   "img002.jpeg",
			},
			Type: "Image",
			Details: catalog.MediaDetails{
				DateTime: time.Date(2021, 2, 20, 0, 0, 0, 0, time.UTC),
			},
			Signature: catalog.MediaSignature{
				SignatureSha256: "4d37f8780f5f5f14b914683b1fd36a9a567f5ea63a835b76100d9970303d6ad6",
				SignatureSize:   32000,
			},
		},
		{
			Location: catalog.MediaLocation{
				FolderName: a.jan21,
				Filename:   "img003.jpeg",
			},
			Type: "Image",
			Details: catalog.MediaDetails{
				DateTime: time.Date(2021, 1, 12, 0, 0, 0, 0, time.UTC),
			},
			Signature: catalog.MediaSignature{
				SignatureSha256: "77f218b4deaab40c47d21799f74a5c400b413d597e3f8926ef7d00572b8bb3d2",
				SignatureSize:   16384,
			},
		},
	}
	err = a.repo.InsertMedias(a.medias)
	a.NoError(err, "failed media initialisation")

	return err
}

func (a *MediaCrudTestSuite) fullPathNames(medias []*catalog.CreateMediaRequest) []string {
	names := make([]string, 0, len(medias))
	for _, a := range medias {
		names = append(names, path.Join(a.Location.FolderName, a.Location.Filename))
	}

	return names
}

func (a *MediaCrudTestSuite) TestFindAlbums() {
	albums, err := a.repo.FindAllAlbums()
	if a.NoError(err) {
		names := make([]string, 0, len(albums))
		for _, a := range albums {
			names = append(names, a.FolderName)
		}

		a.Equal([]string{"/media/2021-jan", "/media/2021-feb", "/media/2021-mar"}, names, "it should list all albums no matter how many medias are also stored")
	}
}

func (a *MediaCrudTestSuite) TestFindMedias() {
	allTime := catalog.TimeRange{}
	tests := []struct {
		name       string
		folderName string
		size       int64
		timeRange  catalog.TimeRange
		medias     [][]string
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
			[][]string{{"/media/2021-jan/img001.jpeg", "/media/2021-jan/img003.jpeg"}},
		},
		{
			"it should paginate with 1 item on each page",
			a.jan21,
			1,
			allTime,
			[][]string{{"/media/2021-jan/img001.jpeg"}, {"/media/2021-jan/img003.jpeg"}, {}},
		},
		{
			"it should paginate with 2 item on each page (last empty)",
			a.jan21,
			2,
			allTime,
			[][]string{{"/media/2021-jan/img001.jpeg", "/media/2021-jan/img003.jpeg"}, {}},
		},
		{
			"it should filter on the date to only get medias between 2 dates",
			a.jan21,
			42,
			newDateRange("2021-01-12", "2021-01-13"),
			[][]string{{"/media/2021-jan/img003.jpeg"}},
		},
	}

	for _, tt := range tests {
		var pages [][]string

		medias, err := a.repo.FindMedias(tt.folderName, catalog.FindMediaFilter{PageRequest: catalog.PageRequest{Size: tt.size}, TimeRange: tt.timeRange})
		if a.NoError(err, tt.name) {
			pages = append(pages, extractFilenames(tt.folderName, medias.Content))

			for medias.NextPage != "" {
				log.WithField("NextPage", medias.NextPage).Infoln("Request next page")
				medias, err = a.repo.FindMedias(tt.folderName, catalog.FindMediaFilter{PageRequest: catalog.PageRequest{Size: tt.size, NextPage: medias.NextPage}})
				if !a.NoError(err, tt.name) {
					return
				}
				pages = append(pages, extractFilenames(tt.folderName, medias.Content))
			}

			a.Equal(tt.medias, pages, tt.name)
		}
	}
}

func (a *MediaCrudTestSuite) TestFindMedias_AllDetails() {
	name := "it should find a media with all its details"
	medias, err := a.repo.FindMedias(a.jan21, catalog.FindMediaFilter{PageRequest: catalog.PageRequest{Size: 1}})
	if a.NoError(err, name) {
		a.Len(medias.Content, 1, name)
		a.Equal(&catalog.MediaMeta{
			Signature: a.medias[0].Signature,
			Filename:  a.medias[0].Location.Filename,
			Type:      a.medias[0].Type,
			Details:   a.medias[0].Details,
		}, medias.Content[0])
	}
}

func (a *MediaCrudTestSuite) TestDeleteNonEmpty() {
	err := a.repo.DeleteEmptyAlbum(a.jan21)
	a.Equal(catalog.NotEmptyError, err, "it should not delete an album with images in it")
}

func (a *MediaCrudTestSuite) TestFindExistingSignatures() {
	exiting := []*catalog.MediaSignature{
		{SignatureSha256: "4d37f8780f5f5f14b914683b1fd36a9a567f5ea63a835b76100d9970303d6ad6", SignatureSize: 32000},
		{SignatureSha256: "dc58865da1228b7a187693c702905d00d6a59439a07d52f2a8e7ae43764b55b9", SignatureSize: 16384},
	}
	search := make([]*catalog.MediaSignature, 0, dynamoReadBatchSize*2+20)
	for i := 0; i < dynamoReadBatchSize*2+20; i++ {
		search = append(search, &catalog.MediaSignature{
			SignatureSha256: fmt.Sprintf("%064d", i),
			SignatureSize:   42,
		})
	}

	signatures, err := a.repo.FindExistingSignatures(search)
	if a.NoError(err) {
		a.Empty(signatures, "it should not find any of non-existing signature")
	} else {
		return
	}

	search[42] = exiting[0]
	search[69] = exiting[1]
	signatures, err = a.repo.FindExistingSignatures(search)
	if a.NoError(err) {
		a.Equal(exiting, signatures, "it should filter out any non exiting signature to keep the only 2 that exist")
	}
}

func extractFilenames(albumFolderName string, medias []*catalog.MediaMeta) []string {
	filenames := make([]string, 0, len(medias))
	for _, m := range medias {
		filenames = append(filenames, path.Join(albumFolderName, m.Filename))
	}

	return filenames
}

func newDateRange(start, end string) catalog.TimeRange {
	return catalog.TimeRange{
		Start: mustParseDate(start),
		End:   mustParseDate(end),
	}
}
