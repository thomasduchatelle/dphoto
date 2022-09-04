package common

import (
	"github.com/thomasduchatelle/dphoto/domain/accessadapters/oauth"
)

func Or(expressions ...func(claims oauth.Claims) error) func(claims oauth.Claims) error {
	return func(claims oauth.Claims) error {
		for _, expr := range expressions {
			err := expr(claims)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func CanReadAsOwner(owner string) func(claims oauth.Claims) error {
	return func(claims oauth.Claims) error {
		return claims.IsOwnerOf(owner)
	}
}

func CanReadAlbum(owner, album string) func(claims oauth.Claims) error {
	return func(claims oauth.Claims) error {
		//granted, err := access.IsGranted("foobar", access.Resource{
		//	Type:  access.AlbumResource,
		//	Owner: owner,
		//	Id:    album,
		//})
		return nil
	}
}

func CanReadMedia(owner, mediaId string) func(claims oauth.Claims) error {
	return func(claims oauth.Claims) error {
		//var details catalog.MediaMeta
		//var err error
		//details, err = catalog.GetMediaDetails(owner, mediaId)
		//if err != nil {
		//	return errors.Wrapf(err, "couldn't read media %s/%s metadata", owner, mediaId)
		//}
		//
		//granted, err := access.IsGrantedToAny("foobar", access.Resource{
		//	Type:  access.AlbumResource,
		//	Owner: owner,
		//	Id:    details.FolderName,
		//}, access.Resource{
		//	Type:  access.MediaResource,
		//	Owner: owner,
		//	Id:    mediaId,
		//})

		return nil
	}
}
