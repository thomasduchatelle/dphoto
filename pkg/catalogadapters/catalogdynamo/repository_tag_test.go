package catalogdynamo

import (
	"context"
	"slices"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/suite"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
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
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays.Owner, s.mediaInJan, s.favorite))

	s.Equal([]string{"favorite"}, s.tagsOfMediaInAlbum(s.holidays, s.mediaInJan),
		"the album medias should carry the tag without a per-media query")
}

func (s *TagTestSuite) TestAddTag_isIdempotent() {
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays.Owner, s.mediaInJan, s.favorite))
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays.Owner, s.mediaInJan, s.favorite))

	s.Equal([]string{"favorite"}, s.tagsOfMediaInAlbum(s.holidays, s.mediaInJan))
}

func (s *TagTestSuite) TestAddTag_rejectsAMediaThatIsNotOwned() {
	err := s.repo.AddTag(context.Background(), s.holidays.Owner, catalog.MediaId("unknown-media"), s.favorite)
	s.ErrorIs(err, catalog.MediaNotFoundError, "tagging a media without a #METADATA row must be rejected")
}

func (s *TagTestSuite) TestRemoveTag_leavesNoResidualItem() {
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays.Owner, s.mediaInJan, s.favorite))
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays.Owner, s.mediaInJan, s.peopleLuke))

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

func (s *TagTestSuite) TestAddTag_concurrentTagsOnTheSameMediaAllSurvive() {
	start := make(chan struct{})
	tags := []catalog.Tag{s.favorite, s.peopleLuke, mustTag("people:leia"), mustTag("places:tatooine")}

	var wg sync.WaitGroup
	errs := make([]error, len(tags))
	for i, tag := range tags {
		wg.Add(1)
		go func(i int, tag catalog.Tag) {
			defer wg.Done()
			<-start
			errs[i] = s.repo.AddTag(context.Background(), s.holidays.Owner, s.mediaInJan, tag)
		}(i, tag)
	}
	close(start)
	wg.Wait()

	for _, err := range errs {
		s.Require().NoError(err)
	}
	s.Equal([]string{"favorite", "people:leia", "people:luke", "places:tatooine"}, s.tagsOfMediaInAlbum(s.holidays, s.mediaInJan),
		"the String Set ADD must merge concurrent writes, no tag is clobbered")
}

func (s *TagTestSuite) TestFindMediasByTag_dropsAMediaWhoseMetadataNoLongerCarriesTheTag() {
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays.Owner, s.mediaInJan, s.favorite))

	s.givenOrphanTagSibling(s.holidays.Owner, s.mediaInJul, s.favorite)

	medias, err := s.repo.FindMediasByTag(context.Background(), s.holidays.Owner, s.favorite)
	s.Require().NoError(err)

	s.Require().Len(medias, 1, "only the media whose #METADATA source of truth carries the tag is returned")
	s.Equal(s.mediaInJan, medias[0].Id)
}

// givenOrphanTagSibling writes a TAG# sibling (and its GSI entry) for a media whose #METADATA
// does not carry the tag, reproducing index drift (e.g. a partial-failure during RemoveTag).
func (s *TagTestSuite) givenOrphanTagSibling(owner ownermodel.Owner, mediaId catalog.MediaId, tag catalog.Tag) {
	item, err := marshalTag(owner, mediaId, tag, 1)
	s.Require().NoError(err)
	_, err = s.repo.client.PutItem(context.Background(), &dynamodb.PutItemInput{Item: item, TableName: &s.repo.table})
	s.Require().NoError(err)
}

