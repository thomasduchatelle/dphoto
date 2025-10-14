# Golang Coding Standards

## Testing Standards

Use idiomatic Go testing patterns with table-driven tests for comprehensive coverage.

### Table-Driven Tests

Structure tests using a slice of test cases with descriptive names:

```go
func TestCatalogAuthorizer_IsAuthorisedToViewMedia(t *testing.T) {
    userId1 := usermodel.UserId("user-1")
    owner1 := ownermodel.Owner("owner-1")
    mediaId1 := catalog.MediaId("media-1")
    albumId1 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("/folder-1")}
    currentUserAsOwner1 := usermodel.CurrentUser{UserId: userId1, Owner: &owner1}
    currentUserAsVisitor := usermodel.CurrentUser{UserId: userId1}

    isAnAccessForbiddenError := func(t assert.TestingT, err error, i ...interface{}) bool {
        return assert.ErrorIs(t, err, aclcore.AccessForbiddenError)
    }

    type fields struct {
        HasPermissionPort  HasPermissionPort
        CatalogQueriesPort CatalogQueriesPort
    }
    type args struct {
        ctx         context.Context
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
            name: "it should GRANT access to the media owner",
            fields: fields{
                HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{},
                CatalogQueriesPort: &catalog.MediaQueriesInMemory{
                    Medias: []catalog.InMemoryMedia{
                        catalog.NewInMemoryMedia(mediaId1, albumId1),
                    },
                },
            },
            args: args{
                ctx:         context.Background(),
                currentUser: currentUserAsOwner1,
                owner:       owner1,
                mediaId:     mediaId1,
            },
            wantErr: assert.NoError,
        },
        {
            name: "it should DENY access to a visitor with no permission",
            fields: fields{
                HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{},
                CatalogQueriesPort: &catalog.MediaQueriesInMemory{
                    Medias: []catalog.InMemoryMedia{
                        catalog.NewInMemoryMedia(mediaId1, albumId1),
                    },
                },
            },
            args: args{
                ctx:         context.Background(),
                currentUser: currentUserAsVisitor,
                owner:       owner1,
                mediaId:     mediaId1,
            },
            wantErr: isAnAccessForbiddenError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            f := &CatalogAuthorizer{
                HasPermissionPort:  tt.fields.HasPermissionPort,
                CatalogQueriesPort: tt.fields.CatalogQueriesPort,
            }
            tt.wantErr(t, f.IsAuthorisedToViewMedia(tt.args.ctx, tt.args.currentUser, tt.args.owner, tt.args.mediaId), fmt.Sprintf("IsAuthorisedToViewMedia(%v, %v, %v, %v)", tt.args.ctx, tt.args.currentUser, tt.args.owner, tt.args.mediaId))
        })
    }
}
```

### Test Naming

Use descriptive test names starting with "it should":

- `"it should GRANT access to the media owner"`
- `"it should DENY access to a visitor with no permission"`
- `"it should match a route which doesn't have any path parameters"`
- `"it should fail when no route matches the path"`

### Test Structure

1. **Setup**: Define test fixtures and helper functions at the top
2. **Type Definitions**: Use `fields` for struct dependencies, `args` for function parameters
3. **Test Cases**: Array of structs with `name`, input, and expected output
4. **Execution**: Loop through cases using `t.Run()` for isolated subtests

### Error Assertions

Use `assert.ErrorAssertionFunc` for flexible error checking:

```go
isAnAccessForbiddenError := func(t assert.TestingT, err error, i ...interface{}) bool {
    return assert.ErrorIs(t, err, aclcore.AccessForbiddenError)
}

tests := []struct {
    name    string
    wantErr assert.ErrorAssertionFunc
}{
    {
        name:    "it should succeed",
        wantErr: assert.NoError,
    },
    {
        name:    "it should fail with forbidden error",
        wantErr: isAnAccessForbiddenError,
    },
}
```

### In-Memory Test Implementations

Use in-memory implementations for ports, not mocks:

```go
fields: fields{
    HasPermissionPort: &aclcore.ScopeReadRepositoryInMemory{
        Scopes: []*aclcore.Scope{
            {Type: aclcore.MainOwnerScope, GrantedTo: userId1, ResourceOwner: owner1},
        },
    },
    CatalogQueriesPort: &catalog.MediaQueriesInMemory{
        Medias: []catalog.InMemoryMedia{
            catalog.NewInMemoryMedia(mediaId1, albumId1),
        },
    },
}
```

## API Endpoints with Authorization

DPhoto uses a Lambda authorizer pattern that validates permissions before handlers execute.

