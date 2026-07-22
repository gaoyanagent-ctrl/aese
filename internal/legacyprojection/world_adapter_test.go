package legacyprojection

import (
	"github.com/industrial-ai/iaos-aese/internal/scenariopack"
	"testing"
)

func TestToWorldEventsPreservesLegacy(t *testing.T) {
	pack, err := scenariopack.Load("../../scenario-packs/hctm")
	if err != nil {
		t.Fatal(err)
	}
	events := ToWorldEvents(pack.Stories[0])
	if len(events) != 22 {
		t.Fatalf("events=%d", len(events))
	}
	if events[0].EventID != "legacy-"+pack.Stories[0].Events.Events[0].EventID {
		t.Fatal("identity not preserved")
	}
}
