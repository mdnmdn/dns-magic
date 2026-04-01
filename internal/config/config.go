package config

import (
	"fmt"
	"time"
)

type Config struct {
	DefaultProvider string                    `toml:"default_provider" json:"default_provider" yaml:"default_provider" toon:"default_provider"`
	RequestTimeout  Duration                  `toml:"request_timeout" json:"request_timeout" yaml:"request_timeout" toon:"request_timeout"`
	Providers       map[string]ProviderConfig `toml:"providers" json:"providers" yaml:"providers" toon:"providers"`
}

type ProviderConfig struct {
	Type      string   `toml:"type" json:"type" yaml:"type" toon:"type"`
	APIKey    string   `toml:"api_key" json:"api_key" yaml:"api_key" toon:"api_key"`
	APISecret string   `toml:"api_secret" json:"api_secret" yaml:"api_secret" toon:"api_secret"`
	APIToken  string   `toml:"api_token" json:"api_token" yaml:"api_token" toon:"api_token"`
	BaseURL   string   `toml:"base_url" json:"base_url" yaml:"base_url" toon:"base_url"`
	ShopperID string   `toml:"shopper_id" json:"shopper_id" yaml:"shopper_id" toon:"shopper_id"`
	OTE       bool     `toml:"ote" json:"ote" yaml:"ote" toon:"ote"`
	Timeout   Duration `toml:"timeout" json:"timeout" yaml:"timeout" toon:"timeout"`
}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		d.Duration = 0
		return nil
	}

	parsed, err := time.ParseDuration(string(text))
	if err != nil {
		return fmt.Errorf("parse duration %q: %w", string(text), err)
	}

	d.Duration = parsed
	return nil
}

func (d Duration) MarshalText() ([]byte, error) {
	if d.Duration == 0 {
		return []byte("0s"), nil
	}

	return []byte(d.Duration.String()), nil
}

func (c Config) Provider(alias string) (ProviderConfig, error) {
	provider, ok := c.Providers[alias]
	if !ok {
		return ProviderConfig{}, fmt.Errorf("provider %q not found in config", alias)
	}

	if provider.Type == "" {
		return ProviderConfig{}, fmt.Errorf("provider %q is missing type", alias)
	}

	return provider, nil
}
