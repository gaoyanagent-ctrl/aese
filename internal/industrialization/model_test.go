package industrialization

import "testing"

func TestDeterminismAndRestore(t *testing.T) {
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
func TestPPAPGateFailsClosed(t *testing.T) {
	x := BuildTrace()
	x.Frames[10].PPAPStatus = "submitted"
	if Validate(x) == nil {
		t.Fatal("accepted without customer PPAP")
	}
}
