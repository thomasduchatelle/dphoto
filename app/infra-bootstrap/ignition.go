package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/dns"
	"github.com/thomasduchatelle/dphoto/domain/dnsadapters/letsencrypt"
	"github.com/thomasduchatelle/dphoto/domain/dnsadapters/route53_adapter"
)

func main() {
	domain := flag.String("domain", "", "Domain for which request and install the certificate")
	email := flag.String("email", "", "Email to own the SSL certificate")
	environment := flag.String("env", "", "DPhoto environment, used for SSM name and tags")
	force := flag.Bool("force", false, "force re-generating the certificate (to be used in case of domain change)")
	googleClientId := flag.String("google-client-id", "", "Google client ID managed in https://console.developers.google.com/apis/credentials (optional if already set)")

	flag.Parse()

	sess := session.Must(session.NewSession())

	err := updateSSMParameter(sess, *googleClientId, *environment, *domain)
	if err != nil {
		panic(err)
	}

	err = createAndInstallSSL(sess, domain, email, environment, force)
	if err != nil {
		panic(err)
	}
}

func updateSSMParameter(sess *session.Session, googleClientId, environment, domain string) error {
	ssmClient := ssm.New(sess)

	ssmName := fmt.Sprintf("/dphoto/%s/googleLogin/clientId", environment)
	value, err := readSSM(ssmClient, ssmName)
	if err != nil {
		return err
	}

	switch {
	case value == "" && googleClientId == "":
		return errors.Errorf("google-client-id is required, SSM key '%s' is not set.", ssmName)

	case googleClientId == "":
		log.WithField("Domain", domain).Infof("skipping Google client ID: already set to %s", value)

	case value == googleClientId:
		log.WithField("Domain", domain).Infof("skipping Google client ID: unchanged.")

	default:
		log.WithField("Domain", domain).Infof("Updating google-clientid from '%s' to '%s'", value, googleClientId)
		_, err = ssmClient.PutParameter(&ssm.PutParameterInput{
			DataType:  aws.String("text"),
			Name:      &ssmName,
			Type:      aws.String("String"),
			Value:     &googleClientId,
			Overwrite: aws.Bool(value != ""),
		})
		return errors.Wrapf(err, "set SSM paramater to %s failed", googleClientId)
	}

	return nil
}

func readSSM(ssmClient *ssm.SSM, key string) (string, error) {
	parameter, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(key),
	})
	if err != nil {
		if awsError, ok := err.(awserr.Error); ok && awsError.Code() == ssm.ErrCodeParameterNotFound {
			return "", nil
		}

		return "", err
	}

	return *parameter.Parameter.Value, nil
}

func createAndInstallSSL(sess *session.Session, domain *string, email *string, environment *string, force *bool) error {
	if *domain == "" || *email == "" || *environment == "" {
		flag.PrintDefaults()
		return nil
	}

	dns.CertificateManager = route53_adapter.NewCertificateManager(sess, map[string]string{
		"Application": "dphoto-app",
		"Environment": *environment,
		"CreatedBy":   "lambda",
	}, *environment)
	dns.CertificateAuthority = letsencrypt.NewCertificateAuthority()

	return dns.RenewCertificate(*email, *domain, *force)
}
