---
name: expose-rest-api
description: Expose a REST API endpoint implemented in go. Required to work on the APIs (`api/`).
---

# How to expose a REST API Endpoints

The REST API is hosted on AWS lambdas - one handler and function by operation - exposed through the AWS API Gateway (HTTP v2), and is secured by a custom lambda
authorizer.

The steps to expose a new endpoint are described below.

## 1. handler deployed as a lambda

```go
// api/lambdas/<function name>/main.go
package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
)

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	ctx := context.Background()

	// (optional) if authorisation is required, the CurrentUser must be present from the context from the Authorizer
	currentUser, err := common.GetCurrentUserFromContext(&request)
	if err != nil {
		return common.UnauthorizedResponse(err.Error())
	}

	// parse path parameters
	argParser := common.NewArgParser(&request)
	owner := ownermodel.Owner(argParser.ReadPathParameterString("owner"))
	folderName := catalog.NewFolderName(argParser.ReadPathParameterString("folderName"))

	// parse query parameters
	width := parser.ReadQueryParameterInt("w", false)

	// if any parameter is missing or malformed, returns a 400 error
	if parser.HasViolations() {
		return parser.BadRequest()
	}

	// ... do something by executing a function from a package in `pkg/`.

	return common.NoContent()
}
```

## 2. authorisation logic

Add the authorisation rules in `api/lambdas/authorizer/main.go`. Example of configuration:

```go
// list-medias
{
    Route: Route{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/medias", Method: "GET"},
    Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
        albumId := catalog.NewAlbumIdFromStrings(pathParams["owner"], pathParams["folderName"])
        err := authoriser.IsAuthorisedToListMedias(ctx, user, albumId)
        if errors.Is(err, catalogacl.ErrAccessDenied) {
            return aclcore.AccessForbiddenError
        }
        return err
    },
}
```

Each new operation should have its authorisation logic in a separate logic like `IsAuthorisedToListMedias`, and needs to be implemented in the package
`pkg/acl/catalogacl`.

Example of authorisation rule:

```go
package catalogacl

import (
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

func (a *CatalogAuthorizer) IsAuthorisedToListMedias(ctx context.Context, userId usermodel.CurrentUser, albumId catalog.AlbumId) error {
	if userId.Owner != nil && *userId.Owner == albumId.Owner {
		return nil
	}

	permissions, err := a.HasPermissionPort.FindScopesByIdCtx(ctx, aclcore.ScopeId{
		Type:          aclcore.AlbumVisitorScope,
		GrantedTo:     userId.UserId,
		ResourceOwner: albumId.Owner,
		ResourceId:    albumId.FolderName.String(),
	})
	if err != nil {
		return errors.Wrapf(err, "failed to check permissions for user %s on album %s", userId.UserId, albumId)
	}

	if len(permissions) > 0 {
		return nil
	}

	return errors.Wrapf(ErrAccessDenied, "user %s is not authorised to list medias from album %s", userId.UserId, albumId)
}
```

## 3. Deployment with CDK

Each subdomain of dphoto -- archive, catalog, and backup -- have their own CDK construct that provision their respective endpoints.

Example of CDK script `deployments/cdk/lib/archive/archive-endpoints-construct.ts`:

```typescript
const getMedia = createSingleRouteEndpoint(this, 'GetMedia', {
    environmentName: props.environmentName,
    functionName: 'get-media', // must match the name of the handler package "api/lambdas/<function name>/main.go"
    httpApi: props.httpApi,
    path: '/api/v1/owners/{owner}/medias/{mediaId}/{filename}', // AWS API Gateway route id 
    method: apigatewayv2.HttpMethod.GET,
    memorySize: 1024, // ignore to use a sensible default
    timeout: Duration.seconds(29), // ignore to use a sensible default (max is 29s)
    authorizer: props.queryParamAuthorizer,
});

// principle of least priviledge - only the access that the process requires is granted
props.catalogStore.grantReadAccess(getMedia.lambda);
props.archiveStore.grantReadAccessToRawAndCacheMedias(getMedia.lambda);
props.archivist.grantAccessToAsyncArchivist(getMedia.lambda);
```
