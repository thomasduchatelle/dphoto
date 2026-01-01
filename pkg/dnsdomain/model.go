package dnsdomain

import (
	"context"
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
	FindCertificate(ctx context.Context, domain string) (*ExistingCertificate, error)

	// InstallCertificate creates the certificate in ACM (or similar) ; empty 'id' will create a new certificate
	InstallCertificate(ctx context.Context, id string, certificate CompleteCertificate) error

	// EnsureSSMParameter checks if the SSM parameter exists and updates it if it doesn't match the certificate ARN
	EnsureSSMParameter(ctx context.Context, certificateArn string) error
}

type CertificateAuthority interface {
	RequestCertificate(ctx context.Context, email, domain string) (*CompleteCertificate, error)
}
