package catalogviews

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

func TestCoverProjection_OnCoversChanged(t *testing.T) {
	coverA := catalog.Cover{MediaId: "media-a", Filename: "a.jpg", Origin: catalog.CoverOriginRandom}
	coverB := catalog.Cover{MediaId: "media-b", Filename: "b.jpg", Origin: catalog.CoverOriginRandom}
	coverC := catalog.Cover{MediaId: "media-c", Filename: "c.jpg", Origin: catalog.CoverOriginCherryPicked}
	coverD := catalog.Cover{MediaId: "media-d", Filename: "d.jpg", Origin: catalog.CoverOriginRandom}
	full := []catalog.Cover{coverA, coverB, coverC, coverD}

	seededOwnerAndVisitor := func() *AlbumSummaryInMemoryRepository {
		return &AlbumSummaryInMemoryRepository{
			Summaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 3},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 3},
					Availability: VisitorAvailability(visitorUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumBeta, Name: "Beta", Start: feb24, End: mar24, MediaCount: 5},
					Availability: OwnerAvailability(ownerUserId),
				},
			},
		}
	}
	seededOwnerAndVisitorWithFullCovers := func() *AlbumSummaryInMemoryRepository {
		repo := seededOwnerAndVisitor()
		for i := range repo.Summaries {
			if repo.Summaries[i].AlbumSummary.AlbumId.IsEqual(albumAlpha) {
				repo.Summaries[i].AlbumSummary.Covers = full
			}
		}
		return repo
	}

	type fields struct {
		Repository *AlbumSummaryInMemoryRepository
	}
	tests := []struct {
		name       string
		fields     fields
		albumId    catalog.AlbumId
		covers     []catalog.Cover
		expectRepo []UserAlbumSummary
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:    "it should fan out a full cover set to every viewer row of the album",
			fields:  fields{Repository: seededOwnerAndVisitor()},
			albumId: albumAlpha,
			covers:  full,
			expectRepo: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 3, Covers: full},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 3, Covers: full},
					Availability: VisitorAvailability(visitorUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumBeta, Name: "Beta", Start: feb24, End: mar24, MediaCount: 5},
					Availability: OwnerAvailability(ownerUserId),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "it should clear the covers when the new set is empty",
			fields:  fields{Repository: seededOwnerAndVisitorWithFullCovers()},
			albumId: albumAlpha,
			covers:  nil,
			expectRepo: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 3},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumAlpha, Name: "Alpha", Start: jan24, End: feb24, MediaCount: 3},
					Availability: VisitorAvailability(visitorUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumBeta, Name: "Beta", Start: feb24, End: mar24, MediaCount: 5},
					Availability: OwnerAvailability(ownerUserId),
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projection := &CoverProjection{Repository: tt.fields.Repository}
			err := projection.OnCoversChanged(context.Background(), tt.albumId, tt.covers)
			if tt.wantErr(t, err) {
				assert.ElementsMatch(t, tt.expectRepo, tt.fields.Repository.Summaries)
			}
		})
	}
}
