package catalogviewsdynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"testing"
)

func TestAlbumViewRepository_InsertAlbumSize(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)
	userId := usermodel.NewUserId("user-1")
	albumId1 := catalog.AlbumId{
		Owner:      "owner",
		FolderName: catalog.NewFolderName("album-1"),
	}

	type args struct {
		albumSizes []catalogviews.AlbumSize
	}
	tests := []struct {
		name    string
		args    args
		before  []map[string]types.AttributeValue
		after   []map[string]types.AttributeValue
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should not save anything if content is empty",
			args: args{
				albumSizes: nil,
			},
			before:  nil,
			after:   nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should not save anything if there is no user on the album size",
			args: args{
				albumSizes: []catalogviews.AlbumSize{
					{
						AlbumId:    albumId1,
						MediaCount: 42,
						Users:      nil,
					},
				},
			},
			before:  nil,
			after:   nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should save the album size for the owner",
			args: args{
				albumSizes: []catalogviews.AlbumSize{
					{
						AlbumId:    albumId1,
						MediaCount: 42,
						Users:      []catalogviews.Availability{catalogviews.OwnerAvailability(userId)},
					},
				},
			},
			before: nil,
			after: []map[string]types.AttributeValue{
				albumSizeItem(userId, "OWNED", albumId1, "42"),
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn := dyn.Subtest(t)

			if !assert.NoError(t, dyn.WithDbContent(ctx, tt.before)) {
				return
			}

			repository := &AlbumViewRepository{
				Client:    dyn.Client,
				TableName: dyn.Table,
			}
			err := repository.InsertAlbumSize(ctx, tt.args.albumSizes)
			if !tt.wantErr(t, err) {
				return
			}

			_, err = dyn.EqualContent(ctx, tt.after)
			assert.NoError(t, err)
		})
	}
}

func albumSizeItem(user usermodel.UserId, accessType string, albumId catalog.AlbumId, count string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"PK":              &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s#ALBUMS_VIEW", user)},
		"SK":              &types.AttributeValueMemberS{Value: fmt.Sprintf("%s#%s#%s#COUNT", accessType, albumId.Owner.Value(), albumId.FolderName.String())},
		"AlbumOwner":      &types.AttributeValueMemberS{Value: albumId.Owner.Value()},
		"AlbumFolderName": &types.AttributeValueMemberS{Value: albumId.FolderName.String()},
		"Count":           &types.AttributeValueMemberN{Value: count},
	}
}
