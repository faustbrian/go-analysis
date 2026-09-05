package analysis

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigRejectsPathResolutionFailure(t *testing.T) {
	t.Parallel()

	_, err := loadConfig("analysis.yml", nil, func(string) (string, error) {
		return "", errors.New("working directory unavailable")
	}, readConfiguration)
	if err == nil || !strings.Contains(err.Error(), "resolve configuration path") {
		t.Fatalf("loadConfig() error = %v, want path resolution error", err)
	}
}

func TestReadConfigurationAcceptsExactSizeLimit(t *testing.T) {
	t.Parallel()

	contents := append(
		[]byte("version: 1\n#"),
		bytes.Repeat([]byte{'x'}, maxConfigurationBytes-len("version: 1\n#"))...,
	)
	path := filepath.Join(t.TempDir(), "analysis.yml")
	if err := os.WriteFile(path, contents, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	got, err := readConfiguration(path)
	if err != nil {
		t.Fatalf("readConfiguration() error = %v", err)
	}
	if len(got) != maxConfigurationBytes {
		t.Fatalf("len(readConfiguration()) = %d, want %d",
			len(got), maxConfigurationBytes)
	}
}

func TestReadConfigurationRejectsCanceledContextBeforeOpen(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := readConfigurationContext(ctx, filepath.Join(t.TempDir(), "missing.yml"))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("readConfiguration() error = %v, want cancellation", err)
	}
}

func TestReadConfigurationContentsStopsAfterCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	reader := &cancelingReader{cancel: cancel}
	_, err := readConfigurationContents(readerWithContext(ctx, reader))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("readConfigurationContents() error = %v, want cancellation", err)
	}
	if reader.reads != 1 {
		t.Fatalf("reader calls = %d, want 1", reader.reads)
	}
}

func TestReadConfigurationContentsDoesNotReadCanceledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	reader := &cancelingReader{cancel: cancel}
	_, err := readConfigurationContents(readerWithContext(ctx, reader))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("readConfigurationContents() error = %v, want cancellation", err)
	}
	if reader.reads != 0 {
		t.Fatalf("reader calls = %d, want 0", reader.reads)
	}
}

type cancelingReader struct {
	cancel context.CancelFunc
	reads  int
}

func (reader *cancelingReader) Read(buffer []byte) (int, error) {
	reader.reads++
	reader.cancel()
	return copy(buffer, "version: 1\n"), nil
}

func TestGeneratedPolicyEnforcesPathLimit(t *testing.T) {
	t.Parallel()

	paths := make([]string, maxGeneratedPaths)
	for index := range paths {
		paths[index] = fmt.Sprintf("generated/file-%d.go", index)
	}
	config := Config{
		Version: 1,
		Generated: GeneratedPolicy{
			Exclude: true,
			Paths:   paths,
		},
	}
	if err := config.Validate(nil); err != nil {
		t.Fatalf("Validate(exact generated path limit) error = %v", err)
	}
	config.Generated.Paths = append(config.Generated.Paths, "generated/overflow.go")
	if err := config.Validate(nil); err == nil {
		t.Fatal("Validate() accepted generated paths above limit")
	}
}

func TestConfigAllowsExceptionsDifferingByRuleOrPackage(t *testing.T) {
	t.Parallel()

	config := Config{
		Version: 1,
		Exceptions: []PolicyException{
			{
				Rule: "security/no-unsafe", Package: "example.com/service/a",
				Path: "internal/bridge.go", Reason: "reviewed bridge",
			},
			{
				Rule: "context/no-background", Package: "example.com/service/a",
				Path: "internal/bridge.go", Reason: "reviewed worker",
			},
			{
				Rule: "security/no-unsafe", Package: "example.com/service/b",
				Path: "internal/bridge.go", Reason: "reviewed bridge",
			},
		},
	}
	if err := config.Validate([]string{
		"security/no-unsafe",
		"context/no-background",
	}); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
