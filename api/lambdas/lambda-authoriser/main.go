package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func Handler(request events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	ctx := context.Background()

	// Extract token from Authorization header
	token, err := extractToken(request)
	if err != nil {
		log.WithError(err).Warn("Failed to extract token")
		return denyResponse(), nil
	}

	// Decode and validate JWT token
	claims, err := common.AccessTokenDecoder().Decode(token)
	if err != nil {
		log.WithError(err).Warn("Failed to decode token")
		return denyResponse(), nil
	}

	user := claims.AsCurrentUser()

	// Check permissions based on the route
	err = checkPermissions(ctx, user, request.RouteKey, request.RawPath)
	if err != nil {
		log.WithError(err).Warnf("Permission denied for user %s on route %s", user.UserId, request.RouteKey)
		return denyResponse(), nil
	}

	// Return allow response with user context
	return allowResponse(claims), nil
}

func extractToken(request events.APIGatewayV2CustomAuthorizerV2Request) (string, error) {
	// Try Authorization header first
	if authHeader, ok := request.Headers["authorization"]; ok {
		bearerPrefix := "bearer "
		if strings.HasPrefix(strings.ToLower(authHeader), bearerPrefix) {
			return authHeader[len(bearerPrefix):], nil
		}
	}

	// Try query parameter for image loading
	if token, ok := request.QueryStringParameters["access_token"]; ok {
		return token, nil
	}

	var authHeaders []string
	for _, header := range request.Headers {
		authHeaders = append(authHeaders, header)
	}
	return "", errors.Errorf("no authorization token found [%s]", strings.Join(authHeaders, ", "))
}

func checkPermissions(ctx context.Context, user usermodel.CurrentUser, routeKey, rawPath string) error {
	authoriser := pkgfactory.AclCatalogAuthoriser(ctx)

	// Parse route and check permissions based on endpoint
	parts := strings.Fields(routeKey)
	if len(parts) < 2 {
		return errors.Errorf("invalid route key: %s", routeKey)
	}

	method := parts[0]
	path := parts[1]

	// Extract path parameters from rawPath
	pathParams := extractPathParameters(path, rawPath)

	switch {
	case method == "GET" && path == "/api/v1/albums":
		// list-albums - no specific permission check needed (user is authenticated)
		return nil

	case method == "POST" && path == "/api/v1/albums":
		// create-album
		_, err := authoriser.CanCreateAlbum(ctx, user)
		return err

	case method == "DELETE" && strings.HasPrefix(path, "/api/v1/owners/") && strings.Contains(path, "/albums/") && !strings.Contains(path, "/shares/"):
		// delete-album
		owner := pathParams["owner"]
		folderName := pathParams["folderName"]
		if owner == "" || folderName == "" {
			return errors.New("missing owner or folderName")
		}
		albumId := catalog.NewAlbumIdFromStrings(owner, folderName)
		return authoriser.CanDeleteAlbum(ctx, user, albumId)

	case method == "PUT" && strings.Contains(path, "/albums/") && strings.HasSuffix(path, "/dates"):
		// amend-album-dates
		owner := pathParams["owner"]
		folderName := pathParams["folderName"]
		if owner == "" || folderName == "" {
			return errors.New("missing owner or folderName")
		}
		albumId := catalog.NewAlbumIdFromStrings(owner, folderName)
		return authoriser.CanAmendAlbumDates(ctx, user, albumId)

	case method == "PUT" && strings.Contains(path, "/albums/") && strings.HasSuffix(path, "/name"):
		// amend-album-name
		owner := pathParams["owner"]
		folderName := pathParams["folderName"]
		if owner == "" || folderName == "" {
			return errors.New("missing owner or folderName")
		}
		albumId := catalog.NewAlbumIdFromStrings(owner, folderName)
		return authoriser.CanRenameAlbum(ctx, user, albumId)

	case method == "GET" && strings.Contains(path, "/albums/") && strings.HasSuffix(path, "/medias"):
		// list-medias
		owner := pathParams["owner"]
		folderName := pathParams["folderName"]
		if owner == "" || folderName == "" {
			return errors.New("missing owner or folderName")
		}
		albumId := catalog.NewAlbumIdFromStrings(owner, folderName)
		err := authoriser.IsAuthorisedToListMedias(ctx, user, albumId)
		if errors.Is(err, catalogacl.ErrAccessDenied) {
			return aclcore.AccessForbiddenError
		}
		return err

	case method == "GET" && strings.Contains(path, "/medias/"):
		// get-media
		owner := pathParams["owner"]
		mediaId := pathParams["mediaId"]
		if owner == "" || mediaId == "" {
			return errors.New("missing owner or mediaId")
		}
		err := authoriser.IsAuthorisedToViewMedia(ctx, user, ownermodel.Owner(owner), catalog.MediaId(mediaId))
		if errors.Is(err, aclcore.AccessForbiddenError) {
			return err
		}
		return err

	case method == "GET" && path == "/api/v1/owners":
		// list-owners - no specific permission check needed (user is authenticated)
		return nil

	case method == "GET" && path == "/api/v1/users":
		// list-users - no specific permission check needed (user is authenticated)
		return nil

	case (method == "PUT" || method == "DELETE") && strings.Contains(path, "/shares/"):
		// share-album
		owner := pathParams["owner"]
		folderName := pathParams["folderName"]
		if owner == "" || folderName == "" {
			return errors.New("missing owner or folderName")
		}
		albumId := catalog.NewAlbumIdFromStrings(owner, folderName)
		return authoriser.CanShareAlbum(ctx, user, albumId)

	default:
		return errors.Errorf("unknown route: %s %s", method, path)
	}
}

func extractPathParameters(template, actualPath string) map[string]string {
	params := make(map[string]string)

	// Split both paths into segments
	templateParts := strings.Split(strings.Trim(template, "/"), "/")
	actualParts := strings.Split(strings.Trim(actualPath, "/"), "/")

	if len(templateParts) != len(actualParts) {
		return params
	}

	for i, part := range templateParts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			paramName := strings.Trim(part, "{}")
			params[paramName] = actualParts[i]
		}
	}

	return params
}

func allowResponse(claims aclcore.Claims) events.APIGatewayV2CustomAuthorizerSimpleResponse {
	contextMap := map[string]interface{}{
		"userId": claims.Subject.Value(),
	}

	if claims.Owner != nil {
		contextMap["owner"] = claims.Owner.Value()
	}

	// Convert scopes to a JSON string
	scopes := make([]string, 0, len(claims.Scopes))
	for scope := range claims.Scopes {
		scopes = append(scopes, scope)
	}
	if len(scopes) > 0 {
		scopesJSON, _ := json.Marshal(scopes)
		contextMap["scopes"] = string(scopesJSON)
	}

	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
		Context:      contextMap,
	}
}

func denyResponse() events.APIGatewayV2CustomAuthorizerSimpleResponse {
	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: false,
	}
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
