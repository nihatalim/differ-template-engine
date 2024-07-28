package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

type ApplicationConfig struct {
	Port       string
	Postgresql Postgresql
	Clients    Clients
}

type APIConfig struct {
	Host    string
	Timeout time.Duration
	Retry   int
}

type Clients struct {
	Nodiffer APIConfig
}

type Postgresql struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type Secret struct {
	Username string
	Password string
}

func getEnv() string {
	env := os.Getenv("GO_ENV")
	if env == "" {
		return "stage"
	}

	return env
}

func New() (*ApplicationConfig, error) {
	env := getEnv()

	cfg := &ApplicationConfig{}

	v := viper.New()
	v.SetConfigFile(fmt.Sprintf("resources/%s.yml", env))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	sub := v.Sub(env)

	if err := sub.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if env == "dev" {
		cfg.Postgresql.Host = os.Getenv("POSTGRESQL_HOST")
		cfg.Postgresql.Port = os.Getenv("POSTGRESQL_PORT")
		cfg.Postgresql.Database = os.Getenv("POSTGRESQL_DB")
		cfg.Postgresql.Username = os.Getenv("POSTGRESQL_USERNAME")
		cfg.Postgresql.Password = os.Getenv("POSTGRESQL_PASSWORD")
		return cfg, nil
	}

	secret, err := readSecrets(v)
	if err != nil {
		return nil, err
	}

	cfg.Postgresql.Username = secret.Username
	cfg.Postgresql.Password = secret.Password

	return cfg, nil
}

func readSecrets(viperInstance *viper.Viper) (*Secret, error) {
	viperInstance.SetConfigFile("configs/secrets.json")

	if err := viperInstance.ReadInConfig(); err != nil {
		return nil, err
	}

	secret := Secret{}

	if err := viperInstance.Unmarshal(&secret); err != nil {
		return nil, fmt.Errorf("failed to read kafka secret configuration %w", err)
	}

	return &secret, nil
}
