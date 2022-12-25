package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gophercloud/utils/env"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/ephoto/pkg/dns"
	"github.com/thomasduchatelle/ephoto/pkg/dnsadapters/letsencrypt"
	"github.com/thomasduchatelle/ephoto/pkg/dnsadapters/route53_adapter"
)

func Handler() error {
	domain := env.Getenv("DPHOTO_DOMAIN")
	err := dns.RenewCertificate(env.Getenv("DPHOTO_CERTIFICATE_EMAIL"), domain, false)
	if err != nil {
		log.WithField("Domain", domain).WithError(err).Errorln("Renewing SSL certificate failed.")
	}
	return err
}

func main() {
	sess := session.Must(session.NewSession())
	environment := env.Getenv("DPHOTO_ENVIRONMENT")
	dns.CertificateManager = route53_adapter.NewCertificateManager(sess, map[string]string{
		"Application": "dphoto-app",
		"Environment": environment,
		"CreatedBy":   "lambda",
	}, environment)
	dns.CertificateAuthority = letsencrypt.NewCertificateAuthority()

	lambda.Start(Handler)
}
