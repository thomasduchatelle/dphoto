package route53_adapter

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/thomasduchatelle/dphoto/pkg/dnsdomain"
	"strings"
)

type manager struct {
	acmClient   *acm.ACM
	environment string
	ssmClient   *ssm.SSM
	tags        map[string]string
}

// NewCertificateManager creates an adapter to use Route53
func NewCertificateManager(sess *session.Session, tags map[string]string, environment string) dnsdomain.CertificateManager {
	return &manager{
		acmClient:   acm.New(sess),
		environment: environment,
		ssmClient:   ssm.New(sess),
		tags:        tags,
	}
}

func (m *manager) FindCertificate(domain string) (*dnsdomain.ExistingCertificate, error) {
	certificates, err := m.acmClient.ListCertificates(&acm.ListCertificatesInput{
		MaxItems: aws.Int64(1000),
	})
	if err != nil {
		return nil, err
	}

	for _, c := range certificates.CertificateSummaryList {
		if *c.DomainName == domain {
			certificate, err := m.acmClient.DescribeCertificate(&acm.DescribeCertificateInput{
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

func (m *manager) InstallCertificate(id string, certificate dnsdomain.CompleteCertificate) error {
	importCertificateInput := &acm.ImportCertificateInput{
		Certificate:      onlyFirstCertificate(certificate.Certificate),
		CertificateArn:   m.awsString(id),
		CertificateChain: certificate.Chain,
		PrivateKey:       certificate.PrivateKey,
	}

	putParameterInput := &ssm.PutParameterInput{
		DataType: aws.String("text"),
		Name:     aws.String(fmt.Sprintf("/dphoto/%s/acm/domainCertARN", m.environment)),
		Type:     aws.String("String"),
	}

	for key, value := range m.tags {
		if id == "" {
			// Tagging is not permitted on re-import.
			importCertificateInput.Tags = append(importCertificateInput.Tags, &acm.Tag{
				Key:   &key,
				Value: &value,
			})
		}

		putParameterInput.Tags = append(putParameterInput.Tags, &ssm.Tag{
			Key:   &key,
			Value: &value,
		})
	}

	cer, err := m.acmClient.ImportCertificate(importCertificateInput)
	if err != nil {
		return err
	}

	parameter, err := m.ssmClient.GetParameter(&ssm.GetParameterInput{
		Name: putParameterInput.Name,
	})
	notFound := false
	if err != nil {
		if awsError, ok := err.(awserr.Error); ok && awsError.Code() == ssm.ErrCodeParameterNotFound {
			notFound = true
		} else {
			return err
		}
	}

	if notFound || parameter.Parameter.Value != putParameterInput.Value {
		if !notFound {
			// To update tags for an existing parameter, please use AddTagsToResource or RemoveTagsFromResource
			putParameterInput.Tags = nil
		}
		putParameterInput.Overwrite = aws.Bool(!notFound)
		putParameterInput.Value = cer.CertificateArn

		_, err = m.ssmClient.PutParameter(putParameterInput)
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
