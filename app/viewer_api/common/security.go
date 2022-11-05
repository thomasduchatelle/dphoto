package common

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/accesscontrol"
	"github.com/thomasduchatelle/dphoto/domain/catalogacl"
	"strings"
)

// RequiresAuthorisation read the bearer token, parse it, and creates a Catalog View from it (enforcing catalog ACL).
//
// accesscontrol.AccessForbiddenError errors will be converted into 403 responses.
func RequiresAuthorisation(request *events.APIGatewayProxyRequest, process func(catalogView catalogacl.View) (Response, error), authorisationRules ...func(authContext catalogacl.View) error, ) (Response, error) {
	_, err := readToken(request)
	if err != nil {
		response, _ := NewJsonResponse(401, map[string]string{
			"error": fmt.Sprintf("access denied: %s", err.Error()),
		}, nil)
		return response, nil
	}

	catalogView := catalogacl.NewUserView("TODO") // TODO token should be decoded and used !
	if err != nil {
		response, _ := NewJsonResponse(401, map[string]string{
			"error": fmt.Sprintf("access denied: %s", err.Error()),
		}, nil)
		return response, nil
	}

	for i, rule := range authorisationRules {
		err = rule(catalogView)
		if err != nil {
			log.Warnf("Failed to authorise request at rule %d/%d %+v: %s", i+1, len(authorisationRules), request, err.Error())
			response, _ := NewJsonResponse(403, map[string]string{
				"error": fmt.Sprintf("access forbidden: %s", err.Error()),
			}, nil)
			return response, nil
		}
	}

	response, err := process(catalogView)
	if errors.Is(err, accesscontrol.AccessForbiddenError) {
		response, _ = NewJsonResponse(403, map[string]string{
			"error": fmt.Sprintf("access forbidden: %s", err.Error()),
		}, nil)
		return response, nil
	}
	return response, err
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
