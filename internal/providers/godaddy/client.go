package godaddy

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
	defaultBaseURL    = "https://api.godaddy.com"
	defaultOTEBaseURL = "https://api.ote-godaddy.com"
)

type Client struct {
	alias      string
	apiKey     string
	apiSecret  string
	shopperID  string
	baseURL    string
	httpClient *http.Client
}

var _ providers.Provider = (*Client)(nil)

func New(alias string, cfg config.ProviderConfig, defaultTimeout config.Duration) (providers.Provider, error) {
	if cfg.APIKey == "" || cfg.APISecret == "" {
		return nil, fmt.Errorf("godaddy provider %q requires api_key and api_secret", alias)
	}

	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		if cfg.OTE {
			baseURL = defaultOTEBaseURL
		} else {
			baseURL = defaultBaseURL
		}
	}

	timeout := cfg.Timeout.Duration
	if timeout == 0 {
		timeout = defaultTimeout.Duration
	}
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	return &Client{
		alias:     alias,
		apiKey:    cfg.APIKey,
		apiSecret: cfg.APISecret,
		shopperID: cfg.ShopperID,
		baseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (c *Client) ListDomains(ctx context.Context, opts providers.DomainListOptions) ([]dns.DomainSummary, error) {
	query := url.Values{}
	addMulti(query, "statuses", opts.Statuses)
	addMulti(query, "statusGroups", opts.StatusGroups)
	addMulti(query, "includes", opts.Includes)
	if opts.Limit > 0 {
		query.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Marker != "" {
		query.Set("marker", opts.Marker)
	}
	if opts.ModifiedSince != "" {
		query.Set("modifiedDate", opts.ModifiedSince)
	}

	var response []godaddyDomainSummary
	if err := c.get(ctx, "/v1/domains", opts.ShopperID, query, &response); err != nil {
		return nil, err
	}

	domains := make([]dns.DomainSummary, 0, len(response))
	for _, item := range response {
		domains = append(domains, item.toDNS())
	}

	return domains, nil
}

func (c *Client) ListRecords(ctx context.Context, opts providers.RecordListOptions) ([]dns.Record, error) {
	if strings.TrimSpace(opts.Domain) == "" {
		return nil, fmt.Errorf("domain is required")
	}
	if opts.Name != "" && opts.Type == "" {
		return nil, fmt.Errorf("record name filter requires --type")
	}

	recordType := providers.NormalizeRecordType(opts.Type)
	path := fmt.Sprintf("/v1/domains/%s/records", url.PathEscape(opts.Domain))
	if recordType != "" {
		path += "/" + url.PathEscape(recordType)
	}
	if opts.Name != "" {
		path += "/" + url.PathEscape(opts.Name)
	}

	query := url.Values{}
	if opts.Offset > 0 {
		query.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.Limit > 0 {
		query.Set("limit", strconv.Itoa(opts.Limit))
	}

	var response []godaddyRecord
	if err := c.get(ctx, path, opts.ShopperID, query, &response); err != nil {
		return nil, err
	}

	records := make([]dns.Record, 0, len(response))
	for _, item := range response {
		record := item.toDNS(opts.Domain)
		if recordType != "" && record.Type == "" {
			record.Type = recordType
		}
		records = append(records, record)
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

func (c *Client) get(ctx context.Context, path string, shopperID string, query url.Values, out any) error {
	requestURL := c.baseURL + path
	if encoded := query.Encode(); encoded != "" {
		requestURL += "?" + encoded
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "sso-key "+c.apiKey+":"+c.apiSecret)
	req.Header.Set("Accept", "application/json")
	if shopperID = strings.TrimSpace(shopperID); shopperID == "" {
		shopperID = c.shopperID
	}
	if shopperID != "" {
		req.Header.Set("X-Shopper-Id", shopperID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 32*1024))
		return fmt.Errorf("godaddy %s %s: status %d: %s", req.Method, path, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if out == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode %s response: %w", path, err)
	}

	return nil
}

func addMulti(values url.Values, key string, items []string) {
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			values.Add(key, item)
		}
	}
}

type godaddyDomainSummary struct {
	Domain      string   `json:"domain"`
	Status      string   `json:"status"`
	Expires     string   `json:"expires"`
	CreatedAt   string   `json:"createdAt"`
	Locked      bool     `json:"locked"`
	Privacy     bool     `json:"privacy"`
	RenewAuto   bool     `json:"renewAuto"`
	NameServers []string `json:"nameServers"`
}

func (d godaddyDomainSummary) toDNS() dns.DomainSummary {
	return dns.DomainSummary{
		Domain:      d.Domain,
		Status:      d.Status,
		ExpiresAt:   parseTime(d.Expires),
		CreatedAt:   parseTime(d.CreatedAt),
		Locked:      d.Locked,
		Privacy:     d.Privacy,
		RenewAuto:   d.RenewAuto,
		NameServers: d.NameServers,
	}
}

type godaddyRecord struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Data     string `json:"data"`
	TTL      int    `json:"ttl"`
	Priority int    `json:"priority"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Service  string `json:"service"`
	Weight   int    `json:"weight"`
}

func (r godaddyRecord) toDNS(domain string) dns.Record {
	return dns.Record{
		Domain:   domain,
		Name:     r.Name,
		Type:     r.Type,
		Data:     r.Data,
		TTL:      r.TTL,
		Priority: r.Priority,
		Port:     r.Port,
		Protocol: r.Protocol,
		Service:  r.Service,
		Weight:   r.Weight,
	}
}

func parseTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err == nil {
		return parsed
	}

	return time.Time{}
}
