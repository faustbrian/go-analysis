package analysis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
)

// Report is the stable machine-readable result of one analysis run.
type Report struct {
	// ToolVersion is caller-supplied provenance text, not inferred from a running
	// binary.
	ToolVersion string `json:"tool_version"`
	// Root anchors absolute diagnostic/suppression paths and is not serialized.
	Root string `json:"-"`
	// Rules is the rule inventory; writers sort an outer copy by ID without
	// validating metadata.
	Rules []Rule `json:"rules"`
	// Diagnostics contains retained findings; writers sort and normalize copied
	// records.
	Diagnostics []Diagnostic `json:"diagnostics"`
	// Exceptions records policy exception attribution and use; writers do not
	// normalize its Path field.
	Exceptions []PolicyException `json:"exceptions"`
	// Suppressions records inline directive attribution, normalized filenames and
	// use.
	Suppressions []Suppression `json:"suppressions"`
}

// WriteJSON emits a newline-terminated report with sorted outer inventories and
// normalized diagnostic/suppression paths. It omits source/snippet fields, but
// caller messages, attribution, metadata and exception Path remain unsanitized.
// Validation precedes writing; writer errors may leave partial output. The writer
// is borrowed, must be non-nil, and is never flushed or closed. Nested rule maps
// remain shared with the input; do not mutate them during serialization. No byte
// or entry limit is enforced here.
func WriteJSON(writer io.Writer, report Report) error {
	normalized, err := normalizeReport(report)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(true)
	if err := encoder.Encode(normalized); err != nil {
		return fmt.Errorf("encode JSON report: %w", err)
	}

	return nil
}

// WriteSARIF emits newline-terminated SARIF 2.1.0 with sorted inventories,
// relative diagnostic/suppression paths and no source/snippet fields. It shares
// WriteJSON's borrowed-writer, partial-output, mutation and arbitrary-text limits.
// Unknown diagnostic rule IDs receive warning level rather than being rejected;
// metadata and SARIF schema conformance are not validated by this writer.
func WriteSARIF(writer io.Writer, report Report) error {
	normalized, err := normalizeReport(report)
	if err != nil {
		return err
	}
	document := buildSARIF(normalized)
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(true)
	if err := encoder.Encode(document); err != nil {
		return fmt.Errorf("encode SARIF report: %w", err)
	}

	return nil
}

func normalizeReport(report Report) (Report, error) {
	normalized := report
	normalized.Rules = append([]Rule{}, report.Rules...)
	normalized.Diagnostics = append([]Diagnostic{}, report.Diagnostics...)
	normalized.Exceptions = append([]PolicyException{}, report.Exceptions...)
	normalized.Suppressions = append([]Suppression{}, report.Suppressions...)
	for index := range normalized.Diagnostics {
		filename, err := reportPath(report.Root, normalized.Diagnostics[index].Filename)
		if err != nil {
			return Report{}, err
		}
		normalized.Diagnostics[index].Filename = filename
	}
	for index := range normalized.Suppressions {
		filename, err := reportPath(report.Root, normalized.Suppressions[index].Filename)
		if err != nil {
			return Report{}, err
		}
		normalized.Suppressions[index].Filename = filename
	}
	sort.Slice(normalized.Rules, func(left, right int) bool {
		return ruleLess(normalized.Rules[left], normalized.Rules[right])
	})
	sort.Slice(normalized.Diagnostics, func(left, right int) bool {
		return diagnosticLess(
			normalized.Diagnostics[left],
			normalized.Diagnostics[right],
		)
	})
	sort.Slice(normalized.Suppressions, func(left, right int) bool {
		return suppressionLess(
			normalized.Suppressions[left],
			normalized.Suppressions[right],
		)
	})
	sort.Slice(normalized.Exceptions, func(left, right int) bool {
		return exceptionLess(normalized.Exceptions[left], normalized.Exceptions[right])
	})

	return normalized, nil
}

func ruleLess(left, right Rule) bool {
	return left.ID < right.ID
}

func diagnosticLess(left, right Diagnostic) bool {
	return diagnosticKey(left) < diagnosticKey(right)
}

func suppressionLess(left, right Suppression) bool {
	return suppressionKey(left) < suppressionKey(right)
}

func exceptionLess(left, right PolicyException) bool {
	return exceptionKey(left) < exceptionKey(right)
}

func reportPath(root, filename string) (string, error) {
	return reportPathWithRel(root, filename, filepath.Rel)
}

