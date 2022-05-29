package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	"github.com/thomasduchatelle/dphoto/domain/backup"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogmodel"
	"github.com/thomasduchatelle/dphoto/domain/oauth"
	"time"
)

func Handler(request events.APIGatewayProxyRequest) (common.Response, error) {
	parser := common.NewArgParser(&request)
	owner := parser.ReadPathParameterString("owner")
	signature, err := common.DecodeMediaId(parser.ReadPathParameterString("encodedId"))
	width := parser.ReadQueryParameterInteger("w", false)

	if parser.HasViolations() {
		return parser.BadRequest()
	}
	if err != nil {
		return common.BadRequest(map[string]string{
			"error": fmt.Sprintf("invalid signature: %s", err),
		})
	}

	if resp, deny := common.ValidateRequest(&request, oauth.NewAuthoriseQuery("owner").WithOwner(owner, "READ")); deny {
		return resp, nil
	}

	//if etag, ok := request.Headers["If-None-Match"]; ok && strings.HasPrefix(etag, etagPrefix) {
	//	return common.Response{
	//		StatusCode: 304, // unchanged
	//	}, nil
	//}

	locations, err := catalog.GetMediaLocations(owner, *signature)
	if err != nil {
		if errors.As(err, catalogmodel.MediaNotFoundError) {
			return common.NotFound(map[string]string{
				"error": "no locations found for this media",
			})
		}

		return common.InternalError(err)
	}

	if width == 0 {
		return RedirectResponse(owner, locations, 15*time.Minute) // TODO use the expiration of the JWT
	}

	return ResizedResponse(owner, locations, width)
}

func RedirectResponse(owner string, locations []*catalogmodel.MediaLocation, expires time.Duration) (common.Response, error) {
	url, err := backup.GetPreAuthorisedUrl(owner, locations, expires)
	if err != nil {
		return common.InternalError(err)
	}

	log.WithField("Location", locations[0].Filename).Infof("redirecting to S3 with pre-signed URL")
	return common.Response{
		StatusCode: 307,
		Headers: map[string]string{
			"Location": url,
		},
	}, nil
}

func ResizedResponse(owner string, locations []*catalogmodel.MediaLocation, width int) (common.Response, error) {
	content, contentType, err := backup.GetMediaContent(owner, locations, width)
	if err != nil {
		return common.InternalError(err)
	}

	log.WithField("Location", locations[0].Filename).Infof("Image size %dkb with width=%dpx", len(content)/1024, width)
	return common.Response{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":  contentType,
			"Cache-Control": fmt.Sprintf("max-age=%d", 3600*24),
		},
		Body:            base64.StdEncoding.EncodeToString(content),
		IsBase64Encoded: true,
	}, nil
}

func main() {
	common.Bootstrap()

	lambda.Start(Handler)
}
