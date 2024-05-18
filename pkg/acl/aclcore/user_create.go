package aclcore

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type CreateUser struct {
	ScopesReader ScopesReader
	ScopeWriter  ScopeWriter
}

// CreateUser create a user capable of backup as 'owner', or update an existing owner to be 'owner'
func (c *CreateUser) CreateUser(email, ownerOptional string) error {
	ctx := context.TODO()

	userId := usermodel.NewUserId(email)
	if err := userId.IsValid(); err != nil {
		return err
	}

	owner := ownermodel.Owner(ownerOptional)
	if owner == "" {
		owner = ownermodel.Owner(email)
	}

	scopes, err := c.ScopesReader.ListScopesByUser(ctx, userId, MainOwnerScope)
	if err != nil {
		return err
	}

	presents := false
	if len(scopes) > 0 {
		var idsToDelete []ScopeId
		for _, scope := range scopes {
			if scope.Type == MainOwnerScope && scope.GrantedTo == userId && scope.ResourceOwner == owner && scope.ResourceId == "" {
				presents = true

			} else if scope.Type == MainOwnerScope && scope.GrantedTo == userId {
				idsToDelete = append(idsToDelete, ScopeId{
					Type:          scope.Type,
					GrantedTo:     scope.GrantedTo,
					ResourceOwner: scope.ResourceOwner,
					ResourceId:    scope.ResourceId,
				})
			}
		}

		if len(idsToDelete) > 0 {
			log.WithField("User", email).Infof("Revoking invalid scopes: %+v", idsToDelete)
			err = c.ScopeWriter.DeleteScopes(idsToDelete...)
		}
		if err != nil {
			return err
		}
	}

	if presents {
		return nil
	}

	log.WithField("Owner", owner).Infof("Creating new user %s, backup capable as owner '%s'", email, owner)

	return c.ScopeWriter.SaveIfNewScope(Scope{
		Type:          MainOwnerScope,
		GrantedAt:     TimeFunc(),
		GrantedTo:     userId,
		ResourceOwner: owner,
	})
}
