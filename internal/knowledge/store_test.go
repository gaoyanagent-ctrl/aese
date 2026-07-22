package knowledge

import (
	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
	"testing"
	"time"
)

func TestKnowledgeNotVisibleBeforeObservation(t *testing.T) {
	s := New()
	actor := worldcontract.StableRef{Namespace: "hctm", Type: "position", Code: "EQUIPMENT-ENGINEER-SZ"}
	r := worldcontract.Knowledge{SchemaVersion: "1.0", KnowledgeID: "k1", TenantID: "t", WorldRunID: "r", BranchID: "main", ActorRef: actor, FactRef: worldcontract.StableRef{Namespace: "hctm", Type: "world_event", Code: "e1"}, ObservedAt: "2026-07-08T10:16:00+08:00", ValidAt: "2026-07-08T10:15:00+08:00", SourceRef: "sensor:s1", Confidence: "0.95", VisibilityScope: "assigned_recipients"}
	if err := s.Learn(r); err != nil {
		t.Fatal(err)
	}
	before, _ := time.Parse(time.RFC3339, "2026-07-08T10:15:59+08:00")
	if len(s.Visible(actor, before)) != 0 {
		t.Fatal("future knowledge leaked")
	}
	after := before.Add(time.Second)
	if len(s.Visible(actor, after)) != 1 {
		t.Fatal("knowledge unavailable")
	}
}
