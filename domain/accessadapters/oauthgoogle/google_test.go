package oauthgoogle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestOAuth2ConfigReader_Read is requiring WEB access <- WARNING
func TestOAuth2ConfigReader_Read(t *testing.T) {
	a := assert.New(t)

	// this test requires access to internet (google API)
	gotName, _, err := NewGoogle().Read()
	if a.NoError(err) {
		a.Equal("accounts.google.com", gotName)

		// KID are dynamic and can't be asserting here
	}
}
