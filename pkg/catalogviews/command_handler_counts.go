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

// AlbumMediaCountDiff carries a delta to apply on the Count attribute of a set of viewer rows
// for a single album. Used by the media insert/delete flow.
type AlbumMediaCountDiff struct {
	AlbumId        catalog.AlbumId
	Users          []Availability
	MediaCountDiff int // MediaCountDiff is the difference between the number of media added, or removed, to the album
}

// AlbumMediaCountForUsers carries an absolute Count value to SET on a set of viewer rows for a
// single album. Used by the re-count-after-transfer path.
type AlbumMediaCountForUsers struct {
	AlbumId    catalog.AlbumId
	Users      []Availability
	MediaCount int
}

// AlbumSummaryRepository is the write+read contract of the album-list view. Its five write
// primitives have disjoint attribute footprints so display-field updates and count updates can
// coexist on the same row without clobbering each other.
type AlbumSummaryRepository interface {
	// ListSummariesForUser returns every album summary row visible to a user.
	ListSummariesForUser(ctx context.Context, userId usermodel.UserId) ([]UserAlbumSummary, error)
	// PutSummaries writes full-row upserts (Put) — writes all attributes including Count.
	PutSummaries(ctx context.Context, summaries []AlbumSummaryForUsers) error
	// SetDisplayFields updates AlbumName/AlbumStart/AlbumEnd on the rows of the given users for the
	// given album. It does NOT touch Count.
	SetDisplayFields(ctx context.Context, albumId catalog.AlbumId, users []Availability, name string, start, end time.Time) error
	// IncrementCounts applies `ADD Count :d` on each viewer row. It does NOT touch display fields.
	IncrementCounts(ctx context.Context, updates []AlbumMediaCountDiff) error
	// SetCounts applies `SET Count = :c` on each viewer row. It does NOT touch display fields.
	SetCounts(ctx context.Context, updates []AlbumMediaCountForUsers) error
	// DeleteRow deletes one viewer row for the given album.
	DeleteRow(ctx context.Context, availability Availability, albumId catalog.AlbumId) error
	// DeleteAllRowsForAlbum deletes every viewer row for the given album.
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

type CommandHandlerAlbumSize struct {
	MediaCounterPort              MediaCounterPort
	ListUserWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort
	ViewWriteRepository           AlbumSummaryRepository
}

func (c *CommandHandlerAlbumSize) OnTransferredMedias(ctx context.Context, transfers catalog.TransferredMedias) error {
	var albumIds []catalog.AlbumId
	for albumId := range transfers.Transfers {
		albumIds = append(albumIds, albumId)
	}
	for _, albumId := range transfers.FromAlbums {
		albumIds = append(albumIds, albumId)
	}

	reCounter := &AlbumReCounter{
		ListUserWhoCanAccessAlbumPort: c.ListUserWhoCanAccessAlbumPort,
		MediaCounterPort:              c.MediaCounterPort,
	}
	return reCounter.ReCountMedias(ctx, albumIds, new(LoggingPutSummariesObserver), c.ViewWriteRepository)
}

func (c *CommandHandlerAlbumSize) OnMediasInserted(ctx context.Context, medias map[catalog.AlbumId][]catalog.MediaId) error {
	if len(medias) == 0 {
		return nil
	}

	var albumIds []catalog.AlbumId
	for albumId := range medias {
		albumIds = append(albumIds, albumId)
	}

	availabilities, err := c.ListUserWhoCanAccessAlbumPort.ListUsersWhoCanAccessAlbum(ctx, albumIds...)
	if err != nil {
		return err
	}

	var updates []AlbumMediaCountDiff
	for albumId, mediaIds := range medias {
		availability, _ := availabilities[albumId]
		updates = append(updates, AlbumMediaCountDiff{
			AlbumId:        albumId,
			Users:          availability,
			MediaCountDiff: len(mediaIds),
		})
	}

	return c.ViewWriteRepository.IncrementCounts(ctx, updates)
}

func (c *CommandHandlerAlbumSize) AlbumShared(ctx context.Context, albumId catalog.AlbumId, userId usermodel.UserId) error {
	counts, err := c.MediaCounterPort.CountMedia(ctx, albumId)
	if err != nil {
		return err
	}

	count, _ := counts[albumId]

	return c.ViewWriteRepository.PutSummaries(ctx, []AlbumSummaryForUsers{
		{
			AlbumSummary: AlbumSummary{
				AlbumId:    albumId,
				MediaCount: count,
			},
			Users: []Availability{VisitorAvailability(userId)},
		},
	})
}

func (c *CommandHandlerAlbumSize) AlbumUnShared(ctx context.Context, albumId catalog.AlbumId, userId usermodel.UserId) error {
	return c.ViewWriteRepository.DeleteRow(ctx, VisitorAvailability(userId), albumId)
}

type AlbumReCounter struct {
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

	var summaries []AlbumSummaryForUsers
	for _, albumId := range albumIds {
		availableTo, _ := availabilities[albumId]
		count, _ := counts[albumId]

		summaries = append(summaries, AlbumSummaryForUsers{
			AlbumSummary: AlbumSummary{
				AlbumId:    albumId,
				MediaCount: count,
			},
			Users: availableTo,
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

