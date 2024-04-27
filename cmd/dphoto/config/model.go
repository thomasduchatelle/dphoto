package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/spf13/viper"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/awsfactory"
)

type Config interface {
	Get(key string) interface{}
	GetString(key string) string
	GetStringOrDefault(key string, defaultValue string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetIntOrDefault(key string, defaultValue int) int
	GetAWSV2Config() aws.Config
	GetAWSFactory() *awsfactory.AWSFactory
}

type viperConfig struct {
	*viper.Viper
	AWSFactory *awsfactory.AWSFactory
}

func (v *viperConfig) GetAWSV2Config() aws.Config {
	return v.AWSFactory.Cfg
}

func (v *viperConfig) GetAWSFactory() *awsfactory.AWSFactory {
	return v.AWSFactory
}

func (v *viperConfig) GetStringOrDefault(key string, defaultValue string) string {
	value := v.Viper.GetString(key)
	if value == "" {
		value = defaultValue
	}

	return value
}

func (v *viperConfig) GetIntOrDefault(key string, defaultValue int) int {
	value := v.Viper.GetInt(key)
	if value == 0 {
		value = defaultValue
	}

	return value
}
