package dns_test

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/dnsdomain"
)

const generatedArn = "arn::generated"

type InMemoryCertificate struct {
	ExistingCertificate dnsdomain.ExistingCertificate
	CompleteCertificate *dnsdomain.CompleteCertificate
}
type CertificateManagerInMemory struct {
	Certificates map[string]InMemoryCertificate
	SSMParameter string // ID (ARN) of the installed certificate
}

// NewCertificateManagerInMemory creates an in memory implementation where the first certificate is considered as the active one.
func NewCertificateManagerInMemory(certs ...dnsdomain.ExistingCertificate) *CertificateManagerInMemory {
	m := &CertificateManagerInMemory{
		Certificates: make(map[string]InMemoryCertificate),
	}
	for i, c := range certs {
		m.Certificates[c.Domain] = InMemoryCertificate{
			ExistingCertificate: c,
			CompleteCertificate: nil,
		}
		if i == 0 {
			m.SSMParameter = c.ID
		}
	}
	return m
}

func (m *CertificateManagerInMemory) FindCertificate(_ context.Context, domain string) (*dnsdomain.ExistingCertificate, error) {
	cert, ok := m.Certificates[domain]
	if !ok {
		return nil, dnsdomain.CertificateNotFoundError
	}
	existing := cert.ExistingCertificate
	return &existing, nil
}

func (m *CertificateManagerInMemory) InstallCertificate(_ context.Context, id string, certificate dnsdomain.CompleteCertificate) error {
	if id == "" {
		id = generatedArn
	}
	for domain, entry := range m.Certificates {
		if entry.ExistingCertificate.ID == id {
			body := certificate
			entry.CompleteCertificate = &body
			entry.ExistingCertificate.ID = id
			m.Certificates[domain] = entry
			return nil
		}
	}
	body := certificate
	m.Certificates[""] = InMemoryCertificate{
		ExistingCertificate: dnsdomain.ExistingCertificate{ID: id},
		CompleteCertificate: &body,
	}
	return nil
}

func (m *CertificateManagerInMemory) EnsureSSMParameter(_ context.Context, certificateArn string) error {
	m.SSMParameter = certificateArn
	return nil
}

type CertificateRequest struct {
	Email  string
	Domain string
}

type CertificateAuthorityInMemory struct {
	NextCertificate dnsdomain.CompleteCertificate
	Requested       []CertificateRequest
}

func NewCertificateAuthorityInMemory(next dnsdomain.CompleteCertificate) *CertificateAuthorityInMemory {
	return &CertificateAuthorityInMemory{NextCertificate: next}
}

func (a *CertificateAuthorityInMemory) RequestCertificate(_ context.Context, email, domain string) (*dnsdomain.CompleteCertificate, error) {
	a.Requested = append(a.Requested, CertificateRequest{Email: email, Domain: domain})
	cert := a.NextCertificate
	return &cert, nil
}
