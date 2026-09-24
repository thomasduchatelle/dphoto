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
	"time"
)

const (
	OwnerAvailability   = "OWNED"
	VisitorAvailability = "VISITOR"
)

var (
	displayFieldsAlbum1Start = time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	displayFieldsAlbum1End   = time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
)

func TestAlbumViewRepository_PutSummaries(t *testing.T) {
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
		summaries []catalogviews.AlbumSummaryForUsers
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
				summaries: nil,
			},
			before:  nil,
			after:   nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should not save anything if there is no user on the summary",
			args: args{
				summaries: []catalogviews.AlbumSummaryForUsers{
					{
						AlbumSummary: catalogviews.AlbumSummary{AlbumId: albumId1, MediaCount: 42},
						Users:        nil,
					},
				},
			},
			before:  nil,
			after:   nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should save the summary for the owner",
			args: args{
				summaries: []catalogviews.AlbumSummaryForUsers{
					{
						AlbumSummary: catalogviews.AlbumSummary{AlbumId: albumId1, MediaCount: 42},
						Users:        []catalogviews.Availability{catalogviews.OwnerAvailability(userId1)},
					},
				},
			},
			before: nil,
			after: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, "OWNED", albumId1).withCount(42).build(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should save the summary with display fields when they are set",
			args: args{
				summaries: []catalogviews.AlbumSummaryForUsers{
					{
						AlbumSummary: catalogviews.AlbumSummary{
							AlbumId:    albumId1,
							MediaCount: 42,
							Name:       "January 2024",
							Start:      displayFieldsAlbum1Start,
							End:        displayFieldsAlbum1End,
						},
						Users: []catalogviews.Availability{catalogviews.OwnerAvailability(userId1)},
					},
				},
			},
			before: nil,
			after: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, "OWNED", albumId1).
					withCount(42).
					withDisplayFields("January 2024", displayFieldsAlbum1Start, displayFieldsAlbum1End).
					build(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should save the summary for the shared user",
			args: args{
				summaries: []catalogviews.AlbumSummaryForUsers{
					{
						AlbumSummary: catalogviews.AlbumSummary{AlbumId: albumId1, MediaCount: 42},
						Users:        []catalogviews.Availability{catalogviews.VisitorAvailability(userId1)},
					},
				},
			},
			before: nil,
			after: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, "VISITOR", albumId1).withCount(42).build(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should save several summaries for multiple users",
			args: args{
				summaries: []catalogviews.AlbumSummaryForUsers{
					{
						AlbumSummary: catalogviews.AlbumSummary{AlbumId: albumId1, MediaCount: 42},
						Users:        []catalogviews.Availability{catalogviews.OwnerAvailability(userId1), catalogviews.VisitorAvailability(userId2), catalogviews.VisitorAvailability(userId3)},
					},
					{
						AlbumSummary: catalogviews.AlbumSummary{AlbumId: albumId2, MediaCount: 24},
						Users:        []catalogviews.Availability{catalogviews.OwnerAvailability(userId2), catalogviews.VisitorAvailability(userId1)},
					},
				},
			},
			before: nil,
			after: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, "OWNED", albumId1).withCount(42).build(),
				albumSummaryItemBuilder(userId1, "VISITOR", albumId2).withCount(24).build(),
				albumSummaryItemBuilder(userId2, "OWNED", albumId2).withCount(24).build(),
				albumSummaryItemBuilder(userId2, "VISITOR", albumId1).withCount(42).build(),
				albumSummaryItemBuilder(userId3, "VISITOR", albumId1).withCount(42).build(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should overwrite the row if the summary already exists",
			args: args{
				summaries: []catalogviews.AlbumSummaryForUsers{
					{
						AlbumSummary: catalogviews.AlbumSummary{AlbumId: albumId1, MediaCount: 42},
						Users:        []catalogviews.Availability{catalogviews.OwnerAvailability(userId1)},
					},
				},
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, "OWNED", albumId1).withCount(24).build(),
			},
			after: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, "OWNED", albumId1).withCount(42).build(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should delete a legacy #COUNT-suffixed row when writing the new row for the same user/album",
			args: args{
				summaries: []catalogviews.AlbumSummaryForUsers{
					{
						AlbumSummary: catalogviews.AlbumSummary{AlbumId: albumId1, MediaCount: 42},
						Users:        []catalogviews.Availability{catalogviews.OwnerAvailability(userId1)},
					},
				},
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, "OWNED", albumId1).withLegacyCountSuffix().withCount(24).build(),
			},
			after: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, "OWNED", albumId1).withCount(42).build(),
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
			err := repository.PutSummaries(ctx, tt.args.summaries)
			if !tt.wantErr(t, err) {
				return
			}

			_, err = dyn.EqualContent(ctx, tt.after)
			assert.NoError(t, err)
		})
	}
}

