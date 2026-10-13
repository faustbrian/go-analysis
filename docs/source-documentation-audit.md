# Source documentation audit

## Pre-edit outcome, owner and checks

Base `a818edc194bd5938fc6c4a2e6bbe2246a431d3e7`, clean detached main-source
worktree. This worker owns the standalone Analysis documentation pass; the
original unrelated Dependabot checkout remains untouched. Parent owns shared
enforcement, remote delivery and aggregate evidence. This is Tier A source prose:
no analyzer logic, configuration, suppression/want annotations, generated-header
classification, fixtures, dependencies, pins, API or frozen history changes.

Outcome: truthful policy construction, conservative analyzer coverage, configuration
and report trust boundaries, suppression lifetime, driver loading/execution,
ownership/cancellation/failure limits and assurance fixture intent. Direct consumers
are standalone analyzer users, the CLI/driver, the policy registry, report writers
and the catalogued coverage test module; inspect each actual dependency contract.
Named risks: overstated ownership/taint/control-flow proof, false-positive scope,
mutable caller options, generated-file trust, diagnostic/source leakage, blocking
IO/callback cancellation, resource limits and suppression temporal semantics.
Smallest checks: all 267 tracked Go AST/significant tokens/literals/directives and
fixture annotations preserved; native rendered package docs, configured apidiff,
existing relevant examples and local docs/spec/links; static source check, pinned
changed-content/staged secrets and ordinary complete final-diff truth review.
No generator execution: all 23 generated headers are modeled analyzer fixtures,
not actual generated production output or claims of an available source producer.
A truthful hand-maintained fixture-role explanation belongs with the test harness,
not a forged recipe or altered modeled header. Go commands await explicit parent
slot with three owned disposable caches and all-path cleanup. No broad analyzer,
fuzz/race/performance/corpus/service/release campaign for unchanged behavior.

## Reproducible full inventory

`git ls-files "*.go"` identifies 267 tracked Go files in 110 directories: 69 test
files and 165 testdata inputs (including the coverage module). modules.json and
packages.json declare the root module's 29 packages and one coverage fixture
module/package. Initial authenticated structural check finds 255 exported comment
omissions, 23 generator-review findings for modeled headers, and two marker reviews
for context.TODO calls in fixture source. These are missing documentation versus
fixture semantics, not a quality score. No package-comment omission is reported;
existing canonical comments still require semantic review.

