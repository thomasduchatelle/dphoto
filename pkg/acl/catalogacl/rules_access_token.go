package catalogacl

import (
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// OptimiseRulesWithAccessToken creates an adapter catalogacl -> aclcore with pre-authorised scopes from the Access Token
func OptimiseRulesWithAccessToken(delegate CatalogRules, claims aclcore.Claims) CatalogRules {
	return &rulesWithAccessToken{
		CatalogRules: delegate,
		owner:        claims.Owner,
	}
}

type rulesWithAccessToken struct {
	CatalogRules
	owner ownermodel.Owner
}

func (a *rulesWithAccessToken) Owner() (*ownermodel.Owner, error) {
	if a.owner == "" {
		return nil, nil // TODO Highly irregular code style !
	}
	return &a.owner, nil
}

func (a *rulesWithAccessToken) CanListMediasFromAlbum(albumId catalog.AlbumId) error {
	if a.owner == albumId.Owner {
		return nil
	}

	return a.CatalogRules.CanListMediasFromAlbum(albumId)
}

func (a *rulesWithAccessToken) CanReadMedia(owner ownermodel.Owner, mediaId catalog.MediaId) error {
	if owner == a.owner {
		return nil
	}

	return a.CatalogRules.CanReadMedia(owner, mediaId)
}

func (a *rulesWithAccessToken) CanManageAlbum(albumId catalog.AlbumId) error {
	if a.owner == albumId.Owner {
		return nil
	}

	return a.CatalogRules.CanManageAlbum(albumId)
}
