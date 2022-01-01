package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/tencentyun/scf-go-lib/events"
	"github.com/thomasduchatelle/dphoto/domain/oauthmodel"
	"testing"
)

type AuthenticateFunc func(tokenString string) (oauthmodel.Authentication, oauthmodel.Identity, error)

func TestOauthRouting(t *testing.T) {
	a := assert.New(t)

	noopAuthenticate := func(tokenString string) (oauthmodel.Authentication, oauthmodel.Identity, error) {
		return oauthmodel.Authentication{}, oauthmodel.Identity{}, errors.Errorf("NOOP - call not expected")
	}

	tests := []struct {
		name         string
		request      events.APIGatewayRequest
		authenticate AuthenticateFunc
		wantCode     int
		wantBody     map[string]interface{}
	}{
		{
			"it should get the authorization header and remove the 'Bearer' prefix",
			events.APIGatewayRequest{
				Headers: map[string]string{
					"authorization": "Bearer qwertyuiop",
				},
			},
			func(tokenString string) (oauthmodel.Authentication, oauthmodel.Identity, error) {
				if tokenString == "qwertyuiop" {
					return oauthmodel.Authentication{
							AccessToken: "asdfghjkl",
							ExpiresIn:   42,
						},
						oauthmodel.Identity{
							Email:   "tony@stark.com",
							Name:    "Ironman",
							Picture: "https://stark.com/ceo",
						},
						nil
				}

				return oauthmodel.Authentication{}, oauthmodel.Identity{}, errors.Errorf("%s token is not expected.", tokenString)
			},
			200,
			map[string]interface{}{
				"access_token": "asdfghjkl",
				"expires_in":   float64(42),
				"identity": map[string]interface{}{
					"name":    "Ironman",
					"email":   "tony@stark.com",
					"picture": "https://stark.com/ceo",
				},
			},
		},

		{
			"it should reject with 401 when authorisation header is missing",
			events.APIGatewayRequest{
				Headers: nil,
			},
			noopAuthenticate,
			401,
			nil,
		},

		{
			"it should reject with 401 when authorisation header is not of BEARER type",
			events.APIGatewayRequest{
				Headers: map[string]string{
					"authorization": "something else qwerty",
				},
			},
			noopAuthenticate,
			400,
			nil,
		},

		{
			"it should reject with 403 when Authenticate method return an error",
			events.APIGatewayRequest{
				Headers: map[string]string{
					"authorization": "bearer qwerty",
				},
			},
			noopAuthenticate,
			403,
			map[string]interface{}{
				"error": "NOOP - call not expected",
			},
		},
	}

	for _, tt := range tests {
		oauthAuthenticate = tt.authenticate
		got, err := Handler(tt.request)

		if a.NoError(err) {
			if a.Equal(tt.wantCode, got.StatusCode, fmt.Sprintf("%s\n\nBody: %s", tt.name, got.Body)) {
				var gotBody map[string]interface{}
				err = json.Unmarshal([]byte(got.Body), &gotBody)

				if a.NoError(err, tt.name) && len(tt.wantBody) > 0 {
					a.Equal(tt.wantBody, gotBody, tt.name)
				}
			}
		}
	}
}