type summaryItemBuilder struct {
	item map[string]types.AttributeValue
}

func albumSummaryItemBuilder(user usermodel.UserId, accessType string, albumId catalog.AlbumId) *summaryItemBuilder {
	return &summaryItemBuilder{
		item: map[string]types.AttributeValue{
			"PK":               &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s#ALBUMS_VIEW", user)},
			"SK":               &types.AttributeValueMemberS{Value: fmt.Sprintf("%s#%s#%s", accessType, albumId.Owner.Value(), albumId.FolderName.String())},
			"AlbumOwner":       &types.AttributeValueMemberS{Value: albumId.Owner.Value()},
			"AlbumFolderName":  &types.AttributeValueMemberS{Value: albumId.FolderName.String()},
			"AvailabilityType": &types.AttributeValueMemberS{Value: accessType},
			"UserId":           &types.AttributeValueMemberS{Value: user.Value()},
			"AlbumViewIndexPK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ALBUM#%s#%s#ALBUMS_VIEW", albumId.Owner.Value(), albumId.FolderName.String())},
		},
	}
}

func (b *summaryItemBuilder) withCount(count int) *summaryItemBuilder {
	b.item["Count"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", count)}
	return b
}

func (b *summaryItemBuilder) withLegacyCountSuffix() *summaryItemBuilder {
	sk := b.item["SK"].(*types.AttributeValueMemberS).Value
	b.item["SK"] = &types.AttributeValueMemberS{Value: sk + "#COUNT"}
	delete(b.item, "AlbumViewIndexPK")
	return b
}

func (b *summaryItemBuilder) withDisplayFields(name string, start, end time.Time) *summaryItemBuilder {
	b.item["AlbumName"] = &types.AttributeValueMemberS{Value: name}
	b.item["AlbumStart"] = &types.AttributeValueMemberS{Value: start.UTC().Format(time.RFC3339)}
	b.item["AlbumEnd"] = &types.AttributeValueMemberS{Value: end.UTC().Format(time.RFC3339)}
	return b
}

func (b *summaryItemBuilder) build() map[string]types.AttributeValue {
	return b.item
}

func TestAlbumViewRepository_DeleteRow(t *testing.T) {
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
			name: "it should do nothing if the row didn't exist",
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
			name: "it should not delete the row if it is for another user",
			args: args{
				ctx:          context.Background(),
				availability: catalogviews.VisitorAvailability(userId1),
				albumId:      albumId1,
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(usermodel.NewUserId("user-2"), visitorType, albumId1).withCount(42).build(),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(usermodel.NewUserId("user-2"), visitorType, albumId1).withCount(42).build(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not delete the row if it is for another album",
			args: args{
				ctx:          context.Background(),
				availability: catalogviews.VisitorAvailability(userId1),
				albumId:      albumId1,
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, visitorType, catalog.AlbumId{
					Owner:      "owner1",
					FolderName: catalog.NewFolderName("/album-2"),
				}).withCount(42).build(),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, visitorType, catalog.AlbumId{
					Owner:      "owner1",
					FolderName: catalog.NewFolderName("/album-2"),
				}).withCount(42).build(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not delete the row if it is for the owner",
			args: args{
				ctx:          context.Background(),
				availability: catalogviews.VisitorAvailability(userId1),
				albumId:      albumId1,
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, "OWNED", albumId1).withCount(42).build(),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, "OWNED", albumId1).withCount(42).build(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should delete the row if it exists",
			args: args{
				ctx:          context.Background(),
				availability: catalogviews.VisitorAvailability(userId1),
				albumId:      albumId1,
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, visitorType, albumId1).withCount(42).build(),
			},
			wantAfter: nil,
			wantErr:   assert.NoError,
		},
		{
			name: "it should also delete the legacy #COUNT-suffixed row",
			args: args{
				ctx:          context.Background(),
				availability: catalogviews.VisitorAvailability(userId1),
				albumId:      albumId1,
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, visitorType, albumId1).withLegacyCountSuffix().withCount(42).build(),
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
			err = a.DeleteRow(tt.args.ctx, tt.args.availability, tt.args.albumId)
			if tt.wantErr(t, err, fmt.Sprintf("DeleteRow(%v, %v, %v)", tt.args.ctx, tt.args.availability, tt.args.albumId)) {
				dyn.MustBool(dyn.EqualContent(tt.args.ctx, tt.wantAfter))
			}
		})
	}
}

