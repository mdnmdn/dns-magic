# Project Structure

## Overview
`dns-magic` is a Go command-line tool for inspecting DNS state and managing records through provider APIs such as GoDaddy. The module path is `github.com/mdnmdn/dns-magic`. The CLI should use Cobra for command wiring and a TOML configuration file for provider credentials, defaults, and runtime settings.

## Repository Layout
- `cmd/dns-magic/main.go`: binary entrypoint. Keep startup minimal and delegate to `internal/app`.
- `internal/app`: application bootstrap. This layer assembles config, services, and the root Cobra command.
- `internal/cli`: Cobra command constructors and flag definitions.
- `internal/config`: TOML-backed config structs and loading helpers.
- `internal/dns`: provider-agnostic DNS models and lookup helpers.
- `internal/providers`: provider interface and shared provider plumbing.
- `internal/providers/godaddy`: GoDaddy-specific API client and record translation.
- `internal/output`: output renderers for table, json, yaml, markdown, and toon.
- `_docs`: project documentation. All design and contributor docs belong here.
- `config.example.toml`: example config showing provider settings and defaults.
- `justfile`: developer shortcuts for formatting, testing, building, and running common CLI flows.

## Command Surface
The initial command set is:

```text
dns-magic check <name> --dns <dns ip>
dns-magic list <domain> --type cname --provider godaddy
dns-magic search <query> [--type a] [--provider godaddy]
dns-magic add <domain> <name> --type a --value 1.2.3.4 --provider godaddy
dns-magic update <domain> <name> --type a --value 1.2.3.5 --provider godaddy
dns-magic remove <domain> <name> --type a --provider godaddy
```

Notes:
- `check` performs resolver-based inspection and may optionally target a specific DNS server IP with `--dns`.
- `list`, `search`, `add`, `update`, and `remove` operate through a provider selected by `--provider` or by the config default.
- Record `--type` values should remain provider-agnostic at the CLI boundary even if provider adapters need translation.

## Configuration
All settings live in a TOML file. The expected shape is:

```toml
default_provider = "godaddy:customer1"
request_timeout = "10s"

[providers."godaddy:customer1"]
type = "godaddy"
api_key = "..."
api_secret = "..."
base_url = "https://api.godaddy.com"
shopper_id = ""
timeout = "10s"
```

Keep secrets out of source control. Command handlers should receive typed config values from `internal/config` rather than reading environment variables directly.

Provider aliases are the CLI-facing identity. Keep provider contracts and GoDaddy-specific requirements documented in `_docs/dns-provider.md` and `_docs/provider-godaddy.md`.

## Implementation Rules
- Keep Cobra command handlers thin; delegate business logic to internal packages.
- Keep provider-specific behavior under `internal/providers/<provider>`.
- Add tests beside the package they validate.
- When the command surface or package layout changes, update this document and `AGENTS.md` in the same change.
- Keep `justfile` aligned with the current module layout so contributors have stable entrypoints for build and test tasks.
