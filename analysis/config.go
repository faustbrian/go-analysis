package analysis

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"
)

const (
	maxConfigurationBytes = 1 << 20
	maxGeneratedPaths     = 4096
)

// Config is a mutable organization policy snapshot. Version 1 is required;
// Validate checks shared structure, while analyzer constructors check typed
// policies. Callers own nested maps/slices and must prevent concurrent mutation.
type Config struct {
	// Version selects schema version 1; the zero configuration is invalid.
	Version int `yaml:"version" json:"version"`
	// Root is the absolute directory derived from the config path, not a YAML field.
	Root string `yaml:"-" json:"root"`
	// Packages declares import classification, boundaries and blocking API policies.
	Packages []PackagePolicy `yaml:"packages,omitempty" json:"packages,omitempty"`
	// Entrypoints lists exact packages exempt from root-context and process-control
	// rules.
	Entrypoints []string `yaml:"entrypoints,omitempty" json:"entrypoints,omitempty"`
	// InitPackages lists exact packages permitted to declare init functions.
	InitPackages []string `yaml:"init_packages,omitempty" json:"init_packages,omitempty"`
	// ContextOwners lists exact packages permitted to store context implementations.
	ContextOwners []string `yaml:"context_owners,omitempty" json:"context_owners,omitempty"`
	// HTTPTimeoutExceptions lists exact packages exempt from HTTP timeout
	// construction checks.
	HTTPTimeoutExceptions []string `yaml:"http_timeout_exceptions,omitempty" json:"http_timeout_exceptions,omitempty"`
	// Constructors identifies declarations checked for syntactically executed
	// goroutine startup.
	Constructors []ConstructorPolicy `yaml:"constructors,omitempty" json:"constructors,omitempty"`
	// ResourceConstructors identifies call results that must not be discarded.
	ResourceConstructors []ResourceConstructorPolicy `yaml:"resource_constructors,omitempty" json:"resource_constructors,omitempty"`
	// Transactions identifies constructors requiring recognized immediate rollback
	// ownership.
	Transactions []TransactionPolicy `yaml:"transactions,omitempty" json:"transactions,omitempty"`
	// LockSensitiveCalls selects resolved callback or IO calls forbidden under
	// definitely-held locks.
	LockSensitiveCalls []LockSensitiveCallPolicy `yaml:"lock_sensitive_calls,omitempty" json:"lock_sensitive_calls,omitempty"`
	// SensitiveTypes identifies named secret-bearing types, not untyped secret
	// values.
	SensitiveTypes []SensitiveTypePolicy `yaml:"sensitive_types,omitempty" json:"sensitive_types,omitempty"`
	// SensitiveSinks selects typed disclosure argument positions.
	SensitiveSinks []SensitiveSinkPolicy `yaml:"sensitive_sinks,omitempty" json:"sensitive_sinks,omitempty"`
	// ForbiddenAPIs declares callable migration policies and remediation text.
	ForbiddenAPIs []ForbiddenAPIPolicy `yaml:"forbidden_apis,omitempty" json:"forbidden_apis,omitempty"`
	// BackendClients restricts dependency trees to approved adapter packages.
	BackendClients []BackendClientPolicy `yaml:"backend_clients,omitempty" json:"backend_clients,omitempty"`
	// MutableGlobals selects package trees whose composite globals are governed.
	MutableGlobals []MutableGlobalPolicy `yaml:"mutable_globals,omitempty" json:"mutable_globals,omitempty"`
	// InterfaceProviderPackages selects provider trees forbidden to declare exported
	// value interfaces.
	InterfaceProviderPackages []string `yaml:"interface_provider_packages,omitempty" json:"interface_provider_packages,omitempty"`
	// InterfaceNames selects exported value-interface name affixes and exemptions.
	InterfaceNames []InterfaceNamePolicy `yaml:"interface_names,omitempty" json:"interface_names,omitempty"`
	// MetricLabelTypes identifies named high-cardinality label-value types.
	MetricLabelTypes []MetricLabelTypePolicy `yaml:"metric_label_types,omitempty" json:"metric_label_types,omitempty"`
	// MetricLabelSinks selects metric call argument positions for label-value
	// checks.
	MetricLabelSinks []MetricLabelSinkPolicy `yaml:"metric_label_sinks,omitempty" json:"metric_label_sinks,omitempty"`
	// MetricLabelNameTypes identifies named untrusted label-name types.
	MetricLabelNameTypes []MetricLabelNameTypePolicy `yaml:"metric_label_name_types,omitempty" json:"metric_label_name_types,omitempty"`
	// MetricLabelNameSinks selects metric call argument positions for label-name
	// checks.
	MetricLabelNameSinks []MetricLabelSinkPolicy `yaml:"metric_label_name_sinks,omitempty" json:"metric_label_name_sinks,omitempty"`
	// BackendErrorBoundaries selects exported package error boundaries.
	BackendErrorBoundaries []string `yaml:"backend_error_boundaries,omitempty" json:"backend_error_boundaries,omitempty"`
	// BackendErrorSources identifies exact backend result positions for bounded SSA
	// tracing.
	BackendErrorSources []BackendErrorFlowPolicy `yaml:"backend_error_sources,omitempty" json:"backend_error_sources,omitempty"`
	// BackendErrorPassthroughs declares explicitly trusted wrapper result-to-
	// argument flow.
	BackendErrorPassthroughs []BackendErrorPassthroughPolicy `yaml:"backend_error_passthroughs,omitempty" json:"backend_error_passthroughs,omitempty"`
	// GoroutineFanout selects package trees and recognized static loop-launch
	// ceilings.
	GoroutineFanout []GoroutineFanoutPolicy `yaml:"goroutine_fanout,omitempty" json:"goroutine_fanout,omitempty"`
	// AdapterRoots is reserved and must be nil; even an empty non-nil list is
	// rejected.
	AdapterRoots []string `yaml:"adapter_roots,omitempty" json:"adapter_roots,omitempty"`
	// Layers declares named architectural import directions.
	Layers []DirectionPolicy `yaml:"layers,omitempty" json:"layers,omitempty"`
	// Contexts declares named bounded-context import directions.
	Contexts []DirectionPolicy `yaml:"contexts,omitempty" json:"contexts,omitempty"`
	// Generated authorizes exact header-bearing file exclusions at driver report
	// collection.
	Generated GeneratedPolicy `yaml:"generated,omitempty" json:"generated,omitempty"`
	// Rules overrides known rule status/severity; nil retains metadata defaults.
	Rules map[string]RulePolicy `yaml:"rules,omitempty" json:"rules,omitempty"`
	// Exceptions declares attributable exact rule/package and optional file
	// exceptions.
	Exceptions []PolicyException `yaml:"exceptions,omitempty" json:"exceptions,omitempty"`
}

