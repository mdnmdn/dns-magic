package config

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

func Load(path string) (Config, error) {
	var cfg Config

	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return Config{}, err
	}

	if cfg.RequestTimeout.Duration == 0 {
		cfg.RequestTimeout.Duration = 10 * time.Second
	}

	if cfg.DefaultProvider == "" && len(cfg.Providers) == 1 {
		for alias := range cfg.Providers {
			cfg.DefaultProvider = alias
		}
	}

	if cfg.DefaultProvider == "" {
		return Config{}, fmt.Errorf("default_provider is required")
	}

	return cfg, nil
}
