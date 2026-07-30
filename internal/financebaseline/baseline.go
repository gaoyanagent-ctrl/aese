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
	LedgerFoundation       LedgerFoundation       `json:"ledger_foundation"`
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
type LedgerFoundation struct {
	Charts     []ChartOfAccounts `json:"charts"`
	Accounts   []Account         `json:"accounts"`
	Calendars  []FiscalCalendar  `json:"calendars"`
	Books      []AccountingBook  `json:"books"`
	LedgerSets []LedgerSet       `json:"ledger_sets"`
}
type ChartOfAccounts struct {
	Code                   string `json:"chart_code"`
	Name                   string `json:"chart_name"`
	DataSetCode            string `json:"data_set_code"`
	Role                   string `json:"chart_role"`
	AccountingStandardCode string `json:"accounting_standard_code"`
	Status                 string `json:"status"`
	Version                int    `json:"version"`
}
type Account struct {
	ChartCode     string `json:"chart_code"`
	Code          string `json:"account_code"`
	Name          string `json:"account_name"`
	Class         string `json:"account_class"`
	NormalBalance string `json:"normal_balance"`
	Status        string `json:"status"`
}
type FiscalCalendar struct {
	Code       string         `json:"calendar_code"`
	Name       string         `json:"calendar_name"`
	Type       string         `json:"calendar_type"`
	FiscalYear int            `json:"fiscal_year"`
	Status     string         `json:"status"`
	Periods    []FiscalPeriod `json:"periods"`
}
type FiscalPeriod struct {
	Code     string `json:"period_code"`
	Number   int    `json:"period_number"`
	StartsOn string `json:"starts_on"`
	EndsOn   string `json:"ends_on"`
	Status   string `json:"status"`
}
type AccountingBook struct {
	Code                   string `json:"book_code"`
	Name                   string `json:"book_name"`
	LegalEntityCode        string `json:"legal_entity_code"`
	Role                   string `json:"book_role"`
	AccountingStandardCode string `json:"accounting_standard_code"`
	FunctionalCurrencyCode string `json:"functional_currency_code"`
	ChartCode              string `json:"chart_code"`
	CalendarCode           string `json:"calendar_code"`
	BalancingSegmentRule   string `json:"balancing_segment_rule"`
	Status                 string `json:"status"`
	EffectiveFrom          string `json:"effective_from"`
	Version                int    `json:"version"`
}
type LedgerSet struct {
	Code    string            `json:"ledger_set_code"`
	Name    string            `json:"ledger_set_name"`
	Status  string            `json:"status"`
	Members []LedgerSetMember `json:"members"`
}
type LedgerSetMember struct {
	BookCode string `json:"book_code"`
	Role     string `json:"member_role"`
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

var requiredCapabilities = []string{"finance.enterprise.structure.configure", "finance.reference.data.configure", "finance.ledger.foundation.configure", "finance.organization.configure", "accounting.book.activate", "chart.of.accounts.activate", "capital.contribution.post", "finance.opening.readiness.evaluate"}
var stableCode = regexp.MustCompile(`^[A-Z][A-Z0-9_-]{1,79}$`)
var accountCode = regexp.MustCompile(`^[A-Z0-9][A-Z0-9._-]{1,39}$`)

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
	if b.SchemaVersion != "1.2" || b.Currency != "CNY" || b.AmountScale != 2 || b.Timezone != "Asia/Shanghai" {
		return fmt.Errorf("invalid finance baseline header")
	}
	if err := validateOrganizationFoundation(b.OrganizationFoundation); err != nil {
		return err
	}
	if err := validateLedgerFoundation(b.LedgerFoundation, b.OrganizationFoundation); err != nil {
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

func validateLedgerFoundation(f LedgerFoundation, organizations OrganizationFoundation) error {
	orgs := map[string]Organization{}
	for _, organization := range organizations.Organizations {
		orgs[organization.Code] = organization
	}
	dataSets := map[string]ReferenceDataSet{}
	for _, set := range organizations.DataSets {
		dataSets[set.Code] = set
	}
	charts := map[string]ChartOfAccounts{}
	for _, chart := range f.Charts {
		set, setExists := dataSets[chart.DataSetCode]
		if !stableCode.MatchString(chart.Code) || strings.TrimSpace(chart.Name) == "" ||
			!setExists || !contains(set.DataTypes, "gl_account") || chart.Role != "operating" ||
			strings.TrimSpace(chart.AccountingStandardCode) == "" || chart.Status != "active" || chart.Version < 1 {
			return fmt.Errorf("invalid chart of accounts %q", chart.Code)
		}
		if _, exists := charts[chart.Code]; exists {
			return fmt.Errorf("duplicate chart of accounts %s", chart.Code)
		}
		charts[chart.Code] = chart
	}
	accountClasses := map[string]bool{"asset": true, "liability": true, "equity": true, "revenue": true, "expense": true}
	accountKeys := map[string]bool{}
	for _, account := range f.Accounts {
		if _, exists := charts[account.ChartCode]; !exists || !accountCode.MatchString(account.Code) ||
			strings.TrimSpace(account.Name) == "" || !accountClasses[account.Class] ||
			(account.NormalBalance != "debit" && account.NormalBalance != "credit") || account.Status != "active" {
			return fmt.Errorf("invalid account %s/%s", account.ChartCode, account.Code)
		}
		key := account.ChartCode + "/" + account.Code
		if accountKeys[key] {
			return fmt.Errorf("duplicate account %s", key)
		}
		accountKeys[key] = true
	}
	calendars := map[string]FiscalCalendar{}
	for _, calendar := range f.Calendars {
		if !stableCode.MatchString(calendar.Code) || strings.TrimSpace(calendar.Name) == "" ||
			calendar.Type != "natural_month" || calendar.FiscalYear < 2000 ||
			calendar.Status != "active" || len(calendar.Periods) != 12 {
			return fmt.Errorf("invalid fiscal calendar %q", calendar.Code)
		}
		expectedStart := time.Date(calendar.FiscalYear, 1, 1, 0, 0, 0, 0, time.UTC)
		for index, period := range calendar.Periods {
			start, startErr := time.Parse("2006-01-02", period.StartsOn)
			end, endErr := time.Parse("2006-01-02", period.EndsOn)
			if startErr != nil || endErr != nil || period.Number != index+1 ||
				period.Code != fmt.Sprintf("%04d-%02d", calendar.FiscalYear, index+1) ||
				!start.Equal(expectedStart) || !end.Equal(start.AddDate(0, 1, -1)) ||
				(period.Status != "open" && period.Status != "future" && period.Status != "closed") {
				return fmt.Errorf("invalid fiscal period %s/%s", calendar.Code, period.Code)
			}
			expectedStart = end.AddDate(0, 0, 1)
		}
		if _, exists := calendars[calendar.Code]; exists {
			return fmt.Errorf("duplicate fiscal calendar %s", calendar.Code)
		}
		calendars[calendar.Code] = calendar
	}
	roles := map[string]bool{"primary": true, "local_statutory": true, "reporting": true, "tax": true, "management": true}
	books := map[string]AccountingBook{}
	for _, book := range f.Books {
		legal, legalExists := orgs[book.LegalEntityCode]
		chart, chartExists := charts[book.ChartCode]
		if !stableCode.MatchString(book.Code) || strings.TrimSpace(book.Name) == "" ||
			!legalExists || legal.Type != "legal_entity" || !roles[book.Role] ||
			!chartExists || chart.AccountingStandardCode != book.AccountingStandardCode ||
			calendars[book.CalendarCode].Code == "" || book.FunctionalCurrencyCode != "CNY" ||
			book.BalancingSegmentRule != "legal_entity" || book.Status != "active" || book.Version < 1 {
			return fmt.Errorf("invalid accounting book %q", book.Code)
		}
		if _, err := time.Parse("2006-01-02", book.EffectiveFrom); err != nil {
			return fmt.Errorf("invalid accounting book effective date %s", book.Code)
		}
		if _, exists := books[book.Code]; exists {
			return fmt.Errorf("duplicate accounting book %s", book.Code)
		}
		books[book.Code] = book
	}
	sets := map[string]bool{}
	for _, set := range f.LedgerSets {
		if !stableCode.MatchString(set.Code) || strings.TrimSpace(set.Name) == "" ||
			set.Status != "active" || len(set.Members) == 0 || sets[set.Code] {
			return fmt.Errorf("invalid ledger set %q", set.Code)
		}
		sets[set.Code] = true
		memberBooks := map[string]bool{}
		for _, member := range set.Members {
			if _, exists := books[member.BookCode]; !exists ||
				(member.Role != "primary" && member.Role != "secondary") || memberBooks[member.BookCode] {
				return fmt.Errorf("invalid ledger set member %s/%s", set.Code, member.BookCode)
			}
			memberBooks[member.BookCode] = true
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
