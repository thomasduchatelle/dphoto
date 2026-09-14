package dns_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/dns"
	"github.com/thomasduchatelle/dphoto/pkg/dnsdomain"
)

const (
	testDomain = "dphoto.example.com"
	testEmail  = "dphoto@example.com"
	testArn    = "arn::132456"
	staleArn   = "arn::stale"
)

var cannedCertificate = dnsdomain.CompleteCertificate{
	Certificate: []byte("cert-123"),
	Chain:       []byte("chain-123"),
	PrivateKey:  []byte("private-key-123"),
}

func TestRenewCertificate(t *testing.T) {
	type args struct {
		email  string
		domain string
		forced bool
	}

	validExpiry := time.Now().Add(dns.MinimumExpiryDelay * 2)
	expiringExpiry := time.Now().Add(dns.MinimumExpiryDelay - time.Hour)

	tests := []struct {
		name                  string
		certificateManager    *CertificateManagerInMemory
		args                  args
		wantErr               assert.ErrorAssertionFunc
		expectDomainsContains map[string]InMemoryCertificate
		expectSSMParameter    string
		expectCARequests      []CertificateRequest
	}{
		{
			name: "it should ensure the SSM parameter points at the existing certificate when it is still valid",
			certificateManager: NewCertificateManagerInMemory(dnsdomain.ExistingCertificate{
				ID:     testArn,
				Domain: testDomain,
				Expiry: validExpiry,
			}),
			args:    args{email: testEmail, domain: testDomain, forced: false},
			wantErr: assert.NoError,
			expectDomainsContains: map[string]InMemoryCertificate{
				testDomain: {
					ExistingCertificate: dnsdomain.ExistingCertificate{ID: testArn, Domain: testDomain, Expiry: validExpiry},
				},
			},
			expectSSMParameter: testArn,
			expectCARequests:   nil,
		},
		{
			name: "it should overwrite the SSM parameter when it holds a stale ARN and the existing certificate is still valid",
			certificateManager: func() *CertificateManagerInMemory {
				m := NewCertificateManagerInMemory(dnsdomain.ExistingCertificate{
					ID:     testArn,
					Domain: testDomain,
					Expiry: validExpiry,
				})
				m.SSMParameter = staleArn
				return m
			}(),
			args:    args{email: testEmail, domain: testDomain, forced: false},
			wantErr: assert.NoError,
			expectDomainsContains: map[string]InMemoryCertificate{
				testDomain: {
					ExistingCertificate: dnsdomain.ExistingCertificate{ID: testArn, Domain: testDomain, Expiry: validExpiry},
				},
			},
			expectSSMParameter: testArn,
			expectCARequests:   nil,
		},
		{
			name: "it should install a new certificate and override the existing one when it is about to expire",
			certificateManager: NewCertificateManagerInMemory(dnsdomain.ExistingCertificate{
				ID:     testArn,
				Domain: testDomain,
				Expiry: expiringExpiry,
			}),
			args:    args{email: testEmail, domain: testDomain, forced: false},
			wantErr: assert.NoError,
			expectDomainsContains: map[string]InMemoryCertificate{
				testDomain: {
					ExistingCertificate: dnsdomain.ExistingCertificate{ID: testArn, Domain: testDomain, Expiry: expiringExpiry},
					CompleteCertificate: &cannedCertificate,
				},
			},
			expectSSMParameter: testArn,
			expectCARequests:   []CertificateRequest{{Email: testEmail, Domain: testDomain}},
		},
		{
			name:               "it should install a new certificate when none exists",
			certificateManager: NewCertificateManagerInMemory(),
			args:               args{email: testEmail, domain: testDomain, forced: false},
			wantErr:            assert.NoError,
			expectDomainsContains: map[string]InMemoryCertificate{
				"": {
					ExistingCertificate: dnsdomain.ExistingCertificate{ID: generatedArn},
					CompleteCertificate: &cannedCertificate,
				},
			},
			expectSSMParameter: "",
			expectCARequests:   []CertificateRequest{{Email: testEmail, Domain: testDomain}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := NewCertificateAuthorityInMemory(cannedCertificate)
			dns.CertificateManager = tt.certificateManager
			dns.CertificateAuthority = ca

			err := dns.RenewCertificate(tt.args.email, tt.args.domain, tt.args.forced)
			if !tt.wantErr(t, err) {
				return
			}

			assert.Equal(t, tt.expectSSMParameter, tt.certificateManager.SSMParameter, "SSM parameter value")
			assert.Equal(t, tt.expectCARequests, ca.Requested, "certificate authority requests")
			assert.Equal(t, tt.expectDomainsContains, tt.certificateManager.Certificates, "certificates in repository")
		})
	}
}
