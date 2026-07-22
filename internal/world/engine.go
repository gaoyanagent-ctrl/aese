// Package world composes virtual time, event ordering and pure reducers.
package world

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/rules"
	"github.com/industrial-ai/iaos-aese/internal/simevent"
	"github.com/industrial-ai/iaos-aese/internal/simtime"
	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
)

type Engine struct {
	run          worldcontract.WorldRun
	clock        *simtime.Clock
	queue        *simevent.Queue
	rules        *rules.Registry
	state        rules.State
	log          []worldcontract.WorldEvent
	baseSequence int64
}

func New(run worldcontract.WorldRun, initial rules.State, events []simevent.Scheduled) (*Engine, error) {
	if err := run.Validate(); err != nil {
		return nil, err
	}
	at, err := time.Parse(time.RFC3339, run.SimTime)
	if err != nil {
		return nil, err
	}
	clock, err := simtime.New(at)
	if err != nil {
		return nil, err
	}
	registry, err := rules.New(run.RulesVersion)
	if err != nil {
		return nil, err
	}
	e := &Engine{run: run, clock: clock, queue: simevent.New(), rules: registry, state: clone(initial), log: []worldcontract.WorldEvent{}}
	for _, event := range events {
		if err := e.queue.Schedule(event, at); err != nil {
			return nil, err
		}
	}
	seen := map[string]struct{}{}
	for _, event := range simevent.Ordered(events) {
		if event.CausationID != "" {
			if _, ok := seen[event.CausationID]; !ok {
				return nil, fmt.Errorf("event %s causation_id must reference an earlier event", event.EventID)
			}
		}
		seen[event.EventID] = struct{}{}
	}
	return e, nil
}

func NewFromSnapshot(run worldcontract.WorldRun, snapshot worldcontract.Snapshot, remaining []simevent.Scheduled) (*Engine, error) {
	if err := snapshot.Validate(); err != nil {
		return nil, err
	}
	if snapshot.TenantID != run.TenantID || snapshot.WorldRunID != run.WorldRunID || snapshot.BranchID != run.BranchID || snapshot.RulesVersion != run.RulesVersion {
		return nil, fmt.Errorf("snapshot does not belong to run/rules")
	}
	var state rules.State
	if err := json.Unmarshal(snapshot.State, &state); err != nil {
		return nil, fmt.Errorf("snapshot state: %w", err)
	}
	hash, err := worldcontract.CanonicalHash(state)
	if err != nil {
		return nil, err
	}
	if hash != snapshot.StateHash {
		return nil, fmt.Errorf("snapshot state hash mismatch")
	}
	run.SimTime = snapshot.SimTime
	engine, err := New(run, state, remaining)
	if err != nil {
		return nil, err
	}
	engine.baseSequence = snapshot.ThroughSequence
	return engine, nil
}
func clone(state rules.State) rules.State {
	out := rules.State{}
	for k, v := range state {
		out[k] = append(json.RawMessage(nil), v...)
	}
	return out
}
func (e *Engine) Step() (worldcontract.WorldEvent, bool, error) {
	scheduled, ok := e.queue.Pop()
	if !ok {
		return worldcontract.WorldEvent{}, false, nil
	}
	at, _ := scheduled.Time()
	if err := e.clock.Step(at); err != nil {
		return worldcontract.WorldEvent{}, false, err
	}
	before, err := worldcontract.CanonicalHash(e.state)
	if err != nil {
		return worldcontract.WorldEvent{}, false, err
	}
	next, err := e.rules.Reduce(e.state, scheduled)
	if err != nil {
		return worldcontract.WorldEvent{}, false, err
	}
	after, err := worldcontract.CanonicalHash(next)
	if err != nil {
		return worldcontract.WorldEvent{}, false, err
	}
	event := worldcontract.WorldEvent{SchemaVersion: worldcontract.SchemaVersion, EventID: scheduled.EventID, EventType: scheduled.EventType, TenantID: e.run.TenantID, WorldRunID: e.run.WorldRunID, BranchID: e.run.BranchID, SimOccurredAt: scheduled.SimOccurredAt, RecordedAt: e.run.CreatedAt, Priority: scheduled.Priority, Sequence: e.baseSequence + int64(len(e.log)+1), CorrelationID: scheduled.CorrelationID, CausationID: scheduled.CausationID, IdempotencyKey: scheduled.IdempotencyKey, RulesVersion: e.rules.Version(), SubjectRef: scheduled.SubjectRef, PayloadType: scheduled.PayloadType, Payload: scheduled.Payload, StateHashBefore: before, StateHashAfter: after}
	e.state = next
	e.log = append(e.log, event)
	return event, true, nil
}
func (e *Engine) RunUntil(until time.Time) error {
	for {
		next, ok := e.queue.Peek()
		if !ok {
			return nil
		}
		at, _ := next.Time()
		if at.After(until) {
			return e.clock.Step(until)
		}
		if _, _, err := e.Step(); err != nil {
			return err
		}
	}
}
func (e *Engine) RunAll() error {
	for e.queue.Len() > 0 {
		if _, _, err := e.Step(); err != nil {
			return err
		}
	}
	return nil
}
func (e *Engine) Log() []worldcontract.WorldEvent {
	return append([]worldcontract.WorldEvent(nil), e.log...)
}
func (e *Engine) State() rules.State { return clone(e.state) }
func (e *Engine) Snapshot() (worldcontract.Snapshot, error) {
	hash, err := worldcontract.CanonicalHash(e.state)
	if err != nil {
		return worldcontract.Snapshot{}, err
	}
	sequence := e.baseSequence + int64(len(e.log))
	return worldcontract.Snapshot{SchemaVersion: worldcontract.SchemaVersion, SnapshotID: fmt.Sprintf("%s-%s-%d", e.run.WorldRunID, e.run.BranchID, sequence), TenantID: e.run.TenantID, WorldRunID: e.run.WorldRunID, BranchID: e.run.BranchID, ThroughSequence: sequence, SimTime: e.clock.Now().Format(time.RFC3339), RulesVersion: e.rules.Version(), StateHash: hash, State: e.mustState(), CreatedAt: e.run.CreatedAt}, nil
}
func (e *Engine) mustState() json.RawMessage { data, _ := json.Marshal(e.state); return data }
func Replay(run worldcontract.WorldRun, initial rules.State, events []simevent.Scheduled, expected []worldcontract.WorldEvent) (*Engine, error) {
	engine, err := New(run, initial, events)
	if err != nil {
		return nil, err
	}
	if err := engine.RunAll(); err != nil {
		return nil, err
	}
	actual := engine.Log()
	if len(actual) != len(expected) {
		return nil, fmt.Errorf("event count mismatch: got %d want %d", len(actual), len(expected))
	}
	for i := range actual {
		if actual[i].EventID != expected[i].EventID || actual[i].StateHashBefore != expected[i].StateHashBefore || actual[i].StateHashAfter != expected[i].StateHashAfter {
			return nil, fmt.Errorf("event %d replay mismatch", i+1)
		}
	}
	return engine, nil
}
