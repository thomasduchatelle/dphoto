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
)

var cannedCertificate = dnsdomain.CompleteCertificate{
	Certificate: []byte("cert-123"),
	Chain:       []byte("chain-123"),
	PrivateKey:  []byte("private-key-123"),
}

func TestRenewCertificate(t *testing.T) {
	type fields struct {
		CertificateManager   *CertificateManagerInMemory
		CertificateAuthority *CertificateAuthorityInMemory
	}
	type args struct {
		email  string
		domain string
		forced bool
	}
	installedAt := func(id string) *string { return &id }

	tests := []struct {
		name                string
		fields              fields
		args                args
		wantErr             assert.ErrorAssertionFunc
		expectSSMEnsuredFor string
		expectInstalledAt   *string
	}{
		{
			name: "it should not create a new certificate if one already exists",
			fields: fields{
				CertificateManager: NewCertificateManagerInMemory(dnsdomain.ExistingCertificate{
					ID:     testArn,
					Domain: testDomain,
					Expiry: time.Now().Add(dns.MinimumExpiryDelay * 2),
				}),
				CertificateAuthority: NewCertificateAuthorityInMemory(cannedCertificate),
			},
			args:                args{email: testEmail, domain: testDomain, forced: false},
			wantErr:             assert.NoError,
			expectSSMEnsuredFor: testArn,
		},
		{
			name: "it should create a new certificate if the existing one is about to expire, and override it",
			fields: fields{
				CertificateManager: NewCertificateManagerInMemory(dnsdomain.ExistingCertificate{
					ID:     testArn,
					Domain: testDomain,
					Expiry: time.Now().Add(dns.MinimumExpiryDelay - time.Hour),
				}),
				CertificateAuthority: NewCertificateAuthorityInMemory(cannedCertificate),
			},
			args:              args{email: testEmail, domain: testDomain, forced: false},
			wantErr:           assert.NoError,
			expectInstalledAt: installedAt(testArn),
		},
		{
			name: "it should create a new certificate if none were there",
			fields: fields{
				CertificateManager:   NewCertificateManagerInMemory(),
				CertificateAuthority: NewCertificateAuthorityInMemory(cannedCertificate),
			},
			args:              args{email: testEmail, domain: testDomain, forced: false},
			wantErr:           assert.NoError,
			expectInstalledAt: installedAt(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dns.CertificateManager = tt.fields.CertificateManager
			dns.CertificateAuthority = tt.fields.CertificateAuthority

			err := dns.RenewCertificate(tt.args.email, tt.args.domain, tt.args.forced)
			if !tt.wantErr(t, err) {
				return
			}

			if tt.expectSSMEnsuredFor != "" {
				assert.True(t, tt.fields.CertificateManager.IsSSMEnsured(tt.expectSSMEnsuredFor), "expected SSM parameter ensured for %s", tt.expectSSMEnsuredFor)
			}
			if tt.expectInstalledAt == nil {
				assert.Empty(t, tt.fields.CertificateManager.InstalledContent, "expected no certificate installed")
			} else {
				installed, ok := tt.fields.CertificateManager.Installed(*tt.expectInstalledAt)
				if assert.True(t, ok, "expected a certificate installed at %q", *tt.expectInstalledAt) {
					assert.Equal(t, cannedCertificate, installed)
				}
			}
		})
	}
}