- `analysis`: `config.go`, `config_fuzz_test.go`, `config_internal_test.go`, `config_test.go`, `diagnostic.go`, `exception.go`, `exception_internal_test.go`, `exception_test.go`, `metadata.go`, `metadata_test.go`, `report.go`, `report_benchmark_test.go`, `report_fuzz_test.go`, `report_internal_test.go`, `report_test.go`, `suppression.go`, `suppression_fuzz_test.go`, `suppression_test.go`.
- `analysistestkit`: `benchmark_test.go`, `doc.go`, `fixtures_test.go`, `fuzz_test.go`, `harness_test.go`, `race_disabled_test.go`, `race_enabled_test.go`, `suggested_fix_test.go`.
- `analyzers/backenderror`: `backenderror.go`, `backenderror_internal_test.go`, `backenderror_test.go`.
- `analyzers/backenderror/testdata/src/backend`: `backend.go`.
- `analyzers/backenderror/testdata/src/boundary`: `alias.go`, `boundary.go`, `buildtag.go`, `generated.go`, `semantic_near_misses.go`.
- `analyzers/backenderror/testdata/src/boundarytree/child`: `child.go`.
- `analyzers/backenderror/testdata/src/boundaryunconfigured`: `unconfigured.go`.
- `analyzers/backenderror/testdata/src/outside`: `outside.go`.
- `analyzers/backenderror/testdata/src/wrapper`: `wrapper.go`.
- `analyzers/blockingcontext`: `blockingcontext.go`, `blockingcontext_internal_test.go`, `blockingcontext_test.go`.
- `analyzers/blockingcontext/testdata/src/blockingapi`: `blockingapi.go`, `buildtag.go`, `generated.go`, `semantic_near_misses.go`.
- `analyzers/blockingcontext/testdata/src/otherapi`: `otherapi.go`.
- `analyzers/cleanupownership`: `cleanupownership.go`, `cleanupownership_internal_test.go`, `cleanupownership_test.go`.
- `analyzers/cleanupownership/testdata/src/approved`: `approved.go`.
- `analyzers/cleanupownership/testdata/src/consumer`: `buildtag.go`, `consumer.go`, `dot.go`, `generated.go`, `semantic_near_misses.go`.
- `analyzers/cleanupownership/testdata/src/resourceapi`: `resourceapi.go`.
- `analyzers/constructorgoroutine`: `constructorgoroutine.go`, `constructorgoroutine_internal_test.go`, `constructorgoroutine_test.go`.
- `analyzers/constructorgoroutine/testdata/src/constructors`: `buildtag.go`, `constructors.go`, `generated.go`, `semantic_near_misses.go`.
- `analyzers/constructorgoroutine/testdata/src/unconfigured`: `unconfigured.go`.
- `analyzers/forbiddenapi`: `forbiddenapi.go`, `forbiddenapi_test.go`.
- `analyzers/forbiddenapi/testdata/src/approved`: `approved.go`.
- `analyzers/forbiddenapi/testdata/src/consumer`: `buildtag.go`, `consumer.go`, `dot.go`, `generated.go`, `semantic_near_misses.go`.
- `analyzers/forbiddenapi/testdata/src/legacy`: `legacy.go`.
- `analyzers/globalgoroutine`: `globalgoroutine.go`, `globalgoroutine_internal_test.go`, `globalgoroutine_test.go`.
- `analyzers/globalgoroutine/testdata/src/globalgoroutine`: `buildtag.go`, `generated.go`, `globalgoroutine.go`, `semantic_near_misses.go`.
- `analyzers/globalgoroutine/testdata/src/globalother`: `globalother.go`.
- `analyzers/goroutinefanout`: `goroutinefanout.go`, `goroutinefanout_internal_test.go`, `goroutinefanout_test.go`.
- `analyzers/goroutinefanout/testdata/src/fanout`: `alias.go`, `buildtag.go`, `fanout.go`, `generated.go`, `semantic_near_misses.go`.
- `analyzers/goroutinefanout/testdata/src/fanoutoutside`: `outside.go`.
- `analyzers/goroutinefanout/testdata/src/fanouttree/child`: `child.go`.
- `analyzers/goroutinefanout/testdata/src/fanoutunconfigured`: `unconfigured.go`.
- `analyzers/httpclienttimeout`: `httpclienttimeout.go`, `httpclienttimeout_internal_test.go`, `httpclienttimeout_test.go`.
- `analyzers/httpclienttimeout/testdata/src/approved`: `approved.go`.
- `analyzers/httpclienttimeout/testdata/src/httptimeout`: `dot.go`, `generated.go`, `httptimeout.go`, `semantic_near_misses.go`.
- `analyzers/importboundary`: `importboundary.go`, `importboundary_internal_test.go`, `importboundary_test.go`.
- `analyzers/importboundary/testdata/src/arch/adapter`: `adapter.go`.
- `analyzers/importboundary/testdata/src/arch/aliasuse`: `alias.go`, `dot.go`.
- `analyzers/importboundary/testdata/src/arch/application`: `application.go`, `buildtag.go`, `generated.go`, `semantic_near_misses.go`.
- `analyzers/importboundary/testdata/src/arch/application/internal`: `internal.go`.
- `analyzers/importboundary/testdata/src/arch/backend/client`: `client.go`.
- `analyzers/importboundary/testdata/src/arch/backend/subclient`: `subclient.go`.
- `analyzers/importboundary/testdata/src/arch/backendconsumer`: `consumer.go`, `dot.go`.
- `analyzers/importboundary/testdata/src/arch/catalog`: `catalog.go`.
- `analyzers/importboundary/testdata/src/arch/contextonly`: `contextonly.go`.
- `analyzers/importboundary/testdata/src/arch/domain`: `domain.go`.
- `analyzers/importboundary/testdata/src/arch/infra`: `infra.go`.
- `analyzers/importboundary/testdata/src/arch/infra/approved`: `approved.go`.
- `analyzers/importboundary/testdata/src/arch/layeronly`: `layeronly.go`.
- `analyzers/importboundary/testdata/src/arch/shared`: `shared.go`.
- `analyzers/importboundary/testdata/src/arch/unclassified`: `unclassified.go`.
- `analyzers/interfacenaming`: `interfacenaming.go`, `interfacenaming_internal_test.go`, `interfacenaming_test.go`.
- `analyzers/interfacenaming/testdata/src/interfaces/other`: `other.go`.
- `analyzers/interfacenaming/testdata/src/interfaces/portapi`: `portapi.go`.
- `analyzers/interfacenaming/testdata/src/interfaces/provider`: `buildtag.go`, `generated.go`, `provider.go`, `semantic_near_misses.go`.
- `analyzers/interfacenaming/testdata/src/interfaces/provider/child`: `child.go`.
- `analyzers/interfaceplacement`: `interfaceplacement.go`, `interfaceplacement_internal_test.go`, `interfaceplacement_test.go`.
- `analyzers/interfaceplacement/testdata/src/interfaces/consumer`: `consumer.go`.
- `analyzers/interfaceplacement/testdata/src/interfaces/portapi`: `portapi.go`.
- `analyzers/interfaceplacement/testdata/src/interfaces/provider`: `buildtag.go`, `generated.go`, `provider.go`, `semantic_near_misses.go`.
- `analyzers/interfaceplacement/testdata/src/interfaces/provider/child`: `child.go`.
- `analyzers/interfaceplacement/testdata/src/interfaces/unconfigured`: `unconfigured.go`.
- `analyzers/lockacrosscall`: `lockacrosscall.go`, `lockacrosscall_internal_test.go`, `lockacrosscall_test.go`.
- `analyzers/lockacrosscall/testdata/src/approved`: `approved.go`.
- `analyzers/lockacrosscall/testdata/src/consumer`: `buildtag.go`, `consumer.go`, `dot.go`, `generated.go`, `semantic_near_misses.go`.
- `analyzers/lockacrosscall/testdata/src/ioapi`: `ioapi.go`.
- `analyzers/metriccardinality`: `metriccardinality.go`, `metriccardinality_internal_test.go`, `metriccardinality_test.go`.
- `analyzers/metriccardinality/testdata/src/labelmodel`: `labelmodel.go`.
- `analyzers/metriccardinality/testdata/src/metricconsumer`: `buildtag.go`, `consumer.go`, `dot.go`, `generated.go`, `semantic_near_misses.go`.
- `analyzers/metriccardinality/testdata/src/metricnameconsumer`: `buildtag.go`, `consumer.go`, `dot.go`, `generated.go`, `semantic_near_misses.go`.
- `analyzers/metriccardinality/testdata/src/metricsink`: `metricsink.go`.
- `analyzers/metriccardinality/testdata/src/metricunconfigured`: `unconfigured.go`.
- `analyzers/mutableglobal`: `mutableglobal.go`, `mutableglobal_internal_test.go`, `mutableglobal_test.go`.
- `analyzers/mutableglobal/testdata/src/mutableglobals`: `buildtag.go`, `generated.go`, `globals.go`, `semantic_near_misses.go`.
- `analyzers/mutableglobal/testdata/src/otherglobals`: `otherglobals.go`.
- `analyzers/mutableglobal/testdata/src/stateapi`: `stateapi.go`.
- `analyzers/mutableglobal/testdata/src/unconfigured`: `unconfigured.go`.
- `analyzers/nobackground`: `nobackground.go`, `nobackground_internal_test.go`, `nobackground_test.go`.
- `analyzers/nobackground/testdata/src/contextroot`: `contextroot.go`.
- `analyzers/nobackground/testdata/src/contextuse`: `buildtag.go`, `contextuse.go`, `contextuse_test.go`, `dot.go`, `generated.go`, `semantic_near_misses.go`.
- `analyzers/nodefaulthttp`: `nodefaulthttp.go`, `nodefaulthttp_internal_test.go`, `nodefaulthttp_test.go`.
- `analyzers/nodefaulthttp/testdata/src/httpdefaults`: `buildtag.go`, `dot.go`, `generated.go`, `httpdefaults.go`, `semantic_near_misses.go`.
- `analyzers/nodefaulthttp/testdata/src/httpother`: `httpother.go`.
- `analyzers/noinit`: `noinit.go`, `noinit_test.go`.
- `analyzers/noinit/testdata/src/initroot`: `initroot.go`.
- `analyzers/noinit/testdata/src/inituse`: `buildtag.go`, `generated.go`, `inituse.go`, `semantic_near_misses.go`.
- `analyzers/noprocesscontrol`: `noprocesscontrol.go`, `noprocesscontrol_internal_test.go`, `noprocesscontrol_test.go`.
- `analyzers/noprocesscontrol/testdata/src/approved`: `approved.go`.
- `analyzers/noprocesscontrol/testdata/src/processuse`: `buildtag.go`, `generated.go`, `processuse.go`, `processuse_test.go`, `semantic_near_misses.go`.
- `analyzers/nostoredcontext`: `nostoredcontext.go`, `nostoredcontext_internal_test.go`, `nostoredcontext_test.go`.
- `analyzers/nostoredcontext/testdata/src/contextowner`: `contextowner.go`.
- `analyzers/nostoredcontext/testdata/src/contextstore`: `buildtag.go`, `contextstore.go`, `generated.go`, `semantic_near_misses.go`.
- `analyzers/nostoredcontext/testdata/src/nocontext`: `nocontext.go`.
- `analyzers/nounsafe`: `nounsafe.go`, `nounsafe_internal_test.go`, `nounsafe_test.go`.
- `analyzers/nounsafe/testdata/src/safeuse`: `safeuse.go`.
- `analyzers/nounsafe/testdata/src/unsafeuse`: `buildtag.go`, `generated.go`, `semantic_near_misses.go`, `unsafeuse.go`.
- `analyzers/sensitivesink`: `sensitivesink.go`, `sensitivesink_internal_test.go`, `sensitivesink_test.go`.
- `analyzers/sensitivesink/testdata/src/approved`: `approved.go`.
- `analyzers/sensitivesink/testdata/src/consumer`: `buildtag.go`, `consumer.go`, `dot.go`, `generated.go`, `semantic_near_misses.go`.
- `analyzers/sensitivesink/testdata/src/secretmodel`: `secretmodel.go`.
- `analyzers/sensitivesink/testdata/src/sinkapi`: `sinkapi.go`.
- `analyzers/transactionrollback`: `transactionrollback.go`, `transactionrollback_internal_test.go`, `transactionrollback_test.go`.
- `analyzers/transactionrollback/testdata/src/txapi`: `txapi.go`.
- `analyzers/transactionrollback/testdata/src/txconsumer`: `buildtag.go`, `consumer.go`, `dot.go`, `generated.go`, `semantic_near_misses.go`.
- `analyzers/transactionrollback/testdata/src/txunconfigured`: `unconfigured.go`.
- `cmd/golib-analysis`: `main.go`, `main_test.go`.
- `examples/quickstart`: `main.go`.
- `internal/driver`: `driver.go`, `driver_internal_test.go`, `driver_test.go`.
- `internal/version`: `version.go`.
- `policy`: `builtin.go`, `builtin_test.go`, `docs_test.go`, `registry.go`, `registry_test.go`.
- `testdata/coverage/sample`: `sample.go`.

