package dns_test

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/dnsdomain"
)

type CertificateManagerInMemory struct {
	Certificates     map[string]dnsdomain.ExistingCertificate
	InstalledContent map[string]dnsdomain.CompleteCertificate
	SSMEnsured       map[string]bool
}

func NewCertificateManagerInMemory() *CertificateManagerInMemory {
	return &CertificateManagerInMemory{
		Certificates:     make(map[string]dnsdomain.ExistingCertificate),
		InstalledContent: make(map[string]dnsdomain.CompleteCertificate),
		SSMEnsured:       make(map[string]bool),
	}
}

func (c *CertificateManagerInMemory) FindCertificate(ctx context.Context, domain string) (*dnsdomain.ExistingCertificate, error) {
	cert, ok := c.Certificates[domain]
	if !ok {
		return nil, dnsdomain.CertificateNotFoundError
	}
	return &cert, nil
}

func (c *CertificateManagerInMemory) InstallCertificate(ctx context.Context, id string, certificate dnsdomain.CompleteCertificate) error {
	c.InstalledContent[id] = certificate
	return nil
}

func (c *CertificateManagerInMemory) EnsureSSMParameter(ctx context.Context, certificateArn string) error {
	c.SSMEnsured[certificateArn] = true
	return nil
}

func (c *CertificateManagerInMemory) IsSSMEnsured(arn string) bool {
	return c.SSMEnsured[arn]
}

func (c *CertificateManagerInMemory) Installed(arn string) (dnsdomain.CompleteCertificate, bool) {
	cert, ok := c.InstalledContent[arn]
	return cert, ok
}

type CertificateRequest struct {
	Email  string
	Domain string
}

type CertificateAuthorityInMemory struct {
	NextCertificate *dnsdomain.CompleteCertificate
	Requested       []CertificateRequest
}

func (c *CertificateAuthorityInMemory) RequestCertificate(ctx context.Context, email, domain string) (*dnsdomain.CompleteCertificate, error) {
	c.Requested = append(c.Requested, CertificateRequest{Email: email, Domain: domain})
	return c.NextCertificate, nil
}

func (c *CertificateAuthorityInMemory) RequestedFor() []CertificateRequest {
	return c.Requested
}
