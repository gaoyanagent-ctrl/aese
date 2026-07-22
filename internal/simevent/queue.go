// Package simevent provides stable ordering for scheduled world events.
package simevent

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
)

type Scheduled struct {
	EventID        string                  `json:"event_id"`
	EventType      string                  `json:"event_type"`
	SimOccurredAt  string                  `json:"sim_occurred_at"`
	Priority       int                     `json:"priority"`
	CorrelationID  string                  `json:"correlation_id"`
	CausationID    string                  `json:"causation_id,omitempty"`
	IdempotencyKey string                  `json:"idempotency_key"`
	SubjectRef     worldcontract.StableRef `json:"subject_ref"`
	PayloadType    string                  `json:"payload_type"`
	Payload        json.RawMessage         `json:"payload"`
	index          int
}

func (s Scheduled) Time() (time.Time, error) { return time.Parse(time.RFC3339, s.SimOccurredAt) }
func (s Scheduled) Validate() error {
	if s.EventID == "" || s.IdempotencyKey == "" || s.EventType == "" || s.CorrelationID == "" || s.PayloadType == "" {
		return fmt.Errorf("event identity fields are required")
	}
	if s.SubjectRef.Namespace == "" || s.SubjectRef.Type == "" || s.SubjectRef.Code == "" {
		return fmt.Errorf("subject_ref is required")
	}
	var object map[string]any
	if len(s.Payload) == 0 || json.Unmarshal(s.Payload, &object) != nil || object == nil {
		return fmt.Errorf("payload must be a JSON object")
	}
	_, err := s.Time()
	return err
}
func Ordered(events []Scheduled) []Scheduled {
	out := append([]Scheduled(nil), events...)
	sort.Slice(out, func(i, j int) bool {
		a, _ := out[i].Time()
		b, _ := out[j].Time()
		if !a.Equal(b) {
			return a.Before(b)
		}
		if out[i].Priority != out[j].Priority {
			return out[i].Priority < out[j].Priority
		}
		return out[i].EventID < out[j].EventID
	})
	return out
}

type items []Scheduled

func (q items) Len() int { return len(q) }
func (q items) Less(i, j int) bool {
	a, _ := q[i].Time()
	b, _ := q[j].Time()
	if !a.Equal(b) {
		return a.Before(b)
	}
	if q[i].Priority != q[j].Priority {
		return q[i].Priority < q[j].Priority
	}
	return q[i].EventID < q[j].EventID
}
func (q items) Swap(i, j int) { q[i], q[j] = q[j], q[i]; q[i].index = i; q[j].index = j }
func (q *items) Push(v any)   { s := v.(Scheduled); s.index = len(*q); *q = append(*q, s) }
func (q *items) Pop() any     { old := *q; n := len(old); s := old[n-1]; *q = old[:n-1]; return s }

type Queue struct {
	items items
	ids   map[string]struct{}
	keys  map[string]struct{}
}

func New() *Queue { return &Queue{ids: map[string]struct{}{}, keys: map[string]struct{}{}} }
func (q *Queue) Schedule(event Scheduled, now time.Time) error {
	if err := event.Validate(); err != nil {
		return fmt.Errorf("event %s: %w", event.EventID, err)
	}
	at, err := event.Time()
	if err != nil {
		return fmt.Errorf("event %s time: %w", event.EventID, err)
	}
	if at.Before(now) {
		return fmt.Errorf("event %s is before current virtual time", event.EventID)
	}
	if _, ok := q.ids[event.EventID]; ok {
		return fmt.Errorf("duplicate event_id %s", event.EventID)
	}
	if _, ok := q.keys[event.IdempotencyKey]; ok {
		return fmt.Errorf("duplicate idempotency_key %s", event.IdempotencyKey)
	}
	q.ids[event.EventID] = struct{}{}
	q.keys[event.IdempotencyKey] = struct{}{}
	heap.Push(&q.items, event)
	return nil
}
func (q *Queue) Peek() (Scheduled, bool) {
	if len(q.items) == 0 {
		return Scheduled{}, false
	}
	return q.items[0], true
}
func (q *Queue) Pop() (Scheduled, bool) {
	if len(q.items) == 0 {
		return Scheduled{}, false
	}
	return heap.Pop(&q.items).(Scheduled), true
}
func (q *Queue) Len() int { return len(q.items) }
