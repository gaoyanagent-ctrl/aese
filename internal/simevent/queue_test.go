package simevent

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
)

func TestStableOrderAndDuplicates(t *testing.T) {
	now := time.Date(2026, 7, 8, 0, 0, 0, 0, time.UTC)
	q := New()
	add := func(id string, p int, at time.Time) error {
		return q.Schedule(Scheduled{EventID: id, EventType: "state.set.v1", SimOccurredAt: at.Format(time.RFC3339), Priority: p, CorrelationID: "c", IdempotencyKey: "k-" + id, SubjectRef: worldcontract.StableRef{Namespace: "hctm", Type: "equipment", Code: "E1"}, PayloadType: "state.set.v1", Payload: json.RawMessage(`{}`)}, now)
	}
	for _, v := range []struct {
		id string
		p  int
		d  time.Duration
	}{{"b", 2, time.Hour}, {"c", 1, time.Hour}, {"a", 1, time.Hour}, {"early", 9, 0}} {
		if err := add(v.id, v.p, now.Add(v.d)); err != nil {
			t.Fatal(err)
		}
	}
	want := []string{"early", "a", "c", "b"}
	for _, id := range want {
		got, _ := q.Pop()
		if got.EventID != id {
			t.Fatalf("got %s want %s", got.EventID, id)
		}
	}
	if err := add("past", 1, now.Add(-time.Second)); err == nil {
		t.Fatal("past accepted")
	}
	if err := add("b", 1, now); err == nil {
		t.Fatal("duplicate accepted")
	}
}
