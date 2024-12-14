package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP     HTTPConfig
	Database DatabaseConfig
	Auth     AuthConfig
	SMTP     SMTPConfig
	UseMockEmail bool `env:"USE_MOCK_EMAIL" envDefault:"true"`
}

type HTTPConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}


type AuthConfig struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		HTTP: HTTPConfig{
			Host:         viper.GetString("HTTP_HOST"),
			Port:         viper.GetString("HTTP_PORT"),
			ReadTimeout:  viper.GetDuration("HTTP_READ_TIMEOUT"),
			WriteTimeout: viper.GetDuration("HTTP_WRITE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Auth: AuthConfig{
			JWTSecret:       viper.GetString("JWT_SECRET"),
			AccessTokenTTL:  viper.GetDuration("ACCESS_TOKEN_TTL"),
			RefreshTokenTTL: viper.GetDuration("REFRESH_TOKEN_TTL"),
		},
		SMTP: SMTPConfig{
			Host:     viper.GetString("SMTP_HOST"),
			Port:     viper.GetString("SMTP_PORT"),
			Username: viper.GetString("SMTP_USERNAME"),
			Password: viper.GetString("SMTP_PASSWORD"),
			From:     viper.GetString("SMTP_FROM"),
		},
		UseMockEmail: viper.GetBool("USE_MOCK_EMAIL"),
	}, nil
}
