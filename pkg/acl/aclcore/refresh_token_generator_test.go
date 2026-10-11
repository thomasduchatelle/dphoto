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

	const email = usermodel.UserId("tony@stark.com")
	refreshDuration := map[aclcore.RefreshTokenPurpose]time.Duration{
		aclcore.RefreshTokenPurposeWeb: 1*time.Hour + 2*time.Minute,
	}

	type fields struct {
		RefreshTokenRepository *RefreshTokenRepositoryInMemory
		RefreshDuration        map[aclcore.RefreshTokenPurpose]time.Duration
	}
	type args struct {
		spec aclcore.RefreshTokenSpec
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantLen    int
		expectSpec aclcore.RefreshTokenSpec
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "it should generate a key with the expiry time derived from the purpose",
			fields: fields{
				RefreshTokenRepository: NewRefreshTokenRepositoryInMemory(),
				RefreshDuration:        refreshDuration,
			},
			args: args{
				spec: aclcore.RefreshTokenSpec{
					Email:               email,
					RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
					Scopes:              []string{"ironman"},
				},
			},
			wantLen: length,
			expectSpec: aclcore.RefreshTokenSpec{
				Email:               email,
				RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
				AbsoluteExpiryTime:  time.Date(2021, 12, 24, 1, 2, 0, 0, time.UTC),
				Scopes:              []string{"ironman"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should generate a token with default '1 hour' expiry time when the purpose is unknown",
			fields: fields{
				RefreshTokenRepository: NewRefreshTokenRepositoryInMemory(),
				RefreshDuration:        refreshDuration,
			},
			args: args{
				spec: aclcore.RefreshTokenSpec{
					Email: email,
				},
			},
			wantLen: length,
			expectSpec: aclcore.RefreshTokenSpec{
				Email:              email,
				AbsoluteExpiryTime: time.Date(2021, 12, 24, 1, 0, 0, 0, time.UTC),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should preserve an explicit absolute expiry time",
			fields: fields{
				RefreshTokenRepository: NewRefreshTokenRepositoryInMemory(),
				RefreshDuration:        refreshDuration,
			},
			args: args{
				spec: aclcore.RefreshTokenSpec{
					Email:               email,
					RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
					AbsoluteExpiryTime:  time.Date(2021, 12, 31, 23, 59, 59, 999, time.UTC),
				},
			},
			wantLen: length,
			expectSpec: aclcore.RefreshTokenSpec{
				Email:               email,
				RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
				AbsoluteExpiryTime:  time.Date(2021, 12, 31, 23, 59, 59, 999, time.UTC),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
			generator := &aclcore.RefreshTokenGenerator{
				RefreshTokenRepository: tt.fields.RefreshTokenRepository,
				RefreshDuration:        tt.fields.RefreshDuration,
			}

			got, err := generator.GenerateRefreshToken(tt.args.spec)
			if !tt.wantErr(t, err, fmt.Sprintf("GenerateRefreshToken(%v)", tt.args.spec)) {
				return
			}
			assert.Lenf(t, got, tt.wantLen, "GenerateRefreshToken(%v)", tt.args.spec)

			storedSpec, err := tt.fields.RefreshTokenRepository.FindRefreshToken(got)
			if assert.NoError(t, err, "FindRefreshToken(%v)", got) {
				assert.Equalf(t, tt.expectSpec, *storedSpec, "stored spec for token %v", got)
			}
		})
	}
}
