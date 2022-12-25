package aclcore

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccessTokenDecoder_Decode(t *testing.T) {
	config := OAuthConfig{
		Issuer:           "https://dphoto.unit.test",
		ValidityDuration: 12 * time.Second,
		SecretJwtKey:     []byte("DPhotoJwtSecret"),
	}

	tests := []struct {
		name        string
		accessToken string
		want        Claims
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:        "it should accept a valid token and parse the Claims",
			accessToken: "eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwczovL2RwaG90by51bml0LnRlc3QiLCJzdWIiOiJ0b255QHN0YXJrLmNvbSIsImF1ZCI6WyJodHRwczovL2RwaG90by51bml0LnRlc3QiXSwiZXhwIjoxNjcxODQwMDA0LCJpYXQiOjE2NzE4NDAwMDAsInNjb3BlcyI6Im93bmVyOmlyb25tYW4gYXBpOmFkbWluIn0.num5Agz1j1m86QUJy27J8ON-nOUd-Myjah3TvzJGiWA",
			want: Claims{
				Subject: "tony@stark.com",
				Scopes: map[string]interface{}{
					"owner:ironman": nil,
					"api:admin":     nil,
				},
				Owner: "ironman",
			},
			wantErr: assert.NoError,
		},
		{
			name:        "it should reject a token with wrong signature",
			accessToken: "eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwczovL2RwaG90by51bml0LnRlc3QiLCJzdWIiOiJ0b255QHN0YXJrLmNvbSIsImF1ZCI6WyJodHRwczovL2RwaG90by51bml0LnRlc3QiXSwiZXhwIjoxNjcxODQwMDA0LCJpYXQiOjE2NzE4NDAwMDAsInNjb3BlcyI6Im93bmVyOmlyb25tYW4gYXBpOmFkbWluIn0.L6rWdGtpj3pGUNdna8kt_MX7ClXeLjSv90WkDCxmOZs",
			want:        Claims{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, AccessUnauthorisedError, i) &&
					assert.Contains(t, err.Error(), "signature is invalid")
			},
		},
		{
			name:        "it should reject a token with wrong issuer",
			accessToken: "eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwczovL2p1c3RpY2VsZWFndWUudW5pdC50ZXN0Iiwic3ViIjoidG9ueUBzdGFyay5jb20iLCJhdWQiOlsiaHR0cHM6Ly9kcGhvdG8udW5pdC50ZXN0Il0sImV4cCI6MTY3MTg0MDAwNCwiaWF0IjoxNjcxODQwMDAwLCJzY29wZXMiOiJvd25lcjppcm9ubWFuIGFwaTphZG1pbiJ9.S8MEodM1xc-waVF8okpjGcqduJEh3FjrBXN9SV_awAY",
			want:        Claims{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, AccessUnauthorisedError, i) &&
					assert.Contains(t, err.Error(), "issuer not accepted")
			},
		},
		{
			name:        "it should reject a token with wrong audience",
			accessToken: "eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwczovL2RwaG90by51bml0LnRlc3QiLCJzdWIiOiJ0b255QHN0YXJrLmNvbSIsImF1ZCI6WyJodHRwczovL2p1c3RpY2VsZWFndWUudW5pdC50ZXN0Il0sImV4cCI6MTY3MTg0MDAwNCwiaWF0IjoxNjcxODQwMDAwLCJzY29wZXMiOiJvd25lcjppcm9ubWFuIGFwaTphZG1pbiJ9.CgTvSgIz_-jG54Co2d-mW-3BnM-3jQ_XrJrgKkKVRec",
			want:        Claims{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, AccessUnauthorisedError, i) &&
					assert.Contains(t, err.Error(), "not in the audience list")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccessTokenDecoder{
				Config: config,
				Now: func() time.Time {
					return time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC)
				},
			}
			got, err := a.Decode(tt.accessToken)
			if !tt.wantErr(t, err, fmt.Sprintf("Decode(%v)", tt.accessToken)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Decode(%v)", tt.accessToken)
		})
	}
}
