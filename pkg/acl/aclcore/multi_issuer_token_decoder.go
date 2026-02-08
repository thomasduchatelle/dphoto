package aclcore

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// MultiIssuerTokenDecoder tries multiple token decoders in sequence
// Implements TokenDecoder interface
type MultiIssuerTokenDecoder struct {
	Decoders []TokenDecoder
}

func (m *MultiIssuerTokenDecoder) Decode(accessToken string) (Claims, error) {
	if len(m.Decoders) == 0 {
		return Claims{}, errors.New("no token decoders configured")
	}

	var lastErr error
	for i, decoder := range m.Decoders {
		claims, err := decoder.Decode(accessToken)
		if err == nil {
			log.Debugf("Token successfully decoded by decoder #%d (%T)", i, decoder)
			return claims, nil
		}
		lastErr = err
		log.Debugf("Decoder #%d (%T) failed: %v", i, decoder, err)
	}

	// All decoders failed
	return Claims{}, errors.Wrapf(AccessUnauthorisedError, "all token decoders failed (last error: %s)", lastErr.Error())
}
