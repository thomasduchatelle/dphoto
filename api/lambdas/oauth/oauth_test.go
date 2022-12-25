package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tencentyun/scf-go-lib/events"
	"github.com/thomasduchatelle/ephoto/mocks"
	"github.com/thomasduchatelle/ephoto/pkg/acl/aclcore"
	"testing"
)

func TestOauthRouting(t *testing.T) {
	noopAuthenticate := func(t *testing.T, oauthMock *mocks.SSOAuthenticator) {
		oauthMock.On("AuthenticateFromExternalIDProvider", mock.Anything).Maybe().Return(nil, nil, aclcore.NotPreregisteredError)
	}

	tests := []struct {
		name      string
		request   events.APIGatewayRequest
		initMocks func(t *testing.T, oauthMock *mocks.SSOAuthenticator)
		wantCode  int
		wantBody  map[string]interface{}
	}{
		{
			name: "it should get the authorization header and remove the 'Bearer' prefix",
			request: events.APIGatewayRequest{
				Headers: map[string]string{
					"authorization": "Bearer qwertyuiop",
				},
			},
			initMocks: func(t *testing.T, oauthMock *mocks.SSOAuthenticator) {
				oauthMock.On("AuthenticateFromExternalIDProvider", "qwertyuiop").Once().Return(
					&aclcore.Authentication{
						AccessToken: "asdfghjkl",
						ExpiresIn:   42,
					},
					&aclcore.Identity{
						Email:   "tony@stark.com",
						Name:    "Ironman",
						Picture: "https://stark.com/ceo",
					},
					nil,
				)
			},
			wantCode: 200,
			wantBody: map[string]interface{}{
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
			name: "it should reject with 401 when authorisation header is missing",
			request: events.APIGatewayRequest{
				Headers: nil,
			},
			initMocks: noopAuthenticate,
			wantCode:  401,
		},
		{
			name: "it should reject with 401 when authorisation header is not of BEARER type",
			request: events.APIGatewayRequest{
				Headers: map[string]string{
					"authorization": "something else qwerty",
				},
			},
			initMocks: noopAuthenticate,
			wantCode:  400,
		},
		{
			name: "it should reject with 403 when Authenticate method return an error",
			request: events.APIGatewayRequest{
				Headers: map[string]string{
					"authorization": "bearer qwerty",
				},
			},
			initMocks: noopAuthenticate,
			wantCode:  403,
			wantBody: map[string]interface{}{
				"code":  "oauth.user-not-preregistered",
				"error": "user must be pre-registered",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			authenticator = mocks.NewSSOAuthenticator(t)
			tt.initMocks(t, authenticator.(*mocks.SSOAuthenticator))

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

		})
	}
}