func TestAlbumViewRepository_ListSummariesForUser(t *testing.T) {
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
		want    []catalogviews.UserAlbumSummary
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return an empty list if there is no row",
			args: args{
				user: userId1,
			},
			before:  nil,
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should return the summaries for the owner and visitor",
			args: args{
				user: userId1,
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, "OWNED", albumId1).withCount(42).build(),
				albumSummaryItemBuilder(userId1, "VISITOR", albumId2).withCount(10).
					withDisplayFields("February 2024", displayFieldsAlbum1Start, displayFieldsAlbum1End).build(),
				albumSummaryItemBuilder(userId2, "OWNED", albumId3).withCount(5).build(),
			},
			want: []catalogviews.UserAlbumSummary{
				{Availability: catalogviews.OwnerAvailability(userId1), AlbumSummary: catalogviews.AlbumSummary{AlbumId: albumId1, MediaCount: 42}},
				{Availability: catalogviews.VisitorAvailability(userId1), AlbumSummary: catalogviews.AlbumSummary{AlbumId: albumId2, MediaCount: 10, Name: "February 2024", Start: displayFieldsAlbum1Start, End: displayFieldsAlbum1End}},
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
			got, err := a.ListSummariesForUser(ctx, tt.args.user)
			if !tt.wantErr(t, err, fmt.Sprintf("ListSummariesForUser(%v, %v)", ctx, tt.args.user)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ListSummariesForUser(%v, %v)", ctx, tt.args.user)
		})
	}
}

func TestAlbumViewRepository_IncrementCountForAllViewers(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)

	userId1 := usermodel.NewUserId("user-1")
	userId2 := usermodel.NewUserId("user-2")
	owner1 := ownermodel.Owner("owner1")
	albumId1 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("/album-1")}
	albumId2 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("/album-2")}

	type args struct {
		ctx     context.Context
		updates []catalogviews.AlbumCountDiff
	}
	tests := []struct {
		name      string
		args      args
		before    []map[string]types.AttributeValue
		wantAfter []map[string]types.AttributeValue
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should do nothing if there is no update",
			args: args{
				ctx:     ctx,
				updates: nil,
			},
			before:    nil,
			wantAfter: nil,
			wantErr:   assert.NoError,
		},
		{
			name: "it should do nothing when no viewer row exists for the album",
			args: args{
				ctx: ctx,
				updates: []catalogviews.AlbumCountDiff{
					{AlbumId: albumId1, MediaCountDiff: 2},
				},
			},
			before:    nil,
			wantAfter: nil,
			wantErr:   assert.NoError,
		},
		{
			name: "it should increment the count on every viewer row of the album",
			args: args{
				ctx: ctx,
				updates: []catalogviews.AlbumCountDiff{
					{AlbumId: albumId1, MediaCountDiff: 2},
				},
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).withCount(42).build(),
				albumSummaryItemBuilder(userId2, VisitorAvailability, albumId1).withCount(42).build(),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).withCount(44).build(),
				albumSummaryItemBuilder(userId2, VisitorAvailability, albumId1).withCount(44).build(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should NOT clobber the display fields on an existing row",
			args: args{
				ctx: ctx,
				updates: []catalogviews.AlbumCountDiff{
					{AlbumId: albumId1, MediaCountDiff: 3},
				},
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).
					withCount(10).
					withDisplayFields("January 2024", displayFieldsAlbum1Start, displayFieldsAlbum1End).
					build(),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).
					withCount(13).
					withDisplayFields("January 2024", displayFieldsAlbum1Start, displayFieldsAlbum1End).
					build(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should increment each album independently",
			args: args{
				ctx: ctx,
				updates: []catalogviews.AlbumCountDiff{
					{AlbumId: albumId1, MediaCountDiff: 2},
					{AlbumId: albumId2, MediaCountDiff: 3},
				},
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).withCount(2).build(),
				albumSummaryItemBuilder(userId2, VisitorAvailability, albumId1).withCount(2).build(),
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId2).withCount(3).build(),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).withCount(4).build(),
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId2).withCount(6).build(),
				albumSummaryItemBuilder(userId2, VisitorAvailability, albumId1).withCount(4).build(),
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
			err = a.IncrementCountForAllViewers(tt.args.ctx, tt.args.updates)
			if tt.wantErr(t, err, fmt.Sprintf("IncrementCountForAllViewers(%v, %v)", tt.args.ctx, tt.args.updates)) {
				dyn.MustBool(dyn.EqualContent(tt.args.ctx, tt.wantAfter))
			}
		})
	}
}

