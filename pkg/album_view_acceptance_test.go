package pkg

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func TestAlbumViewAcceptance(t *testing.T) {
	suite.Run(t, new(AlbumViewTestSuite))
}

// AlbumViewTestSuite drives the AlbumView through every catalog operation that must fan out
// to the per-user summary rows: create, insert medias, share, amend dates, rename in place,
// delete. Each Test<N>_ step reads back through pkgfactory.AlbumView(...).ListAlbums and
// asserts the view reflects the change.
type AlbumViewTestSuite struct {
	suite.Suite
	owner       ownermodel.Owner
	visitor     ownermodel.Owner
	ownerUser   usermodel.CurrentUser
	visitorUser usermodel.CurrentUser
	albumId     catalog.AlbumId
	start       time.Time
	end         time.Time
}

func (s *AlbumViewTestSuite) SetupSuite() {
	ctx := context.Background()

	if err := initForLocalstack(ctx); !assert.NoError(s.T(), err) {
		return
	}

	owner, err := createRandomUser(ctx)
	if !assert.NoError(s.T(), err) {
		return
	}
	s.owner = owner
	s.ownerUser = usermodel.CurrentUser{UserId: usermodel.NewUserId(owner.Value()), Owner: &s.owner}

	visitor, err := createRandomUser(ctx)
	if !assert.NoError(s.T(), err) {
		return
	}
	s.visitor = visitor
	s.visitorUser = usermodel.CurrentUser{UserId: usermodel.NewUserId(visitor.Value())}

	s.start = time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC)
	s.end = time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC)
}

func (s *AlbumViewTestSuite) Test10_CreateAlbum_ownerSeesEmptyAlbumInView() {
	ctx := context.Background()
	t := s.T()

	albumIdPtr, err := pkgfactory.Catalog().CreateAlbumCase(ctx).Create(ctx, catalog.CreateAlbumRequest{
		Owner: s.owner,
		Name:  "Spring 2026",
		Start: s.start,
		End:   s.end,
	})
	if !assert.NoError(t, err) {
		return
	}
	s.albumId = *albumIdPtr

	albums := listOwnerAlbums(ctx, t, s.ownerUser)
	visible := findAlbum(albums, s.albumId)
	if assert.NotNil(t, visible, "owner should see the freshly created album") {
		assert.Equal(t, "Spring 2026", visible.Name)
		assert.Equal(t, s.start, visible.Start)
		assert.Equal(t, s.end, visible.End)
		assert.Equal(t, 0, visible.MediaCount)
		assert.True(t, visible.OwnedByCurrentUser)
	}
}

func (s *AlbumViewTestSuite) Test20_InsertMedias_ownerCountReflectsInsert() {
	ctx := context.Background()
	t := s.T()

	medias := []catalog.CreateMediaRequest{
		newMediaRequest(s.albumId, "spring-a", s.start.Add(24*time.Hour)),
		newMediaRequest(s.albumId, "spring-b", s.start.Add(48*time.Hour)),
	}

	err := pkgfactory.InsertMediasCase(ctx).Insert(ctx, s.owner, medias)
	if !assert.NoError(t, err) {
		return
	}

	albums := listOwnerAlbums(ctx, t, s.ownerUser)
	visible := findAlbum(albums, s.albumId)
	if assert.NotNil(t, visible, "owner should still see the album") {
		assert.Equal(t, 2, visible.MediaCount)
	}
}

func (s *AlbumViewTestSuite) Test30_ShareAlbum_visitorSeesAlbumWithCount() {
	ctx := context.Background()
	t := s.T()

	err := pkgfactory.AclCatalogShare(ctx).ShareAlbumWith(ctx, s.albumId, s.visitorUser.UserId)
	if !assert.NoError(t, err) {
		return
	}

	albums := listAllAlbums(ctx, t, s.visitorUser)
	visible := findAlbum(albums, s.albumId)
	if assert.NotNil(t, visible, "visitor should see the album after share") {
		assert.Equal(t, "Spring 2026", visible.Name)
		assert.Equal(t, s.start, visible.Start)
		assert.Equal(t, s.end, visible.End)
		assert.Equal(t, 2, visible.MediaCount)
		assert.False(t, visible.OwnedByCurrentUser)
	}
}

func (s *AlbumViewTestSuite) Test40_AmendAlbumDates_bothSeeNewDates() {
	ctx := context.Background()
	t := s.T()

	newStart := time.Date(2026, time.February, 15, 0, 0, 0, 0, time.UTC)
	newEnd := time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC)

	err := pkgfactory.Catalog().AmendAlbumDatesCase(ctx).AmendAlbumDates(ctx, s.albumId, newStart, newEnd)
	if !assert.NoError(t, err) {
		return
	}
	s.start = newStart
	s.end = newEnd

	ownerAlbums := listOwnerAlbums(ctx, t, s.ownerUser)
	ownerVisible := findAlbum(ownerAlbums, s.albumId)
	if assert.NotNil(t, ownerVisible, "owner should still see the album after date amend") {
		assert.Equal(t, newStart, ownerVisible.Start)
		assert.Equal(t, newEnd, ownerVisible.End)
	}

	visitorAlbums := listAllAlbums(ctx, t, s.visitorUser)
	visitorVisible := findAlbum(visitorAlbums, s.albumId)
	if assert.NotNil(t, visitorVisible, "visitor should still see the album after date amend") {
		assert.Equal(t, newStart, visitorVisible.Start)
		assert.Equal(t, newEnd, visitorVisible.End)
	}
}

