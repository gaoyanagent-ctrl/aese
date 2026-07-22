package experiment

import "testing"

func TestEvidenceDeterministicIsolatedAndPaired(t *testing.T) {
	d := DefaultDefinition()
	e, err := BuildEvidence(d)
	if err != nil {
		t.Fatal(err)
	}
	if !e.StrategyEvidenceReady || len(e.Runs) != 60 {
		t.Fatalf("unexpected evidence %#v", e)
	}
	pairs := map[string]string{}
	ids := map[string]bool{}
	for _, r := range e.Runs {
		if r.ProductionWrites != 0 || r.ParentHash != d.ParentHash {
			t.Fatal("isolation violated")
		}
		if ids[r.RunID] {
			t.Fatal("duplicate run")
		}
		ids[r.RunID] = true
		if x := pairs[r.PairKey]; x != "" && x != r.DrawHash {
			t.Fatal("CRN violated")
		}
		pairs[r.PairKey] = r.DrawHash
	}
	first := e.EvidenceHash
	for i := 0; i < 100; i++ {
		h, err := ReplayHash(d)
		if err != nil || h != first {
			t.Fatalf("replay %d %s %v", i, h, err)
		}
	}
}

func TestNamedStreamsIndependent(t *testing.T) {
	seed := uint64(99)
	before := []uint64{draw(seed, "nominal", "demand", 1), draw(seed, "nominal", "supplier", 1)}
	_ = draw(seed, "nominal", "payment", 99)
	after := []uint64{draw(seed, "nominal", "demand", 1), draw(seed, "nominal", "supplier", 1)}
	if before[0] != after[0] || before[1] != after[1] {
		t.Fatal("named stream drift")
	}
}

func TestValidationFailsClosed(t *testing.T) {
	cases := []Definition{DefaultDefinition(), DefaultDefinition(), DefaultDefinition()}
	cases[0].HorizonWeeks = 0
	cases[1].Checkpoint = "unknown"
	cases[2].Policies[0].SafetyStock = -1
	for _, d := range cases {
		if Validate(d) == nil {
			t.Fatal("expected invalid")
		}
	}
}
