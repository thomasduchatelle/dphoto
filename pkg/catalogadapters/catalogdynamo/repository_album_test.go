package catalogdynamo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"testing"
	"time"
)

type AlbumCrudTestSuite struct {
	suite.Suite
	suffix string
	repo   *Repository
	owner  ownermodel.Owner
}

func TestRepositoryAlbum(t *testing.T) {
	suite.Run(t, new(AlbumCrudTestSuite))
}

func (a *AlbumCrudTestSuite) SetupSuite() {
	dyn := dynamotestutils.NewTestContext(context.Background(), a.T())
	err := appdynamodb.CreateTableIfNecessary(context.Background(), dyn.Table, dyn.Client, true)
	if !assert.NoError(a.T(), err) {
		assert.FailNow(a.T(), err.Error())
	}

	a.owner = "UNITTEST#1"
	a.repo = &Repository{
		client: dyn.Client,
		table:  dyn.Table,
	}
}

func (a *AlbumCrudTestSuite) TestInsertAndFind() {
	folderName := catalog.NewFolderName("/christmas")

	err := a.repo.InsertAlbum(context.TODO(), catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: folderName,
		},
		Name:  "Christmas",
		Start: mustParseDate("2020-12-24"),
		End:   mustParseDate("2020-12-26"),
	})
	if !a.NoError(err, "it should insert a new album in DB") {
		return
	}

	name := "it should find previously saved album"
	found, err := a.repo.FindAlbumByIds(context.TODO(), catalog.AlbumId{Owner: a.owner, FolderName: folderName})
	if a.NoError(err, name) && a.Len(found, 1, name) {
		a.Equal(&catalog.Album{
			AlbumId: catalog.AlbumId{
				Owner:      a.owner,
				FolderName: folderName,
			},
			Name:  "Christmas",
			Start: mustParseDate("2020-12-24"),
			End:   mustParseDate("2020-12-26"),
		}, found[0], name)
	}
}

func (a *AlbumCrudTestSuite) TestInsertTwiceFails() {
	folderName := catalog.NewFolderName("New Year")

	err := a.repo.InsertAlbum(context.TODO(), catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: folderName,
		},
		Name:  "New Year",
		Start: mustParseDate("2020-12-31"),
		End:   mustParseDate("2021-01-01"),
	})
	if !a.NoError(err, "it should insert a new album in DB") {
		return
	}

	err = a.repo.InsertAlbum(context.TODO(), catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: folderName,
		},
		Name:  "New Year Again",
		Start: mustParseDate("2020-12-31"),
		End:   mustParseDate("2021-01-01"),
	})
	log.WithField("Error", err).Infoln("insert twice fails")
	a.Error(err, "it should fail to override an existing album")
}

func (a *AlbumCrudTestSuite) TestFindNotFound() {
	ttName := "it should return [?, AlbumNotFoundError] when searched album do not exists"
	albums, err := a.repo.FindAlbumByIds(context.TODO(), catalog.AlbumId{Owner: a.owner, FolderName: "_donotexist"})
	if a.NoError(err, ttName) {
		a.Empty(albums)
	}
}

func TestRepository_CountMediasBySelectors(t *testing.T) {
	const owner = "ironman"
	mediaId1 := "media-1"
	mediaId2 := "media-2"
	jan24 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	may24 := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	jun24 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	may24the8th := "2024-05-08"
	feb24the12th := "2024-02-12"
	album1 := catalog.AlbumId{
		Owner:      owner,
		FolderName: catalog.NewFolderName("/album-1"),
	}
	album2 := catalog.AlbumId{
		Owner:      owner,
		FolderName: catalog.NewFolderName("/album-2"),
	}

	type args struct {
		owner     ownermodel.Owner
		selectors []catalog.MediaSelector
	}
	tests := []struct {
		name          string
		withDbContent []map[string]types.AttributeValue
		args          args
		want          int
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "it should return 0 when no selector is provided",
			args: args{
				owner:     owner,
				selectors: nil,
			},
			want:    0,
			wantErr: assert.NoError,
		},
		{
			name: "it should return 1 when only one media is in the album",
			withDbContent: []map[string]types.AttributeValue{
				mediaEntry(album1, mediaId1, may24the8th),
			},
			args: args{
				owner: owner,
				selectors: []catalog.MediaSelector{
					{
						FromAlbums: []catalog.AlbumId{album1},
						Start:      may24,
						End:        jun24,
					},
				},
			},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name: "it should return 1 when one media is in the date range, not the other",
			withDbContent: []map[string]types.AttributeValue{
				mediaEntry(album1, mediaId1, may24the8th),
				mediaEntry(album1, mediaId2, feb24the12th),
			},
			args: args{
				owner: owner,
				selectors: []catalog.MediaSelector{
					{
						FromAlbums: []catalog.AlbumId{album1},
						Start:      may24,
						End:        jun24,
					},
				},
			},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name: "it should return 2, one from each selector on different albums",
			withDbContent: []map[string]types.AttributeValue{
				mediaEntry(album1, mediaId1, may24the8th),
				mediaEntry(album2, mediaId2, feb24the12th),
			},
			args: args{
				owner: owner,
				selectors: []catalog.MediaSelector{
					{
						FromAlbums: []catalog.AlbumId{album1},
						Start:      may24,
						End:        jun24,
					},
					{
						FromAlbums: []catalog.AlbumId{album2},
						Start:      jan24,
						End:        may24,
					},
				},
			},
			want:    2,
			wantErr: assert.NoError,
		},
		{
			name: "it should return 2 from the first selector from 2 albums",
			withDbContent: []map[string]types.AttributeValue{
				mediaEntry(album1, mediaId1, may24the8th),
				mediaEntry(album2, mediaId2, feb24the12th),
			},
			args: args{
				owner: owner,
				selectors: []catalog.MediaSelector{
					{
						FromAlbums: []catalog.AlbumId{album1, album2},
						Start:      jan24,
						End:        jun24,
					},
				},
			},
			want:    2,
			wantErr: assert.NoError,
		},
	}

	dyn := dynamotestutils.NewTestContext(context.Background(), t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn = dyn.Subtest(t)

			err := dyn.WithDbContent(dyn.Ctx, tt.withDbContent)
			if !assert.NoError(t, err, "WithDbContent") {
				return
			}

			r := &Repository{
				client: dyn.Client,
				table:  dyn.Table,
			}
			got, err := r.CountMediasBySelectors(dyn.Ctx, tt.args.owner, tt.args.selectors)
			if !tt.wantErr(t, err, fmt.Sprintf("CountMediasBySelectors(%v, %v)", tt.args.owner, tt.args.selectors)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CountMediasBySelectors(%v, %v)", tt.args.owner, tt.args.selectors)
		})
	}
}

