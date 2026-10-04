package catalogdynamo

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

func TestRepository_FindCoversByAlbum(t *testing.T) {
	albumId := catalog.AlbumId{
		Owner:      "ironman",
		FolderName: catalog.NewFolderName("/avengers-1"),
	}

	dyn := dynamotestutils.NewTestContext(context.Background(), t)

	type args struct {
		albumId catalog.AlbumId
	}
	tests := []struct {
		name    string
		args    args
		before  []map[string]types.AttributeValue
		want    []catalog.Cover
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "it should return nil when the album has no cover record",
			args:    args{albumId: albumId},
			before:  nil,
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should return the covers preserving stored order",
			args: args{albumId: albumId},
			before: []map[string]types.AttributeValue{
				coverEntry(albumId,
					catalog.Cover{MediaId: "media-1", Filename: "a.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-2", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
				),
			},
			want: []catalog.Cover{
				{MediaId: "media-1", Filename: "a.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-2", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn = dyn.Subtest(t)

			err := dyn.WithDbContent(dyn.Ctx, tt.before)
			if !assert.NoError(t, err, "WithDbContent") {
				return
			}

			r := &Repository{client: dyn.Client, table: dyn.Table}

			got, err := r.FindCoversByAlbum(context.Background(), tt.args.albumId)
			if !tt.wantErr(t, err, fmt.Sprintf("FindCoversByAlbum(%v)", tt.args.albumId)) {
				return
			}
			assert.Equal(t, tt.want, got, "FindCoversByAlbum(%v)", tt.args.albumId)
		})
	}
}

func TestRepository_FindCoversByAlbums(t *testing.T) {
	avengersId := catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/avengers-1")}
	stealthId := catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/stealth")}
	otherId := catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/other")}

	dyn := dynamotestutils.NewTestContext(context.Background(), t)

	type args struct {
		albumIds []catalog.AlbumId
	}
	tests := []struct {
		name    string
		args    args
		before  []map[string]types.AttributeValue
		want    map[catalog.AlbumId][]catalog.Cover
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "it should return a nil map when no albumId is requested",
			args:    args{albumIds: nil},
			before:  nil,
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name:    "it should return an empty map when none of the requested albums has a cover record",
			args:    args{albumIds: []catalog.AlbumId{avengersId, stealthId}},
			before:  nil,
			want:    map[catalog.AlbumId][]catalog.Cover{},
			wantErr: assert.NoError,
		},
		{
			name: "it should return the covers of every requested album that has a record and skip the others",
			args: args{albumIds: []catalog.AlbumId{avengersId, stealthId, otherId}},
			before: []map[string]types.AttributeValue{
				coverEntry(avengersId,
					catalog.Cover{MediaId: "media-1", Filename: "a.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-2", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
				),
				coverEntry(stealthId,
					catalog.Cover{MediaId: "media-9", Filename: "z.jpg", Origin: catalog.CoverOriginRandom},
				),
			},
			want: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "a.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-2", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
				},
				stealthId: {
					{MediaId: "media-9", Filename: "z.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn = dyn.Subtest(t)

			err := dyn.WithDbContent(dyn.Ctx, tt.before)
			if !assert.NoError(t, err, "WithDbContent") {
				return
			}

			r := &Repository{client: dyn.Client, table: dyn.Table}

			got, err := r.FindCoversByAlbums(context.Background(), tt.args.albumIds...)
			if !tt.wantErr(t, err, fmt.Sprintf("FindCoversByAlbums(%v)", tt.args.albumIds)) {
				return
			}
			assert.Equal(t, tt.want, got, "FindCoversByAlbums(%v)", tt.args.albumIds)
		})
	}
}

func TestRepository_SaveCovers(t *testing.T) {
	albumId := catalog.AlbumId{
		Owner:      "ironman",
		FolderName: catalog.NewFolderName("/avengers-1"),
	}

	dyn := dynamotestutils.NewTestContext(context.Background(), t)

	type args struct {
		albumId catalog.AlbumId
		covers  []catalog.Cover
	}
	tests := []struct {
		name    string
		args    args
		before  []map[string]types.AttributeValue
		after   []map[string]types.AttributeValue
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should insert the covers record when none exists",
			args: args{
				albumId: albumId,
				covers: []catalog.Cover{
					{MediaId: "media-1", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-2", Filename: "b.jpg", Origin: catalog.CoverOriginCherryPicked},
				},
			},
			before: nil,
			after: []map[string]types.AttributeValue{
				coverEntry(albumId,
					catalog.Cover{MediaId: "media-1", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-2", Filename: "b.jpg", Origin: catalog.CoverOriginCherryPicked},
				),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should overwrite the covers record with the whole new set",
			args: args{
				albumId: albumId,
				covers: []catalog.Cover{
					{MediaId: "media-9", Filename: "z.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			before: []map[string]types.AttributeValue{
				coverEntry(albumId,
					catalog.Cover{MediaId: "media-1", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-2", Filename: "b.jpg", Origin: catalog.CoverOriginCherryPicked},
				),
			},
			after: []map[string]types.AttributeValue{
				coverEntry(albumId,
					catalog.Cover{MediaId: "media-9", Filename: "z.jpg", Origin: catalog.CoverOriginRandom},
				),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should delete the record when saving an empty cover set",
			args: args{albumId: albumId, covers: nil},
			before: []map[string]types.AttributeValue{
				coverEntry(albumId,
					catalog.Cover{MediaId: "media-1", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
				),
			},
			after:   nil,
			wantErr: assert.NoError,
		},
		{
			name:    "it should be a no-op when saving an empty cover set on an album with no cover record (idempotent delete)",
			args:    args{albumId: albumId, covers: nil},
			before:  nil,
			after:   nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should reject saving more than 4 covers",
			args: args{
				albumId: albumId,
				covers: []catalog.Cover{
					{MediaId: "m1", Filename: "1.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "m2", Filename: "2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "m3", Filename: "3.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "m4", Filename: "4.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "m5", Filename: "5.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			before: nil,
			after:  nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.TooManyCoversErr, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn = dyn.Subtest(t)

			err := dyn.WithDbContent(dyn.Ctx, tt.before)
			if !assert.NoError(t, err, "WithDbContent") {
				return
			}

			r := &Repository{client: dyn.Client, table: dyn.Table}

			err = r.SaveCovers(context.Background(), tt.args.albumId, tt.args.covers)
			if !tt.wantErr(t, err, fmt.Sprintf("SaveCovers(%v, %v)", tt.args.albumId, tt.args.covers)) {
				return
			}
			_, err = dyn.EqualContent(dyn.Ctx, tt.after)
			assert.NoError(t, err, "AssertDbContent")
		})
	}
}

func TestRepository_MoveCovers(t *testing.T) {
	fromId := catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/avengers-1")}
	toId := catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/avengers-renamed")}
	otherId := catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/stealth")}

	dyn := dynamotestutils.NewTestContext(context.Background(), t)

	type args struct {
		from catalog.AlbumId
		to   catalog.AlbumId
	}
	tests := []struct {
		name    string
		args    args
		before  []map[string]types.AttributeValue
		after   []map[string]types.AttributeValue
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "it should be a no-op when the source album has no cover record",
			args:    args{from: fromId, to: toId},
			before:  nil,
			after:   nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should move the cover record from the source to the destination preserving the cover list verbatim",
			args: args{from: fromId, to: toId},
			before: []map[string]types.AttributeValue{
				coverEntry(fromId,
					catalog.Cover{MediaId: "media-1", Filename: "a.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-2", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
				),
			},
			after: []map[string]types.AttributeValue{
				coverEntry(toId,
					catalog.Cover{MediaId: "media-1", Filename: "a.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-2", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
				),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should leave other album cover records untouched when moving one album's record",
			args: args{from: fromId, to: toId},
			before: []map[string]types.AttributeValue{
				coverEntry(fromId,
					catalog.Cover{MediaId: "media-1", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
				),
				coverEntry(otherId,
					catalog.Cover{MediaId: "media-9", Filename: "z.jpg", Origin: catalog.CoverOriginRandom},
				),
			},
			after: []map[string]types.AttributeValue{
				coverEntry(toId,
					catalog.Cover{MediaId: "media-1", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
				),
				coverEntry(otherId,
					catalog.Cover{MediaId: "media-9", Filename: "z.jpg", Origin: catalog.CoverOriginRandom},
				),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should be a no-op when from and to are the same album",
			args: args{from: fromId, to: fromId},
			before: []map[string]types.AttributeValue{
				coverEntry(fromId,
					catalog.Cover{MediaId: "media-1", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
				),
			},
			after: []map[string]types.AttributeValue{
				coverEntry(fromId,
					catalog.Cover{MediaId: "media-1", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
				),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn = dyn.Subtest(t)

			err := dyn.WithDbContent(dyn.Ctx, tt.before)
			if !assert.NoError(t, err, "WithDbContent") {
				return
			}

			r := &Repository{client: dyn.Client, table: dyn.Table}

			err = r.MoveCovers(context.Background(), tt.args.from, tt.args.to)
			if !tt.wantErr(t, err, fmt.Sprintf("MoveCovers(%v, %v)", tt.args.from, tt.args.to)) {
				return
			}
			_, err = dyn.EqualContent(dyn.Ctx, tt.after)
			assert.NoError(t, err, "AssertDbContent")
		})
	}
}

func coverEntry(albumId catalog.AlbumId, covers ...catalog.Cover) map[string]types.AttributeValue {
	items := make([]types.AttributeValue, 0, len(covers))
	for _, cover := range covers {
		items = append(items, &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"MediaId":  &types.AttributeValueMemberS{Value: string(cover.MediaId)},
				"Filename": &types.AttributeValueMemberS{Value: cover.Filename},
				"Origin":   &types.AttributeValueMemberS{Value: string(cover.Origin)},
			},
		})
	}
	return map[string]types.AttributeValue{
		"PK":              &types.AttributeValueMemberS{Value: fmt.Sprintf("%s#ALBUM", albumId.Owner)},
		"SK":              &types.AttributeValueMemberS{Value: fmt.Sprintf("ALBUM#%s#COVERS", albumId.FolderName)},
		"AlbumOwner":      &types.AttributeValueMemberS{Value: albumId.Owner.Value()},
		"AlbumFolderName": &types.AttributeValueMemberS{Value: albumId.FolderName.String()},
		"Covers":          &types.AttributeValueMemberL{Value: items},
	}
}
