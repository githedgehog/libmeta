set shell := ["bash", "-euo", "pipefail", "-c"]

import "hack/tools.just"

# Print list of available recipes
default:
  @just --list

export CGO_ENABLED := "0"

_gotools:
  go fmt ./...
  go vet {{go_flags}} ./...

# Called in CI
_lint: _license_headers _gotools

# Generate, lint, test and build everything
all: gen lint lint-gha test && version

# Run linters against code (incl. license headers)
lint: _lint _golangci_lint
  {{golangci_lint}} run --show-stats ./...

# Run golangci-lint to attempt to fix issues
lint-fix: _lint _golangci_lint
  {{golangci_lint}} run --show-stats --fix ./...

go_flags := ""
go_build := "go build " + go_flags

_kube_gen:
  # Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject implementations
  {{controller_gen}} object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Generate docs, code/manifests, things to embed, etc
gen: _kube_gen

# Update golden files
test-update path="./...": gen && test
  UPDATE=true go test {{go_flags}} `go list {{path}} | grep -v /e2e`
