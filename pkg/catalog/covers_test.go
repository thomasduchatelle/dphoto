package catalog_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

var deterministicRandomiser = catalog.RandomiserFunc(func(upperBound, n int) []int {
	if n > upperBound {
		n = upperBound
	}
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}
	return indices
})

var (
	coverAlbumId = catalog.AlbumId{
		Owner:      "ironman",
		FolderName: catalog.NewFolderName("/avengers"),
	}
	image1 = &catalog.MediaMeta{Id: "media-1", Filename: "photo-1.jpg", Type: catalog.MediaTypeImage}
	image2 = &catalog.MediaMeta{Id: "media-2", Filename: "photo-2.jpg", Type: catalog.MediaTypeImage}
	image3 = &catalog.MediaMeta{Id: "media-3", Filename: "photo-3.jpg", Type: catalog.MediaTypeImage}
	image4 = &catalog.MediaMeta{Id: "media-4", Filename: "photo-4.jpg", Type: catalog.MediaTypeImage}
	image5 = &catalog.MediaMeta{Id: "media-5", Filename: "photo-5.jpg", Type: catalog.MediaTypeImage}
	video1 = &catalog.MediaMeta{Id: "video-1", Filename: "clip-1.mp4", Type: catalog.MediaTypeVideo}
	other1 = &catalog.MediaMeta{Id: "other-1", Filename: "note-1.txt", Type: catalog.MediaTypeOther}
)

func TestCoverMaintenance_Reconcile(t *testing.T) {
	type fields struct {
		CoverRepository     *CoverRepositoryInMemory
		MediaReadRepository *MediaReadRepositoryInMemory
	}
	type args struct {
		albumId AlbumIdOrDefault
		added   []*catalog.MediaMeta
		removed []catalog.MediaId
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		expectSavedCovers   []catalog.Cover
		expectNotifications []coversChangedNotification
		wantErr             assert.ErrorAssertionFunc
	}{
		{
			name: "it should fill empty slots from the added candidates and notify observers",
			fields: fields{
				CoverRepository:     NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args: args{added: []*catalog.MediaMeta{image1, image2, image3, image4}},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			expectNotifications: []coversChangedNotification{{AlbumId: coverAlbumId, Covers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			}}},
			wantErr: assert.NoError,
		},
		{
			name: "it should cap at MaxCoversPerAlbum when more candidates are provided",
			fields: fields{
				CoverRepository:     NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args: args{added: []*catalog.MediaMeta{image1, image2, image3, image4, image5}},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			expectNotifications: []coversChangedNotification{{AlbumId: coverAlbumId, Covers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			}}},
			wantErr: assert.NoError,
		},
		{
			name: "it should drop existing RANDOM covers and redraw them from the added set",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(coverAlbumId,
					catalog.Cover{MediaId: "media-98", Filename: "stale.jpg", Origin: catalog.CoverOriginRandom},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args: args{added: []*catalog.MediaMeta{image1, image2, image3, image4}},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			expectNotifications: []coversChangedNotification{{AlbumId: coverAlbumId, Covers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			}}},
			wantErr: assert.NoError,
		},
		{
			name: "it should keep CHERRY_PICKED covers untouched",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(coverAlbumId,
					catalog.Cover{MediaId: "media-99", Filename: "starred.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-98", Filename: "random-existing.jpg", Origin: catalog.CoverOriginRandom},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args: args{added: []*catalog.MediaMeta{image1, image2, image3}},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-99", Filename: "starred.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
			},
			expectNotifications: []coversChangedNotification{{AlbumId: coverAlbumId, Covers: []catalog.Cover{
				{MediaId: "media-99", Filename: "starred.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
			}}},
			wantErr: assert.NoError,
		},
		{
			name: "it should strip a CHERRY_PICKED cover whose media is in the removed set",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(coverAlbumId,
					catalog.Cover{MediaId: "media-99", Filename: "starred.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-98", Filename: "stay.jpg", Origin: catalog.CoverOriginCherryPicked},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args: args{removed: []catalog.MediaId{"media-99"}},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-98", Filename: "stay.jpg", Origin: catalog.CoverOriginCherryPicked},
			},
			expectNotifications: []coversChangedNotification{{AlbumId: coverAlbumId, Covers: []catalog.Cover{
				{MediaId: "media-98", Filename: "stay.jpg", Origin: catalog.CoverOriginCherryPicked},
			}}},
			wantErr: assert.NoError,
		},
		{
			name: "it should fall back to the album's medias when the added set is empty and slots must be filled",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
						coverAlbumId: {image1, video1, image2},
					},
				},
			},
			args: args{added: nil},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
			},
			expectNotifications: []coversChangedNotification{{AlbumId: coverAlbumId, Covers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
			}}},
			wantErr: assert.NoError,
		},
		{
			name: "it should fall back to the album's medias when the added set is insufficient to fill the slots",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
						coverAlbumId: {image1, image2, image3, image4, image5},
					},
				},
			},
			args: args{added: []*catalog.MediaMeta{image1}},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			expectNotifications: []coversChangedNotification{{AlbumId: coverAlbumId, Covers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			}}},
			wantErr: assert.NoError,
		},
		{
			name: "it should skip videos and OTHER medias in both the added set and the fallback",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{coverAlbumId: {video1, other1}},
				},
			},
			args: args{added: []*catalog.MediaMeta{video1, other1, image1, image2}},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
			},
			expectNotifications: []coversChangedNotification{{AlbumId: coverAlbumId, Covers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
			}}},
			wantErr: assert.NoError,
		},
		{
			name: "it should be a no-op when the kept set is already full of CHERRY_PICKED covers",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(coverAlbumId,
					catalog.Cover{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args: args{added: []*catalog.MediaMeta{image1, image2}},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
			},
			expectNotifications: nil,
			wantErr:             assert.NoError,
		},
		{
			name: "it should produce an empty set and persist it when no eligible candidate is available",
			fields: fields{
				CoverRepository:     NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args:                args{added: []*catalog.MediaMeta{video1, other1}},
			expectSavedCovers:   nil,
			expectNotifications: nil,
			wantErr:             assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := &coversChangedObserverFake{}
			reconcile := &catalog.CoverMaintenance{
				CoverRepository:     tt.fields.CoverRepository,
				MediaReadRepository: tt.fields.MediaReadRepository,
				Randomiser:          deterministicRandomiser,
				Observers:           []catalog.CoversChangedObserver{observer},
			}

			albumId := tt.args.albumId.resolve()
			err := reconcile.Reconcile(context.Background(), albumId, tt.args.added, tt.args.removed)
			if !tt.wantErr(t, err, fmt.Sprintf("Reconcile(%v, added=%v, removed=%v)", albumId, tt.args.added, tt.args.removed)) {
				return
			}
			assert.Equal(t, tt.expectSavedCovers, tt.fields.CoverRepository.Covers[albumId], "covers stored for album")
			assert.Equal(t, tt.expectNotifications, observer.Notifications, "notifications fired to observers")
		})
	}
}

