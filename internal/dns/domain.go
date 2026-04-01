package dns

import "time"

type DomainSummary struct {
	Domain      string    `json:"domain" yaml:"domain" toon:"domain"`
	Status      string    `json:"status,omitempty" yaml:"status,omitempty" toon:"status,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty" yaml:"expires_at,omitempty" toon:"expires_at,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty" toon:"created_at,omitempty"`
	Locked      bool      `json:"locked" yaml:"locked" toon:"locked"`
	Privacy     bool      `json:"privacy" yaml:"privacy" toon:"privacy"`
	RenewAuto   bool      `json:"renew_auto" yaml:"renew_auto" toon:"renew_auto"`
	NameServers []string  `json:"name_servers,omitempty" yaml:"name_servers,omitempty" toon:"name_servers,omitempty"`
}
