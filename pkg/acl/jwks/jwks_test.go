package jwks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestOAuth2ConfigReader_Read is requiring WEB access <- WARNING
func TestOAuth2ConfigReader_Read(t *testing.T) {
	a := assert.New(t)

	// this test requires access to internet (google API)
	config, err := LoadIssuerConfig("https://accounts.google.com/.well-known/openid-configuration")
	if a.NoError(err) {
		iss, ok := config["accounts.google.com"]
		if assert.Truef(t, ok, "got %+v", config) {
			assert.Equal(t, "https://accounts.google.com/.well-known/openid-configuration", iss.ConfigSource)
			// KID are dynamic and can't be asserting here
		}
	}
}
