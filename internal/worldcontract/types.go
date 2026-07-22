// Package worldcontract defines the machine-readable boundary of the AESE World.
// It contains contracts only; reducers and runtime behavior belong to later slices.
package worldcontract

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const SchemaVersion = "1.0"

type StableRef struct {
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	Code      string `json:"code"`
}

type Producer struct {
	System    string `json:"system"`
	Component string `json:"component"`
	Version   string `json:"version"`
}

type WorldRun struct {
	SchemaVersion    string `json:"schema_version"`
	WorldRunID       string `json:"world_run_id"`
	TenantID         string `json:"tenant_id"`
	WorldPackKey     string `json:"world_pack_key"`
	WorldPackVersion string `json:"world_pack_version"`
	RulesVersion     string `json:"rules_version"`
	BranchID         string `json:"branch_id"`
	Timezone         string `json:"timezone"`
	Seed             string `json:"seed"`
	Status           string `json:"status"`
	SimTime          string `json:"sim_time"`
	CreatedAt        string `json:"created_at"`
}

type WorldEvent struct {
	SchemaVersion   string          `json:"schema_version"`
	EventID         string          `json:"event_id"`
	EventType       string          `json:"event_type"`
	TenantID        string          `json:"tenant_id"`
	WorldRunID      string          `json:"world_run_id"`
	BranchID        string          `json:"branch_id"`
	SimOccurredAt   string          `json:"sim_occurred_at"`
	RecordedAt      string          `json:"recorded_at"`
	Priority        int             `json:"priority"`
	Sequence        int64           `json:"sequence"`
	CorrelationID   string          `json:"correlation_id"`
	CausationID     string          `json:"causation_id,omitempty"`
	IdempotencyKey  string          `json:"idempotency_key"`
	RulesVersion    string          `json:"rules_version"`
	SubjectRef      StableRef       `json:"subject_ref"`
	PayloadType     string          `json:"payload_type"`
	Payload         json.RawMessage `json:"payload"`
	StateHashBefore string          `json:"state_hash_before"`
	StateHashAfter  string          `json:"state_hash_after"`
}

type Snapshot struct {
	SchemaVersion   string          `json:"schema_version"`
	SnapshotID      string          `json:"snapshot_id"`
	TenantID        string          `json:"tenant_id"`
	WorldRunID      string          `json:"world_run_id"`
	BranchID        string          `json:"branch_id"`
	ThroughSequence int64           `json:"through_sequence"`
	SimTime         string          `json:"sim_time"`
	RulesVersion    string          `json:"rules_version"`
	StateHash       string          `json:"state_hash"`
	State           json.RawMessage `json:"state"`
	CreatedAt       string          `json:"created_at"`
}

type Knowledge struct {
	SchemaVersion   string    `json:"schema_version"`
	KnowledgeID     string    `json:"knowledge_id"`
	TenantID        string    `json:"tenant_id"`
	WorldRunID      string    `json:"world_run_id"`
	BranchID        string    `json:"branch_id"`
	ActorRef        StableRef `json:"actor_ref"`
	FactRef         StableRef `json:"fact_ref"`
	ObservedAt      string    `json:"observed_at"`
	ValidAt         string    `json:"valid_at"`
	SourceRef       string    `json:"source_ref"`
	Confidence      string    `json:"confidence"`
	VisibilityScope string    `json:"visibility_scope"`
	Supersedes      string    `json:"supersedes,omitempty"`
}

type Discrepancy struct {
	SchemaVersion string    `json:"schema_version"`
	DiscrepancyID string    `json:"discrepancy_id"`
	TenantID      string    `json:"tenant_id"`
	WorldRunID    string    `json:"world_run_id"`
	BranchID      string    `json:"branch_id"`
	SubjectRef    StableRef `json:"subject_ref"`
	Kind          string    `json:"kind"`
	WorldFactRef  string    `json:"world_fact_ref"`
	IAOSRecordRef string    `json:"iaos_record_ref,omitempty"`
	KnowledgeRef  string    `json:"knowledge_ref,omitempty"`
	Status        string    `json:"status"`
	DetectedAt    string    `json:"detected_at"`
	ClosedAt      string    `json:"closed_at,omitempty"`
	CorrelationID string    `json:"correlation_id"`
}

type Envelope struct {
	SchemaVersion    string          `json:"schema_version"`
	MessageID        string          `json:"message_id"`
	Kind             string          `json:"kind"`
	TenantID         string          `json:"tenant_id"`
	WorldPackKey     string          `json:"world_pack_key"`
	WorldPackVersion string          `json:"world_pack_version"`
	WorldRunID       string          `json:"world_run_id"`
	BranchID         string          `json:"branch_id"`
	SimOccurredAt    string          `json:"sim_occurred_at"`
	RecordedAt       string          `json:"recorded_at,omitempty"`
	CorrelationID    string          `json:"correlation_id"`
	CausationID      string          `json:"causation_id,omitempty"`
	IdempotencyKey   string          `json:"idempotency_key"`
	Producer         Producer        `json:"producer"`
	SubjectRef       StableRef       `json:"subject_ref"`
	PayloadType      string          `json:"payload_type"`
	Payload          json.RawMessage `json:"payload"`
}

type Observation Envelope
type Intent Envelope
type CommittedOutcome Envelope

func ParseStrict[T any](data []byte) (T, error) {
	var value T
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return value, fmt.Errorf("decode contract: %w", err)
	}
	if err := ensureEOF(decoder); err != nil {
		return value, err
	}
	if validator, ok := any(&value).(interface{ Validate() error }); ok {
		if err := validator.Validate(); err != nil {
			return value, fmt.Errorf("validate contract: %w", err)
		}
	}
	return value, nil
}

