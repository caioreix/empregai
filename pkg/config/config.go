package config

import (
	"time"

	"github.com/spf13/viper"
)

// Server config
type Server struct {
	SSL          bool
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	JWTSecret    string
}

// Session config
type Session struct {
	BasePrefix string
	Name       string
	Duration   time.Duration
}

// Cookie config
type Cookie struct {
	Name     string
	Domain   string
	MaxAge   int
	Secure   bool
	Path     string
	HTTPOnly bool
}

// Postgresql config
type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  bool
	Driver   string
}

// Redis config
type Redis struct {
	Addr         string
	MinIdleConns int
	PoolSize     int
	PoolTimeout  time.Duration
	Password     string
	DB           int
}

// Config centralizer
type Config struct {
	Server   Server
	Session  Session
	Cookie   Cookie
	Postgres Postgres
	Redis    Redis
}

func LoadConfig(fileName string, filePath string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(fileName)
	v.AddConfigPath(filePath)
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (c *Config, err error) {
	err = v.Unmarshal(&c)
	return
}
