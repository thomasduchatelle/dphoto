package aclcore_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func TestRefreshTokenGenerator_GenerateRefreshToken(t1 *testing.T) {
	const length = 92
	aclcore.TimeFunc = func() time.Time {
		return time.Date(2021, 12, 24, 0, 0, 0, 0, time.UTC)
	}

	email := usermodel.UserId("tony@stark.com")
	refreshDuration := map[aclcore.RefreshTokenPurpose]time.Duration{
		aclcore.RefreshTokenPurposeWeb: 1*time.Hour + 2*time.Minute,
	}

	tests := []struct {
		name     string
		spec     aclcore.RefreshTokenSpec
		wantSpec aclcore.RefreshTokenSpec
		wantLen  int
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "it should generate a key with default expiry time from the purpose",
			spec: aclcore.RefreshTokenSpec{
				Email:               email,
				RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
				Scopes:              []string{"ironman"},
			},
			wantSpec: aclcore.RefreshTokenSpec{
				Email:               email,
				RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
				AbsoluteExpiryTime:  time.Date(2021, 12, 24, 1, 2, 0, 0, time.UTC),
				Scopes:              []string{"ironman"},
			},
			wantLen: length,
			wantErr: assert.NoError,
		},
		{
			name: "it should generate a token with default '1 hour' expiry time when no purpose is provided",
			spec: aclcore.RefreshTokenSpec{
				Email: email,
			},
			wantSpec: aclcore.RefreshTokenSpec{
				Email:              email,
				AbsoluteExpiryTime: time.Date(2021, 12, 24, 1, 0, 0, 0, time.UTC),
			},
			wantLen: length,
			wantErr: assert.NoError,
		},
		{
			name: "it should keep the specified absolute expiry time",
			spec: aclcore.RefreshTokenSpec{
				Email:               email,
				RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
				AbsoluteExpiryTime:  time.Date(2021, 12, 31, 23, 59, 59, 999, time.UTC),
			},
			wantSpec: aclcore.RefreshTokenSpec{
				Email:               email,
				RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
				AbsoluteExpiryTime:  time.Date(2021, 12, 31, 23, 59, 59, 999, time.UTC),
			},
			wantLen: length,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
			repository := NewRefreshTokenRepositoryInMemory()
			generator := &aclcore.RefreshTokenGenerator{
				RefreshTokenRepository: repository,
				RefreshDuration:        refreshDuration,
			}

			got, err := generator.GenerateRefreshToken(tt.spec)
			if !tt.wantErr(t, err, fmt.Sprintf("GenerateRefreshToken(%v)", tt.spec)) {
				return
			}
			assert.Lenf(t, got, tt.wantLen, "GenerateRefreshToken(%v)", tt.spec)

			gotSpec, findErr := repository.FindRefreshToken(got)
			if assert.NoError(t, findErr, "FindRefreshToken(%q)", got) {
				assert.Equal(t, tt.wantSpec, *gotSpec, "stored spec")
			}
		})
	}
}
