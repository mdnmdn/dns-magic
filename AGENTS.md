# Repository Guidelines

## Project Structure & Module Organization
This repository is a Go CLI module for DNS inspection and provider-backed record management. The canonical structure, package responsibilities, and command map are defined in [_docs/project-structure.md](/Users/mdn/works/projects/dns-magic-cli/_docs/project-structure.md). Provider contracts and GoDaddy-specific behavior are defined in [_docs/dns-provider.md](/Users/mdn/works/projects/dns-magic-cli/_docs/dns-provider.md) and [_docs/provider-godaddy.md](/Users/mdn/works/projects/dns-magic-cli/_docs/provider-godaddy.md). Read those documents before changing provider code or CLI flags.

Use `cmd/dns-magic` for the binary entrypoint, keep Cobra command wiring in `internal/cli`, place configuration loading in `internal/config`, put DNS lookup logic in `internal/dns`, isolate provider implementations under `internal/providers/<provider>`, and keep output serialization in `internal/output`. Keep user-facing documentation in `_docs/`; do not spread design notes across the repository root.

## Build, Test, and Development Commands
Use `just` to discover tasks in `justfile`. Use `just test` to run the full test suite and `go test ./... -run TestName` for a single test when tests are added. Use `just build` to compile the CLI entrypoint and `just fmt` to format Go files.

If Cobra or TOML dependencies change, update `go.mod` and `go.sum` together, typically via `just tidy`. Keep command names, flags, and examples aligned with the structure documented in [_docs/project-structure.md](/Users/mdn/works/projects/dns-magic-cli/_docs/project-structure.md).

## Coding Style & Naming Conventions
Follow standard Go formatting with `gofmt`. Keep packages small and purpose-specific; avoid placing provider-specific logic in generic DNS packages. Use lowercase package names, exported identifiers only for cross-package APIs, and Cobra command constructors in the `New<Name>Command` form.

Represent settings through typed config structs loaded from TOML. Provider credentials and endpoints belong in config parsing or provider packages, not in command handlers.

## Documentation Maintenance
Treat [_docs/project-structure.md](/Users/mdn/works/projects/dns-magic-cli/_docs/project-structure.md) as the source of truth for architecture and CLI layout. When commands, packages, or config shape change, update that document in the same change set as the code.
