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

// TagRepositoryInMemory is a Fake catalog.TagRepository backed by a map keyed by the media
// partition and the normalised tag. It records the AlbumId stored alongside each tag so a
// test can read the persisted state back.
type TagRepositoryInMemory struct {
	tags map[catalog.MediaId]map[string]catalog.AlbumId
}

type tagSeed struct {
	albumId catalog.AlbumId
	mediaId catalog.MediaId
	tag     catalog.Tag
}

func taggedWith(albumId catalog.AlbumId, mediaId catalog.MediaId, tag catalog.Tag) tagSeed {
	return tagSeed{albumId: albumId, mediaId: mediaId, tag: tag}
}

func NewTagRepositoryInMemory(seeds ...tagSeed) *TagRepositoryInMemory {
	repo := &TagRepositoryInMemory{tags: make(map[catalog.MediaId]map[string]catalog.AlbumId)}
	for _, seed := range seeds {
		repo.put(seed.albumId, seed.mediaId, seed.tag)
	}
	return repo
}

func (r *TagRepositoryInMemory) put(albumId catalog.AlbumId, mediaId catalog.MediaId, tag catalog.Tag) {
	if r.tags[mediaId] == nil {
		r.tags[mediaId] = make(map[string]catalog.AlbumId)
	}
	r.tags[mediaId][tag.String()] = albumId
}

func (r *TagRepositoryInMemory) AddTag(_ context.Context, albumId catalog.AlbumId, mediaId catalog.MediaId, tag catalog.Tag) error {
	r.put(albumId, mediaId, tag)
	return nil
}

func (r *TagRepositoryInMemory) RemoveTag(_ context.Context, _ ownermodel.Owner, mediaId catalog.MediaId, tag catalog.Tag) error {
	if tags, ok := r.tags[mediaId]; ok {
		delete(tags, tag.String())
		if len(tags) == 0 {
			delete(r.tags, mediaId)
		}
	}
	return nil
}

func (r *TagRepositoryInMemory) FindMediasByTag(_ context.Context, owner ownermodel.Owner, tag catalog.Tag) ([]*catalog.MediaWithAlbum, error) {
	var medias []*catalog.MediaWithAlbum
	for mediaId, tags := range r.tags {
		if albumId, ok := tags[tag.String()]; ok && albumId.Owner == owner {
			medias = append(medias, &catalog.MediaWithAlbum{
				MediaMeta: catalog.MediaMeta{Id: mediaId},
				AlbumId:   albumId,
			})
		}
	}
	return medias, nil
}

// storedTags projects the fake into a flat list of (mediaId, tag) pairs so a test can assert
// the full persisted state without caring about album details.
func (r *TagRepositoryInMemory) storedTags() map[catalog.MediaId][]string {
	result := make(map[catalog.MediaId][]string)
	for mediaId, tags := range r.tags {
		for tag := range tags {
			result[mediaId] = append(result[mediaId], tag)
		}
	}
	return result
}

func TestTagService_AddTag(t *testing.T) {
	type fields struct {
		MediaReadRepository catalog.MediaReadRepository
		TagRepository       *TagRepositoryInMemory
	}
	type args struct {
		owner   ownermodel.Owner
		mediaId catalog.MediaId
		tag     catalog.Tag
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		expectStoredTag map[catalog.MediaId][]string
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "it should store the tag as a sibling of the media on its album",
			fields: fields{
				MediaReadRepository: NewCatalogInMemory(withMedias(tagHolidays, &catalog.MediaMeta{Id: tagMediaId})),
				TagRepository:       NewTagRepositoryInMemory(),
			},
			args:            args{owner: tagOwner, mediaId: tagMediaId, tag: favoriteTag},
			expectStoredTag: map[catalog.MediaId][]string{tagMediaId: {"favorite"}},
			wantErr:         assert.NoError,
		},
		{
			name: "it should be idempotent when the same tag is added twice",
			fields: fields{
				MediaReadRepository: NewCatalogInMemory(withMedias(tagHolidays, &catalog.MediaMeta{Id: tagMediaId})),
				TagRepository:       NewTagRepositoryInMemory(taggedWith(tagHolidays, tagMediaId, favoriteTag)),
			},
			args:            args{owner: tagOwner, mediaId: tagMediaId, tag: favoriteTag},
			expectStoredTag: map[catalog.MediaId][]string{tagMediaId: {"favorite"}},
			wantErr:         assert.NoError,
		},
		{
			name: "it should DENY tagging a media that is not owned by the owner",
			fields: fields{
				MediaReadRepository: NewCatalogInMemory(),
				TagRepository:       NewTagRepositoryInMemory(),
			},
			args:            args{owner: tagOwner, mediaId: tagMediaId, tag: favoriteTag},
			expectStoredTag: map[catalog.MediaId][]string{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.MediaNotFoundError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := catalog.NewTagService(tt.fields.MediaReadRepository, tt.fields.TagRepository)
			err := service.AddTag(context.Background(), tt.args.owner, tt.args.mediaId, tt.args.tag)
			if tt.wantErr(t, err) {
				assert.Equal(t, tt.expectStoredTag, tt.fields.TagRepository.storedTags())
			}
		})
	}
}