// ConstructorPolicy identifies exact constructors with explicit lifecycle rules.
type ConstructorPolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
	// Symbols lists configured Function or Type.Method names.
	Symbols []string `yaml:"symbols" json:"symbols"`
}

// ResourceConstructorPolicy identifies an explicit cleanup ownership contract.
type ResourceConstructorPolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
	// Symbol names a Function or Type.Method; matching uses the resolved Go object.
	Symbol string `yaml:"symbol" json:"symbol"`
	// CleanupResult selects the zero-based cleanup result; binding it does not prove
	// eventual cleanup.
	CleanupResult int `yaml:"cleanup_result" json:"cleanup_result"`
	// AllowedPackages lists exact importing/calling package exemptions; nil grants
	// none.
	AllowedPackages []string `yaml:"allowed_packages,omitempty" json:"allowed_packages,omitempty"`
}

// TransactionPolicy identifies one transaction constructor and its rollback
// ownership contract.
type TransactionPolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
	// Symbol names a Function or Type.Method; matching uses the resolved Go object.
	Symbol string `yaml:"symbol" json:"symbol"`
	// Result selects the zero-based result position; zero selects the first result.
	Result int `yaml:"result" json:"result"`
	// RollbackMethod names the instance method that must establish rollback
	// ownership.
	RollbackMethod string `yaml:"rollback_method" json:"rollback_method"`
}

// LockSensitiveCallPolicy identifies a callback or I/O call forbidden under locks.
type LockSensitiveCallPolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
	// Symbol names a Function or Type.Method; matching uses the resolved Go object.
	Symbol string `yaml:"symbol" json:"symbol"`
	// AllowedPackages lists exact importing/calling package exemptions; nil grants
	// none.
	AllowedPackages []string `yaml:"allowed_packages,omitempty" json:"allowed_packages,omitempty"`
}

// SensitiveTypePolicy identifies one named secret-bearing type.
type SensitiveTypePolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
	// Name names the exact declared Go type, with aliases resolved during analysis.
	Name string `yaml:"name" json:"name"`
}

