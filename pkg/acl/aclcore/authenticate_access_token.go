package aclcore

import log "github.com/sirupsen/logrus"

// RefreshTokenAuthenticator use a known identity token issued by a known and trusted identity provider (google, facebook, ...) to create an access token
type RefreshTokenAuthenticator struct {
	AccessTokenGenerator   IAccessTokenGenerator
	RefreshTokenGenerator  IRefreshTokenGenerator
	RefreshTokenRepository RefreshTokenRepository
	IdentityDetailsStore   IdentityDetailsStore
}

type IAccessTokenGenerator interface {
	GenerateAccessToken(email string) (*Authentication, error)
}

type IRefreshTokenGenerator interface {
	GenerateRefreshToken(spec RefreshTokenSpec) (string, error)
}

func (s *RefreshTokenAuthenticator) AuthenticateFromAccessToken(refreshToken string) (*Authentication, *Identity, error) {
	spec, err := s.RefreshTokenRepository.FindRefreshToken(refreshToken)
	if err != nil {
		return nil, nil, err // can be an InvalidRefreshTokenError
	}

	if spec.AbsoluteExpiryTime.Before(TimeFunc()) {
		if deletedTokens, err := s.RefreshTokenRepository.HouseKeepRefreshToken(); err != nil {
			log.Infof("housekeeping - %d expired refresh token have been deleted", deletedTokens)
		}

		return nil, nil, ExpiredRefreshTokenError
	}

	identity, err := s.IdentityDetailsStore.FindIdentity(spec.Email)
	if err == IdentityDetailsNotFoundError {
		identity = &Identity{
			Email:   spec.Email,
			Name:    spec.Email,
			Picture: "",
		}
	} else if err != nil {
		return nil, nil, err
	}

	authentication, err := s.AccessTokenGenerator.GenerateAccessToken(spec.Email)
	if err != nil {
		return nil, nil, err
	}

	authentication.RefreshToken, err = s.RefreshTokenGenerator.GenerateRefreshToken(*spec)
	if err != nil {
		return nil, nil, err
	}

	err = s.RefreshTokenRepository.DeleteRefreshToken(refreshToken)
	return authentication, identity, err
}
