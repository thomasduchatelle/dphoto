package dnsdomain

import (
	"fmt"
	"time"
)

var (
	CertificateNotFoundError = fmt.Errorf("certificate not found")
)

// ExistingCertificate represents a certificate already existing in the infrastructure
type ExistingCertificate struct {
	ID     string // ID is the ARN in the case of AWS
	Domain string
	Expiry time.Time
}

// CompleteCertificate contains PEM encoded certificate
type CompleteCertificate struct {
	Certificate []byte
	Chain       []byte
	PrivateKey  []byte
}

type CertificateManager interface {
	// FindCertificate returns CertificateNotFoundError if the certificate hasn't been found.
	FindCertificate(domain string) (*ExistingCertificate, error)

	// InstallCertificate creates the certificate in ACM (or similar) ; empty 'id' will create a new certificate
	InstallCertificate(id string, certificate CompleteCertificate) error
}

type CertificateAuthority interface {
	RequestCertificate(email, domain string) (*CompleteCertificate, error)
}
