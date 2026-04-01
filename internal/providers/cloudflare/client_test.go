package cloudflare

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mdnmdn/dns-magic/internal/config"
	"github.com/mdnmdn/dns-magic/internal/providers"
)

func TestListDomains(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer my-token" {
			t.Fatalf("unexpected authorization header: %s", got)
		}
		if got := r.URL.Path; got != "/zones" {
			t.Fatalf("unexpected path: %s", got)
		}
		if got := r.URL.Query().Get("per_page"); got != "10" {
			t.Fatalf("unexpected per_page: %s", got)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":[{"id":"zone-123","name":"example.com","status":"active","name_servers":["ns1.cloudflare.com"]}]}`))
	}))
	defer server.Close()

	client, err := New("cloudflare:test", config.ProviderConfig{
		APIToken: "my-token",
		BaseURL:  server.URL,
	}, config.Duration{Duration: 5 * time.Second})
	if err != nil {
		t.Fatal(err)
	}

	domains, err := client.ListDomains(context.Background(), providers.DomainListOptions{
		Limit: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(domains) != 1 || domains[0].Domain != "example.com" {
		t.Fatalf("unexpected domains: %#v", domains)
	}
}

func TestListRecords(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer my-token" {
			t.Fatalf("unexpected authorization header: %s", got)
		}

		if r.URL.Path == "/zones" {
			if got := r.URL.Query().Get("name"); got != "example.com" {
				t.Fatalf("unexpected zone name query: %s", got)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"result":[{"id":"zone-123","name":"example.com"}]}`))
			return
		}

		if r.URL.Path == "/zones/zone-123/dns_records" {
			if got := r.URL.Query().Get("type"); got != "A" {
				t.Fatalf("unexpected type query: %s", got)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"result":[{"id":"rec-1","name":"www.example.com","type":"A","content":"1.2.3.4","ttl":300}]}`))
			return
		}

		t.Fatalf("unexpected path: %s", r.URL.Path)
	}))
	defer server.Close()

	client, err := New("cloudflare:test", config.ProviderConfig{
		APIToken: "my-token",
		BaseURL:  server.URL,
	}, config.Duration{Duration: 5 * time.Second})
	if err != nil {
		t.Fatal(err)
	}

	records, err := client.ListRecords(context.Background(), providers.RecordListOptions{
		Domain: "example.com",
		Type:   "A",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].Data != "1.2.3.4" {
		t.Fatalf("unexpected records: %#v", records)
	}
}
