# GoDaddy Provider Requirements

## Scope
The GoDaddy provider implements provider alias resolution, account-level domain listing, delegated reseller access, and DNS record listing for `dns-magic`.

## Authentication and Environments
GoDaddy requests use:

- `Authorization: sso-key <api_key>:<api_secret>`
- `X-Shopper-Id: <shopper_id>` only when delegated access is needed

Base URLs:

- OTE: `https://api.ote-godaddy.com`
- Production: `https://api.godaddy.com`

Config fields:

- `type = "godaddy"`
- `api_key`
- `api_secret`
- optional `base_url`
- optional `shopper_id`
- optional `ote`
- optional `timeout`

## Delegated Access
GoDaddy documents two API user models:

- Self-serve: operate on your own account and ignore `X-Shopper-Id`
- Reseller: operate on behalf of a customer and send the customer shopper id in `X-Shopper-Id`

In this project:

- a provider alias may include a default `shopper_id`
- `--shopper-id` overrides the configured value for a single command
- multiple GoDaddy accounts or delegated customers are represented as distinct provider aliases such as `godaddy:customer1`

## Supported Operations
### `list-domain`
Maps to `GET /v1/domains` and supports:

- `--status`
- `--status-group`
- `--include`
- `--limit`
- `--marker`
- `--modified-since`
- `--shopper-id`

### `list <domain>`
Maps to GoDaddy DNS record retrieval and supports:

- `--type`
- `--name` with `--type`
- `--offset`
- `--limit`
- `--shopper-id`

The public GoDaddy Domains Swagger documents record retrieval under operation `recordGet` and clearly documents the type+name path. The broad path rendering is inconsistent across the published spec, so the implementation should keep the request shaping isolated in the GoDaddy adapter and document the inconsistency rather than exposing it to commands.

Supported normalized record fields:

- `domain`
- `name`
- `type`
- `data`
- `ttl`
- `priority`
- `port`
- `protocol`
- `service`
- `weight`

## Rate Limits and Operational Notes
GoDaddy states a default limit of 60 requests per minute and notes that production access to management and DNS APIs may require qualifying accounts. Expect `401`, `403`, `404`, `409`, `422`, `429`, and `5xx` responses and preserve the upstream details in errors.

## Reference URLs
- https://developer.godaddy.com/doc/endpoint/domains
- https://developer.godaddy.com/swagger/swagger_domains.json
- https://developer.godaddy.com/doc/endpoint/dns
- https://developer.godaddy.com/getstarted?source=post_page---------------------------
- https://pkg.go.dev/github.com/toon-format/toon-go
- https://github.com/toon-format/toon-go
