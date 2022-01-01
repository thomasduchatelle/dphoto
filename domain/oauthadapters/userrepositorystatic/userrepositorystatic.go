package userrepositorystatic

import (
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/oauthmodel"
)

type StaticUserRepository struct{}

func (s *StaticUserRepository) FindUserRoles(email string) (*oauthmodel.UserRoles, error) {
	if email == "tomdush@gmail.com" {
		return &oauthmodel.UserRoles{
			IsUserManager: true,
			Owners: map[string]string{
				email: "ADMIN",
			},
		}, nil
	}

	return nil, errors.Errorf("%s is registered user", email)
}

func New() oauthmodel.UserRepository {
	return new(StaticUserRepository)
}