// SensitiveSinkPolicy identifies typed-sensitive positions on one callable.
type SensitiveSinkPolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
	// Symbol names a Function or Type.Method; matching uses the resolved Go object.
	Symbol string `yaml:"symbol" json:"symbol"`
	// Arguments lists zero-based call-expression argument positions; method-
	// expression receivers participate.
	Arguments []int `yaml:"arguments,omitempty" json:"arguments,omitempty"`
	// VariadicFrom selects all argument positions at or after this non-negative
	// index; nil disables the tail.
	VariadicFrom *int `yaml:"variadic_from,omitempty" json:"variadic_from,omitempty"`
	// AllowedPackages lists exact importing/calling package exemptions; nil grants
	// none.
	AllowedPackages []string `yaml:"allowed_packages,omitempty" json:"allowed_packages,omitempty"`
}

// MetricLabelTypePolicy identifies one named high-cardinality label type.
type MetricLabelTypePolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
	// Name names the exact declared Go type, with aliases resolved during analysis.
	Name string `yaml:"name" json:"name"`
}

// MetricLabelNameTypePolicy identifies one named attacker-controlled label-name
// type.
type MetricLabelNameTypePolicy = MetricLabelTypePolicy

// MetricLabelSinkPolicy identifies label positions on one metric callable.
type MetricLabelSinkPolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
	// Symbol names a Function or Type.Method; matching uses the resolved Go object.
	Symbol string `yaml:"symbol" json:"symbol"`
	// Arguments lists zero-based call-expression argument positions; method-
	// expression receivers participate.
	Arguments []int `yaml:"arguments,omitempty" json:"arguments,omitempty"`
	// VariadicFrom selects all argument positions at or after this non-negative
	// index; nil disables the tail.
	VariadicFrom *int `yaml:"variadic_from,omitempty" json:"variadic_from,omitempty"`
}

// BackendErrorFlowPolicy identifies one backend callable result.
type BackendErrorFlowPolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
	// Symbol names a Function or Type.Method; matching uses the resolved Go object.
	Symbol string `yaml:"symbol" json:"symbol"`
	// Result selects the zero-based result position; zero selects the first result.
	Result int `yaml:"result" json:"result"`
}

// BackendErrorPassthroughPolicy identifies a wrapper that preserves errors.
type BackendErrorPassthroughPolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
	// Symbol names a Function or Type.Method; matching uses the resolved Go object.
	Symbol string `yaml:"symbol" json:"symbol"`
	// Result selects the zero-based result position; zero selects the first result.
	Result int `yaml:"result" json:"result"`
	// Arguments lists zero-based explicit parameter positions; the SSA method
	// receiver is excluded, including method-expression calls.
	Arguments []int `yaml:"arguments,omitempty" json:"arguments,omitempty"`
	// VariadicFrom selects all argument positions at or after this non-negative
	// index; nil disables the tail.
	VariadicFrom *int `yaml:"variadic_from,omitempty" json:"variadic_from,omitempty"`
}

// GoroutineFanoutPolicy configures a package tree's static fan-out ceiling.
type GoroutineFanoutPolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
	// MaxStatic is the permitted statically expanded fan-out, from 1 through 1024;
	// zero is invalid.
	MaxStatic int `yaml:"max_static" json:"max_static"`
}

// ForbiddenAPIPolicy configures one exact organization API migration.
type ForbiddenAPIPolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
	// Symbol names a Function or Type.Method; matching uses the resolved Go object.
	Symbol string `yaml:"symbol" json:"symbol"`
	// Replacement is non-empty remediation text, not an automatically applied code
	// fix.
	Replacement string `yaml:"replacement" json:"replacement"`
	// AllowedPackages lists exact importing/calling package exemptions; nil grants
	// none.
	AllowedPackages []string `yaml:"allowed_packages,omitempty" json:"allowed_packages,omitempty"`
}

// BackendClientPolicy restricts a backend package tree to approved adapters.
type BackendClientPolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
	// AllowedPackages lists approved exact import paths or trailing /... trees;
	// analyzer construction requires at least one.
	AllowedPackages []string `yaml:"allowed_packages" json:"allowed_packages"`
}

// MutableGlobalPolicy governs composite globals in one package tree.
type MutableGlobalPolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
}

