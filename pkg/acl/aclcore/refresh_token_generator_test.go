package aclcore_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"testing"
	"time"
)

func TestRefreshTokenGenerator_GenerateRefreshToken(t1 *testing.T) {
	const length = 92
	aclcore.TimeFunc = func() time.Time {
		return time.Date(2021, 12, 24, 0, 0, 0, 0, time.UTC)
	}

	email := "tony@stark.com"
	refreshDuration := map[aclcore.RefreshTokenPurpose]time.Duration{
		aclcore.RefreshTokenPurposeWeb: 1*time.Hour + 2*time.Minute,
	}

	type fields struct {
		AccessTokenRepository func(t *testing.T) aclcore.RefreshTokenRepository
		RefreshDuration       map[aclcore.RefreshTokenPurpose]time.Duration
	}
	type args struct {
		spec aclcore.RefreshTokenSpec
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantLen int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should generate a key with specified expiry time",
			fields: fields{
				AccessTokenRepository: func(t *testing.T) aclcore.RefreshTokenRepository {
					repository := mocks.NewAccessTokenRepository(t)
					repository.On("StoreRefreshToken", mock.Anything, aclcore.RefreshTokenSpec{
						Email:               email,
						RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
						AbsoluteExpiryTime:  time.Date(2021, 12, 24, 1, 2, 0, 0, time.UTC),
						Scopes:              []string{"ironman"},
					}).Return(nil)

					return repository
				},
				RefreshDuration: refreshDuration,
			},
			args: args{
				spec: aclcore.RefreshTokenSpec{
					Email:               email,
					RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
					Scopes:              []string{"ironman"},
				},
			},
			wantLen: length,
			wantErr: assert.NoError,
		},
		{
			name: "it should generate a token with default '1 hour' expiry time",
			fields: fields{
				AccessTokenRepository: func(t *testing.T) aclcore.RefreshTokenRepository {
					repository := mocks.NewAccessTokenRepository(t)
					repository.On("StoreRefreshToken", mock.Anything, aclcore.RefreshTokenSpec{
						Email:              email,
						AbsoluteExpiryTime: time.Date(2021, 12, 24, 1, 0, 0, 0, time.UTC),
					}).Return(nil)

					return repository
				},
				RefreshDuration: refreshDuration,
			},
			args: args{
				spec: aclcore.RefreshTokenSpec{
					Email: email,
				},
			},
			wantLen: length,
			wantErr: assert.NoError,
		},
		{
			name: "it should generate a key with specified expiry time",
			fields: fields{
				AccessTokenRepository: func(t *testing.T) aclcore.RefreshTokenRepository {
					repository := mocks.NewAccessTokenRepository(t)
					repository.On("StoreRefreshToken", mock.Anything, aclcore.RefreshTokenSpec{
						Email:               email,
						RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
						AbsoluteExpiryTime:  time.Date(2021, 12, 31, 23, 59, 59, 999, time.UTC),
					}).Return(nil)

					return repository
				},
				RefreshDuration: refreshDuration,
			},
			args: args{
				spec: aclcore.RefreshTokenSpec{
					Email:               email,
					RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
					AbsoluteExpiryTime:  time.Date(2021, 12, 31, 23, 59, 59, 999, time.UTC),
				},
			},
			wantLen: length,
			wantErr: assert.NoError,
		},
		// it should generate a token with same absolute time
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
			generator := &aclcore.RefreshTokenGenerator{
				RefreshTokenRepository: tt.fields.AccessTokenRepository(t),
				RefreshDuration:        tt.fields.RefreshDuration,
			}

			got, err := generator.GenerateRefreshToken(tt.args.spec)
			if !tt.wantErr(t, err, fmt.Sprintf("GenerateRefreshToken(%v)", tt.args.spec)) {
				return
			}
			assert.Lenf(t, got, tt.wantLen, "GenerateRefreshToken(%v)", tt.args.spec)
		})
	}
}
