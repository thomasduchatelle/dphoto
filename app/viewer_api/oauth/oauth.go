package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/tencentyun/scf-go-lib/events"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	"github.com/thomasduchatelle/dphoto/domain/oauth"
	"strings"
)

const bearerPrefix = "bearer "

type identityDTO struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

var (
	oauthAuthenticate       = oauth.Authenticate
	isNotPreregisteredError = oauth.IsNotPreregisteredError
	isInvalidTokenError     = oauth.IsInvalidTokenError
)

func Handler(request events.APIGatewayRequest) (common.Response, error) {
	authorisation, ok := request.Headers["authorization"]
	if !ok {
		return common.NewJsonResponse(401, map[string]string{
			"error": "Authorization header is required.",
		}, nil)
	}

	if !strings.HasPrefix(strings.Trim(strings.ToLower(authorisation), " "), bearerPrefix) {
		return common.BadRequest(map[string]string{
			"error": "Authorization header must be of Bearer type.",
		})
	}

	tokenString := strings.Trim(authorisation, " ")[len(bearerPrefix):]
	authentication, identity, err := oauthAuthenticate(tokenString)
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
	case isNotPreregisteredError(err):
		return "oauth.user-not-preregistered", 403
	case isInvalidTokenError(err):
		return "oauth.invalid-token", 403
	default:
		return "", 500
	}
}

func main() {
	common.BootstrapOAuthDomain()

	lambda.Start(Handler)
}
