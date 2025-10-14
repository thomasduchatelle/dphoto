# Golang Coding Standards

## Testing Standards

### Table-Driven Tests

Use idiomatic Go testing with a slice of test cases:

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

### Test Structure

1. **Test fixtures** at top (reusable test data)
2. **Type definitions**: `fields` for struct dependencies, `args` for function parameters
3. **Test cases**: Array with `name`, inputs, expected outputs
4. **Execution**: Loop with `t.Run()` for isolated subtests

**Test naming**: Use descriptive names starting with "it should":
- `"it should GRANT access to the media owner"`
- `"it should DENY access to a visitor with no permission"`

**Error assertions**: Use `assert.ErrorAssertionFunc` for flexible checking:

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

**In-Memory implementations**: Use in-memory ports, not mocks:

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

DPhoto uses Lambda authorizers to validate permissions before handlers execute.

**Architecture**:

```
Request → API Gateway → Lambda Authorizer → Lambda Handler → Response
                        (token + permissions) (uses context)
```

### Lambda Authorizer Pattern

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

### Route Configuration

Configure routes declaratively with authorization logic:

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

### Route Matching

Extract path parameters using regex patterns:

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

### Handler Pattern

Read user from context (already authorized):

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

### Adding an Authorized Endpoint

**1. Define route in authorizer** (`api/lambdas/authorizer/main.go`):

```go
{
    Route: Route{Pattern: "/api/v1/owners/{owner}/albums/{folderName}", Method: "DELETE"},
    Authorize: func(ctx context.Context, authoriser *catalogacl.CatalogAuthorizer, user usermodel.CurrentUser, pathParams map[string]string) error {
        albumId := catalog.NewAlbumIdFromStrings(pathParams["owner"], pathParams["folderName"])
        return authoriser.CanDeleteAlbum(ctx, user, albumId)
    },
}
```

**2. Create Lambda handler**:

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

**3. Wire in CDK** (`deployments/cdk/lib/.../endpoints-construct.ts`):

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

`CatalogAuthorizer` provides permission checks:

- `IsAuthorisedToListMedias(ctx, user, albumId)` - View album
- `IsAuthorisedToViewMedia(ctx, user, owner, mediaId)` - View media
- `CanCreateAlbum(ctx, user)` - Create album (returns owner)
- `CanDeleteAlbum(ctx, user, albumId)` - Delete album
- `CanShareAlbum(ctx, user, albumId)` - Share album
- `CanAmendAlbumDates(ctx, user, albumId)` - Update dates
- `CanRenameAlbum(ctx, user, albumId)` - Rename album

Returns `nil` for allowed, error for denied.