// AlbumIdOrDefault lets a case leave the args.albumId field at its zero value and
// resolve to the shared coverAlbumId when the test runs.
type AlbumIdOrDefault struct {
	albumId *catalog.AlbumId
}

func (a AlbumIdOrDefault) resolve() catalog.AlbumId {
	if a.albumId == nil {
		return coverAlbumId
	}
	return *a.albumId
}

type coversChangedNotification struct {
	AlbumId catalog.AlbumId
	Covers  []catalog.Cover
}

type coversChangedObserverFake struct {
	Notifications []coversChangedNotification
}

func (c *coversChangedObserverFake) OnCoversChanged(_ context.Context, albumId catalog.AlbumId, covers []catalog.Cover) error {
	c.Notifications = append(c.Notifications, coversChangedNotification{AlbumId: albumId, Covers: covers})
	return nil
}

func coversFor(albumId catalog.AlbumId, covers ...catalog.Cover) CoverRepositorySeed {
	return CoverRepositorySeed{AlbumId: albumId, Covers: covers}
}

// findAlbumByOwnerPortFake lets tests return a canned list of albums or an error for a
// specific owner, in a fully deterministic order.
type findAlbumByOwnerPortFake struct {
	AlbumsByOwner map[ownermodel.Owner][]*catalog.Album
	Err           error
}

func (f *findAlbumByOwnerPortFake) FindAlbumsByOwner(_ context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.AlbumsByOwner[owner], nil
}

// reconcileCoversPortFake records every (albumId, added, removed) tuple Reconcile is
// called with, and lets a case pre-register per-album errors to simulate partial
// failures.
type reconcileCoversPortFake struct {
	Reconciled []reconcileCoversCall
	Errors     map[catalog.AlbumId]error
}

