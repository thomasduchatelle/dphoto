package catalogviews

import (
	"context"
	"fmt"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"strings"
	"time"
)

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

// AlbumSummary is the per-album projection carried by the album-list view: identity, display
// fields and the media count. It replaces the old AlbumSize.
type AlbumSummary struct {
	AlbumId    catalog.AlbumId
	MediaCount int
	Name       string
	Start      time.Time
	End        time.Time
}

// AlbumSummaryForUsers is a full projection to be written for a set of users (owner + visitors).
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
	IncrementCountForAllViewers(ctx context.Context, updates []AlbumCountDiff) error
	SetCountForAllViewers(ctx context.Context, updates []AlbumCount) error
	RenameAlbum(ctx context.Context, existingId, renamedId catalog.AlbumId, newName string) error
	DeleteRow(ctx context.Context, availability Availability, albumId catalog.AlbumId) error
	DeleteAllRowsForAlbum(ctx context.Context, albumId catalog.AlbumId) error
}

type PutSummariesPort interface {
	PutSummaries(ctx context.Context, summaries []AlbumSummaryForUsers) error
}

type DeleteRowPort interface {
	DeleteRow(ctx context.Context, availability Availability, albumId catalog.AlbumId) error
}

type ListUserWhoCanAccessAlbumPort interface {
	ListUsersWhoCanAccessAlbum(ctx context.Context, albumId ...catalog.AlbumId) (map[catalog.AlbumId][]Availability, error)
}

type AlbumReCounter struct {
	FindAlbumsByIdsPort           FindAlbumsByIdsPort
	ListUserWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort
	MediaCounterPort              MediaCounterPort
}

func (c *AlbumReCounter) ReCountMedias(ctx context.Context, albumIds []catalog.AlbumId, observers ...PutSummariesPort) error {
	if len(albumIds) == 0 {
		return nil
	}

	availabilities, err := c.ListUserWhoCanAccessAlbumPort.ListUsersWhoCanAccessAlbum(ctx, albumIds...)
	if err != nil {
		return err
	}

	counts, err := c.MediaCounterPort.CountMedia(ctx, albumIds...)
	if err != nil {
		return err
	}

	albums, err := c.FindAlbumsByIdsPort.FindAlbumsById(ctx, albumIds)
	if err != nil {
		return err
	}

	albumsById := make(map[catalog.AlbumId]*catalog.Album, len(albums))
	for _, album := range albums {
		albumsById[album.AlbumId] = album
	}

	var summaries []AlbumSummaryForUsers
	for _, albumId := range albumIds {
		album, present := albumsById[albumId]
		if !present {
			continue
		}

		summaries = append(summaries, AlbumSummaryForUsers{
			AlbumSummary: AlbumSummary{
				AlbumId:    albumId,
				MediaCount: counts[albumId],
				Name:       album.Name,
				Start:      album.Start,
				End:        album.End,
			},
			Users: availabilities[albumId],
		})
	}

	for _, observer := range observers {
		err = observer.PutSummaries(ctx, summaries)
		if err != nil {
			return err
		}
	}

	return nil
}
