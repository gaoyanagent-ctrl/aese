package worldcontract

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestM9FrozenContractMatchesWorldEnvelopeAndRejectsBrokenFixture(t *testing.T) {
	root := filepath.Join("..", "..")
	raw, err := os.ReadFile(filepath.Join(root, "docs", "contracts", "m9-native-incorporation-contract.json"))
	if err != nil {
		t.Fatal(err)
	}
	var lock struct {
		SchemaVersion string   `json:"schema_version"`
		TenantID      string   `json:"tenant_id"`
		States        []string `json:"states"`
		TerminalState string   `json:"terminal_state"`
	}
	if err := json.Unmarshal(raw, &lock); err != nil {
		t.Fatal(err)
	}
	if lock.SchemaVersion != SchemaVersion || lock.TenantID != "tenant-hctm-genesis" || len(lock.States) != 13 || lock.TerminalState != "enterprise_operational_ready" {
		t.Fatalf("incompatible M9 lock: %+v", lock)
	}
	broken, err := os.ReadFile(filepath.Join(root, "world-contracts", "broken-fixtures", "m9-observation-old-schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseStrict[Observation](broken); err == nil {
		t.Fatal("old schema fixture must fail closed")
	}
}
