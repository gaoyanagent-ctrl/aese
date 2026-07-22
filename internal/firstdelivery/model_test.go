package firstdelivery

import "testing"

func TestDeterministicCommercialCycle(t *testing.T) {
	x := BuildTrace()
	if e := Validate(x); e != nil {
		t.Fatal(e)
	}
	h, _ := Hash(x)
	for i := 0; i < 100; i++ {
		g, _ := Hash(BuildTrace())
		if g != h {
			t.Fatal("nondeterministic")
		}
	}
	b, _ := Snapshot(x)
	y, e := Restore(b)
	if e != nil {
		t.Fatal(e)
	}
	g, _ := Hash(y)
	if g != h {
		t.Fatal("restore")
	}
}
func TestQuantityAndFinanceFailClosed(t *testing.T) {
	x := BuildTrace()
	x.Frames[12].Accepted = 11700
	if Validate(x) == nil {
		t.Fatal("short delivery closed")
	}
	x = BuildTrace()
	x.Frames[12].AR = m("300.00")
	if Validate(x) == nil {
		t.Fatal("unsettled AR closed")
	}
}
