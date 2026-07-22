package capabilitybuild

import (
	"testing"
)

func TestTraceDeterministic(t *testing.T) {
	x := BuildTrace()
	if e := Validate(x); e != nil {
		t.Fatal(e)
	}
	h, _ := Hash(x)
	for i := 0; i < 100; i++ {
		got, _ := Hash(BuildTrace())
		if got != h {
			t.Fatal("nondeterministic")
		}
	}
	b, _ := Snapshot(x)
	restored, e := Restore(b)
	if e != nil {
		t.Fatal(e)
	}
	got, _ := Hash(restored)
	if got != h {
		t.Fatal("restore changed hash")
	}
}
func TestHardConstraintsAndGate(t *testing.T) {
	o := Options()
	if o[0].Feasible || !o[1].Feasible || o[2].Feasible {
		t.Fatal("acquisition constraints")
	}
	x := BuildTrace()
	x.Frames[9].Workers = nil
	if Validate(x) == nil {
		t.Fatal("accepted missing workforce")
	}
}
