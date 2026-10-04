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

func TestBackfillCovers_BackfillForOwner(t *testing.T) {
	owner := ownermodel.Owner("ironman")
	avengersAlbum := &catalog.Album{AlbumId: avengersId, Name: "Avengers"}
	stealthAlbum := &catalog.Album{AlbumId: stealthId, Name: "Stealth"}
	avengersCovers := []catalog.Cover{{MediaId: "m-a", Filename: "a.jpg", Origin: catalog.CoverOriginRandom}}
	stealthCovers := []catalog.Cover{{MediaId: "m-s", Filename: "s.jpg", Origin: catalog.CoverOriginRandom}}

	listFailure := errors.New("list exploded")
	randomiseFailure := errors.New("randomise exploded")
	observerFailure := errors.New("observer exploded")

	type fields struct {
		FindAlbumByOwnerPort *findAlbumByOwnerPortFake
		CoverService         *coverServicePortFake
		Observer             *coverBackfillObserverFake
	}
	type args struct {
		owner ownermodel.Owner
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		want                map[catalog.AlbumId][]catalog.Cover
		wantErr             assert.ErrorAssertionFunc
		expectRandomiseCall *randomiseCall
		expectObserverSeen  []map[catalog.AlbumId][]catalog.Cover
	}{
		{
			name: "it should call Randomise(stable=true) with every album of the owner and return the resulting map",
			fields: fields{
				FindAlbumByOwnerPort: &findAlbumByOwnerPortFake{AlbumsByOwner: map[ownermodel.Owner][]*catalog.Album{owner: {avengersAlbum, stealthAlbum}}},
				CoverService: &coverServicePortFake{
					RandomiseResult: map[catalog.AlbumId][]catalog.Cover{avengersId: avengersCovers, stealthId: stealthCovers},
				},
				Observer: &coverBackfillObserverFake{},
			},
			args:                args{owner: owner},
			want:                map[catalog.AlbumId][]catalog.Cover{avengersId: avengersCovers, stealthId: stealthCovers},
			wantErr:             assert.NoError,
			expectRandomiseCall: &randomiseCall{Stable: true, AlbumIds: []catalog.AlbumId{avengersId, stealthId}},
			expectObserverSeen: []map[catalog.AlbumId][]catalog.Cover{
				{avengersId: avengersCovers, stealthId: stealthCovers},
			},
		},
		{
			name: "it should notify every observer with the empty map when no cover was changed",
			fields: fields{
				FindAlbumByOwnerPort: &findAlbumByOwnerPortFake{AlbumsByOwner: map[ownermodel.Owner][]*catalog.Album{owner: {avengersAlbum}}},
				CoverService:         &coverServicePortFake{},
				Observer:             &coverBackfillObserverFake{},
			},
			args:                args{owner: owner},
			want:                nil,
			wantErr:             assert.NoError,
			expectRandomiseCall: &randomiseCall{Stable: true, AlbumIds: []catalog.AlbumId{avengersId}},
			expectObserverSeen:  []map[catalog.AlbumId][]catalog.Cover{nil},
		},
		{
			name: "it should return a nil map when the owner has no album",
			fields: fields{
				FindAlbumByOwnerPort: &findAlbumByOwnerPortFake{},
				CoverService:         &coverServicePortFake{},
				Observer:             &coverBackfillObserverFake{},
			},
			args:                args{owner: owner},
			want:                nil,
			wantErr:             assert.NoError,
			expectRandomiseCall: &randomiseCall{Stable: true, AlbumIds: []catalog.AlbumId{}},
			expectObserverSeen:  []map[catalog.AlbumId][]catalog.Cover{nil},
		},
		{
			name: "it should wrap and return the error when FindAlbumsByOwner fails",
			fields: fields{
				FindAlbumByOwnerPort: &findAlbumByOwnerPortFake{Err: listFailure},
				CoverService:         &coverServicePortFake{},
				Observer:             &coverBackfillObserverFake{},
			},
			args: args{owner: owner},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, listFailure)
			},
			expectRandomiseCall: nil,
			expectObserverSeen:  nil,
		},
		{
			name: "it should wrap and return the error when Randomise fails and skip the observers",
			fields: fields{
				FindAlbumByOwnerPort: &findAlbumByOwnerPortFake{AlbumsByOwner: map[ownermodel.Owner][]*catalog.Album{owner: {avengersAlbum}}},
				CoverService:         &coverServicePortFake{Err: randomiseFailure},
				Observer:             &coverBackfillObserverFake{},
			},
			args: args{owner: owner},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, randomiseFailure)
			},
			expectRandomiseCall: &randomiseCall{Stable: true, AlbumIds: []catalog.AlbumId{avengersId}},
			expectObserverSeen:  nil,
		},
		{
			name: "it should return the error when an observer fails",
			fields: fields{
				FindAlbumByOwnerPort: &findAlbumByOwnerPortFake{AlbumsByOwner: map[ownermodel.Owner][]*catalog.Album{owner: {avengersAlbum}}},
				CoverService: &coverServicePortFake{
					RandomiseResult: map[catalog.AlbumId][]catalog.Cover{avengersId: avengersCovers},
				},
				Observer: &coverBackfillObserverFake{Err: observerFailure},
			},
			args: args{owner: owner},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, observerFailure)
			},
			expectRandomiseCall: &randomiseCall{Stable: true, AlbumIds: []catalog.AlbumId{avengersId}},
			expectObserverSeen: []map[catalog.AlbumId][]catalog.Cover{
				{avengersId: avengersCovers},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backfill := &catalog.BackfillCovers{
				FindAlbumByOwnerPort: tt.fields.FindAlbumByOwnerPort,
				CoverService:         tt.fields.CoverService,
				Observers:            []catalog.CoverBackfillObserver{tt.fields.Observer},
			}

			got, err := backfill.BackfillForOwner(context.Background(), tt.args.owner)
			if !tt.wantErr(t, err, fmt.Sprintf("BackfillForOwner(%s)", tt.args.owner)) {
				return
			}
			assert.Equal(t, tt.want, got, "BackfillForOwner result")
			if tt.expectRandomiseCall == nil {
				assert.Empty(t, tt.fields.CoverService.RandomiseCalls, "CoverService.Randomise calls")
			} else {
				assert.Equal(t, []randomiseCall{*tt.expectRandomiseCall}, tt.fields.CoverService.RandomiseCalls, "CoverService.Randomise calls")
			}
			assert.Equal(t, tt.expectObserverSeen, tt.fields.Observer.Notifications, "observer notifications")
		})
	}
}

