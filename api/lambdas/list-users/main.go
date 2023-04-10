package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"strings"
)

type UserDetailsDTO struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	emails := strings.Split(request.QueryStringParameters["emails"], ",")

	return common.RequiresCatalogACL(&request, func(claims aclcore.Claims, rules catalogacl.CatalogRules) (common.Response, error) {
		log.Infof("list identities %s", strings.Join(emails, ", "))

		identities, err := common.GetIdentityQueries().FindIdentities(emails)
		if err != nil {
			return common.InternalError(err)
		}

		identitiesDTO := make([]UserDetailsDTO, len(identities), len(identities))
		for i, identity := range identities {
			identitiesDTO[i] = UserDetailsDTO{
				Name:    identity.Name,
				Email:   identity.Email,
				Picture: identity.Picture,
			}
		}
		return common.Ok(identitiesDTO)
	})
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
