package aclcore

type IdentityQueriesIdentityRepository interface {
	FindIdentities(emails []string) ([]*Identity, error)
}

type IdentityQueriesScopeRepository interface {
	ListScopesByOwners(owners []string, types ...ScopeType) ([]*Scope, error)
}

type IdentityQueries struct {
	IdentityRepository IdentityQueriesIdentityRepository
	ScopeRepository    IdentityQueriesScopeRepository
}

func (i *IdentityQueries) FindIdentities(emails []string) ([]*Identity, error) {
	return i.IdentityRepository.FindIdentities(emails)
}

func (i *IdentityQueries) FindOwnerIdentities(owners []string) (map[string][]*Identity, error) {
	scopes, err := i.ScopeRepository.ListScopesByOwners(owners, MainOwnerScope)
	if err != nil || len(scopes) == 0 {
		return nil, err
	}

	var emails []string
	for _, scope := range scopes {
		emails = append(emails, scope.GrantedTo)
	}

	identities, err := i.IdentityRepository.FindIdentities(emails)

	ownerIdentities := make(map[string][]*Identity)
	for _, scope := range scopes {
		identity := &Identity{
			Email: scope.GrantedTo,
			Name:  scope.GrantedTo,
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