func (s *AlbumViewTestSuite) Test50_RenameAlbumInPlace_bothSeeNewName() {
	ctx := context.Background()
	t := s.T()

	err := pkgfactory.Catalog().RenameAlbumCase(ctx).RenameAlbum(ctx, catalog.RenameAlbumRequest{
		CurrentId:    s.albumId,
		NewName:      "Spring Blossoms 2026",
		RenameFolder: false,
	})
	if !assert.NoError(t, err) {
		return
	}

	ownerAlbums := listOwnerAlbums(ctx, t, s.ownerUser)
	ownerVisible := findAlbum(ownerAlbums, s.albumId)
	if assert.NotNil(t, ownerVisible, "owner should still see the album after in-place rename") {
		assert.Equal(t, "Spring Blossoms 2026", ownerVisible.Name)
	}

	visitorAlbums := listAllAlbums(ctx, t, s.visitorUser)
	visitorVisible := findAlbum(visitorAlbums, s.albumId)
	if assert.NotNil(t, visitorVisible, "visitor should still see the album after in-place rename") {
		assert.Equal(t, "Spring Blossoms 2026", visitorVisible.Name)
	}
}

func (s *AlbumViewTestSuite) Test60_DeleteAlbum_neitherSeesTheAlbum() {
	ctx := context.Background()
	t := s.T()

	fallbackId, err := pkgfactory.Catalog().CreateAlbumCase(ctx).Create(ctx, catalog.CreateAlbumRequest{
		Owner: s.owner,
		Name:  "Catch-all 2026",
		Start: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC),
	})
	if !assert.NoError(t, err, "fallback album to absorb transferred medias") {
		return
	}

	err = pkgfactory.Catalog().CreateAlbumDeleteCase(ctx).DeleteAlbum(ctx, s.albumId)
	if !assert.NoError(t, err) {
		return
	}

	ownerAlbums := listOwnerAlbums(ctx, t, s.ownerUser)
	assert.Nil(t, findAlbum(ownerAlbums, s.albumId), "owner should no longer see the deleted album")
	if catchAll := findAlbum(ownerAlbums, *fallbackId); assert.NotNil(t, catchAll, "fallback album should absorb the transferred medias") {
		assert.Equal(t, 2, catchAll.MediaCount, "fallback album should carry the transferred medias")
	}

	visitorAlbums := listAllAlbums(ctx, t, s.visitorUser)
	assert.Nil(t, findAlbum(visitorAlbums, s.albumId), "visitor should no longer see the deleted album")
}

func listOwnerAlbums(ctx context.Context, t *testing.T, user usermodel.CurrentUser) []*catalogviews.VisibleAlbum {
	t.Helper()
	albums, err := pkgfactory.AlbumView(ctx).ListAlbums(ctx, user, catalogviews.ListAlbumsFilter{OnlyDirectlyOwned: true})
	assert.NoError(t, err)
	return albums
}

func listAllAlbums(ctx context.Context, t *testing.T, user usermodel.CurrentUser) []*catalogviews.VisibleAlbum {
	t.Helper()
	albums, err := pkgfactory.AlbumView(ctx).ListAlbums(ctx, user, catalogviews.ListAlbumsFilter{})
	assert.NoError(t, err)
	return albums
}

func findAlbum(albums []*catalogviews.VisibleAlbum, albumId catalog.AlbumId) *catalogviews.VisibleAlbum {
	for _, album := range albums {
		if album.AlbumId.IsEqual(albumId) {
			return album
		}
	}
	return nil
}

func newMediaRequest(albumId catalog.AlbumId, suffix string, captured time.Time) catalog.CreateMediaRequest {
	hash := sha256.Sum256([]byte(albumId.String() + "/" + suffix))
	signature := catalog.MediaSignature{
		SignatureSha256: hex.EncodeToString(hash[:]),
		SignatureSize:   len(suffix),
	}
	mediaId, _ := catalog.GenerateMediaId(signature)
	return catalog.CreateMediaRequest{
		Id:         mediaId,
		Signature:  signature,
		FolderName: albumId.FolderName,
		Filename:   suffix + ".jpg",
		Type:       catalog.MediaType("IMAGE"),
		Details: catalog.MediaDetails{
			DateTime: captured,
		},
	}
}
