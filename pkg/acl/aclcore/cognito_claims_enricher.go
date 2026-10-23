package aclcore

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// CognitoClaimsEnricher resolves the DPhoto permissions of a user to be stamped into the Cognito access token
// at token-generation time (pre-token-generation trigger). It lets the Lambda Authorizer trust the token claims
// instead of resolving the email and scopes on every request.
type CognitoClaimsEnricher struct {
	ScopesReader ScopesReader
}

// EnrichedClaims are the custom claims added to a Cognito access token so the Lambda Authorizer can build the
// Claims without any UserInfo call nor database lookup.
type EnrichedClaims struct {
	Email  string
	Owner  string // Owner is empty when the user is not a main owner
	Scopes string
}

// Enrich resolves the scopes of the user identified by its email. It returns NotPreregisteredError when the user
// has no scope at all, so the token generation can be denied for unknown users.
func (e *CognitoClaimsEnricher) Enrich(ctx context.Context, email string) (EnrichedClaims, error) {
	userId := usermodel.NewUserId(email)
	if err := userId.IsValid(); err != nil {
		return EnrichedClaims{}, err
	}

	scopeStrings, _, owner, err := LoadUserScopes(ctx, e.ScopesReader, userId)
	if err != nil {
		return EnrichedClaims{}, errors.Wrapf(err, "failed to resolve scopes for %s", email)
	}

	ownerValue := ""
	if owner != nil {
		ownerValue = owner.Value()
	}

	return EnrichedClaims{
		Email:  userId.Value(),
		Owner:  ownerValue,
		Scopes: strings.Join(scopeStrings, " "),
	}, nil
}
