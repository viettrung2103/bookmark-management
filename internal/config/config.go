package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/viettrung2103/bookmark-management/internal/service"
)

type Config struct {
	AppPort     string `default:"8080" envconfig:"APP_PORT"`
	ServiceName string `default:"bookmark_service" envconfig:"SERVICE_NAME""`
	InstanceId  string `default:""  envconfig:"INSTANCE_ID"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := envconfig.Process("api", cfg)
	if err != nil {
		return nil, err
	}

	if cfg.InstanceId == "" {
		genIdSvc := service.NewGenId()
		cfg.InstanceId = genIdSvc.GenerateId()
	}

	return cfg, err
}
