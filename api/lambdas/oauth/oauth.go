package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tencentyun/scf-go-lib/events"
	common2 "github.com/thomasduchatelle/ephoto/api/lambdas/common"
	"github.com/thomasduchatelle/ephoto/pkg/acl/aclcore"
	"strings"
)

const bearerPrefix = "bearer "

type identityDTO struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

var (
	authenticator SSOAuthenticator
)

type SSOAuthenticator interface {
	AuthenticateFromExternalIDProvider(identityJWT string) (*aclcore.Authentication, *aclcore.Identity, error)
}

func Handler(request events.APIGatewayRequest) (common2.Response, error) {
	authorisation, ok := request.Headers["authorization"]
	authorisation = strings.Trim(authorisation, " ")
	if !ok {
		return common2.NewJsonResponse(401, map[string]string{
			"error": "Authorization header is required.",
		}, nil)
	}

	if !strings.HasPrefix(strings.ToLower(authorisation), bearerPrefix) {
		return common2.BadRequest(map[string]string{
			"error": "Authorization header must be of Bearer type.",
		})
	}

	tokenString := authorisation[len(bearerPrefix):]
	authentication, identity, err := authenticator.AuthenticateFromExternalIDProvider(tokenString)
	if err != nil {
		log.WithError(err).Infof("Authentication rejected: %s", tokenString)

		code, status := lookupCode(err)
		return common2.NewJsonResponse(status, map[string]string{
			"code":  code,
			"error": err.Error(),
		}, nil)
	}

	return common2.Ok(map[string]interface{}{
		"access_token": authentication.AccessToken,
		"expires_in":   authentication.ExpiresIn,
		"identity": identityDTO{
			Email:   identity.Email,
			Name:    identity.Name,
			Picture: identity.Picture,
		},
	})
}

func lookupCode(err error) (string, int) {
	switch {
	case errors.Is(err, aclcore.NotPreregisteredError):
		return "oauth.user-not-preregistered", 403
	case errors.Is(err, aclcore.InvalidTokenError) || errors.Is(err, aclcore.InvalidTokenExplicitError):
		return "oauth.invalid-token", 403
	default:
		return "", 500
	}
}

func main() {
	authenticator = common2.NewSSOAuthenticator()

	lambda.Start(Handler)
}
