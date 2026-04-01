package providers

import (
	"context"
	"fmt"
	"strings"
)

import (
	"github.com/mdnmdn/dns-magic/internal/config"
	"github.com/mdnmdn/dns-magic/internal/dns"
)

type Provider interface {
	ListDomains(ctx context.Context, opts DomainListOptions) ([]dns.DomainSummary, error)
	ListRecords(ctx context.Context, opts RecordListOptions) ([]dns.Record, error)
	SearchRecords(ctx context.Context, query string, recordType string) ([]dns.Record, error)
	AddRecord(ctx context.Context, domain string, record dns.Record) error
	UpdateRecord(ctx context.Context, domain string, record dns.Record) error
	RemoveRecord(ctx context.Context, domain string, name string, recordType string) error
}

type DomainListOptions struct {
	ShopperID     string
	Statuses      []string
	StatusGroups  []string
	Includes      []string
	Limit         int
	Marker        string
	ModifiedSince string
}

type RecordListOptions struct {
	Domain    string
	Type      string
	Name      string
	ShopperID string
	Offset    int
	Limit     int
}

func ResolveProvider(cfg config.Config, alias string) (string, config.ProviderConfig, error) {
	if alias == "" {
		alias = cfg.DefaultProvider
	}

	providerCfg, err := cfg.Provider(alias)
	if err != nil {
		return "", config.ProviderConfig{}, err
	}

	return alias, providerCfg, nil
}

func NormalizeRecordType(recordType string) string {
	return strings.ToUpper(strings.TrimSpace(recordType))
}

type Factory func(alias string, cfg config.ProviderConfig, defaultTimeout config.Duration) (Provider, error)

type Registry struct {
	factories map[string]Factory
}

func NewRegistry() *Registry {
	return &Registry{
		factories: map[string]Factory{},
	}
}

func (r *Registry) Register(providerType string, factory Factory) {
	r.factories[providerType] = factory
}

func (r *Registry) Build(alias string, cfg config.Config) (Provider, error) {
	resolvedAlias, providerCfg, err := ResolveProvider(cfg, alias)
	if err != nil {
		return nil, err
	}

	factory, ok := r.factories[providerCfg.Type]
	if !ok {
		return nil, fmt.Errorf("provider type %q is not registered", providerCfg.Type)
	}

	return factory(resolvedAlias, providerCfg, cfg.RequestTimeout)
}
