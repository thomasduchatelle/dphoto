package common

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/oauth"
	"strings"
)

// ValidateRequest make sure authorisation is valid for the expected actions. Return [response, TRUE] when the consumer is not authorised.
func ValidateRequest(request events.APIGatewayProxyRequest, query *oauth.AuthoriseQuery) (Response, bool) {
	tokenString, err := readToken(request)
	if err != nil {
		response, _ := NewJsonResponse(401, map[string]string{
			"error": fmt.Sprintf("access denied: %s", err.Error()),
		}, nil)
		return response, true
	}

	_, err = oauth.Authorise(tokenString, query)
	if err != nil {
		response, _ := NewJsonResponse(403, map[string]string{
			"error": fmt.Sprintf("access denied: %s", err.Error()),
		}, nil)
		return response, true
	}

	return Response{}, false
}

func readToken(request events.APIGatewayProxyRequest) (string, error) {
	authorisation, ok := request.Headers["authorization"]
	if !ok {
		return "", errors.Errorf("no authorization header")
	}

	bearerPrefix := "bearer "
	if !strings.HasPrefix(strings.ToLower(authorisation), bearerPrefix) {
		return "", errors.Errorf("authorization header should start by 'Bearer '")
	}

	return authorisation[len(bearerPrefix):], nil
}
