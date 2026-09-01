# analysis

[![CI](https://github.com/faustbrian/go-analysis/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-analysis/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-analysis/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-analysis.svg)](https://pkg.go.dev/github.com/faustbrian/go-analysis)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-analysis?sort=semver)](https://github.com/faustbrian/go-analysis/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`analysis` is a deterministic `go/analysis` policy suite for Go
organizations. It enforces repository architecture, context propagation,
lifecycle ownership, HTTP ownership, secure API migration, and typed secret
handling, plus shared mutable state policy, without replacing the compiler,
`go vet`, Staticcheck,
golangci-lint, gosec, govulncheck, CodeQL, race tests, fuzzing, or NilAway.

The v1 API is stable. Every shipped rule is advisory by default until corpus
evidence supports an explicit blocking promotion.

Use the [documentation index](docs/README.md) for configuration, rule design,
rollout, security, compatibility, and maintenance guidance. The package-owned
guides are [contributor guide](CONTRIBUTING.md), [security policy](SECURITY.md),
[changelog](CHANGELOG.md), [complete rule catalog](docs/rules.md),
[command, API, SARIF, and performance reference](docs/reference.md),
[SARIF specification decisions](docs/specification-decisions.md),
[rule governance and conflict matrix](docs/governance.md),
[organization hardening evidence](docs/hardening.md),
[repository rollout and advisory promotion](docs/rollout.md),
[corpus precision and performance](docs/corpus.md),
[release process](docs/release.md), [compatibility policy](docs/compatibility.md),
[custom-rule design](docs/custom-rules.md), and [FAQ](docs/faq.md).

## Requirements

- Go 1.26 or newer
- No target program execution and no configuration plugins

## Five-minute quickstart

Build the pinned local binary:

```sh
mkdir -p .build
go build -trimpath -o .build/golib-analysis ./cmd/golib-analysis
```

### Standalone analyzer

The raw multichecker runs config-free rules and configured rule packages with
empty policy:

```sh
./.build/golib-analysis ./...
```

Organization policy should use the configured reporting command:

```sh
./.build/golib-analysis validate-config analysis.yml
./.build/golib-analysis check -config analysis.yml -format json ./...
./.build/golib-analysis check -config analysis.yml -format sarif ./... \
  > analysis.sarif
```

When policy is owned in a canonical checkout, synchronize it explicitly and
make drift a local and CI failure:

```sh
./.build/golib-analysis sync-policy update \
  ../mono/policies/service.yml analysis.yml
./.build/golib-analysis sync-policy check \
  ../mono/policies/service.yml analysis.yml
```

`LOCAL_POLICY` defaults to `analysis.yml`. Both commands are offline. The
canonical file is validated before an update, and `check` requires exact byte
identity so formatting or comment drift is also reviewable.

`check` exits 0 when no blocking finding remains, 1 for blocking findings, and
2 for invalid arguments, invalid policy, loading failures, or analyzer errors.
Advisory diagnostics never change the exit status to 1.

Print the exact embedded release version with:

```sh
./.build/golib-analysis version
```

### go vet vettool

```sh
go vet -vettool="$PWD/.build/golib-analysis" ./...
```

The vettool interface has no YAML policy channel. It therefore runs the
config-free rules and the empty-policy form of configured analyzers. Use
`golib-analysis check` when repository policy is required. Go vet treats every
emitted diagnostic as a failing result and has no advisory-status channel; use
configured `check` when advisory versus blocking behavior must be preserved.

## Configuration

Configuration uses strict, versioned YAML with explicit architecture,
lifecycle, security, observability, exception, and suppression policy. Unknown,
ambiguous, stale, or unbounded policy is rejected. See the
[configuration reference](docs/configuration.md) for the complete schema and
examples.
## Inventory and reports

```sh
./.build/golib-analysis rules
```

`golib-analysis rules` emits stable JSON metadata and governance ownership for
every rule. Configured `check` reports use repository-relative slash paths,
stable sorting, no source snippets, and no absolute repository paths. See the
[rule catalog](docs/rules.md) for rationale, examples, configuration, and tool
ownership.

## Development

Run `make check`. See [CONTRIBUTING.md](CONTRIBUTING.md) for focused and release
verification.

## Scope

This project does not format code, fork the compiler, execute target programs,
or claim ownership, borrow, or data-race proof. NilAway remains separately
pinned and advisory; its exit status must not be hidden when reports are later
normalized.