func TestAlbumViewRepository_SetCountForAllViewers(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)

	userId1 := usermodel.NewUserId("user-1")
	userId2 := usermodel.NewUserId("user-2")
	albumId1 := catalog.AlbumId{Owner: "owner1", FolderName: catalog.NewFolderName("/album-1")}

	type args struct {
		ctx     context.Context
		updates []catalogviews.AlbumCount
	}
	tests := []struct {
		name      string
		args      args
		before    []map[string]types.AttributeValue
		wantAfter []map[string]types.AttributeValue
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should do nothing when no viewer row exists for the album",
			args: args{
				ctx:     ctx,
				updates: []catalogviews.AlbumCount{{AlbumId: albumId1, MediaCount: 7}},
			},
			before:    nil,
			wantAfter: nil,
			wantErr:   assert.NoError,
		},
		{
			name: "it should set the count on every viewer row of the album",
			args: args{
				ctx:     ctx,
				updates: []catalogviews.AlbumCount{{AlbumId: albumId1, MediaCount: 7}},
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).withCount(42).build(),
				albumSummaryItemBuilder(userId2, VisitorAvailability, albumId1).withCount(42).build(),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).withCount(7).build(),
				albumSummaryItemBuilder(userId2, VisitorAvailability, albumId1).withCount(7).build(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should NOT clobber the display fields on an existing row",
			args: args{
				ctx:     ctx,
				updates: []catalogviews.AlbumCount{{AlbumId: albumId1, MediaCount: 7}},
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).
					withCount(10).
					withDisplayFields("January 2024", displayFieldsAlbum1Start, displayFieldsAlbum1End).
					build(),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).
					withCount(7).
					withDisplayFields("January 2024", displayFieldsAlbum1Start, displayFieldsAlbum1End).
					build(),
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
			err = a.SetCountForAllViewers(tt.args.ctx, tt.args.updates)
			if tt.wantErr(t, err, fmt.Sprintf("SetCountForAllViewers(%v, %v)", tt.args.ctx, tt.args.updates)) {
				dyn.MustBool(dyn.EqualContent(tt.args.ctx, tt.wantAfter))
			}
		})
	}
}