func TestTagService_RemoveTag(t *testing.T) {
	type fields struct {
		TagRepository *TagRepositoryInMemory
	}
	type args struct {
		owner   ownermodel.Owner
		mediaId catalog.MediaId
		tag     catalog.Tag
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		expectStoredTag map[catalog.MediaId][]string
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "it should leave no residual item after removing the tag",
			fields: fields{
				TagRepository: NewTagRepositoryInMemory(taggedWith(tagHolidays, tagMediaId, favoriteTag)),
			},
			args:            args{owner: tagOwner, mediaId: tagMediaId, tag: favoriteTag},
			expectStoredTag: map[catalog.MediaId][]string{},
			wantErr:         assert.NoError,
		},
		{
			name: "it should be idempotent when removing a tag that is absent",
			fields: fields{
				TagRepository: NewTagRepositoryInMemory(),
			},
			args:            args{owner: tagOwner, mediaId: tagMediaId, tag: favoriteTag},
			expectStoredTag: map[catalog.MediaId][]string{},
			wantErr:         assert.NoError,
		},
		{
			name: "it should remove only the targeted tag and keep the others",
			fields: fields{
				TagRepository: NewTagRepositoryInMemory(
					taggedWith(tagHolidays, tagMediaId, favoriteTag),
					taggedWith(tagHolidays, tagMediaId, peopleLukeTg),
				),
			},
			args:            args{owner: tagOwner, mediaId: tagMediaId, tag: favoriteTag},
			expectStoredTag: map[catalog.MediaId][]string{tagMediaId: {"people:luke"}},
			wantErr:         assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := catalog.NewTagService(nil, tt.fields.TagRepository)
			err := service.RemoveTag(context.Background(), tt.args.owner, tt.args.mediaId, tt.args.tag)
			if tt.wantErr(t, err) {
				assert.Equal(t, tt.expectStoredTag, tt.fields.TagRepository.storedTags())
			}
		})
	}
}

func TestTagService_ListMediasByTag(t *testing.T) {
	otherMediaId := catalog.MediaId("media-2")
	otherAlbum := catalog.NewAlbumIdFromStrings("tony@stark.com", "/2024-07_Summer")

	type fields struct {
		TagRepository *TagRepositoryInMemory
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
			name: "it should return every media carrying the tag across albums",
			fields: fields{
				TagRepository: NewTagRepositoryInMemory(
					taggedWith(tagHolidays, tagMediaId, favoriteTag),
					taggedWith(otherAlbum, otherMediaId, favoriteTag),
					taggedWith(tagHolidays, tagMediaId, peopleLukeTg),
				),
			},
			args: args{owner: tagOwner, tag: favoriteTag},
			wantMedias: []*catalog.MediaWithAlbum{
				{MediaMeta: catalog.MediaMeta{Id: tagMediaId}, AlbumId: tagHolidays},
				{MediaMeta: catalog.MediaMeta{Id: otherMediaId}, AlbumId: otherAlbum},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return no media when no media carries the tag",
			fields: fields{
				TagRepository: NewTagRepositoryInMemory(taggedWith(tagHolidays, tagMediaId, favoriteTag)),
			},
			args:       args{owner: tagOwner, tag: peopleLukeTg},
			wantMedias: nil,
			wantErr:    assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := catalog.NewTagService(nil, tt.fields.TagRepository)
			got, err := service.ListMediasByTag(context.Background(), tt.args.owner, tt.args.tag)
			if tt.wantErr(t, err) {
				assert.ElementsMatch(t, tt.wantMedias, got)
			}
		})
	}
}
