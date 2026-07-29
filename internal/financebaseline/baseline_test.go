package financebaseline

import (
	"path/filepath"
	"testing"
)

func TestHCTMFinanceGovernanceBaseline(t *testing.T) {
	baseline, err := Load(filepath.Join("..", "..", "scenario-packs", "hctm", "finance-governance-baseline.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(baseline.Objects) != 9 || len(baseline.Controls) != 5 {
		t.Fatalf("unexpected inventory: %d objects %d controls", len(baseline.Objects), len(baseline.Controls))
	}
}
func TestFinanceBaselineRejectsMissingCompensation(t *testing.T) {
	baseline, err := Load(filepath.Join("..", "..", "scenario-packs", "hctm", "finance-governance-baseline.json"))
	if err != nil {
		t.Fatal(err)
	}
	baseline.Controls[0].Compensation = ""
	if err := Validate(baseline); err == nil {
		t.Fatal("control without compensation was accepted")
	}
}
