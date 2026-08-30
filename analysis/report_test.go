package analysis_test

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	shared "github.com/faustbrian/go-analysis/analysis"
)

func TestWriteJSONIsDeterministicWithoutMutatingInput(t *testing.T) {
	t.Parallel()

	report := shared.Report{
		ToolVersion: "0.1.0",
		Rules: []shared.Rule{
			{ID: "security/no-unsafe"},
			{ID: "context/no-background"},
		},
		Diagnostics: []shared.Diagnostic{
			{Rule: "security/no-unsafe", Filename: "z.go", Line: 9, Message: "unsafe"},
			{Rule: "context/no-background", Filename: "a.go", Line: 2, Message: "context"},
		},
		Exceptions: []shared.PolicyException{
			{Rule: "security/no-unsafe", Package: "z/package", Path: "z.go"},
			{Rule: "context/no-background", Package: "a/package", Path: "a.go"},
		},
		Suppressions: []shared.Suppression{
			{Rule: "security/no-unsafe", Filename: "z.go", DirectiveLine: 8},
			{Rule: "context/no-background", Filename: "a.go", DirectiveLine: 1},
		},
	}
	originalFirst := report.Diagnostics[0]
	originalException := report.Exceptions[0]
	var first bytes.Buffer
	if err := shared.WriteJSON(&first, report); err != nil {
		t.Fatalf("WriteJSON() error = %v", err)
	}
	report.Diagnostics[0], report.Diagnostics[1] = report.Diagnostics[1], report.Diagnostics[0]
	report.Exceptions[0], report.Exceptions[1] = report.Exceptions[1], report.Exceptions[0]
	var second bytes.Buffer
	if err := shared.WriteJSON(&second, report); err != nil {
		t.Fatalf("WriteJSON() second error = %v", err)
	}
	if first.String() != second.String() {
		t.Fatalf("JSON output is order-dependent:\n%s\n%s", first.String(), second.String())
	}
	var normalized shared.Report
	if err := json.Unmarshal(first.Bytes(), &normalized); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if normalized.Rules[0].ID != "context/no-background" ||
		normalized.Diagnostics[0].Filename != "a.go" ||
		normalized.Exceptions[0].Rule != "context/no-background" ||
		normalized.Suppressions[0].Filename != "a.go" {
		t.Fatalf("normalized order = %#v", normalized)
	}
	if originalFirst != report.Diagnostics[1] {
		t.Fatal("WriteJSON() mutated caller diagnostics")
	}
	if originalException != report.Exceptions[1] {
		t.Fatal("WriteJSON() mutated caller exceptions")
	}
}