func reportPathWithRel(
	root string,
	filename string,
	rel func(string, string) (string, error),
) (string, error) {
	if !filepath.IsAbs(filename) {
		cleaned := filepath.Clean(filename)
		if cleaned == ".." ||
			strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
			return "", errors.New("diagnostic path escapes the report root")
		}
		return filepath.ToSlash(cleaned), nil
	}
	if root == "" {
		return "", errors.New("absolute diagnostic path requires a report root")
	}
	relative, err := rel(root, filename)
	if err != nil {
		return "", fmt.Errorf("make report path relative: %w", err)
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", errors.New("diagnostic path escapes the report root")
	}

	return filepath.ToSlash(relative), nil
}

func diagnosticKey(diagnostic Diagnostic) string {
	return fmt.Sprintf("%s\x00%012d\x00%012d\x00%s\x00%s",
		diagnostic.Filename,
		diagnostic.Line,
		diagnostic.Column,
		diagnostic.Rule,
		diagnostic.Message,
	)
}

func suppressionKey(suppression Suppression) string {
	return fmt.Sprintf("%s\x00%012d\x00%s",
		suppression.Filename,
		suppression.DirectiveLine,
		suppression.Rule,
	)
}

func exceptionKey(exception PolicyException) string {
	return fmt.Sprintf("%s\x00%s\x00%s", exception.Rule, exception.Package, exception.Path)
}

type sarifDocument struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool       sarifTool          `json:"tool"`
	Results    []sarifResult      `json:"results"`
	Properties sarifRunProperties `json:"properties"`
}

type sarifRunProperties struct {
	Exceptions   []PolicyException `json:"exceptions"`
	Suppressions []Suppression     `json:"suppressions"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string                `json:"name"`
	Version        string                `json:"version"`
	InformationURI string                `json:"informationUri"`
	Rules          []sarifRuleDescriptor `json:"rules"`
}

type sarifRuleDescriptor struct {
	ID               string         `json:"id"`
	ShortDescription sarifMessage   `json:"shortDescription"`
	Help             sarifMessage   `json:"help"`
	Properties       map[string]any `json:"properties"`
}

type sarifResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   sarifMessage    `json:"message"`
	Locations []sarifLocation `json:"locations"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           sarifRegion           `json:"region"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn,omitempty"`
}

func buildSARIF(report Report) sarifDocument {
	descriptors := make([]sarifRuleDescriptor, 0, len(report.Rules))
	levels := make(map[string]string, len(report.Rules))
	for _, rule := range report.Rules {
		levels[rule.ID] = sarifLevel(rule.Severity)
		descriptors = append(descriptors, sarifRuleDescriptor{
			ID:               rule.ID,
			ShortDescription: sarifMessage{Text: rule.Rationale},
			Help:             sarifMessage{Text: rule.Remediation},
			Properties: map[string]any{
				"category":          rule.Category,
				"defaultStatus":     rule.DefaultStatus,
				"introducedVersion": rule.IntroducedVersion,
				"severity":          rule.Severity,
			},
		})
	}
	results := make([]sarifResult, 0, len(report.Diagnostics))
	for _, diagnostic := range report.Diagnostics {
		level := levels[diagnostic.Rule]
		if level == "" {
			level = "warning"
		}
		results = append(results, sarifResult{
			RuleID:  diagnostic.Rule,
			Level:   level,
			Message: sarifMessage{Text: diagnostic.Message},
			Locations: []sarifLocation{{PhysicalLocation: sarifPhysicalLocation{
				ArtifactLocation: sarifArtifactLocation{URI: diagnostic.Filename},
				Region: sarifRegion{
					StartLine:   diagnostic.Line,
					StartColumn: diagnostic.Column,
				},
			}}},
		})
	}

	return sarifDocument{
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Version: "2.1.0",
		Runs: []sarifRun{{
			Tool: sarifTool{Driver: sarifDriver{
				Name:           "analysis",
				Version:        report.ToolVersion,
				InformationURI: "https://github.com/faustbrian/go-analysis",
				Rules:          descriptors,
			}},
			Results: results,
			Properties: sarifRunProperties{
				Exceptions:   report.Exceptions,
				Suppressions: report.Suppressions,
			},
		}},
	}
}

func sarifLevel(severity Severity) string {
	switch severity {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	default:
		return "note"
	}
}
