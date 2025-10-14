API handlers in `api/lambdas` are each authorising the request using this code:

```golang
package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	// validate the access token (JWT) and extract the principal (user) - it would respond with 401 if the token is not found or invalid
	return common.RequiresAuthenticated(&request, func(user usermodel.CurrentUser) (common.Response, error) {
		albumId := catalog.AlbumId{Owner: ownermodel.Owner(owner), FolderName: catalog.NewFolderName(folderName)}

		// verify the permission the user has about this resource, here "CanAmendAlbumDates"
		if err := pkgfactory.AclCatalogAuthoriser(ctx).CanAmendAlbumDates(ctx, user, albumId); err != nil {
			// return with 403 if not authorised
			if errors.Is(err, aclcore.AccessForbiddenError) {
				return common.ForbiddenResponse(err.Error())
			}
			return common.InternalError(err)
		}

		// process the request if authorised and return
		// ...

		return common.NoContent()
	})
}
```

Your objective is to create a lambda authorizer for each of the following handler:

- amend-album-dates
- amend-album-name
- create-album
- delete-album
- get-media
- list-albums
- list-medias
- list-owners
- list-users
- share-album

You'll find the mapping to the API paths in `deployments/cdk/lib/access/user-endpoints-construct.ts`,
`deployments/cdk/lib/catalog/catalog-endpoints-construct.ts`, `deployments/cdk/lib/archive/archive-endpoints-construct.ts`.

Acceptance criteria:

* a new lambda is created and call `RequiresAuthenticated` followed by the appropriate method from `AclCatalogAuthoriser` based on the path
    * in case of failure, it will return either 401 or 403
    * a context with the properties in `usermodel.CurrentUser` is returned so it can be accessed by the handler
* the new lambda is deployed using CDK and used as authorizer for each handler
    * CDK coding standards: `docs/principles_cdk.md`
    * a test must validate that all API route have an authorizer attached, except those in the whitelist (the ones not listed above)
* the code compile and tests are passing
