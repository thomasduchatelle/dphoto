package catalogviews

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
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
		name          string
		fields        fields
		current       []UserAlbumSummary
		args          args
		wantSummaries []UserAlbumSummary
		wantErr       assert.ErrorAssertionFunc
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
			current: []UserAlbumSummary{
				{AlbumSummary: AlbumSummary{AlbumId: album1, MediaCount: 9}, Availability: OwnerAvailability(userId1)},
				{AlbumSummary: AlbumSummary{AlbumId: album2, MediaCount: 2}, Availability: VisitorAvailability(userId1)},
				{AlbumSummary: AlbumSummary{AlbumId: album2, MediaCount: 2}, Availability: VisitorAvailability(userId2)},
				{AlbumSummary: AlbumSummary{AlbumId: album2, MediaCount: 2}, Availability: VisitorAvailability(userId3)},
			},
			wantSummaries: []UserAlbumSummary{
				{AlbumSummary: AlbumSummary{AlbumId: album1, MediaCount: 1}, Availability: OwnerAvailability(userId1)},   // drift = wrong count
				{AlbumSummary: AlbumSummary{AlbumId: album1, MediaCount: 1}, Availability: VisitorAvailability(userId2)}, // drift = missing
				{AlbumSummary: AlbumSummary{AlbumId: album2, MediaCount: 2}, Availability: OwnerAvailability(userId1)},   // drift = wrong availability type
				{AlbumSummary: AlbumSummary{AlbumId: album2, MediaCount: 2}, Availability: VisitorAvailability(userId3)}, // no drift
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
			current:       nil,
			wantSummaries: nil,
			args: args{
				owner: owner1,
				dry:   true,
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &AlbumSummaryInMemoryRepository{Summaries: tt.current}

			reconciler := NewDriftReconciler(
				tt.fields.findAlbumByOwnerPort,
				repository,
				tt.fields.listUserWhoCanAccessAlbumPort,
				tt.fields.mediaCounterPort,
				DriftOptionDryMode(tt.args.dry, repository),
			)

			err := reconciler.Reconcile(context.Background(), tt.args.owner)
			if tt.wantErr(t, err, fmt.Sprintf("Reconcile(%v, %v)", tt.args.owner, tt.args.dry)) {
				assert.ElementsMatch(t, tt.wantSummaries, repository.Summaries, "Reconcile(%v, %v) ; A=Expected ; B=Got", tt.args.owner, tt.args.dry)
			}
		})
	}
}

func TestDriftDetector_PutSummaries(t *testing.T) {
	userId1 := usermodel.UserId("user1")
	userId2 := usermodel.UserId("user2")
	albumId1 := catalog.AlbumId{Owner: "owner1", FolderName: "/folder-1"}
	albumId2 := catalog.AlbumId{Owner: "owner2", FolderName: "/folder-2"}
	user1Album1Owner := UserAlbumSummary{AlbumSummary: AlbumSummary{AlbumId: albumId1, MediaCount: 1}, Availability: OwnerAvailability(userId1)}
	user2Album1Visitor := UserAlbumSummary{AlbumSummary: AlbumSummary{AlbumId: albumId1, MediaCount: 1}, Availability: VisitorAvailability(userId2)}
	user2Album2Owner := UserAlbumSummary{AlbumSummary: AlbumSummary{AlbumId: albumId2, MediaCount: 2}, Availability: OwnerAvailability(userId2)}

	type fields struct {
		GetCurrentAlbumSummariesPort GetCurrentAlbumSummariesPort
	}
	type args struct {
		ctx       context.Context
		albumSize []AlbumSummaryForUsers
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
				GetCurrentAlbumSummariesPort: nil,
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
				GetCurrentAlbumSummariesPort: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{user1Album1Owner},
				},
			},
			args: args{
				ctx: context.Background(),
				albumSize: []AlbumSummaryForUsers{
					user1Album1Owner.ToSummaryForUsers(),
				},
			},
			wantDrifts: nil,
			wantErr:    assert.NoError,
		},
		{
			name: "it should detect missing album size as owner",
			fields: fields{
				GetCurrentAlbumSummariesPort: &AlbumSummaryInMemoryRepository{},
			},
			args: args{
				ctx: context.Background(),
				albumSize: []AlbumSummaryForUsers{
					user1Album1Owner.ToSummaryForUsers(),
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
				GetCurrentAlbumSummariesPort: &AlbumSummaryInMemoryRepository{},
			},
			args: args{
				ctx: context.Background(),
				albumSize: []AlbumSummaryForUsers{
					user2Album1Visitor.ToSummaryForUsers(),
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
				GetCurrentAlbumSummariesPort: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						user1Album1Owner,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				albumSize: []AlbumSummaryForUsers{
					{
						AlbumSummary: AlbumSummary{AlbumId: albumId1, MediaCount: 1},
						Users:        []Availability{OwnerAvailability(userId1), VisitorAvailability(userId2)},
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
				GetCurrentAlbumSummariesPort: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						{AlbumSummary: AlbumSummary{AlbumId: user1Album1Owner.AlbumSummary.AlbumId, MediaCount: 9}, Availability: user1Album1Owner.Availability},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				albumSize: []AlbumSummaryForUsers{
					user1Album1Owner.ToSummaryForUsers(),
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
				GetCurrentAlbumSummariesPort: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						user1Album1Owner,
						user2Album1Visitor,
						user2Album2Owner,
					},
				},
			},
			args: args{
				ctx: context.Background(),
				albumSize: []AlbumSummaryForUsers{
					user1Album1Owner.ToSummaryForUsers(),
					user2Album2Owner.ToSummaryForUsers(),
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
				GetCurrentAlbumSummariesPort: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						{AlbumSummary: user1Album1Owner.AlbumSummary, Availability: VisitorAvailability(userId1)},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				albumSize: []AlbumSummaryForUsers{
					user1Album1Owner.ToSummaryForUsers(),
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
				GetCurrentAlbumSummariesPort: tt.fields.GetCurrentAlbumSummariesPort,
				DriftObservers:               []DriftObserver{new(LoggerDriftObserver), observer},
			}

			err := d.PutSummaries(tt.args.ctx, tt.args.albumSize)
			if tt.wantErr(t, err, fmt.Sprintf("PutSummaries(%v, %v)", tt.args.ctx, tt.args.albumSize)) {
				assert.Equal(t, tt.wantDrifts, observer.Drifts, fmt.Sprintf("PutSummaries(%v, %v)", tt.args.ctx, tt.args.albumSize))
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
