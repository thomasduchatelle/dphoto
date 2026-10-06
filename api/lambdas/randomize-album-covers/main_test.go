package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

var (
	ownerStark    = ownermodel.Owner("stark")
	folderAvenger = "avengers"
)

type fakeRandomizer struct {
	covers []catalog.Cover
	err    error
}

func (f *fakeRandomizer) Randomize(_ context.Context, _ catalog.AlbumId) ([]catalog.Cover, error) {
	return f.covers, f.err
}

func ownerRequest() events.APIGatewayV2HTTPRequest {
	return events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{
			"owner":      ownerStark.Value(),
			"folderName": folderAvenger,
		},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			Authorizer: &events.APIGatewayV2HTTPRequestContextAuthorizerDescription{
				Lambda: map[string]interface{}{
					"userId": "tony@stark.com",
					"owner":  ownerStark.Value(),
				},
			},
		},
	}
}

func TestHandle(t *testing.T) {
	type args struct {
		request    events.APIGatewayV2HTTPRequest
		randomizer coverRandomizer
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantBody       string
	}{
		{
			name: "it should return 200 with the covers freshly computed when the owner triggers the refresh",
			args: args{
				request: ownerRequest(),
				randomizer: &fakeRandomizer{covers: []catalog.Cover{
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				}},
			},
			wantStatusCode: 200,
			wantBody:       `{"covers":[{"mediaId":"media-1","filename":"photo-1.jpg","origin":"CHERRY_PICKED"},{"mediaId":"media-2","filename":"photo-2.jpg","origin":"RANDOM"}]}`,
		},
		{
			name: "it should return 200 with an empty covers list when the album has no media to pick from",
			args: args{
				request:    ownerRequest(),
				randomizer: &fakeRandomizer{covers: nil},
			},
			wantStatusCode: 200,
			wantBody:       `{"covers":[]}`,
		},
		{
			name: "it should return 403 when the use case denies access to a visitor",
			args: args{
				request:    ownerRequest(),
				randomizer: &fakeRandomizer{err: aclcore.AccessForbiddenError},
			},
			wantStatusCode: 403,
		},
		{
			name: "it should return 404 when the album does not exist",
			args: args{
				request:    ownerRequest(),
				randomizer: &fakeRandomizer{err: errors.Wrap(catalog.AlbumNotFoundErr, "boom")},
			},
			wantStatusCode: 404,
		},
		{
			name: "it should return 500 when the use case fails with an unexpected error",
			args: args{
				request:    ownerRequest(),
				randomizer: &fakeRandomizer{err: errors.New("database on fire")},
			},
			wantStatusCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := handle(context.Background(), tt.args.request, tt.args.randomizer)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatusCode, response.StatusCode)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, response.Body)
			}
		})
	}
}
