# Release process

## Local candidate verification

Install the exact stable patch release recorded in `.go-version`. `make check`
delegates to the pinned shared tooling contract, which checks the active Go
toolchain before analysis. `make workflows` validates the pinned workflow.

Choose a semantic version without a leading `v`, update the changelog and any
rule `introduced_version` metadata, then run:

```sh
make check
make ci
golib release check
```

`./scripts/verify-release.sh <version>` packages the candidate twice and
compares every byte. It checks six CGO-disabled targets: Linux, macOS, and
Windows on amd64 and arm64.
Every ZIP contains only the versioned directory, executable, README, changelog,
and security policy. It validates sorted SHA-256 checksums, exact archive
contents, and the host executable's reported version.

To inspect a candidate without publishing it, use an empty output directory:

```sh
golib release dry-run
```

The release command refuses invalid semantic versions and non-empty output
directories. Builds use `CGO_ENABLED=0`, `-trimpath`, `-buildvcs=false`, and a
linker-injected version. Source file timestamps and ZIP metadata are normalized
before checksums are created.

## Tag publication

After local verification, independent review, and green exact-source CI, a
maintainer may create a signed `vX.Y.Z` tag pointing at the reviewed main
commit and publish a GitHub release. The current repository workflow runs on
pull requests, main pushes, schedules, and manual dispatch; it does not run on
tags or publish archives. A release is source-only unless a separate, reviewed
publication process builds and uploads the verified archives and their
`checksums.txt`. Do not describe locally built archives as CI-attested assets.

Signing is owned by the maintainer and publishing environment; the repository
does not manufacture signing identity. Consumers should verify the trusted tag
signature and public Go module checksum. If archives are actually published,
consumers should also verify the corresponding archive checksum.

## Rollback and replacement

Published tags and any attached artifacts are immutable. If a candidate is
wrong, publish a new patch version with a changelog entry. Do not replace
archives or checksums under an existing version. A withdrawn release may be
marked clearly, but its tag and any artifacts remain available for audit
unless security response requires a documented exception.
