package aclcore

import "github.com/pkg/errors"

type CoreRules struct {
	ScopeReader ScopesReader
	Email       string
}

// Owner returns empty if the user own nothing, or the identifier of its owner
func (a *CoreRules) Owner() (string, error) {
	scopes, err := a.ScopeReader.ListUserScopes(a.Email, MainOwnerScope)
	if err != nil {
		return "", err
	}

	if len(scopes) == 0 {
		return "", errors.Errorf("%s is not a main user.", a.Email)
	}

	return scopes[0].ResourceOwner, nil
}
