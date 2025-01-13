package catalogacl

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

var (
	ErrAccessDenied = errors.New("access denied")
)

type HasPermissionPort interface {
	FindScopesByIdCtx(ctx context.Context, ids ...aclcore.ScopeId) ([]*aclcore.Scope, error)
}

type CatalogQueriesPort interface {
	FindMediaOwnership(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId) (*catalog.AlbumId, error)
}

type CatalogAuthorizer struct {
	HasPermissionPort  HasPermissionPort
	CatalogQueriesPort CatalogQueriesPort
}

func (a *CatalogAuthorizer) IsAuthorisedToListMedias(ctx context.Context, userId usermodel.CurrentUser, albumId catalog.AlbumId) error {
	if userId.Owner != nil && *userId.Owner == albumId.Owner {
		return nil
	}

	permissions, err := a.HasPermissionPort.FindScopesByIdCtx(ctx, aclcore.ScopeId{
		Type:          aclcore.AlbumVisitorScope,
		GrantedTo:     userId.UserId,
		ResourceOwner: albumId.Owner,
		ResourceId:    albumId.FolderName.String(),
	})
	if err != nil {
		return errors.Wrapf(err, "failed to check permissions for user %s on album %s", userId.UserId, albumId)
	}

	if len(permissions) > 0 {
		return nil
	}

	return errors.Wrapf(ErrAccessDenied, "user %s is not authorised to list medias from album %s", userId.UserId, albumId)
}

func (a *CatalogAuthorizer) IsAuthorisedToViewMedia(ctx context.Context, currentUser usermodel.CurrentUser, owner ownermodel.Owner, mediaId catalog.MediaId) error {
	if currentUser.Owner != nil && *currentUser.Owner == owner {
		return nil
	}

	albumId, err := a.CatalogQueriesPort.FindMediaOwnership(ctx, owner, mediaId)
	if err != nil {
		return errors.Wrapf(aclcore.AccessForbiddenError, err.Error())
	}

	scopes, err := a.HasPermissionPort.FindScopesByIdCtx(
		ctx,
		aclcore.ScopeId{
			Type:          aclcore.MainOwnerScope,
			GrantedTo:     currentUser.UserId,
			ResourceOwner: owner,
		},
		aclcore.ScopeId{
			Type:          aclcore.AlbumVisitorScope,
			GrantedTo:     currentUser.UserId,
			ResourceOwner: owner,
			ResourceId:    albumId.FolderName.String(),
		},
		aclcore.ScopeId{
			Type:          aclcore.MediaVisitorScope,
			GrantedTo:     currentUser.UserId,
			ResourceOwner: owner,
			ResourceId:    mediaId.Value(),
		},
	)
	if err != nil {
		return err
	}
	if len(scopes) == 0 {
		return errors.Wrapf(aclcore.AccessForbiddenError, "reading media %s/%s has been denied.", owner, mediaId)
	}

	return nil
}

func (a *CatalogAuthorizer) CanShareAlbum(ctx context.Context, user usermodel.CurrentUser, albumId catalog.AlbumId) error {
	if user.Owner != nil && *user.Owner == albumId.Owner {
		return nil
	}

	permissions, err := a.HasPermissionPort.FindScopesByIdCtx(ctx, aclcore.ScopeId{
		Type:          aclcore.MainOwnerScope,
		GrantedTo:     user.UserId,
		ResourceOwner: albumId.Owner,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to check permissions for user %s on album %s", user.UserId, albumId)
	}
	if len(permissions) > 0 {
		return nil
	}

	return aclcore.AccessForbiddenError
}

func (a *CatalogAuthorizer) CanCreateAlbum(ctx context.Context, user usermodel.CurrentUser, owner ownermodel.Owner) error {
	if user.Owner != nil && *user.Owner == owner {
		return nil
	}

	permissions, err := a.HasPermissionPort.FindScopesByIdCtx(ctx, aclcore.ScopeId{
		Type:          aclcore.MainOwnerScope,
		GrantedTo:     user.UserId,
		ResourceOwner: owner,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to check permissions for user %s for owner %s", user.UserId, owner)
	}
	if len(permissions) > 0 {
		return nil
	}

	return aclcore.AccessForbiddenError
}
