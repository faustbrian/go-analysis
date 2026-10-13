package driver_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	shared "github.com/faustbrian/go-analysis/analysis"
	"github.com/faustbrian/go-analysis/analyzers/nounsafe"
	"github.com/faustbrian/go-analysis/internal/driver"
	"github.com/faustbrian/go-analysis/policy"
)

func TestRunPropagatesConfiguredSeverityWithoutChangingBlockingPolicy(t *testing.T) {
	registry, err := policy.Builtin()
	if err != nil {
		t.Fatal(err)
	}
	baseline := nounsafe.Rule
	for _, test := range []struct {
		name     string
		severity shared.Severity
		status   shared.Status
		level    string
	}{
		{"default", "", "", "error"},
		{"blocking info", shared.SeverityInfo, shared.StatusBlocking, "note"},
		{"blocking warning", shared.SeverityWarning, shared.StatusBlocking, "warning"},
		{"advisory error", shared.SeverityError, shared.StatusAdvisory, "error"},
		{"advisory info", shared.SeverityInfo, shared.StatusAdvisory, "note"},
		{"disabled warning", shared.SeverityWarning, shared.StatusDisabled, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := writeModule(t, map[string]string{"sample/sample.go": "package sample\nimport _ \"unsafe\"\n"})
			configuration := "version: 1\n"
			if test.severity != "" || test.status != "" {
				configuration += "rules:\n  security/no-unsafe:\n"
				if test.status != "" {
					configuration += fmt.Sprintf("    status: %s\n", test.status)
				}
				if test.severity != "" {
					configuration += fmt.Sprintf("    severity: %s\n", test.severity)
				}
			}
			config := writePolicy(t, root, configuration)
			result, err := driver.Run(context.Background(), driver.Options{ConfigPath: config, Patterns: []string{"./..."}, Sequential: true})
			if err != nil {
				t.Fatal(err)
			}
			severity := test.severity
			if severity == "" {
				severity = baseline.Severity
			}
			status := test.status
			if status == "" {
				status = baseline.DefaultStatus
			}
			wantDiagnostics := 1
			if status == shared.StatusDisabled {
				wantDiagnostics = 0
			}
			if len(result.Report.Diagnostics) != wantDiagnostics || result.Blocking != (status == shared.StatusBlocking) {
				t.Errorf("diagnostics=%d blocking=%t; status=%s", len(result.Report.Diagnostics), result.Blocking, status)
			}
			if len(result.Report.Rules) != len(registry.Entries()) {
				t.Errorf("rule inventory=%d; want all %d rules", len(result.Report.Rules), len(registry.Entries()))
			}
			assertRule := func(rules []shared.Rule) {
				t.Helper()
				for _, rule := range rules {
					if rule.ID == nounsafe.Rule.ID {
						if rule.Severity != severity || rule.DefaultStatus != baseline.DefaultStatus {
							t.Errorf("reported severity=%s defaultStatus=%s; want %s, %s", rule.Severity, rule.DefaultStatus, severity, baseline.DefaultStatus)
						}
						return
					}
				}
				t.Error("governed no-unsafe rule missing from report")
			}
			assertRule(result.Report.Rules)
			var output bytes.Buffer
			if err := shared.WriteJSON(&output, result.Report); err != nil {
				t.Fatal(err)
			}
			var decoded shared.Report
			if err := json.Unmarshal(output.Bytes(), &decoded); err != nil {
				t.Fatal(err)
			}
			assertRule(decoded.Rules)
			output.Reset()
			if err := shared.WriteSARIF(&output, result.Report); err != nil {
				t.Fatal(err)
			}
			var sarif struct {
				Runs []struct {
					Tool struct {
						Driver struct {
							Rules []struct {
								ID         string
								Properties struct {
									Severity      shared.Severity
									DefaultStatus shared.Status
								}
							}
						}
					}
					Results []struct {
						RuleID string
						Level  string
					}
				}
			}
			if err := json.Unmarshal(output.Bytes(), &sarif); err != nil {
				t.Fatal(err)
			}
			if len(sarif.Runs) != 1 {
				t.Fatalf("SARIF runs=%d", len(sarif.Runs))
			}
			run := sarif.Runs[0]
			if len(run.Tool.Driver.Rules) != len(result.Report.Rules) {
				t.Error("SARIF dropped governed descriptors")
			}
			found := false
			previous := ""
			for _, descriptor := range run.Tool.Driver.Rules {
				if descriptor.ID <= previous {
					t.Error("SARIF descriptors are not deterministically ordered")
				}
				previous = descriptor.ID
				if descriptor.ID == nounsafe.Rule.ID {
					found = true
					if descriptor.Properties.Severity != severity || descriptor.Properties.DefaultStatus != baseline.DefaultStatus {
						t.Errorf("SARIF descriptor=%+v; want severity=%s and unchanged default status", descriptor, severity)
					}
				}
			}
			if !found {
				t.Error("SARIF omitted configured rule descriptor")
			}
			if len(run.Results) != wantDiagnostics {
				t.Fatalf("SARIF results=%d; want %d", len(run.Results), wantDiagnostics)
			}
			if wantDiagnostics != 0 && (run.Results[0].RuleID != nounsafe.Rule.ID || run.Results[0].Level != test.level) {
				t.Errorf("SARIF result=%+v; want no-unsafe level %s", run.Results[0], test.level)
			}
			if nounsafe.Rule.Severity != baseline.Severity {
				t.Error("configuration mutated exported analyzer rule metadata")
			}
			for _, entry := range registry.Entries() {
				if entry.Rule.ID == nounsafe.Rule.ID && entry.Rule.Severity != baseline.Severity {
					t.Error("configuration mutated the built-in rule registry")
				}
			}
		})
	}
}
