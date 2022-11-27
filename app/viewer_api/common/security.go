package common

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/accesscontrol"
	"github.com/thomasduchatelle/dphoto/domain/catalogacl"
	"github.com/thomasduchatelle/dphoto/domain/catalogacladapters/catalogacllogic"
	"strings"
)

// RequiresCatalogACL parses the token and instantiate the ACL extension for 'catalog' domain
//
// accesscontrol.AccessUnauthorisedError errors will be converted into 401 responses
// accesscontrol.AccessForbiddenError errors will be converted into 403 responses
// any other errors will be converted into 500
func RequiresCatalogACL(request *events.APIGatewayProxyRequest, process func(claims accesscontrol.Claims, catalogACL catalogacl.AccessControlAdapter) (Response, error)) (Response, error) {
	token, err := readToken(request)
	if err != nil {
		return UnauthorizedResponse(err.Error())
	}

	claims, err := jwtDecoder.Decode(token)
	if err != nil {
		return UnauthorizedResponse(err.Error())
	}

	catalogACL := catalogacllogic.NewAccessControlAdapterFromToken(grantRepository, claims)

	response, err := process(claims, catalogACL)

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

// RequiresCatalogView is based on RequiresCatalogACL but creates a convenient Catalog View
func RequiresCatalogView(request *events.APIGatewayProxyRequest, process func(catalogView *catalogacl.View) (Response, error)) (Response, error) {
	return RequiresCatalogACL(request, func(claims accesscontrol.Claims, catalogACL catalogacl.AccessControlAdapter) (Response, error) {
		view := &catalogacl.View{
			UserEmail:     claims.Subject,
			AccessControl: catalogacllogic.NewAccessControlAdapterFromToken(grantRepository, claims),
		}

		return process(view)
	})
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
