package analysis

import "testing"

// FuzzDecodeConfig bounds arbitrary YAML to 1 MiB and exercises strict decoding
// for panic freedom; it does not assert acceptance or a YAML conformance oracle.
func FuzzDecodeConfig(f *testing.F) {
	f.Add("version: 1\n")
	f.Add("version: 1\nunknown: true\n")
	f.Add("version: 1\n---\nversion: 1\n")
	f.Add("version: [\n")

	f.Fuzz(func(t *testing.T, contents string) {
		if len(contents) > maxConfigurationBytes {
			t.Skip()
		}
		_, _ = decodeConfig(
			[]byte(contents),
			"/repository",
			[]string{"security/no-unsafe"},
		)
	})
}
