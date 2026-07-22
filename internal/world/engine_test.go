package world

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/rules"
	"github.com/industrial-ai/iaos-aese/internal/simevent"
	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
)

func testRun() worldcontract.WorldRun {
	return worldcontract.WorldRun{SchemaVersion: "1.0", WorldRunID: "r1", TenantID: "t1", WorldPackKey: "p", WorldPackVersion: "1", RulesVersion: "rules-0.1.0", BranchID: "main", Timezone: "Asia/Shanghai", Seed: "1", Status: "paused", SimTime: "2026-07-08T10:00:00+08:00", CreatedAt: "2026-07-22T04:00:00Z"}
}
func scheduled(id, at, key string) simevent.Scheduled {
	return simevent.Scheduled{EventID: id, EventType: "state.changed.v1", SimOccurredAt: at, Priority: 1, CorrelationID: "corr", IdempotencyKey: "idem-" + id, SubjectRef: worldcontract.StableRef{Namespace: "hctm", Type: "equipment", Code: "LAS-WLD-02"}, PayloadType: "state.set.v1", Payload: json.RawMessage(`{"key":"` + key + `","value":"degraded"}`)}
}
func TestDeterministicRunReplayAndSnapshot(t *testing.T) {
	events := []simevent.Scheduled{scheduled("b", "2026-07-08T10:01:00+08:00", "b"), scheduled("a", "2026-07-08T10:01:00+08:00", "a")}
	first, err := New(testRun(), rules.State{}, events)
	if err != nil {
		t.Fatal(err)
	}
	if err := first.RunAll(); err != nil {
		t.Fatal(err)
	}
	second, err := New(testRun(), rules.State{}, events)
	if err != nil {
		t.Fatal(err)
	}
	if err := second.RunAll(); err != nil {
		t.Fatal(err)
	}
	a, _ := worldcontract.CanonicalHash(first.Log())
	b, _ := worldcontract.CanonicalHash(second.Log())
	if a != b {
		t.Fatalf("non-deterministic logs %s %s", a, b)
	}
	if first.Log()[0].EventID != "a" {
		t.Fatalf("unstable tie break: %s", first.Log()[0].EventID)
	}
	if _, err := Replay(testRun(), rules.State{}, events, first.Log()); err != nil {
		t.Fatal(err)
	}
	snap, _ := first.Snapshot()
	if snap.ThroughSequence != 2 {
		t.Fatalf("snapshot sequence %d", snap.ThroughSequence)
	}
}
func TestRunUntilAndUnknownRuleFailClosed(t *testing.T) {
	events := []simevent.Scheduled{scheduled("a", "2026-07-08T10:01:00+08:00", "a"), scheduled("b", "2026-07-08T10:03:00+08:00", "b")}
	engine, _ := New(testRun(), rules.State{}, events)
	until, _ := time.Parse(time.RFC3339, "2026-07-08T10:02:00+08:00")
	if err := engine.RunUntil(until); err != nil {
		t.Fatal(err)
	}
	if len(engine.Log()) != 1 {
		t.Fatalf("got %d events", len(engine.Log()))
	}
	bad := scheduled("x", "2026-07-08T10:04:00+08:00", "x")
	bad.PayloadType = "unknown.v1"
	engine, _ = New(testRun(), rules.State{}, []simevent.Scheduled{bad})
	if err := engine.RunAll(); err == nil {
		t.Fatal("unknown rule accepted")
	}
}
func TestBundleLoads(t *testing.T) {
	bundle, err := LoadBundle("../../world-contracts/runtime-example")
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Events) != 2 {
		t.Fatalf("events %d", len(bundle.Events))
	}
}

func TestOneHundredRunsAreIdentical(t *testing.T) {
	events := []simevent.Scheduled{scheduled("year", "2027-07-08T10:01:00+08:00", "year"), scheduled("second", "2026-07-08T10:00:01+08:00", "second")}
	var expected string
	for i := 0; i < 100; i++ {
		engine, err := New(testRun(), rules.State{}, events)
		if err != nil {
			t.Fatal(err)
		}
		if err := engine.RunAll(); err != nil {
			t.Fatal(err)
		}
		hash, _ := worldcontract.CanonicalHash(engine.Log())
		if i == 0 {
			expected = hash
		} else if hash != expected {
			t.Fatalf("run %d differs", i+1)
		}
	}
}

func TestRestoreFromSnapshotContinuesSequence(t *testing.T) {
	first := scheduled("a", "2026-07-08T10:01:00+08:00", "a")
	second := scheduled("b", "2026-07-08T10:02:00+08:00", "b")
	engine, _ := New(testRun(), rules.State{}, []simevent.Scheduled{first})
	if err := engine.RunAll(); err != nil {
		t.Fatal(err)
	}
	snapshot, _ := engine.Snapshot()
	restored, err := NewFromSnapshot(testRun(), snapshot, []simevent.Scheduled{second})
	if err != nil {
		t.Fatal(err)
	}
	if err := restored.RunAll(); err != nil {
		t.Fatal(err)
	}
	if restored.Log()[0].Sequence != 2 {
		t.Fatalf("sequence %d", restored.Log()[0].Sequence)
	}
	snapshot.StateHash = "sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	if _, err := NewFromSnapshot(testRun(), snapshot, nil); err == nil {
		t.Fatal("corrupt snapshot accepted")
	}
}

func TestCausationMustReferenceEarlierEvent(t *testing.T) {
	cause := scheduled("cause", "2026-07-08T10:02:00+08:00", "cause")
	effect := scheduled("effect", "2026-07-08T10:01:00+08:00", "effect")
	effect.CausationID = cause.EventID
	if _, err := New(testRun(), rules.State{}, []simevent.Scheduled{cause, effect}); err == nil {
		t.Fatal("forward causation accepted")
	}
}
