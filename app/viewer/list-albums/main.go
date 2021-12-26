package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/app/viewer/common"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
)

func Handler(ctx context.Context) (common.Response, error) {
	if err := common.ConnectCatalog("tomdush@gmail.com"); err != nil {
		return common.InternalError(err), nil
	}

	albums, err := catalog.FindAllAlbumsWithStats()
	if err != nil {
		common.InternalError(errors.Wrapf(err, "failed to fetch albums"))
	}

	return common.Ok(albums), nil
}

func main() {
	lambda.Start(Handler)
}
