# Organization hardening evidence

This page is the package-owned index for rule precision, promotion, and release
evidence. Shared repository verification is pinned in [`.golib.yaml`](../.golib.yaml)
and executed by `go-library-tools`; private corpus and mutation artifacts are
kept under `.verification` when they are suitable for publication.

Every shipped rule remains advisory until its configured corpus findings,
precision fixtures, coverage, mutation evidence, and ownership decision are
reviewed. No rule is promoted by severity alone.

## Rule matrix

| Rule | Evidence status |
| --- | --- |
| `api/backend-error-boundary` | Advisory; package fixtures and mutation evidence required for promotion |
| `api/forbidden-call` | Advisory; package fixtures and mutation evidence required for promotion |
| `api/interface-naming` | Advisory; package fixtures and mutation evidence required for promotion |
| `api/interface-placement` | Advisory; package fixtures and mutation evidence required for promotion |
| `architecture/import-boundary` | Advisory; package fixtures and mutation evidence required for promotion |
| `context/blocking-api-context` | Advisory; package fixtures and mutation evidence required for promotion |
| `context/no-background` | Advisory; package fixtures and mutation evidence required for promotion |
| `context/no-stored-context` | Advisory; package fixtures and mutation evidence required for promotion |
| `http/client-timeout` | Advisory; package fixtures and mutation evidence required for promotion |
| `http/no-default-client` | Advisory; package fixtures and mutation evidence required for promotion |
| `lifecycle/cleanup-ownership` | Advisory; package fixtures and mutation evidence required for promotion |
| `lifecycle/lock-across-call` | Advisory; package fixtures and mutation evidence required for promotion |
| `lifecycle/no-constructor-goroutine` | Advisory; package fixtures and mutation evidence required for promotion |
| `lifecycle/no-global-goroutine` | Advisory; package fixtures and mutation evidence required for promotion |
| `lifecycle/no-init` | Advisory; package fixtures and mutation evidence required for promotion |
| `lifecycle/no-process-control` | Advisory; package fixtures and mutation evidence required for promotion |
| `lifecycle/unbounded-goroutine-fanout` | Advisory; package fixtures and mutation evidence required for promotion |
| `lifecycle/transaction-rollback` | Advisory; package fixtures and mutation evidence required for promotion |
| `observability/dynamic-label-name` | Advisory; package fixtures and mutation evidence required for promotion |
| `observability/high-cardinality-label` | Advisory; package fixtures and mutation evidence required for promotion |
| `safety/no-mutable-global` | Advisory; package fixtures and mutation evidence required for promotion |
| `security/no-unsafe` | Advisory; package fixtures and mutation evidence required for promotion |
| `security/sensitive-sink` | Advisory; package fixtures and mutation evidence required for promotion |

## Evidence boundaries

Compatibility snapshots under `compat` are reviewed as exact files. Mutation
bootstrap archives and ledgers under `.verification/mutation` are preserved by
exact content identity during tooling migrations. Corpus and performance
reports may contain repository metadata and must be reviewed before publication.

The release gate must be fresh at the candidate commit. Local success does not
prove hosted CI or release readiness; the pinned workflow remains authoritative.
