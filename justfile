set shell := ["zsh", "-lc"]

default:
  @just --list

fmt:
  gofmt -w cmd/dns-magic/main.go internal/app/app.go internal/cli/*.go internal/config/*.go internal/dns/*.go internal/providers/*.go internal/providers/godaddy/*.go internal/providers/cloudflare/*.go

test:
  go test ./...

build:
  go build ./cmd/dns-magic

run *args:
  go run ./cmd/dns-magic {{args}}

tidy:
  go mod tidy

check name dns="":
  if [[ -n "{{dns}}" ]]; then go run ./cmd/dns-magic check "{{name}}" --dns "{{dns}}"; else go run ./cmd/dns-magic check "{{name}}"; fi

list domain type="" provider="":
  cmd=(go run ./cmd/dns-magic list "{{domain}}")
  if [[ -n "{{type}}" ]]; then cmd+=(--type "{{type}}"); fi
  if [[ -n "{{provider}}" ]]; then cmd+=(--provider "{{provider}}"); fi
  "${cmd[@]}"

search query type="" provider="":
  cmd=(go run ./cmd/dns-magic search "{{query}}")
  if [[ -n "{{type}}" ]]; then cmd+=(--type "{{type}}"); fi
  if [[ -n "{{provider}}" ]]; then cmd+=(--provider "{{provider}}"); fi
  "${cmd[@]}"

list-domain provider="" output="table":
  cmd=(go run ./cmd/dns-magic list-domain --output "{{output}}")
  if [[ -n "{{provider}}" ]]; then cmd+=(--provider "{{provider}}"); fi
  "${cmd[@]}"
