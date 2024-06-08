package common

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"strings"
)

func RequiresAuthenticated(request *events.APIGatewayV2HTTPRequest, process func(user usermodel.CurrentUser) (Response, error)) (Response, error) {
	token, err := readToken(request)
	if err != nil {
		return UnauthorizedResponse(err.Error())
	}

	claims, err := jwtDecoder.Decode(token)
	if err != nil {
		return UnauthorizedResponse(err.Error())
	}

	return process(claims.AsCurrentUser())
}

func HandleError(err error) (Response, error) {
	switch {
	case errors.Is(err, aclcore.AccessUnauthorisedError):
		return UnauthorizedResponse(err.Error())

	case errors.Is(err, aclcore.AccessForbiddenError):
		return ForbiddenResponse(err.Error())

	case err != nil:
		return InternalError(err)

	default:
		return Response{}, nil
	}
}

func readToken(request *events.APIGatewayV2HTTPRequest) (string, error) {
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
