package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig   `mapstructure:",squash"`
	DB     DatabaseConfig `mapstructure:",squash"`
	Auth   AuthConfig     `mapstructure:",squash"`
}

type AuthConfig struct {
	AccessTokenSecret        string `mapstructure:"ACCESS_TOKEN_SECRET"`
	AccessTokenExpiryMinutes int    `mapstructure:"ACCESS_TOKEN_EXPIRY_MINUTES"`
	RefreshTokenExpiryDays   int    `mapstructure:"REFRESH_TOKEN_EXPIRY_DAYS"`
}

type ServerConfig struct {
	Address string `mapstructure:"ADDRESS"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSLMODE"`
}

func LoadConfig() (Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
