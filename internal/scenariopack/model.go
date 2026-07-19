package scenariopack

import "time"

const SupportedSchemaVersion = "1.0.0"

type Manifest struct {
	SchemaVersion  string     `json:"schema_version"`
	PackKey        string     `json:"pack_key"`
	PackVersion    string     `json:"pack_version"`
	DisplayName    string     `json:"display_name"`
	Timezone       string     `json:"timezone"`
	TenantTemplate string     `json:"tenant_template"`
	MasterData     []EntryRef `json:"master_data"`
	Stories        []StoryRef `json:"stories"`
}

type EntryRef struct {
	Key  string `json:"key,omitempty"`
	Path string `json:"path"`
}

type StoryRef struct {
	Key              string `json:"key"`
	InitialState     string `json:"initial_state"`
	Events           string `json:"events"`
	ExpectedOutcomes string `json:"expected_outcomes"`
}

type RecordSet struct {
	SchemaVersion string           `json:"schema_version"`
	Entity        string           `json:"entity"`
	NaturalKey    []string         `json:"natural_key"`
	Records       []map[string]any `json:"records"`
	Source        string           `json:"-"`
}

type InitialState struct {
	SchemaVersion string           `json:"schema_version"`
	StoryKey      string           `json:"story_key"`
	RecordSets    []RecordSet      `json:"record_sets,omitempty"`
	Records       []map[string]any `json:"records,omitempty"`
	Source        string           `json:"-"`
}

type EventSequence struct {
	SchemaVersion string  `json:"schema_version"`
	StoryKey      string  `json:"story_key"`
	CorrelationID string  `json:"correlation_id,omitempty"`
	Events        []Event `json:"events"`
	Source        string  `json:"-"`
}

type Event struct {
	EventID         string         `json:"event_id"`
	EventType       string         `json:"event_type"`
	Subject         string         `json:"subject,omitempty"`
	Timestamp       string         `json:"timestamp"`
	TenantID        string         `json:"tenant_id,omitempty"`
	CorrelationID   string         `json:"correlation_id,omitempty"`
	CausationID     string         `json:"causation_id,omitempty"`
	IdempotencyKey  string         `json:"idempotency_key,omitempty"`
	AggregateType   string         `json:"aggregate_type,omitempty"`
	AggregateID     string         `json:"aggregate_id,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
	Payload         map[string]any `json:"payload"`
	ConcurrentGroup string         `json:"concurrent_group,omitempty"`
}

func (e Event) Correlation() string {
	return first(e.CorrelationID, stringValue(e.Metadata, "correlation_id"))
}
func (e Event) Causation() string {
	return first(e.CausationID, stringValue(e.Metadata, "causation_id"))
}
func (e Event) Idempotency() string {
	return first(e.IdempotencyKey, stringValue(e.Metadata, "idempotency_key"))
}
func (e Event) Time() (time.Time, error) { return time.Parse(time.RFC3339, e.Timestamp) }

type ExpectedOutcomes struct {
	SchemaVersion  string         `json:"schema_version"`
	StoryKey       string         `json:"story_key"`
	Assertions     []Assertion    `json:"assertions"`
	IAOSAssertions []Assertion    `json:"iaos_assertions,omitempty"`
	Summary        map[string]any `json:"summary,omitempty"`
	Source         string         `json:"-"`
}

type Assertion struct {
	Key      string         `json:"key"`
	Type     string         `json:"type,omitempty"`
	Entity   string         `json:"entity,omitempty"`
	Match    map[string]any `json:"match,omitempty"`
	Field    string         `json:"field,omitempty"`
	Operator string         `json:"operator,omitempty"`
	Expected any            `json:"expected,omitempty"`
}

type Story struct {
	Ref      StoryRef
	Initial  InitialState
	Events   EventSequence
	Expected ExpectedOutcomes
}

type Pack struct {
	Root       string
	Manifest   Manifest
	RecordSets []RecordSet
	Stories    []Story
}

func first(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func stringValue(values map[string]any, key string) string {
	if values == nil {
		return ""
	}
	value, _ := values[key].(string)
	return value
}
