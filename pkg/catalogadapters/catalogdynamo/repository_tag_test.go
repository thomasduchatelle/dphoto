package catalogdynamo

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

type TagTestSuite struct {
	suite.Suite
	repo       *Repository
	owner      string
	holidays   catalog.AlbumId
	summer     catalog.AlbumId
	mediaInJan catalog.MediaId
	mediaInJul catalog.MediaId
	favorite   catalog.Tag
	peopleLuke catalog.Tag
}

func TestTagSuite(t *testing.T) {
	suite.Run(t, new(TagTestSuite))
}

func mustTag(raw string) catalog.Tag {
	tag, err := catalog.NewTag(raw)
	if err != nil {
		panic(err)
	}
	return tag
}

func (s *TagTestSuite) SetupTest() {
	dyn := dynamotestutils.NewTestContext(context.Background(), s.T())
	if err := appdynamodb.CreateTableIfNecessary(context.Background(), dyn.Table, dyn.Client, true); err != nil {
		s.FailNow(err.Error())
	}

	s.repo = &Repository{client: dyn.Client, table: dyn.Table}
	s.owner = "tony@stark.com"
	s.holidays = catalog.NewAlbumIdFromStrings(s.owner, "/2024-01_Holiday")
	s.summer = catalog.NewAlbumIdFromStrings(s.owner, "/2024-07_Summer")
	s.favorite = mustTag("favorite")
	s.peopleLuke = mustTag("people:luke")

	for _, albumId := range []catalog.AlbumId{s.holidays, s.summer} {
		if err := s.repo.InsertAlbum(context.Background(), catalog.Album{
			AlbumId: albumId,
			Name:    albumId.FolderName.String(),
			Start:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			End:     time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		}); err != nil {
			s.FailNow(err.Error())
		}
	}

	janSig := catalog.MediaSignature{SignatureSha256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", SignatureSize: 10}
	julSig := catalog.MediaSignature{SignatureSha256: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", SignatureSize: 20}
	s.mediaInJan = mustGenerateMediaId(catalog.GenerateMediaId(janSig))
	s.mediaInJul = mustGenerateMediaId(catalog.GenerateMediaId(julSig))

	err := s.repo.InsertMedias(context.Background(), s.holidays.Owner, []catalog.CreateMediaRequest{
		{Id: s.mediaInJan, Signature: janSig, FolderName: s.holidays.FolderName, Filename: "jan.jpg", Type: catalog.MediaTypeImage, Details: catalog.MediaDetails{DateTime: time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC)}},
		{Id: s.mediaInJul, Signature: julSig, FolderName: s.summer.FolderName, Filename: "jul.jpg", Type: catalog.MediaTypeImage, Details: catalog.MediaDetails{DateTime: time.Date(2024, 7, 2, 10, 0, 0, 0, time.UTC)}},
	})
	if err != nil {
		s.FailNow(err.Error())
	}
}

func (s *TagTestSuite) tagsOfMediaInAlbum(albumId catalog.AlbumId, mediaId catalog.MediaId) []string {
	medias, err := s.repo.FindMedias(context.Background(), catalog.NewFindMediaRequest(albumId.Owner).WithAlbum(albumId.FolderName))
	s.Require().NoError(err)
	for _, media := range medias {
		if media.Id == mediaId {
			return tagStringsSorted(media.Tags)
		}
	}
	return nil
}

func (s *TagTestSuite) TestAddTag_isStoredWithTheMediaAndIncludedInAlbumMedias() {
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays, s.mediaInJan, s.favorite))

	s.Equal([]string{"favorite"}, s.tagsOfMediaInAlbum(s.holidays, s.mediaInJan),
		"the album medias should carry the tag without a per-media query")
}

func (s *TagTestSuite) TestAddTag_isIdempotent() {
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays, s.mediaInJan, s.favorite))
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays, s.mediaInJan, s.favorite))

	s.Equal([]string{"favorite"}, s.tagsOfMediaInAlbum(s.holidays, s.mediaInJan))
}

func (s *TagTestSuite) TestRemoveTag_leavesNoResidualItem() {
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays, s.mediaInJan, s.favorite))
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays, s.mediaInJan, s.peopleLuke))

	s.Require().NoError(s.repo.RemoveTag(context.Background(), s.holidays.Owner, s.mediaInJan, s.favorite))

	s.Equal([]string{"people:luke"}, s.tagsOfMediaInAlbum(s.holidays, s.mediaInJan),
		"only the removed tag should disappear")

	medias, err := s.repo.FindMediasByTag(context.Background(), s.holidays.Owner, s.favorite)
	s.Require().NoError(err)
	s.Empty(medias, "no media should carry the removed tag anymore")
}

func (s *TagTestSuite) TestRemoveTag_isIdempotentWhenTagIsAbsent() {
	s.Require().NoError(s.repo.RemoveTag(context.Background(), s.holidays.Owner, s.mediaInJan, s.favorite))
	s.Empty(s.tagsOfMediaInAlbum(s.holidays, s.mediaInJan))
}

func (s *TagTestSuite) TestFindMediasByTag_returnsEveryMatchingMediaWithItsAlbumAndTagsAcrossAlbums() {
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays, s.mediaInJan, s.favorite))
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays, s.mediaInJan, s.peopleLuke))
	s.Require().NoError(s.repo.AddTag(context.Background(), s.summer, s.mediaInJul, s.favorite))

	medias, err := s.repo.FindMediasByTag(context.Background(), s.holidays.Owner, s.favorite)
	s.Require().NoError(err)

	sort.Slice(medias, func(i, j int) bool { return medias[i].Id < medias[j].Id })
	s.Require().Len(medias, 2)

	byMedia := map[catalog.MediaId]*catalog.MediaWithAlbum{}
	for _, m := range medias {
		byMedia[m.Id] = m
	}

	jan := byMedia[s.mediaInJan]
	s.Require().NotNil(jan)
	s.Equal(s.holidays, jan.AlbumId)
	s.Equal("jan.jpg", jan.Filename)
	s.Equal([]string{"favorite", "people:luke"}, tagStringsSorted(jan.Tags), "the full tag set is returned, not only the queried tag")

	jul := byMedia[s.mediaInJul]
	s.Require().NotNil(jul)
	s.Equal(s.summer, jul.AlbumId)
	s.Equal([]string{"favorite"}, tagStringsSorted(jul.Tags))
}

func tagStringsSorted(tags []catalog.Tag) []string {
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		result = append(result, tag.String())
	}
	sort.Strings(result)
	return result
}
