# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected module's current docs.
2. Run `make inventory` and the narrow baseline gate for the module.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Repository verification is configured in [`.golib.yaml`](.golib.yaml) and the
workflow is pinned to the same released `go-library-tools` version. Use
`make inventory` to inspect the manifests, `golib config validate` to validate
the contract, and `make workflows` to validate the workflow locally.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design), update
the affected stable decision entries, and complete the Specification Decisions
section of the pull request template. An unresolved interpretation or stale
source pin is release-blocking; peer behavior cannot silently select policy.
The root SARIF decisions are recorded in the
[specification decision register](docs/specification-decisions.md).

Required release mutation gates must finish with zero surviving viable mutants.

Do not fork shared verification workflows, add permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Use Go 1.27.2; public module minima remain Go 1.27.0. An isolated checkout
needs no sibling repositories or coordination goal files. Begin with root
deterministic tests and one active target when configuration decoding changes,
using a disposable cache per bounded session:

```bash
(
  set -eu
  GOCACHE="$(mktemp -d)"
  export GOCACHE
  trap 'rm -rf "${GOCACHE}"' EXIT
  GOWORK=off go test ./...
  GOWORK=off go test ./analysis -run '^$' -fuzz '^FuzzDecodeConfig$' -fuzztime=2s -timeout=2m
)
```

Deterministic tests execute fuzz seeds/corpus, not active fuzzing. The
intentional `testdata/coverage` fixture remains a separate test module, not
a releasable library. Install the public v2.1.3 binary archive for your platform
from [go-library-tools releases](https://github.com/faustbrian/go-library-tools/releases/tag/v2.1.3),
checking `checksums.txt` against `.golib.yaml` and the archive against that
checksum set. Put the extracted `golib` binary on `PATH`, or select it through
`GOLIB`; do not build Tools from source. CI uses this public release while
checking out its immutable reusable workflow source.

Run adjacent commands in one disposable-cache session as shown above:

```bash
make inventory
make check
make ci
make specification-check
golib mutation --module . --package ./analysis
```

`make check` and `make ci` select the bounded contract. A non-race, once-only
one-minute report-writer check rejects absolute/traversal diagnostic and
suppression paths and verifies SARIF contains no source/snippet fields.
It is not certification of every analyzer, workspace or source-disclosure
path. Base-free checks run seeds/corpus; changed-source CI selects causal
active fuzz smoke. Optional mutation selects one production package, one
worker and a twenty-minute deadline. Offline specification checking is
ordinary; no external services or test tags are declared.

Select full qualification and release rehearsal explicitly:

```bash
make check-full
make specification-check-online
make -f verification/package.mk package-specific
make -f verification/package.mk compatibility
golib api check
golib release dry-run --module .
```

Full checks retain all four fuzz targets at 10000x/20m, the 1s/20m benchmark,
package-specific/toolchain/corpus/performance/reproducibility checks, policy
documentation, compatibility and report/CLI conformance declarations. The
coverage fixture, mutation imports and API baseline remain unchanged.
Scheduled CI uses the shared bounded route, not an automatic full campaign.
Release packaging remains explicitly selected as described in
[the release guide](docs/release.md). Report unavailable or failing commands;
bounded checks alone do not establish release qualification.

## Selected Monthly Assurance

The owner-level [selected-assurance caller](.github/workflows/assurance.yml) is
configured for main only, on day 3 of each month at 05:07 UTC. It selects one
of the 4 current declared root-module fuzz targets using the shared UTC-week
rotation, with 30 seconds of active exploration. The selected risks are:
configuration/suppression parsing, source-free JSON/SARIF report writing and
bounded in-memory analyzer execution.
The testdata coverage fixture stays a separate test module. Selected fuzzing
does not certify every analyzer, compiler/toolchain combination or workspace
source-disclosure boundary; full compatibility and package-specific checks
remain explicit.

Base-free ordinary checks execute deterministic seeds, not active exploration;
committed-range CI selects affected active smoke. Manual assurance requires an
exact anchored selector from the root fuzz operation in `.golib.yaml`, for
example `^FuzzReportWriters$`.
Empty, unknown and malformed selections refuse. The shared job bounds setup
and execution to ten minutes, active workers and CPU to two, owned memory to
2 GiB, native output to 1 MiB and private-cache lifetime. Failure, cancellation,
exhausted budgets and cleanup refusal remain attributable failures. This lane
does not run mutation or replace Required, original full fuzz budgets,
explicit `make check-full`, or relevant release qualification.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
