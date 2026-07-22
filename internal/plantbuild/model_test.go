package plantbuild

import "testing"

func TestPlantBuildDeterminismAndFailures(t *testing.T) {
	tr := BuildTrace()
	if e := Validate(tr); e != nil {
		t.Fatal(e)
	}
	h := Hash(tr)
	for i := 0; i < 100; i++ {
		if Hash(BuildTrace()) != h {
			t.Fatal("nondeterministic")
		}
	}
	r := Ranked()
	if r[0].SiteCode != "SITE-SZ-NORTH-LEASED-SHELL" || r[1].Feasible || r[2].Feasible {
		t.Fatalf("bad ranking %+v", r)
	}
	bad := BuildTrace()
	bad.Frames[9].Paid = money("16000000.00")
	if Validate(bad) == nil {
		t.Fatal("overpayment accepted")
	}
	snap, e := Snapshot(tr)
	if e != nil {
		t.Fatal(e)
	}
	restored, e := Restore(snap)
	if e != nil || Hash(restored) != h {
		t.Fatal("restore failed")
	}
	if _, e = Restore([]byte(`{"bad":1}`)); e == nil {
		t.Fatal("bad snapshot accepted")
	}
	for _, mode := range []string{"human", "agent"} {
		if e := Authorize(Operator{"ACTOR-HCTM-CFO-01", mode, map[string]bool{"genesis.payment.approve": true}}, "genesis.payment.approve", true, true); e != nil {
			t.Fatal(e)
		}
	}
	if Authorize(Operator{"A", "agent", map[string]bool{"genesis.payment.approve": true}}, "genesis.payment.approve", false, true) == nil {
		t.Fatal("unaccepted payment allowed")
	}
}
