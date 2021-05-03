// Package config create AWS session and distribute it to any package interested once all config files has been loaded.
package config

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Listener function is called once config is loaded
type Listener func(Config)

var listeners []Listener

func init() {
	viper.SetConfigName("dphoto")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.dphoto")
	viper.AddConfigPath("/etc/dphoto/")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}

// Listen registers a Listener that will be invoked when configuration will be provided.
func Listen(listener Listener) {
	listeners = append(listeners, listener)
}

// Connect must be called by main function, it dispatches the config to all components requiring it.
func Connect() {
	// use explicit config to avoid creating unwanted environment
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(viper.GetString("aws.key"), viper.GetString("aws.secret"), viper.GetString("aws.token")),
		Endpoint:    awsString(viper.GetString("aws.endpoint")),
		Region:      aws.String(viper.GetString("aws.region")),
	}))

	config := &viperConfig{
		Viper:      viper.GetViper(),
		awsSession: sess,
	}

	for _, l := range listeners {
		l(config)
	}
	log.Debugf("Config > %d adapters connected", len(listeners))
}

func awsString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