// InterfaceNamePolicy governs exported value-interface names in one package tree.
type InterfaceNamePolicy struct {
	// Package is the configured import path, not a filesystem path.
	Package string `yaml:"package" json:"package"`
	// RequiredPrefix is the required case-sensitive name prefix; empty imposes none.
	RequiredPrefix string `yaml:"required_prefix,omitempty" json:"required_prefix,omitempty"`
	// RequiredSuffix is the required case-sensitive name suffix; empty imposes none.
	RequiredSuffix string `yaml:"required_suffix,omitempty" json:"required_suffix,omitempty"`
	// AllowedNames lists exact exported interface names exempt from affix checks.
	AllowedNames []string `yaml:"allowed_names,omitempty" json:"allowed_names,omitempty"`
}

// PackagePolicy assigns architecture constraints to matching import paths.
type PackagePolicy struct {
	// Pattern selects an exact import path or trailing /... package tree.
	Pattern string `yaml:"pattern" json:"pattern"`
	// Layer names a configured layer; empty leaves layer classification unset.
	Layer string `yaml:"layer,omitempty" json:"layer,omitempty"`
	// Context names a configured bounded context; empty leaves context
	// classification unset.
	Context string `yaml:"context,omitempty" json:"context,omitempty"`
	// AllowImports lists exceptions to this policy's denied imports, not a general
	// import allowlist.
	AllowImports []string `yaml:"allow_imports,omitempty" json:"allow_imports,omitempty"`
	// DenyImports lists forbidden exact imports or trailing /... trees.
	DenyImports []string `yaml:"deny_imports,omitempty" json:"deny_imports,omitempty"`
	// BlockingFunctions lists exported declarations requiring a context parameter.
	BlockingFunctions []string `yaml:"blocking_functions,omitempty" json:"blocking_functions,omitempty"`
}

// DirectionPolicy names one architecture partition and permitted dependencies.
type DirectionPolicy struct {
	// Name names the architectural partition used by package classifications.
	Name string `yaml:"name" json:"name"`
	// MayImport lists permitted different partitions; imports within the same name
	// are allowed.
	MayImport []string `yaml:"may_import,omitempty" json:"may_import,omitempty"`
}

// GeneratedPolicy controls exact trusted generated Go file exclusions.
type GeneratedPolicy struct {
	// Exclude enables exclusions only for explicit Paths whose parsed AST has a
	// generated header.
	Exclude bool `yaml:"exclude,omitempty" json:"exclude,omitempty"`
	// Paths lists at most 4096 exact repository-relative .go paths; no globs or
	// traversal.
	Paths []string `yaml:"paths,omitempty" json:"paths,omitempty"`
}

// RulePolicy overrides metadata defaults for one known rule.
type RulePolicy struct {
	// Status overrides enforcement mode; empty retains the default, blocking
	// requires Promotion.
	Status Status `yaml:"status,omitempty" json:"status,omitempty"`
	// Severity selects analyzer metadata severity; empty retains the default.
	// The current driver report inventory retains built-in severities, so this
	// override is not propagated to its JSON rule metadata or SARIF levels.
	Severity Severity `yaml:"severity,omitempty" json:"severity,omitempty"`
	// Options must be nil; use typed top-level policies instead of this unsupported
	// map.
	Options map[string]any `yaml:"options,omitempty" json:"options,omitempty"`
	// Promotion is required exactly for an explicit blocking override; evidence is
	// attribution, not executed proof.
	Promotion *PromotionPolicy `yaml:"promotion,omitempty" json:"promotion,omitempty"`
}

// PromotionPolicy attributes an explicit versioned blocking override; its
// evidence is validated as non-empty text, not authenticated proof.
type PromotionPolicy struct {
	// Version is the required numeric major.minor.patch promotion version.
	Version string `yaml:"version" json:"version"`
	// Evidence is non-empty attribution text; no URL or artifact authentication is
	// performed.
	Evidence string `yaml:"evidence" json:"evidence"`
}

// PolicyException is a narrow, reviewed exception to one rule.
type PolicyException struct {
	// Rule is the exact known rule ID to filter.
	Rule string `yaml:"rule" json:"rule"`
	// Package is the exact clean import path, not a package pattern.
	Package string `yaml:"package" json:"package"`
	// Path optionally narrows the exception to a clean repository-relative file.
	Path string `yaml:"path,omitempty" json:"path,omitempty"`
	// Reason is required explanatory text and is emitted in report inventories.
	Reason string `yaml:"reason" json:"reason"`
	// Expires is an optional YYYY-MM-DD UTC date, valid through that entire day.
	Expires string `yaml:"expires,omitempty" json:"expires,omitempty"`
	// Issue is optional attribution text; no issue lookup is performed.
	Issue string `yaml:"issue,omitempty" json:"issue,omitempty"`
	// Used marks consumption; caller-supplied true is retained by
	// ApplyPolicyExceptions.
	Used bool `yaml:"-" json:"used"`
}

