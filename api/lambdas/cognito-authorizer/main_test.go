package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineRequiredGroup(t *testing.T) {
	tests := []struct {
		name     string
		routeKey string
		expected string
	}{
		{
			name:     "Admin endpoint - users",
			routeKey: "GET /api/users",
			expected: "admins",
		},
		{
			name:     "Admin endpoint - admin path",
			routeKey: "POST /api/admin/something",
			expected: "admins",
		},
		{
			name:     "Owner endpoint - albums",
			routeKey: "GET /api/albums",
			expected: "owners",
		},
		{
			name:     "Owner endpoint - medias",
			routeKey: "POST /api/medias",
			expected: "owners",
		},
		{
			name:     "Owner endpoint - upload",
			routeKey: "PUT /api/upload",
			expected: "owners",
		},
		{
			name:     "Visitor endpoint - shared albums",
			routeKey: "GET /api/albums/shared",
			expected: "visitors",
		},
		{
			name:     "Visitor endpoint - timeline",
			routeKey: "GET /api/timeline",
			expected: "visitors",
		},
		{
			name:     "Unknown endpoint",
			routeKey: "GET /api/unknown",
			expected: "owners",
		},
		{
			name:     "Invalid route key",
			routeKey: "invalid",
			expected: "admins",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineRequiredGroup(tt.routeKey)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasRequiredPermission(t *testing.T) {
	tests := []struct {
		name           string
		userGroups     []string
		requiredGroup  string
		expectedAccess bool
	}{
		{
			name:           "Admin can access admin endpoints",
			userGroups:     []string{"admins"},
			requiredGroup:  "admins",
			expectedAccess: true,
		},
		{
			name:           "Admin can access owner endpoints",
			userGroups:     []string{"admins"},
			requiredGroup:  "owners",
			expectedAccess: true,
		},
		{
			name:           "Admin can access visitor endpoints",
			userGroups:     []string{"admins"},
			requiredGroup:  "visitors",
			expectedAccess: true,
		},
		{
			name:           "Owner can access owner endpoints",
			userGroups:     []string{"owners"},
			requiredGroup:  "owners",
			expectedAccess: true,
		},
		{
			name:           "Owner can access visitor endpoints",
			userGroups:     []string{"owners"},
			requiredGroup:  "visitors",
			expectedAccess: true,
		},
		{
			name:           "Owner cannot access admin endpoints",
			userGroups:     []string{"owners"},
			requiredGroup:  "admins",
			expectedAccess: false,
		},
		{
			name:           "Visitor can access visitor endpoints",
			userGroups:     []string{"visitors"},
			requiredGroup:  "visitors",
			expectedAccess: true,
		},
		{
			name:           "Visitor cannot access owner endpoints",
			userGroups:     []string{"visitors"},
			requiredGroup:  "owners",
			expectedAccess: false,
		},
		{
			name:           "Visitor cannot access admin endpoints",
			userGroups:     []string{"visitors"},
			requiredGroup:  "admins",
			expectedAccess: false,
		},
		{
			name:           "User with no groups cannot access anything",
			userGroups:     []string{},
			requiredGroup:  "visitors",
			expectedAccess: false,
		},
		{
			name:           "User with multiple groups (admin + owner)",
			userGroups:     []string{"admins", "owners"},
			requiredGroup:  "admins",
			expectedAccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasRequiredPermission(tt.userGroups, tt.requiredGroup)
			assert.Equal(t, tt.expectedAccess, result)
		})
	}
}