func TestAlbumViewRepository_SetDisplayFieldsForAllViewers(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)

	userId1 := usermodel.NewUserId("user-1")
	userId2 := usermodel.NewUserId("user-2")
	albumId1 := catalog.AlbumId{Owner: "owner1", FolderName: catalog.NewFolderName("/album-1")}

	type args struct {
		ctx     context.Context
		albumId catalog.AlbumId
		name    string
		start   time.Time
		end     time.Time
	}
	tests := []struct {
		name      string
		args      args
		before    []map[string]types.AttributeValue
		wantAfter []map[string]types.AttributeValue
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should do nothing when no viewer row exists for the album",
			args: args{
				ctx:     ctx,
				albumId: albumId1,
				name:    "January 2024",
				start:   displayFieldsAlbum1Start,
				end:     displayFieldsAlbum1End,
			},
			before:    nil,
			wantAfter: nil,
			wantErr:   assert.NoError,
		},
		{
			name: "it should set the display fields on every viewer row of the album",
			args: args{
				ctx:     ctx,
				albumId: albumId1,
				name:    "January 2024",
				start:   displayFieldsAlbum1Start,
				end:     displayFieldsAlbum1End,
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).withCount(11).build(),
				albumSummaryItemBuilder(userId2, VisitorAvailability, albumId1).withCount(11).build(),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).withCount(11).
					withDisplayFields("January 2024", displayFieldsAlbum1Start, displayFieldsAlbum1End).build(),
				albumSummaryItemBuilder(userId2, VisitorAvailability, albumId1).withCount(11).
					withDisplayFields("January 2024", displayFieldsAlbum1Start, displayFieldsAlbum1End).build(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should NOT clobber the Count on an existing row",
			args: args{
				ctx:     ctx,
				albumId: albumId1,
				name:    "January 2024",
				start:   displayFieldsAlbum1Start,
				end:     displayFieldsAlbum1End,
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).withCount(42).build(),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).withCount(42).
					withDisplayFields("January 2024", displayFieldsAlbum1Start, displayFieldsAlbum1End).build(),
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
			err = a.SetDisplayFieldsForAllViewers(tt.args.ctx, tt.args.albumId, tt.args.name, tt.args.start, tt.args.end)
			if tt.wantErr(t, err, fmt.Sprintf("SetDisplayFieldsForAllViewers(%v, %v)", tt.args.ctx, tt.args.albumId)) {
				dyn.MustBool(dyn.EqualContent(tt.args.ctx, tt.wantAfter))
			}
		})
	}
}

func TestAlbumViewRepository_DeleteAllRowsForAlbum(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)

	userId1 := usermodel.NewUserId("user-1")
	userId2 := usermodel.NewUserId("user-2")
	userId3 := usermodel.NewUserId("user-3")
	albumId1 := catalog.AlbumId{Owner: "owner1", FolderName: catalog.NewFolderName("/album-1")}
	albumId2 := catalog.AlbumId{Owner: "owner1", FolderName: catalog.NewFolderName("/album-2")}

	type args struct {
		ctx     context.Context
		albumId catalog.AlbumId
	}
	tests := []struct {
		name      string
		args      args
		before    []map[string]types.AttributeValue
		wantAfter []map[string]types.AttributeValue
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should delete every viewer row for the album",
			args: args{
				ctx:     ctx,
				albumId: albumId1,
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).withCount(42).build(),
				albumSummaryItemBuilder(userId2, VisitorAvailability, albumId1).withCount(42).build(),
				albumSummaryItemBuilder(userId3, VisitorAvailability, albumId1).withCount(42).build(),
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId2).withCount(7).build(),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId2).withCount(7).build(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should do nothing if no rows exist for the album",
			args: args{
				ctx:     ctx,
				albumId: albumId1,
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId2).withCount(7).build(),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId2).withCount(7).build(),
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
			err = a.DeleteAllRowsForAlbum(tt.args.ctx, tt.args.albumId)
			if tt.wantErr(t, err, fmt.Sprintf("DeleteAllRowsForAlbum(%v, %v)", tt.args.ctx, tt.args.albumId)) {
				dyn.MustBool(dyn.EqualContent(tt.args.ctx, tt.wantAfter))
			}
		})
	}
}

