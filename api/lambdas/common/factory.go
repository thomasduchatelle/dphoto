package common

import (
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

func GetIdentityQueries() *aclcore.IdentityQueries {
	return &aclcore.IdentityQueries{
		IdentityRepository: getIdentityDetailsStore(),
		ScopeRepository:    grantRepository,
	}
}
