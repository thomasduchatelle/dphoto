package route53_adapter

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	acmtypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/thomasduchatelle/dphoto/pkg/dnsdomain"
	"strings"
)

type manager struct {
	acmClient   *acm.Client
	environment string
	ssmClient   *ssm.Client
	tags        map[string]string
}

// NewCertificateManager creates an adapter to use Route53
func NewCertificateManager(cfg aws.Config, tags map[string]string, environment string) dnsdomain.CertificateManager {
	return &manager{
		acmClient:   acm.NewFromConfig(cfg),
		environment: environment,
		ssmClient:   ssm.NewFromConfig(cfg),
		tags:        tags,
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

	putParameterInput := &ssm.PutParameterInput{
		DataType: aws.String("text"),
		Name:     aws.String(fmt.Sprintf("/dphoto/%s/acm/domainCertARN", m.environment)),
		Type:     ssmtypes.ParameterTypeString,
	}

	for key, value := range m.tags {
		if id == "" {
			// Tagging is not permitted on re-import.
			importCertificateInput.Tags = append(importCertificateInput.Tags, acmtypes.Tag{
				Key:   &key,
				Value: &value,
			})
		}

		putParameterInput.Tags = append(putParameterInput.Tags, ssmtypes.Tag{
			Key:   &key,
			Value: &value,
		})
	}

	cer, err := m.acmClient.ImportCertificate(ctx, importCertificateInput)
	if err != nil {
		return err
	}

	parameter, err := m.ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name: putParameterInput.Name,
	})

	var parameterNotFound *ssmtypes.ParameterNotFound
	notFound := errors.As(err, &parameterNotFound)
	if err != nil && !notFound {
		return err
	}

	if notFound || parameter.Parameter.Value != putParameterInput.Value {
		if !notFound {
			// "To update tags for an existing parameter, please use AddTagsToResource or RemoveTagsFromResource"
			putParameterInput.Tags = nil
		}
		putParameterInput.Overwrite = aws.Bool(!notFound)
		putParameterInput.Value = cer.CertificateArn

		_, err = m.ssmClient.PutParameter(ctx, putParameterInput)
		return err
	}

	return nil
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

func (m *manager) awsString(id string) *string {
	if id == "" {
		return nil
	}

	return &id
}
