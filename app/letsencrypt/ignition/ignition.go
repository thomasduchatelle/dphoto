package main

import (
	"flag"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/thomasduchatelle/dphoto/domain/dns"
	"github.com/thomasduchatelle/dphoto/domain/dnsadapters/letsencrypt"
	"github.com/thomasduchatelle/dphoto/domain/dnsadapters/route53_adapter"
)

func main() {
	domain := flag.String("domain", "", "Domain for which request and install the certificate")
	email := flag.String("email", "", "Email to own the SSL certificate")
	environment := flag.String("env", "", "DPhoto environment, used for SSM name and tags")
	force := flag.Bool("force", false, "force re-generating the certificate (to be used in case of domain change)")

	flag.Parse()

	if *domain == "" || *email == "" || *environment == "" {
		flag.PrintDefaults()
		return
	}

	sess := session.Must(session.NewSession())
	dns.CertificateManager = route53_adapter.NewCertificateManager(sess, map[string]string{
		"Application": "dphoto-app",
		"Environment": *environment,
		"CreatedBy":   "lambda",
	}, *environment)
	dns.CertificateAuthority = letsencrypt.NewCertificateAuthority()

	err := dns.RenewCertificate(*email, *domain, *force)
	if err != nil {
		panic(err)
	}
}
