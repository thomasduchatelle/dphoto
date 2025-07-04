package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gophercloud/utils/env"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/dns"
	"github.com/thomasduchatelle/dphoto/pkg/dnsadapters/letsencrypt"
	"github.com/thomasduchatelle/dphoto/pkg/dnsadapters/route53_adapter"
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
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(errors.Wrapf(err, "invalid credentials").Error())
	}
	environment := env.Getenv("DPHOTO_ENVIRONMENT")
	ssmKeyCertificateArn := env.Getenv("SSM_KEY_CERTIFICATE_ARN")
	dns.CertificateManager = route53_adapter.NewCertificateManager(
		cfg,
		map[string]string{
			"Application": "dphoto-app",
			"Environment": environment,
			"CreatedBy":   "lambda",
		},
		environment,
		ssmKeyCertificateArn,
	)
	dns.CertificateAuthority = letsencrypt.NewCertificateAuthority()

	lambda.Start(Handler)
}
