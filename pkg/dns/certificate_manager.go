// Package dns expose functions to renew SSL certificates on AWS
package dns

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/dnsdomain"
	"time"
)

const MinimumExpiryDelay = 30 * 24 * time.Hour

var (
	CertificateManager   dnsdomain.CertificateManager
	CertificateAuthority dnsdomain.CertificateAuthority
)

func RenewCertificate(email, domain string, forced bool) error {
	logCtx := log.WithField("Domain", domain)
	logCtx.Infoln("checking SSL certificate validity...")

	id := ""
	ctx := context.TODO()

	existing, err := CertificateManager.FindCertificate(ctx, domain)
	if err != nil && !errors.Is(err, dnsdomain.CertificateNotFoundError) {
		return errors.Wrapf(err, "find existing certificate failed for domain %s", domain)
	} else if err == nil {
		id = existing.ID
	}

	if forced || existing == nil || existing.Expiry.Before(time.Now().Add(MinimumExpiryDelay)) {
		logCtx.Infoln("Renewing SSL certificate")

		cert, err := CertificateAuthority.RequestCertificate(ctx, email, domain)
		if err != nil {
			return err
		}

		err = CertificateManager.InstallCertificate(ctx, id, *cert)
		if err == nil {
			logCtx.Infoln("Certificate installed.")
		}
		return err
	} else {
		logCtx.WithField("CertARN", existing.ID).Infof("Certificate present and valid until %s", existing.Expiry.Format("02/01/2006 15:04:05"))
	}

	return nil
}
