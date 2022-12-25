package aclcore

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

type CreateUser struct {
	ScopesReader ScopesReader
	ScopeWriter  ScopeWriter
}

// CreateUser create a user capable of backup as 'owner', or update an existing owner to be 'owner'
func (c *CreateUser) CreateUser(email, ownerOptional string) error {
	email = strings.Trim(email, " ")
	if email == "" {
		return InvalidUserEmailError
	}

	owner := ownerOptional
	if owner == "" {
		owner = email
	}

	scopes, err := c.ScopesReader.ListUserScopes(email, MainOwnerScope)
	if err != nil {
		return err
	}

	presents := false
	if len(scopes) > 0 {
		var idsToDelete []ScopeId
		for _, scope := range scopes {
			if scope.Type == MainOwnerScope && scope.GrantedTo == email && scope.ResourceOwner == owner && scope.ResourceId == "" {
				presents = true

			} else if scope.Type == MainOwnerScope && scope.GrantedTo == email {
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
		GrantedTo:     email,
		ResourceOwner: owner,
	})
}
