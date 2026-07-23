package iaos

import (
	"encoding/json"
	"testing"

	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
)

func TestReconcileConvergedAndFailureClasses(t *testing.T) {
	entry := func(id, kind, correlation string) worldcontract.Envelope {
		return worldcontract.Envelope{MessageID: id, Kind: kind, CorrelationID: correlation}
	}
	ok := Reconcile([]worldcontract.Envelope{entry("i", "intent", "c"), entry("o", "observation", "c"), entry("r", "committed_outcome", "c")})
	if !ok.Converged {
		t.Fatalf("issues=%v", ok.Issues)
	}
	bad := Reconcile([]worldcontract.Envelope{entry("i", "intent", "lag"), entry("o", "observation", "orphan"), entry("r", "committed_outcome", "terminal"), entry("i", "intent", "lag")})
	if bad.Converged || len(bad.Issues) < 4 {
		t.Fatalf("report=%+v", bad)
	}
	left, right := entry("o2", "observation", "hash"), entry("r2", "committed_outcome", "hash")
	intent := entry("i2", "intent", "hash")
	left.Payload, right.Payload = json.RawMessage(`{"state_hash":"a"}`), json.RawMessage(`{"state_hash":"b"}`)
	mismatch := Reconcile([]worldcontract.Envelope{intent, left, right})
	if mismatch.Converged || mismatch.Issues[0].Kind != "hash_mismatch" {
		t.Fatalf("expected hash mismatch: %+v", mismatch)
	}
	shuffled := Reconcile([]worldcontract.Envelope{entry("r3", "committed_outcome", "recover"), entry("i3", "intent", "recover"), entry("o3", "observation", "recover")})
	if !shuffled.Converged {
		t.Fatalf("delivery order must not affect durable convergence: %+v", shuffled)
	}
	lagging := Reconcile([]worldcontract.Envelope{entry("i4", "intent", "delayed")})
	recovered := Reconcile([]worldcontract.Envelope{entry("i4", "intent", "delayed"), entry("o4", "observation", "delayed"), entry("r4", "committed_outcome", "delayed")})
	if lagging.Converged || !recovered.Converged {
		t.Fatalf("delayed delivery did not converge: before=%+v after=%+v", lagging, recovered)
	}
}
