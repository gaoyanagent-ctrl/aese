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
	if len(baseline.Objects) != 9 || len(baseline.Controls) != 8 {
		t.Fatalf("unexpected inventory: %d objects %d controls", len(baseline.Objects), len(baseline.Controls))
	}
	if len(baseline.OrganizationFoundation.Organizations) != 6 ||
		len(baseline.OrganizationFoundation.DataSets) != 2 ||
		len(baseline.OrganizationFoundation.Assignments) != 9 ||
		len(baseline.OrganizationFoundation.AccessTemplates) != 3 {
		t.Fatalf("unexpected organization foundation: %#v", baseline.OrganizationFoundation)
	}
	if len(baseline.LedgerFoundation.Charts) != 1 ||
		len(baseline.LedgerFoundation.Accounts) != 2 ||
		len(baseline.LedgerFoundation.Calendars) != 1 ||
		len(baseline.LedgerFoundation.Calendars[0].Periods) != 12 ||
		len(baseline.LedgerFoundation.Books) != 1 ||
		len(baseline.LedgerFoundation.LedgerSets) != 1 {
		t.Fatalf("unexpected ledger foundation: %#v", baseline.LedgerFoundation)
	}
}
func TestFinanceBaselineRejectsPeriodGap(t *testing.T) {
	baseline, err := Load(filepath.Join("..", "..", "scenario-packs", "hctm", "finance-governance-baseline.json"))
	if err != nil {
		t.Fatal(err)
	}
	baseline.LedgerFoundation.Calendars[0].Periods[6].StartsOn = "2026-07-02"
	if err := Validate(baseline); err == nil {
		t.Fatal("fiscal calendar with a period gap was accepted")
	}
}
func TestFinanceBaselineRejectsUnknownBookLegalEntity(t *testing.T) {
	baseline, err := Load(filepath.Join("..", "..", "scenario-packs", "hctm", "finance-governance-baseline.json"))
	if err != nil {
		t.Fatal(err)
	}
	baseline.LedgerFoundation.Books[0].LegalEntityCode = "LE-MISSING"
	if err := Validate(baseline); err == nil {
		t.Fatal("accounting book with unknown legal entity was accepted")
	}
}
func TestFinanceBaselineRejectsCrossReferenceErrors(t *testing.T) {
	baseline, err := Load(filepath.Join("..", "..", "scenario-packs", "hctm", "finance-governance-baseline.json"))
	if err != nil {
		t.Fatal(err)
	}
	baseline.OrganizationFoundation.Assignments[0].SetCode = "MISSING"
	if err := Validate(baseline); err == nil {
		t.Fatal("assignment referencing missing Data Set was accepted")
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
