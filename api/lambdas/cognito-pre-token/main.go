package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	_ "github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

// Handler is the Cognito Pre Token Generation (V2) trigger. It resolves the DPhoto scopes and owner of the user
// signing in and stamps them into the access token so the Lambda Authorizer can trust the token without any
// UserInfo call nor database lookup. A user with no DPhoto scope is rejected (login denied).
func Handler(ctx context.Context, event events.CognitoEventUserPoolsPreTokenGenV2) (events.CognitoEventUserPoolsPreTokenGenV2, error) {
	email := event.Request.UserAttributes["email"]

	enricher := &aclcore.CognitoClaimsEnricher{
		ScopesReader: pkgfactory.AclRepository(ctx),
	}

	enriched, err := enricher.Enrich(ctx, email)
	if err != nil {
		if errors.Is(err, aclcore.NotPreregisteredError) {
			log.WithField("Email", email).Warn("Denying token generation: user is not pre-registered")
		} else {
			log.WithError(err).WithField("Email", email).Error("Failed to enrich token claims")
		}
		return event, err
	}

	claims := map[string]string{
		"dphoto:email":  enriched.Email,
		"dphoto:scopes": enriched.Scopes,
	}
	if enriched.Owner != "" {
		claims["dphoto:owner"] = enriched.Owner
	}

	event.Response.ClaimsAndScopeOverrideDetails = events.ClaimsAndScopeOverrideDetails{
		AccessTokenGeneration: events.AccessTokenGeneration{
			ClaimsToAddOrOverride: claims,
		},
	}

	log.WithField("Email", enriched.Email).WithField("Owner", enriched.Owner).Info("Enriched access token claims")
	return event, nil
}

func main() {
	lambda.Start(Handler)
}
