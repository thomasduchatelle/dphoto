package catalog_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

var (
	tagOwner     = ownermodel.Owner("tony@stark.com")
	tagMediaId   = catalog.MediaId("media-1")
	tagHolidays  = catalog.NewAlbumIdFromStrings("tony@stark.com", "/2024-01_Holiday")
	favoriteTag  = mustTag("favorite")
	peopleLukeTg = mustTag("people:luke")
)

func mustTag(raw string) catalog.Tag {
	tag, err := catalog.NewTag(raw)
	if err != nil {
		panic(err)
	}
	return tag
}

// taggedMedia builds a seed media already carrying the given tags, so a case can start from a
// pre-tagged state.
func taggedMedia(mediaId catalog.MediaId, tags ...catalog.Tag) *catalog.MediaMeta {
	return &catalog.MediaMeta{Id: mediaId, Tags: tags}
}

func TestTagService_AddTag(t *testing.T) {
	type fields struct {
		Catalog *CatalogInMemory
	}
	type args struct {
		owner   ownermodel.Owner
		mediaId catalog.MediaId
		tag     catalog.Tag
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		expectTags map[catalog.MediaId][]string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "it should store the tag on the media",
			fields: fields{
				Catalog: NewCatalogInMemory(withMedias(tagHolidays, &catalog.MediaMeta{Id: tagMediaId})),
			},
			args:       args{owner: tagOwner, mediaId: tagMediaId, tag: favoriteTag},
			expectTags: map[catalog.MediaId][]string{tagMediaId: {"favorite"}},
			wantErr:    assert.NoError,
		},
		{
			name: "it should be idempotent when the same tag is added twice",
			fields: fields{
				Catalog: NewCatalogInMemory(withMedias(tagHolidays, taggedMedia(tagMediaId, favoriteTag))),
			},
			args:       args{owner: tagOwner, mediaId: tagMediaId, tag: favoriteTag},
			expectTags: map[catalog.MediaId][]string{tagMediaId: {"favorite"}},
			wantErr:    assert.NoError,
		},
		{
			name: "it should DENY tagging a media that is not owned by the owner",
			fields: fields{
				Catalog: NewCatalogInMemory(),
			},
			args:       args{owner: tagOwner, mediaId: tagMediaId, tag: favoriteTag},
			expectTags: map[catalog.MediaId][]string{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.MediaNotFoundError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := catalog.NewTagService(tt.fields.Catalog)
			err := service.AddTag(context.Background(), tt.args.owner, tt.args.mediaId, tt.args.tag)
			if tt.wantErr(t, err) {
				assert.Equal(t, tt.expectTags, tt.fields.Catalog.TagsByMedia())
			}
		})
	}
}

func TestTagService_RemoveTag(t *testing.T) {
	type fields struct {
		Catalog *CatalogInMemory
	}
	type args struct {
		owner   ownermodel.Owner
		mediaId catalog.MediaId
		tag     catalog.Tag
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		expectTags map[catalog.MediaId][]string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "it should leave no residual tag after removing the only tag",
			fields: fields{
				Catalog: NewCatalogInMemory(withMedias(tagHolidays, taggedMedia(tagMediaId, favoriteTag))),
			},
			args:       args{owner: tagOwner, mediaId: tagMediaId, tag: favoriteTag},
			expectTags: map[catalog.MediaId][]string{},
			wantErr:    assert.NoError,
		},
		{
			name: "it should be idempotent when removing a tag that is absent",
			fields: fields{
				Catalog: NewCatalogInMemory(withMedias(tagHolidays, &catalog.MediaMeta{Id: tagMediaId})),
			},
			args:       args{owner: tagOwner, mediaId: tagMediaId, tag: favoriteTag},
			expectTags: map[catalog.MediaId][]string{},
			wantErr:    assert.NoError,
		},
		{
			name: "it should remove only the targeted tag and keep the others",
			fields: fields{
				Catalog: NewCatalogInMemory(withMedias(tagHolidays, taggedMedia(tagMediaId, favoriteTag, peopleLukeTg))),
			},
			args:       args{owner: tagOwner, mediaId: tagMediaId, tag: favoriteTag},
			expectTags: map[catalog.MediaId][]string{tagMediaId: {"people:luke"}},
			wantErr:    assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := catalog.NewTagService(tt.fields.Catalog)
			err := service.RemoveTag(context.Background(), tt.args.owner, tt.args.mediaId, tt.args.tag)
			if tt.wantErr(t, err) {
				assert.Equal(t, tt.expectTags, tt.fields.Catalog.TagsByMedia())
			}
		})
	}
}

func TestTagService_ListMediasByTag(t *testing.T) {
	otherMediaId := catalog.MediaId("media-2")
	otherAlbum := catalog.NewAlbumIdFromStrings("tony@stark.com", "/2024-07_Summer")

	type fields struct {
		Catalog *CatalogInMemory
	}
	type args struct {
		owner ownermodel.Owner
		tag   catalog.Tag
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantMedias []*catalog.MediaWithAlbum
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "it should return every media carrying the tag across albums with its full tag set",
			fields: fields{
				Catalog: NewCatalogInMemory(
					withMedias(tagHolidays, taggedMedia(tagMediaId, favoriteTag, peopleLukeTg)),
					withMedias(otherAlbum, taggedMedia(otherMediaId, favoriteTag)),
				),
			},
			args: args{owner: tagOwner, tag: favoriteTag},
			wantMedias: []*catalog.MediaWithAlbum{
				{MediaMeta: catalog.MediaMeta{Id: tagMediaId, Tags: []catalog.Tag{favoriteTag, peopleLukeTg}}, AlbumId: tagHolidays},
				{MediaMeta: catalog.MediaMeta{Id: otherMediaId, Tags: []catalog.Tag{favoriteTag}}, AlbumId: otherAlbum},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return no media when no media carries the tag",
			fields: fields{
				Catalog: NewCatalogInMemory(withMedias(tagHolidays, taggedMedia(tagMediaId, favoriteTag))),
			},
			args:       args{owner: tagOwner, tag: peopleLukeTg},
			wantMedias: nil,
			wantErr:    assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := catalog.NewTagService(tt.fields.Catalog)
			got, err := service.ListMediasByTag(context.Background(), tt.args.owner, tt.args.tag)
			if tt.wantErr(t, err) {
				assert.ElementsMatch(t, tt.wantMedias, got)
			}
		})
	}
}
