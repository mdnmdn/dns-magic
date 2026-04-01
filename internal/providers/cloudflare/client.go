package cloudflare

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mdnmdn/dns-magic/internal/config"
	"github.com/mdnmdn/dns-magic/internal/dns"
	"github.com/mdnmdn/dns-magic/internal/providers"
)

const (
	defaultBaseURL = "https://api.cloudflare.com/client/v4"
)

type Client struct {
	alias      string
	apiToken   string
	baseURL    string
	httpClient *http.Client
}

var _ providers.Provider = (*Client)(nil)

func New(alias string, cfg config.ProviderConfig, defaultTimeout config.Duration) (providers.Provider, error) {
	if cfg.APIToken == "" {
		return nil, fmt.Errorf("cloudflare provider %q requires api_token", alias)
	}

	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	timeout := cfg.Timeout.Duration
	if timeout == 0 {
		timeout = defaultTimeout.Duration
	}
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	return &Client{
		alias:    alias,
		apiToken: cfg.APIToken,
		baseURL:  baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (c *Client) ListDomains(ctx context.Context, opts providers.DomainListOptions) ([]dns.DomainSummary, error) {
	query := url.Values{}
	if len(opts.Statuses) > 0 {
		query.Set("status", opts.Statuses[0])
	}
	if opts.Limit > 0 {
		query.Set("per_page", strconv.Itoa(opts.Limit))
	}
	// Cloudflare uses 'page' and 'per_page'. Mapping 'marker' to 'page' if numeric, otherwise just first page.
	page := 1
	if opts.Marker != "" {
		if p, err := strconv.Atoi(opts.Marker); err == nil {
			page = p
		}
	}
	query.Set("page", strconv.Itoa(page))

	var response cloudflareZonesResponse
	if err := c.get(ctx, "/zones", query, &response); err != nil {
		return nil, err
	}

	domains := make([]dns.DomainSummary, 0, len(response.Result))
	for _, item := range response.Result {
		domains = append(domains, item.toDNS())
	}

	return domains, nil
}

func (c *Client) ListRecords(ctx context.Context, opts providers.RecordListOptions) ([]dns.Record, error) {
	if strings.TrimSpace(opts.Domain) == "" {
		return nil, fmt.Errorf("domain is required")
	}

	zoneID, err := c.getZoneID(ctx, opts.Domain)
	if err != nil {
		return nil, fmt.Errorf("resolve zone id for %q: %w", opts.Domain, err)
	}

	query := url.Values{}
	if opts.Type != "" {
		query.Set("type", providers.NormalizeRecordType(opts.Type))
	}
	if opts.Name != "" {
		query.Set("name", opts.Name)
	}
	if opts.Limit > 0 {
		query.Set("per_page", strconv.Itoa(opts.Limit))
	}
	// Mapping offset to page/per_page is tricky, let's keep it simple for now
	if opts.Offset > 0 && opts.Limit > 0 {
		page := (opts.Offset / opts.Limit) + 1
		query.Set("page", strconv.Itoa(page))
	}

	path := fmt.Sprintf("/zones/%s/dns_records", zoneID)
	var response cloudflareRecordsResponse
	if err := c.get(ctx, path, query, &response); err != nil {
		return nil, err
	}

	records := make([]dns.Record, 0, len(response.Result))
	for _, item := range response.Result {
		records = append(records, item.toDNS(opts.Domain))
	}

	return records, nil
}

func (c *Client) SearchRecords(ctx context.Context, query string, recordType string) ([]dns.Record, error) {
	return nil, fmt.Errorf("search is not implemented for provider %q", c.alias)
}

func (c *Client) AddRecord(ctx context.Context, domain string, record dns.Record) error {
	return fmt.Errorf("add is not implemented for provider %q", c.alias)
}

func (c *Client) UpdateRecord(ctx context.Context, domain string, record dns.Record) error {
	return fmt.Errorf("update is not implemented for provider %q", c.alias)
}

func (c *Client) RemoveRecord(ctx context.Context, domain string, name string, recordType string) error {
	return fmt.Errorf("remove is not implemented for provider %q", c.alias)
}

func (c *Client) getZoneID(ctx context.Context, domain string) (string, error) {
	query := url.Values{}
	query.Set("name", domain)

	var response cloudflareZonesResponse
	if err := c.get(ctx, "/zones", query, &response); err != nil {
		return "", err
	}

	for _, zone := range response.Result {
		if zone.Name == domain {
			return zone.ID, nil
		}
	}

	return "", fmt.Errorf("zone not found for domain %q", domain)
}

func (c *Client) get(ctx context.Context, path string, query url.Values, out any) error {
	requestURL := c.baseURL + path
	if encoded := query.Encode(); encoded != "" {
		requestURL += "?" + encoded
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 32*1024))
		return fmt.Errorf("cloudflare %s %s: status %d: %s", req.Method, path, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if out == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode %s response: %w", path, err)
	}

	return nil
}

type cloudflareZonesResponse struct {
	Result []cloudflareZone `json:"result"`
}

type cloudflareZone struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Paused      bool      `json:"paused"`
	Type        string    `json:"type"`
	NameServers []string  `json:"name_servers"`
	ModifiedOn  time.Time `json:"modified_on"`
	CreatedOn   time.Time `json:"created_on"`
}

func (z cloudflareZone) toDNS() dns.DomainSummary {
	return dns.DomainSummary{
		Domain:      z.Name,
		Status:      z.Status,
		CreatedAt:   z.CreatedOn,
		NameServers: z.NameServers,
	}
}

type cloudflareRecordsResponse struct {
	Result []cloudflareRecord `json:"result"`
}

type cloudflareRecord struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	TTL       int    `json:"ttl"`
	Priority  int    `json:"priority"`
	Proxied   bool   `json:"proxied"`
	CreatedOn string `json:"created_on"`
}

func (r cloudflareRecord) toDNS(domain string) dns.Record {
	return dns.Record{
		Domain:   domain,
		Name:     r.Name,
		Type:     r.Type,
		Data:     r.Content,
		TTL:      r.TTL,
		Priority: r.Priority,
	}
}