func required(values ...string) error {
	for _, value := range values {
		if value == "" {
			return fmt.Errorf("required string is empty")
		}
	}
	return nil
}
func validRef(ref StableRef) error { return required(ref.Namespace, ref.Type, ref.Code) }
func validRaw(value json.RawMessage) error {
	var object map[string]any
	if len(value) == 0 || json.Unmarshal(value, &object) != nil || object == nil {
		return fmt.Errorf("payload/state must be a JSON object")
	}
	return nil
}

func (v *WorldRun) Validate() error {
	if v.SchemaVersion != SchemaVersion || v.Timezone != "Asia/Shanghai" || v.BranchID != "main" {
		return fmt.Errorf("unsupported schema_version, timezone, or branch")
	}
	if err := required(v.WorldRunID, v.TenantID, v.WorldPackKey, v.WorldPackVersion, v.RulesVersion, v.Seed, v.Status); err != nil {
		return err
	}
	return ValidateRFC3339(v.SimTime, v.CreatedAt)
}
func (v *WorldEvent) Validate() error {
	if v.SchemaVersion != SchemaVersion || v.BranchID != "main" || v.Sequence < 0 {
		return fmt.Errorf("unsupported schema_version/branch or negative sequence")
	}
	if err := required(v.EventID, v.EventType, v.TenantID, v.WorldRunID, v.CorrelationID, v.IdempotencyKey, v.RulesVersion, v.PayloadType, v.StateHashBefore, v.StateHashAfter); err != nil {
		return err
	}
	if err := validRef(v.SubjectRef); err != nil {
		return err
	}
	if err := validRaw(v.Payload); err != nil {
		return err
	}
	return ValidateRFC3339(v.SimOccurredAt, v.RecordedAt)
}
func (v *Snapshot) Validate() error {
	if v.SchemaVersion != SchemaVersion || v.BranchID != "main" || v.ThroughSequence < 0 {
		return fmt.Errorf("unsupported schema_version/branch or negative sequence")
	}
	if err := required(v.SnapshotID, v.TenantID, v.WorldRunID, v.RulesVersion, v.StateHash); err != nil {
		return err
	}
	if err := validRaw(v.State); err != nil {
		return err
	}
	return ValidateRFC3339(v.SimTime, v.CreatedAt)
}
func (v *Knowledge) Validate() error {
	if v.SchemaVersion != SchemaVersion || v.BranchID != "main" {
		return fmt.Errorf("unsupported schema_version or branch")
	}
	if err := required(v.KnowledgeID, v.TenantID, v.WorldRunID, v.SourceRef, v.Confidence, v.VisibilityScope); err != nil {
		return err
	}
	if err := validRef(v.ActorRef); err != nil {
		return err
	}
	if err := validRef(v.FactRef); err != nil {
		return err
	}
	return ValidateRFC3339(v.ObservedAt, v.ValidAt)
}
func (v *Discrepancy) Validate() error {
	if v.SchemaVersion != SchemaVersion || v.BranchID != "main" {
		return fmt.Errorf("unsupported schema_version or branch")
	}
	if err := required(v.DiscrepancyID, v.TenantID, v.WorldRunID, v.Kind, v.WorldFactRef, v.Status, v.CorrelationID); err != nil {
		return err
	}
	if err := validRef(v.SubjectRef); err != nil {
		return err
	}
	times := []string{v.DetectedAt}
	if v.ClosedAt != "" {
		times = append(times, v.ClosedAt)
	}
	return ValidateRFC3339(times...)
}
func validateEnvelope(v *Envelope, kind string, recorded bool) error {
	if v.SchemaVersion != SchemaVersion || v.BranchID != "main" || v.Kind != kind {
		return fmt.Errorf("unsupported schema_version, branch, or kind")
	}
	if err := required(v.MessageID, v.TenantID, v.WorldPackKey, v.WorldPackVersion, v.WorldRunID, v.CorrelationID, v.IdempotencyKey, v.PayloadType, v.Producer.System, v.Producer.Component, v.Producer.Version); err != nil {
		return err
	}
	if recorded && v.RecordedAt == "" {
		return fmt.Errorf("recorded_at is required")
	}
	if err := validRef(v.SubjectRef); err != nil {
		return err
	}
	if err := validRaw(v.Payload); err != nil {
		return err
	}
	times := []string{v.SimOccurredAt}
	if v.RecordedAt != "" {
		times = append(times, v.RecordedAt)
	}
	return ValidateRFC3339(times...)
}
func (v *Observation) Validate() error {
	e := Envelope(*v)
	return validateEnvelope(&e, "observation", false)
}
func (v *Intent) Validate() error { e := Envelope(*v); return validateEnvelope(&e, "intent", true) }
func (v *CommittedOutcome) Validate() error {
	e := Envelope(*v)
	return validateEnvelope(&e, "committed_outcome", true)
}

func ensureEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("decode contract: multiple JSON values")
		}
		return fmt.Errorf("decode contract trailing data: %w", err)
	}
	return nil
}

func CanonicalHash(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("marshal canonical JSON: %w", err)
	}
	var normalized any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&normalized); err != nil {
		return "", fmt.Errorf("normalize canonical JSON: %w", err)
	}
	canonical, err := json.Marshal(normalized)
	if err != nil {
		return "", fmt.Errorf("marshal normalized JSON: %w", err)
	}
	sum := sha256.Sum256(canonical)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func ValidateRFC3339(values ...string) error {
	for _, value := range values {
		if _, err := time.Parse(time.RFC3339, value); err != nil {
			return fmt.Errorf("%q is not RFC3339: %w", value, err)
		}
	}
	return nil
}
