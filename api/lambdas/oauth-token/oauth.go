package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tencentyun/scf-go-lib/events"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"net/url"
	"strings"
)

type identityDTO struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

var (
	ssoAuthenticator     SSOAuthenticator
	refreshAuthenticator RefreshAuthenticator
)

type SSOAuthenticator interface {
	AuthenticateFromExternalIDProvider(identityJWT string, refreshTokenPurpose aclcore.RefreshTokenPurpose) (*aclcore.Authentication, *aclcore.Identity, error)
}

type RefreshAuthenticator interface {
	AuthenticateFromAccessToken(refreshToken string) (*aclcore.Authentication, *aclcore.Identity, error)
}

func Handler(request events.APIGatewayRequest) (common.Response, error) {
	body, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		return common.InternalError(errors.New("lambda body was expected to be base 64 encoded"))
	}
	grantType, scopes, attributes, err := parseBody(request.Headers, string(body))
	if err != nil {
		log.WithError(err).Infof("invalid request: %s", request.Body)
		return common.BadRequest(err.Error())
	}

	switch grantType {
	case "identity":
		return authenticateFromIdentityToken(scopes, attributes)

	case "refresh_token":
		return authenticateFromRefreshToken(scopes, attributes)

	default:
		log.Warnf("invalid grant type received: '%s'", grantType)
		return common.BadRequest(fmt.Sprintf("%s grant type is not supported. Supported are only 'identity' and 'refresh_token'", grantType))
	}
}

func parseBody(headers map[string]string, body string) (string, []string, map[string]string, error) {
	contentType, _ := headers["content-type"]
	if contentType != "application/x-www-form-urlencoded" {
		return "", nil, nil, errors.Errorf("'%s' encoding is not supported. Supported encoding are: application/x-www-form-urlencoded", contentType)
	}

	const grantTypeKey = "grant_type"
	query, err := url.ParseQuery(body)
	if err != nil {
		return "", nil, nil, errors.Wrapf(err, "body is invalid [%s]", body)
	} else if !query.Has(grantTypeKey) {
		return "", nil, nil, errors.Errorf("body must define the grant_type %+v", query)
	}

	args := make(map[string]string)
	for k, v := range query {
		if len(v) > 0 && k != grantTypeKey && k != "scope" {
			args[strings.ToLower(k)] = v[0]
		}
	}

	return query.Get(grantTypeKey), strings.Split(query.Get("scope"), " "), args, nil
}

func authenticateFromIdentityToken(scopes []string, attributes map[string]string) (common.Response, error) {
	identityToken, ok := attributes["identity_token"]
	identityToken = strings.Trim(identityToken, " ")
	if !ok || identityToken == "" {
		return common.NewJsonResponse(400, map[string]string{
			"error": "'identity_token' is required'",
		}, nil)
	}

	authentication, identity, err := ssoAuthenticator.AuthenticateFromExternalIDProvider(identityToken, aclcore.RefreshTokenPurposeWeb)
	if err != nil {
		log.WithError(err).Infof("Identity token rejected: %s", identityToken)

		code, status := lookupCode(err)
		return common.NewJsonResponse(status, map[string]string{
			"code":  code,
			"error": err.Error(),
		}, nil)
	}

	return toOkResponse(authentication, identity)
}

func authenticateFromRefreshToken(scopes []string, attributes map[string]string) (common.Response, error) {
	refreshToken, ok := attributes["refresh_token"]
	refreshToken = strings.Trim(refreshToken, " ")
	if !ok || refreshToken == "" {
		return common.NewJsonResponse(400, map[string]string{
			"error": "'refresh_token' is required'",
		}, nil)
	}

	tokens, identity, err := refreshAuthenticator.AuthenticateFromAccessToken(refreshToken)
	if err != nil {
		log.WithError(err).Infof("Refresh token rejected")

		code, status := lookupCode(err)
		return common.NewJsonResponse(status, map[string]string{
			"code":  code,
			"error": err.Error(),
		}, nil)
	}

	return toOkResponse(tokens, identity)
}

func toOkResponse(authentication *aclcore.Authentication, identity *aclcore.Identity) (common.Response, error) {
	return common.Ok(map[string]interface{}{
		"token_type":    "Bearer",
		"access_token":  authentication.AccessToken,
		"refresh_token": authentication.RefreshToken,
		"expires_in":    authentication.ExpiresIn,
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
	case errors.Is(err, aclcore.ExpiredRefreshTokenError):
		return "oauth.refresh.expired", 403
	case errors.Is(err, aclcore.InvalidRefreshTokenError):
		return "oauth.refresh.invalid", 403
	default:
		return "", 500
	}
}

func main() {
	ssoAuthenticator, refreshAuthenticator = common.NewAuthenticators()

	lambda.Start(Handler)
}
