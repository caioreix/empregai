package config

import (
	"time"

	"github.com/spf13/viper"
)

type Session struct {
	BasePrefix string
	Duration   time.Duration
}

type Config struct {
	Session Session
}

func LoadConfig(fileName string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(fileName)
	v.SetConfigFile(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (c *Config, err error) {
	err = v.Unmarshal(c)
	return
}
