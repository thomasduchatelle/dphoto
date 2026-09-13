---
name: go
description: Golang coding standards for DPhoto including table-driven testing patterns, REST API endpoints, and deployment as lambda with CDK. Required skill to work on the backend (`pkg/`), APIs (`api/`), and CLI (`cmd/`).
---

# Golang Coding Standards for DPhoto

## How to write a test

Use the golang idiomatic **table-driven** testing with a slice of test cases. The **canonical
structure of each case is fixed** — do not omit or rename fields on a whim:

1. **`name`** — mandatory, always. Descriptive names starting with "it should" and summarising the
   trigger and the expectation. Examples:
    * `"it should GRANT access to the media owner"`
    * `"it should DENY access to a visitor with no permission"`

2. **`fields`** — the actual fields of the structure under test. Mirror the struct literally.
   **Ignore only if the code under test is a plain function, not a struct.**
    * Do NOT drop `fields` because the tests happen to share the same collaborators — the fields
      must still be declared so each case can override any one of them when it needs to.
    * Use concrete Fake types (not the interface) when `expectX` below needs to read state back
      through a test-only accessor.

3. **`args`** — the arguments passed to the function under test. **Mandatory even if there is only
   a single argument.** Ignore only if the function takes no argument at all.

4. **`want`** — the first returned value (other than `error`). Ignore if the function returns only
   an error, or if more than one non-error value is returned (use `wantX` instead). Asserted with
   deep equals.

5. **`wantX`** — use when the function returns multiple non-error values; replace `X` with the
   parameter name (e.g. `wantAlbum`, `wantMedias`). Asserted with deep equals.

6. **`wantErr`** — of type `assert.ErrorAssertionFunc`
   (`func(TestingT, error, ...interface{}) bool`). Ignore if the function does not return an
   error.
    * If no error is expected: `assert.NoError`.
    * If a specific error is expected: an anonymous function using `assert.ErrorIs` /
      `assert.ErrorAs` (see the full example).

7. **`expectX`** — use only when a Fake dependency must be asserted after the call, and only
   when the assertion is genuinely meaningful for the test. Examples: `expectEventsFired []Event`,
   `expectStoredAlbum *Album`.
    * **Never** use `expectX` as a leaked "was-called" assertion (e.g. `expectCalled int`,
      `expectMockCalls int`). If you find yourself counting calls, redesign the test around state.

### Fake reuse across cases

- Instantiate **one "happy path" version of each Fake once, before the table**, and let every
  case reuse it by default. This keeps each case focused on what is specific to it.
- Override a Fake in a single case **only** when that case needs to exercise a specific branch
  (empty state, seeded record, error injection, …). In that case, build the Fake inline inside
  the case's `fields` (or via a small helper), leaving all other cases on the shared instance.
- Fixtures shared across cases (IDs, sample records, canned certificates, …) also live at the
  top of the test file. Each case declares only what is specific to it.

### Testing example

```go
func TestCatalogAuthorizer_IsAuthorisedToViewMedia(t *testing.T) {
    userId := usermodel.UserId("user-1")
    owner := ownermodel.Owner("owner-1")
    mediaId := catalog.MediaId("media-1")

    isAnAccessForbiddenError := func(t assert.TestingT, err error, i ...interface{}) bool {
        return assert.ErrorIs(t, err, aclcore.AccessForbiddenError)
    }

    // Shared happy-path Fake, reused by every case unless a case overrides it in `fields`.
    scopes := &ScopeRepositoryInMemory{}

    type fields struct {
        HasPermissionPort catalogacl.HasPermissionPort
    }
    type args struct {
        currentUser usermodel.CurrentUser
        owner       ownermodel.Owner
        mediaId     catalog.MediaId
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        wantErr assert.ErrorAssertionFunc
    }{
        {
            name:   "it should GRANT access to the media owner",
            fields: fields{HasPermissionPort: scopes},
            args: args{
                currentUser: usermodel.CurrentUser{UserId: userId, Owner: &owner},
                owner:       owner,
                mediaId:     mediaId,
            },
            wantErr: assert.NoError,
        },
        {
            name:   "it should DENY access to a visitor with no permission",
            fields: fields{HasPermissionPort: scopes},
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
            authorizer := &catalogacl.CatalogAuthorizer{HasPermissionPort: tt.fields.HasPermissionPort}
            tt.wantErr(t, authorizer.IsAuthorisedToViewMedia(ctx, tt.args.currentUser, tt.args.owner, tt.args.mediaId))
        })
    }
}
```

## How to expose a REST API Endpoints

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