func mediaEntry(albumId catalog.AlbumId, mediaId1 string, mediaDateTime string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"PK":           &types.AttributeValueMemberS{Value: fmt.Sprintf("%s#MEDIA#%s", albumId.Owner, mediaId1)},
		"SK":           &types.AttributeValueMemberS{Value: "#METADATA"},
		"AlbumIndexPK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s#%s", albumId.Owner, albumId.FolderName)},
		"AlbumIndexSK": &types.AttributeValueMemberS{Value: fmt.Sprintf("MEDIA#%s#%s", mediaDateTime, mediaId1)},
	}
}

func TestRepository_AmendDates(t *testing.T) {
	albumId1 := catalog.AlbumId{
		Owner:      "ironman",
		FolderName: catalog.NewFolderName("/album-1"),
	}
	album1PK := &types.AttributeValueMemberS{Value: "ironman#ALBUM"}
	album1SK := &types.AttributeValueMemberS{Value: "ALBUM#/album-1"}

	dyn := dynamotestutils.NewTestContext(context.Background(), t)

	type args struct {
		albumId catalog.AlbumId
		start   time.Time
		end     time.Time
	}
	tests := []struct {
		name    string
		args    args
		before  []map[string]types.AttributeValue
		after   []map[string]types.AttributeValue
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should update an album that exists",
			args: args{
				albumId: albumId1,
				start:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				end:     time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			},
			before: []map[string]types.AttributeValue{
				{
					"PK":         album1PK,
					"SK":         album1SK,
					"AlbumStart": &types.AttributeValueMemberS{Value: "2024-02-01T00:00:00Z"},
					"AlbumEnd":   &types.AttributeValueMemberS{Value: "2024-05-01T00:00:00Z"},
				},
			},
			after: []map[string]types.AttributeValue{
				{
					"PK":         album1PK,
					"SK":         album1SK,
					"AlbumStart": &types.AttributeValueMemberS{Value: "2024-01-01T00:00:00Z"},
					"AlbumEnd":   &types.AttributeValueMemberS{Value: "2024-06-01T00:00:00Z"},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should raise an error if the album doesn't exists",
			args: args{
				albumId: albumId1,
				start:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				end:     time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			},
			before: nil,
			after:  nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundError, i...)
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

			r := &Repository{
				client: dyn.Client,
				table:  dyn.Table,
			}

			err = r.AmendDates(context.Background(), tt.args.albumId, tt.args.start, tt.args.end)
			if tt.wantErr(t, err, fmt.Sprintf("AmendDates(%v, %v, %v, %v)", context.Background(), tt.args.albumId, tt.args.start, tt.args.end)) {
				_, err := dyn.EqualContent(dyn.Ctx, tt.after)
				assert.NoError(t, err, "AssertDbContent")
			}
		})
	}
}

func TestRepository_UpdateAlbumName(t *testing.T) {
	albumId1 := catalog.AlbumId{
		Owner:      "ironman",
		FolderName: catalog.NewFolderName("/album-1"),
	}
	album1PK := &types.AttributeValueMemberS{Value: "ironman#ALBUM"}
	album1SK := &types.AttributeValueMemberS{Value: "ALBUM#/album-1"}

	dyn := dynamotestutils.NewTestContext(context.Background(), t)

	type args struct {
		albumId catalog.AlbumId
		newName string
	}
	tests := []struct {
		name    string
		args    args
		before  []map[string]types.AttributeValue
		after   []map[string]types.AttributeValue
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should update the name of an album that exists",
			args: args{
				albumId: albumId1,
				newName: "New Name",
			},
			before: []map[string]types.AttributeValue{
				{
					"PK":        album1PK,
					"SK":        album1SK,
					"AlbumName": &types.AttributeValueMemberS{Value: "Old Name"},
				},
			},
			after: []map[string]types.AttributeValue{
				{
					"PK":        album1PK,
					"SK":        album1SK,
					"AlbumName": &types.AttributeValueMemberS{Value: "New Name"},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should raise an error if the media doesn't exists",
			args: args{
				albumId: albumId1,
				newName: "New Name",
			},
			before: nil,
			after:  nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.MediaNotFoundError, i...)
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

			r := &Repository{
				client: dyn.Client,
				table:  dyn.Table,
			}

			err = r.UpdateAlbumName(context.Background(), tt.args.albumId, tt.args.newName)
			if tt.wantErr(t, err, fmt.Sprintf("UpdateAlbumName(%v, %v, %v)", context.Background(), tt.args.albumId, tt.args.newName)) {
				_, err := dyn.EqualContent(dyn.Ctx, tt.after)
				assert.NoError(t, err, "AssertDbContent")
			}
		})
	}
}
