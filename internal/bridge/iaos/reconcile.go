package iaos

import (
	"encoding/json"
	"sort"

	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
)

type ReconciliationIssue struct {
	Kind          string `json:"kind"`
	MessageID     string `json:"message_id,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	Detail        string `json:"detail"`
}

type ReconciliationReport struct {
	Converged bool                  `json:"converged"`
	Counts    map[string]int        `json:"counts"`
	Issues    []ReconciliationIssue `json:"issues"`
}

// Reconcile is a deterministic, read-only comparison over the durable bridge
// journal. Notification order is irrelevant; convergence is established by
// correlation chains and stable message identifiers.
func Reconcile(entries []worldcontract.Envelope) ReconciliationReport {
	report := ReconciliationReport{Converged: true, Counts: map[string]int{"intent": 0, "observation": 0, "committed_outcome": 0}, Issues: []ReconciliationIssue{}}
	seen := map[string]bool{}
	byCorrelation := map[string]map[string]int{}
	hashes := map[string]map[string]string{}
	for _, entry := range entries {
		report.Counts[entry.Kind]++
		if seen[entry.MessageID] {
			report.Issues = append(report.Issues, ReconciliationIssue{"duplicate", entry.MessageID, entry.CorrelationID, "message_id repeated"})
		}
		seen[entry.MessageID] = true
		if byCorrelation[entry.CorrelationID] == nil {
			byCorrelation[entry.CorrelationID] = map[string]int{}
		}
		byCorrelation[entry.CorrelationID][entry.Kind]++
		var payload map[string]any
		if json.Unmarshal(entry.Payload, &payload) == nil {
			value, _ := payload["state_hash"].(string)
			if value == "" {
				value, _ = payload["canonical_hash"].(string)
			}
			if value != "" {
				if hashes[entry.CorrelationID] == nil {
					hashes[entry.CorrelationID] = map[string]string{}
				}
				if prior := hashes[entry.CorrelationID][entry.Kind]; prior != "" && prior != value {
					report.Issues = append(report.Issues, ReconciliationIssue{"hash_mismatch", entry.MessageID, entry.CorrelationID, "same-kind canonical hashes disagree"})
				}
				hashes[entry.CorrelationID][entry.Kind] = value
			}
		}
	}
	for correlation, kinds := range byCorrelation {
		if kinds["intent"] > 0 && kinds["observation"] == 0 {
			report.Issues = append(report.Issues, ReconciliationIssue{"lagging", "", correlation, "intent has no observation"})
		}
		if kinds["observation"] > 0 && kinds["intent"] == 0 {
			report.Issues = append(report.Issues, ReconciliationIssue{"missing", "", correlation, "observation has no intent"})
		}
		if kinds["observation"] > 0 && kinds["committed_outcome"] == 0 {
			report.Issues = append(report.Issues, ReconciliationIssue{"lagging", "", correlation, "observation has no committed outcome"})
		}
		if kinds["committed_outcome"] > 0 && kinds["observation"] == 0 {
			report.Issues = append(report.Issues, ReconciliationIssue{"terminal_conflict", "", correlation, "outcome has no observation"})
		}
		if observation, outcome := hashes[correlation]["observation"], hashes[correlation]["committed_outcome"]; observation != "" && outcome != "" && observation != outcome {
			report.Issues = append(report.Issues, ReconciliationIssue{"hash_mismatch", "", correlation, "observation and committed outcome hashes disagree"})
		}
	}
	sort.Slice(report.Issues, func(i, j int) bool {
		a, b := report.Issues[i], report.Issues[j]
		if a.CorrelationID != b.CorrelationID {
			return a.CorrelationID < b.CorrelationID
		}
		if a.Kind != b.Kind {
			return a.Kind < b.Kind
		}
		return a.MessageID < b.MessageID
	})
	report.Converged = len(report.Issues) == 0
	return report
}
