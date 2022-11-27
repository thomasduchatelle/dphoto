package common

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/accesscontrol"
	"github.com/thomasduchatelle/dphoto/domain/catalogacl"
	"github.com/thomasduchatelle/dphoto/domain/catalogacladapters/cacl2ac"
	"strings"
)

// RequiresCatalogACL read the bearer token, parse it, and creates a Catalog View from it (enforcing catalog ACL).
//
// accesscontrol.AccessUnauthorisedError errors will be converted into 401 responses
// accesscontrol.AccessForbiddenError errors will be converted into 403 responses
// any other errors will be converted into 500
func RequiresCatalogACL(request *events.APIGatewayProxyRequest, process func(catalogView *catalogacl.View) (Response, error)) (Response, error) {
	token, err := readToken(request)
	if err != nil {
		return UnauthorizedResponse(err.Error())
	}

	claims, err := jwtDecoder.Decode(token)
	if err != nil {
		return UnauthorizedResponse(err.Error())
	}

	catalogView := &catalogacl.View{
		UserEmail:     claims.Subject,
		AccessControl: cacl2ac.NewAccessControlAdapterFromToken(grantRepository, claims),
	}

	response, err := process(catalogView)

	switch {
	case errors.Is(err, accesscontrol.AccessUnauthorisedError):
		return UnauthorizedResponse(err.Error())

	case errors.Is(err, accesscontrol.AccessForbiddenError):
		return ForbiddenResponse(err.Error())

	case err != nil:
		return InternalError(err)

	default:
		return response, nil
	}
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

func UnauthorizedResponse(message string) (Response, error) {
	response, _ := NewJsonResponse(401, map[string]string{
		"error": fmt.Sprintf("access unauthorised: %s", message),
	}, nil)
	return response, nil
}

func ForbiddenResponse(message string) (Response, error) {
	response, _ := NewJsonResponse(403, map[string]string{
		"error": fmt.Sprintf("access forbidden: %s", message),
	}, nil)
	return response, nil
}
