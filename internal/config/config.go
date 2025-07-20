package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"time"
)

type Config_DB struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string
}

type Config_Server struct {
	ServerPort string
}

type Config_Redis struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	TTL           time.Duration
}

type CacheConfig struct {
	TTL time.Duration `yaml:"ttl"`
}

type AuthConfig struct {
	AccessTokenTTL time.Duration `yaml:"accessTokenTTL"`
	JWTSecret      string
}

type AppConfig struct {
	Cache CacheConfig `yaml:"cache"`
	Auth  AuthConfig  `yaml:"auth"`
}

type Config struct {
	Config_DB
	Config_Server
	Config_Redis
	AuthConfig
}

func LoadYaml(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg, err := LoadYaml("configs/main.yml")
	if err != nil {
		return nil, fmt.Errorf("error loading main.yml file: %w", err)
	}

	return &Config{
		Config_DB: Config_DB{
			DBUser:     os.Getenv("DB_USER"),
			DBPassword: os.Getenv("DB_PASSWORD"),
			DBName:     os.Getenv("DB_NAME"),
			DBPort:     os.Getenv("DB_PORT"),
			DBSSLMode:  os.Getenv("DB_SSLMODE"),
		},
		Config_Server: Config_Server{
			ServerPort: os.Getenv("SERVER_PORT"),
		},
		Config_Redis: Config_Redis{
			RedisAddr:     os.Getenv("REDIS_ADDR"),
			RedisPassword: os.Getenv("REDIS_PASSWORD"),
			RedisDB: func() int {
				db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
				return db
			}(),
			TTL: cfg.Cache.TTL,
		},
		AuthConfig: AuthConfig{
			AccessTokenTTL: cfg.Auth.AccessTokenTTL,
			JWTSecret:      os.Getenv("JWT_SECRET"),
		},
	}, nil
}