type randomiseCall struct {
	Stable   bool
	AlbumIds []catalog.AlbumId
}

type coverServicePortFake struct {
	RandomiseCalls  []randomiseCall
	RandomiseResult map[catalog.AlbumId][]catalog.Cover
	Err             error
}

func (c *coverServicePortFake) Randomise(_ context.Context, stable bool, albumIds ...catalog.AlbumId) (map[catalog.AlbumId][]catalog.Cover, error) {
	c.RandomiseCalls = append(c.RandomiseCalls, randomiseCall{Stable: stable, AlbumIds: albumIds})
	if c.Err != nil {
		return nil, c.Err
	}
	if len(c.RandomiseResult) == 0 {
		return nil, nil
	}
	return c.RandomiseResult, nil
}

func (c *coverServicePortFake) StableRefresh(_ context.Context, _ catalog.TransferredMedias) (map[catalog.AlbumId][]catalog.Cover, error) {
	return nil, nil
}

type coverBackfillObserverFake struct {
	Notifications []map[catalog.AlbumId][]catalog.Cover
	Err           error
}

func (o *coverBackfillObserverFake) OnCoverBackfilled(_ context.Context, coversByAlbumId map[catalog.AlbumId][]catalog.Cover) error {
	o.Notifications = append(o.Notifications, coversByAlbumId)
	return o.Err
}

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
