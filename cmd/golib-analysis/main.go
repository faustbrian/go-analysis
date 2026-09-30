// Command analysis runs the governed analyzer suite as a standalone vettool.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	shared "github.com/faustbrian/go-analysis/analysis"
	"github.com/faustbrian/go-analysis/analyzers/backenderror"
	"github.com/faustbrian/go-analysis/analyzers/blockingcontext"
	"github.com/faustbrian/go-analysis/analyzers/cleanupownership"
	"github.com/faustbrian/go-analysis/analyzers/constructorgoroutine"
	"github.com/faustbrian/go-analysis/analyzers/forbiddenapi"
	"github.com/faustbrian/go-analysis/analyzers/globalgoroutine"
	"github.com/faustbrian/go-analysis/analyzers/goroutinefanout"
	"github.com/faustbrian/go-analysis/analyzers/httpclienttimeout"
	"github.com/faustbrian/go-analysis/analyzers/importboundary"
	"github.com/faustbrian/go-analysis/analyzers/interfacenaming"
	"github.com/faustbrian/go-analysis/analyzers/interfaceplacement"
	"github.com/faustbrian/go-analysis/analyzers/lockacrosscall"
	"github.com/faustbrian/go-analysis/analyzers/metriccardinality"
	"github.com/faustbrian/go-analysis/analyzers/mutableglobal"
	"github.com/faustbrian/go-analysis/analyzers/nobackground"
	"github.com/faustbrian/go-analysis/analyzers/nodefaulthttp"
	"github.com/faustbrian/go-analysis/analyzers/noinit"
	"github.com/faustbrian/go-analysis/analyzers/noprocesscontrol"
	"github.com/faustbrian/go-analysis/analyzers/nostoredcontext"
	"github.com/faustbrian/go-analysis/analyzers/nounsafe"
	"github.com/faustbrian/go-analysis/analyzers/sensitivesink"
	"github.com/faustbrian/go-analysis/analyzers/transactionrollback"
	"github.com/faustbrian/go-analysis/internal/driver"
	"github.com/faustbrian/go-analysis/internal/version"
	"github.com/faustbrian/go-analysis/policy"
	toolanalysis "golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

const (
	commandPackage     = "github.com/faustbrian/go-analysis/cmd/golib-analysis"
	maxPolicySyncBytes = 1 << 20
)

func main() {
	runCommand(
		context.Background(),
		os.Args[1:],
		os.Stdout,
		os.Stderr,
		os.Exit,
		multichecker.Main,
	)
}

func runCommand(
	ctx context.Context,
	arguments []string,
	output io.Writer,
	errorOutput io.Writer,
	exit func(int),
	runAnalyzers func(...*toolanalysis.Analyzer),
) {
	handled, blocking, err := runCheck(ctx, arguments, output)
	if handled {
		if err != nil {
			exitWithCommandError(errorOutput, err, exit)
			return
		}
		if blocking {
			exit(1)
		}
		return
	}
	handled, err = runUtilityContext(ctx, arguments, output)
	if handled {
		if err != nil {
			exitWithCommandError(errorOutput, err, exit)
		}
		return
	}
	runAnalyzers(analyzers()...)
}

func exitWithCommandError(output io.Writer, err error, exit func(int)) {
	if _, writeErr := fmt.Fprintln(output, err); writeErr != nil {
		exit(2)
		return
	}
	exit(2)
}

