package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogaclview"
	"strings"
)

type OwnerDTO struct {
	ID    string           `json:"id"`
	Name  string           `json:"name"`
	Users []UserDetailsDTO `json:"users"`
}

type UserDetailsDTO struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	ownerIds := strings.Split(request.QueryStringParameters["ids"], ",")

	return common.RequiresCatalogView(&request, func(catalogView *catalogaclview.View) (common.Response, error) {
		log.Infof("list identities of owners: %s", strings.Join(ownerIds, ", "))
		var dtos []OwnerDTO

		if len(ownerIds) == 0 || len(ownerIds) == 1 && ownerIds[0] == "" {
			return common.Ok(dtos)
		}

		owners, err := common.GetIdentityQueries().FindOwnerIdentities(ownerIds)
		if err != nil {
			return common.InternalError(err)
		}

		for id, identities := range owners {
			names := make([]string, len(identities), len(identities))
			identityDTOs := make([]UserDetailsDTO, len(identities), len(identities))
			for i, identity := range identities {
				names[i] = identity.Name
				identityDTOs[i] = UserDetailsDTO{
					Name:    identity.Name,
					Email:   identity.Email,
					Picture: identity.Picture,
				}
			}

			dtos = append(dtos, OwnerDTO{
				ID:    id,
				Name:  strings.Join(names, ", "),
				Users: identityDTOs,
			})
		}

		return common.Ok(dtos)
	})
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
