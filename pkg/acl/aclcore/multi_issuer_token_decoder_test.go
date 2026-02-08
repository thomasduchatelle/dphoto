package aclcore

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// Mock decoder for testing
type mockTokenDecoder struct {
	shouldSucceed bool
	claims        Claims
	err           error
}

func (m *mockTokenDecoder) Decode(accessToken string) (Claims, error) {
	if m.shouldSucceed {
		return m.claims, nil
	}
	return Claims{}, m.err
}

func TestMultiIssuerTokenDecoder_Decode(t *testing.T) {
	t.Run("First Decoder Succeeds", func(t *testing.T) {
		expectedClaims := Claims{
			Subject: usermodel.NewUserId("user@example.com"),
			Scopes:  map[string]interface{}{"api:admin": nil},
		}

		decoder1 := &mockTokenDecoder{
			shouldSucceed: true,
			claims:        expectedClaims,
		}
		decoder2 := &mockTokenDecoder{
			shouldSucceed: false,
			err:           errors.New("should not be called"),
		}

		multiDecoder := &MultiIssuerTokenDecoder{
			Decoders: []TokenDecoder{decoder1, decoder2},
		}

		claims, err := multiDecoder.Decode("test-token")
		assert.NoError(t, err)
		assert.Equal(t, expectedClaims.Subject, claims.Subject)
		assert.Equal(t, expectedClaims.Scopes, claims.Scopes)
	})

	t.Run("Second Decoder Succeeds", func(t *testing.T) {
		expectedClaims := Claims{
			Subject: usermodel.NewUserId("cognito@example.com"),
			Scopes:  map[string]interface{}{"visitor": nil},
		}

		decoder1 := &mockTokenDecoder{
			shouldSucceed: false,
			err:           errors.Wrap(AccessUnauthorisedError, "internal token invalid"),
		}
		decoder2 := &mockTokenDecoder{
			shouldSucceed: true,
			claims:        expectedClaims,
		}

		multiDecoder := &MultiIssuerTokenDecoder{
			Decoders: []TokenDecoder{decoder1, decoder2},
		}

		claims, err := multiDecoder.Decode("test-token")
		assert.NoError(t, err)
		assert.Equal(t, expectedClaims.Subject, claims.Subject)
		assert.Equal(t, expectedClaims.Scopes, claims.Scopes)
	})

	t.Run("All Decoders Fail", func(t *testing.T) {
		decoder1 := &mockTokenDecoder{
			shouldSucceed: false,
			err:           errors.Wrap(AccessUnauthorisedError, "decoder1 failed"),
		}
		decoder2 := &mockTokenDecoder{
			shouldSucceed: false,
			err:           errors.Wrap(AccessUnauthorisedError, "decoder2 failed"),
		}

		multiDecoder := &MultiIssuerTokenDecoder{
			Decoders: []TokenDecoder{decoder1, decoder2},
		}

		_, err := multiDecoder.Decode("test-token")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, AccessUnauthorisedError))
		assert.Contains(t, err.Error(), "all token decoders failed")
	})

	t.Run("Empty Decoders List", func(t *testing.T) {
		multiDecoder := &MultiIssuerTokenDecoder{
			Decoders: []TokenDecoder{},
		}

		_, err := multiDecoder.Decode("test-token")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no token decoders configured")
	})

	t.Run("Multiple Decoders With Owner", func(t *testing.T) {
		owner := ownermodel.Owner("test-owner")
		expectedClaims := Claims{
			Subject: usermodel.NewUserId("user@example.com"),
			Scopes:  map[string]interface{}{"owner:test-owner": nil},
			Owner:   &owner,
		}

		decoder1 := &mockTokenDecoder{
			shouldSucceed: false,
			err:           errors.Wrap(AccessUnauthorisedError, "failed"),
		}
		decoder2 := &mockTokenDecoder{
			shouldSucceed: true,
			claims:        expectedClaims,
		}
		decoder3 := &mockTokenDecoder{
			shouldSucceed: false,
			err:           errors.New("should not be called"),
		}

		multiDecoder := &MultiIssuerTokenDecoder{
			Decoders: []TokenDecoder{decoder1, decoder2, decoder3},
		}

		claims, err := multiDecoder.Decode("test-token")
		assert.NoError(t, err)
		assert.Equal(t, expectedClaims.Subject, claims.Subject)
		assert.NotNil(t, claims.Owner)
		assert.Equal(t, *expectedClaims.Owner, *claims.Owner)
	})
}
