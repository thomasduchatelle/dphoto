package uploaders

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"
	"github.com/thomasduchatelle/dphoto/mocks"
	"sort"
	"testing"
	"time"
)

const isoDatePattern = "2006-01-02"
const owner = "unittest"

var mediaDate = time.Date(2021, 4, 27, 10, 16, 22, 0, time.UTC)

func TestUploader_Upload(t *testing.T) {
	a := assert.New(t)

	catalogProxy := new(mocks.CatalogProxyAdapter)
	onlineStorage := new(mocks.OnlineStorageAdapter)
	postFilter := new(mocks.PostAnalyseFilter)

	medias := []*backupmodel.AnalysedMedia{
		{
			FoundMedia: backupmodel.NewInmemoryMedia("image_001.jpg", 42, mediaDate),
			Type:       backupmodel.MediaTypeImage,
			Signature:  &backupmodel.FullMediaSignature{Sha256: "00000001", Size: 42},
			Details:    &backupmodel.MediaDetails{DateTime: mustParseDate("2021-03-27")},
		},
		{
			FoundMedia: backupmodel.NewInmemoryMedia("video_002.mkv", 4200, mediaDate),
			Type:       backupmodel.MediaTypeVideo,
			Signature:  &backupmodel.FullMediaSignature{Sha256: "00000002", Size: 4200},
			Details:    &backupmodel.MediaDetails{DateTime: mustParseDate("2021-04-02")},
		},
		{
			FoundMedia: backupmodel.NewInmemoryMedia("image_003.jpg", 42, mediaDate),
			Type:       backupmodel.MediaTypeImage,
			Signature:  &backupmodel.FullMediaSignature{Sha256: "00000003", Size: 42},
			Details:    &backupmodel.MediaDetails{DateTime: mustParseDate("2021-04-04")},
		},
		{
			FoundMedia: backupmodel.NewInmemoryMedia("image_004.jpg", 42, mediaDate),
			Type:       backupmodel.MediaTypeImage,
			Signature:  &backupmodel.FullMediaSignature{Sha256: "00000004", Size: 42},
			Details:    &backupmodel.MediaDetails{DateTime: mustParseDate("2021-04-05")},
		},
		{
			FoundMedia: backupmodel.NewInmemoryMedia("image_005.jpg", 42, mediaDate),
			Type:       backupmodel.MediaTypeImage,
			Signature:  &backupmodel.FullMediaSignature{Sha256: "00000005", Size: 42},
			Details:    &backupmodel.MediaDetails{DateTime: mustParseDate("2021-04-12")},
		},
		{
			FoundMedia: backupmodel.NewInmemoryMedia("image_006.jpg", 32, mediaDate),
			Type:       backupmodel.MediaTypeOther,
			Signature:  &backupmodel.FullMediaSignature{Sha256: "00000006", Size: 32},
			Details:    &backupmodel.MediaDetails{DateTime: mustParseDate("2021-04-13")},
		},
		{
			FoundMedia: backupmodel.NewInmemoryMedia("image_001_again.jpg", 42, mediaDate),
			Type:       backupmodel.MediaTypeImage,
			Signature:  &backupmodel.FullMediaSignature{Sha256: "00000001", Size: 42},
			Details:    &backupmodel.MediaDetails{DateTime: mustParseDate("2021-03-26")},
		},
	}

	catalogProxy.On("FindAllAlbums", owner).Return([]*catalog.Album{
		{owner, "Easter", "/2021-04_easter", mustParseDate("2021-04-04"), mustParseDate("2021-04-05")},
	}, nil)

	catalogProxy.On("Create", catalog.CreateAlbum{
		Owner:            owner,
		Name:             "Q1 2021",
		Start:            mustParseDate("2021-01-01"),
		End:              mustParseDate("2021-04-01"),
		ForcedFolderName: "/2021-Q1",
	}).Return(nil).Once()
	catalogProxy.On("Create", catalog.CreateAlbum{
		Owner:            owner,
		Name:             "Q2 2021",
		Start:            mustParseDate("2021-04-01"),
		End:              mustParseDate("2021-07-01"),
		ForcedFolderName: "/2021-Q2",
	}).Return(nil).Once()

	signatureRequest := make([]*catalog.MediaSignature, len(medias)-1) // dynamoDB do not support duplicates even in queries

	for i, sign := range medias {
		if i != 6 {
			signatureRequest[i] = &catalog.MediaSignature{SignatureSha256: sign.Signature.Sha256, SignatureSize: int(sign.Signature.Size)}
		}
	}
	catalogProxy.On("FindSignatures", owner, signatureRequest).Return([]*catalog.MediaSignature{signatureRequest[4]}, nil).Once()

	postFilter.On("AcceptAnalysedMedia", medias[5], "/2021-Q2").Return(false)
	postFilter.On("AcceptAnalysedMedia", mock.Anything, mock.Anything).Return(true)

	// EXPECTATION 1/2
	expectedCreateMediaRequest := []catalog.CreateMediaRequest{
		{
			Location: catalog.MediaLocation{
				FolderName: "/2021-Q1",
				Filename:   "2021-03-27_00-00-00_ONLINE.jpg",
			},
			Type:      "IMAGE",
			Details:   catalog.MediaDetails{DateTime: medias[0].Details.DateTime},
			Signature: *signatureRequest[0],
		},
		{
			Location: catalog.MediaLocation{
				FolderName: "/2021-Q2",
				Filename:   "2021-04-02_00-00-00_00000002.mkv",
			},
			Type:      "VIDEO",
			Details:   catalog.MediaDetails{DateTime: medias[1].Details.DateTime},
			Signature: *signatureRequest[1],
		},
		{
			Location: catalog.MediaLocation{
				FolderName: "/2021-04_easter",
				Filename:   "2021-04-04_00-00-00_00000003.jpg",
			},
			Type:      "IMAGE",
			Details:   catalog.MediaDetails{DateTime: medias[2].Details.DateTime},
			Signature: *signatureRequest[2],
		},
		{
			Location: catalog.MediaLocation{
				FolderName: "/2021-Q2",
				Filename:   "2021-04-05_00-00-00_00000004.jpg",
			},
			Type:      "IMAGE",
			Details:   catalog.MediaDetails{DateTime: medias[3].Details.DateTime},
			Signature: *signatureRequest[3],
		},
	}
	catalogProxy.On("InsertMedias", owner, mock.Anything).Return(func(owner string, actual []catalog.CreateMediaRequest) error {
		sort.Slice(actual, func(i, j int) bool {
			return actual[i].Location.Filename < actual[j].Location.Filename
		})
		sort.Slice(expectedCreateMediaRequest, func(i, j int) bool {
			return expectedCreateMediaRequest[i].Location.Filename < expectedCreateMediaRequest[j].Location.Filename
		})

		a.Equal(expectedCreateMediaRequest, actual, "InsertMedias should be called with the right list of medias, no matter the order.")

		return nil
	}).Once()

	// EXPECTATION 2/2
	onlineStorage.On("UploadFile", owner, mock.Anything, "/2021-Q1", "2021-03-27_00-00-00_00000001.jpg").Return("2021-03-27_00-00-00_ONLINE.jpg", nil).Once()
	onlineStorage.On("UploadFile", owner, mock.Anything, "/2021-Q2", "2021-04-02_00-00-00_00000002.mkv").Return("2021-04-02_00-00-00_00000002.mkv", nil).Once()
	onlineStorage.On("UploadFile", owner, mock.Anything, "/2021-04_easter", "2021-04-04_00-00-00_00000003.jpg").Return("2021-04-04_00-00-00_00000003.jpg", nil).Once()
	onlineStorage.On("UploadFile", owner, mock.Anything, "/2021-Q2", "2021-04-05_00-00-00_00000004.jpg").Return("2021-04-05_00-00-00_00000004.jpg", nil).Once()

	uploader, err := NewUploader(catalogProxy, onlineStorage, owner, postFilter)
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	err = uploader.Upload(medias, make(chan *backupmodel.ProgressEvent, 42))
	if a.NoError(err) {
		catalogProxy.AssertExpectations(t)
		onlineStorage.AssertExpectations(t)
	}
}

func mustParseDate(date string) time.Time {
	parse, err := time.Parse(isoDatePattern, date)
	if err != nil {
		panic(err)
	}

	return parse
}
