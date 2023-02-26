package aclcore

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

type RefreshTokenGenerator struct {
	RefreshTokenRepository RefreshTokenRepository
	RefreshDuration        map[RefreshTokenPurpose]time.Duration
}

func (t *RefreshTokenGenerator) GenerateRefreshToken(spec RefreshTokenSpec) (string, error) {
	accessToken, err := CreateString(StringParams{
		Length:     92,
		Upper:      true,
		MinUpper:   3,
		Lower:      true,
		MinLower:   3,
		Numeric:    true,
		MinNumeric: 3,
		Special:    true,
		MinSpecial: 3,
	})
	if err != nil {
		return "", errors.Wrapf(err, "generating random access token")
	}

	if spec.AbsoluteExpiryTime.IsZero() {
		now := TimeFunc()
		lifeDuration := t.expiration(spec.RefreshTokenPurpose)
		spec.AbsoluteExpiryTime = now.Add(lifeDuration)
	}

	err = t.RefreshTokenRepository.StoreRefreshToken(string(accessToken), spec)
	return string(accessToken), err
}

func (t *RefreshTokenGenerator) expiration(purpose RefreshTokenPurpose) time.Duration {
	if duration, ok := t.RefreshDuration[purpose]; ok {
		return duration
	}

	log.Warnf("%s refresh token purpose is not defined, falling back on 1 hour.", purpose)
	return 1 * time.Hour
}