func TestWriteJSONRelativizesRepositoryPaths(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	report := shared.Report{
		Root: root,
		Diagnostics: []shared.Diagnostic{{
			Filename: filepath.Join(root, "internal", "service.go"),
		}},
		Suppressions: []shared.Suppression{
			{Filename: filepath.Join(root, "z.go"), DirectiveLine: 2},
			{Filename: filepath.Join(root, "a.go"), DirectiveLine: 1},
		},
	}
	var output bytes.Buffer
	if err := shared.WriteJSON(&output, report); err != nil {
		t.Fatalf("WriteJSON() error = %v", err)
	}
	var got shared.Report
	if err := json.Unmarshal(output.Bytes(), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got.Diagnostics[0].Filename != "internal/service.go" ||
		got.Suppressions[0].Filename != "a.go" {
		t.Fatalf("normalized report = %#v", got)
	}
}

func TestWritersEncodeEmptyInventoriesAsArrays(t *testing.T) {
	t.Parallel()

	var jsonOutput bytes.Buffer
	if err := shared.WriteJSON(&jsonOutput, shared.Report{}); err != nil {
		t.Fatalf("WriteJSON() error = %v", err)
	}
	for _, field := range []string{"rules", "diagnostics", "exceptions", "suppressions"} {
		if !strings.Contains(jsonOutput.String(), `"`+field+`":[]`) {
			t.Errorf("JSON %s is not an empty array: %s", field, jsonOutput.String())
		}
	}

	var sarifOutput bytes.Buffer
	if err := shared.WriteSARIF(&sarifOutput, shared.Report{}); err != nil {
		t.Fatalf("WriteSARIF() error = %v", err)
	}
	for _, field := range []string{"rules", "results", "exceptions", "suppressions"} {
		if !strings.Contains(sarifOutput.String(), `"`+field+`":[]`) {
			t.Errorf("SARIF %s is not an empty array: %s", field, sarifOutput.String())
		}
	}
}

func TestWritersRejectPathsThatCouldExposeWorkspace(t *testing.T) {
	t.Parallel()

	absolute := filepath.Join(t.TempDir(), "secret.go")
	reports := []shared.Report{
		{Diagnostics: []shared.Diagnostic{{Filename: absolute}}},
		{Suppressions: []shared.Suppression{{Filename: absolute}}},
		{Diagnostics: []shared.Diagnostic{{Filename: "../secret.go"}}},
		{Suppressions: []shared.Suppression{{Filename: "../secret.go"}}},
		{
			Root:        t.TempDir(),
			Diagnostics: []shared.Diagnostic{{Filename: absolute}},
		},
	}
	for _, report := range reports {
		var output bytes.Buffer
		if err := shared.WriteJSON(&output, report); err == nil {
			t.Fatal("WriteJSON() accepted unsafe absolute path")
		}
		if err := shared.WriteSARIF(&output, report); err == nil {
			t.Fatal("WriteSARIF() accepted unsafe absolute path")
		}
	}
}

func TestWriteSARIFIncludesStableRulesAndNoSource(t *testing.T) {
	t.Parallel()

	report := shared.Report{
		ToolVersion: "0.1.0",
		Rules: []shared.Rule{
			{
				ID:                "security/no-unsafe",
				Category:          shared.CategorySecurity,
				Severity:          shared.SeverityError,
				DefaultStatus:     shared.StatusAdvisory,
				Rationale:         "Unsafe bypasses language guarantees.",
				Remediation:       "Use a safe API.",
				IntroducedVersion: "0.1.0",
			},
			{ID: "context/no-background", Severity: shared.SeverityWarning},
			{ID: "api/broad-interface", Severity: shared.SeverityInfo},
			{ID: "unknown/severity", Severity: shared.Severity("unexpected")},
		},
		Diagnostics: []shared.Diagnostic{
			{Rule: "security/no-unsafe", Filename: "internal/unsafe.go", Line: 4, Column: 2, Message: "unsafe import"},
			{Rule: "context/no-background", Filename: "context.go", Line: 2, Message: "context"},
			{Rule: "api/broad-interface", Filename: "api.go", Line: 3, Message: "interface"},
			{Rule: "unknown/rule", Filename: "unknown.go", Line: 1, Message: "unknown"},
			{Rule: "unknown/severity", Filename: "severity.go", Line: 1, Message: "severity"},
		},
		Exceptions: []shared.PolicyException{{
			Rule: "security/no-unsafe", Package: "example.com/service", Used: true,
		}},
		Suppressions: []shared.Suppression{{
			Rule: "context/no-background", Filename: "context.go", Used: true,
		}},
	}
	var output bytes.Buffer
	if err := shared.WriteSARIF(&output, report); err != nil {
		t.Fatalf("WriteSARIF() error = %v", err)
	}
	if strings.Contains(output.String(), "source") || strings.Contains(output.String(), "snippet") {
		t.Fatalf("SARIF unexpectedly embeds source: %s", output.String())
	}
	var document map[string]any
	if err := json.Unmarshal(output.Bytes(), &document); err != nil {
		t.Fatalf("SARIF JSON error = %v", err)
	}
	if document["version"] != "2.1.0" {
		t.Fatalf("SARIF version = %#v", document["version"])
	}
	if document["$schema"] != "https://json.schemastore.org/sarif-2.1.0.json" {
		t.Fatalf("SARIF schema = %#v", document["$schema"])
	}
	runs := document["runs"].([]any)
	if len(runs) != 1 {
		t.Fatalf("SARIF runs = %#v", runs)
	}
	run := runs[0].(map[string]any)
	driver := run["tool"].(map[string]any)["driver"].(map[string]any)
	if driver["name"] != "analysis" || driver["version"] != "0.1.0" ||
		driver["informationUri"] != "https://github.com/faustbrian/go-analysis" {
		t.Fatalf("SARIF driver = %#v", driver)
	}
	descriptors := driver["rules"].([]any)
	if len(descriptors) != 4 {
		t.Fatalf("SARIF rule descriptors = %#v", descriptors)
	}
	wantRuleIDs := []string{
		"api/broad-interface",
		"context/no-background",
		"security/no-unsafe",
		"unknown/severity",
	}
	for index, raw := range descriptors {
		if got := raw.(map[string]any)["id"]; got != wantRuleIDs[index] {
			t.Fatalf("SARIF rule descriptor %d ID = %#v, want %q", index, got, wantRuleIDs[index])
		}
	}
	securityDescriptor := descriptors[2].(map[string]any)
	securityProperties := securityDescriptor["properties"].(map[string]any)
	if securityDescriptor["shortDescription"].(map[string]any)["text"] != "Unsafe bypasses language guarantees." ||
		securityDescriptor["help"].(map[string]any)["text"] != "Use a safe API." ||
		securityProperties["category"] != "security" ||
		securityProperties["severity"] != "error" ||
		securityProperties["defaultStatus"] != "advisory" ||
		securityProperties["introducedVersion"] != "0.1.0" {
		t.Fatalf("SARIF security descriptor = %#v", securityDescriptor)
	}
	properties := run["properties"].(map[string]any)
	if len(properties["exceptions"].([]any)) != 1 ||
		len(properties["suppressions"].([]any)) != 1 {
		t.Fatalf("SARIF inventory properties = %#v", properties)
	}
	results := runs[0].(map[string]any)["results"].([]any)
	wantLevels := map[string]string{
		"security/no-unsafe":    "error",
		"context/no-background": "warning",
		"api/broad-interface":   "note",
		"unknown/rule":          "warning",
		"unknown/severity":      "note",
	}
	for _, raw := range results {
		result := raw.(map[string]any)
		ruleID := result["ruleId"].(string)
		want, exists := wantLevels[ruleID]
		if !exists {
			t.Fatalf("unexpected SARIF rule = %q", ruleID)
		}
		if result["level"] != want {
			t.Errorf("SARIF level for %s = %#v, want %q", ruleID, result["level"], want)
		}
		delete(wantLevels, ruleID)
	}
	if len(wantLevels) != 0 {
		t.Fatalf("SARIF results omitted rules = %#v", wantLevels)
	}
	first := results[0].(map[string]any)
	locations := first["locations"].([]any)
	physical := locations[0].(map[string]any)["physicalLocation"].(map[string]any)
	artifact := physical["artifactLocation"].(map[string]any)
	region := physical["region"].(map[string]any)
	if artifact["uri"] != "api.go" || region["startLine"] != float64(3) {
		t.Fatalf("first SARIF location = %#v", physical)
	}
	if len(locations) != 1 || region["startColumn"] != nil {
		t.Fatalf("first SARIF location shape = %#v", physical)
	}
}

func TestWritersPropagateOutputFailure(t *testing.T) {
	t.Parallel()

	writer := failingWriter{}
	if err := shared.WriteJSON(writer, shared.Report{}); err == nil {
		t.Fatal("WriteJSON() error = nil")
	}
	if err := shared.WriteSARIF(writer, shared.Report{}); err == nil {
		t.Fatal("WriteSARIF() error = nil")
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, bytes.ErrTooLarge
}
