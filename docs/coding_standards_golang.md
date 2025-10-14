# Golang Coding Standards

## How to write a test ?

Use the golang idiomatic **table-driven** testing with a slice of test cases with for each:

1. **name**: descriptive names starting by "it should" and summarising the trigger and the expectation. Examples:
    * "it should GRANT access to the media owner"`
    * "it should DENY access to a visitor with no permission"`

2. **fields**: structure of the fields of the structure under test (ignore if a function is tested)
    * use fake in-memory implementation rather than mocks to stub the dependencies

3. **args**: structure of the function parameters

4. **want**: of the type of the first returned argument

5. **wantErr**: of the type `type ErrorAssertionFunc func(TestingT, error, ...interface{}) bool`
    * if no error are expected: `assert.NoError`
    * if an error is expected: an anonymous function (see the full example)

Write at the beginning of the file reusable fixture: each test should only declare what is specific to the test case.

### Testing example

```
func TestCatalogAuthorizer_IsAuthorisedToViewMedia(t *testing.T) {
    userId := usermodel.UserId("user-1")
    owner := ownermodel.Owner("owner-1")
    mediaId := catalog.MediaId("media-1")

    isAnAccessForbiddenError := func(t assert.TestingT, err error, i ...interface{}) bool {
        return assert.ErrorIs(t, err, aclcore.AccessForbiddenError)
    }

    type args struct {
        currentUser usermodel.CurrentUser
        owner       ownermodel.Owner
        mediaId     catalog.MediaId
    }
    tests := []struct {
        name    string
        args    args
        wantErr assert.ErrorAssertionFunc
    }{
        {
            name: "it should GRANT access to the media owner",
            args: args{
                currentUser: usermodel.CurrentUser{UserId: userId, Owner: &owner},
                owner:       owner,
                mediaId:     mediaId,
            },
            wantErr: assert.NoError,
        },
        {
            name: "it should DENY access to a visitor with no permission",
            args: args{
                currentUser: usermodel.CurrentUser{UserId: userId},
                owner:       owner,
                mediaId:     mediaId,
            },
            wantErr: isAnAccessForbiddenError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.wantErr(t, authorizer.IsAuthorisedToViewMedia(ctx, tt.args.currentUser, tt.args.owner, tt.args.mediaId))
        })
    }
}
```

## How to expose a REST API Endpoints ?

The REST API is hosted on AWS lambdas - one handler and function by operation - exposed through the AWS API Gateway (HTTP v2), and is secured by a custom lambda
authorizer.

The steps to expose a new endpoint are described below.

### 1. handler deployed as a lambda

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

### 2. authorisation logic

Add the authorisation rules in `api/lambdas/authorizer/main.go`. Example of configuration:

```
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

### 3. Deployment with CDK

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
