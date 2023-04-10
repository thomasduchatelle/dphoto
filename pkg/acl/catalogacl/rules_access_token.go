package catalogacl

import (
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
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
	owner string
}

func (a *rulesWithAccessToken) Owner() (string, error) {
	return a.owner, nil
}

func (a *rulesWithAccessToken) CanListMediasFromAlbum(owner string, folderName string) error {
	if a.owner == owner {
		return nil
	}

	return a.CatalogRules.CanListMediasFromAlbum(owner, folderName)
}

func (a *rulesWithAccessToken) CanReadMedia(owner string, mediaId string) error {
	if a.owner == owner {
		return nil
	}

	return a.CatalogRules.CanReadMedia(owner, mediaId)
}

func (a *rulesWithAccessToken) CanManageAlbum(owner string, folderName string) error {
	if a.owner == owner {
		return nil
	}

	return a.CatalogRules.CanManageAlbum(owner, folderName)
}
