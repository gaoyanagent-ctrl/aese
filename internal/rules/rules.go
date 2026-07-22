// Package rules contains versioned pure world reducers.
package rules

import (
	"encoding/json"
	"fmt"

	"github.com/industrial-ai/iaos-aese/internal/simevent"
)

type State map[string]json.RawMessage
type Reducer func(State, simevent.Scheduled) (State, error)
type Registry struct {
	version  string
	reducers map[string]Reducer
}
const SupportedVersion = "rules-0.1.0"

func New(version string) (*Registry, error) {
	if version != SupportedVersion {
		return nil, fmt.Errorf("unsupported rules version %q", version)
	}
	r := &Registry{version: version, reducers: map[string]Reducer{}}
	r.reducers["state.set.v1"] = setState
	return r, nil
}
func (r *Registry) Version() string { return r.version }
func (r *Registry) Reduce(state State, event simevent.Scheduled) (State, error) {
	reduce, ok := r.reducers[event.PayloadType]
	if !ok {
		return nil, fmt.Errorf("unsupported payload_type %q for rules %s", event.PayloadType, r.version)
	}
	return reduce(clone(state), event)
}
func clone(state State) State {
	out := make(State, len(state))
	for k, v := range state {
		out[k] = append(json.RawMessage(nil), v...)
	}
	return out
}
func setState(state State, event simevent.Scheduled) (State, error) {
	var payload struct {
		Key   string          `json:"key"`
		Value json.RawMessage `json:"value"`
	}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return nil, fmt.Errorf("state.set payload: %w", err)
	}
	if payload.Key == "" || len(payload.Value) == 0 {
		return nil, fmt.Errorf("state.set requires key and value")
	}
	state[payload.Key] = append(json.RawMessage(nil), payload.Value...)
	return state, nil
}
