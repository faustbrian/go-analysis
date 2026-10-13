// Package version provides the release version shared by reports and commands.
package version

// Value is overridden through -ldflags when building a release artifact. It is
// a mutable provenance label, not evidence of a tag, signature or compatibility;
// do not mutate it concurrently with report/command use.
var Value = "0.1.0"
