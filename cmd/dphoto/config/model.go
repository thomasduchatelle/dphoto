package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/spf13/viper"
)

type Config interface {
	Get(key string) interface{}
	GetString(key string) string
	GetStringOrDefault(key string, defaultValue string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetIntOrDefault(key string, defaultValue int) int
	GetAWSV2Config() aws.Config
}

type viperConfig struct {
	*viper.Viper
	awsConfig aws.Config
}

func (v *viperConfig) GetAWSV2Config() aws.Config {
	return v.awsConfig
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
