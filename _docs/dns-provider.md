# DNS Provider Interface

## Purpose
The provider layer isolates CLI commands from DNS vendor APIs. Commands should deal in provider-neutral request options and normalized output models, while each provider adapter handles authentication, endpoint selection, request shaping, error translation, and vendor quirks.

## Provider Identity and Config
Providers are selected by alias, not by type. The config file stores provider instances under `providers."<alias>"`, for example:

```toml
default_provider = "godaddy:customer1"

[providers."godaddy:customer1"]
type = "godaddy"
api_key = "..."
api_secret = "..."
base_url = "https://api.godaddy.com"
shopper_id = ""
timeout = "10s"
```

This allows multiple accounts or delegated shoppers for the same provider type. `default_provider` must reference one configured alias.

## Interface Contract
The provider interface currently supports:

- `ListDomains(ctx, opts)` for account-level domain inventory.
- `ListRecords(ctx, opts)` for zone-level DNS record listing.

Request options are provider-neutral:

- `DomainListOptions`: `shopper_id`, `statuses`, `status_groups`, `includes`, `limit`, `marker`, `modified_since`
- `RecordListOptions`: `domain`, optional `type`, optional `name`, optional `shopper_id`, `offset`, `limit`

Current implementations include:
- `internal/providers/godaddy` for GoDaddy-specific API client and record translation.
- `internal/providers/cloudflare` for Cloudflare-specific API client and record translation.

The CLI must never construct vendor-specific headers directly. Delegated access is always expressed through request options and resolved by the provider adapter.

## Output and Models
Providers return normalized `dns.DomainSummary` and `dns.Record` values. Rendering is handled separately and must support `table`, `json`, `yaml`, `markdown`, and `toon`.

`toon` output is implemented through the TOON Go package rather than custom text serialization:

- https://pkg.go.dev/github.com/toon-format/toon-go
- https://github.com/toon-format/toon-go

## Error Handling
Providers should surface actionable errors with request path context and upstream status details. Authentication, delegated-access, validation, and rate-limit failures must remain visible to the CLI caller; do not collapse them into generic “request failed” messages.
