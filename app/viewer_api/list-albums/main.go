package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/oauth"
	"time"
)

type Album struct {
	Name       string    `json:"name"`
	Owner      string    `json:"owner"`
	FolderName string    `json:"folderName"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	TotalCount int       `json:"totalCount"`
}

func Handler(request events.APIGatewayProxyRequest) (common.Response, error) {
	owner, _ := request.PathParameters["owner"]
	if resp, deny := common.ValidateRequest(&request, oauth.NewAuthoriseQuery("owner").WithOwner(owner, "READ")); deny {
		return resp, nil
	}

	albums, err := catalog.FindAllAlbums(owner)
	if err != nil {
		return common.InternalError(errors.Wrapf(err, "failed to fetch albums"))
	}

	restAlbums := make([]Album, len(albums))
	for i, a := range albums {
		restAlbums[i] = Album{
			End:        a.End,
			FolderName: a.FolderName,
			Name:       a.Name,
			Owner:      owner,
			Start:      a.Start,
			TotalCount: a.TotalCount,
		}
	}
	return common.Ok(restAlbums)
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
