package incorporation

import "testing"

func TestCampaignDeterminismAndInvariants(t *testing.T) {
	tr := BuildTrace()
	if err := Validate(tr); err != nil {
		t.Fatal(err)
	}
	h := Hash(tr)
	for i := 0; i < 100; i++ {
		if got := Hash(BuildTrace()); got != h {
			t.Fatalf("run %d changed hash", i)
		}
	}
	if DecideRegistration(false) != "rejected" {
		t.Fatal("incomplete registration approved")
	}
	if CanSubmitBudget(true, false, true) {
		t.Fatal("unaccepted CEO submitted budget")
	}
	broken := BuildTrace()
	broken.Frames[len(broken.Frames)-1].Investor.Balance.Value = "1.00"
	if Validate(broken) == nil {
		t.Fatal("cash imbalance accepted")
	}
	broken = BuildTrace()
	broken.Frames[len(broken.Frames)-1].Governance.MandateActive = false
	if Validate(broken) == nil {
		t.Fatal("expired mandate accepted")
	}
}

func TestSnapshotRestoreAndUnifiedOperatorContract(t *testing.T) {
	snapshot,err:=Snapshot(BuildTrace());if err!=nil{t.Fatal(err)}
	restored,err:=Restore(snapshot);if err!=nil||Hash(restored)!=Hash(BuildTrace()){t.Fatalf("restore failed: %v",err)}
	if _,err=Restore([]byte(`{"broken":true}`));err==nil{t.Fatal("corrupt snapshot accepted")}
	frame:=BuildTrace().Frames[6];perms:=map[string]bool{"genesis.budget.submit":true}
	for _,mode:=range[]string{"human","agent"}{if err:=Authorize(Operator{"ACTOR-HCTM-CEO-01",mode,perms},"genesis.budget.submit",frame);err!=nil{t.Fatalf("%s rejected: %v",mode,err)}}
	frame.Governance.MandateActive=false;if Authorize(Operator{"ACTOR-HCTM-CEO-01","agent",perms},"genesis.budget.submit",frame)==nil{t.Fatal("expired mandate accepted")}
}
