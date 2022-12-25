package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	common2 "github.com/thomasduchatelle/ephoto/api/lambdas/common"
	"github.com/thomasduchatelle/ephoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/ephoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/ephoto/pkg/archive"
)

const (
	responseMaxContent = 6 * 1024 * 1024
)

func Handler(request events.APIGatewayProxyRequest) (common2.Response, error) {
	parser := common2.NewArgParser(&request)
	owner := parser.ReadPathParameterString("owner")
	mediaId := parser.ReadPathParameterString("mediaId")
	width := parser.ReadQueryParameterInt("w", false)

	if parser.HasViolations() {
		return parser.BadRequest()
	}

	return common2.RequiresCatalogACL(&request, func(claims aclcore.Claims, rules catalogacl.CatalogRules) (common2.Response, error) {
		if err := rules.CanReadMedia(owner, mediaId); err != nil {
			return common2.Response{}, err
		}

		if width == 0 {
			return redirectTo(archive.GetMediaOriginalURL(owner, mediaId))
		}

		content, contentType, err := archive.GetResizedImage(owner, mediaId, width, responseMaxContent)
		if err == archive.NotFoundError {
			return common2.NotFound(nil)
		}
		if err == archive.MediaOverflowError {
			log.WithField("Owner", owner).Infof("Media %s/%s with width=%d is over max allowed payload. Redirecting.", owner, mediaId, width)
			return redirectTo(archive.GetResizedImageURL(owner, mediaId, width))
		}
		if err != nil {
			return common2.InternalError(err)
		}

		return common2.Response{
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

func redirectTo(url string, err error) (common2.Response, error) {
	if err == archive.NotFoundError {
		return common2.NotFound(nil)
	}
	if err != nil {
		log.WithError(err).Error("GetMediaOriginalURL failed with", err)
	}

	return common2.Response{
		StatusCode: 307,
		Headers: map[string]string{
			"Location": url,
		},
	}, nil
}

func main() {
	common2.BootstrapCatalogAndArchiveDomains()

	lambda.Start(Handler)
}
