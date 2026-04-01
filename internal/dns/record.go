package dns

type Record struct {
	Domain   string `json:"domain,omitempty" yaml:"domain,omitempty" toon:"domain,omitempty"`
	Name     string `json:"name" yaml:"name" toon:"name"`
	Type     string `json:"type" yaml:"type" toon:"type"`
	Data     string `json:"data" yaml:"data" toon:"data"`
	TTL      int    `json:"ttl,omitempty" yaml:"ttl,omitempty" toon:"ttl,omitempty"`
	Priority int    `json:"priority,omitempty" yaml:"priority,omitempty" toon:"priority,omitempty"`
	Port     int    `json:"port,omitempty" yaml:"port,omitempty" toon:"port,omitempty"`
	Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty" toon:"protocol,omitempty"`
	Service  string `json:"service,omitempty" yaml:"service,omitempty" toon:"service,omitempty"`
	Weight   int    `json:"weight,omitempty" yaml:"weight,omitempty" toon:"weight,omitempty"`
}
