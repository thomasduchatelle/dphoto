package jwks

import "github.com/thomasduchatelle/dphoto/domain/acl/aclcore"

// LoadIssuerConfig reads the OpenId configuration from the provided URLs to provide OAuth2 config used to validate a JWT
func LoadIssuerConfig(issuerConfigUrls ...string) (map[string]aclcore.OAuth2IssuerConfig, error) {
	issuers := make(map[string]aclcore.OAuth2IssuerConfig)
	for _, openIdConfigUrl := range issuerConfigUrls {
		issuerId, config, err := readUrl(openIdConfigUrl)
		if err != nil {
			return nil, err
		}

		issuers[issuerId] = config
	}

	return issuers, nil
}
