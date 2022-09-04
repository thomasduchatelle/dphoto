package oauth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/access"
	"testing"
	"time"
)

func TestAuthorise(t *testing.T) {
	config := Config{
		TrustedIssuers:   nil,
		Issuer:           "unit-tests.dphoto.com",
		ValidityDuration: "10s",
		SecretJwtKey:     []byte("UnitSecret"),
	}

	tonyToken := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ1bml0LXRlc3RzLmRwaG90by5jb20iLCJzdWIiOiJ0b255QHN0YXJrLmNvbSIsImF1ZCI6WyJ1bml0LXRlc3RzLmRwaG90by5jb20iXSwiZXhwIjoxNjQxNDAzMzYzLCJpYXQiOjE2NDEzNzQ1NjMsImp0aSI6IjA5N2ZmOThlLTZlMDktMTFlYy1hZDQ3LTBlZmVkYTc2NDgwNyIsInNjb3BlcyI6ImFwaTphZG1pbiBvd25lcjpzZWxmOnRvbnlAc3RhcmsuY29tIn0.xiQ9QTTqsonINswIaed8kUYqiCaLcmi3yG4TNNKHPzaGHDYvov6OWdtsT0PwN9WD84BtJBqt4M_zGbbLNp5Kfw"

	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "it should grant access to a fully authorised consumer",
			args:    args{tonyToken},
			wantErr: assert.NoError,
		},
		{
			name: "it should deny access to a JWT with wrong signature",
			args: args{
				"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ1bml0LXRlc3RzLmRwaG90by5jb20iLCJzdWIiOiJ0b255QHN0YXJrLmNvbSIsImF1ZCI6WyJ1bml0LXRlc3RzLmRwaG90by5jb20iXSwiZXhwIjoxNjQxNDAzMzYzLCJpYXQiOjE2NDEzNzQ1NjMsImp0aSI6IjA5N2ZmOThlLTZlMDktMTFlYy1hZDQ3LTBlZmVkYTc2NDgwNyIsIm93bmVycyI6eyJpcm9ubWFuIjoiQURNSU4iLCJodWxrIjoiV1JJVEUifSwic2NvcGVzIjoiYWRtaW4gb3duZXIgY2VvIn0._oYEn9W_5uQrrNsv4phyKKr81Og3SWibjo-i6ubNp4i76-UHxCAlDILh-lgLdT2-PF7e7vV62kQOFPQe-N5Hjg",
			},
			wantErr: withErrorMessageContaining("invalid token"),
		},
		{
			name: "it should deny access to a JWT with wrong issuer",
			args: args{
				"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJoYWNrZXIuY29tIiwic3ViIjoidG9ueUBzdGFyay5jb20iLCJhdWQiOlsidW5pdC10ZXN0cy5kcGhvdG8uY29tIl0sImV4cCI6MTY0MTQwMzM2MywiaWF0IjoxNjQxMzc0NTYzLCJqdGkiOiIwOTdmZjk4ZS02ZTA5LTExZWMtYWQ0Ny0wZWZlZGE3NjQ4MDciLCJvd25lcnMiOnsiaXJvbm1hbiI6IkFETUlOIiwiaHVsayI6IldSSVRFIn0sInNjb3BlcyI6ImFkbWluIG93bmVyIGNlbyJ9.jrxM3dm9iv-6YZbYhzfSdaRpHmI371sAZxEIMh0MXYt-D3nhi7HZi6unWLmbRPpSGiUEzPjykA8Wjgi4_r_KeA",
			},
			wantErr: withErrorMessageContaining("issuer hacker.com not accepted"),
		},
		{
			name: "it should deny access to a JWT with wrong audience",
			args: args{
				"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ1bml0LXRlc3RzLmRwaG90by5jb20iLCJzdWIiOiJ0b255QHN0YXJrLmNvbSIsImF1ZCI6WyJ1bml0LXRlc3RzLmhhY2tlci5jb20iLCJoYWNrZXIuY29tIl0sImV4cCI6MTY0MTQwMzM2MywiaWF0IjoxNjQxMzc0NTYzLCJqdGkiOiIwOTdmZjk4ZS02ZTA5LTExZWMtYWQ0Ny0wZWZlZGE3NjQ4MDciLCJvd25lcnMiOnsiaXJvbm1hbiI6IkFETUlOIiwiaHVsayI6IldSSVRFIn0sInNjb3BlcyI6ImFkbWluIG93bmVyIGNlbyJ9.hCbgb00hKa2OEOnxGQzXu2cwAPW4QeDjZi4QGXHBDTpZgSLPtM3Zar__mKU1rCU-k7YhzhmtBztDWlXOnqBFnA",
			},
			wantErr: withErrorMessageContaining("not in the audience list"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewOverride(
				config,
				func() time.Time {
					return time.Unix(1641403360, 0)
				},
				func(email string, resourceType ...access.ResourceType) ([]*access.Resource, error) {
					return nil, errors.Errorf("Unexpected call to ListGrants")
				},
			).(*oauth)
			jwt.TimeFunc = adapter.now

			err := adapter.Authorise(tt.args.tokenString, func(claims Claims) error {
				return nil
			})

			tt.wantErr(t, err, tt.name)
		})
	}
}

func withErrorMessageContaining(contains string) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.Error(t, err, i) && assert.Contains(t, err.Error(), contains, i)
	}
}
