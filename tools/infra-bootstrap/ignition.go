package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/dns"
	"github.com/thomasduchatelle/dphoto/pkg/dnsadapters/letsencrypt"
	"github.com/thomasduchatelle/dphoto/pkg/dnsadapters/route53_adapter"
)

func main() {
	domain := flag.String("domain", "", "Domain for which request and install the certificate")
	email := flag.String("email", "", "Email to own the SSL certificate")
	environment := flag.String("env", "", "DPhoto environment, used for SSM name and tags")
	force := flag.Bool("force", false, "force re-generating the certificate (to be used in case of domain change)")
	googleClientId := flag.String("google-client-id", "", "Google client ID managed in https://console.developers.google.com/apis/credentials (optional if already set)")

	flag.Parse()

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(errors.Wrapf(err, "invalid credetials for ignition"))
	}

	err = updateSSMParameter(ctx, cfg, *googleClientId, *environment, *domain)
	if err != nil {
		panic(err)
	}

	err = createAndInstallSSL(cfg, domain, email, environment, force)
	if err != nil {
		panic(err)
	}
}

func updateSSMParameter(ctx context.Context, cfg aws.Config, googleClientId, environment, domain string) error {
	ssmClient := ssm.NewFromConfig(cfg)

	ssmName := fmt.Sprintf("/dphoto/%s/googleLogin/clientId", environment)
	value, err := readSSM(ctx, ssmClient, ssmName)
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
		_, err = ssmClient.PutParameter(ctx, &ssm.PutParameterInput{
			DataType:  aws.String("text"),
			Name:      &ssmName,
			Type:      types.ParameterTypeString,
			Value:     &googleClientId,
			Overwrite: aws.Bool(value != ""),
		})
		return errors.Wrapf(err, "set SSM paramater to %s failed", googleClientId)
	}

	return nil
}

func readSSM(ctx context.Context, ssmClient *ssm.Client, key string) (string, error) {
	parameter, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(key),
	})

	var parameterNotFound *types.ParameterNotFound
	if errors.As(err, &parameterNotFound) {
		return "", nil

	} else if err != nil {
		return "", err
	}

	return *parameter.Parameter.Value, nil
}

func createAndInstallSSL(cfg aws.Config, domain *string, email *string, environment *string, force *bool) error {
	if *domain == "" || *email == "" || *environment == "" {
		flag.PrintDefaults()
		return nil
	}

	dns.CertificateManager = route53_adapter.NewCertificateManager(cfg, map[string]string{
		"Application": "dphoto-app",
		"Environment": *environment,
		"CreatedBy":   "lambda",
	}, *environment, "")
	dns.CertificateAuthority = letsencrypt.NewCertificateAuthority()

	return dns.RenewCertificate(*email, *domain, *force)
}
