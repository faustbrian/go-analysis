// Package noinit requires explicit construction instead of package init.
//
// Use New or Analyzer to report init declarations outside exact approved
// packages. Package variable initialization remains a separate concern; the
// analyzer does not execute initialization.
//
// Analyzer descriptors and Rule metadata are exported mutable records. Configure
// before execution and do not mutate descriptors or shared metadata concurrently.
// Checks use supplied syntax/types; package loading and policy enforcement belong
// to the caller or driver, not this package.
package noinit

import (
	"go/ast"

	shared "github.com/faustbrian/go-analysis/analysis"
	"golang.org/x/tools/go/analysis"
)

const ruleID = "lifecycle/no-init"

// Options configures packages with a reviewed need for init.
type Options struct {
	// AllowedPackages lists exact importing/calling package exemptions; nil grants
	// none.
	AllowedPackages []string
}

// Rule is the stable metadata for explicit construction policy.
var Rule = shared.Rule{
	ID:                ruleID,
	Category:          shared.CategoryLifecycle,
	Severity:          shared.SeverityError,
	DefaultStatus:     shared.StatusAdvisory,
	Rationale:         "Package init hides ordering, failure, and lifecycle ownership.",
	Remediation:       "Expose an explicit constructor called by the composition root.",
	IntroducedVersion: "0.1.0",
	Configuration: shared.ConfigurationSchema{Properties: map[string]shared.ConfigurationProperty{
		"allowed_packages": {
			Type:        shared.ConfigurationArray,
			Description: "Exact package paths with reviewed init requirements.",
		},
	}},
}

// Analyzer reports package init declarations outside configured exceptions.
var Analyzer = New(Options{})

// New constructs an independently configured analyzer.
// Configuration collections are copied into the analyzer; pass syntax/type data
// is borrowed for each run. No caller collection is mutated.
func New(options Options) *analysis.Analyzer {
	allowed := make(map[string]struct{}, len(options.AllowedPackages))
	for _, packagePath := range options.AllowedPackages {
		allowed[packagePath] = struct{}{}
	}

	return &analysis.Analyzer{
		Name: "noinit",
		Doc:  Rule.Rationale,
		Run: func(pass *analysis.Pass) (any, error) {
			if _, ok := allowed[pass.Pkg.Path()]; ok {
				return nil, nil
			}
			for _, file := range pass.Files {
				for _, declaration := range file.Decls {
					function, ok := declaration.(*ast.FuncDecl)
					if !ok || function.Name.Name != "init" {
						continue
					}
					pass.Reportf(
						function.Name.Pos(),
						"%s: package init hides construction and lifecycle ownership",
						ruleID,
					)
				}
			}
			return nil, nil
		},
	}
}