## Native generation, tags and missing versus unknown contracts

The API baseline is root apidiff at api/baseline.txt, not full GoDoc prose;
The native scripts/compatibility.sh also owns a full GoDoc text snapshot at
compat/public-api.txt and the rules.json inventory. Regenerate only prose through
that actual producer; rules.json, api/baseline.txt, corpus and typed specification
history remain byte-identical.
The native docs operation runs policy TestDocumentation; select applicable
Examples only if present. Read actual driver generated-file handling, literal
fixture headers, want comments and suppression parsers before any source change.
The race-enabled/disabled test helpers and analyzer buildtag fixtures are modeled
variants, not newly qualified platform/race behavior. No generator route should
be executed on analyzed target code.

Public field omissions need semantic documentation. Questions requiring source
investigation include option deep-copy boundaries, nil defaults and registry
aliasing; package-load/checker cooperation versus real cancellation; entry/count
caps versus preallocation; raw loader/error and configured metadata disclosure;
call-local syntactic/SSA analysis limits versus interprocedural proof; exception
and inline-suppression date/line association; and literal generated-header input
classification. Unanswered policy/rationale will remain explicit rather than be
invented. Full semantic coverage and native outcomes will be recorded below.

## Semantic pass and selected native request

Every tracked source/test/fixture directory and declaration was inventoried,
including all 22 analyzer implementations and 23 active rule descriptors, 165
modeled source inputs, 69 test files, four fuzz targets and three benchmarks.
Private policy compilation/object identity, closure execution scope, root SSA
construction and bounded traces, CFG predecessor must-intersection, loop/channel/
WaitGroup recognizers, generated trust filters, report admission, relative paths,
policy synchronization snapshot selection and prerequisite caching were traced.
Tests were read by case/assertion boundary and high-risk bodies; this is source
audit, not a claim all runtime tests were executed.

