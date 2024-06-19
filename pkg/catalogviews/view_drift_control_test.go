package catalogviews

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"testing"
)

func TestNewDriftReconcilerAcceptance(t *testing.T) {
	owner1 := ownermodel.Owner("owner1")
	album1 := catalog.AlbumId{Owner: owner1, FolderName: "/folder-1"}
	album2 := catalog.AlbumId{Owner: owner1, FolderName: "/folder-2"}
	userId1 := usermodel.UserId("user1")
	userId2 := usermodel.UserId("user2")
	userId3 := usermodel.UserId("user3")

	type fields struct {
		findAlbumByOwnerPort          FindAlbumByOwnerPort
		listUserWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort
		mediaCounterPort              MediaCounterPort
	}
	type args struct {
		owner ownermodel.Owner
		dry   bool
	}
	tests := []struct {
		name      string
		fields    fields
		current   []UserAlbumSize
		args      args
		wantSizes []UserAlbumSize
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should not fail when no album is found for the owner",
			fields: fields{
				findAlbumByOwnerPort: stubFindAlbumByOwnerPort(),
			},
			args: args{
				owner: owner1,
				dry:   false,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should reconcile the 3 different types of drifts",
			fields: fields{
				findAlbumByOwnerPort: stubFindAlbumByOwnerPort(&catalog.Album{AlbumId: album1}, &catalog.Album{AlbumId: album2}),
				listUserWhoCanAccessAlbumPort: &ListUserWhoCanAccessAlbumPortFake{
					Values: map[catalog.AlbumId][]Availability{
						album1: {OwnerAvailability(userId1), VisitorAvailability(userId2)},
						album2: {OwnerAvailability(userId1), VisitorAvailability(userId3)},
					},
				},
				mediaCounterPort: &MediaCounterPortFake{
					album1: 1,
					album2: 2,
				},
			},
			current: []UserAlbumSize{
				{AlbumSize: AlbumSize{AlbumId: album1, MediaCount: 9}, Availability: OwnerAvailability(userId1)},
				{AlbumSize: AlbumSize{AlbumId: album2, MediaCount: 2}, Availability: VisitorAvailability(userId1)},
				{AlbumSize: AlbumSize{AlbumId: album2, MediaCount: 2}, Availability: VisitorAvailability(userId2)},
				{AlbumSize: AlbumSize{AlbumId: album2, MediaCount: 2}, Availability: VisitorAvailability(userId3)},
			},
			wantSizes: []UserAlbumSize{
				{AlbumSize: AlbumSize{AlbumId: album1, MediaCount: 1}, Availability: OwnerAvailability(userId1)},   // drift = wrong count
				{AlbumSize: AlbumSize{AlbumId: album1, MediaCount: 1}, Availability: VisitorAvailability(userId2)}, // drift = missing
				{AlbumSize: AlbumSize{AlbumId: album2, MediaCount: 2}, Availability: OwnerAvailability(userId1)},   // drift = wrong availability type
				{AlbumSize: AlbumSize{AlbumId: album2, MediaCount: 2}, Availability: VisitorAvailability(userId3)}, // no drift
			},
			args: args{
				owner: owner1,
				dry:   false,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not do anything on dry mode",
			fields: fields{
				findAlbumByOwnerPort: stubFindAlbumByOwnerPort(&catalog.Album{AlbumId: album1}),
				listUserWhoCanAccessAlbumPort: &ListUserWhoCanAccessAlbumPortFake{
					Values: map[catalog.AlbumId][]Availability{
						album1: {OwnerAvailability(userId1)},
					},
				},
				mediaCounterPort: &MediaCounterPortFake{
					album1: 1,
				},
			},
			current:   nil,
			wantSizes: nil,
			args: args{
				owner: owner1,
				dry:   true,
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &AlbumSizeInMemoryRepository{Sizes: tt.current}

			reconciler := NewDriftReconciler(
				tt.fields.findAlbumByOwnerPort,
				repository,
				tt.fields.listUserWhoCanAccessAlbumPort,
				tt.fields.mediaCounterPort,
				DriftOptionDryMode(tt.args.dry, repository),
			)

			err := reconciler.Reconcile(context.Background(), tt.args.owner)
			if tt.wantErr(t, err, fmt.Sprintf("Reconcile(%v, %v)", tt.args.owner, tt.args.dry)) {
				assert.ElementsMatch(t, tt.wantSizes, repository.Sizes, "Reconcile(%v, %v) ; A=Expected ; B=Got", tt.args.owner, tt.args.dry)
			}
		})
	}
}

func TestDriftMeasurer_InsertAlbumSize(t *testing.T) {
	userId1 := usermodel.UserId("user1")
	userId2 := usermodel.UserId("user2")
	albumId1 := catalog.AlbumId{Owner: "owner1", FolderName: "/folder-1"}
	albumId2 := catalog.AlbumId{Owner: "owner2", FolderName: "/folder-2"}
	user1Album1Owner := UserAlbumSize{AlbumSize: AlbumSize{AlbumId: albumId1, MediaCount: 1}, Availability: OwnerAvailability(userId1)}
	user2Album1Visitor := UserAlbumSize{AlbumSize: AlbumSize{AlbumId: albumId1, MediaCount: 1}, Availability: VisitorAvailability(userId2)}
	user2Album2Owner := UserAlbumSize{AlbumSize: AlbumSize{AlbumId: albumId2, MediaCount: 2}, Availability: OwnerAvailability(userId2)}

	type fields struct {
		GetCurrentAlbumSizesPort GetCurrentAlbumSizesPort
	}
	type args struct {
		ctx       context.Context
		albumSize []MultiUserAlbumSize
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantDrifts []Drift
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "it should not detect any drift if everything is empty",
			fields: fields{
				GetCurrentAlbumSizesPort: nil,
			},
			args: args{
				ctx:       context.Background(),
				albumSize: nil,
			},
			wantDrifts: nil,
			wantErr:    assert.NoError,
		},
		{
			name: "it should not detect any drift if current match the expected",
			fields: fields{
				GetCurrentAlbumSizesPort: &AlbumSizeInMemoryRepository{
					Sizes: []UserAlbumSize{user1Album1Owner},
				},
			},
			args: args{
				ctx: context.Background(),
				albumSize: []MultiUserAlbumSize{
					user1Album1Owner.ToMultiUser(),
				},
			},
			wantDrifts: nil,
			wantErr:    assert.NoError,
		},
		{
			name: "it should detect missing album size as owner",
			fields: fields{
				GetCurrentAlbumSizesPort: &AlbumSizeInMemoryRepository{},
			},
			args: args{
				ctx: context.Background(),
				albumSize: []MultiUserAlbumSize{
					user1Album1Owner.ToMultiUser(),
				},
			},
			wantDrifts: []Drift{
				NewMissingDrift(user1Album1Owner),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should detect missing album size as visitor",
			fields: fields{
				GetCurrentAlbumSizesPort: &AlbumSizeInMemoryRepository{},
			},
			args: args{
				ctx: context.Background(),
				albumSize: []MultiUserAlbumSize{
					user2Album1Visitor.ToMultiUser(),
				},
			},
			wantDrifts: []Drift{
				NewMissingDrift(user2Album1Visitor),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should detect missing album size when albums is shared to multiple users",
			fields: fields{
				GetCurrentAlbumSizesPort: &AlbumSizeInMemoryRepository{
					Sizes: []UserAlbumSize{
						user1Album1Owner,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				albumSize: []MultiUserAlbumSize{
					{
						AlbumSize: AlbumSize{AlbumId: albumId1, MediaCount: 1},
						Users:     []Availability{OwnerAvailability(userId1), VisitorAvailability(userId2)},
					},
				},
			},
			wantDrifts: []Drift{
				NewMissingDrift(user2Album1Visitor),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should detect different album size",
			fields: fields{
				GetCurrentAlbumSizesPort: &AlbumSizeInMemoryRepository{
					Sizes: []UserAlbumSize{
						{AlbumSize: AlbumSize{AlbumId: user1Album1Owner.AlbumSize.AlbumId, MediaCount: 9}, Availability: user1Album1Owner.Availability},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				albumSize: []MultiUserAlbumSize{
					user1Album1Owner.ToMultiUser(),
				},
			},
			wantDrifts: []Drift{
				NewOverrideDrift(user1Album1Owner),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should detect a size that is still present but shouldn't be.",
			fields: fields{
				GetCurrentAlbumSizesPort: &AlbumSizeInMemoryRepository{
					Sizes: []UserAlbumSize{
						user1Album1Owner,
						user2Album1Visitor,
						user2Album2Owner,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				albumSize: []MultiUserAlbumSize{
					user1Album1Owner.ToMultiUser(),
					user2Album2Owner.ToMultiUser(), // user2 is not synced if not explicitly in the expected list
				},
			},
			wantDrifts: []Drift{
				NewNotExpectedDrift(VisitorAvailability(userId2), albumId1),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should detect when a user is not at the right level of availability",
			fields: fields{
				GetCurrentAlbumSizesPort: &AlbumSizeInMemoryRepository{
					Sizes: []UserAlbumSize{
						{AlbumSize: user1Album1Owner.AlbumSize, Availability: VisitorAvailability(userId1)},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				albumSize: []MultiUserAlbumSize{
					user1Album1Owner.ToMultiUser(),
				},
			},
			wantDrifts: []Drift{
				NewNotExpectedDrift(VisitorAvailability(userId1), albumId1),
				NewMissingDrift(user1Album1Owner),
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := new(DriftObserverFake)
			d := &DriftDetector{
				GetCurrentAlbumSizesPort: tt.fields.GetCurrentAlbumSizesPort,
				DriftObservers:           []DriftObserver{new(LoggerDriftObserver), observer},
			}

			err := d.InsertAlbumSize(tt.args.ctx, tt.args.albumSize)
			if tt.wantErr(t, err, fmt.Sprintf("InsertAlbumSize(%v, %v)", tt.args.ctx, tt.args.albumSize)) {
				assert.Equal(t, tt.wantDrifts, observer.Drifts, fmt.Sprintf("InsertAlbumSize(%v, %v)", tt.args.ctx, tt.args.albumSize))
			}
		})
	}
}

type DriftObserverFake struct {
	Drifts []Drift
}

func (d *DriftObserverFake) OnDetectedDrifts(ctx context.Context, drifts []Drift) error {
	d.Drifts = append(d.Drifts, drifts...)
	return nil
}
