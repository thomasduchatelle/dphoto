package common

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/ephoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/ephoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/ephoto/pkg/acl/catalogaclview"
	"strings"
)

// RequiresCatalogACL parses the token and instantiate the ACL extension for 'catalog' domain
//
// accesscontrol.AccessUnauthorisedError errors will be converted into 401 responses
// accesscontrol.AccessForbiddenError errors will be converted into 403 responses
// any other errors will be converted into 500
func RequiresCatalogACL(request *events.APIGatewayProxyRequest, process func(claims aclcore.Claims, rules catalogacl.CatalogRules) (Response, error)) (Response, error) {
	token, err := readToken(request)
	if err != nil {
		return UnauthorizedResponse(err.Error())
	}

	claims, err := jwtDecoder.Decode(token)
	if err != nil {
		return UnauthorizedResponse(err.Error())
	}

	catalogRules := catalogacl.OptimiseRulesWithAccessToken(catalogacl.NewCatalogRules(grantRepository, new(mediaAlbumResolver), claims.Subject), claims)

	response, err := process(claims, catalogRules)

	switch {
	case errors.Is(err, aclcore.AccessUnauthorisedError):
		return UnauthorizedResponse(err.Error())

	case errors.Is(err, aclcore.AccessForbiddenError):
		return ForbiddenResponse(err.Error())

	case err != nil:
		return InternalError(err)

	default:
		return response, nil
	}
}

// RequiresCatalogView is based on RequiresCatalogACL but creates a convenient Catalog View
func RequiresCatalogView(request *events.APIGatewayProxyRequest, process func(catalogView *catalogaclview.View) (Response, error)) (Response, error) {
	return RequiresCatalogACL(request, func(claims aclcore.Claims, rules catalogacl.CatalogRules) (Response, error) {
		view := &catalogaclview.View{
			UserEmail:      claims.Subject,
			CatalogRules:   rules,
			CatalogAdapter: new(catalogAdapter),
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
