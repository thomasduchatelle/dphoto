package letsencrypt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/route53"
	"github.com/go-acme/lego/v4/registration"
	"github.com/thomasduchatelle/dphoto/domain/dnsdomain"
)

type legoAdapter struct {
}

type LegoUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *LegoUser) GetEmail() string {
	return u.Email
}
func (u LegoUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *LegoUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func NewCertificateAuthority() dnsdomain.CertificateAuthority {
	return &legoAdapter{}
}

func (l *legoAdapter) RequestCertificate(email, domain string) (*dnsdomain.CompleteCertificate, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	config := lego.NewConfig(&LegoUser{
		Email: email,
		key:   privateKey,
	})
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	provider, err := route53.NewDNSProvider()
	if err != nil {
		return nil, err
	}

	err = client.Challenge.SetDNS01Provider(provider)
	if err != nil {
		return nil, err
	}

	_, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}

	cer, err := client.Certificate.Obtain(certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  false,
	})
	if err != nil {
		return nil, err
	}

	return &dnsdomain.CompleteCertificate{
		Certificate: cer.Certificate,
		Chain:       cer.IssuerCertificate,
		PrivateKey:  cer.PrivateKey,
	}, nil
}
