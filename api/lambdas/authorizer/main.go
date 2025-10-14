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

	// Try cookie 'dphoto-access-token'
	if cookieHeader, ok := request.Headers["cookie"]; ok {
		cookies := strings.Split(cookieHeader, ";")
		for _, cookie := range cookies {
			cookie = strings.TrimSpace(cookie)
			if strings.HasPrefix(cookie, "dphoto-access-token=") {
				return strings.TrimPrefix(cookie, "dphoto-access-token="), nil
			}
		}
	}

	// Try query parameter 'access_token'
	if token, ok := request.QueryStringParameters["access_token"]; ok && token != "" {
		return token, nil
	}

	var authHeaders []string
	for _, header := range request.Headers {
		authHeaders = append(authHeaders, header)
	}
	return "", errors.Errorf("no authorization token found [%s]", strings.Join(authHeaders, ", "))
}

// Define all supported routes
var supportedRoutes = []Route{
	// Catalog endpoints - read-only
	{Pattern: "/api/v1/albums", Method: "GET"},                                            // list-albums
	{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/medias", Method: "GET"},         // list-medias
	
	// Catalog endpoints - mutations
	{Pattern: "/api/v1/albums", Method: "POST"},                                           // create-album
	{Pattern: "/api/v1/owners/{owner}/albums/{folderName}", Method: "DELETE"},             // delete-album
	{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/dates", Method: "PUT"},          // amend-album-dates
	{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/name", Method: "PUT"},           // amend-album-name
	
	// Archive endpoints
	{Pattern: "/api/v1/owners/{owner}/medias/{mediaId}/{filename}", Method: "GET"},        // get-media
	
	// Access control endpoints
	{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/shares/{email}", Method: "PUT"},    // share-album (grant)
	{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/shares/{email}", Method: "DELETE"}, // share-album (revoke)
	
	// User endpoints
	{Pattern: "/api/v1/owners", Method: "GET"},                                            // list-owners
	{Pattern: "/api/v1/users", Method: "GET"},                                             // list-users
}

func checkPermissions(ctx context.Context, user usermodel.CurrentUser, routeKey, rawPath string) error {
	authoriser := pkgfactory.AclCatalogAuthoriser(ctx)

	// Parse route key to extract method and template path
	parts := strings.Fields(routeKey)
	if len(parts) < 2 {
		return errors.Errorf("invalid route key: %s", routeKey)
	}

	method := parts[0]
	
	// Match the route using the new route matching system
	matched, err := MatchRoute(supportedRoutes, method, rawPath)
	if err != nil {
		return errors.Wrapf(err, "route %s %s", method, rawPath)
	}

	// Check permissions based on the matched route pattern
	switch matched.Route.Pattern {
	case "/api/v1/albums":
		if method == "GET" {
			// list-albums - no specific permission check needed (user is authenticated)
			return nil
		} else if method == "POST" {
			// create-album
			_, err := authoriser.CanCreateAlbum(ctx, user)
			return err
		}

	case "/api/v1/owners/{owner}/albums/{folderName}":
		// delete-album
		owner := matched.PathParams["owner"]
		folderName := matched.PathParams["folderName"]
		albumId := catalog.NewAlbumIdFromStrings(owner, folderName)
		return authoriser.CanDeleteAlbum(ctx, user, albumId)

	case "/api/v1/owners/{owner}/albums/{folderName}/dates":
		// amend-album-dates
		owner := matched.PathParams["owner"]
		folderName := matched.PathParams["folderName"]
		albumId := catalog.NewAlbumIdFromStrings(owner, folderName)
		return authoriser.CanAmendAlbumDates(ctx, user, albumId)

	case "/api/v1/owners/{owner}/albums/{folderName}/name":
		// amend-album-name
		owner := matched.PathParams["owner"]
		folderName := matched.PathParams["folderName"]
		albumId := catalog.NewAlbumIdFromStrings(owner, folderName)
		return authoriser.CanRenameAlbum(ctx, user, albumId)

	case "/api/v1/owners/{owner}/albums/{folderName}/medias":
		// list-medias
		owner := matched.PathParams["owner"]
		folderName := matched.PathParams["folderName"]
		albumId := catalog.NewAlbumIdFromStrings(owner, folderName)
		err := authoriser.IsAuthorisedToListMedias(ctx, user, albumId)
		if errors.Is(err, catalogacl.ErrAccessDenied) {
			return aclcore.AccessForbiddenError
		}
		return err

	case "/api/v1/owners/{owner}/medias/{mediaId}/{filename}":
		// get-media
		owner := matched.PathParams["owner"]
		mediaId := matched.PathParams["mediaId"]
		err := authoriser.IsAuthorisedToViewMedia(ctx, user, ownermodel.Owner(owner), catalog.MediaId(mediaId))
		if errors.Is(err, aclcore.AccessForbiddenError) {
			return err
		}
		return err

	case "/api/v1/owners":
		// list-owners - no specific permission check needed (user is authenticated)
		return nil

	case "/api/v1/users":
		// list-users - no specific permission check needed (user is authenticated)
		return nil

	case "/api/v1/owners/{owner}/albums/{folderName}/shares/{email}":
		// share-album (PUT = grant, DELETE = revoke)
		owner := matched.PathParams["owner"]
		folderName := matched.PathParams["folderName"]
		albumId := catalog.NewAlbumIdFromStrings(owner, folderName)
		return authoriser.CanShareAlbum(ctx, user, albumId)
	}

	return errors.Errorf("unknown route pattern: %s", matched.Route.Pattern)
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
