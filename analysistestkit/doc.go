// Package analysistestkit owns internal test-only analyzer robustness, fixture
// precision and benchmark harnesses; it exports no execution API. The harness
// parses/types modeled source but does not execute modeled target statements.
//
// Generated headers in analyzer testdata and driver test source strings are
// intentional hand-maintained inputs to ast.IsGenerated and exclusion policy,
// not claims that production output needs an available source generator. Their
// exact header, build-tag, want and suppression semantics must remain intact.
// Fixture tests require diagnostics even in generated/tagged inputs; driver
// exclusions additionally require an exact trusted path. Semantic near-miss
// fixtures model aliases, closures, embedding, generics and interfaces without
// claiming a complete language or false-positive oracle. No regeneration or
// modeled go:generate execution is required or authorized by these headers.
package analysistestkit
