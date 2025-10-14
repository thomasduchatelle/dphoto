package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchRoute(t *testing.T) {
	tests := []struct {
		name           string
		routes         []Route
		method         string
		path           string
		wantPattern    string
		wantMethod     string
		wantPathParams map[string]string
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "it should match a route which doesn't have any path parameters",
			routes: []Route{
				{Pattern: "/api/v1/albums", Method: "GET"},
				{Pattern: "/api/v1/owners", Method: "GET"},
			},
			method:         "GET",
			path:           "/api/v1/albums",
			wantPattern:    "/api/v1/albums",
			wantMethod:     "GET",
			wantPathParams: map[string]string{},
			wantErr:        assert.NoError,
		},
		{
			name: "it should match a route with a single path parameter",
			routes: []Route{
				{Pattern: "/api/v1/owners/{owner}", Method: "GET"},
			},
			method:      "GET",
			path:        "/api/v1/owners/john-doe",
			wantPattern: "/api/v1/owners/{owner}",
			wantMethod:  "GET",
			wantPathParams: map[string]string{
				"owner": "john-doe",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should match a route with multiple path parameters",
			routes: []Route{
				{Pattern: "/api/v1/owners/{owner}/albums/{folderName}", Method: "DELETE"},
			},
			method:      "DELETE",
			path:        "/api/v1/owners/john-doe/albums/vacation-2023",
			wantPattern: "/api/v1/owners/{owner}/albums/{folderName}",
			wantMethod:  "DELETE",
			wantPathParams: map[string]string{
				"owner":      "john-doe",
				"folderName": "vacation-2023",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should match the correct method for the same path",
			routes: []Route{
				{Pattern: "/api/v1/albums", Method: "GET"},
				{Pattern: "/api/v1/albums", Method: "POST"},
			},
			method:         "POST",
			path:           "/api/v1/albums",
			wantPattern:    "/api/v1/albums",
			wantMethod:     "POST",
			wantPathParams: map[string]string{},
			wantErr:        assert.NoError,
		},
		{
			name: "it should fail when no route matches the path",
			routes: []Route{
				{Pattern: "/api/v1/albums", Method: "GET"},
			},
			method:  "GET",
			path:    "/api/v1/nonexistent",
			wantErr: assert.Error,
		},
		{
			name: "it should fail when the method does not match",
			routes: []Route{
				{Pattern: "/api/v1/albums", Method: "GET"},
			},
			method:  "POST",
			path:    "/api/v1/albums",
			wantErr: assert.Error,
		},
		{
			name: "it should match a route with nested parameters",
			routes: []Route{
				{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/medias", Method: "GET"},
				{Pattern: "/api/v1/owners/{owner}/medias/{mediaId}/{filename}", Method: "GET"},
			},
			method:      "GET",
			path:        "/api/v1/owners/alice/albums/summer-2024/medias",
			wantPattern: "/api/v1/owners/{owner}/albums/{folderName}/medias",
			wantMethod:  "GET",
			wantPathParams: map[string]string{
				"owner":      "alice",
				"folderName": "summer-2024",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should match a route with multiple parameters and file extension",
			routes: []Route{
				{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/medias", Method: "GET"},
				{Pattern: "/api/v1/owners/{owner}/medias/{mediaId}/{filename}", Method: "GET"},
			},
			method:      "GET",
			path:        "/api/v1/owners/bob/medias/media123/photo.jpg",
			wantPattern: "/api/v1/owners/{owner}/medias/{mediaId}/{filename}",
			wantMethod:  "GET",
			wantPathParams: map[string]string{
				"owner":    "bob",
				"mediaId":  "media123",
				"filename": "photo.jpg",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should match a route with dates path suffix",
			routes: []Route{
				{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/dates", Method: "PUT"},
			},
			method:      "PUT",
			path:        "/api/v1/owners/test-owner/albums/my-album/dates",
			wantPattern: "/api/v1/owners/{owner}/albums/{folderName}/dates",
			wantMethod:  "PUT",
			wantPathParams: map[string]string{
				"owner":      "test-owner",
				"folderName": "my-album",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should match a route with name path suffix",
			routes: []Route{
				{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/name", Method: "PUT"},
			},
			method:      "PUT",
			path:        "/api/v1/owners/test-owner/albums/my-album/name",
			wantPattern: "/api/v1/owners/{owner}/albums/{folderName}/name",
			wantMethod:  "PUT",
			wantPathParams: map[string]string{
				"owner":      "test-owner",
				"folderName": "my-album",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should match shares route with PUT method",
			routes: []Route{
				{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/shares/{email}", Method: "PUT"},
				{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/shares/{email}", Method: "DELETE"},
			},
			method:      "PUT",
			path:        "/api/v1/owners/owner1/albums/album1/shares/user@example.com",
			wantPattern: "/api/v1/owners/{owner}/albums/{folderName}/shares/{email}",
			wantMethod:  "PUT",
			wantPathParams: map[string]string{
				"owner":      "owner1",
				"folderName": "album1",
				"email":      "user@example.com",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should match shares route with DELETE method",
			routes: []Route{
				{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/shares/{email}", Method: "PUT"},
				{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/shares/{email}", Method: "DELETE"},
			},
			method:      "DELETE",
			path:        "/api/v1/owners/owner1/albums/album1/shares/user@example.com",
			wantPattern: "/api/v1/owners/{owner}/albums/{folderName}/shares/{email}",
			wantMethod:  "DELETE",
			wantPathParams: map[string]string{
				"owner":      "owner1",
				"folderName": "album1",
				"email":      "user@example.com",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should handle special characters in path parameters",
			routes: []Route{
				{Pattern: "/api/v1/owners/{owner}/albums/{folderName}", Method: "GET"},
			},
			method:      "GET",
			path:        "/api/v1/owners/john-doe-123/albums/vacation-2023-summer",
			wantPattern: "/api/v1/owners/{owner}/albums/{folderName}",
			wantMethod:  "GET",
			wantPathParams: map[string]string{
				"owner":      "john-doe-123",
				"folderName": "vacation-2023-summer",
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched, err := MatchRoute(tt.routes, tt.method, tt.path)

			if !tt.wantErr(t, err, fmt.Sprintf("MatchRoute(%v, %s, %s)", tt.routes, tt.method, tt.path)) {
				return
			}

			if err == nil {
				assert.Equal(t, tt.wantPattern, matched.Route.Pattern, "Route pattern should match")
				assert.Equal(t, tt.wantMethod, matched.Route.Method, "Route method should match")
				assert.Equal(t, tt.wantPathParams, matched.PathParams, "Path parameters should match")
			}
		})
	}
}