// LoadConfig reads exactly one strict YAML policy document.
//
// Deprecated: use LoadConfigContext for cooperative filesystem cancellation.
func LoadConfig(path string, knownRules []string) (*Config, error) {
	return loadConfig(path, knownRules, filepath.Abs, readConfiguration)
}

// LoadConfigContext reads exactly one strict YAML policy document. It observes
// cancellation before and around file reads and does not retain ctx. A blocked
// os.Open/Read and buffered YAML decoding are not preempted. The file is closed
// on return; up to 1 MiB of configuration is admitted. ctx must be non-nil.
func LoadConfigContext(
	ctx context.Context,
	path string,
	knownRules []string,
) (*Config, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}
	return loadConfig(
		path,
		knownRules,
		filepath.Abs,
		func(path string) ([]byte, error) {
			return readConfigurationContext(ctx, path)
		},
	)
}

func loadConfig(
	path string,
	knownRules []string,
	resolve func(string) (string, error),
	read func(string) ([]byte, error),
) (*Config, error) {
	absolute, err := resolve(path)
	if err != nil {
		return nil, fmt.Errorf("resolve configuration path: %w", err)
	}
	contents, err := read(absolute)
	if err != nil {
		return nil, err
	}

	return ParseConfigBytes(absolute, contents, knownRules)
}

// ParseConfigBytes validates one bounded strict YAML policy snapshot. The path
// determines the root for repository-relative configuration and is not read.
// Input above 1 MiB, unknown YAML fields, multiple documents or invalid shared
// policy are rejected with nil configuration. YAML decoding can allocate before
// semantic validation; this byte cap is not a total memory/CPU budget. Contents
// and knownRules are borrowed only during the call; the decoded records are new.
func ParseConfigBytes(path string, contents []byte, knownRules []string) (*Config, error) {
	return parseConfigBytes(path, contents, knownRules, filepath.Abs)
}

func parseConfigBytes(
	path string,
	contents []byte,
	knownRules []string,
	resolve func(string) (string, error),
) (*Config, error) {
	if len(contents) > maxConfigurationBytes {
		return nil, fmt.Errorf("configuration exceeds %d bytes", maxConfigurationBytes)
	}
	absolute, err := resolve(path)
	if err != nil {
		return nil, fmt.Errorf("resolve configuration path: %w", err)
	}
	return decodeConfig(contents, filepath.Dir(absolute), knownRules)
}

func readConfiguration(path string) ([]byte, error) {
	file, err := os.Open(path) // #nosec G304 -- explicit policy path
	if err != nil {
		return nil, fmt.Errorf("read configuration: %w", err)
	}
	defer func() { _ = file.Close() }()

	return readConfigurationContents(io.LimitReader(file, maxConfigurationBytes+1))
}

func readConfigurationContext(ctx context.Context, path string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("read configuration: %w", err)
	}
	file, err := os.Open(path) // #nosec G304 -- explicit policy path
	if err != nil {
		return nil, fmt.Errorf("read configuration: %w", err)
	}
	defer func() { _ = file.Close() }()

	return readConfigurationContents(io.LimitReader(
		readerWithContext(ctx, file),
		maxConfigurationBytes+1,
	))
}

func readerWithContext(ctx context.Context, reader io.Reader) io.Reader {
	return contextReadFunc(func(buffer []byte) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
		count, err := reader.Read(buffer)
		select {
		case <-ctx.Done():
			return count, ctx.Err()
		default:
			return count, err
		}
	})
}

func readConfigurationContents(reader io.Reader) ([]byte, error) {
	contents, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read configuration: %w", err)
	}
	if len(contents) > maxConfigurationBytes {
		return nil, fmt.Errorf("configuration exceeds %d bytes", maxConfigurationBytes)
	}
	return contents, nil
}

type contextReadFunc func([]byte) (int, error)

func (read contextReadFunc) Read(buffer []byte) (int, error) { return read(buffer) }

func decodeConfig(contents []byte, root string, knownRules []string) (*Config, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(contents))
	decoder.KnownFields(true)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("decode configuration: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err != nil {
			return nil, fmt.Errorf("decode trailing configuration: %w", err)
		}
		return nil, errors.New("configuration must contain exactly one document")
	}

	config.Root = root
	if err := config.Validate(knownRules); err != nil {
		return nil, err
	}

	return &config, nil
}

