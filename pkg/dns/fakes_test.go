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

func NewCertificateManagerInMemory(certs ...dnsdomain.ExistingCertificate) *CertificateManagerInMemory {
	m := &CertificateManagerInMemory{
		Certificates:     make(map[string]dnsdomain.ExistingCertificate),
		InstalledContent: make(map[string]dnsdomain.CompleteCertificate),
		SSMEnsured:       make(map[string]bool),
	}
	for _, c := range certs {
		m.Certificates[c.Domain] = c
	}
	return m
}

func (m *CertificateManagerInMemory) FindCertificate(_ context.Context, domain string) (*dnsdomain.ExistingCertificate, error) {
	cert, ok := m.Certificates[domain]
	if !ok {
		return nil, dnsdomain.CertificateNotFoundError
	}
	return &cert, nil
}

func (m *CertificateManagerInMemory) InstallCertificate(_ context.Context, id string, certificate dnsdomain.CompleteCertificate) error {
	m.InstalledContent[id] = certificate
	return nil
}

func (m *CertificateManagerInMemory) EnsureSSMParameter(_ context.Context, certificateArn string) error {
	m.SSMEnsured[certificateArn] = true
	return nil
}

func (m *CertificateManagerInMemory) IsSSMEnsured(arn string) bool {
	return m.SSMEnsured[arn]
}

func (m *CertificateManagerInMemory) Installed(arn string) (dnsdomain.CompleteCertificate, bool) {
	cert, ok := m.InstalledContent[arn]
	return cert, ok
}

type CertificateAuthorityInMemory struct {
	NextCertificate dnsdomain.CompleteCertificate
}

func NewCertificateAuthorityInMemory(next dnsdomain.CompleteCertificate) *CertificateAuthorityInMemory {
	return &CertificateAuthorityInMemory{NextCertificate: next}
}

func (a *CertificateAuthorityInMemory) RequestCertificate(_ context.Context, _, _ string) (*dnsdomain.CompleteCertificate, error) {
	cert := a.NextCertificate
	return &cert, nil
}
