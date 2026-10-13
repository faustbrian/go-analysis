package analysis

// Diagnostic is the stable reporting form of one policy violation.
type Diagnostic struct {
	// Rule identifies the reporting rule without implicitly validating that ID.
	Rule string `json:"rule"`
	// Package is the analyzed import path, when available.
	Package string `json:"package,omitempty"`
	// Filename is a source path; report writers reject traversal and relativize
	// absolute paths against Root.
	Filename string `json:"filename"`
	// Line is a one-based source line; writers do not validate its range.
	Line int `json:"line"`
	// Column is a one-based source byte column; zero means unspecified.
	Column int `json:"column,omitempty"`
	// Message is emitted verbatim after JSON escaping; callers must omit secrets or
	// source excerpts.
	Message string `json:"message"`
	// Rationale is optional policy context and is serialized by WriteJSON.
	Rationale string `json:"rationale,omitempty"`
	// Remediation is optional guidance and is serialized by WriteJSON.
	Remediation string `json:"remediation,omitempty"`
}
