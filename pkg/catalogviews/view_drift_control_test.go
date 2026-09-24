package catalogviews

import (
	"context"
	"fmt"
	"testing"
	"time"

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

	album1Start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	album1End := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	album2Start := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	album2End := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

	canonicalAlbum1 := &catalog.Album{AlbumId: album1, Name: "Album One", Start: album1Start, End: album1End}
	canonicalAlbum2 := &catalog.Album{AlbumId: album2, Name: "Album Two", Start: album2Start, End: album2End}
	album1Summary := AlbumSummary{AlbumId: album1, MediaCount: 1, Name: "Album One", Start: album1Start, End: album1End}
	album2Summary := AlbumSummary{AlbumId: album2, MediaCount: 2, Name: "Album Two", Start: album2Start, End: album2End}

	staleAlbum2Visitor1 := UserAlbumSummary{AlbumSummary: AlbumSummary{AlbumId: album2, MediaCount: 2, Name: "Album Two", Start: album2Start, End: album2End}, Availability: VisitorAvailability(userId1)}
	staleAlbum2Visitor2 := UserAlbumSummary{AlbumSummary: AlbumSummary{AlbumId: album2, MediaCount: 2, Name: "Album Two", Start: album2Start, End: album2End}, Availability: VisitorAvailability(userId2)}

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
		name                        string
		fields                      fields
		current                     []UserAlbumSummary
		args                        args
		wantSummaries               []UserAlbumSummary
		wantDrifts                  []Drift
		expectLegacyCleanupForUsers []usermodel.UserId
		wantErr                     assert.ErrorAssertionFunc
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
			wantDrifts: nil,
			wantErr:    assert.NoError,
		},
		{
			name: "it should reconcile the 3 different types of drifts",
			fields: fields{
				findAlbumByOwnerPort: stubFindAlbumByOwnerPort(canonicalAlbum1, canonicalAlbum2),
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
				{AlbumSummary: AlbumSummary{AlbumId: album1, MediaCount: 9, Name: "Album One", Start: album1Start, End: album1End}, Availability: OwnerAvailability(userId1)},
				staleAlbum2Visitor1,
				staleAlbum2Visitor2,
				{AlbumSummary: AlbumSummary{AlbumId: album2, MediaCount: 2, Name: "Album Two", Start: album2Start, End: album2End}, Availability: VisitorAvailability(userId3)},
			},
			wantSummaries: []UserAlbumSummary{
				{AlbumSummary: album1Summary, Availability: OwnerAvailability(userId1)},   // fixed count
				{AlbumSummary: album1Summary, Availability: VisitorAvailability(userId2)}, // missing added
				{AlbumSummary: album2Summary, Availability: OwnerAvailability(userId1)},   // swapped from visitor
				{AlbumSummary: album2Summary, Availability: VisitorAvailability(userId3)}, // untouched
			},
			wantDrifts: []Drift{
				NewOverrideDrift(UserAlbumSummary{AlbumSummary: album1Summary, Availability: OwnerAvailability(userId1)}),
				NewDeletedDrift(staleAlbum2Visitor1),
				NewMissingDrift(UserAlbumSummary{AlbumSummary: album2Summary, Availability: OwnerAvailability(userId1)}),
				NewDeletedDrift(staleAlbum2Visitor2),
				NewMissingDrift(UserAlbumSummary{AlbumSummary: album1Summary, Availability: VisitorAvailability(userId2)}),
			},
			expectLegacyCleanupForUsers: []usermodel.UserId{userId1, userId2, userId3},
			args: args{
				owner: owner1,
				dry:   false,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not do anything on dry mode",
			fields: fields{
				findAlbumByOwnerPort: stubFindAlbumByOwnerPort(canonicalAlbum1),
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
			wantDrifts: []Drift{
				NewMissingDrift(UserAlbumSummary{AlbumSummary: album1Summary, Availability: OwnerAvailability(userId1)}),
			},
			args: args{
				owner: owner1,
				dry:   true,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should rebuild display fields from the canonical album",
			fields: fields{
				findAlbumByOwnerPort: stubFindAlbumByOwnerPort(canonicalAlbum1),
				listUserWhoCanAccessAlbumPort: &ListUserWhoCanAccessAlbumPortFake{
					Values: map[catalog.AlbumId][]Availability{
						album1: {OwnerAvailability(userId1)},
					},
				},
				mediaCounterPort: &MediaCounterPortFake{
					album1: 1,
				},
			},
			current: []UserAlbumSummary{
				{AlbumSummary: AlbumSummary{AlbumId: album1, MediaCount: 1, Name: "Stale Name", Start: album1End, End: album1Start}, Availability: OwnerAvailability(userId1)},
			},
			wantSummaries: []UserAlbumSummary{
				{AlbumSummary: album1Summary, Availability: OwnerAvailability(userId1)},
			},
			wantDrifts: []Drift{
				NewOverrideDrift(UserAlbumSummary{AlbumSummary: album1Summary, Availability: OwnerAvailability(userId1)}),
			},
			expectLegacyCleanupForUsers: []usermodel.UserId{userId1},
			args: args{
				owner: owner1,
				dry:   false,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should backfill missing display fields on a legacy row",
			fields: fields{
				findAlbumByOwnerPort: stubFindAlbumByOwnerPort(canonicalAlbum1),
				listUserWhoCanAccessAlbumPort: &ListUserWhoCanAccessAlbumPortFake{
					Values: map[catalog.AlbumId][]Availability{
						album1: {OwnerAvailability(userId1)},
					},
				},
				mediaCounterPort: &MediaCounterPortFake{
					album1: 1,
				},
			},
			current: []UserAlbumSummary{
				{AlbumSummary: AlbumSummary{AlbumId: album1, MediaCount: 1}, Availability: OwnerAvailability(userId1)},
			},
			wantSummaries: []UserAlbumSummary{
				{AlbumSummary: album1Summary, Availability: OwnerAvailability(userId1)},
			},
			wantDrifts: []Drift{
				NewOverrideDrift(UserAlbumSummary{AlbumSummary: album1Summary, Availability: OwnerAvailability(userId1)}),
			},
			expectLegacyCleanupForUsers: []usermodel.UserId{userId1},
			args: args{
				owner: owner1,
				dry:   false,
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

			drifts, err := reconciler.Reconcile(context.Background(), tt.args.owner)
			if tt.wantErr(t, err, fmt.Sprintf("Reconcile(%v, %v)", tt.args.owner, tt.args.dry)) {
				assert.ElementsMatch(t, tt.wantSummaries, repository.Summaries, "Reconcile(%v, %v) ; A=Expected ; B=Got", tt.args.owner, tt.args.dry)
				assert.ElementsMatch(t, tt.wantDrifts, drifts, "Reconcile(%v, %v) drifts ; A=Expected ; B=Got", tt.args.owner, tt.args.dry)
				assert.ElementsMatch(t, tt.expectLegacyCleanupForUsers, repository.LegacyRowsCleanedForUsers, "Reconcile(%v, %v) legacy cleanup ; A=Expected ; B=Got", tt.args.owner, tt.args.dry)
			}
		})
	}
}

func TestDriftDetector_Detect(t *testing.T) {
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
	tests := []struct {
		name       string
		fields     fields
		expected   []AlbumSummaryForUsers
		wantDrifts []Drift
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "it should not detect any drift if everything is empty",
			fields: fields{
				GetCurrentAlbumSummariesPort: nil,
			},
			expected:   nil,
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
			expected: []AlbumSummaryForUsers{
				user1Album1Owner.ToSummaryForUsers(),
			},
			wantDrifts: nil,
			wantErr:    assert.NoError,
		},
		{
			name: "it should detect a missing row for the album owner",
			fields: fields{
				GetCurrentAlbumSummariesPort: &AlbumSummaryInMemoryRepository{},
			},
			expected: []AlbumSummaryForUsers{
				user1Album1Owner.ToSummaryForUsers(),
			},
			wantDrifts: []Drift{
				NewMissingDrift(user1Album1Owner),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should detect a missing row for a visitor",
			fields: fields{
				GetCurrentAlbumSummariesPort: &AlbumSummaryInMemoryRepository{},
			},
			expected: []AlbumSummaryForUsers{
				user2Album1Visitor.ToSummaryForUsers(),
			},
			wantDrifts: []Drift{
				NewMissingDrift(user2Album1Visitor),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should detect a missing visitor row when the owner is already up to date",
			fields: fields{
				GetCurrentAlbumSummariesPort: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{user1Album1Owner},
				},
			},
			expected: []AlbumSummaryForUsers{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumId1, MediaCount: 1},
					Users:        []Availability{OwnerAvailability(userId1), VisitorAvailability(userId2)},
				},
			},
			wantDrifts: []Drift{
				NewMissingDrift(user2Album1Visitor),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should detect a stale count on an existing row",
			fields: fields{
				GetCurrentAlbumSummariesPort: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						{AlbumSummary: AlbumSummary{AlbumId: albumId1, MediaCount: 9}, Availability: OwnerAvailability(userId1)},
					},
				},
			},
			expected: []AlbumSummaryForUsers{
				user1Album1Owner.ToSummaryForUsers(),
			},
			wantDrifts: []Drift{
				NewOverrideDrift(user1Album1Owner),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should detect stale display fields on an existing row",
			fields: fields{
				GetCurrentAlbumSummariesPort: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						{
							AlbumSummary: AlbumSummary{
								AlbumId:    albumId1,
								MediaCount: 1,
								Name:       "Stale Name",
								Start:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
								End:        time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
							},
							Availability: OwnerAvailability(userId1),
						},
					},
				},
			},
			expected: []AlbumSummaryForUsers{
				{
					AlbumSummary: AlbumSummary{
						AlbumId:    albumId1,
						MediaCount: 1,
						Name:       "Canonical Name",
						Start:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						End:        time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
					},
					Users: []Availability{OwnerAvailability(userId1)},
				},
			},
			wantDrifts: []Drift{
				NewOverrideDrift(UserAlbumSummary{
					AlbumSummary: AlbumSummary{
						AlbumId:    albumId1,
						MediaCount: 1,
						Name:       "Canonical Name",
						Start:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						End:        time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
					},
					Availability: OwnerAvailability(userId1),
				}),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should flag an orphan row for deletion",
			fields: fields{
				GetCurrentAlbumSummariesPort: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						user1Album1Owner,
						user2Album1Visitor,
						user2Album2Owner,
					},
				},
			},
			expected: []AlbumSummaryForUsers{
				user1Album1Owner.ToSummaryForUsers(),
				user2Album2Owner.ToSummaryForUsers(),
			},
			wantDrifts: []Drift{
				NewDeletedDrift(user2Album1Visitor),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should flag an availability swap as delete + missing",
			fields: fields{
				GetCurrentAlbumSummariesPort: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						{AlbumSummary: user1Album1Owner.AlbumSummary, Availability: VisitorAvailability(userId1)},
					},
				},
			},
			expected: []AlbumSummaryForUsers{
				user1Album1Owner.ToSummaryForUsers(),
			},
			wantDrifts: []Drift{
				NewDeletedDrift(UserAlbumSummary{AlbumSummary: user1Album1Owner.AlbumSummary, Availability: VisitorAvailability(userId1)}),
				NewMissingDrift(user1Album1Owner),
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DriftDetector{GetCurrentAlbumSummariesPort: tt.fields.GetCurrentAlbumSummariesPort}

			got, _, err := d.Detect(context.Background(), tt.expected)
			if tt.wantErr(t, err, fmt.Sprintf("Detect(%v)", tt.expected)) {
				assert.ElementsMatch(t, tt.wantDrifts, got, "Detect(%v)", tt.expected)
			}
		})
	}
}