func (s *TagTestSuite) TestAddTag_siblingSurvivesAStaleConcurrentRemove() {
	// Deterministic reconstruction of the dangerous interleave of a concurrent AddTag +
	// RemoveTag of the same (media, tag): RemoveTag's set-DELETE serialises FIRST and mints the
	// lower version (removeVersion), then AddTag fully completes with the higher version, and
	// finally RemoveTag's stale sibling-delete runs LAST. Without the TagVersion guard this last
	// delete would strand the tag in #METADATA while erasing it from the index, making a real tag
	// invisible on the Tags page. The guard must reject it.
	removeVersion, err := s.repo.updateMediaTags(context.Background(), s.holidays.Owner, s.mediaInJan,
		expression.Delete(expression.Name("Tags"), expression.Value(stringSet(s.favorite))))
	s.Require().NoError(err)

	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays.Owner, s.mediaInJan, s.favorite))

	s.Require().NoError(s.repo.deleteTagSibling(context.Background(), s.holidays.Owner, s.mediaInJan, s.favorite, removeVersion),
		"a stale delete losing the version race must be a no-op success")

	s.True(slices.Contains(s.tagsOfMediaInAlbum(s.holidays, s.mediaInJan), s.favorite.String()),
		"the tag won the set, so it is in the #METADATA source of truth")
	medias, err := s.repo.FindMediasByTag(context.Background(), s.holidays.Owner, s.favorite)
	s.Require().NoError(err)
	s.True(slices.ContainsFunc(medias, func(m *catalog.MediaWithAlbum) bool { return m.Id == s.mediaInJan }),
		"the stale delete must NOT have erased the sibling written by the newer AddTag")
}

func (s *TagTestSuite) TestPutTagSibling_staleWriteDoesNotOverwriteANewerSibling() {
	// A stale AddTag's sibling-put (low version) must not clobber a sibling already written by a
	// newer operation (high version). This exercises the TagVersion < :v REJECT branch of the
	// guard directly: the sibling is present at the newer version when the stale put runs.
	staleVersion, err := s.repo.updateMediaTags(context.Background(), s.holidays.Owner, s.mediaInJan,
		expression.Add(expression.Name("Tags"), expression.Value(stringSet(s.favorite))))
	s.Require().NoError(err)

	newerVersion, err := s.repo.updateMediaTags(context.Background(), s.holidays.Owner, s.mediaInJan,
		expression.Add(expression.Name("Tags"), expression.Value(stringSet(s.favorite))))
	s.Require().NoError(err)
	s.Require().NoError(s.repo.putTagSibling(context.Background(), s.holidays.Owner, s.mediaInJan, s.favorite, newerVersion))

	s.Require().NoError(s.repo.putTagSibling(context.Background(), s.holidays.Owner, s.mediaInJan, s.favorite, staleVersion),
		"a stale put losing the version race must be a no-op success")

	s.Equal(newerVersion, s.tagSiblingVersion(s.holidays.Owner, s.mediaInJan, s.favorite),
		"the stale put must NOT have overwritten the sibling stamped by the newer operation")
}

// tagSiblingVersion reads the TagVersion stamped on a media's TAG# sibling (-1 if absent).
func (s *TagTestSuite) tagSiblingVersion(owner ownermodel.Owner, mediaId catalog.MediaId, tag catalog.Tag) int64 {
	key, err := attributevalue.MarshalMap(TagPrimaryKey(owner, mediaId, tag))
	s.Require().NoError(err)
	out, err := s.repo.client.GetItem(context.Background(), &dynamodb.GetItemInput{Key: key, TableName: &s.repo.table})
	s.Require().NoError(err)
	if len(out.Item) == 0 {
		return -1
	}
	var record TagRecord
	s.Require().NoError(attributevalue.UnmarshalMap(out.Item, &record))
	return record.TagVersion
}

func (s *TagTestSuite) TestFindMediasByTag_returnsEveryMatchingMediaWithItsAlbumAndTagsAcrossAlbums() {
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays.Owner, s.mediaInJan, s.favorite))
	s.Require().NoError(s.repo.AddTag(context.Background(), s.holidays.Owner, s.mediaInJan, s.peopleLuke))
	s.Require().NoError(s.repo.AddTag(context.Background(), s.summer.Owner, s.mediaInJul, s.favorite))

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