Canonical package surfaces now distinguish each rule's configured/static scope,
non-goals, nil/inactive defaults and mutable analyzer metadata. Constructors copy
configuration into compiled structures; generic records and NewRegistry nested
metadata retain caller ownership. Configuration caps admission at 1 MiB, with
YAML allocation before semantic checks; file cancellation cooperates around reads,
not blocked filesystem operations or decode CPU. The driver context reaches loader
but not checker. Its emission ceiling applies before exclusions/suppression and
stops emission, not remaining computation. Report writers impose no independent
size/count limit, omit source fields but preserve arbitrary supplied strings,
borrow writers, and can leave partial output after writer failure. SARIF 2.1.0
position requirements and genuine pinned decisions remain unchanged.

Suppression targets follow the whole parsed comment group, with inclusive UTC
expiry; direct application helpers retain caller Used flags and do not repeat
all validation. Generated headers and context.TODO marker findings are modeled
semantic inputs, not debt/generator failures. Their exact 23 headers, literal
source strings, build constraints, want and suppression annotations remain intact;
no fixture is manually rewritten and no alleged producer is run.

Source-verified existing limitation: driver buildSpecs applies severity overrides
only to analyzer metadata copies while Report.Rules is built from the original
registry. JSON rule metadata and SARIF levels thus retain built-in severities;
status overrides do affect Blocking. Public seams/current reference explain the
limitation without endorsing it or rewriting normative expectations. No behavioral
probe or repair has been performed. Unanswered maintainer questions are desired
severity propagation, input-metadata ownership hardening, interruptible checker
work, full direct-writer admission and deeper analyzer coverage; those require
separate behavioral scope, not invented guarantees or conformance claims.

