package oauth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/oauthmodel"
	"testing"
	"time"
)

func TestAuthorise(t *testing.T) {
	a := assert.New(t)

	Now = func() time.Time {
		return time.Unix(1641403360, 0)
	}
	jwt.TimeFunc = Now
	Config = oauthmodel.Config{
		TrustedIssuers:   nil,
		Issuer:           "unit-tests.dphoto.com",
		ValidityDuration: "10s",
		SecretJwtKey:     []byte("UnitSecret"),
	}

	tonyToken := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ1bml0LXRlc3RzLmRwaG90by5jb20iLCJzdWIiOiJ0b255QHN0YXJrLmNvbSIsImF1ZCI6WyJ1bml0LXRlc3RzLmRwaG90by5jb20iXSwiZXhwIjoxNjQxNDAzMzYzLCJpYXQiOjE2NDEzNzQ1NjMsImp0aSI6IjA5N2ZmOThlLTZlMDktMTFlYy1hZDQ3LTBlZmVkYTc2NDgwNyIsIm93bmVycyI6eyJpcm9ubWFuIjoiQURNSU4iLCJodWxrIjoiV1JJVEUifSwic2NvcGVzIjoiYWRtaW4gb3duZXIgY2VvIn0.iM5MKMnd_kq1H_mH4BaEsK2Mmwb4f1BYQGgZBQboEAoiHX1ZZq4d6FlF2V5oUUxeYDdJcp4I17mERlvfuGSLXg"

	type args struct {
		tokenString string
		query       *AuthoriseQuery
	}
	tests := []struct {
		name    string
		args    args
		want    *oauthmodel.Claims
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"it should grant access to a fully authorised consumer",
			args{
				tonyToken,
				NewAuthoriseQuery("admin", "owner", "ceo").WithOwner("ironman", "ADMIN").WithOwner("hulk", "READ"),
			},
			nil,
			assert.NoError,
		},
		{
			"it should grant access to a valid JWT with default query",
			args{
				tonyToken,
				&AuthoriseQuery{},
			},
			nil,
			assert.NoError,
		},
		{
			"it should grant access to a valid JWT with nil query",
			args{
				tonyToken,
				nil,
			},
			nil,
			assert.NoError,
		},
		{
			"it should deny access to a valid JWT without the right scope",
			args{
				tonyToken,
				NewAuthoriseQuery("GOD"),
			},
			nil,
			withErrorMessageContaining("scopes too restrictive"),
		},
		{
			"it should deny access to a valid JWT without the owner permissions",
			args{
				tonyToken,
				NewAuthoriseQuery().WithOwner("hulk", "ADMIN"),
			},
			nil,
			withErrorMessageContaining("access to owners"),
		},
		{
			"it should deny access to a valid JWT without the owner listed at all",
			args{
				tonyToken,
				NewAuthoriseQuery().WithOwner("blackwidow", "READ"),
			},
			nil,
			withErrorMessageContaining("access to owners"),
		},
		{
			"it should deny access to a JWT with wrong signature",
			args{
				"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ1bml0LXRlc3RzLmRwaG90by5jb20iLCJzdWIiOiJ0b255QHN0YXJrLmNvbSIsImF1ZCI6WyJ1bml0LXRlc3RzLmRwaG90by5jb20iXSwiZXhwIjoxNjQxNDAzMzYzLCJpYXQiOjE2NDEzNzQ1NjMsImp0aSI6IjA5N2ZmOThlLTZlMDktMTFlYy1hZDQ3LTBlZmVkYTc2NDgwNyIsIm93bmVycyI6eyJpcm9ubWFuIjoiQURNSU4iLCJodWxrIjoiV1JJVEUifSwic2NvcGVzIjoiYWRtaW4gb3duZXIgY2VvIn0._oYEn9W_5uQrrNsv4phyKKr81Og3SWibjo-i6ubNp4i76-UHxCAlDILh-lgLdT2-PF7e7vV62kQOFPQe-N5Hjg",
				nil,
			},
			nil,
			withErrorMessageContaining("invalid token"),
		},
		{
			"it should deny access to a JWT with wrong issuer",
			args{
				"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJoYWNrZXIuY29tIiwic3ViIjoidG9ueUBzdGFyay5jb20iLCJhdWQiOlsidW5pdC10ZXN0cy5kcGhvdG8uY29tIl0sImV4cCI6MTY0MTQwMzM2MywiaWF0IjoxNjQxMzc0NTYzLCJqdGkiOiIwOTdmZjk4ZS02ZTA5LTExZWMtYWQ0Ny0wZWZlZGE3NjQ4MDciLCJvd25lcnMiOnsiaXJvbm1hbiI6IkFETUlOIiwiaHVsayI6IldSSVRFIn0sInNjb3BlcyI6ImFkbWluIG93bmVyIGNlbyJ9.jrxM3dm9iv-6YZbYhzfSdaRpHmI371sAZxEIMh0MXYt-D3nhi7HZi6unWLmbRPpSGiUEzPjykA8Wjgi4_r_KeA",
				nil,
			},
			nil,
			withErrorMessageContaining("issuer hacker.com not accepted"),
		},
		{
			"it should deny access to a JWT with wrong audience",
			args{
				"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ1bml0LXRlc3RzLmRwaG90by5jb20iLCJzdWIiOiJ0b255QHN0YXJrLmNvbSIsImF1ZCI6WyJ1bml0LXRlc3RzLmhhY2tlci5jb20iLCJoYWNrZXIuY29tIl0sImV4cCI6MTY0MTQwMzM2MywiaWF0IjoxNjQxMzc0NTYzLCJqdGkiOiIwOTdmZjk4ZS02ZTA5LTExZWMtYWQ0Ny0wZWZlZGE3NjQ4MDciLCJvd25lcnMiOnsiaXJvbm1hbiI6IkFETUlOIiwiaHVsayI6IldSSVRFIn0sInNjb3BlcyI6ImFkbWluIG93bmVyIGNlbyJ9.hCbgb00hKa2OEOnxGQzXu2cwAPW4QeDjZi4QGXHBDTpZgSLPtM3Zar__mKU1rCU-k7YhzhmtBztDWlXOnqBFnA",
				nil,
			},
			nil,
			withErrorMessageContaining("not in the audience list"),
		},
	}

	for _, tt := range tests {
		got, err := Authorise(tt.args.tokenString, tt.args.query)

		name := fmt.Sprintf("%s [Authorise(%+v, %+v)]", tt.name, tt.args.tokenString, tt.args.query)

		if !tt.wantErr(t, err, name) {
			return
		}

		if tt.want != nil {
			a.Equalf(tt.want, got, name)
		}
	}
}

func withErrorMessageContaining(contains string) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.Error(t, err, i) && assert.Contains(t, err.Error(), contains, i)
	}
}
