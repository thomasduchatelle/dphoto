package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tencentyun/scf-go-lib/events"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	"github.com/thomasduchatelle/dphoto/domain/accesscontrol"
	"strings"
)

const bearerPrefix = "bearer "

type identityDTO struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

var (
	authenticator *accesscontrol.SSOAuthenticator
)

func Handler(request events.APIGatewayRequest) (common.Response, error) {
	authorisation, ok := request.Headers["authorization"]
	authorisation = strings.Trim(authorisation, " ")
	if !ok {
		return common.NewJsonResponse(401, map[string]string{
			"error": "Authorization header is required.",
		}, nil)
	}

	if !strings.HasPrefix(strings.ToLower(authorisation), bearerPrefix) {
		return common.BadRequest(map[string]string{
			"error": "Authorization header must be of Bearer type.",
		})
	}

	tokenString := authorisation[len(bearerPrefix):]
	authentication, identity, err := authenticator.AuthenticateFromExternalIDProvider(tokenString)
	if err != nil {
		log.WithError(err).Infof("Authentication rejected: %+v", request)

		code, status := lookupCode(err)
		return common.NewJsonResponse(status, map[string]string{
			"code":  code,
			"error": err.Error(),
		}, nil)
	}

	return common.Ok(map[string]interface{}{
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
	case errors.Is(err, accesscontrol.NotPreregisteredError):
		return "oauth.user-not-preregistered", 403
	case errors.Is(err, accesscontrol.InvalidTokenError) || errors.Is(err, accesscontrol.InvalidTokenExplicitError):
		return "oauth.invalid-token", 403
	default:
		return "", 500
	}
}

func main() {
	authenticator = common.NewSSOAuthenticator()

	lambda.Start(Handler)
}
