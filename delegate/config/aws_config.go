package config

import "github.com/aws/aws-sdk-go/aws/session"

var AwsConfig *DefaultAwsConfigSupplier

func init() {
	AwsConfig = newDefaultAwsConfigSupplier()
}

func newDefaultAwsConfigSupplier() *DefaultAwsConfigSupplier {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &DefaultAwsConfigSupplier{
		session: sess,
	}
}

type DefaultAwsConfigSupplier struct {
	session *session.Session
}

func (d *DefaultAwsConfigSupplier) GetSession() *session.Session {
	return d.session
}
