# Command and API reference

## Commands

`golib-analysis check -config <path> [-root <absolute-path>]
[-format json|sarif] [-sequential] [packages...]` loads packages once, applies
configured analyzers, and writes one deterministically ordered report. Package
patterns default to `./...`. JSON is the default format. `-sequential` is a
determinism and diagnosis aid; parallel and sequential reports must be
identical.

By default, the configuration directory is the target repository root. Pass
`-root` when a canonical policy is intentionally outside the target checkout.
The override MUST be absolute so package loading, exception paths, generated
files, and report paths remain independent of the invocation directory.

Exit code 0 means no blocking finding, 1 means at least one blocking finding,
and 2 means arguments, configuration, loading, analysis, or report writing
failed. Advisory findings never produce exit code 1.

`golib-analysis validate-config <path>` performs strict decoding and all analyzer
policy validation without loading target packages. Success prints
`configuration valid`.

`golib-analysis sync-policy check <canonical> <local>` validates the canonical
policy and requires the local policy to be byte-identical. Drift and missing
files fail the command. Both files are limited to 1 MiB. `sync-policy update`
validates the exact bounded canonical snapshot it copies to the local path;
later changes to the canonical path cannot substitute unvalidated bytes.
Neither mode uses the network,
loads plugins, or executes configuration. The repository Makefile delegates
verification to the shared tooling contract; policy synchronization remains an
explicit `golib-analysis` command.

`golib-analysis rules` writes the stable JSON inventory of rules, owners, metadata,
configuration properties, and overlap decisions. It accepts no arguments.

`golib-analysis version` writes the exact semantic version embedded at build time.
Release archives inject their tag version; ordinary local builds report the
current development version.

Without a subcommand the binary is a standard multichecker. Use it directly or
as `go vet -vettool=<binary>`. This mode has no YAML policy channel and runs
configured analyzers with empty organization policy. Go vet has no advisory
exit status: any emitted diagnostic fails the invocation. Use configured
`check` whenever governed advisory and blocking statuses must be preserved.

## Reports

JSON contains `tool_version`, sorted rule metadata, diagnostics, reviewed policy
exceptions, and parsed suppressions. A diagnostic contains its stable rule ID,
repository-relative filename, one-based line and column, message, rationale and
remediation. Severity/default status live in the rule inventory, not each finding.
Reports omit source/snippet fields, but arbitrary message, metadata and attribution
text is not redacted. Direct writer callers own valid positions and report sizes.

One run accepts at most 100,000 diagnostics and 100,000 source suppressions.
Diagnostic emission is stopped at the limit before checker accumulation;
checker work itself continues. The ceiling counts findings before generated
exclusions and suppressions, not only retained findings.
Report collection independently rechecks the same bound. Exceeding either
limit fails analysis instead of allocating an unbounded report. SSA backend
error traces stop after 256 values, goroutine proofs stop after 1,024 static
iterations, and lock must-analysis rejects functions above 4,096 CFG blocks or
256 lock identities.

SARIF output is version 2.1.0. Rule rationale and remediation populate the tool
driver descriptors; rule-inventory severity maps to SARIF level; exceptions and
suppressions remain run properties for audit. Artifact URIs are
repository-relative and use forward slashes. Paths that escape the report root
are rejected rather than exposed.
The exact supported profile, maintained-peer differences, and direct-call
position precondition are recorded in the
[SARIF specification decision register](specification-decisions.md).

For CI, write SARIF to a file and let the CI platform upload it:

```sh
./.build/golib-analysis check -config analysis.yml -format sarif ./... \
  > analysis.sarif
```

The project does not upload reports itself and does not require repository or
cloud credentials.

Configured severity overrides currently modify analyzer metadata copies but do
not reach the driver's built-in `Report.Rules` inventory or SARIF levels. Status
overrides do control blocking. This is an existing propagation limitation, not
an intentional interpretation of SARIF or a conformance result.

Package loading observes context; checker execution does not accept it. A
cancelled configuration read cannot preempt an already blocked filesystem call
or buffered YAML decoding.

## Public Go packages

