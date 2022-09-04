package common

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/accessadapters/oauth"
	"strings"
)

// ValidateRequest make sure authorisation is valid for the expected actions. Return [response, TRUE] when the consumer is not authorised.
func ValidateRequest(request *events.APIGatewayProxyRequest, authoriseFct func(claims oauth.Claims) error) (Response, bool) {
	tokenString, err := readToken(request)
	if err != nil {
		response, _ := NewJsonResponse(401, map[string]string{
			"error": fmt.Sprintf("access denied: %s", err.Error()),
		}, nil)
		return response, true
	}

	err = OAuthClient.Authorise(tokenString, authoriseFct)
	if err != nil {
		response, _ := NewJsonResponse(403, map[string]string{
			"error": fmt.Sprintf("access forbidden: %s", err.Error()),
		}, nil)
		return response, true
	}

	return Response{}, false
}

func readToken(request *events.APIGatewayProxyRequest) (string, error) {
	authorisation, ok := request.Headers["authorization"]
	if !ok {
		// allow to pass tokens in the Query parameters for images loaded in <img /> tags
		if token, withAccessToken := request.QueryStringParameters["access_token"]; withAccessToken {
			authorisation = "bearer " + token

		} else {
			return "", errors.Errorf("no authorization header")
		}
	}

	bearerPrefix := "bearer "
	if !strings.HasPrefix(strings.ToLower(authorisation), bearerPrefix) {
		return "", errors.Errorf("authorization header should start by 'Bearer '")
	}

	return authorisation[len(bearerPrefix):], nil
}
