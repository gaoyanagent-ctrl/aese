package assurance

import "testing"

func TestCycleAndThreeDispositions(t *testing.T) {
	c := BuildCycle()
	if err := Validate(c); err != nil {
		t.Fatal(err)
	}
	if c.Decision.Disposition != "renewed" || c.InjectedDecisions[0].Disposition != "reexperiment_required" || c.InjectedDecisions[1].Disposition != "retired" {
		t.Fatal("dispositions")
	}
	h := c.CycleHash
	for i := 0; i < 100; i++ {
		if BuildCycle().CycleHash != h {
			t.Fatal("hash drift")
		}
	}
}
func TestQualityAndLeakageFailClosed(t *testing.T) {
	c := BuildCycle()
	c.Missingness = 1
	if Validate(c) == nil {
		t.Fatal("missing accepted")
	}
	c = BuildCycle()
	c.Calibration.HoldoutLocked = false
	if Validate(c) == nil {
		t.Fatal("leak accepted")
	}
}
