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
	"testing"
	"time"
)

type AlbumCrudTestSuite struct {
	suite.Suite
	suffix string
	repo   *Repository
	owner  catalog.Owner
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

func (a *AlbumCrudTestSuite) TestDeleteEmpty() {
	folderName := catalog.NewFolderName("ToBeDeleted")
	albumId := catalog.AlbumId{
		Owner:      a.owner,
		FolderName: folderName,
	}

	err := a.repo.InsertAlbum(context.TODO(), catalog.Album{
		AlbumId: albumId,
		Name:    "ToBeDeleted",
		Start:   mustParseDate("2020-12-24"),
		End:     mustParseDate("2020-12-26"),
	})
	if !a.NoError(err, "it should insert an album to delete") {
		return
	}

	err = a.repo.DeleteEmptyAlbum(context.TODO(), albumId)
	a.NoError(err, "it should delete an album that do not have any medias")
}

func (a *AlbumCrudTestSuite) TestUpdate() {
	folderName := "Update1"

	err := a.repo.InsertAlbum(context.TODO(), catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: catalog.NewFolderName("Update1"),
		},
		Name:  folderName,
		Start: mustParseDate("2020-12-01"),
		End:   mustParseDate("2021-01-31"),
	})
	if !a.NoError(err, "it should insert an album to update") {
		return
	}

	update := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: catalog.NewFolderName("Update1"),
		},
		Name:  "Another Name",
		Start: mustParseDate("2021-01-01"),
		End:   mustParseDate("2021-02-01"),
	}
	err = a.repo.UpdateAlbum(context.TODO(), update)
	name := "it should update an exiting album"
	if a.NoError(err, name) {
		updated, err := a.repo.FindAlbumByIds(context.TODO(), catalog.AlbumId{Owner: a.owner, FolderName: catalog.NewFolderName(folderName)})
		if a.NoError(err, name) && a.Len(updated, 1) {
			a.Equal(&update, updated[0], name)
		}
	}
}

func (a *AlbumCrudTestSuite) TestUpdateNotExisting() {
	folderName := catalog.NewFolderName("_do_not_exist")

	update := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      a.owner,
			FolderName: folderName,
		},
		Name:  "Another Name",
		Start: mustParseDate("2021-01-01"),
		End:   mustParseDate("2021-02-01"),
	}
	err := a.repo.UpdateAlbum(context.TODO(), update)
	a.Error(err, "it should fail to update an album that do not exist.")
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
		owner     catalog.Owner
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
