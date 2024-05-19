package catalogviews

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type AlbumSize struct {
	AlbumId    catalog.AlbumId
	Users      []usermodel.UserId
	MediaCount int
}

type ViewWriteRepository interface {
	InsertAlbumSize(ctx context.Context, albumSize []AlbumSize) error
}

type ListUserWhoCanAccessAlbumPort interface {
	ListUserWhoCanAccessAlbum(ctx context.Context, albumId ...catalog.AlbumId) (map[catalog.AlbumId][]usermodel.UserId, error)
}

type CommandHandlerAlbumSize struct {
	MediaCounterPort              MediaCounterPort
	ListUserWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort
	ViewWriteRepository           ViewWriteRepository
}

func (c *CommandHandlerAlbumSize) OnTransferredMedias(ctx context.Context, transfers catalog.TransferredMedias) error {
	if len(transfers) == 0 {
		return nil
	}

	var albumIds []catalog.AlbumId
	for albumId := range transfers {
		albumIds = append(albumIds, albumId)
	}

	availabilities, err := c.ListUserWhoCanAccessAlbumPort.ListUserWhoCanAccessAlbum(ctx, albumIds...)
	if err != nil {
		return err
	}

	counts, err := c.MediaCounterPort.CountMedia(ctx, albumIds...)
	if err != nil {
		return err
	}

	var albumSizes []AlbumSize
	for albumId := range transfers {
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

func (c *CommandHandlerAlbumSize) OnMediasInserted(ctx context.Context, medias map[catalog.AlbumId][]catalog.MediaId) error {
	//TODO implement me
	panic("implement me")
}

// TODO Everything about the album should be deleted if the album is deleted
