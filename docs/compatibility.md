# Compatibility policy

The public compatibility surface includes rule IDs and trigger semantics,
metadata and owners, configuration keys and validation, suppression syntax,
JSON and SARIF fields, command arguments and exit codes, and exported Go APIs.
The [SARIF specification decision register](specification-decisions.md)
governs the standard-backed portion of that surface.

Semantic versioning governs releases: additive
rules and optional fields are minor changes, compatible fixes are patch changes,
and removals, renames, stricter existing policy, or report/schema breakage require
a major version.

`golib api check` compares the current exported documentation and complete rule
inventory with reviewed files under `compat`. An intentional change uses
`./scripts/compatibility.sh update`; reviewers inspect both snapshots and the
changelog.
Snapshot regeneration does not itself make a change compatible.

The minimum Go version is the `go` directive in `go.mod`. A release is tested
with that toolchain contract. The standalone configured command and vettool are
supported; a golangci-lint module plugin remains optional and is not part of the
compatibility surface until its lifecycle is proven stable and reproducible.
