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

// AuthorizationFunc defines the function signature for route authorization
type AuthorizationFunc func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error

// AuthorizedRoute represents a route with its authorization logic
type AuthorizedRoute struct {
	Route
	Authorize AuthorizationFunc
}

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

// Define all supported routes with their authorization logic
var supportedRoutes = []AuthorizedRoute{
	// Catalog endpoints - read-only
	{
		Route: Route{Pattern: "/api/v1/albums", Method: "GET"},
		Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
			// list-albums - no specific permission check needed (user is authenticated)
			return nil
		},
	},
	{
		Route: Route{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/medias", Method: "GET"},
		Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
			// list-medias
			albumId := catalog.NewAlbumIdFromStrings(pathParams["owner"], pathParams["folderName"])
			err := authoriser.IsAuthorisedToListMedias(ctx, user, albumId)
			if errors.Is(err, catalogacl.ErrAccessDenied) {
				return aclcore.AccessForbiddenError
			}
			return err
		},
	},

	// Catalog endpoints - mutations
	{
		Route: Route{Pattern: "/api/v1/albums", Method: "POST"},
		Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
			// create-album
			_, err := authoriser.CanCreateAlbum(ctx, user)
			return err
		},
	},
	{
		Route: Route{Pattern: "/api/v1/owners/{owner}/albums/{folderName}", Method: "DELETE"},
		Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
			// delete-album
			albumId := catalog.NewAlbumIdFromStrings(pathParams["owner"], pathParams["folderName"])
			return authoriser.CanDeleteAlbum(ctx, user, albumId)
		},
	},
	{
		Route: Route{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/dates", Method: "PUT"},
		Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
			// amend-album-dates
			albumId := catalog.NewAlbumIdFromStrings(pathParams["owner"], pathParams["folderName"])
			return authoriser.CanAmendAlbumDates(ctx, user, albumId)
		},
	},
	{
		Route: Route{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/name", Method: "PUT"},
		Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
			// amend-album-name
			albumId := catalog.NewAlbumIdFromStrings(pathParams["owner"], pathParams["folderName"])
			return authoriser.CanRenameAlbum(ctx, user, albumId)
		},
	},

	// Archive endpoints
	{
		Route: Route{Pattern: "/api/v1/owners/{owner}/medias/{mediaId}/{filename}", Method: "GET"},
		Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
			// get-media
			err := authoriser.IsAuthorisedToViewMedia(ctx, user, ownermodel.Owner(pathParams["owner"]), catalog.MediaId(pathParams["mediaId"]))
			if errors.Is(err, aclcore.AccessForbiddenError) {
				return err
			}
			return err
		},
	},

	// Access control endpoints
	{
		Route: Route{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/shares/{email}", Method: "PUT"},
		Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
			// share-album (grant)
			albumId := catalog.NewAlbumIdFromStrings(pathParams["owner"], pathParams["folderName"])
			return authoriser.CanShareAlbum(ctx, user, albumId)
		},
	},
	{
		Route: Route{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/shares/{email}", Method: "DELETE"},
		Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
			// share-album (revoke)
			albumId := catalog.NewAlbumIdFromStrings(pathParams["owner"], pathParams["folderName"])
			return authoriser.CanShareAlbum(ctx, user, albumId)
		},
	},

	// User endpoints
	{
		Route: Route{Pattern: "/api/v1/owners", Method: "GET"},
		Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
			// list-owners - no specific permission check needed (user is authenticated)
			return nil
		},
	},
	{
		Route: Route{Pattern: "/api/v1/users", Method: "GET"},
		Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
			// list-users - no specific permission check needed (user is authenticated)
			return nil
		},
	},
}

func checkPermissions(ctx context.Context, user usermodel.CurrentUser, routeKey, rawPath string) error {
	authoriser := pkgfactory.AclCatalogAuthoriser(ctx)

	// Parse route key to extract method and template path
	parts := strings.Fields(routeKey)
	if len(parts) < 2 {
		return errors.Errorf("invalid route key: %s", routeKey)
	}

	method := parts[0]

	// Extract route patterns from authorized routes
	routes := make([]Route, len(supportedRoutes))
	for i, ar := range supportedRoutes {
		routes[i] = ar.Route
	}

	// Match the route using the route matching system
	matched, err := MatchRoute(routes, method, rawPath)
	if err != nil {
		return errors.Wrapf(err, "route %s %s", method, rawPath)
	}

	// Find the corresponding authorized route and execute its authorization logic
	for _, ar := range supportedRoutes {
		if ar.Route.Pattern == matched.Route.Pattern && ar.Route.Method == matched.Route.Method {
			return ar.Authorize(ctx, authoriser, user, matched.PathParams)
		}
	}

	return errors.Errorf("no authorization logic found for route pattern: %s", matched.Route.Pattern)
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
