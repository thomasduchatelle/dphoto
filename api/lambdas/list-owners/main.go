package main

import (
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
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

	return common.RequiresAuthenticated(&request, func(_ usermodel.CurrentUser) (common.Response, error) {
		log.Infof("list identities of owners: %s", strings.Join(ownerIds, ", "))
		var dtos []OwnerDTO

		if len(ownerIds) == 0 || len(ownerIds) == 1 && ownerIds[0] == "" {
			return common.Ok(dtos)
		}

		var owners []ownermodel.Owner
		for _, id := range ownerIds {
			owners = append(owners, ownermodel.Owner(id))
		}

		identities, err := common.GetIdentityQueries().FindOwnerIdentities(owners)
		if err != nil {
			return common.InternalError(err)
		}

		for id, identities := range identities {
			names := make([]string, len(identities), len(identities))
			identityDTOs := make([]UserDetailsDTO, len(identities), len(identities))
			for i, identity := range identities {
				names[i] = identity.Name
				identityDTOs[i] = UserDetailsDTO{
					Name:    identity.Name,
					Email:   identity.Email.Value(),
					Picture: identity.Picture,
				}
			}

			dtos = append(dtos, OwnerDTO{
				ID:    id.Value(),
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