### Architecture Overview

```
Request → API Gateway → Lambda Authorizer → Lambda Handler → Response
                        (validates token)    (uses context)
                        (checks permissions)
```

### Lambda Authorizer

The authorizer validates JWT tokens and checks permissions:

```go
func Handler(request events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
    ctx := context.Background()

    // Extract token from Authorization header, cookie, or query parameter
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
```

### Route Authorization Configuration

Define routes with authorization logic in a declarative way:

```go
type AuthorizationFunc func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error

type AuthorizedRoute struct {
    Route
    Authorize AuthorizationFunc
}

var supportedRoutes = []AuthorizedRoute{
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
    {
        Route: Route{Pattern: "/api/v1/albums", Method: "POST"},
        Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
            // create-album
            _, err := authoriser.CanCreateAlbum(ctx, user)
            return err
        },
    },
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
}
```

### Route Matching with Path Parameters

Use regex-based pattern matching to extract path parameters:

```go
type Route struct {
    Pattern string // Pattern like "/api/v1/owners/{owner}/albums/{folderName}"
    Method  string // HTTP method like "GET", "POST", etc.
}

type MatchedRoute struct {
    Route      Route
    PathParams map[string]string
}

func MatchRoute(routes []Route, method, path string) (*MatchedRoute, error) {
    // Compile patterns and match against method and path
    // Returns matched route with extracted parameters
}
```

### Lambda Handler Pattern

Handlers read user information from the context passed by the authorizer:

```go
func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
    ctx := context.Background()
    owner := request.PathParameters["owner"]
    folderName := request.PathParameters["folderName"]

    albumId := catalog.NewAlbumIdFromStrings(owner, folderName)

    return common.RequiresAuthenticated(&request, func(user usermodel.CurrentUser) (common.Response, error) {
        log.Infof("list medias for album %s/%s", owner, folderName)
        
        medias, err := pkgfactory.CatalogMediaQueries(ctx).ListMedias(ctx, albumId)
        if err != nil {
            return common.InternalError(err)
        }

        resp := make([]Media, len(medias))
        for i, media := range medias {
            resp[i] = Media{
                Id:       string(media.Id),
                Type:     string(media.Type),
                Filename: media.Filename,
                Time:     media.Details.DateTime,
            }
        }

        return common.Ok(resp)
    })
}
```

### Error Handling

Use sentinel errors for authorization failures:

```go
var (
    ErrAccessDenied = errors.New("access denied")
)

// In authorizer
err := authoriser.IsAuthorisedToViewMedia(ctx, user, owner, mediaId)
if errors.Is(err, aclcore.AccessForbiddenError) {
    return err
}

// Common error handling
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
```

### Adding a New Authorized Endpoint

1. **Define the route in the authorizer**:

```go
{
    Route: Route{Pattern: "/api/v1/owners/{owner}/albums/{folderName}", Method: "DELETE"},
    Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
        albumId := catalog.NewAlbumIdFromStrings(pathParams["owner"], pathParams["folderName"])
        return authoriser.CanDeleteAlbum(ctx, user, albumId)
    },
}
```

2. **Create the Lambda handler**:

```go
func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
    ctx := context.Background()
    
    return common.RequiresAuthenticated(&request, func(user usermodel.CurrentUser) (common.Response, error) {
        // Authorization already checked by authorizer
        // Implement business logic
        return common.NoContent()
    })
}
```

3. **Wire it in the CDK**:

```typescript
const deleteAlbum = createSingleRouteEndpoint(this, 'DeleteAlbum', {
    environmentName: props.environmentName,
    httpApi: props.httpApi,
    functionName: 'delete-album',
    path: '/api/v1/owners/{owner}/albums/{folderName}',
    method: apigatewayv2.HttpMethod.DELETE,
    authorizer: props.authorizer,
});
props.catalogStore.grantReadWriteAccess(deleteAlbum.lambda);
```

### Authorization Methods

The `CatalogAuthorizer` provides methods for different permission checks:

- `IsAuthorisedToListMedias(ctx, user, albumId)` - View album contents
- `IsAuthorisedToViewMedia(ctx, user, owner, mediaId)` - View specific media
- `CanCreateAlbum(ctx, user)` - Create new album (returns owner)
- `CanDeleteAlbum(ctx, user, albumId)` - Delete album
- `CanShareAlbum(ctx, user, albumId)` - Share album with others
- `CanAmendAlbumDates(ctx, user, albumId)` - Update album dates
- `CanRenameAlbum(ctx, user, albumId)` - Rename album

Each method returns `nil` for allowed actions or an error for denied access.
