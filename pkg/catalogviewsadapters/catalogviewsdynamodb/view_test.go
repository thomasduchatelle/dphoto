package catalogviewsdynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"testing"
)

const (
	OwnerAvailability   = "OWNED"
	VisitorAvailability = "VISITOR"
)

func TestAlbumViewRepository_InsertAlbumSize(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)
	userId1 := usermodel.NewUserId("user-1")
	userId2 := usermodel.NewUserId("user-2")
	userId3 := usermodel.NewUserId("user-3")
	albumId1 := catalog.AlbumId{
		Owner:      "owner1",
		FolderName: catalog.NewFolderName("album-1"),
	}
	albumId2 := catalog.AlbumId{
		Owner:      "owner2",
		FolderName: catalog.NewFolderName("album-2"),
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
						Users:      []catalogviews.Availability{catalogviews.OwnerAvailability(userId1)},
					},
				},
			},
			before: nil,
			after: []map[string]types.AttributeValue{
				albumSizeItem(userId1, "OWNED", albumId1, "42"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should save the album size for the shared user",
			args: args{
				albumSizes: []catalogviews.AlbumSize{
					{
						AlbumId:    albumId1,
						MediaCount: 42,
						Users:      []catalogviews.Availability{catalogviews.VisitorAvailability(userId1)},
					},
				},
			},
			before: nil,
			after: []map[string]types.AttributeValue{
				albumSizeItem(userId1, "VISITOR", albumId1, "42"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should save several album sizes for multiple users",
			args: args{
				albumSizes: []catalogviews.AlbumSize{
					{
						AlbumId:    albumId1,
						MediaCount: 42,
						Users:      []catalogviews.Availability{catalogviews.OwnerAvailability(userId1), catalogviews.VisitorAvailability(userId2), catalogviews.VisitorAvailability(userId3)},
					},
					{
						AlbumId:    albumId2,
						MediaCount: 24,
						Users:      []catalogviews.Availability{catalogviews.OwnerAvailability(userId2), catalogviews.VisitorAvailability(userId1)},
					},
				},
			},
			before: nil,
			after: []map[string]types.AttributeValue{
				albumSizeItem(userId1, "OWNED", albumId1, "42"),
				albumSizeItem(userId1, "VISITOR", albumId2, "24"),
				albumSizeItem(userId2, "OWNED", albumId2, "24"),
				albumSizeItem(userId2, "VISITOR", albumId1, "42"),
				albumSizeItem(userId3, "VISITOR", albumId1, "42"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should update the count if the album size already exists",
			args: args{
				albumSizes: []catalogviews.AlbumSize{
					{
						AlbumId:    albumId1,
						MediaCount: 42,
						Users:      []catalogviews.Availability{catalogviews.OwnerAvailability(userId1)},
					},
				},
			},
			before: []map[string]types.AttributeValue{
				albumSizeItem(userId1, "OWNED", albumId1, "24"),
			},
			after: []map[string]types.AttributeValue{
				albumSizeItem(userId1, "OWNED", albumId1, "42"),
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

func TestAlbumViewRepository_DeleteAlbumSize(t *testing.T) {
	dyn := dynamotestutils.NewTestContext(context.Background(), t)
	const visitorType = "VISITOR"

	userId1 := usermodel.NewUserId("user-1")
	albumId1 := catalog.AlbumId{
		Owner:      "owner1",
		FolderName: catalog.NewFolderName("/album-1"),
	}

	type args struct {
		ctx          context.Context
		availability catalogviews.Availability
		albumId      catalog.AlbumId
	}
	tests := []struct {
		name      string
		args      args
		before    []map[string]types.AttributeValue
		wantAfter []map[string]types.AttributeValue
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should do nothing if the counts didn't exist",
			args: args{
				ctx:          context.Background(),
				availability: catalogviews.VisitorAvailability(userId1),
				albumId:      albumId1,
			},
			before:    nil,
			wantAfter: nil,
			wantErr:   assert.NoError,
		},
		{
			name: "it should not delete the count if it is for another user",
			args: args{
				ctx:          context.Background(),
				availability: catalogviews.VisitorAvailability(userId1),
				albumId:      albumId1,
			},
			before: []map[string]types.AttributeValue{
				albumSizeItem(usermodel.NewUserId("user-2"), visitorType, albumId1, "42"),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSizeItem(usermodel.NewUserId("user-2"), visitorType, albumId1, "42"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not delete the count if it is for another album",
			args: args{
				ctx:          context.Background(),
				availability: catalogviews.VisitorAvailability(userId1),
				albumId:      albumId1,
			},
			before: []map[string]types.AttributeValue{
				albumSizeItem(userId1, visitorType, catalog.AlbumId{
					Owner:      "owner1",
					FolderName: catalog.NewFolderName("/album-2"),
				}, "42"),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSizeItem(userId1, visitorType, catalog.AlbumId{
					Owner:      "owner1",
					FolderName: catalog.NewFolderName("/album-2"),
				}, "42"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not delete the count if it is for the owner",
			args: args{
				ctx:          context.Background(),
				availability: catalogviews.VisitorAvailability(userId1),
				albumId:      albumId1,
			},
			before: []map[string]types.AttributeValue{
				albumSizeItem(userId1, "OWNED", albumId1, "42"),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSizeItem(userId1, "OWNED", albumId1, "42"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should delete the count if it exists",
			args: args{
				ctx:          context.Background(),
				availability: catalogviews.VisitorAvailability(userId1),
				albumId:      albumId1,
			},
			before: []map[string]types.AttributeValue{
				albumSizeItem(userId1, visitorType, albumId1, "42"),
			},
			wantAfter: nil,
			wantErr:   assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn = dyn.Subtest(t)

			err := dyn.WithDbContent(dyn.Ctx, tt.before)
			if !assert.NoError(t, err) {
				return
			}

			a := &AlbumViewRepository{
				Client:    dyn.Client,
				TableName: dyn.Table,
			}
			err = a.DeleteAlbumSize(tt.args.ctx, tt.args.availability, tt.args.albumId)
			if tt.wantErr(t, err, fmt.Sprintf("DeleteAlbumSize(%v, %v, %v)", tt.args.ctx, tt.args.availability, tt.args.albumId)) {
				dyn.MustBool(dyn.EqualContent(tt.args.ctx, tt.wantAfter))
			}
		})
	}
}

func TestAlbumViewRepository_GetAvailabilitiesByUser(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)
	userId1 := usermodel.NewUserId("user-1")
	userId2 := usermodel.NewUserId("user-2")
	albumId1 := catalog.AlbumId{
		Owner:      "owner1",
		FolderName: catalog.NewFolderName("album-1"),
	}
	albumId2 := catalog.AlbumId{
		Owner:      "owner2",
		FolderName: catalog.NewFolderName("album-2"),
	}
	albumId3 := catalog.AlbumId{
		Owner:      "owner3",
		FolderName: catalog.NewFolderName("album-3"),
	}

	type args struct {
		user usermodel.UserId
	}
	tests := []struct {
		name    string
		args    args
		before  []map[string]types.AttributeValue
		want    []catalogviews.AlbumSize
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return an empty list if there is no album size",
			args: args{
				user: userId1,
			},
			before:  nil,
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should return the album sizes for the owner and visitor",
			args: args{
				user: userId1,
			},
			before: []map[string]types.AttributeValue{
				albumSizeItem(userId1, "OWNED", albumId1, "42"),
				albumSizeItem(userId1, "VISITOR", albumId2, "10"),
				albumSizeItem(userId2, "OWNED", albumId3, "5"),
			},
			want: []catalogviews.AlbumSize{
				{
					AlbumId:    albumId1,
					MediaCount: 42,
				},
				{
					AlbumId:    albumId2,
					MediaCount: 10,
				},
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

			a := &AlbumViewRepository{
				Client:    dyn.Client,
				TableName: dyn.Table,
			}
			got, err := a.GetAvailabilitiesByUser(ctx, tt.args.user)
			if !tt.wantErr(t, err, fmt.Sprintf("GetAvailabilitiesByUser(%v, %v)", ctx, tt.args.user)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAvailabilitiesByUser(%v, %v)", ctx, tt.args.user)
		})
	}
}

func TestAlbumViewRepository_UpdateAlbumSize(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)

	userId1 := usermodel.NewUserId("user-1")
	userId2 := usermodel.NewUserId("user-2")
	owner1 := ownermodel.Owner("owner1")
	owner2 := ownermodel.Owner("owner2")
	albumId1 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("/album-1")}
	albumId2 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("/album-2")}
	albumId3 := catalog.AlbumId{Owner: owner2, FolderName: catalog.NewFolderName("/album-3")}

	type args struct {
		ctx        context.Context
		albumSizes []catalogviews.AlbumSizeDiff
	}
	tests := []struct {
		name      string
		args      args
		before    []map[string]types.AttributeValue
		wantAfter []map[string]types.AttributeValue
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should do nothing if there is no album size",
			args: args{
				ctx:        ctx,
				albumSizes: nil,
			},
			before:    nil,
			wantAfter: nil,
			wantErr:   assert.NoError,
		},
		{
			name: "it should update the count for the owner",
			args: args{
				ctx: ctx,
				albumSizes: []catalogviews.AlbumSizeDiff{
					{
						AlbumId:        albumId1,
						Users:          []catalogviews.Availability{catalogviews.OwnerAvailability(userId1)},
						MediaCountDiff: 2,
					},
				},
			},
			before: []map[string]types.AttributeValue{
				albumSizeItem(userId1, OwnerAvailability, albumId1, "42"),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSizeItem(userId1, OwnerAvailability, albumId1, "44"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should create the count if it doesn't exist yet",
			args: args{
				ctx: ctx,
				albumSizes: []catalogviews.AlbumSizeDiff{
					{
						AlbumId:        albumId1,
						Users:          []catalogviews.Availability{catalogviews.OwnerAvailability(userId1)},
						MediaCountDiff: 2,
					},
				},
			},
			before: nil,
			wantAfter: []map[string]types.AttributeValue{
				albumSizeItem(userId1, OwnerAvailability, albumId1, "2"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should support a mix of albums and users that exists and don't",
			args: args{
				ctx: ctx,
				albumSizes: []catalogviews.AlbumSizeDiff{
					{
						AlbumId:        albumId1,
						Users:          []catalogviews.Availability{catalogviews.OwnerAvailability(userId1), catalogviews.VisitorAvailability(userId2)},
						MediaCountDiff: 2,
					},
					{
						AlbumId:        albumId2,
						Users:          []catalogviews.Availability{catalogviews.OwnerAvailability(userId1)},
						MediaCountDiff: 3,
					},
					{
						AlbumId:        albumId3,
						Users:          []catalogviews.Availability{catalogviews.OwnerAvailability(userId2), catalogviews.VisitorAvailability(userId1)},
						MediaCountDiff: 5,
					},
				},
			},
			before: []map[string]types.AttributeValue{
				albumSizeItem(userId1, OwnerAvailability, albumId1, "2"),
				albumSizeItem(userId2, VisitorAvailability, albumId1, "2"),
				albumSizeItem(userId1, OwnerAvailability, albumId2, "3"),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSizeItem(userId1, OwnerAvailability, albumId1, "4"),
				albumSizeItem(userId1, OwnerAvailability, albumId2, "6"),
				albumSizeItem(userId1, VisitorAvailability, albumId3, "5"),
				albumSizeItem(userId2, OwnerAvailability, albumId3, "5"),
				albumSizeItem(userId2, VisitorAvailability, albumId1, "4"),
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn = dyn.Subtest(t)

			err := dyn.WithDbContent(ctx, tt.before)
			if !assert.NoError(t, err) {
				return
			}

			a := &AlbumViewRepository{
				Client:    dyn.Client,
				TableName: dyn.Table,
			}
			err = a.UpdateAlbumSize(tt.args.ctx, tt.args.albumSizes)
			if tt.wantErr(t, err, fmt.Sprintf("UpdateAlbumSize(%v, %v)", tt.args.ctx, tt.args.albumSizes)) {
				dyn.MustBool(dyn.EqualContent(tt.args.ctx, tt.wantAfter))
			}
		})
	}
}
