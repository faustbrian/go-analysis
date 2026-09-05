# Configuration

Applications load policy files with
`analysis.LoadConfigContext(ctx, path, knownRules)`. Cancellation is observed
before and between filesystem operations. The older `analysis.LoadConfig`
entry point remains available for v1 compatibility but cannot accept a caller
context and is deprecated.
The replacement call and frozen support interval are recorded in
[`DEPRECATION.md`](../DEPRECATION.md#analysisloadconfig).

Configuration is strict YAML. Unknown fields, unknown rule IDs, multiple YAML
documents, malformed analyzer policy, and unsupported versions are rejected.
Paths and package patterns are resolved from the directory containing the
configuration file, not the invocation directory. Canonical policy runners use
the explicit absolute `-root` override when policy is stored outside the
analyzed checkout.

```yaml
version: 1

entrypoints:
  - example.com/service/cmd/service

init_packages:
  - example.com/service/internal/runtimeinit

context_owners:
  - example.com/service/internal/requeststate

generated:
  exclude: true
  paths:
    - internal/protocol/client.gen.go

layers:
  - name: domain
    may_import: [shared]
  - name: infrastructure
    may_import: [domain, shared]
  - name: shared

contexts:
  - name: orders
    may_import: [shared]
  - name: shared

packages:
  - pattern: example.com/service/domain/...
    layer: domain
    context: orders
    deny_imports:
      - example.com/service/infrastructure/...
    allow_imports:
      - example.com/service/infrastructure/approved
  - pattern: example.com/service/internal/repository
    blocking_functions:
      - Repository.Load

constructors:
  - package: example.com/service/internal/worker
    symbols:
      - New
      - Builder.Build

transactions:
  - package: database/sql
    symbol: DB.BeginTx
    result: 0
    rollback_method: Rollback

forbidden_apis:
  - package: example.com/legacy/backend
    symbol: NewClient
    replacement: example.com/service/internal/adapter.NewClient
    allowed_packages:
      - example.com/service/internal/adapter

backend_clients:
  - package: database/sql
    allowed_packages:
      - example.com/service/internal/adapter/sql/...

mutable_globals:
  - package: example.com/service/internal/...

interface_provider_packages:
  - example.com/service/internal/providers/...

interface_names:
  - package: example.com/service/internal/ports/...
    required_suffix: Port
    allowed_names:
      - Compatibility

metric_label_types:
  - package: example.com/service/internal/model
    name: UserID

metric_label_sinks:
  - package: example.com/service/internal/metrics
    symbol: Counter.Label
    arguments: [0]

metric_label_name_types:
  - package: example.com/service/internal/request
    name: MetricName

metric_label_name_sinks:
  - package: example.com/service/internal/metrics
    symbol: Counter.LabelName
    arguments: [0]

backend_error_boundaries:
  - example.com/service/api/...

backend_error_sources:
  - package: example.com/backend
    symbol: Client.Load
    result: 1

backend_error_passthroughs:
  - package: fmt
    symbol: Errorf
    result: 0
    variadic_from: 1

goroutine_fanout:
  - package: example.com/service/internal/worker/...
    max_static: 8

exceptions:
  - rule: security/no-unsafe
    package: example.com/service/internal/bridge
    path: internal/bridge/abi.go
    reason: reviewed operating-system ABI bridge
    issue: SEC-42
    expires: 2027-01-31

sensitive_types:
  - package: example.com/security
    name: Token

sensitive_sinks:
  - package: log/slog
    symbol: Logger.Log
    arguments: [2]
    variadic_from: 3
    allowed_packages:
      - example.com/service/internal/audit

rules:
  architecture/import-boundary:
    status: blocking
    promotion:
      version: 0.1.0
      evidence: ARCH-101 reviewed owned-corpus baseline
  security/sensitive-sink:
    status: advisory
    severity: error
```

Statuses are `disabled`, `advisory`, or `blocking`. Every blocking override
requires a semantic `promotion.version` and a non-empty `promotion.evidence`
record. Severities are `info`, `warning`, or `error`. A severity affects
reports; only status controls the blocking exit code.

Layer and bounded-context names use lower-kebab-case. Each `may_import` list
names other partitions that may be imported; imports within the same layer or
context are implicitly allowed. The configured layer and context graphs must
be acyclic; validation reports a deterministic cycle trace before packages are
loaded. Package classifications must not overlap. Direction checks run only
when both importer and imported package are classified for that dimension, so
unclassified dependencies are not guessed at. Explicit `deny_imports` takes
precedence at the same import location.

`backend_clients` reverses the boundary declaration: each backend package tree
names the non-overlapping adapter package trees permitted to import it. This
prevents a new service package from bypassing the adapter merely because it
was not yet named by source-oriented package policy. Both client and adapter
entries use exact import paths or a trailing `/...` package-tree suffix.

`transactions` identifies exact transaction constructors, their zero-based
transaction result, and the rollback method that must be deferred immediately
after the constructor's exact terminating error guard. This closes transaction
ownership without duplicating `sqlclosecheck` query and statement closure.

`mutable_globals` explicitly opts exact package paths or package trees into
typed shared-state enforcement. The unconfigured analyzer is inactive, so the
vettool never turns an advisory organization policy into an implicit blanket
style gate. Reviewed declarations use the ordinary exact suppression or
reviewed configuration-exception inventory rather than a separate global
allowlist.

`interface_provider_packages` identifies exact provider packages or trailing
`/...` trees where exported runtime interfaces would invert ownership. Generic
constraint-only interfaces remain accepted; consumer packages are not inferred
from directory names.

`metric_label_types` names organization types known to carry unbounded values,
and `metric_label_sinks` identifies zero-based metric-label positions on exact
functions or `Type.Method` symbols. A direct conversion such as
`string(userID)` preserves the source-type evidence; an explicit bucketing
function is the reviewed boundary that removes it.

`metric_label_name_types` names types proven to carry attacker-controlled label
names, while `metric_label_name_sinks` identifies the exact name positions.
Map request values through a fixed allowlist; do not turn them into metric
schema.

`backend_error_boundaries` identifies exported API package trees.
`backend_error_sources` identifies zero-based results from exact backend
callables, while `backend_error_passthroughs` identifies wrappers that preserve
those errors. An unlisted mapping helper is an explicit classification
boundary; configure only wrappers proven to preserve backend identity.

`goroutine_fanout` opts exact packages or trailing `/...` trees into an
advisory ceiling for loop-expanded goroutine launches. `max_static` accepts
constant-sized work while runtime-sized loops require a worker pool or a
statically proven synchronization bound.

Rule-specific behavior uses typed top-level policy such as `constructors` and
`sensitive_sinks`; arbitrary `rules.<id>.options` are rejected. The reserved
`adapter_roots` field is rejected because unassociated global roots cannot
express which dependency each adapter owns; use typed `backend_clients`
contracts instead.

Generated files remain analyzed by default. `generated.exclude: true` excludes
only exact repository-relative files listed in `generated.paths` that also use
Go's generated-file convention: a matching
`// Code generated ... DO NOT EDIT.` header before the package declaration.
The header, explicit exclusion, and exact trusted path are all required.
Wildcards, directories, traversal, duplicate paths, and non-Go files are
rejected. A filename such as `generated.go` or a forged header at an unlisted
path has no special treatment, so its diagnostics and suppression directives
are still validated. Diagnostics and suppression directives from an explicitly
trusted generated file are omitted together.

Configuration exceptions are central, reviewed policy records rather than
source suppressions. Every exception names one known rule, one exact package,
a non-empty reason, and optionally one exact repository-relative slash path,
an issue, and an expiry date. Package patterns, absolute paths, traversal,
duplicates, overlaps, expired entries, and exceptions that match no diagnostic
are rejected. Configuration exceptions take precedence over source directives;
a directive targeting an excepted finding is therefore rejected as stale.
Applied exceptions are retained in deterministic JSON and SARIF inventories.

## Suppressions

A suppression applies only to a diagnostic on the immediately following line.
It must name one known rule and contain a non-empty reason:

```go
//analysis:ignore security/no-unsafe -- required ABI bridge; issue=SEC-42; expires=2027-01-31
import _ "unsafe"
```

Supported metadata is `issue=<reference>` and `expires=YYYY-MM-DD`.
Malformed, unknown, duplicate, expired, or unused suppressions fail analysis.
JSON and SARIF runs retain a used-suppression inventory for audit.
