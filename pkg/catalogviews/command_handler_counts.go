package catalogviews

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
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

type AlbumSize struct {
	AlbumId    catalog.AlbumId
	Users      []Availability
	MediaCount int
}

type ViewWriteRepository interface {
	InsertAlbumSize(ctx context.Context, albumSize []AlbumSize) error
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
	for albumId := range transfers {
		albumIds = append(albumIds, albumId)
	}

	return c.updateUserViews(ctx, albumIds)
}

func (c *CommandHandlerAlbumSize) OnMediasInserted(ctx context.Context, medias map[catalog.AlbumId][]catalog.MediaId) error {
	var albumIds []catalog.AlbumId
	for albumId := range medias {
		albumIds = append(albumIds, albumId)
	}

	return c.updateUserViews(ctx, albumIds)
}

func (c *CommandHandlerAlbumSize) updateUserViews(ctx context.Context, albumIds []catalog.AlbumId) error {
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

	var albumSizes []AlbumSize
	for _, albumId := range albumIds {
		availableTo, _ := availabilities[albumId]
		count, _ := counts[albumId]

		albumSizes = append(albumSizes, AlbumSize{
			AlbumId:    albumId,
			Users:      availableTo,
			MediaCount: count,
		})
	}

	return c.ViewWriteRepository.InsertAlbumSize(ctx, albumSizes)
}

// TODO Everything about the album should be deleted if the album is deleted
