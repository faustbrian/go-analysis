// Package noprocesscontrol restricts process termination to entrypoints.
//
// Use New or Analyzer for builtin panic and selected os.Exit/log.Fatal calls.
// Exact approved packages, .test packages and _test.go files are exempt.
// Function-value indirection and every possible process-termination API are not
// covered.
//
// Analyzer descriptors and Rule metadata are exported mutable records. Configure
// before execution and do not mutate descriptors or shared metadata concurrently.
// Checks use supplied syntax/types; package loading and policy enforcement belong
// to the caller or driver, not this package.
package noprocesscontrol

import (
	"go/ast"
	"go/types"
	"strings"

	shared "github.com/faustbrian/go-analysis/analysis"
	"golang.org/x/tools/go/analysis"
)

const ruleID = "lifecycle/no-process-control"

// Options configures approved composition-root packages.
type Options struct {
	// AllowedPackages lists exact importing/calling package exemptions; nil grants
	// none.
	AllowedPackages []string
}

// Rule is the stable metadata for process-control policy.
var Rule = shared.Rule{
	ID:                ruleID,
	Category:          shared.CategoryLifecycle,
	Severity:          shared.SeverityError,
	DefaultStatus:     shared.StatusAdvisory,
	Rationale:         "Library process termination bypasses caller-owned error handling and cleanup.",
	Remediation:       "Return a classified error and let an approved entrypoint terminate the process.",
	IntroducedVersion: "0.1.0",
	Configuration: shared.ConfigurationSchema{Properties: map[string]shared.ConfigurationProperty{
		"allowed_packages": {
			Type:        shared.ConfigurationArray,
			Description: "Exact package import paths allowed to control the process.",
		},
	}},
}

// Analyzer reports process-control calls outside approved entrypoints.
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
		Name: "noprocesscontrol",
		Doc:  Rule.Rationale,
		Run: func(pass *analysis.Pass) (any, error) {
			if _, ok := allowed[pass.Pkg.Path()]; ok ||
				strings.HasSuffix(pass.Pkg.Path(), ".test") {
				return nil, nil
			}
			for _, file := range pass.Files {
				if strings.HasSuffix(pass.Fset.Position(file.Pos()).Filename, "_test.go") {
					continue
				}
				ast.Inspect(file, func(node ast.Node) bool {
					call, ok := node.(*ast.CallExpr)
					if !ok {
						return true
					}
					reportCall(pass, call)
					return true
				})
			}
			return nil, nil
		},
	}
}

func reportCall(pass *analysis.Pass, call *ast.CallExpr) {
	switch function := call.Fun.(type) {
	case *ast.Ident:
		builtin, ok := pass.TypesInfo.Uses[function].(*types.Builtin)
		if ok && builtin.Name() == "panic" {
			pass.Reportf(function.Pos(), "%s: panic is restricted to approved entrypoints", ruleID)
		}
	case *ast.SelectorExpr:
		called, ok := pass.TypesInfo.Uses[function.Sel].(*types.Func)
		if !ok || called.Pkg() == nil {
			return
		}
		packagePath := called.Pkg().Path()
		if packagePath == "os" && called.Name() == "Exit" {
			pass.Reportf(function.Pos(), "%s: os.Exit is restricted to approved entrypoints", ruleID)
		}
		if packagePath == "log" && isFatal(called.Name()) {
			pass.Reportf(function.Pos(), "%s: log.%s is restricted to approved entrypoints", ruleID, called.Name())
		}
	}
}

func isFatal(name string) bool {
	return name == "Fatal" || name == "Fatalf" || name == "Fatalln"
}