`analysis` is the reporting and policy model. `LoadConfigContext` strictly
decodes a versioned file against known rule IDs while observing caller
cancellation before and between filesystem operations. The deprecated
`LoadConfig` entry point remains as a context-free v1 compatibility wrapper.
`ParseConfigBytes` strictly validates an in-memory policy snapshot using the
given path to determine its repository-relative root without reading that path.
`ParseSuppressions`, `ApplySuppressions`, and `ApplyPolicyExceptions` implement
the auditable exception pipeline. `WriteJSON` and `WriteSARIF` normalize and serialize a
`Report`. `Rule`, `Diagnostic`, `Suppression`, configuration policy types, and
their enums are stable data contracts.

`policy.NewRegistry` validates unique rule IDs, owners, metadata, and overlap
records. Construction retains nested caller maps and overlap slices; do not
mutate them while using the registry. `Entries` copies those collections.
`policy.Builtin` borrows nested exported analyzer rule metadata.

Each `analyzers/<rule>` package exports `Rule` and `Analyzer`. Configurable
packages also export typed `Options` and `New`; `New` rejects ambiguous or
unbounded policy before analysis begins. Callers embedding analyzers should
construct them once and share the package-loading driver rather than load the
same target independently for every rule.

`analysistestkit` provides shared fixtures, a no-panic fuzz target, and aggregate
benchmarks. It is intended for this repository's rule corpus; external custom
rules may use `golang.org/x/tools/go/analysis/analysistest` directly.

`internal/driver` is intentionally not a public API. Invoke the command for
configured analysis and machine-readable reports.

## Architecture policy design

Classify non-overlapping package trees by layer and bounded context. Direction
edges are allowlists: same-partition imports are accepted and cross-partition
imports require an explicit `may_import` edge. Unclassified packages are not
guessed. Add exact `deny_imports` only for narrower package policy.

Backend ownership is target-oriented: each `backend_clients` entry names the
restricted dependency and the adapter package trees permitted to import it.
This prevents a new source package from bypassing policy simply because it was
not classified yet. Keep domain ports free of infrastructure types and place
backend construction inside the approved adapter.

Policy patterns are exact import paths or one trailing `/...`. Overlapping
classifications, client policies, and ambiguous exceptions are configuration
errors; declaration order never establishes precedence.

## Performance and determinism

Run `go test -bench=. ./analysis ./analysistestkit` for per-analyzer, aggregate, and 1,000-diagnostic JSON and
SARIF nanoseconds, bytes, and allocations. Run it on the same machine and Go
version when comparing changes.
Per-analyzer and aggregate allocation budgets run without race or coverage
instrumentation, both of which change allocation accounting. Every shipped
analyzer must have an explicit budget. The ordinary test phase enforces those
budgets before instrumented coverage starts, while the shared race gate
independently proves race safety. `./scripts/performance.sh corpus/performance.tsv`
additionally enforces manifest budgets
for cold and warm wall time and peak resident memory and writes the observation
to `.build/performance.tsv`. Organization policy should cover representative
small, large-library, and service modules.

Run the same check concurrently and with `-sequential` when investigating
nondeterminism. Reports sort rules, diagnostics, suppressions, and exceptions;
configuration maps and package scheduling must not affect bytes emitted.

Run `golib api check` for the configured comment-insensitive API baseline at
`api/baseline.txt`. Separately, `./scripts/compatibility.sh check` compares native
full GoDoc and the rule inventory under `compat`. Comment changes require native
`./scripts/compatibility.sh update` for the documentation snapshot, preserving
unchanged `rules.json`. Intentional contract changes require review of both diffs,
compatibility guidance, and a version decision. Run
`./scripts/reproducible-build.sh` to build the
CGO-disabled, trimpath release binary twice and compare its bytes; the command
prints the resulting SHA-256 checksum.

## Troubleshooting

- Exit 2: run `validate-config`, then inspect the exact load or analysis error.
- A policy has no effect in vettool mode: use `check -config`; vettool has no
  YAML channel.
- Paths differ by invocation directory: pass the same config file and, for an
  external canonical policy, the same absolute `-root` target.
- Generated code is still analyzed: `generated.exclude: true`, an exact entry
  in `generated.paths`, and a valid Go generated header before the package
  clause are all required.
- A suppression is rejected: keep it at the diagnostic location, name one known
  rule ID, and provide a non-empty reason; remove duplicate or expired entries.
- JSON and SARIF differ between runs: reproduce with `-sequential` and retain
  both outputs as a determinism defect fixture.
- A mature linter reports the same concern: consult the conflict matrix and
  disable only the explicitly documented duplicate configuration.
