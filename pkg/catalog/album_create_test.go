package catalog_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

func TestCreateAlbum_Create(t *testing.T) {
	const owner = "tonystark"
	apr28 := time.Date(2024, 4, 28, 0, 0, 0, 0, time.UTC)
	may01 := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)

	createdAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      ownermodel.Owner(owner),
			FolderName: catalog.NewFolderName("/2024-04_Ironman_1"),
		},
		Name:  "Ironman 1",
		Start: apr28,
		End:   may01,
	}
	standardRequest := catalog.CreateAlbumRequest{
		Owner: ownermodel.Owner(owner),
		Name:  createdAlbum.Name,
		Start: createdAlbum.Start,
		End:   createdAlbum.End,
	}

	lifetimeAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      ownermodel.Owner(owner),
			FolderName: catalog.NewFolderName("/lifetime"),
		},
		Name:  "lifetime",
		Start: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	transferFromLifetime := catalog.MediaTransferRecords{
		createdAlbum.AlbumId: []catalog.MediaSelector{{
			FromAlbums: []catalog.AlbumId{lifetimeAlbum.AlbumId},
			Start:      apr28,
			End:        may01,
		}},
	}

	transferredFromLifetime := func() catalog.TransferredMedias {
		var ids []catalog.MediaId
		for day := apr28; day.Before(may01); day = day.AddDate(0, 0, 1) {
			ids = append(ids, fakeMediaId(lifetimeAlbum.AlbumId, day))
		}
		return catalog.TransferredMedias{
			Transfers:  map[catalog.AlbumId][]catalog.MediaId{createdAlbum.AlbumId: ids},
			FromAlbums: []catalog.AlbumId{lifetimeAlbum.AlbumId},
		}
	}

	mediasInCreatedAlbumAfterTransfer := func() []*catalog.MediaMeta {
		var medias []*catalog.MediaMeta
		for day := apr28; day.Before(may01); day = day.AddDate(0, 0, 1) {
			id := fakeMediaId(lifetimeAlbum.AlbumId, day)
			medias = append(medias, &catalog.MediaMeta{
				Id:       id,
				Filename: string(id) + ".jpg",
				Type:     catalog.MediaTypeImage,
			})
		}
		return medias
	}

	m1 := fakeMediaId(lifetimeAlbum.AlbumId, apr28)
	m2 := fakeMediaId(lifetimeAlbum.AlbumId, apr28.AddDate(0, 0, 1))
	m3 := fakeMediaId(lifetimeAlbum.AlbumId, apr28.AddDate(0, 0, 2))

	testError := errors.New("TEST error")

	repositoryWithLifetime := func() *AlbumRepositoryInMemory {
		freshAlbum := *lifetimeAlbum
		return NewAlbumRepositoryInMemory(&freshAlbum)
	}

	type fields struct {
		AlbumRepository     catalog.TimelineRepository
		CoverRepository     *CoverRepositoryInMemory
		MediaReadRepository *MediaReadRepositoryInMemory
	}
	type args struct {
		request catalog.CreateAlbumRequest
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		expectStoredAlbumIds  []catalog.AlbumId
		expectTransferRecords []catalog.MediaTransferRecords
		expectCreatedEvents   []catalog.AlbumCreated
		expectSavedCovers     map[catalog.AlbumId][]catalog.Cover
		wantErr               assert.ErrorAssertionFunc
	}{
		{
			name: "it should reject an invalid request (empty name) without touching the repository",
			fields: fields{
				AlbumRepository:     NewAlbumRepositoryInMemory(),
				CoverRepository:     NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: ownermodel.Owner(owner),
					Name:  "",
					Start: apr28,
					End:   may01,
				},
			},
			expectStoredAlbumIds: []catalog.AlbumId{},
			expectSavedCovers:    map[catalog.AlbumId][]catalog.Cover{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNameMandatoryErr, i...)
			},
		},
		{
			name: "it should reject a request whose forced folder name is already taken",
			fields: fields{
				AlbumRepository:     repositoryWithLifetime(),
				CoverRepository:     NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner:            ownermodel.Owner(owner),
					Name:             "A different name",
					Start:            apr28,
					End:              may01,
					ForcedFolderName: lifetimeAlbum.AlbumId.FolderName.String(),
				},
			},
			expectStoredAlbumIds: []catalog.AlbumId{lifetimeAlbum.AlbumId},
			expectSavedCovers:    map[catalog.AlbumId][]catalog.Cover{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumFolderNameAlreadyTakenErr, i...)
			},
		},
		{
			name: "it should insert the album without transferring any media when no other album overlaps, and not touch covers",
			fields: fields{
				AlbumRepository:     NewAlbumRepositoryInMemory(),
				CoverRepository:     NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args:                  args{request: standardRequest},
			expectStoredAlbumIds:  []catalog.AlbumId{createdAlbum.AlbumId},
			expectTransferRecords: []catalog.MediaTransferRecords{nil},
			expectCreatedEvents: []catalog.AlbumCreated{{
				CreatedAlbum:      createdAlbum,
				TransferredMedias: catalog.TransferredMedias{Transfers: map[catalog.AlbumId][]catalog.MediaId{}},
			}},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{},
			wantErr:           assert.NoError,
		},
		{
			name: "it should wire the CoverService into BulkCreateAlbum: inherit the source's CHERRY_PICKED cover into the new album, fill the rest from the transferred medias, and refill the source's stripped slot from its remaining medias",
			fields: fields{
				AlbumRepository: repositoryWithLifetime(),
				CoverRepository: NewCoverRepositoryInMemory(coversFor(lifetimeAlbum.AlbumId,
					catalog.Cover{MediaId: m1, Filename: string(m1) + ".jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "lifetime-survivor-a", Filename: "survivor-a.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "lifetime-survivor-b", Filename: "survivor-b.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "lifetime-survivor-c", Filename: "survivor-c.jpg", Origin: catalog.CoverOriginRandom},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
						createdAlbum.AlbumId: mediasInCreatedAlbumAfterTransfer(),
						lifetimeAlbum.AlbumId: {
							{Id: "lifetime-survivor-a", Filename: "survivor-a.jpg", Type: catalog.MediaTypeImage},
							{Id: "lifetime-survivor-b", Filename: "survivor-b.jpg", Type: catalog.MediaTypeImage},
							{Id: "lifetime-survivor-c", Filename: "survivor-c.jpg", Type: catalog.MediaTypeImage},
							{Id: "lifetime-backfill", Filename: "backfill.jpg", Type: catalog.MediaTypeImage},
						},
					},
				},
			},
			args:                  args{request: standardRequest},
			expectStoredAlbumIds:  []catalog.AlbumId{lifetimeAlbum.AlbumId, createdAlbum.AlbumId},
			expectTransferRecords: []catalog.MediaTransferRecords{transferFromLifetime},
			expectCreatedEvents: []catalog.AlbumCreated{{
				CreatedAlbum:      createdAlbum,
				TransferredMedias: transferredFromLifetime(),
				Covers: map[catalog.AlbumId][]catalog.Cover{
					createdAlbum.AlbumId: {
						{MediaId: m1, Filename: string(m1) + ".jpg", Origin: catalog.CoverOriginCherryPicked},
						{MediaId: m2, Filename: string(m2) + ".jpg", Origin: catalog.CoverOriginRandom},
						{MediaId: m3, Filename: string(m3) + ".jpg", Origin: catalog.CoverOriginRandom},
					},
					lifetimeAlbum.AlbumId: {
						{MediaId: "lifetime-survivor-a", Filename: "survivor-a.jpg", Origin: catalog.CoverOriginRandom},
						{MediaId: "lifetime-survivor-b", Filename: "survivor-b.jpg", Origin: catalog.CoverOriginRandom},
						{MediaId: "lifetime-survivor-c", Filename: "survivor-c.jpg", Origin: catalog.CoverOriginRandom},
						{MediaId: "lifetime-backfill", Filename: "backfill.jpg", Origin: catalog.CoverOriginRandom},
					},
				},
			}},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				createdAlbum.AlbumId: {
					{MediaId: m1, Filename: string(m1) + ".jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: m2, Filename: string(m2) + ".jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: m3, Filename: string(m3) + ".jpg", Origin: catalog.CoverOriginRandom},
				},
				lifetimeAlbum.AlbumId: {
					{MediaId: "lifetime-survivor-a", Filename: "survivor-a.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "lifetime-survivor-b", Filename: "survivor-b.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "lifetime-survivor-c", Filename: "survivor-c.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "lifetime-backfill", Filename: "backfill.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not insert any album, transfer any media, or fire any event if the list of existing albums cannot be read",
			fields: fields{
				AlbumRepository:     failingLoadTimeline(repositoryWithLifetime(), testError),
				CoverRepository:     NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args:                 args{request: standardRequest},
			expectStoredAlbumIds: []catalog.AlbumId{lifetimeAlbum.AlbumId},
			expectSavedCovers:    map[catalog.AlbumId][]catalog.Cover{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testError, i...)
			},
		},
		{
			name: "it should not transfer any media or fire any event if the album insertion fails",
			fields: fields{
				AlbumRepository:     failingInsertAlbumForCreate(repositoryWithLifetime(), testError),
				CoverRepository:     NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args:                 args{request: standardRequest},
			expectStoredAlbumIds: []catalog.AlbumId{lifetimeAlbum.AlbumId},
			expectSavedCovers:    map[catalog.AlbumId][]catalog.Cover{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testError, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transferService := &TransferMediasServiceFake{}
			observer := &AlbumCreatedObserverInMemory{}
			coverService := &catalog.CoverService{
				CoverRepository:     tt.fields.CoverRepository,
				MediaReadRepository: tt.fields.MediaReadRepository,
				Randomiser:          deterministicRandomiser,
			}

			createAlbum := catalog.NewAlbumCreate(
				tt.fields.AlbumRepository,
				transferService,
				coverService,
				observer,
			)

			_, err := createAlbum.Create(context.Background(), tt.args.request)
			if !tt.wantErr(t, err, fmt.Sprintf("Create(%v)", tt.args.request)) {
				return
			}

			repository := underlyingCreateRepository(tt.fields.AlbumRepository)
			storedIds := make([]catalog.AlbumId, 0, len(repository.Albums))
			for id := range repository.Albums {
				storedIds = append(storedIds, id)
			}
			assert.ElementsMatch(t, tt.expectStoredAlbumIds, storedIds, "stored albums")
			assert.Equal(t, tt.expectTransferRecords, transferService.Records, "records passed to TransferMedias")
			assert.Equal(t, tt.expectCreatedEvents, observer.Events, "AlbumCreated events fired")
			assert.Equal(t, tt.expectSavedCovers, tt.fields.CoverRepository.Covers, "covers persisted")
		})
	}
}

func failingInsertAlbumForCreate(memory *AlbumRepositoryInMemory, err error) catalog.TimelineRepository {
	return &failingInsertAlbumForCreateInterceptor{
		AlbumRepositoryInMemory: memory,
		err:                     err,
	}
}

type failingInsertAlbumForCreateInterceptor struct {
	*AlbumRepositoryInMemory
	err error
}

func (i *failingInsertAlbumForCreateInterceptor) InsertAlbum(_ context.Context, _ catalog.Album) error {
	return i.err
}

func failingLoadTimeline(memory *AlbumRepositoryInMemory, err error) catalog.TimelineRepository {
	return &failingLoadTimelineInterceptor{
		AlbumRepositoryInMemory: memory,
		err:                     err,
	}
}

type failingLoadTimelineInterceptor struct {
	*AlbumRepositoryInMemory
	err error
}

func (i *failingLoadTimelineInterceptor) LoadTimeline(_ context.Context, _ ownermodel.Owner) (*catalog.TimelineAggregate, error) {
	return nil, i.err
}

func underlyingCreateRepository(port catalog.TimelineRepository) *AlbumRepositoryInMemory {
	switch v := port.(type) {
	case *AlbumRepositoryInMemory:
		return v
	case *failingInsertAlbumForCreateInterceptor:
		return v.AlbumRepositoryInMemory
	case *failingLoadTimelineInterceptor:
		return v.AlbumRepositoryInMemory
	}
	return nil
}
