package gameprojection

import (
	"testing"

	"github.com/industrial-ai/iaos-aese/internal/iaosclient"
)

func TestFromIAOSUsesCommittedWorkItemsAndWorldCursor(t *testing.T) {
	trace := iaosclient.IncorporationTrace{
		SchemaVersion: "1.0", CaseCode: "INC-LIVE-001", Verified: true,
		State: iaosclient.IncorporationState{
			CaseCode: "INC-LIVE-001", TenantID: "tenant-live", State: "registration_submitted",
			CommitmentMinor: 100000000, Currency: "CNY", ProposedName: "澄流热管理有限公司",
		},
		Journal:        []iaosclient.IncorporationJournalEntry{{Sequence: 1, CorrelationID: "corr-live", CreatedAt: "2026-07-27T10:00:00Z"}},
		WorldExchanges: []iaosclient.IncorporationWorldExchange{{Cursor: 42, MessageID: "intent-1", Kind: "intent", CorrelationID: "corr-live", RecordedAt: "2026-07-27T10:00:01Z"}},
	}
	items := []iaosclient.IncorporationWorkItem{
		{Sequence: 1, Capability: "incorporation.case.open", TaskType: "human_task", Participant: "founder-principal", Effective: "completed"},
		{Sequence: 2, Capability: "founder.resolution.prepare", TaskType: "agent_task", Participant: "incorporation-agent", Effective: "ready"},
	}
	projection, err := FromIAOS(trace, items)
	if err != nil {
		t.Fatal(err)
	}
	if projection.Cursor != 42 || projection.Chapter != "registration" || projection.Lifecycle.Progress != 50 {
		t.Fatalf("unexpected projection %#v", projection)
	}
	if len(projection.WorkItems) != 2 || projection.WorkItems[1].EvidenceRef == "" || projection.Brand.CompanyName != "澄流热管理有限公司" {
		t.Fatalf("projection lost governed facts %#v", projection)
	}
}

func TestFromIAOSRejectsUnverifiedTrace(t *testing.T) {
	_, err := FromIAOS(iaosclient.IncorporationTrace{SchemaVersion: "1.0", CaseCode: "INC-1"}, nil)
	if err == nil {
		t.Fatal("unverified trace accepted")
	}
}
