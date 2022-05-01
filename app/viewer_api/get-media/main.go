package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	"github.com/thomasduchatelle/dphoto/domain/catalogmodel"
	"github.com/thomasduchatelle/dphoto/domain/oauth"
	"strconv"
)

func Handler(request events.APIGatewayProxyRequest) (common.Response, error) {
	owner, _ := request.PathParameters["owner"]

	signature, err := parseSignature(&request)
	if err != nil {
		return common.BadRequest(map[string]string{
			"error": "path is invalid",
		})
	}

	width, err := parseWidth(&request)
	if err != nil {
		return common.BadRequest(map[string]string{
			"error": "width is invalid",
		})
	}

	if resp, deny := common.ValidateRequest(request, oauth.NewAuthoriseQuery("owner").WithOwner(owner, "READ")); deny {
		return resp, nil
	}

}

func parseWidth(request *events.APIGatewayProxyRequest) (width int, err error) {
	if widthString, ok := request.QueryStringParameters["w"]; ok {
		width, err = strconv.Atoi(widthString)
	}
	return
}

func parseSignature(request *events.APIGatewayProxyRequest) (signature catalogmodel.MediaSignature, err error) {
	signature.SignatureSha256, _ = request.PathParameters["signatureHash"]
	if signatureSize, ok := request.PathParameters["signatureSize"]; ok {
		signature.SignatureSize, err = strconv.Atoi(signatureSize)
	}
	return
}

func main() {
	common.Bootstrap()

	lambda.Start(Handler)
}
