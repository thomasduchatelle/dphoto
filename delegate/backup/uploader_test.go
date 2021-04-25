package backup

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"duchatelle.io/dphoto/dphoto/catalog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sort"
	"testing"
	"time"
)

const isoDatePattern = "2006-01-02"

func TestUploader_Upload(t *testing.T) {
	a := assert.New(t)

	catalogProxy := new(MockCatalogProxyAdapter)
	onlineStorage := new(MockOnlineStorageAdapter)

	medias := []*model.AnalysedMedia{
		{
			FoundMedia: newInmemoryMedia("image_001.jpg", 42),
			Type:       model.MediaTypeImage,
			Signature:  &model.FullMediaSignature{Sha256: "00000001", Size: 42},
			Details:    &model.MediaDetails{DateTime: mustParseDate("2021-03-27")},
		},
		{
			FoundMedia: newInmemoryMedia("video_002.mkv", 4200),
			Type:       model.MediaTypeVideo,
			Signature:  &model.FullMediaSignature{Sha256: "00000002", Size: 4200},
			Details:    &model.MediaDetails{DateTime: mustParseDate("2021-04-02")},
		},
		{
			FoundMedia: newInmemoryMedia("image_003.jpg", 42),
			Type:       model.MediaTypeImage,
			Signature:  &model.FullMediaSignature{Sha256: "00000003", Size: 42},
			Details:    &model.MediaDetails{DateTime: mustParseDate("2021-04-04")},
		},
		{
			FoundMedia: newInmemoryMedia("image_004.jpg", 42),
			Type:       model.MediaTypeImage,
			Signature:  &model.FullMediaSignature{Sha256: "00000004", Size: 42},
			Details:    &model.MediaDetails{DateTime: mustParseDate("2021-04-05")},
		},
		{
			FoundMedia: newInmemoryMedia("image_005.jpg", 42),
			Type:       model.MediaTypeImage,
			Signature:  &model.FullMediaSignature{Sha256: "00000005", Size: 42},
			Details:    &model.MediaDetails{DateTime: mustParseDate("2021-04-12")},
		},
		{
			FoundMedia: newInmemoryMedia("image_001_again.jpg", 42),
			Type:       model.MediaTypeImage,
			Signature:  &model.FullMediaSignature{Sha256: "00000001", Size: 42},
			Details:    &model.MediaDetails{DateTime: mustParseDate("2021-03-26")},
		},
	}

	catalogProxy.On("FindAllAlbums").Return([]*catalog.Album{
		{"Easter", "2021-04_easter", mustParseDate("2021-04-04"), mustParseDate("2021-04-05")},
	}, nil)

	catalogProxy.On("Create", catalog.CreateAlbum{
		Name:             "Q1 2021",
		Start:            mustParseDate("2021-01-01"),
		End:              mustParseDate("2021-04-01"),
		ForcedFolderName: "2021-Q1",
	}).Return(nil).Once()
	catalogProxy.On("Create", catalog.CreateAlbum{
		Name:             "Q2 2021",
		Start:            mustParseDate("2021-04-01"),
		End:              mustParseDate("2021-07-01"),
		ForcedFolderName: "2021-Q2",
	}).Return(nil).Once()

	signatureRequest := make([]*catalog.MediaSignature, len(medias))

	for i, sign := range medias {
		signatureRequest[i] = &catalog.MediaSignature{SignatureSha256: sign.Signature.Sha256, SignatureSize: sign.Signature.Size}
	}
	catalogProxy.On("FindSignatures", signatureRequest).Return([]*catalog.MediaSignature{signatureRequest[4]}, nil).Once()

	// EXPECTATION 1/2
	expectedCreateMediaRequest := []catalog.CreateMediaRequest{
		{
			Location: catalog.MediaLocation{
				FolderName: "2021-Q1",
				Filename:   "2021-03-27_00-00-00_ONLINE.jpg",
			},
			Type:      "IMAGE",
			Details:   catalog.MediaDetails{DateTime: medias[0].Details.DateTime},
			Signature: *signatureRequest[0],
		},
		{
			Location: catalog.MediaLocation{
				FolderName: "2021-Q2",
				Filename:   "2021-04-02_00-00-00_00000002.mkv",
			},
			Type:      "VIDEO",
			Details:   catalog.MediaDetails{DateTime: medias[1].Details.DateTime},
			Signature: *signatureRequest[1],
		},
		{
			Location: catalog.MediaLocation{
				FolderName: "2021-04_easter",
				Filename:   "2021-04-04_00-00-00_00000003.jpg",
			},
			Type:      "IMAGE",
			Details:   catalog.MediaDetails{DateTime: medias[2].Details.DateTime},
			Signature: *signatureRequest[2],
		},
		{
			Location: catalog.MediaLocation{
				FolderName: "2021-Q2",
				Filename:   "2021-04-05_00-00-00_00000004.jpg",
			},
			Type:      "IMAGE",
			Details:   catalog.MediaDetails{DateTime: medias[3].Details.DateTime},
			Signature: *signatureRequest[3],
		},
	}
	catalogProxy.On("InsertMedias", mock.Anything).Return(func(actual []catalog.CreateMediaRequest) error {
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
	onlineStorage.On("UploadFile", mock.Anything, "2021-Q1", "2021-03-27_00-00-00_00000001.jpg").Return("2021-03-27_00-00-00_ONLINE.jpg", nil).Once()
	onlineStorage.On("UploadFile", mock.Anything, "2021-Q2", "2021-04-02_00-00-00_00000002.mkv").Return("2021-04-02_00-00-00_00000002.mkv", nil).Once()
	onlineStorage.On("UploadFile", mock.Anything, "2021-04_easter", "2021-04-04_00-00-00_00000003.jpg").Return("2021-04-04_00-00-00_00000003.jpg", nil).Once()
	onlineStorage.On("UploadFile", mock.Anything, "2021-Q2", "2021-04-05_00-00-00_00000004.jpg").Return("2021-04-05_00-00-00_00000004.jpg", nil).Once()

	uploader, err := NewUploader(catalogProxy, onlineStorage)
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	err = uploader.Upload(medias)
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
