package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	sessionv1 "github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"
)

type Config interface {
	Get(key string) interface{}
	GetString(key string) string
	GetStringOrDefault(key string, defaultValue string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetIntOrDefault(key string, defaultValue int) int
	GetAWSSession() *sessionv1.Session
	GetAWSV2Config() aws.Config
}

type viperConfig struct {
	*viper.Viper
	awsSession *sessionv1.Session
	awsConfig  aws.Config
}

func (v *viperConfig) GetAWSSession() *sessionv1.Session {
	return v.awsSession
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
