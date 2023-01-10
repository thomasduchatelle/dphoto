package dns_test

import (
	"github.com/stretchr/testify/assert"
	mocks2 "github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/dns"
	"github.com/thomasduchatelle/dphoto/pkg/dnsdomain"
	"testing"
	"time"
)

func TestRenewCertificate(t *testing.T) {
	a := assert.New(t)

	const (
		domain = "dphoto.example.com"
		email  = "dphoto@example.com"
	)
	tests := []struct {
		name     string
		setMocks func(certManager *mocks2.CertificateManager, certAuthority *mocks2.CertificateAuthority)
	}{
		{"it should not create a new certificate if one already exists", func(certManager *mocks2.CertificateManager, certAuthority *mocks2.CertificateAuthority) {
			certManager.On("FindCertificate", domain).Return(&dnsdomain.ExistingCertificate{
				ID:     "arn::132456",
				Domain: domain,
				Expiry: time.Now().Add(dns.MinimumExpiryDelay * 2),
			}, nil)
		}},
		{"it should create a new certificate if the existing one has or is about to expire, and override it", func(certManager *mocks2.CertificateManager, certAuthority *mocks2.CertificateAuthority) {
			certManager.On("FindCertificate", domain).Return(&dnsdomain.ExistingCertificate{
				ID:     "arn::132456",
				Domain: domain,
				Expiry: time.Now().Add(dns.MinimumExpiryDelay - time.Hour),
			}, nil)

			certManager.On("InstallCertificate", "arn::132456", dnsdomain.CompleteCertificate{
				Certificate: []byte("cert-123"),
				Chain:       []byte("chain-123"),
				PrivateKey:  []byte("private-key-123"),
			}).Return(nil)

			certAuthority.On("RequestCertificate", email, domain).Return(&dnsdomain.CompleteCertificate{
				Certificate: []byte("cert-123"),
				Chain:       []byte("chain-123"),
				PrivateKey:  []byte("private-key-123"),
			}, nil)
		}},
		{"it should create a new certificate if none were there", func(certManager *mocks2.CertificateManager, certAuthority *mocks2.CertificateAuthority) {
			certManager.On("FindCertificate", domain).Return(nil, dnsdomain.CertificateNotFoundError)

			certManager.On("InstallCertificate", "", dnsdomain.CompleteCertificate{
				Certificate: []byte("cert-123"),
				Chain:       []byte("chain-123"),
				PrivateKey:  []byte("private-key-123"),
			}).Return(nil)

			certAuthority.On("RequestCertificate", email, domain).Return(&dnsdomain.CompleteCertificate{
				Certificate: []byte("cert-123"),
				Chain:       []byte("chain-123"),
				PrivateKey:  []byte("private-key-123"),
			}, nil)

		}},
	}

	for _, tt := range tests {
		certManager := new(mocks2.CertificateManager)
		dns.CertificateManager = certManager
		certAuthority := new(mocks2.CertificateAuthority)
		dns.CertificateAuthority = certAuthority

		tt.setMocks(certManager, certAuthority)

		err := dns.RenewCertificate(email, domain, false)
		if a.NoError(err, tt.name) {
			certManager.AssertExpectations(t)
			certAuthority.AssertExpectations(t)
		}
	}
}
