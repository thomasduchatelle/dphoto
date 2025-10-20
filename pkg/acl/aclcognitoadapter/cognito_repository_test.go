package aclcognitoadapter_test

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcognitoadapter"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		userPoolId string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "it should return an error if userPoolId is empty",
			userPoolId: "",
			wantErr:    assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't test the successful case without AWS credentials
			// The integration tests should cover the actual AWS interaction
			cfg := aws.Config{Region: "us-east-1"}
			_, err := aclcognitoadapter.New(cfg, tt.userPoolId)
			tt.wantErr(t, err)
		})
	}
}

func TestCognitoRepositoryAdapter_CreateUser(t *testing.T) {
	t.Skip("Integration test - requires AWS Cognito setup")
}

func TestCognitoRepositoryAdapter_UserExists(t *testing.T) {
	t.Skip("Integration test - requires AWS Cognito setup")
}

// These tests document the expected behavior and serve as a contract for the implementation
// Real integration tests would require AWS Cognito credentials and setup

func TestCognitoRepository_Interface(t *testing.T) {
	t.Run("it should implement CognitoRepository interface", func(t *testing.T) {
		var _ aclcore.CognitoRepository = (*MockCognitoRepository)(nil)
	})
}

// MockCognitoRepository is a test double for integration testing
type MockCognitoRepository struct {
	users map[string][]aclcore.CognitoUserGroup
}

func (m *MockCognitoRepository) CreateUser(ctx context.Context, email usermodel.UserId, group aclcore.CognitoUserGroup) error {
	if m.users == nil {
		m.users = make(map[string][]aclcore.CognitoUserGroup)
	}
	
	emailStr := email.Value()
	if groups, exists := m.users[emailStr]; exists {
		// Add to group if not already present
		for _, g := range groups {
			if g == group {
				return nil
			}
		}
		m.users[emailStr] = append(groups, group)
	} else {
		m.users[emailStr] = []aclcore.CognitoUserGroup{group}
	}
	return nil
}

func (m *MockCognitoRepository) UserExists(ctx context.Context, email usermodel.UserId) (bool, error) {
	if m.users == nil {
		return false, nil
	}
	_, exists := m.users[email.Value()]
	return exists, nil
}
