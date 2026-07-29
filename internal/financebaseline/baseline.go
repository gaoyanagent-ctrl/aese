package financebaseline

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Baseline struct {
	SchemaVersion string    `json:"schema_version"`
	Currency      string    `json:"currency"`
	AmountScale   int       `json:"amount_scale"`
	Timezone      string    `json:"timezone"`
	Objects       []Object  `json:"objects"`
	Controls      []Control `json:"controls"`
}
type Object struct {
	Code             string `json:"code"`
	Milestone        string `json:"milestone"`
	Authority        string `json:"authority"`
	CanonicalEntity  string `json:"canonical_entity"`
	HistoricalSource string `json:"historical_source"`
	MigrationStatus  string `json:"migration_status"`
}
type Control struct {
	Capability          string   `json:"capability"`
	ResponsiblePosition string   `json:"responsible_position"`
	Executor            string   `json:"executor"`
	Permission          string   `json:"permission"`
	ApprovalGate        string   `json:"approval_gate"`
	AmountLimitMinor    string   `json:"amount_limit_minor"`
	DataClassification  string   `json:"data_classification"`
	SoDConflicts        []string `json:"sod_conflicts"`
	Compensation        string   `json:"compensation"`
}

var requiredCapabilities = []string{"finance.organization.configure", "accounting.book.activate", "chart.of.accounts.activate", "capital.contribution.post", "finance.opening.readiness.evaluate"}

func Load(path string) (Baseline, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Baseline{}, err
	}
	var b Baseline
	if err = json.Unmarshal(raw, &b); err != nil {
		return Baseline{}, err
	}
	return b, Validate(b)
}
func Validate(b Baseline) error {
	if b.SchemaVersion != "1.0" || b.Currency != "CNY" || b.AmountScale != 2 || b.Timezone != "Asia/Shanghai" {
		return fmt.Errorf("invalid finance baseline header")
	}
	objects := map[string]bool{}
	for _, o := range b.Objects {
		if o.Code == "" || o.Milestone == "" || o.Authority != "IAOS" || o.CanonicalEntity == "" || o.HistoricalSource == "" {
			return fmt.Errorf("incomplete finance object %q", o.Code)
		}
		if objects[o.Code] {
			return fmt.Errorf("duplicate finance object %s", o.Code)
		}
		objects[o.Code] = true
		if o.MigrationStatus != "migrated" && o.MigrationStatus != "planned" {
			return fmt.Errorf("invalid migration status for %s", o.Code)
		}
	}
	controls := map[string]bool{}
	classes := map[string]bool{"internal": true, "confidential": true, "restricted": true}
	for _, c := range b.Controls {
		if c.Capability == "" || c.ResponsiblePosition == "" || c.Executor == "" || c.Permission == "" || c.ApprovalGate == "" || c.Compensation == "" {
			return fmt.Errorf("incomplete finance control %q", c.Capability)
		}
		if controls[c.Capability] {
			return fmt.Errorf("duplicate finance control %s", c.Capability)
		}
		controls[c.Capability] = true
		if !classes[c.DataClassification] {
			return fmt.Errorf("invalid data classification for %s", c.Capability)
		}
		amount, err := strconv.ParseInt(c.AmountLimitMinor, 10, 64)
		if err != nil || amount < 0 {
			return fmt.Errorf("invalid amount limit for %s", c.Capability)
		}
		if len(c.SoDConflicts) == 0 {
			return fmt.Errorf("missing SoD conflict for %s", c.Capability)
		}
	}
	for _, capability := range requiredCapabilities {
		if !controls[capability] {
			return fmt.Errorf("missing M9 finance control %s", capability)
		}
	}
	return nil
}
