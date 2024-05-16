package aclcore

import (
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type CoreRules struct {
	ScopeReader ScopesReader
	Email       usermodel.UserId
}

// Owner returns empty if the user own nothing, or the identifier of its owner
func (a *CoreRules) Owner() (*ownermodel.Owner, error) {
	scopes, err := a.ScopeReader.ListUserScopes(a.Email, MainOwnerScope)
	if err != nil {
		return nil, err
	}

	if len(scopes) == 0 {
		return nil, errors.Errorf("%s is not a main user.", a.Email)
	}

	return &scopes[0].ResourceOwner, nil
}
