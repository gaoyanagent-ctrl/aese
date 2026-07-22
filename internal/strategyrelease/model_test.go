package strategyrelease

import "testing"

func TestAdoptedAndRollbackPaths(t *testing.T) {
	x := BuildTrace()
	if err := Validate(x); err != nil {
		t.Fatal(err)
	}
	h := x.TraceHash
	for i := 0; i < 100; i++ {
		y := BuildTrace()
		if y.TraceHash != h {
			t.Fatal("hash drift")
		}
	}
}
func TestFailsClosed(t *testing.T) {
	x := BuildTrace()
	x.Frames[3].BusinessWrites = 1
	if Validate(x) == nil {
		t.Fatal("shadow write accepted")
	}
	x = BuildTrace()
	x.RollbackFrames[2].Commitments[0].Status = "open"
	if Validate(x) == nil {
		t.Fatal("open commitment accepted")
	}
}