type reconcileCoversCall struct {
	AlbumId catalog.AlbumId
	Added   []*catalog.MediaMeta
	Removed []catalog.MediaId
}

func (r *reconcileCoversPortFake) Reconcile(_ context.Context, albumId catalog.AlbumId, added []*catalog.MediaMeta, removed []catalog.MediaId) error {
	r.Reconciled = append(r.Reconciled, reconcileCoversCall{AlbumId: albumId, Added: added, Removed: removed})
	if err, ok := r.Errors[albumId]; ok {
		return err
	}
	return nil
}

func TestBackfillCovers_BackfillForOwner(t *testing.T) {
	owner := ownermodel.Owner("ironman")
	avengersId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avengers")}
	stealthId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/stealth")}
	avengers := &catalog.Album{AlbumId: avengersId, Name: "Avengers"}
	stealth := &catalog.Album{AlbumId: stealthId, Name: "Stealth"}

	listFailure := errors.New("list exploded")
	reconcileFailure := errors.New("reconcile exploded")

	ownerWithTwoAlbums := func() *findAlbumByOwnerPortFake {
		return &findAlbumByOwnerPortFake{AlbumsByOwner: map[ownermodel.Owner][]*catalog.Album{
			owner: {avengers, stealth},
		}}
	}

	type fields struct {
		FindAlbumByOwnerPort *findAlbumByOwnerPortFake
		ReconcileCoversPort  *reconcileCoversPortFake
	}
	type args struct {
		owner ownermodel.Owner
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		wantReport          catalog.BackfillReport
		wantErr             assert.ErrorAssertionFunc
		expectReconciledIds []catalog.AlbumId
	}{
		{
			name: "it should reconcile every album of the owner",
			fields: fields{
				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
				ReconcileCoversPort:  &reconcileCoversPortFake{},
			},
			args:                args{owner: owner},
			wantReport:          catalog.BackfillReport{Albums: 2},
			wantErr:             assert.NoError,
			expectReconciledIds: []catalog.AlbumId{avengersId, stealthId},
		},
		{
			name: "it should return an empty report when the owner has no album",
			fields: fields{
				FindAlbumByOwnerPort: &findAlbumByOwnerPortFake{},
				ReconcileCoversPort:  &reconcileCoversPortFake{},
			},
			args:                args{owner: owner},
			wantReport:          catalog.BackfillReport{Albums: 0},
			wantErr:             assert.NoError,
			expectReconciledIds: nil,
		},
		{
			name: "it should continue after a per-album failure and record it in the report",
			fields: fields{
				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
				ReconcileCoversPort: &reconcileCoversPortFake{
					Errors: map[catalog.AlbumId]error{avengersId: reconcileFailure},
				},
			},
			args: args{owner: owner},
			wantReport: catalog.BackfillReport{
				Albums:   2,
				Failures: []catalog.BackfillFailure{{AlbumId: avengersId, Err: reconcileFailure}},
			},
			wantErr:             assert.NoError,
			expectReconciledIds: []catalog.AlbumId{avengersId, stealthId},
		},
		{
			name: "it should return an error when listing the owner's albums fails",
			fields: fields{
				FindAlbumByOwnerPort: &findAlbumByOwnerPortFake{Err: listFailure},
				ReconcileCoversPort:  &reconcileCoversPortFake{},
			},
			args:       args{owner: owner},
			wantReport: catalog.BackfillReport{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, listFailure)
			},
			expectReconciledIds: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backfill := &catalog.BackfillCovers{
				FindAlbumByOwnerPort: tt.fields.FindAlbumByOwnerPort,
				ReconcileCoversPort:  tt.fields.ReconcileCoversPort,
			}

			report, err := backfill.BackfillForOwner(context.Background(), tt.args.owner)
			if !tt.wantErr(t, err, fmt.Sprintf("BackfillForOwner(%s)", tt.args.owner)) {
				return
			}
			assert.Equal(t, tt.wantReport, report, "backfill report")

			reconciledIds := make([]catalog.AlbumId, 0, len(tt.fields.ReconcileCoversPort.Reconciled))
			for _, call := range tt.fields.ReconcileCoversPort.Reconciled {
				reconciledIds = append(reconciledIds, call.AlbumId)
			}
			if tt.expectReconciledIds == nil {
				assert.Empty(t, reconciledIds, "albums reconciled")
			} else {
				assert.Equal(t, tt.expectReconciledIds, reconciledIds, "albums reconciled")
			}
		})
	}
}
