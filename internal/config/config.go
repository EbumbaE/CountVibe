package config

import (
	"github.com/EbumbaE/CountVibe/internal/certificate"
	"github.com/EbumbaE/CountVibe/internal/server"
	"github.com/EbumbaE/CountVibe/internal/session"
	"github.com/EbumbaE/CountVibe/internal/storage/psql"

	"github.com/spf13/viper"
)

type Config struct {
	Pages       map[string]string  `mapstructure:"pages"`
	Server      server.Config      `mapstructure:"server"`
	Session     session.Config     `mapstructure:"session"`
	Certificate certificate.Config `mapstructure:"certificate"`
	Database    psql.Config        `mapstructure:"psql"`
}

func NewConfig() (cfg Config, err error) {
	viper.AddConfigPath("../../internal/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	err = viper.Unmarshal(&cfg)

	return cfg, err
}
