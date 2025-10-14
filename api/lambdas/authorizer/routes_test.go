package main

import (
	"testing"
)

func TestMatchRoute_ShouldMatchRouteWithoutPathParameters(t *testing.T) {
	routes := []Route{
		{Pattern: "/api/v1/albums", Method: "GET"},
		{Pattern: "/api/v1/owners", Method: "GET"},
	}

	matched, err := MatchRoute(routes, "GET", "/api/v1/albums")
	if err != nil {
		t.Fatalf("Expected match, got error: %v", err)
	}

	if matched.Route.Pattern != "/api/v1/albums" {
		t.Errorf("Expected pattern /api/v1/albums, got %s", matched.Route.Pattern)
	}

	if matched.Route.Method != "GET" {
		t.Errorf("Expected method GET, got %s", matched.Route.Method)
	}

	if len(matched.PathParams) != 0 {
		t.Errorf("Expected no path params, got %v", matched.PathParams)
	}
}

func TestMatchRoute_ShouldMatchRouteWithSinglePathParameter(t *testing.T) {
	routes := []Route{
		{Pattern: "/api/v1/owners/{owner}", Method: "GET"},
	}

	matched, err := MatchRoute(routes, "GET", "/api/v1/owners/john-doe")
	if err != nil {
		t.Fatalf("Expected match, got error: %v", err)
	}

	if matched.Route.Pattern != "/api/v1/owners/{owner}" {
		t.Errorf("Expected pattern /api/v1/owners/{owner}, got %s", matched.Route.Pattern)
	}

	if matched.PathParams["owner"] != "john-doe" {
		t.Errorf("Expected owner=john-doe, got %s", matched.PathParams["owner"])
	}

	if len(matched.PathParams) != 1 {
		t.Errorf("Expected 1 path param, got %d", len(matched.PathParams))
	}
}

func TestMatchRoute_ShouldMatchRouteWithMultiplePathParameters(t *testing.T) {
	routes := []Route{
		{Pattern: "/api/v1/owners/{owner}/albums/{folderName}", Method: "DELETE"},
	}

	matched, err := MatchRoute(routes, "DELETE", "/api/v1/owners/john-doe/albums/vacation-2023")
	if err != nil {
		t.Fatalf("Expected match, got error: %v", err)
	}

	if matched.Route.Pattern != "/api/v1/owners/{owner}/albums/{folderName}" {
		t.Errorf("Expected pattern /api/v1/owners/{owner}/albums/{folderName}, got %s", matched.Route.Pattern)
	}

	if matched.PathParams["owner"] != "john-doe" {
		t.Errorf("Expected owner=john-doe, got %s", matched.PathParams["owner"])
	}

	if matched.PathParams["folderName"] != "vacation-2023" {
		t.Errorf("Expected folderName=vacation-2023, got %s", matched.PathParams["folderName"])
	}

	if len(matched.PathParams) != 2 {
		t.Errorf("Expected 2 path params, got %d", len(matched.PathParams))
	}
}

func TestMatchRoute_ShouldMatchCorrectMethodForSamePath(t *testing.T) {
	routes := []Route{
		{Pattern: "/api/v1/albums", Method: "GET"},
		{Pattern: "/api/v1/albums", Method: "POST"},
	}

	// Test GET
	matchedGet, err := MatchRoute(routes, "GET", "/api/v1/albums")
	if err != nil {
		t.Fatalf("Expected GET to match, got error: %v", err)
	}

	if matchedGet.Route.Method != "GET" {
		t.Errorf("Expected method GET, got %s", matchedGet.Route.Method)
	}

	// Test POST
	matchedPost, err := MatchRoute(routes, "POST", "/api/v1/albums")
	if err != nil {
		t.Fatalf("Expected POST to match, got error: %v", err)
	}

	if matchedPost.Route.Method != "POST" {
		t.Errorf("Expected method POST, got %s", matchedPost.Route.Method)
	}
}

func TestMatchRoute_ShouldFailWhenNoRouteMatches(t *testing.T) {
	routes := []Route{
		{Pattern: "/api/v1/albums", Method: "GET"},
	}

	_, err := MatchRoute(routes, "GET", "/api/v1/nonexistent")
	if err == nil {
		t.Error("Expected error for non-matching route, got nil")
	}
}

func TestMatchRoute_ShouldFailWhenMethodDoesNotMatch(t *testing.T) {
	routes := []Route{
		{Pattern: "/api/v1/albums", Method: "GET"},
	}

	_, err := MatchRoute(routes, "POST", "/api/v1/albums")
	if err == nil {
		t.Error("Expected error for non-matching method, got nil")
	}
}

