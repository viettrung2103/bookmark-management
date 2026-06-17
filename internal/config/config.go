package config

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

// Config struct for configuration
type Config struct {
	AppPort     string `default:"8080" envconfig:"APP_PORT"`
	ServiceName string `default:"bookmark_service" envconfig:"SERVICE_NAME"`
	InstanceId  string `default:""  envconfig:"INSTANCE_ID"`
	Hostname    string `default:"localhost:8080" envconfig:"APP_HOSTNAME"`
}

// NewConfig creates a new config
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := envconfig.Process("api", cfg)
	if err != nil {
		return nil, err
	}

	if cfg.InstanceId == "" {
		log.Error().Err(err).Str("from", "config.Newconfig").Msg("failed to create config")

		cfg.InstanceId = uuid.New().String()
	}

	return cfg, err
}
