package common

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// GetCurrentUserFromContext extracts the CurrentUser from the authorizer context
// This should be used when the Lambda Authorizer has already validated the request
func GetCurrentUserFromContext(request *events.APIGatewayV2HTTPRequest) (usermodel.CurrentUser, error) {
	if request.RequestContext.Authorizer == nil || request.RequestContext.Authorizer.Lambda == nil {
		return usermodel.CurrentUser{}, errors.New("no authorizer context found")
	}

	context := request.RequestContext.Authorizer.Lambda

	userIdStr, ok := context["userId"].(string)
	if !ok || userIdStr == "" {
		return usermodel.CurrentUser{}, errors.Errorf("userId not found in authorizer context: %+v", context)
	}

	user := usermodel.CurrentUser{
		UserId: usermodel.UserId(userIdStr),
	}

	// Owner is optional
	if ownerStr, ok := context["owner"].(string); ok && ownerStr != "" {
		owner := ownermodel.Owner(ownerStr)
		user.Owner = &owner
	}

	return user, nil
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
