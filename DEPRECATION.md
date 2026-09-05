# Deprecation Policy

Deprecations MUST identify the replacement, reason, migration steps, and
earliest removal version. Public Go identifiers use a valid `Deprecated:` doc
paragraph and corresponding changelog entry.

At `v1` and later, a supported replacement SHOULD exist for at least one minor
release before removal. Security or correctness defects MAY require faster
removal when continued support would be unsafe; the release notes must explain
the exception.

Silent behavior changes, undocumented aliases, and indefinite deprecated code
are prohibited. Deprecations are checked during compatibility and release
review.

## Active deprecations

### `analysis.LoadConfig`

`LoadConfig` is deprecated because its signature cannot propagate caller
cancellation or deadlines through configuration-file I/O. Replace calls with:

```go
config, err := analysis.LoadConfigContext(ctx, path, knownRules)
```

Callers must pass the context that owns the surrounding command or operation.
`LoadConfig` retains its existing context-free behavior during migration.

Support for `LoadConfig` is frozen for the longer of 180 days after the first
stable release containing this deprecation and two subsequently published
stable minor releases. Removal is permitted only in an explicitly authorized
next major release after known consumers have migrated. No removal version is
currently scheduled.
