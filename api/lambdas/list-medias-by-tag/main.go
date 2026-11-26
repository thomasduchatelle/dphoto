package main

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

type AlbumId struct {
	Owner      string `json:"owner"`
	FolderName string `json:"folderName"`
}

type Media struct {
	Id       string    `json:"id"`
	Type     string    `json:"type"`
	Filename string    `json:"filename"`
	Time     time.Time `json:"time"`
	Source   string    `json:"source"`
	AlbumId  AlbumId   `json:"albumId"`
	Tags     []string  `json:"tags"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	ctx := context.Background()

	// Extract user from authorizer context (already authenticated and authorized by Lambda Authorizer)
	_, err := common.GetCurrentUserFromContext(&request)
	if err != nil {
		return common.UnauthorizedResponse(err.Error())
	}

	// Note: owner-only permission check is already done by the Lambda Authorizer

	argParser := common.NewArgParser(&request)
	owner := ownermodel.Owner(argParser.ReadPathParameterString("owner"))
	rawTag := argParser.ReadPathParameterString("tag")
	if argParser.HasViolations() {
		return argParser.BadRequest()
	}

	tag, err := catalog.NewTag(rawTag)
	if err != nil {
		return common.BadRequest(map[string]string{"error": err.Error()})
	}

	log.Infof("list medias for tag %s/%s", owner, tag)

	medias, err := pkgfactory.TagServiceCase(ctx).ListMediasByTag(ctx, owner, tag)
	if err != nil {
		return common.InternalError(err)
	}

	resp := make([]Media, len(medias), len(medias))
	for i, media := range medias {
		resp[i] = Media{
			Id:       string(media.Id),
			Type:     string(media.Type),
			Filename: media.Filename,
			Time:     media.Details.DateTime,
			Source:   strings.Join([]string{media.Details.Make, media.Details.Model}, " "),
			AlbumId: AlbumId{
				Owner:      media.AlbumId.Owner.String(),
				FolderName: media.AlbumId.FolderName.String(),
			},
			Tags: tagStrings(media.Tags),
		}
	}

	return common.Ok(resp)
}

func tagStrings(tags []catalog.Tag) []string {
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		result = append(result, tag.String())
	}
	return result
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
