package route53_adapter

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	acmtypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/thomasduchatelle/dphoto/pkg/dnsdomain"
)

type manager struct {
	acmClient       *acm.Client
	ssmParameterKey string
	ssmClient       *ssm.Client
	tags            map[string]string
}

// NewCertificateManager creates an adapter to use Route53 ; 'environment' is deprecated, the SSM parameter name should be used instead
func NewCertificateManager(cfg aws.Config, tags map[string]string, environment string, ssmParameterKey string) dnsdomain.CertificateManager {
	key := ssmParameterKey
	if key == "" {
		key = fmt.Sprintf("/dphoto/%s/acm/domainCertARN", environment)
	}

	return &manager{
		acmClient:       acm.NewFromConfig(cfg),
		ssmParameterKey: key,
		ssmClient:       ssm.NewFromConfig(cfg),
		tags:            tags,
	}
}

func (m *manager) FindCertificate(ctx context.Context, domain string) (*dnsdomain.ExistingCertificate, error) {
	certificates, err := m.acmClient.ListCertificates(ctx, &acm.ListCertificatesInput{
		MaxItems: aws.Int32(1000),
	})
	if err != nil {
		return nil, err
	}

	for _, c := range certificates.CertificateSummaryList {
		if *c.DomainName == domain {
			certificate, err := m.acmClient.DescribeCertificate(ctx, &acm.DescribeCertificateInput{
				CertificateArn: c.CertificateArn,
			})
			if err != nil {
				return nil, err
			}

			return &dnsdomain.ExistingCertificate{
				ID:     *certificate.Certificate.CertificateArn,
				Domain: *c.DomainName,
				Expiry: *certificate.Certificate.NotAfter,
			}, nil
		}
	}

	return nil, dnsdomain.CertificateNotFoundError
}

func (m *manager) InstallCertificate(ctx context.Context, id string, certificate dnsdomain.CompleteCertificate) error {
	importCertificateInput := &acm.ImportCertificateInput{
		Certificate:      onlyFirstCertificate(certificate.Certificate),
		CertificateArn:   m.awsString(id),
		CertificateChain: certificate.Chain,
		PrivateKey:       certificate.PrivateKey,
	}

	// Tagging is not permitted on re-import.
	if id == "" {
		for key, value := range m.tags {
			importCertificateInput.Tags = append(importCertificateInput.Tags, acmtypes.Tag{
				Key:   &key,
				Value: &value,
			})
		}
	}

	cer, err := m.acmClient.ImportCertificate(ctx, importCertificateInput)
	if err != nil {
		return err
	}

	return m.EnsureSSMParameter(ctx, *cer.CertificateArn)
}

func onlyFirstCertificate(certificate []byte) []byte {
	var cer []string
	for _, line := range strings.Split(string(certificate), "\n") {
		cer = append(cer, line)
		if strings.Trim(line, " ") == "-----END CERTIFICATE-----" {
			return []byte(strings.Join(cer, "\n"))
		}
	}

	return certificate
}

func (m *manager) EnsureSSMParameter(ctx context.Context, certificateArn string) error {
	putParameterInput := &ssm.PutParameterInput{
		DataType: aws.String("text"),
		Name:     aws.String(m.ssmParameterKey),
		Type:     ssmtypes.ParameterTypeString,
		Value:    aws.String(certificateArn),
	}

	notFound, needsUpdating, err := m.isNewOrNeedUpdating(ctx, *putParameterInput.Name, certificateArn)
	if err != nil {
		return err
	}
	if notFound {
		// Tags cannot be updated via PutParameter, must use AddTagsToResource or RemoveTagsFromResource ; they will remain out of sync.
		for key, value := range m.tags {
			putParameterInput.Tags = append(putParameterInput.Tags, ssmtypes.Tag{
				Key:   &key,
				Value: &value,
			})
		}

		_, err = m.ssmClient.PutParameter(ctx, putParameterInput)
		return err
	}

	if needsUpdating {
		putParameterInput.Overwrite = aws.Bool(true)

		_, err = m.ssmClient.PutParameter(ctx, putParameterInput)
		return err
	}

	return nil
}

func (m *manager) isNewOrNeedUpdating(ctx context.Context, parameterName string, expectedValue string) (notFound bool, needsUpdating bool, err error) {
	parameter, err := m.ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name: &parameterName,
	})

	var parameterNotFound *ssmtypes.ParameterNotFound
	notFound = errors.As(err, &parameterNotFound)
	if notFound {
		err = nil
	}
	if err != nil || notFound {
		return
	}

	needsUpdating = *parameter.Parameter.Value != expectedValue
	return
}

func (m *manager) awsString(id string) *string {
	if id == "" {
		return nil
	}

	return &id
}
