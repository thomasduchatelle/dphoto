// Package config create AWS session and distribute it (alongside other config value) to any package interested.
package config

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awsv1 "github.com/aws/aws-sdk-go/aws"
	credentialsv1 "github.com/aws/aws-sdk-go/aws/credentials"
	sessionv1 "github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"os"
	"path"
)

// Listener function is called once config is loaded
type Listener func(Config)

var (
	ForcedConfigFile string // ForcedConfigFile is the path to the file of the config to use (instead of defaulting to ./dphoto.yml, $HOME/.dphoto/dphoto.yml, ...)
	Environment      string // Environment is used as suffix for the config file name.
	listeners        []Listener
	config           *viperConfig
)

// Listen registers a Listener that will be invoked when configuration will be provided.
func Listen(listener Listener) {
	listeners = append(listeners, listener)
	if config != nil {
		listener(config)
	}
}

// Connect must be called by main function, it dispatches the config to all components requiring it. Set ignite to TRUE to connect to AWS (required for most commands)
func Connect(ignite, createConfigIfNotExist bool) error {
	defaultConfigFile := ForcedConfigFile
	if ForcedConfigFile == "" {
		configFileName := "dphoto"
		if Environment != "" {
			configFileName = fmt.Sprintf("%s-%s", configFileName, Environment)
		}
		viper.SetConfigName(configFileName)
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.dphoto")
		viper.AddConfigPath("/etc/dphoto/")

		defaultConfigFile = os.ExpandEnv("$HOME/.dphoto/") + configFileName + ".yaml"
	} else {
		viper.SetConfigFile(ForcedConfigFile)
	}

	err := viper.ReadInConfig()
	if err != nil {
		if _, isFileNotFound := err.(viper.ConfigFileNotFoundError); isFileNotFound && createConfigIfNotExist {
			printer.Info("Creating default configuration file: %s", defaultConfigFile)
			err = os.MkdirAll(path.Dir(defaultConfigFile), 0600)
			if err != nil {
				return errors.Wrapf(err, "can't create directory for the config file %s", defaultConfigFile)
			}

			_, err = os.Create(defaultConfigFile)
			if err != nil {
				return err
			}

			// read config is re-run to find the now-created file
			return viper.ReadInConfig()

		} else {
			return err
		}
	}

	if ignite {
		// use explicit config to avoid creating unwanted environment
		sess := sessionv1.Must(sessionv1.NewSession(
			&awsv1.Config{
				Endpoint:         awsString(viper.GetString("aws.endpoint")),
				Region:           awsv1.String(viper.GetString("aws.region")),
				S3ForcePathStyle: awsv1.Bool(true), // prevent using S3 bucket name in HOST (incompatible with localstack on macOS)
			},
			awsCredentialsConfig(),
		))
		cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
			awsconfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				endpoint := viper.GetString("aws.endpoint")
				region = viper.GetString("aws.region")
				if endpoint != "" {
					return aws.Endpoint{
						PartitionID:   "aws",
						URL:           endpoint,
						SigningRegion: region,
					}, nil
				}

				// returning EndpointNotFoundError will allow the service to fallback to its default resolution
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			})),
		)
		if err != nil {
			return err
		}

		config = &viperConfig{
			Viper:      viper.GetViper(),
			awsSession: sess,
			awsConfig:  cfg,
		}

		for _, listener := range listeners {
			listener(config)
		}
		log.Debugf("Config > %d adapters connected", len(listeners))
	}

	return nil
}

func awsCredentialsConfig() *awsv1.Config {
	if viper.GetString("aws.key") != "" || viper.GetString("aws.secret") != "" {
		log.Info("AWS Credentials > using static key")
		return &awsv1.Config{
			Credentials: credentialsv1.NewStaticCredentials(viper.GetString("aws.key"), viper.GetString("aws.secret"), viper.GetString("aws.token")),
		}
	}

	log.Info("AWS Credentials > using default config")
	return new(awsv1.Config)
}

func awsString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
