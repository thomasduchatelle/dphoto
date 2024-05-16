package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type PutBodyDTO struct {
	Role string `json:"role"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	owner := request.PathParameters["owner"]
	folderName := request.PathParameters["folderName"]
	email := request.PathParameters["email"]

	return common.RequiresCatalogACL(&request, func(claims aclcore.Claims, rules catalogacl.CatalogRules) (common.Response, error) {
		albumId := catalog.NewAlbumIdFromStrings(owner, folderName)
		userId := usermodel.NewUserId(email)

		if err := rules.CanManageAlbum(albumId); err != nil {
			return common.Response{}, err
		}

		method := request.RequestContext.HTTP.Method
		switch method {
		case "PUT":
			body := new(PutBodyDTO)
			err := json.Unmarshal([]byte(request.Body), body)
			if err != nil {
				return common.BadRequest(err.Error())
			}

			scope, err := translateScope(body.Role)
			if err != nil {
				return common.BadRequest(err.Error())
			}

			err = common.GetShareAlbumCase().ShareAlbumWith(albumId, userId, scope)
			if errors.Is(err, catalog.AlbumNotFoundError) {
				return common.NotFound(fmt.Sprintf("%s/%s hasn't been found", owner, folderName))
			} else if err != nil {
				return common.InternalError(err)
			}

		case "DELETE":
			err := common.GetUnShareAlbumCase().StopSharingAlbum(albumId, userId)
			if err != nil {
				return common.InternalError(err)
			}

		default:
			return common.BadRequest(fmt.Sprintf("%s method is not supported", method))
		}

		return common.NoContent()
	})
}

func translateScope(role string) (aclcore.ScopeType, error) {
	switch role {
	case "visitor":
		return aclcore.AlbumVisitorScope, nil
	case "contributor":
		return aclcore.AlbumContributorScope, nil
	default:
		return "", errors.Errorf("'%s' role is not supported. Expected 'visitor' or 'contributor'", role)
	}
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
