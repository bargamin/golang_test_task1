package config

import (
	"fmt"
	"github.com/joeshaw/envdecode"
)

type ServerConfig struct {
	Host string `env:"SERVER_HOST"`
	Port string `env:"SERVER_PORT"`
}

func (c ServerConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func NewServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}

	if err := envdecode.Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
