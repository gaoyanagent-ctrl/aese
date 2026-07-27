package gameprojection

import (
	"testing"

	"github.com/industrial-ai/iaos-aese/internal/incorporation"
)

func TestFromIncorporationTraceProjectsEveryFrame(t *testing.T) {
	trace := incorporation.BuildTrace()
	for i := range trace.Frames {
		projection, err := FromIncorporationTrace(trace, "INC-DEMO-001", i)
		if err != nil {
			t.Fatalf("frame %d: %v", i, err)
		}
		if projection.Cursor != trace.Frames[i].IAOSCursor || len(projection.EvidenceRefs) == 0 {
			t.Fatalf("frame %d lost source evidence", i)
		}
	}
}

func TestFromIncorporationTraceRejectsUnknownFrame(t *testing.T) {
	if _, err := FromIncorporationTrace(incorporation.BuildTrace(), "INC-DEMO-001", 99); err == nil {
		t.Fatal("expected frame bounds failure")
	}
}
