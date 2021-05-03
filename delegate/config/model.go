package config

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"
)

type Config interface {
	Get(key string) interface{}
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetIntOrDefault(key string, defaultValue int) int
	GetAWSSession() *session.Session
}

type viperConfig struct {
	*viper.Viper
	awsSession *session.Session
}

func (v *viperConfig) GetAWSSession() *session.Session {
	return v.awsSession
}

func (v *viperConfig) GetIntOrDefault(key string, defaultValue int) int {
	value := v.Viper.GetInt(key)
	if value == 0 {
		value = defaultValue
	}

	return value
}
