# SARIF conformance matrix

The root module implements the documented output profile for
[OASIS SARIF 2.1.0](https://docs.oasis-open.org/sarif/sarif/v2.1.0/os/sarif-v2.1.0-os.html).
The [specification decision register](../docs/specification-decisions.md)
defines the supported boundary; it does not claim every optional SARIF object.

[`manifest.tsv`](manifest.tsv) pins the standard, approved errata, official
schema, and maintained-peer sources. [`monitoring.json`](monitoring.json)
requires a digest review every 30 days. A changed source requires review of
every affected decision, public contract, test, compatibility statement, and
changelog entry before behavior changes.

## Decision conformance

| Decision | Authority | Executable evidence | Differential result |
| --- | --- | --- | --- |
| ANALYSIS-SARIF-DEC-001 | `oasis-sarif-source` | `TestWriteSARIFIncludesStableRulesAndNoSource`, `TestWritersEncodeEmptyInventoriesAsArrays`, `TestRunCheckEmitsSARIFForAdvisoryFinding`, `FuzzReportWriters` | Deliberate schema-mirror policy difference from gosec; version agrees with both maintained peers. |
| ANALYSIS-SARIF-DEC-002 | `oasis-sarif-source` | `TestWriteSARIFIncludesStableRulesAndNoSource`, `FuzzReportWriters` | Deliberate descriptor-rich profile between golangci-lint and gosec. |
| ANALYSIS-SARIF-DEC-003 | `oasis-sarif-source` | `TestWriteSARIFIncludesStableRulesAndNoSource`, `FuzzReportWriters` | Deliberate severity and unknown-rule fallback differences. |
| ANALYSIS-SARIF-DEC-004 | `oasis-sarif-source` | `TestWriteSARIFIncludesStableRulesAndNoSource`, `TestWritersRejectPathsThatCouldExposeWorkspace`, `TestRunCheckEmitsSARIFForAdvisoryFinding`, `FuzzReportWriters` | Deliberate source-omission policy; direct invalid positions remain outside the supported profile. |
| ANALYSIS-SARIF-DEC-005 | `oasis-sarif-source` | `TestWriteSARIFIncludesStableRulesAndNoSource`, `TestWritersEncodeEmptyInventoriesAsArrays`, `FuzzReportWriters` | Deliberate run property-bag extension rather than finding-specific result suppressions. |

The source-reviewed comparison is reproducible from
[`maintained-peers.json`](maintained-peers.json). Peer disagreement is
classified as a deliberate policy difference, not as normative authority.