func TestMatchRoute_ShouldMatchRouteWithNestedParameters(t *testing.T) {
	routes := []Route{
		{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/medias", Method: "GET"},
		{Pattern: "/api/v1/owners/{owner}/medias/{mediaId}/{filename}", Method: "GET"},
	}

	// Test first route
	matched1, err := MatchRoute(routes, "GET", "/api/v1/owners/alice/albums/summer-2024/medias")
	if err != nil {
		t.Fatalf("Expected match for first route, got error: %v", err)
	}

	if matched1.PathParams["owner"] != "alice" {
		t.Errorf("Expected owner=alice, got %s", matched1.PathParams["owner"])
	}

	if matched1.PathParams["folderName"] != "summer-2024" {
		t.Errorf("Expected folderName=summer-2024, got %s", matched1.PathParams["folderName"])
	}

	// Test second route
	matched2, err := MatchRoute(routes, "GET", "/api/v1/owners/bob/medias/media123/photo.jpg")
	if err != nil {
		t.Fatalf("Expected match for second route, got error: %v", err)
	}

	if matched2.PathParams["owner"] != "bob" {
		t.Errorf("Expected owner=bob, got %s", matched2.PathParams["owner"])
	}

	if matched2.PathParams["mediaId"] != "media123" {
		t.Errorf("Expected mediaId=media123, got %s", matched2.PathParams["mediaId"])
	}

	if matched2.PathParams["filename"] != "photo.jpg" {
		t.Errorf("Expected filename=photo.jpg, got %s", matched2.PathParams["filename"])
	}
}

func TestMatchRoute_ShouldMatchRouteWithDatesPathSuffix(t *testing.T) {
	routes := []Route{
		{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/dates", Method: "PUT"},
	}

	matched, err := MatchRoute(routes, "PUT", "/api/v1/owners/test-owner/albums/my-album/dates")
	if err != nil {
		t.Fatalf("Expected match, got error: %v", err)
	}

	if matched.PathParams["owner"] != "test-owner" {
		t.Errorf("Expected owner=test-owner, got %s", matched.PathParams["owner"])
	}

	if matched.PathParams["folderName"] != "my-album" {
		t.Errorf("Expected folderName=my-album, got %s", matched.PathParams["folderName"])
	}
}

func TestMatchRoute_ShouldMatchRouteWithNamePathSuffix(t *testing.T) {
	routes := []Route{
		{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/name", Method: "PUT"},
	}

	matched, err := MatchRoute(routes, "PUT", "/api/v1/owners/test-owner/albums/my-album/name")
	if err != nil {
		t.Fatalf("Expected match, got error: %v", err)
	}

	if matched.PathParams["owner"] != "test-owner" {
		t.Errorf("Expected owner=test-owner, got %s", matched.PathParams["owner"])
	}

	if matched.PathParams["folderName"] != "my-album" {
		t.Errorf("Expected folderName=my-album, got %s", matched.PathParams["folderName"])
	}
}

func TestMatchRoute_ShouldMatchSharesRoute(t *testing.T) {
	routes := []Route{
		{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/shares/{email}", Method: "PUT"},
		{Pattern: "/api/v1/owners/{owner}/albums/{folderName}/shares/{email}", Method: "DELETE"},
	}

	// Test PUT
	matchedPut, err := MatchRoute(routes, "PUT", "/api/v1/owners/owner1/albums/album1/shares/user@example.com")
	if err != nil {
		t.Fatalf("Expected PUT to match, got error: %v", err)
	}

	if matchedPut.Route.Method != "PUT" {
		t.Errorf("Expected method PUT, got %s", matchedPut.Route.Method)
	}

	if matchedPut.PathParams["email"] != "user@example.com" {
		t.Errorf("Expected email=user@example.com, got %s", matchedPut.PathParams["email"])
	}

	// Test DELETE
	matchedDelete, err := MatchRoute(routes, "DELETE", "/api/v1/owners/owner1/albums/album1/shares/user@example.com")
	if err != nil {
		t.Fatalf("Expected DELETE to match, got error: %v", err)
	}

	if matchedDelete.Route.Method != "DELETE" {
		t.Errorf("Expected method DELETE, got %s", matchedDelete.Route.Method)
	}
}

func TestMatchRoute_ShouldHandleSpecialCharactersInPathParameters(t *testing.T) {
	routes := []Route{
		{Pattern: "/api/v1/owners/{owner}/albums/{folderName}", Method: "GET"},
	}

	// Test with hyphenated values
	matched, err := MatchRoute(routes, "GET", "/api/v1/owners/john-doe-123/albums/vacation-2023-summer")
	if err != nil {
		t.Fatalf("Expected match, got error: %v", err)
	}

	if matched.PathParams["owner"] != "john-doe-123" {
		t.Errorf("Expected owner=john-doe-123, got %s", matched.PathParams["owner"])
	}

	if matched.PathParams["folderName"] != "vacation-2023-summer" {
		t.Errorf("Expected folderName=vacation-2023-summer, got %s", matched.PathParams["folderName"])
	}
}
