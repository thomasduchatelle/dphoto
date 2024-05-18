package aclcore

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type IdentityQueriesIdentityRepository interface {
	FindIdentities(emails []usermodel.UserId) ([]*Identity, error)
}

type IdentityQueriesScopeRepository interface {
	ListScopesByOwners(ctx context.Context, owners []ownermodel.Owner, types ...ScopeType) ([]*Scope, error)
}

type IdentityQueries struct {
	IdentityRepository IdentityQueriesIdentityRepository
	ScopeRepository    IdentityQueriesScopeRepository
}

func (i *IdentityQueries) FindIdentities(emails []usermodel.UserId) ([]*Identity, error) {
	return i.IdentityRepository.FindIdentities(emails)
}

func (i *IdentityQueries) FindOwnerIdentities(owners []ownermodel.Owner) (map[ownermodel.Owner][]*Identity, error) {
	ctx := context.TODO()
	scopes, err := i.ScopeRepository.ListScopesByOwners(ctx, owners, MainOwnerScope)
	if err != nil || len(scopes) == 0 {
		return nil, err
	}

	var emails []usermodel.UserId
	for _, scope := range scopes {
		emails = append(emails, scope.GrantedTo)
	}

	identities, err := i.IdentityRepository.FindIdentities(emails)

	ownerIdentities := make(map[ownermodel.Owner][]*Identity)
	for _, scope := range scopes {
		identity := &Identity{
			Email: scope.GrantedTo,
			Name:  scope.GrantedTo.Value(),
		}
		for _, i := range identities {
			if i.Email == scope.GrantedTo {
				identity = i
			}
		}

		current, _ := ownerIdentities[scope.ResourceOwner]
		ownerIdentities[scope.ResourceOwner] = append(current, identity)
	}

	return ownerIdentities, err
}
