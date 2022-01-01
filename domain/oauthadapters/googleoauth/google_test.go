package googleoauth

import (
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/oauthmodel"
	"testing"
)

// TestOAuth2ConfigReader_Read is requiring WEB access <- WARNING
func TestOAuth2ConfigReader_Read(t *testing.T) {
	a := assert.New(t)

	// this test requires access to internet (google API)
	config := oauthmodel.Config{}
	err := NewGoogle().Read(&config)
	if a.NoError(err) {
		auth2Config, ok := config.TrustedIssuers["accounts.google.com"]
		a.Truef(ok, "issuer not found in %+v", auth2Config)

		// KID are dynamic and can't be asserting here
	}
}
