package common

import (
	"github.com/thomasduchatelle/dphoto/domain/accesscontrol"
)

func Or(expressions ...func(claims accesscontrol.Claims) error) func(claims accesscontrol.Claims) error {
	return func(claims accesscontrol.Claims) error {
		for _, expr := range expressions {
			err := expr(claims)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func CanReadAsOwner(owner string) func(claims accesscontrol.Claims) error {
	return func(claims accesscontrol.Claims) error {
		return claims.IsOwnerOf(owner)
	}
}

func CanReadAlbum(owner, album string) func(claims accesscontrol.Claims) error {
	return func(claims accesscontrol.Claims) error {
		//granted, err := access.IsGranted("foobar", access.Resource{
		//	Type:  access.AlbumResource,
		//	Owner: owner,
		//	Id:    album,
		//})
		return nil
	}
}

func CanReadMedia(owner, mediaId string) func(claims accesscontrol.Claims) error {
	return func(claims accesscontrol.Claims) error {
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