func runCheck(
	ctx context.Context,
	arguments []string,
	output io.Writer,
) (bool, bool, error) {
	if len(arguments) == 0 || arguments[0] != "check" {
		return false, false, nil
	}
	flags := flag.NewFlagSet("check", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	configPath := flags.String("config", "", "path to analysis policy")
	format := flags.String("format", "json", "report format: json or sarif")
	root := flags.String("root", "", "absolute target repository root")
	sequential := flags.Bool("sequential", false, "disable parallel analysis")
	if err := flags.Parse(arguments[1:]); err != nil {
		return true, false, fmt.Errorf("parse check arguments: %w", err)
	}
	if *configPath == "" {
		return true, false, errors.New("check requires -config <path>")
	}
	if *format != "json" && *format != "sarif" {
		return true, false, fmt.Errorf("unsupported report format %q", *format)
	}
	result, err := driver.Run(ctx, driver.Options{
		ConfigPath: *configPath,
		Root:       *root,
		Patterns:   flags.Args(),
		Sequential: *sequential,
	})
	if err != nil {
		return true, false, err
	}
	if *format == "json" {
		err = shared.WriteJSON(output, result.Report)
	} else {
		err = shared.WriteSARIF(output, result.Report)
	}
	if err != nil {
		return true, false, err
	}

	return true, result.Blocking, nil
}

func analyzers() []*toolanalysis.Analyzer {
	return []*toolanalysis.Analyzer{
		backenderror.Analyzer,
		blockingcontext.Analyzer,
		cleanupownership.Analyzer,
		transactionrollback.Analyzer,
		constructorgoroutine.Analyzer,
		forbiddenapi.Analyzer,
		globalgoroutine.Analyzer,
		goroutinefanout.Analyzer,
		httpclienttimeout.Analyzer,
		importboundary.Analyzer,
		lockacrosscall.Analyzer,
		nobackground.New(nobackground.Options{
			AllowedPackages: []string{commandPackage},
		}),
		nodefaulthttp.Analyzer,
		noinit.New(noinit.Options{AllowedPackages: []string{commandPackage}}),
		noprocesscontrol.New(noprocesscontrol.Options{
			AllowedPackages: []string{commandPackage},
		}),
		nostoredcontext.Analyzer,
		sensitivesink.Analyzer,
		nounsafe.Analyzer,
		mutableglobal.Analyzer,
		interfacenaming.Analyzer,
		interfaceplacement.Analyzer,
		metriccardinality.Analyzer,
		metriccardinality.LabelNameAnalyzer,
	}
}

func runUtility(arguments []string, output io.Writer) (bool, error) {
	return runUtilityContext(context.Background(), arguments, output)
}

func runUtilityContext(
	ctx context.Context,
	arguments []string,
	output io.Writer,
) (bool, error) {
	return runUtilityWithRegistryContext(ctx, arguments, output, policy.Builtin)
}

func runUtilityWithRegistry(
	arguments []string,
	output io.Writer,
	registryFactory func() (*policy.Registry, error),
) (bool, error) {
	return runUtilityWithRegistryContext(
		context.Background(),
		arguments,
		output,
		registryFactory,
	)
}

func runUtilityWithRegistryContext(
	ctx context.Context,
	arguments []string,
	output io.Writer,
	registryFactory func() (*policy.Registry, error),
) (bool, error) {
	if len(arguments) == 0 {
		return false, nil
	}
	switch arguments[0] {
	case "version":
		if len(arguments) != 1 {
			return true, errors.New("usage: golib-analysis version")
		}
		if _, err := fmt.Fprintln(output, version.Value); err != nil {
			return true, fmt.Errorf("write version: %w", err)
		}
		return true, nil
	case "rules":
		if len(arguments) != 1 {
			return true, errors.New("usage: golib-analysis rules")
		}
		registry, err := registryFactory()
		if err != nil {
			return true, fmt.Errorf("build rule inventory: %w", err)
		}
		if err := json.NewEncoder(output).Encode(registry.Entries()); err != nil {
			return true, fmt.Errorf("write rule inventory: %w", err)
		}
		return true, nil
	case "validate-config":
		if len(arguments) != 2 {
			return true, errors.New("usage: golib-analysis validate-config <path>")
		}
		registry, err := registryFactory()
		if err != nil {
			return true, fmt.Errorf("build rule inventory: %w", err)
		}
		config, err := shared.LoadConfigContext(ctx, arguments[1], registry.IDs())
		if err != nil {
			return true, err
		}
		if err := driver.Validate(config); err != nil {
			return true, err
		}
		if _, err := fmt.Fprintln(output, "configuration valid"); err != nil {
			return true, fmt.Errorf("write validation result: %w", err)
		}
		return true, nil
	case "sync-policy":
		return runPolicySync(ctx, arguments, output, registryFactory)
	default:
		return false, nil
	}
}

func runPolicySync(
	ctx context.Context,
	arguments []string,
	output io.Writer,
	registryFactory func() (*policy.Registry, error),
) (bool, error) {
	return runPolicySyncWithDependenciesContext(
		ctx,
		arguments,
		output,
		registryFactory,
		policySyncDependencies{
			loadConfig: shared.LoadConfigContext,
			validate:   driver.Validate,
			readFile: func(path string) ([]byte, error) {
				return readPolicySyncFile(ctx, path)
			},
			writeFile: os.WriteFile,
		},
	)
}

func readPolicySyncFile(ctx context.Context, path string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	file, err := os.Open(path) // #nosec G304 G703 -- operator-selected policy file; this command is not a path sandbox
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	return readPolicySyncContents(ctx, file)
}

func readPolicySyncContents(ctx context.Context, reader io.Reader) ([]byte, error) {
	contents, err := io.ReadAll(io.LimitReader(reader, maxPolicySyncBytes+1))
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(contents) > maxPolicySyncBytes {
		return nil, fmt.Errorf("policy snapshot exceeds %d bytes", maxPolicySyncBytes)
	}
	return contents, nil
}

type policySyncDependencies struct {
	loadConfig func(context.Context, string, []string) (*shared.Config, error)
	validate   func(*shared.Config) error
	readFile   func(string) ([]byte, error)
	writeFile  func(string, []byte, os.FileMode) error
}

func runPolicySyncWithDependencies(
	arguments []string,
	output io.Writer,
	registryFactory func() (*policy.Registry, error),
	dependencies policySyncDependencies,
) (bool, error) {
	return runPolicySyncWithDependenciesContext(
		context.Background(),
		arguments,
		output,
		registryFactory,
		dependencies,
	)
}

func runPolicySyncWithDependenciesContext(
	ctx context.Context,
	arguments []string,
	output io.Writer,
	registryFactory func() (*policy.Registry, error),
	dependencies policySyncDependencies,
) (bool, error) {
	if len(arguments) != 4 ||
		(arguments[1] != "check" && arguments[1] != "update") {
		return true, errors.New(
			"usage: golib-analysis sync-policy <check|update> <canonical> <local>",
		)
	}
	registry, err := registryFactory()
	if err != nil {
		return true, fmt.Errorf("build rule inventory: %w", err)
	}
	config, err := dependencies.loadConfig(ctx, arguments[2], registry.IDs())
	if err != nil {
		return true, fmt.Errorf("validate canonical policy: %w", err)
	}
	if err := dependencies.validate(config); err != nil {
		return true, fmt.Errorf("validate canonical policy: %w", err)
	}
	canonical, err := dependencies.readFile(arguments[2])
	if err != nil {
		return true, fmt.Errorf("read canonical policy: %w", err)
	}
	if len(canonical) > maxPolicySyncBytes {
		return true, fmt.Errorf("canonical policy exceeds %d bytes", maxPolicySyncBytes)
	}
	snapshot, err := shared.ParseConfigBytes(arguments[2], canonical, registry.IDs())
	if err != nil {
		return true, fmt.Errorf("validate canonical policy snapshot: %w", err)
	}
	if err := dependencies.validate(snapshot); err != nil {
		return true, fmt.Errorf("validate canonical policy snapshot: %w", err)
	}
	if arguments[1] == "update" {
		if err := ctx.Err(); err != nil {
			return true, err
		}
		if err := dependencies.writeFile(arguments[3], canonical, 0o644); err != nil {
			return true, fmt.Errorf("write local policy: %w", err)
		}
		if _, err := fmt.Fprintln(output, "policy synchronized"); err != nil {
			return true, fmt.Errorf("write synchronization result: %w", err)
		}
		return true, nil
	}
	local, err := dependencies.readFile(arguments[3])
	if err != nil {
		return true, fmt.Errorf("read local policy: %w", err)
	}
	if len(local) > maxPolicySyncBytes {
		return true, fmt.Errorf("local policy exceeds %d bytes", maxPolicySyncBytes)
	}
	if !bytes.Equal(canonical, local) {
		return true, errors.New(
			"policy drift: local policy differs from canonical policy",
		)
	}
	if _, err := fmt.Fprintln(output, "policy in sync"); err != nil {
		return true, fmt.Errorf("write synchronization result: %w", err)
	}
	return true, nil
}
