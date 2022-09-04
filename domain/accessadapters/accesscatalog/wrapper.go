// Package accesscatalog wraps both access and catalog packages to merge shared resources to the catalog
package accesscatalog

import (
	"github.com/thomasduchatelle/dphoto/domain/access"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
)

func FindAccessibleAlbums(email string) ([]*catalog.Album, error) {
	grants, err := access.ListGrants(email, access.AlbumResource)
	if err != nil {
		return nil, err
	}

	requests := make([]catalog.GetAlbum, len(grants)+1, len(grants)+1)
	requests[0] = catalog.GetAlbum{
		Owner: email,
	}
	for i, grant := range grants {
		requests[i+1] = catalog.GetAlbum{
			Owner:      grant.Owner,
			FolderName: grant.Id,
		}
	}

	sharedMedia, oldest, mostRecent, err := access.CountGrants(email, access.MediaResource)
	if err != nil {
		return nil, err
	}

	albums, err := catalog.FindAlbums(requests...)
	if sharedMedia > 0 {
		sharedWithMe := &catalog.Album{
			Owner:      email,
			Name:       "shared",
			FolderName: "shared",
			Start:      oldest,
			End:        mostRecent,
			TotalCount: sharedMedia,
		}
		albums = append([]*catalog.Album{sharedWithMe}, albums...)
	}

	return albums, err
}
