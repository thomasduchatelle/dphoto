package dns_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/dns"
	"github.com/thomasduchatelle/dphoto/pkg/dnsdomain"
)

func TestRenewCertificate(t *testing.T) {
	const (
		domain = "dphoto.example.com"
		email  = "dphoto@example.com"
		arn    = "arn::132456"
	)

	cannedCertificate := dnsdomain.CompleteCertificate{
		Certificate: []byte("cert-123"),
		Chain:       []byte("chain-123"),
		PrivateKey:  []byte("private-key-123"),
	}

	validExistingCertificate := dnsdomain.ExistingCertificate{
		ID:     arn,
		Domain: domain,
		Expiry: time.Now().Add(dns.MinimumExpiryDelay * 2),
	}

	newHappyCertificateManager := func() *CertificateManagerInMemory {
		certManager := NewCertificateManagerInMemory()
		certManager.Certificates[domain] = validExistingCertificate
		return certManager
	}
	newHappyCertificateAuthority := func() *CertificateAuthorityInMemory {
		return &CertificateAuthorityInMemory{NextCertificate: &cannedCertificate}
	}

	type fields struct {
		CertificateManager   *CertificateManagerInMemory
		CertificateAuthority *CertificateAuthorityInMemory
	}
	type args struct {
		email  string
		domain string
		forced bool
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantSSMEnsuredArn string
		wantInstalledArn  string
		wantRequested     []CertificateRequest
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "it should not create a new certificate if one already exists",
			fields: fields{
				CertificateManager:   newHappyCertificateManager(),
				CertificateAuthority: newHappyCertificateAuthority(),
			},
			args:              args{email: email, domain: domain, forced: false},
			wantSSMEnsuredArn: arn,
			wantErr:           assert.NoError,
		},
		{
			name: "it should create a new certificate if the existing one has or is about to expire, and override it",
			fields: fields{
				CertificateManager: func() *CertificateManagerInMemory {
					certManager := NewCertificateManagerInMemory()
					certManager.Certificates[domain] = dnsdomain.ExistingCertificate{
						ID:     arn,
						Domain: domain,
						Expiry: time.Now().Add(dns.MinimumExpiryDelay - time.Hour),
					}
					return certManager
				}(),
				CertificateAuthority: newHappyCertificateAuthority(),
			},
			args:             args{email: email, domain: domain, forced: false},
			wantInstalledArn: arn,
			wantRequested:    []CertificateRequest{{Email: email, Domain: domain}},
			wantErr:          assert.NoError,
		},
		{
			name: "it should create a new certificate if none were there",
			fields: fields{
				CertificateManager:   NewCertificateManagerInMemory(),
				CertificateAuthority: newHappyCertificateAuthority(),
			},
			args:             args{email: email, domain: domain, forced: false},
			wantInstalledArn: "",
			wantRequested:    []CertificateRequest{{Email: email, Domain: domain}},
			wantErr:          assert.NoError,
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
			if tt.wantSSMEnsuredArn != "" {
				assert.True(t, tt.fields.CertificateManager.IsSSMEnsured(tt.wantSSMEnsuredArn), "expected SSM parameter to be ensured for %s", tt.wantSSMEnsuredArn)
				assert.Empty(t, tt.fields.CertificateAuthority.RequestedFor(), "expected no certificate request to the authority")
			} else {
				installed, ok := tt.fields.CertificateManager.Installed(tt.wantInstalledArn)
				if assert.True(t, ok, "expected a certificate installed for arn %q", tt.wantInstalledArn) {
					assert.Equal(t, cannedCertificate, installed)
				}
				assert.Equal(t, tt.wantRequested, tt.fields.CertificateAuthority.RequestedFor())
			}
		})
	}
}
