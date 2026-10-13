// Package policy owns validated rule governance and lexical inventories. Builtin
// selects the shipped rules; NewRegistry accepts explicit metadata without
// executing analyzers or authenticating owner/promotion/overlap claims.
//
// Construction borrows nested rule Properties maps and Overlaps slices. Treat
// those caller-owned values as immutable while using a registry, including across
// goroutines. Entries returns copies of those collections, not a live view.
package policy

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	shared "github.com/faustbrian/go-analysis/analysis"
)

// Overlap documents another tool covering adjacent or identical semantics.
type Overlap struct {
	// Tool names the adjacent authority; trimmed case-insensitive duplicates are
	// rejected.
	Tool string `json:"tool"`
	// CanonicalAuthority is required explanatory authority text, not executable tool
	// delegation.
	CanonicalAuthority string `json:"canonical_authority"`
	// CompatibleConfig is optional configuration guidance, not validated against the
	// named tool.
	CompatibleConfig string `json:"compatible_configuration,omitempty"`
}

// Entry associates a rule with ownership and promotion evidence.
type Entry struct {
	// Rule contains governed metadata; nested Properties storage is borrowed at
	// construction.
	Rule shared.Rule `json:"rule"`
	// Owner is required attribution text, not an authenticated principal.
	Owner string `json:"owner"`
	// PromotionEvidence is required for blocking defaults, but no linked artifact is
	// fetched or verified.
	PromotionEvidence string `json:"promotion_evidence,omitempty"`
	// Overlaps records adjacent tool authority; its backing slice is borrowed at
	// construction.
	Overlaps []Overlap `json:"overlaps,omitempty"`
}

// Registry is a validated inventory with no mutation methods. NewRegistry
// retains nested input storage, so immutability/concurrent read safety requires
// callers not to mutate the borrowed maps or slices. A zero Registry reads empty;
// nil receiver calls are not supported.
type Registry struct {
	entries map[string]Entry
}

// NewRegistry rejects invalid shared rule metadata, blank owners, blocking
// defaults without promotion text, duplicate IDs and duplicate overlap tools.
// It returns nil on error and does not mutate inputs, but retains each entry
// by value with borrowed nested maps/slices. Evidence strings are not verified.
// Empty input produces an empty registry.
func NewRegistry(entries []Entry) (*Registry, error) {
	registry := &Registry{entries: make(map[string]Entry, len(entries))}
	for _, entry := range entries {
		if err := entry.Rule.Validate(); err != nil {
			return nil, fmt.Errorf("validate rule metadata: %w", err)
		}
		if strings.TrimSpace(entry.Owner) == "" {
			return nil, fmt.Errorf("rule %q requires an owner", entry.Rule.ID)
		}
		if entry.Rule.DefaultStatus == shared.StatusBlocking &&
			strings.TrimSpace(entry.PromotionEvidence) == "" {
			return nil, fmt.Errorf("blocking rule %q requires promotion evidence", entry.Rule.ID)
		}
		overlapTools := make(map[string]struct{}, len(entry.Overlaps))
		for _, overlap := range entry.Overlaps {
			tool := strings.TrimSpace(overlap.Tool)
			if tool == "" ||
				strings.TrimSpace(overlap.CanonicalAuthority) == "" {
				return nil, errors.New("overlap requires tool and canonical authority")
			}
			toolKey := strings.ToLower(tool)
			if _, exists := overlapTools[toolKey]; exists {
				return nil, fmt.Errorf(
					"rule %q duplicates overlap authority %q",
					entry.Rule.ID,
					tool,
				)
			}
			overlapTools[toolKey] = struct{}{}
		}
		if _, exists := registry.entries[entry.Rule.ID]; exists {
			return nil, fmt.Errorf("duplicate rule %q", entry.Rule.ID)
		}
		// Only the record is copied here; nested metadata retains caller ownership.
		registry.entries[entry.Rule.ID] = entry
	}

	return registry, nil
}

// IDs returns rule IDs in stable lexical order.
func (registry *Registry) IDs() []string {
	ids := make([]string, 0, len(registry.entries))
	for id := range registry.entries {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	return ids
}

// Entries returns ID-sorted entry copies, including fresh Overlaps slices and
// Configuration.Properties maps. Callers may mutate those returned collections
// independently; other record fields are scalar values.
func (registry *Registry) Entries() []Entry {
	ids := registry.IDs()
	entries := make([]Entry, 0, len(ids))
	for _, id := range ids {
		entry := registry.entries[id]
		entry.Overlaps = append([]Overlap(nil), entry.Overlaps...)
		properties := entry.Rule.Configuration.Properties
		entry.Rule.Configuration.Properties = make(
			map[string]shared.ConfigurationProperty,
			len(properties),
		)
		for name, property := range properties {
			entry.Rule.Configuration.Properties[name] = property
		}
		entries = append(entries, entry)
	}

	return entries
}
