package bootstrap

import (
	"github.com/thomasduchatelle/dphoto/domain/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/domain/acl/acldynamodb"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
)

func CreateUserCase(cfg config.Config) *aclcore.CreateUser {
	repository := acldynamodb.Must(acldynamodb.New(cfg.GetAWSSession(), cfg.GetString(config.CatalogDynamodbTable), false))
	return &aclcore.CreateUser{
		ScopesReader: repository,
		ScopeWriter:  repository,
	}
}
