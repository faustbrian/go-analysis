// Package nobackground restricts root context creation to composition roots.
//
// Use New or Analyzer for resolved context.Background, TODO and WithoutCancel
// calls. Exact approved packages and _test.go files are exempt. No context is
// created and no body-level propagation is proved.
//
// Analyzer descriptors and Rule metadata are exported mutable records. Configure
// before execution and do not mutate descriptors or shared metadata concurrently.
// Checks use supplied syntax/types; package loading and policy enforcement belong
// to the caller or driver, not this package.
package nobackground

import (
	"go/ast"
	"go/types"
	"strings"

	shared "github.com/faustbrian/go-analysis/analysis"
	"golang.org/x/tools/go/analysis"
)

const ruleID = "context/no-background"

// Options configures packages that own root context creation.
type Options struct {
	// AllowedPackages lists exact importing/calling package exemptions; nil grants
	// none.
	AllowedPackages []string
}

// Rule is the stable metadata for root context policy.
var Rule = shared.Rule{
	ID:                ruleID,
	Category:          shared.CategoryContext,
	Severity:          shared.SeverityError,
	DefaultStatus:     shared.StatusAdvisory,
	Rationale:         "Replacing an available lifecycle with a root context loses cancellation and deadlines.",
	Remediation:       "Propagate the caller context or create roots only in approved composition packages.",
	IntroducedVersion: "0.1.0",
	Configuration: shared.ConfigurationSchema{Properties: map[string]shared.ConfigurationProperty{
		"allowed_packages": {
			Type:        shared.ConfigurationArray,
			Description: "Exact package paths allowed to create root contexts.",
		},
	}},
}

// Analyzer reports root and detached contexts outside approved roots.
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
		Name: "nobackground",
		Doc:  Rule.Rationale,
		Run: func(pass *analysis.Pass) (any, error) {
			if _, ok := allowed[pass.Pkg.Path()]; ok {
				return nil, nil
			}
			for _, file := range pass.Files {
				if strings.HasSuffix(pass.Fset.Position(file.Pos()).Filename, "_test.go") {
					continue
				}
				ast.Inspect(file, func(node ast.Node) bool {
					call, ok := node.(*ast.CallExpr)
					if ok {
						reportCall(pass, call)
					}
					return true
				})
			}
			return nil, nil
		},
	}
}

func reportCall(pass *analysis.Pass, call *ast.CallExpr) {
	var object types.Object
	switch function := call.Fun.(type) {
	case *ast.Ident:
		object = pass.TypesInfo.Uses[function]
	case *ast.SelectorExpr:
		object = pass.TypesInfo.Uses[function.Sel]
	default:
		return
	}
	called, ok := object.(*types.Func)
	if !ok || called.Pkg() == nil || called.Pkg().Path() != "context" {
		return
	}
	if called.Name() != "Background" && called.Name() != "TODO" &&
		called.Name() != "WithoutCancel" {
		return
	}
	pass.Reportf(
		call.Fun.Pos(),
		"%s: context.%s is restricted to approved roots",
		ruleID,
		called.Name(),
	)
}
