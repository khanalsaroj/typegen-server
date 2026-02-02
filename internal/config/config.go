package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Security SecurityConfig
}

type AppConfig struct {
	Name        string
	Version     string
	Environment string
}

type ServerConfig struct {
	Port            int
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int
}

type DatabaseConfig struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
	Filepath        string
}

type SecurityConfig struct {
	RateLimitEnabled bool
	RateLimitRPS     int
	CORSAllowOrigins []string
	DbEncryptionKey  string
}

func Load() (*Config, error) {
	v := viper.New()

	env := getEnv("APP_ENV", "dev")

	v.SetConfigName(fmt.Sprintf(".env.%s", env))
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("./")
	v.AddConfigPath("../")
	v.AddConfigPath("../../")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	setDefaults(v)

	cfg := &Config{
		App: AppConfig{
			Name:        v.GetString("APP_NAME"),
			Version:     v.GetString("APP_VERSION"),
			Environment: env,
		},
		Server: ServerConfig{
			Port:            v.GetInt("SERVER_PORT"),
			ReadTimeout:     v.GetInt("SERVER_READ_TIMEOUT"),
			WriteTimeout:    v.GetInt("SERVER_WRITE_TIMEOUT"),
			ShutdownTimeout: v.GetInt("SERVER_SHUTDOWN_TIMEOUT"),
		},
		Database: DatabaseConfig{
			MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
			MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
			ConnMaxLifetime: v.GetInt("DB_CONN_MAX_LIFETIME"),
			Filepath:        v.GetString("DB_FILE_PATH"),
		},
		Security: SecurityConfig{
			RateLimitEnabled: v.GetBool("RATE_LIMIT_ENABLED"),
			RateLimitRPS:     v.GetInt("RATE_LIMIT_RPS"),
			CORSAllowOrigins: v.GetStringSlice("CORS_ALLOW_ORIGINS"),
			DbEncryptionKey:  v.GetString("DB_ENCRYPTION_KEY"),
		},
	}
	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_NAME", "github.com/khanalsaroj/typegen-server")
	v.SetDefault("APP_VERSION", "1.0.0")
	v.SetDefault("SERVER_PORT", 8080)
	v.SetDefault("SERVER_READ_TIMEOUT", 10)
	v.SetDefault("SERVER_WRITE_TIMEOUT", 10)
	v.SetDefault("SERVER_SHUTDOWN_TIMEOUT", 10)
	v.SetDefault("DB_MAX_IDLE_CONNS", 10)
	v.SetDefault("DB_MAX_OPEN_CONNS", 100)
	v.SetDefault("DB_CONN_MAX_LIFETIME", 60)
	v.SetDefault("DB_FILE_PATH", "./data/database.db")
	v.SetDefault("RATE_LIMIT_RPS", 100)
	v.SetDefault("CORS_ALLOW_ORIGINS", []string{"*"})
	v.SetDefault("DB_ENCRYPTION_KEY", "9f7c8b2d1a4e6c3f9a0b2c5e7d8f1a2b")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
