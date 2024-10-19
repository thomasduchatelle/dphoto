package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

const (
	responseMaxContent = 6 * 1024 * 1024
)

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	ctx := context.Background()

	parser := common.NewArgParser(&request)
	ownerValue := parser.ReadPathParameterString("owner")
	mediaIdValue := parser.ReadPathParameterString("mediaId")
	width := parser.ReadQueryParameterInt("w", false)

	if parser.HasViolations() {
		return parser.BadRequest()
	}

	owner := ownermodel.Owner(ownerValue)
	mediaId := catalog.MediaId(mediaIdValue)

	return common.RequiresAuthenticated(&request, func(user usermodel.CurrentUser) (common.Response, error) {
		err := pkgfactory.AclCatalogAuthoriser(ctx).IsAuthorisedToViewMedia(ctx, user, owner, mediaId)
		if errors.Is(err, aclcore.AccessForbiddenError) {
			return common.ForbiddenResponse(err.Error())
		}
		if err != nil {
			return common.InternalError(err)
		}

		if width == 0 {
			return redirectTo(archive.GetMediaOriginalURL(owner.Value(), mediaId.Value()))
		}

		content, contentType, err := archive.GetResizedImage(owner.Value(), mediaId.Value(), width, responseMaxContent)
		if errors.Is(err, archive.NotFoundError) {
			return common.NotFound(nil)
		}
		if errors.Is(err, archive.MediaOverflowError) {
			log.WithField("Owner", owner).Infof("Media %s/%s with width=%d is over max allowed payload. Redirecting.", owner, mediaId, width)
			return redirectTo(archive.GetResizedImageURL(owner.Value(), mediaId.Value(), width))
		}
		if err != nil {
			return common.InternalError(err)
		}

		based64Encoded := base64.StdEncoding.EncodeToString(content)
		log.WithField("Owner", owner).Infof("Media %s/%s with width=%d is served (%d KB ; base64 = %d KB)", owner, mediaId, width, len(content)/1024, len(based64Encoded)/1024)
		return common.Response{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type":  contentType,
				"Cache-Control": fmt.Sprintf("max-age=%d", 3600*24),
			},
			Body:            based64Encoded,
			IsBase64Encoded: true,
		}, nil
	})
}

func redirectTo(url string, err error) (common.Response, error) {
	if errors.Is(err, archive.NotFoundError) {
		return common.NotFound(nil)
	}
	if err != nil {
		log.WithError(err).Error("GetMediaOriginalURL failed with", err)
	}

	return common.Response{
		StatusCode: 307,
		Headers: map[string]string{
			"Location": url,
		},
	}, nil
}

func main() {
	common.BootstrapCatalogAndArchiveDomains()

	lambda.Start(Handler)
}
