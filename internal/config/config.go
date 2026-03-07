package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	Upload   UploadConfig
}

type UploadConfig struct {
	Dir string
}

type ServerConfig struct {
	Port        string
	Host        string
	Environment string
	JWTSecret   string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	Expiration int
}

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GitHubClientID     string
	GitHubClientSecret string
	RedirectURL        string
}

func Load() (*Config, error) {
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	return &Config{
		Server: ServerConfig{
			Port:        viper.GetString("PORT"),
			Host:        viper.GetString("HOST"),
			Environment: viper.GetString("ENVIRONMENT"),
			JWTSecret:   viper.GetString("JWTSECRET"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("JWT_SECRET"),
			Expiration: viper.GetInt("JWT_EXPIRATION_HOURS"),
		},
		OAuth: OAuthConfig{
			GoogleClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
			GoogleClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
			GitHubClientID:     viper.GetString("GITHUB_CLIENT_ID"),
			GitHubClientSecret: viper.GetString("GITHUB_CLIENT_SECRET"),
			RedirectURL:        viper.GetString("OAUTH_REDIRECT_URL"),
		},
		Upload: UploadConfig{
			Dir: viper.GetString("UPLOAD_DIR"),
		},
	}, nil
}

func (d *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}

func (r *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}