func TestAlbumViewRepository_RenameAlbum(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)

	owner1 := ownermodel.Owner("owner1")
	oldId := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("old-folder")}
	newId := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("new-folder")}
	otherId := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("other")}
	userId1 := usermodel.NewUserId("user-1")
	userId2 := usermodel.NewUserId("user-2")

	type args struct {
		existingId catalog.AlbumId
		renamedId  catalog.AlbumId
		newName    string
	}
	tests := []struct {
		name      string
		args      args
		before    []map[string]types.AttributeValue
		wantAfter []map[string]types.AttributeValue
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should do nothing when no row exists for the existing album",
			args: args{
				existingId: oldId,
				renamedId:  newId,
				newName:    "New Name",
			},
			before:    nil,
			wantAfter: nil,
			wantErr:   assert.NoError,
		},
		{
			name: "it should delete the old rows and recreate them under the renamed id, inheriting count and dates from the owner projection",
			args: args{
				existingId: oldId,
				renamedId:  newId,
				newName:    "New Name",
			},
			before: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, oldId).
					withCount(4).
					withDisplayFields("Old Name", displayFieldsAlbum1Start, displayFieldsAlbum1End).build(),
				albumSummaryItemBuilder(userId2, VisitorAvailability, oldId).
					withCount(4).
					withDisplayFields("Old Name", displayFieldsAlbum1Start, displayFieldsAlbum1End).build(),
				albumSummaryItemBuilder(userId1, OwnerAvailability, otherId).withCount(99).build(),
			},
			wantAfter: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, newId).
					withCount(4).
					withDisplayFields("New Name", displayFieldsAlbum1Start, displayFieldsAlbum1End).build(),
				albumSummaryItemBuilder(userId1, OwnerAvailability, otherId).withCount(99).build(),
				albumSummaryItemBuilder(userId2, VisitorAvailability, newId).
					withCount(4).
					withDisplayFields("New Name", displayFieldsAlbum1Start, displayFieldsAlbum1End).build(),
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
			err = a.RenameAlbum(ctx, tt.args.existingId, tt.args.renamedId, tt.args.newName)
			if tt.wantErr(t, err, fmt.Sprintf("RenameAlbum(%v, %v, %v)", tt.args.existingId, tt.args.renamedId, tt.args.newName)) {
				dyn.MustBool(dyn.EqualContent(ctx, tt.wantAfter))
			}
		})
	}
}

func TestAlbumViewRepository_ListSummariesForUserAndOwners(t *testing.T) {
	dyn := dynamotestutils.NewTestContext(context.Background(), t)

	owner1 := ownermodel.Owner("owner1")
	albumId1 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("album-1")}
	albumId2 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("album-2")}
	owner2 := ownermodel.Owner("owner2")
	albumId3 := catalog.AlbumId{Owner: owner2, FolderName: catalog.NewFolderName("album-3")}
	userId1 := usermodel.NewUserId("user-1")
	userId2 := usermodel.NewUserId("user-2")

	type args struct {
		ctx    context.Context
		userId usermodel.UserId
		owner  []ownermodel.Owner
	}
	tests := []struct {
		name    string
		content []map[string]types.AttributeValue
		args    args
		want    []catalogviews.UserAlbumSummary
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "it should return an empty list if there is no row",
			content: nil,
			args: args{
				ctx:    context.Background(),
				userId: userId1,
				owner:  []ownermodel.Owner{owner1},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should return no album if there is no owner selected",
			content: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).withCount(42).build(),
			},
			args: args{
				ctx:    context.Background(),
				userId: userId1,
				owner:  nil,
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should return the summaries for the owner",
			content: []map[string]types.AttributeValue{
				albumSummaryItemBuilder(userId1, OwnerAvailability, albumId1).withCount(42).build(),
				albumSummaryItemBuilder(userId1, VisitorAvailability, albumId2).withCount(42).build(),
				albumSummaryItemBuilder(userId1, VisitorAvailability, albumId3).withCount(42).build(),
				albumSummaryItemBuilder(userId2, VisitorAvailability, albumId1).withCount(42).build(),
			},
			args: args{
				ctx:    context.Background(),
				userId: userId1,
				owner:  []ownermodel.Owner{owner1},
			},
			want: []catalogviews.UserAlbumSummary{
				{Availability: catalogviews.OwnerAvailability(userId1), AlbumSummary: catalogviews.AlbumSummary{AlbumId: albumId1, MediaCount: 42}},
				{Availability: catalogviews.VisitorAvailability(userId1), AlbumSummary: catalogviews.AlbumSummary{AlbumId: albumId2, MediaCount: 42}},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn = dyn.Subtest(t)

			err := dyn.WithDbContent(tt.args.ctx, tt.content)
			if !assert.NoError(t, err) {
				return
			}

			a := &AlbumViewRepository{
				Client:    dyn.Client,
				TableName: dyn.Table,
			}
			got, err := a.ListSummariesForUserAndOwners(tt.args.ctx, tt.args.userId, tt.args.owner...)
			if !tt.wantErr(t, err, fmt.Sprintf("ListSummariesForUserAndOwners(%v, %v, %v)", tt.args.ctx, tt.args.userId, tt.args.owner)) {
				return
			}
			assert.ElementsMatchf(t, tt.want, got, "ListSummariesForUserAndOwners(%v, %v, %v)", tt.args.ctx, tt.args.userId, tt.args.owner)
		})
	}
}
