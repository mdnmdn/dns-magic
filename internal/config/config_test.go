package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadSupportsProviderAliases(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `
default_provider = "godaddy:customer1"
request_timeout = "15s"

[providers."godaddy:customer1"]
type = "godaddy"
api_key = "key"
api_secret = "secret"
shopper_id = "123"
timeout = "5s"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.DefaultProvider != "godaddy:customer1" {
		t.Fatalf("unexpected default provider: %s", cfg.DefaultProvider)
	}
	if got := cfg.Providers["godaddy:customer1"].Timeout.Duration; got != 5*time.Second {
		t.Fatalf("unexpected provider timeout: %v", got)
	}
	if got := cfg.RequestTimeout.Duration; got != 15*time.Second {
		t.Fatalf("unexpected request timeout: %v", got)
	}
}