Applicable bounded session (awaiting parent grant): Go1.27.2 formatting plus all
267 stdlib AST/token/literal/directive and comment HTML/text identity checks;
30 native owning-module package renders; actual scripts/compatibility.sh update,
repeat stability and check while preserving rules.json/api/baseline.txt; root
native API/docs-local/offline specification checks under declared Tools2.3.0;
existing Example selector (no Example declarations are present), policy docs and
narrow generated/fixture guards only. The quickstart remains unchanged executable
example code, not a newly executed target analysis. No race/active fuzz/allocation/
corpus/performance/release/full module test campaign. All caches/helpers/renders
are task-owned and removed after this bounded session on success or failure.

## Native outcomes and remaining boundaries

Pinned Go1.27.2 ran under GOWORK=off, bounded parallelism, read-only modules and
three task-owned caches. gofmt reports no exception across all 267 tracked Go
files. All 267 executable AST/signatures, raw significant tokens/literal spellings
and directives/fixture annotations are exact against base; no optional-semicolon
normalization or producer/expression delta is needed. All 775 current comment
groups parsed/rendered as native GoDoc HTML/text. All 30 catalogued owning-module
packages rendered, including the unchanged sample coverage fixture: its absent
package commentary is intentional modeled input, not a production API omission.
Canonical fixture-role documentation belongs in the hand-maintained harness.

