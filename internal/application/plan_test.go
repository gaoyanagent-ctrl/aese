package application

import (
	"path/filepath"
	"testing"

	"github.com/industrial-ai/iaos-aese/internal/scenariopack"
)

func TestCompilePlanHasSevenActsAndDeterministicHash(t *testing.T) {
	pack, err := loadPack("..", "..", "scenario-packs", "hctm")
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}
	planned, err := CompilePlan(pack, "order-expedite-01")
	if err != nil {
		t.Fatalf("compile plan: %v", err)
	}
	if planned.ActCount != 7 {
		t.Fatalf("act count = %d, want 7", planned.ActCount)
	}
	if planned.PlanHash == "" {
		t.Fatal("plan hash is required")
	}
	if len(planned.Stages) != 12 { // preflight+initialize+act1..7+analyze+verify+reset
		t.Fatalf("plan stages count = %d, want 12", len(planned.Stages))
	}

	repeat, err := CompilePlan(pack, "order-expedite-01")
	if err != nil {
		t.Fatalf("compile plan repeat: %v", err)
	}
	if planned.PlanHash != repeat.PlanHash {
		t.Fatalf("plan hash mismatch on repeat: %s != %s", planned.PlanHash, repeat.PlanHash)
	}
}

func TestCompilePlanRejectsUnknownStory(t *testing.T) {
	pack, err := loadPack("..", "..", "scenario-packs", "hctm")
	if err != nil {
		t.Fatalf("load pack: %v", err)
	}
	_, err = CompilePlan(pack, "unknown")
	if err == nil {
		t.Fatal("expected unknown story to fail")
	}
}

func TestStateTransitionProgression(t *testing.T) {
	next, err := NextStatus(RunStatusPlanned, RunActionPreflight, RunTransitionContext{CurrentAct: 0, TotalActs: 7})
	if err != nil || next != RunStatusInitializing {
		t.Fatalf("planned -> preflight = %v %v", next, err)
	}
	next, err = NextStatus(next, RunActionInitialize, RunTransitionContext{CurrentAct: 0, TotalActs: 7})
	if err != nil || next != RunStatusReady {
		t.Fatalf("initializing -> initialize = %v %v", next, err)
	}
	next, err = NextStatus(next, RunActionAdvance, RunTransitionContext{CurrentAct: 6, TotalActs: 7})
	if err != nil || next != RunStatusRunning {
		t.Fatalf("ready -> running = %v %v", next, err)
	}
	next, err = NextStatus(next, RunActionAdvance, RunTransitionContext{CurrentAct: 7, TotalActs: 7})
	if err != nil || next != RunStatusAwaitingAnalysis {
		t.Fatalf("running -> awaiting_analysis = %v %v", next, err)
	}
}

func loadPack(parts ...string) (*scenariopack.Pack, error) {
	path := filepath.Join(parts...)
	return scenariopack.Load(path)
}