// Validate checks shared policy structure against governed IDs without mutation.
// It does not authenticate promotion evidence, check file existence/generated
// headers, or validate every analyzer-specific typed policy. A nil knownRules
// list permits no explicit rule override or exception IDs.
func (config Config) Validate(knownRules []string) error {
	if config.Version != 1 {
		return fmt.Errorf("unsupported configuration version %d", config.Version)
	}
	known := make(map[string]struct{}, len(knownRules))
	for _, rule := range knownRules {
		known[rule] = struct{}{}
	}
	for rule, policy := range config.Rules {
		if _, ok := known[rule]; !ok {
			return fmt.Errorf("unknown rule %q", rule)
		}
		if policy.Status != "" && !validStatus(policy.Status) {
			return fmt.Errorf("rule %q has unknown status %q", rule, policy.Status)
		}
		if policy.Severity != "" && !validSeverity(policy.Severity) {
			return fmt.Errorf("rule %q has unknown severity %q", rule, policy.Severity)
		}
		if policy.Options != nil {
			return fmt.Errorf("rule %q options are unsupported; use its typed top-level policy", rule)
		}
		if policy.Status == StatusBlocking {
			if policy.Promotion == nil {
				return fmt.Errorf("blocking rule %q requires promotion evidence", rule)
			}
			if !versionPattern.MatchString(policy.Promotion.Version) {
				return fmt.Errorf("blocking rule %q requires a semantic promotion version", rule)
			}
			if strings.TrimSpace(policy.Promotion.Evidence) == "" {
				return fmt.Errorf("blocking rule %q requires non-empty promotion evidence", rule)
			}
		} else if policy.Promotion != nil {
			return fmt.Errorf("non-blocking rule %q must not declare promotion evidence", rule)
		}
	}
	if config.AdapterRoots != nil {
		return errors.New("adapter_roots are unsupported without a typed backend client policy")
	}
	if config.Generated.Exclude && len(config.Generated.Paths) == 0 {
		return errors.New("generated exclusion requires exact paths")
	}
	if !config.Generated.Exclude && config.Generated.Paths != nil {
		return errors.New("generated paths require exclusion to be enabled")
	}
	if len(config.Generated.Paths) > maxGeneratedPaths {
		return fmt.Errorf("generated exclusion exceeds %d paths", maxGeneratedPaths)
	}
	generatedPaths := make(map[string]struct{}, len(config.Generated.Paths))
	for index, generatedPath := range config.Generated.Paths {
		if !validGeneratedPath(generatedPath) {
			return fmt.Errorf(
				"generated path %d requires an exact repository-relative Go file",
				index,
			)
		}
		if _, exists := generatedPaths[generatedPath]; exists {
			return fmt.Errorf("generated path %d duplicates %q", index, generatedPath)
		}
		generatedPaths[generatedPath] = struct{}{}
	}
	for index, exception := range config.Exceptions {
		if _, ok := known[exception.Rule]; !ok {
			return fmt.Errorf("exception %d names unknown rule %q", index, exception.Rule)
		}
		if !validExceptionPackage(exception.Package) {
			return fmt.Errorf("exception %d requires an exact clean package", index)
		}
		if exception.Path != "" && !validExceptionPath(exception.Path) {
			return fmt.Errorf("exception %d requires an exact repository-relative path", index)
		}
		if strings.TrimSpace(exception.Reason) == "" {
			return fmt.Errorf("exception %d requires a reason", index)
		}
		if exception.Expires != "" {
			if _, err := time.Parse("2006-01-02", exception.Expires); err != nil {
				return fmt.Errorf("exception %d has invalid expiry: %w", index, err)
			}
		}
		if exception.Issue != "" && strings.TrimSpace(exception.Issue) == "" {
			return fmt.Errorf("exception %d has an empty issue", index)
		}
		for previous := range index {
			other := config.Exceptions[previous]
			if other.Rule == exception.Rule && other.Package == exception.Package &&
				(other.Path == exception.Path || other.Path == "" || exception.Path == "") {
				return fmt.Errorf("exception %d overlaps exception %d", index, previous)
			}
		}
	}

	return nil
}

func validGeneratedPath(value string) bool {
	return validExceptionPath(value) && strings.HasSuffix(value, ".go") &&
		!strings.ContainsAny(value, "*") && !strings.Contains(value, "...")
}
