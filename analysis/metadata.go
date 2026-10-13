// Package analysis owns strict policy configuration, rule metadata, attributable
// suppressions and deterministic JSON/SARIF 2.1.0 reporting. LoadConfigContext or
// ParseConfigBytes produces a policy snapshot; analyzer-specific constructors
// validate their own options. This package does not load or execute target code.
//
// Exported records are mutable and caller-owned. Do not mutate shared maps,
// slices or records during validation/reporting. Writers borrow their output and
// copy outer report slices, not nested rule metadata; no callback or IO lifecycle
// is retained. Omitting source/snippet fields is not arbitrary-text redaction:
// configured attribution, diagnostic messages and underlying errors remain data
// the caller must govern. Parser byte limits are not a process-memory budget.
package analysis

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Category groups rules by the engineering policy they enforce.
type Category string

const (
	// CategoryArchitecture groups import and package-boundary policy.
	CategoryArchitecture Category = "architecture"
	// CategoryContext groups context propagation and storage policy.
	CategoryContext Category = "context"
	// CategoryHTTP groups HTTP client ownership policy.
	CategoryHTTP Category = "http"
	// CategoryLifecycle groups construction, resource and goroutine lifetime policy.
	CategoryLifecycle Category = "lifecycle"
	// CategorySecurity groups unsafe boundaries and typed disclosure policy.
	CategorySecurity Category = "security"
	// CategoryAPI groups public API migration and ownership policy.
	CategoryAPI Category = "api"
	// CategorySafety groups shared-state policy.
	CategorySafety Category = "safety"
	// CategoryObservability groups metric schema and cardinality policy.
	CategoryObservability Category = "observability"
)

// Severity determines how prominently a diagnostic is reported.
type Severity string

const (
	// SeverityInfo maps to a SARIF note.
	SeverityInfo Severity = "info"
	// SeverityWarning maps to a SARIF warning.
	SeverityWarning Severity = "warning"
	// SeverityError maps to a SARIF error; it does not itself make a finding blocking.
	SeverityError Severity = "error"
)

// Status controls whether a rule informs developers or blocks a build.
type Status string

const (
	// StatusDisabled omits the rule from driver execution.
	StatusDisabled Status = "disabled"
	// StatusAdvisory retains findings without setting the driver blocking result.
	StatusAdvisory Status = "advisory"
	// StatusBlocking makes a remaining finding fail blocking policy.
	StatusBlocking Status = "blocking"
)

// ConfigurationType is a primitive accepted by a rule's configuration.
type ConfigurationType string

const (
	// ConfigurationString describes textual option metadata.
	ConfigurationString ConfigurationType = "string"
	// ConfigurationBoolean describes boolean option metadata.
	ConfigurationBoolean ConfigurationType = "boolean"
	// ConfigurationInteger describes integral option metadata.
	ConfigurationInteger ConfigurationType = "integer"
	// ConfigurationArray describes collection option metadata.
	ConfigurationArray ConfigurationType = "array"
	// ConfigurationObject describes keyed option metadata.
	ConfigurationObject ConfigurationType = "object"
)

// ConfigurationProperty documents one accepted rule configuration key.
type ConfigurationProperty struct {
	// Type describes the primitive shape in metadata; Rule.Validate does not enforce
	// it.
	Type ConfigurationType `json:"type"`
	// Description is optional explanatory text, not executable validation.
	Description string `json:"description,omitempty"`
	// Required documents requiredness; the generic metadata validator does not
	// enforce it.
	Required bool `json:"required,omitempty"`
}

// ConfigurationSchema documents the complete rule-specific configuration.
type ConfigurationSchema struct {
	// Properties maps option names to descriptive metadata; nil means no declared
	// properties.
	Properties map[string]ConfigurationProperty `json:"properties,omitempty"`
}

// Rule contains stable, machine-readable metadata for one policy rule.
type Rule struct {
	// ID is a stable category/kebab-case key used by policy and suppression
	// matching.
	ID string `json:"id"`
	// Category is one of the supported engineering-policy families.
	Category Category `json:"category"`
	// Severity sets report prominence, independently of blocking status.
	Severity Severity `json:"severity"`
	// DefaultStatus sets default disabled/advisory/blocking enforcement.
	DefaultStatus Status `json:"default_status"`
	// Rationale is required human-readable policy motivation, emitted in reports.
	Rationale string `json:"rationale"`
	// Remediation is required guidance, not an automatically applied fix.
	Remediation string `json:"remediation"`
	// IntroducedVersion is the required numeric major.minor.patch introduction
	// version.
	IntroducedVersion string `json:"introduced_version"`
	// Configuration describes options; Validate does not validate its properties.
	Configuration ConfigurationSchema `json:"configuration"`
}

var (
	ruleIDPattern  = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*/[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)
	versionPattern = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)
)

// Validate rejects invalid ID/category/severity/status, blank rationale or
// remediation, and a non-numeric semantic introduction version. It does not
// validate Configuration properties or authenticate metadata provenance.
func (rule Rule) Validate() error {
	if !ruleIDPattern.MatchString(rule.ID) {
		return errors.New("rule ID must use category/kebab-case syntax")
	}
	if !validCategory(rule.Category) {
		return fmt.Errorf("unknown category %q", rule.Category)
	}
	if !validSeverity(rule.Severity) {
		return fmt.Errorf("unknown severity %q", rule.Severity)
	}
	if !validStatus(rule.DefaultStatus) {
		return fmt.Errorf("unknown default status %q", rule.DefaultStatus)
	}
	if strings.TrimSpace(rule.Rationale) == "" {
		return errors.New("rationale is required")
	}
	if strings.TrimSpace(rule.Remediation) == "" {
		return errors.New("remediation is required")
	}
	if !versionPattern.MatchString(rule.IntroducedVersion) {
		return errors.New("introduced version must be a semantic version")
	}

	return nil
}

func validCategory(category Category) bool {
	switch category {
	case CategoryArchitecture, CategoryContext, CategoryHTTP, CategoryLifecycle,
		CategorySecurity, CategoryAPI, CategorySafety, CategoryObservability:
		return true
	default:
		return false
	}
}

func validSeverity(severity Severity) bool {
	switch severity {
	case SeverityInfo, SeverityWarning, SeverityError:
		return true
	default:
		return false
	}
}

func validStatus(status Status) bool {
	switch status {
	case StatusDisabled, StatusAdvisory, StatusBlocking:
		return true
	default:
		return false
	}
}
