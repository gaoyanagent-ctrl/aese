package rules

import (
	"encoding/json"
	"github.com/industrial-ai/iaos-aese/internal/simevent"
	"testing"
)

func TestReducerIsPureAndVersioned(t *testing.T) {
	registry, err := New("rules-0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	initial := State{"existing": json.RawMessage(`"unchanged"`)}
	event := simevent.Scheduled{PayloadType: "state.set.v1", Payload: json.RawMessage(`{"key":"equipment","value":{"condition":"degraded"}}`)}
	next, err := registry.Reduce(initial, event)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := initial["equipment"]; ok {
		t.Fatal("input state mutated")
	}
	if string(next["existing"]) != `"unchanged"` {
		t.Fatal("existing state lost")
	}
	event.PayloadType = "unknown.v1"
	if _, err := registry.Reduce(initial, event); err == nil {
		t.Fatal("unknown rule accepted")
	}
	if _, err := New("rules-9.9.9"); err == nil {
		t.Fatal("unknown rules version accepted")
	}
}
