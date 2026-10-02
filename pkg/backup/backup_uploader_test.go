package backup

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

func TestUploader_OnMediaCatalogued_completesCovers(t *testing.T) {
	const owner = ownermodel.Owner("ironman")
	album1Reference1 := &CatalogReferenceStub{MediaIdValue: "media-1", AlbumFolderNameValue: "/album1"}
	album1Reference2 := &CatalogReferenceStub{MediaIdValue: "media-2", AlbumFolderNameValue: "/album1"}
	album2Reference := &CatalogReferenceStub{MediaIdValue: "media-3", AlbumFolderNameValue: "/album2"}
	videoReference := &CatalogReferenceStub{MediaIdValue: "media-4", AlbumFolderNameValue: "/album1"}
	otherReference := &CatalogReferenceStub{MediaIdValue: "media-5", AlbumFolderNameValue: "/album3"}

	imageMedia := func(name string) *AnalysedMedia {
		return &AnalysedMedia{
			FoundMedia: NewInMemoryMedia(name, time.Now(), []byte("content")),
			Type:       MediaTypeImage,
			Details:    &MediaDetails{DateTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		}
	}
	videoMedia := &AnalysedMedia{
		FoundMedia: NewInMemoryMedia("clip.mp4", time.Now(), []byte("clip")),
		Type:       MediaTypeVideo,
		Details:    &MediaDetails{DateTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	otherMedia := &AnalysedMedia{
		FoundMedia: NewInMemoryMedia("note.txt", time.Now(), []byte("note")),
		Type:       MediaTypeOther,
		Details:    &MediaDetails{DateTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	type fields struct {
		CompleteCoversPort *CompleteCoversPortFake
	}
	type args struct {
		requests []BackingUpMediaRequest
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		expectCoverRequests map[string][]CoverCandidate
	}{
		{
			name:   "it should pass the inserted IMAGE medias as candidates, grouped by album",
			fields: fields{CompleteCoversPort: NewCompleteCoversPortFake()},
			args: args{
				requests: []BackingUpMediaRequest{
					{AnalysedMedia: imageMedia("photo-1.jpg"), CatalogReference: album1Reference1},
					{AnalysedMedia: imageMedia("photo-2.jpg"), CatalogReference: album1Reference2},
					{AnalysedMedia: imageMedia("photo-3.jpg"), CatalogReference: album2Reference},
				},
			},
			expectCoverRequests: map[string][]CoverCandidate{
				"/album1": {
					{MediaId: "media-1", Filename: "photo-1.jpg"},
					{MediaId: "media-2", Filename: "photo-2.jpg"},
				},
				"/album2": {
					{MediaId: "media-3", Filename: "photo-3.jpg"},
				},
			},
		},
		{
			name:   "it should ignore non-IMAGE medias as cover candidates",
			fields: fields{CompleteCoversPort: NewCompleteCoversPortFake()},
			args: args{
				requests: []BackingUpMediaRequest{
					{AnalysedMedia: videoMedia, CatalogReference: videoReference},
					{AnalysedMedia: otherMedia, CatalogReference: otherReference},
				},
			},
			expectCoverRequests: nil,
		},
		{
			name:   "it should skip albums that received no IMAGE in the batch",
			fields: fields{CompleteCoversPort: NewCompleteCoversPortFake()},
			args: args{
				requests: []BackingUpMediaRequest{
					{AnalysedMedia: imageMedia("photo-1.jpg"), CatalogReference: album1Reference1},
					{AnalysedMedia: videoMedia, CatalogReference: videoReference},
				},
			},
			expectCoverRequests: map[string][]CoverCandidate{
				"/album1": {
					{MediaId: "media-1", Filename: "photo-1.jpg"},
				},
			},
		},
		{
			name:   "it should not call the completion port when no media is inserted at all",
			fields: fields{CompleteCoversPort: NewCompleteCoversPortFake()},
			args: args{
				requests: []BackingUpMediaRequest{},
			},
			expectCoverRequests: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uploader{
				Owner:              owner,
				InsertMediaPort:    &NoopInsertMediaPort{},
				ArchivePort:        &SameNameArchiveMediaPort{},
				CompleteCoversPort: tt.fields.CompleteCoversPort,
			}

			err := u.OnMediaCatalogued(context.Background(), tt.args.requests)

			if assert.NoError(t, err) {
				assert.Equal(t, tt.expectCoverRequests, tt.fields.CompleteCoversPort.CandidatesByAlbum[owner])
			}
		})
	}
}

func TestUploader_OnMediaCatalogued_isCompatibleWithoutCompleteCoversPort(t *testing.T) {
	u := &uploader{
		Owner:           ownermodel.Owner("ironman"),
		InsertMediaPort: &NoopInsertMediaPort{},
		ArchivePort:     &SameNameArchiveMediaPort{},
	}

	err := u.OnMediaCatalogued(context.Background(), []BackingUpMediaRequest{
		{
			AnalysedMedia: &AnalysedMedia{
				FoundMedia: NewInMemoryMedia("photo.jpg", time.Now(), []byte("x")),
				Type:       MediaTypeImage,
				Details:    &MediaDetails{DateTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			CatalogReference: &CatalogReferenceStub{MediaIdValue: "media-1", AlbumFolderNameValue: "/album1"},
		},
	})

	assert.NoError(t, err)
}

type NoopInsertMediaPort struct{}

func (n *NoopInsertMediaPort) IndexMedias(ctx context.Context, owner ownermodel.Owner, requests []*CatalogMediaRequest) error {
	return nil
}

type SameNameArchiveMediaPort struct{}

func (s *SameNameArchiveMediaPort) ArchiveMedia(owner string, media *BackingUpMediaRequest) (string, error) {
	return media.AnalysedMedia.FoundMedia.MediaPath().Filename, nil
}

func NewCompleteCoversPortFake() *CompleteCoversPortFake {
	return &CompleteCoversPortFake{
		CandidatesByAlbum: make(map[ownermodel.Owner]map[string][]CoverCandidate),
	}
}

type CompleteCoversPortFake struct {
	lock              sync.Mutex
	CandidatesByAlbum map[ownermodel.Owner]map[string][]CoverCandidate
}

func (f *CompleteCoversPortFake) CompleteCoversFromCandidates(ctx context.Context, owner ownermodel.Owner, albumFolderName string, candidates []CoverCandidate) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.CandidatesByAlbum[owner] == nil {
		f.CandidatesByAlbum[owner] = make(map[string][]CoverCandidate)
	}
	f.CandidatesByAlbum[owner][albumFolderName] = append(f.CandidatesByAlbum[owner][albumFolderName], candidates...)
	return nil
}
