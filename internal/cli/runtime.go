package cli

import (
	"context"
	"fmt"

	"github.com/mdnmdn/dns-magic/internal/config"
	"github.com/mdnmdn/dns-magic/internal/providers"
	"github.com/mdnmdn/dns-magic/internal/providers/cloudflare"
	"github.com/mdnmdn/dns-magic/internal/providers/godaddy"
)

type Runtime struct {
	ConfigPath string
	registry   *providers.Registry
}

func NewRuntime() *Runtime {
	registry := providers.NewRegistry()
	registry.Register("godaddy", godaddy.New)
	registry.Register("cloudflare", cloudflare.New)

	return &Runtime{
		ConfigPath: "config.toml",
		registry:   registry,
	}
}

func (r *Runtime) LoadConfig() (config.Config, error) {
	return config.Load(r.ConfigPath)
}

func (r *Runtime) Provider(alias string) (providers.Provider, error) {
	cfg, err := r.LoadConfig()
	if err != nil {
		return nil, err
	}

	return r.registry.Build(alias, cfg)
}

func (r *Runtime) MustProvider(ctx context.Context, alias string) (providers.Provider, config.Config, error) {
	cfg, err := r.LoadConfig()
	if err != nil {
		return nil, config.Config{}, err
	}

	provider, err := r.registry.Build(alias, cfg)
	if err != nil {
		return nil, config.Config{}, err
	}

	_ = ctx
	return provider, cfg, nil
}

func requireOutput(format string) error {
	switch format {
	case "", "table", "json", "yaml", "markdown", "toon":
		return nil
	default:
		return fmt.Errorf("unsupported --output %q", format)
	}
}
