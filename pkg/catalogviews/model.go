package catalogviews

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type VisibleAlbum struct {
	catalog.Album
	MediaCount         int                // MediaCount is the number of medias on the album
	Visitors           []usermodel.UserId // Visitors are the users that can see the album ; only visible to the owner of the album
	OwnedByCurrentUser bool               // OwnedByCurrentUser is set to true when the user is an owner of the album
	Covers             []catalog.Cover    // Covers are the 0 to 4 featured medias of the album, denormalised from the canonical cover set
}

type ListAlbumsFilter struct {
	OnlyDirectlyOwned bool // OnlyDirectlyOwned provides a sub-view where only resources directly owned by user are displayed and accessible
}

type Availability struct {
	UserId  usermodel.UserId
	AsOwner bool // AsOwner is true if the user is the owner of the album
}

func VisitorAvailability(userId usermodel.UserId) Availability {
	return Availability{
		UserId: userId,
	}
}

func OwnerAvailability(userId usermodel.UserId) Availability {
	return Availability{
		UserId:  userId,
		AsOwner: true,
	}
}

func (a Availability) String() string {
	availabilityType := "visitor"
	if a.AsOwner {
		availabilityType = "owner"
	}
	return fmt.Sprintf("%s:%s", availabilityType, a.UserId.Value())
}

// AlbumSummary is the per-album projection carried by the main album-list row: identity,
// display fields and the media count. The cover set lives on a sibling row (see
// AlbumCoversForUsers); the adapter merges both rows per album in Covers when a view is
// read.
type AlbumSummary struct {
	AlbumId    catalog.AlbumId
	MediaCount int
	Name       string
	Start      time.Time
	End        time.Time
	Covers     []catalog.Cover
}

// AlbumSummaryForUsers is a full projection to be written for a set of users (owner + visitors).
// When Covers is non-empty, PutSummaries writes the sibling cover row for the listed users in
// the same call so a newly-created viewer row carries its covers without a second fan-out.
// Cross-viewer cover updates (cover changes, backfills) still travel through
// PutCoversForAllViewers so they stay independent from count / display-field writes.
type AlbumSummaryForUsers struct {
	AlbumSummary
	Users []Availability
}

func (a AlbumSummaryForUsers) String() string {
	var users []string
	for _, user := range a.Users {
		users = append(users, user.String())
	}
	return fmt.Sprintf("%s: %d media(s) available to %s", a.AlbumId, a.MediaCount, strings.Join(users, ", "))
}

type AlbumCountDiff struct {
	AlbumId        catalog.AlbumId
	MediaCountDiff int
}

type AlbumCount struct {
	AlbumId    catalog.AlbumId
	MediaCount int
}

type AlbumSummaryRepository interface {
	ListSummariesForUser(ctx context.Context, userId usermodel.UserId) ([]UserAlbumSummary, error)
	PutSummaries(ctx context.Context, summaries []AlbumSummaryForUsers) error
	SetDisplayFieldsForAllViewers(ctx context.Context, albumId catalog.AlbumId, name string, start, end time.Time) error
	PutCoversForAllViewers(ctx context.Context, albumId catalog.AlbumId, covers []catalog.Cover) error
	DeleteCoversForAllViewers(ctx context.Context, albumId catalog.AlbumId) error
	IncrementCountForAllViewers(ctx context.Context, updates []AlbumCountDiff) error
	SetCountForAllViewers(ctx context.Context, updates []AlbumCount) error
	RenameAlbum(ctx context.Context, existingId, renamedId catalog.AlbumId, newName string) error
	DeleteRow(ctx context.Context, availability Availability, albumId catalog.AlbumId) error
	DeleteAllRowsForAlbum(ctx context.Context, albumId catalog.AlbumId) error
}

type PutSummariesPort interface {
	PutSummaries(ctx context.Context, summaries []AlbumSummaryForUsers) error
}

type PutCoversForAllViewersPort interface {
	PutCoversForAllViewers(ctx context.Context, albumId catalog.AlbumId, covers []catalog.Cover) error
}

type DeleteCoversForAllViewersPort interface {
	DeleteCoversForAllViewers(ctx context.Context, albumId catalog.AlbumId) error
}

type DeleteRowPort interface {
	DeleteRow(ctx context.Context, availability Availability, albumId catalog.AlbumId) error
}

type DeleteLegacyRowsForUserPort interface {
	DeleteLegacyRowsForUser(ctx context.Context, userId usermodel.UserId) error
}

type ListUserWhoCanAccessAlbumPort interface {
	ListUsersWhoCanAccessAlbum(ctx context.Context, albumId ...catalog.AlbumId) (map[catalog.AlbumId][]Availability, error)
}