| Native command / boundary | Outcome |
| --- | --- |
| scripts/compatibility.sh update twice, then check | Passed; genuine 25-section public GoDoc snapshot refreshed stably. compat/rules.json and api/baseline.txt remain byte-identical. |
| Declared authenticated Tools2.3.0 api check | Root compatibility passed; coverage fixture module API explicitly not applicable. |
| Declared Tools2.3.0 docs check --local | Passed, including the native policy TestDocumentation selector. |
| Declared Tools2.3.0 specification check | Passed offline: one module, five unchanged SARIF decisions. |
| go test -run '^Example' -count=1 -timeout=2m ./... | Passed compilation/default vet; no Example declarations exist and no runnable Example result is claimed. Quickstart command code is unchanged. |
| go test ./analysistestkit -run '^(TestEveryAnalyzerHasCrossBuildFixtureEvidence\|TestEveryAnalyzerHasSemanticNearMissFixture\|TestShippedAnalyzersDeclareNoFacts)$' -count=1 -timeout=2m | Passed existing fixture/near-miss/no-fact guards, not a new full analyzer precision or concurrency qualification. |
| go test ./internal/driver -run '^TestRun(ExcludesOnlyRecognizedGeneratedFilesWhenConfigured\|RequiresGeneratedHeaderPolicyAndExactPath\|RejectsSuppressionBehindForgedGeneratedHeader)$' -count=1 -timeout=2m | Passed existing exact header/path trust boundaries; fixture source, literal headers and expected diagnostics are unchanged. |

All 165 fixture files remain byte-identical. Dependency graphs, config, catalogs,
actual generated/header/suppression directives, want annotations, rules inventory,
API baseline, corpus/precision/performance inputs, specification pins/manifests,
typed decisions/history and native scripts/CI remain unchanged. No runtime repair,
behavioral probe, race/leak/active fuzz/mutation/allocation/performance/corpus/
service/full-module qualification or release campaign was run.

Disposable evidence mechanics initially overmatched a GoDoc mention of suppression
syntax, then misread a nested catalog directory as repository-relative. Both were
corrected in owned helpers only; unchanged source proof/29 root renders were
reused and the remaining owning-module fixture render completed. Read-only module
cache directories initially refused ordinary deletion; exact owned targets were
made removable immediately and every cache was removed, including all failure
paths. No shared cache or user resource was touched. Final three caches and all
helpers/logs/renders/saved base are removed after capturing evidence; parent-owned
verifiers/scanner remain intact. Source structure, changed-content/staged secrets,
links and ordinary complete final-diff inspection conclude the Tier A batch.
Publication, current-head hosted CI and enforcement remain coordinator-owned.

Ordinary full-diff review found two source-documentation mistakes before commit:
backend-error passthrough argument indices exclude SSA method receivers (unlike
call-expression sink indices), and a nil importboundary.Options.Policies disables
only explicit bans, not configured restrictions/directions. Three affected source
files were corrected without changing runtime. Prior runtime/fixture/example/native
gate evidence remains unchanged. The focused final refresh passed for all three
affected files: AST/signatures/raw tokens/literals/directives exact, 174 comment
groups parsed/HTML rendered, analysis/backenderror/importboundary native GoDocs
inspected, actual compatibility update twice/check passed with stable 25-section
output and frozen rules/API. The other 27 renders and prior applicable native
results are reused for unchanged inputs, not represented as newly executed.
Final three caches, helper/base/logs/renders were removed and absence verified.

Final ordinary complete source-to-base diff review reconciled the corrected
indexing and independent import-policy dimensions; no material finding remains.
Qualified Tools2.4.1 structural check reports zero errors, 23 modeled-header
reviews and three context.TODO semantic marker reviews (two fixture calls and
one truthful package-doc mention), not missing real producer/debt failures.
All 13 changed Markdown local targets resolve. Final explicit staging and pinned
added-content/new-file secret scans protect the coherent documentation commit;
unchanged historical fixtures and normative evidence are not rewritten to make
those checks pass.

## Objective Enforcement

The checksum-pinned public Tools v2.4.0 verifier now supplies
`make source-docs-check` as a prerequisite of the existing local and CI Make
routes. Ordinary shared CI selects strict source-documentation checking and
requires its result. Native prerequisites, operation budgets, security and
release controls are retained. Hosted exact-head evidence is recorded by the
coordination ledger; configuration alone does not establish green delivery.
