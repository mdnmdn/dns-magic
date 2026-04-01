package godaddy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mdnmdn/dns-magic/internal/config"
	"github.com/mdnmdn/dns-magic/internal/providers"
)

func TestListDomainsAddsDelegatedHeaderAndQuery(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "sso-key key:secret" {
			t.Fatalf("unexpected authorization header: %s", got)
		}
		if got := r.Header.Get("X-Shopper-Id"); got != "shopper-123" {
			t.Fatalf("unexpected shopper id: %s", got)
		}
		if got := r.URL.Query().Get("limit"); got != "25" {
			t.Fatalf("unexpected limit: %s", got)
		}
		if got := r.URL.Query()["statuses"]; len(got) != 2 || got[0] != "ACTIVE" || got[1] != "CANCELLED" {
			t.Fatalf("unexpected statuses: %#v", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"domain":"example.com","status":"ACTIVE","locked":true,"privacy":false,"renewAuto":true,"nameServers":["ns1.example.com"]}]`))
	}))
	defer server.Close()

	client, err := New("godaddy:test", config.ProviderConfig{
		APIKey:    "key",
		APISecret: "secret",
		BaseURL:   server.URL,
	}, config.Duration{Duration: 5 * time.Second})
	if err != nil {
		t.Fatal(err)
	}

	domains, err := client.ListDomains(context.Background(), providers.DomainListOptions{
		ShopperID: "shopper-123",
		Statuses:  []string{"ACTIVE", "CANCELLED"},
		Limit:     25,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(domains) != 1 || domains[0].Domain != "example.com" {
		t.Fatalf("unexpected domains: %#v", domains)
	}
}

func TestListRecordsUsesTypeAndNamePath(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Path; got != "/v1/domains/example.com/records/A/www" {
			t.Fatalf("unexpected path: %s", got)
		}
		if got := r.URL.Query().Get("limit"); got != "10" {
			t.Fatalf("unexpected limit: %s", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"name":"www","type":"A","data":"1.2.3.4","ttl":600}]`))
	}))
	defer server.Close()

	client, err := New("godaddy:test", config.ProviderConfig{
		APIKey:    "key",
		APISecret: "secret",
		BaseURL:   server.URL,
	}, config.Duration{Duration: 5 * time.Second})
	if err != nil {
		t.Fatal(err)
	}

	records, err := client.ListRecords(context.Background(), providers.RecordListOptions{
		Domain: "example.com",
		Type:   "a",
		Name:   "www",
		Limit:  10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].Data != "1.2.3.4" {
		t.Fatalf("unexpected records: %#v", records)
	}
}
