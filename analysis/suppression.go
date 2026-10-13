package analysis

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"time"
)

const suppressionPrefix = "//analysis:ignore "

// Suppression attributes an inline exception to the line after its comment group.
// Exported records are mutable; validation/use tracking does not authenticate
// the reason, issue or source author.
type Suppression struct {
	// Rule is the exact known rule ID selected by the directive.
	Rule string `json:"rule"`
	// Filename is the file-set source filename, matched exactly before report
	// normalization.
	Filename string `json:"filename"`
	// DirectiveLine is the one-based line containing the directive.
	DirectiveLine int `json:"directive_line"`
	// TargetLine is the line after the entire containing comment group, not
	// necessarily after this directive.
	TargetLine int `json:"target_line"`
	// Reason is required attribution text and is emitted in reports.
	Reason string `json:"reason"`
	// Expires is an optional UTC expiry date, valid through that entire day; zero
	// means no expiry.
	Expires time.Time `json:"expires,omitempty"`
	// Issue is optional attribution text; no issue lookup is performed.
	Issue string `json:"issue,omitempty"`
	// Used marks consumption; caller-supplied true is retained by ApplySuppressions.
	Used bool `json:"used"`
}

// ParseSuppressions reads exact //analysis:ignore directives from parsed comments.
// Rules must be known; reason is required, metadata keys may not repeat, and
// expires uses an inclusive UTC calendar date. TargetLine follows group.End,
// including other comments in that group. Malformed or duplicate target/rule
// directives return nil inventory and an error. file/fset must be non-nil and
// consistent. The parser borrows them without mutation and has no size limit.
func ParseSuppressions(
	fset *token.FileSet,
	file *ast.File,
	knownRules []string,
	now time.Time,
) ([]Suppression, error) {
	known := make(map[string]struct{}, len(knownRules))
	for _, rule := range knownRules {
		known[rule] = struct{}{}
	}
	var suppressions []Suppression
	seen := make(map[string]struct{})
	for _, group := range file.Comments {
		for _, comment := range group.List {
			if !strings.HasPrefix(comment.Text, suppressionPrefix) {
				continue
			}
			suppression, err := parseSuppression(
				comment.Text,
				fset.Position(comment.Slash),
				fset.Position(group.End()).Line+1,
				known,
				now,
			)
			if err != nil {
				return nil, err
			}
			key := fmt.Sprintf("%s:%d:%s", suppression.Filename,
				suppression.TargetLine, suppression.Rule)
			if _, exists := seen[key]; exists {
				return nil, fmt.Errorf("duplicate suppression for %s", key)
			}
			seen[key] = struct{}{}
			suppressions = append(suppressions, suppression)
		}
	}

	return suppressions, nil
}

func parseSuppression(
	directive string,
	position token.Position,
	targetLine int,
	known map[string]struct{},
	now time.Time,
) (Suppression, error) {
	identity, details, ok := strings.Cut(
		strings.TrimPrefix(directive, suppressionPrefix),
		" -- ",
	)
	if !ok || strings.TrimSpace(identity) == "" {
		return Suppression{}, errors.New("suppression requires a rule and reason")
	}
	if _, exists := known[identity]; !exists {
		return Suppression{}, fmt.Errorf("suppression names unknown rule %q", identity)
	}
	parts := strings.Split(details, ";")
	reason := strings.TrimSpace(parts[0])
	if reason == "" {
		return Suppression{}, errors.New("suppression reason must not be empty")
	}
	suppression := Suppression{
		Rule:          identity,
		Filename:      position.Filename,
		DirectiveLine: position.Line,
		TargetLine:    targetLine,
		Reason:        reason,
	}
	metadata := make(map[string]struct{})
	for _, raw := range parts[1:] {
		key, value, found := strings.Cut(strings.TrimSpace(raw), "=")
		if !found || strings.TrimSpace(value) == "" {
			return Suppression{}, errors.New("suppression metadata must use key=value")
		}
		if _, exists := metadata[key]; exists {
			return Suppression{}, fmt.Errorf("duplicate suppression metadata %q", key)
		}
		metadata[key] = struct{}{}
		switch key {
		case "expires":
			expires, err := time.Parse("2006-01-02", value)
			if err != nil {
				return Suppression{}, fmt.Errorf("parse suppression expiry: %w", err)
			}
			if expires.Before(midnightUTC(now)) {
				return Suppression{}, fmt.Errorf("suppression expired on %s", value)
			}
			suppression.Expires = expires
		case "issue":
			suppression.Issue = value
		default:
			return Suppression{}, fmt.Errorf("unknown suppression metadata %q", key)
		}
	}

	return suppression, nil
}

func midnightUTC(value time.Time) time.Time {
	year, month, day := value.UTC().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// ApplySuppressions filters exact rule/filename/target-line matches and updates
// a copied inventory. It does not revalidate directive syntax, expiry or prior
// Used flags; normally pass a freshly parsed inventory. Duplicate inputs use the
// first match. An unused entry returns nil findings, the updated inventory and
// an error. Inputs are not mutated; retained findings use a new outer slice.
func ApplySuppressions(
	diagnostics []Diagnostic,
	suppressions []Suppression,
) ([]Diagnostic, []Suppression, error) {
	inventory := append([]Suppression(nil), suppressions...)
	remaining := make([]Diagnostic, 0, len(diagnostics))
	for _, diagnostic := range diagnostics {
		matched := false
		for index := range inventory {
			suppression := &inventory[index]
			if suppression.Rule == diagnostic.Rule &&
				suppression.Filename == diagnostic.Filename &&
				suppression.TargetLine == diagnostic.Line {
				suppression.Used = true
				matched = true
				break
			}
		}
		if !matched {
			remaining = append(remaining, diagnostic)
		}
	}
	for _, suppression := range inventory {
		if !suppression.Used {
			return nil, inventory, fmt.Errorf(
				"stale suppression for %s at %s:%d",
				suppression.Rule,
				suppression.Filename,
				suppression.DirectiveLine,
			)
		}
	}

	return remaining, inventory, nil
}
