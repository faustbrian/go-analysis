// Command quickstart demonstrates context-aware policy loading with the public
// analysis and policy packages.
package main

import (
	"context"
	"fmt"
	"io"
	"os"

	shared "github.com/faustbrian/go-analysis/analysis"
	"github.com/faustbrian/go-analysis/policy"
)

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdout); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, arguments []string, output io.Writer) error {
	if len(arguments) != 1 {
		return fmt.Errorf("usage: quickstart <analysis.yml>")
	}
	registry, err := policy.Builtin()
	if err != nil {
		return fmt.Errorf("build rule registry: %w", err)
	}
	config, err := shared.LoadConfigContext(ctx, arguments[0], registry.IDs())
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(
		output,
		"loaded policy version %d with %d governed rules\n",
		config.Version,
		len(registry.IDs()),
	)
	return err
}
