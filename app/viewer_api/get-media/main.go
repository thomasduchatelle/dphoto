package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	"github.com/thomasduchatelle/dphoto/domain/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/domain/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/domain/archive"
)

const (
	responseMaxContent = 6 * 1024 * 1024
)

func Handler(request events.APIGatewayProxyRequest) (common.Response, error) {
	parser := common.NewArgParser(&request)
	owner := parser.ReadPathParameterString("owner")
	mediaId := parser.ReadPathParameterString("mediaId")
	width := parser.ReadQueryParameterInt("w", false)

	if parser.HasViolations() {
		return parser.BadRequest()
	}

	return common.RequiresCatalogACL(&request, func(claims aclcore.Claims, rules catalogacl.CatalogRules) (common.Response, error) {
		if err := rules.CanReadMedia(owner, mediaId); err != nil {
			return common.Response{}, err
		}

		if width == 0 {
			return redirectTo(archive.GetMediaOriginalURL(owner, mediaId))
		}

		content, contentType, err := archive.GetResizedImage(owner, mediaId, width, responseMaxContent)
		if err == archive.NotFoundError {
			return common.NotFound(nil)
		}
		if err == archive.MediaOverflowError {
			log.WithField("Owner", owner).Infof("Media %s/%s with width=%d is over max allowed payload. Redirecting.", owner, mediaId, width)
			return redirectTo(archive.GetResizedImageURL(owner, mediaId, width))
		}
		if err != nil {
			return common.InternalError(err)
		}

		return common.Response{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type":  contentType,
				"Cache-Control": fmt.Sprintf("max-age=%d", 3600*24),
			},
			Body:            base64.StdEncoding.EncodeToString(content),
			IsBase64Encoded: true,
		}, nil
	})
}

func redirectTo(url string, err error) (common.Response, error) {
	if err == archive.NotFoundError {
		return common.NotFound(nil)
	}
	if err != nil {
		log.WithError(err).Error("GetMediaOriginalURL failed with", err)
	}

	return common.Response{
		StatusCode: 307,
		Headers: map[string]string{
			"FolderName": url,
		},
	}, nil
}

func main() {
	common.BootstrapCatalogAndArchiveDomains()

	lambda.Start(Handler)
}
