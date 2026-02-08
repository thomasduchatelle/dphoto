package aclcore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// LoadUserScopes fetches user scopes from the database and returns them in different formats.
// This is the canonical implementation used by both legacy token generation and Cognito token validation.
//
// Logic:
// 1. First tries to fetch ApiScope and MainOwnerScope
// 2. If no high-level scopes found, checks for AlbumVisitorScope/MediaVisitorScope and returns "visitor" scope
// 3. If no scopes at all, returns NotPreregisteredError
//
// Returns:
// - scopeStrings: Array of scope strings for JWT encoding (e.g., ["api:admin", "owner:tony@stark.com"])
// - scopeMap: Map of scope strings for Claims.Scopes field
// - owner: Extracted owner from MainOwnerScope (or nil if not found)
// - error: NotPreregisteredError if user has no scopes, or database error
func LoadUserScopes(ctx context.Context, scopesReader ScopesReader, userId usermodel.UserId) (scopeStrings []string, scopeMap map[string]interface{}, owner *ownermodel.Owner, err error) {
	// First, try to get ApiScope and MainOwnerScope
	grants, err := scopesReader.ListScopesByUser(ctx, userId, ApiScope, MainOwnerScope)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to list scopes for user %s", userId)
	}

	scopeMap = make(map[string]interface{})

	for _, grant := range grants {
		var scopeStr string
		switch grant.Type {
		case ApiScope:
			scopeStr = fmt.Sprintf("api:%s", grant.ResourceId)

		case MainOwnerScope:
			scopeStr = fmt.Sprintf("%s%s", JWTScopeOwnerPrefix, grant.ResourceOwner)
			// Set the owner (use the first MainOwnerScope found)
			if owner == nil {
				ownerValue := grant.ResourceOwner
				owner = &ownerValue
			}
		}

		if scopeStr != "" {
			scopeStrings = append(scopeStrings, scopeStr)
			scopeMap[scopeStr] = nil
		}
	}

	// If we found scopes, return them
	if len(scopeStrings) > 0 {
		return scopeStrings, scopeMap, owner, nil
	}

	// Second chance for visitors: check if user has album/media visitor scopes
	visitorGrants, err := scopesReader.ListScopesByUser(ctx, userId, AlbumVisitorScope, MediaVisitorScope)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to list visitor scopes for user %s", userId)
	}

	if len(visitorGrants) > 0 {
		scopeStrings = []string{"visitor"}
		scopeMap = map[string]interface{}{"visitor": nil}
		return scopeStrings, scopeMap, nil, nil
	}

	// User has no scopes - not pre-registered
	return nil, nil, nil, NotPreregisteredError
}
