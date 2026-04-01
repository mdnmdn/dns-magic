# Cloudflare Provider Requirements

## Scope
The Cloudflare provider implements provider alias resolution, account-level zone (domain) listing, and DNS record listing for `dns-magic`.

## Authentication and Environments
Cloudflare requests use:

- `Authorization: Bearer <api_token>`

Base URL:

- Production: `https://api.cloudflare.com/client/v4`

Config fields:

- `type = "cloudflare"`
- `api_token`
- optional `base_url`
- optional `timeout`

## Supported Operations
### `list-domain`
Maps to `GET /zones` and supports:

- `--status` (mapped to Cloudflare zone status)
- `--limit` (per_page)
- `--marker` (using Cloudflare pagination if needed, but for now simple page)

### `list <domain>`
Maps to Cloudflare DNS record retrieval: `GET /zones/{zone_id}/dns_records`.

Supported normalized record fields:

- `domain`
- `name`
- `type`
- `content` (mapped to `data`)
- `ttl`
- `priority`
- `proxied` (vendor specific, but may be useful)

Cloudflare identifies domains by `zone_id`. The implementation must first resolve the domain name to a `zone_id` using the `/zones` endpoint.

## Rate Limits and Operational Notes
Cloudflare has various rate limits depending on the plan. Standard API rate limit is 1200 requests per 5 minutes.

## Reference URLs
- https://developers.cloudflare.com/api/operations/zones-get-zones
- https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-list-dns-records
