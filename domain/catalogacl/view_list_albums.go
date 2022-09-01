package catalogacl

import (
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/accesscontrol"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
)

type AlbumId struct {
	Owner      string
	FolderName string
}

func (v *viewimpl) ListAlbums(filter ListAlbumsFilter) ([]*AlbumInView, error) {
	accessibleAlbums, err := v.listAccessibleAlbums(filter)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't list accessible albums of '%s'", v.userEmail)
	}

	if !filter.OnlyDirectlyOwned {
		sharedMedia, oldest, mostRecent, err := accesscontrol.CountUserPermissions(v.userEmail, accesscontrol.MediaVisitorScope)
		if err != nil {
			return nil, err
		}

		if sharedMedia > 0 {
			sharedWithMe := &AlbumInView{
				Album: &catalog.Album{
					Owner:      v.userEmail,
					Name:       "Shared with me",
					FolderName: "",
					Start:      oldest,
					End:        mostRecent,
					TotalCount: sharedMedia,
				},
				SharedTo: nil,
				SharedBy: "",
			}
			accessibleAlbums = append([]*AlbumInView{sharedWithMe}, accessibleAlbums...)
		}
	}

	return accessibleAlbums, err
}

func (v *viewimpl) listAccessibleAlbums(filter ListAlbumsFilter) ([]*AlbumInView, error) {
	roles, err := accesscontrol.ScopesReader(v.userEmail, accesscontrol.MainOwnerScope, accesscontrol.AlbumVisitorScope)
	if err != nil {
		return nil, err
	}

	var ownerOf []string
	var requests []catalog.ListAlbumsInput
	sharedByOwner := make(map[AlbumId]string)
	for _, role := range roles {
		switch role.Type {
		case accesscontrol.MainOwnerScope:
			ownerOf = append(ownerOf, role.ResourceOwner)
			requests = append(requests, catalog.ListAlbumsInput{
				Owner: role.ResourceOwner,
			})

		case accesscontrol.AlbumVisitorScope:
			if !filter.OnlyDirectlyOwned {
				sharedByOwner[NewAlbumId(role.ResourceOwner, role.ResourceId)] = role.ResourceOwner
				requests = append(requests, catalog.ListAlbumsInput{
					Owner:      role.ResourceOwner,
					FolderName: role.ResourceId,
				})
			}
		}
	}

	sharing, err := v.listAlbumSharing(ownerOf)
	if err != nil {
		return nil, err
	}

	accessibleAlbums, err := catalog.ListAlbums(requests...)
	if err != nil {
		return nil, err
	}

	var view []*AlbumInView
	for _, album := range accessibleAlbums {
		sharedTo, _ := sharing[NewAlbumId(album.Owner, album.FolderName)]
		sharedBy, _ := sharedByOwner[NewAlbumId(album.Owner, album.FolderName)]
		view = append(view, &AlbumInView{
			Album:    album,
			SharedTo: sharedTo,
			SharedBy: sharedBy,
		})
	}

	return view, err
}

func (v *viewimpl) listAlbumSharing(owners []string) (map[AlbumId][]string, error) {
	permissions, err := accesscontrol.ListResourcesPermissionsByOwner(owners, accesscontrol.AlbumVisitorScope)

	sharingWith := make(map[AlbumId][]string)
	for _, permission := range permissions {
		sharedTo, _ := sharingWith[NewAlbumId(permission.ResourceOwner, permission.ResourceId)]
		sharingWith[NewAlbumId(permission.ResourceOwner, permission.ResourceId)] = append(sharedTo, permission.GrantedTo)
	}

	return sharingWith, err
}

func NewAlbumId(owner, folderName string) AlbumId {
	return AlbumId{Owner: owner, FolderName: folderName}
}
