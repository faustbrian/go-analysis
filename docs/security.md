# Security and threat model

Model version 1 (2026-09-30) describes the analyzer behavior in the signed
`v1.2.0` source at `3bfc234fb561ae1f1c9d45a72a9fe1f2c1849930` and on
`main` until a behavior change supersedes it. This post-release risk
disposition is maintained documentation; it does not alter the immutable tag
or establish that later source revisions have been reviewed.

## Trust boundaries

The analyzer binary, pinned Go toolchain, checked-in policy, and release
workflow are trusted. Target repositories, Go source, generated headers,
package metadata, build tags, configuration input, and report consumers are
potentially untrusted. The tool has the invoking user's filesystem permissions;
it is not a sandbox and must not be run with broader credentials than needed.

`analysis` parses and type-checks target code through Go tooling. It never
executes target binaries, package initializers, tests, generators, or arbitrary
configuration programs. The project does not load analyzer or configuration
plugins. The optional golangci-lint module plugin remains unshipped while its
lifecycle is outside the reproducible compatibility contract.

## Threats and controls

| Threat | Control | Residual risk and operation |
| --- | --- | --- |
| Target-code execution | Analysis uses syntax, types, SSA, CFGs, and package metadata only; no target `go run`, test, generator, or plugin execution | Go package loading may invoke the Go command and read module metadata; run in a least-privilege checkout with an explicit module-download policy |
| Configuration injection | Strict single-document YAML decoding rejects unknown keys and executable configuration | A valid policy can intentionally weaken rules; policy changes require code review and compatibility review |
| Path traversal or disclosure | Analyzer target-relative paths and report paths are cleaned and rejected when they escape the analysis root; `sync-policy` reads or writes only the file paths explicitly selected by the invoking operator | The command is not a filesystem sandbox and uses the operator's permissions; package names and relative paths remain sensitive repository metadata |
| Source or secret disclosure | JSON and SARIF omit source snippets and diagnostic traces; messages use governed metadata rather than source values | A filename, package, rule, or suppression reason may reveal design information; protect reports like source artifacts |
| SARIF or JSON injection | Standard encoders escape untrusted strings; stable schemas and path validation precede emission | Downstream renderers remain separate trust boundaries and must stay patched |
| Forged generated code | Exclusion requires explicit policy, an exact trusted repository-relative path, and a recognized generated header before the package clause; unlisted forged headers and their suppressions remain analyzed | A reviewed generator can still emit unsafe code at an authorized output path; generated output remains subject to generator and supply-chain review |
| Suppression bypass | Directives require an exact known rule, adjacent diagnostic, non-empty reason, unique location, valid optional expiry, and matching finding | A semantically poor but syntactically valid reason requires human review; inventories support audit and trend checks |
| Resource exhaustion | Configuration input, diagnostics, suppressions, SSA traces, static fan-out proofs, corpus entries, and benchmark budgets are bounded | Extremely large valid packages still consume parser and type-checker resources; CI timeouts and representative corpus budgets remain required |
| Dependency or release compromise | Dependencies and tools are pinned, actions use commit SHAs, local archive scripts check reproducibility and checksums, and the current CI is read-only for repository contents | Current releases may be source-only; consumers must verify the trusted tag and public module checksum, plus archive checksums if archives are published, and retain independent dependency, vulnerability, and CodeQL gates |
| Advisory escalation | Rule metadata defaults to advisory; configured reporting separates severity from blocking status; NilAway runs separately with visible advisory status | Raw multichecker and vettool execution use Go vet exit semantics, so use configured `check` when advisory status must be preserved |

## Accepted residual risks

The controls above reduce but do not eliminate these risks. Each owner is
responsible for the stated mitigation in its own environment and for reopening
the disposition when the review condition occurs.

| Residual risk | Owner | Acceptance rationale and mitigation | Review condition |
| --- | --- | --- | --- |
| Package loading and module downloads | Invoking operator or CI owner | Go package metadata is needed for analysis; use least-privilege credentials and an explicit module-download/network policy | Package-loading behavior or download policy changes, or unexpected process/network activity |
| Valid policy weakening | Policy owner and reviewers | Configurability is intentional; review rule, exception, and compatibility changes as security-relevant policy | New policy syntax, rule promotion or exception, or unexpected blocking-status drift |
| Operator-selected paths, permissions, and path metadata | Invoking operator or CI owner | `sync-policy` acts on explicitly selected local files, not in a sandbox; select trusted paths, use least privilege, and restrict metadata exposure | Higher-privilege use, a new sandbox expectation, changed path handling, or metadata disclosure |
| Report secrecy and downstream rendering | Report storage and renderer owners | Useful diagnostics necessarily expose repository metadata; restrict access and retention, and keep renderers patched | New report fields, renderer/export path, public upload, or disclosure finding |
| Authorized generators and suppression reasons | Generator and repository policy owners | Authorized output paths and justified exceptions support repository workflows; review generator output and suppression inventories | Generator or exclusion changes, expired or abused suppression, or an unjustified reason |
| Large valid packages | Invoking CI owner and analyzer maintainers | Analyzer-owned inputs are bounded, but Go parsing and type checking can still exhaust resources; retain CI timeouts and representative corpus budgets | Sustained memory/time budget breach, loader change, or adversarial valid corpus |
| Dependencies and release integrity | Repository and release maintainers | External tools and publication remain trust dependencies; review pins, retain security gates, and verify trusted tag and public module checksum | Dependency, action, toolchain, or release process change; checksum/signature mismatch or compromise |
| Advisory exit semantics | Integrating caller or CI owner | Raw vettool execution intentionally follows Go vet exits; use configured `check` when advisory status is required | Switching invocation mode or changing rule severity/blocking policy |

## Report handling

Reports intentionally contain no source snippets, values, data-flow traces, or
absolute repository paths. Retain JSON, SARIF, exception, and suppression
inventories only as long as required by the organization's engineering and
security evidence policy. Do not publish private report artifacts merely
because they are machine-readable.

## Supply-chain verification

The local release-verification script builds every candidate archive twice,
compares bytes, and emits a SHA-256 checksum manifest. The current CI workflow
does not publish archives or produce release provenance. Consumers should pin
an exact module version and verify its public checksum and trusted tag; for a
release with attached archives, also verify the archive checksum. Keep gosec,
govulncheck, and CodeQL independently enabled. Claim signed provenance only
when an actual trusted release environment produces it.

The blocking CI workflow runs the pinned CodeQL Go query suite with a reviewed
manual `go build -trimpath ./...` step. Go does not support CodeQL's `none`
mode; explicit compilation gives CodeQL complete production-package evidence
without running target binaries, package initializers, tests, generators, or
arbitrary build scripts. Its job receives only read access to repository
contents and write access to code-scanning results. Local release-equivalent
gates continue to run the enabled gosec integration in the pinned
golangci-lint binary and the separately pinned govulncheck command without
hiding either authority behind this project's diagnostics.

`make workflows` tests and enforces the workflow trust boundary locally.
Every external action must use a full commit SHA. The current CI workflow has
read-only contents permission; CodeQL has code-scanning write access. There is
no tagged-release publication job in this repository.

Security defects include target execution, path escape, source disclosure,
suppression bypass, unbounded attacker-controlled analysis, report injection,
and artifact-integrity failures. Follow [the private reporting process](../SECURITY.md)
without attaching proprietary source or diagnostic output.
