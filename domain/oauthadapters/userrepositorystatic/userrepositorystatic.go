package userrepositorystatic

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/oauthmodel"
)

type StaticUserRepository struct{}

func (s *StaticUserRepository) FindUserRoles(email string) (*oauthmodel.UserRoles, error) {
	emailB64 := base64.StdEncoding.EncodeToString([]byte(email))
	if emailB64 == "dG9tZHVzaEBnbWFpbC5jb20=" || emailB64 == "Y2xhaXJlLm1hZ25pZXJAZ21haWwuY29t" {
		return &oauthmodel.UserRoles{
			IsUserManager: true,
			Owners: map[string]string{
				email: "ADMIN",
			},
		}, nil
	}

	return nil, errors.Errorf("%s is registered user", emailB64)
}

func New() oauthmodel.UserRepository {
	return new(StaticUserRepository)
}
