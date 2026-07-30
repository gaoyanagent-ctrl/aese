package financebaseline

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Baseline struct {
	SchemaVersion          string                 `json:"schema_version"`
	Currency               string                 `json:"currency"`
	AmountScale            int                    `json:"amount_scale"`
	Timezone               string                 `json:"timezone"`
	OrganizationFoundation OrganizationFoundation `json:"organization_foundation"`
	Objects                []Object               `json:"objects"`
	Controls               []Control              `json:"controls"`
}
type OrganizationFoundation struct {
	Organizations   []Organization      `json:"organizations"`
	DataSets        []ReferenceDataSet  `json:"data_sets"`
	Assignments     []DataSetAssignment `json:"assignments"`
	AccessTemplates []AccessTemplate    `json:"access_templates"`
}
type Organization struct {
	Code            string `json:"organization_code"`
	Name            string `json:"organization_name"`
	Type            string `json:"organization_type"`
	ParentCode      string `json:"parent_code"`
	LegalEntityCode string `json:"legal_entity_code"`
	Status          string `json:"status"`
	EffectiveFrom   string `json:"effective_from"`
	SourceRef       string `json:"source_ref"`
}
type ReferenceDataSet struct {
	Code          string   `json:"set_code"`
	Name          string   `json:"set_name"`
	ScopeType     string   `json:"scope_type"`
	OwnerCode     string   `json:"owner_organization_code"`
	DataTypes     []string `json:"data_types"`
	Status        string   `json:"status"`
	EffectiveFrom string   `json:"effective_from"`
}
type DataSetAssignment struct {
	SetCode          string `json:"set_code"`
	DataType         string `json:"data_type"`
	OrganizationCode string `json:"determinant_organization_code"`
	AccessMode       string `json:"access_mode"`
	Priority         int    `json:"priority"`
	Status           string `json:"status"`
	EffectiveFrom    string `json:"effective_from"`
}
type AccessTemplate struct {
	SubjectRef       string   `json:"subject_ref"`
	OrganizationCode string   `json:"organization_code"`
	AccessScope      string   `json:"access_scope"`
	Permissions      []string `json:"permissions"`
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

var requiredCapabilities = []string{"finance.enterprise.structure.configure", "finance.reference.data.configure", "finance.organization.configure", "accounting.book.activate", "chart.of.accounts.activate", "capital.contribution.post", "finance.opening.readiness.evaluate"}
var stableCode = regexp.MustCompile(`^[A-Z][A-Z0-9_-]{1,79}$`)

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
	if b.SchemaVersion != "1.1" || b.Currency != "CNY" || b.AmountScale != 2 || b.Timezone != "Asia/Shanghai" {
		return fmt.Errorf("invalid finance baseline header")
	}
	if err := validateOrganizationFoundation(b.OrganizationFoundation); err != nil {
		return err
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

func validateOrganizationFoundation(f OrganizationFoundation) error {
	organizationTypes := map[string]bool{
		"enterprise_group": true, "legal_entity": true, "business_unit": true,
		"site": true, "plant": true, "finance_organization": true,
		"shared_service_center": true, "management_accounting_area": true,
	}
	dataTypes := map[string]bool{
		"gl_account": true, "business_partner": true, "customer": true, "supplier": true,
		"product": true, "payment_term": true, "currency": true, "exchange_rate": true,
		"unit_of_measure": true,
	}
	organizations := map[string]Organization{}
	for _, organization := range f.Organizations {
		if !stableCode.MatchString(organization.Code) || strings.TrimSpace(organization.Name) == "" ||
			!organizationTypes[organization.Type] || organization.Status != "active" ||
			strings.TrimSpace(organization.SourceRef) == "" {
			return fmt.Errorf("invalid organization %q", organization.Code)
		}
		if _, err := time.Parse("2006-01-02", organization.EffectiveFrom); err != nil {
			return fmt.Errorf("invalid organization effective date %s", organization.Code)
		}
		if _, exists := organizations[organization.Code]; exists {
			return fmt.Errorf("duplicate organization %s", organization.Code)
		}
		organizations[organization.Code] = organization
	}
	for _, organization := range f.Organizations {
		if organization.Type == "enterprise_group" {
			if organization.ParentCode != "" {
				return fmt.Errorf("enterprise group %s cannot have parent", organization.Code)
			}
			continue
		}
		if _, exists := organizations[organization.ParentCode]; !exists {
			return fmt.Errorf("organization %s has unknown parent %s", organization.Code, organization.ParentCode)
		}
		if organization.LegalEntityCode != "" {
			legal, exists := organizations[organization.LegalEntityCode]
			if !exists || legal.Type != "legal_entity" {
				return fmt.Errorf("organization %s has invalid legal entity %s", organization.Code, organization.LegalEntityCode)
			}
		}
	}
	sets := map[string]ReferenceDataSet{}
	for _, set := range f.DataSets {
		if !stableCode.MatchString(set.Code) || set.Name == "" || len(set.DataTypes) == 0 ||
			(set.ScopeType != "common" && set.ScopeType != "group" && set.ScopeType != "region" &&
				set.ScopeType != "business_unit" && set.ScopeType != "organization") {
			return fmt.Errorf("invalid reference data set %q", set.Code)
		}
		if _, exists := organizations[set.OwnerCode]; !exists {
			return fmt.Errorf("data set %s has unknown owner %s", set.Code, set.OwnerCode)
		}
		for _, dataType := range set.DataTypes {
			if !dataTypes[dataType] {
				return fmt.Errorf("data set %s has unsupported data type %s", set.Code, dataType)
			}
		}
		if _, exists := sets[set.Code]; exists {
			return fmt.Errorf("duplicate reference data set %s", set.Code)
		}
		sets[set.Code] = set
	}
	for _, assignment := range f.Assignments {
		set, exists := sets[assignment.SetCode]
		if !exists {
			return fmt.Errorf("assignment has unknown data set %s", assignment.SetCode)
		}
		if _, exists := organizations[assignment.OrganizationCode]; !exists {
			return fmt.Errorf("assignment has unknown organization %s", assignment.OrganizationCode)
		}
		if !contains(set.DataTypes, assignment.DataType) ||
			(assignment.AccessMode != "read" && assignment.AccessMode != "reference" && assignment.AccessMode != "extend") ||
			assignment.Priority <= 0 {
			return fmt.Errorf("invalid data set assignment %s/%s", assignment.SetCode, assignment.DataType)
		}
	}
	allowedPermissions := map[string]bool{"read": true, "reference": true, "extend": true, "manage": true}
	for _, access := range f.AccessTemplates {
		if !strings.HasPrefix(access.SubjectRef, "position:") || access.AccessScope != "subtree" && access.AccessScope != "own" {
			return fmt.Errorf("invalid access template %s", access.SubjectRef)
		}
		if _, exists := organizations[access.OrganizationCode]; !exists || len(access.Permissions) == 0 {
			return fmt.Errorf("invalid access template organization %s", access.OrganizationCode)
		}
		for _, permission := range access.Permissions {
			if !allowedPermissions[permission] {
				return fmt.Errorf("invalid access permission %s", permission)
			}
		}
	}
	return nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
