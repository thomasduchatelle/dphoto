package catalogviews

import (
	"context"
	"fmt"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"strings"
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

type AlbumSize struct {
	AlbumId    catalog.AlbumId
	MediaCount int
}

type MultiUserAlbumSize struct {
	AlbumSize
	Users []Availability
}

func (a MultiUserAlbumSize) String() string {
	var users []string
	for _, user := range a.Users {
		users = append(users, user.String())
	}
	return fmt.Sprintf("%s: %d media(s) available to %s", a.AlbumId, a.MediaCount, strings.Join(users, ", "))
}

type AlbumSizeDiff struct {
	AlbumId        catalog.AlbumId
	Users          []Availability
	MediaCountDiff int // MediaCountDiff is the difference between the number of media added, or removed, to the album
}

type InsertAlbumSizePort interface {
	InsertAlbumSize(ctx context.Context, albumSize []MultiUserAlbumSize) error
}

type DeleteAlbumSizePort interface {
	DeleteAlbumSize(ctx context.Context, availability Availability, albumId catalog.AlbumId) error
}

type ViewWriteRepository interface {
	InsertAlbumSizePort
	DeleteAlbumSizePort
	UpdateAlbumSize(ctx context.Context, albumCountUpdates []AlbumSizeDiff) error
}

type ListUserWhoCanAccessAlbumPort interface {
	ListUsersWhoCanAccessAlbum(ctx context.Context, albumId ...catalog.AlbumId) (map[catalog.AlbumId][]Availability, error)
}

type CommandHandlerAlbumSize struct {
	MediaCounterPort              MediaCounterPort
	ListUserWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort
	ViewWriteRepository           ViewWriteRepository
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
	return reCounter.ReCountMedias(ctx, albumIds, new(LoggingInsertAlbumSizeObserver), c.ViewWriteRepository)
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

	var updates []AlbumSizeDiff
	for albumId, mediaIds := range medias {
		availability, _ := availabilities[albumId]
		updates = append(updates, AlbumSizeDiff{
			AlbumId:        albumId,
			Users:          availability,
			MediaCountDiff: len(mediaIds),
		})
	}

	return c.ViewWriteRepository.UpdateAlbumSize(ctx, updates)
}

func (c *CommandHandlerAlbumSize) AlbumShared(ctx context.Context, albumId catalog.AlbumId, userId usermodel.UserId) error {
	counts, err := c.MediaCounterPort.CountMedia(ctx, albumId)
	if err != nil {
		return err
	}

	count, _ := counts[albumId]

	return c.ViewWriteRepository.InsertAlbumSize(ctx, []MultiUserAlbumSize{
		{
			AlbumSize: AlbumSize{
				AlbumId:    albumId,
				MediaCount: count,
			},
			Users: []Availability{VisitorAvailability(userId)},
		},
	})
}

func (c *CommandHandlerAlbumSize) AlbumUnShared(ctx context.Context, albumId catalog.AlbumId, userId usermodel.UserId) error {
	return c.ViewWriteRepository.DeleteAlbumSize(ctx, VisitorAvailability(userId), albumId)
}

type AlbumReCounter struct {
	ListUserWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort
	MediaCounterPort              MediaCounterPort
}

func (c *AlbumReCounter) ReCountMedias(ctx context.Context, albumIds []catalog.AlbumId, observers ...InsertAlbumSizePort) error {
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

	var albumSizes []MultiUserAlbumSize
	for _, albumId := range albumIds {
		availableTo, _ := availabilities[albumId]
		count, _ := counts[albumId]

		albumSizes = append(albumSizes, MultiUserAlbumSize{
			AlbumSize: AlbumSize{
				AlbumId:    albumId,
				MediaCount: count,
			},
			Users: availableTo,
		})
	}

	for _, observer := range observers {
		err = observer.InsertAlbumSize(ctx, albumSizes)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO Everything about the album should be deleted if the album is deleted
