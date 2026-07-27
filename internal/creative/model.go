// Package creative defines candidate-only AI creative contracts. Nothing in
// this package commits IAOS business facts.
package creative

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const SchemaVersion = "1.0"

type FounderIntent struct {
	SchemaVersion string   `json:"schema_version"`
	IntentID      string   `json:"intent_id"`
	TenantID      string   `json:"tenant_id"`
	CaseCode      string   `json:"case_code"`
	RawIdea       string   `json:"raw_idea"`
	Industry      string   `json:"industry"`
	Customers     []string `json:"customers"`
	Offerings     []string `json:"offerings"`
	BrandTraits   []string `json:"brand_traits"`
	CapitalMinor  string   `json:"capital_minor"`
	RiskAppetite  string   `json:"risk_appetite"`
	Assumptions   []string `json:"assumptions"`
	NeedsConfirm  []string `json:"needs_confirmation"`
	CreatedAt     string   `json:"created_at"`
}

type NamingProposal struct {
	ProposalID  string   `json:"proposal_id"`
	ChineseName string   `json:"chinese_name"`
	EnglishName string   `json:"english_name"`
	ShortName   string   `json:"short_name"`
	Rationale   string   `json:"rationale"`
	Slogan      string   `json:"slogan"`
	Keywords    []string `json:"keywords"`
	Primary     string   `json:"primary_color"`
	RiskHints   []string `json:"risk_hints"`
	Status      string   `json:"status"`
}

type BrandBrief struct {
	BriefID       string   `json:"brief_id"`
	IntentID      string   `json:"intent_id"`
	ProposalID    string   `json:"proposal_id"`
	CompanyName   string   `json:"company_name"`
	BrandPromise  string   `json:"brand_promise"`
	VisualMotifs  []string `json:"visual_motifs"`
	Avoid         []string `json:"avoid"`
	PrimaryColor  string   `json:"primary_color"`
	AssetVariants []string `json:"asset_variants"`
}

type CreativeJob struct {
	JobID        string         `json:"job_id"`
	TenantID     string         `json:"tenant_id"`
	CaseCode     string         `json:"case_code"`
	Kind         string         `json:"kind"`
	Status       string         `json:"status"`
	Provider     string         `json:"provider"`
	Model        string         `json:"model"`
	ModelVersion string         `json:"model_version"`
	Prompt       string         `json:"prompt"`
	Seed         string         `json:"seed"`
	Parameters   map[string]any `json:"parameters"`
	CreatedAt    string         `json:"created_at"`
	ContentHash  string         `json:"content_hash"`
}

type BrandAsset struct {
	AssetID       string `json:"asset_id"`
	TenantID      string `json:"tenant_id"`
	CaseCode      string `json:"case_code"`
	CompanyCode   string `json:"company_code"`
	JobID         string `json:"job_id"`
	AssetType     string `json:"asset_type"`
	MIME          string `json:"mime"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	URI           string `json:"uri"`
	ContentHash   string `json:"content_hash"`
	License       string `json:"license"`
	Status        string `json:"status"`
	SelectedBy    string `json:"selected_by,omitempty"`
	SelectionNote string `json:"selection_note,omitempty"`
	CreatedAt     string `json:"created_at"`
}

func (v FounderIntent) Validate() error {
	if v.SchemaVersion != SchemaVersion || v.IntentID == "" || v.TenantID == "" || v.CaseCode == "" {
		return fmt.Errorf("founder intent identity is incomplete")
	}
	if len(strings.TrimSpace(v.RawIdea)) < 12 || v.Industry == "" || len(v.BrandTraits) == 0 {
		return fmt.Errorf("founder idea, industry, and brand traits are required")
	}
	if v.RiskAppetite != "conservative" && v.RiskAppetite != "balanced" && v.RiskAppetite != "aggressive" {
		return fmt.Errorf("invalid risk appetite")
	}
	if _, err := time.Parse(time.RFC3339, v.CreatedAt); err != nil {
		return fmt.Errorf("created_at must be RFC3339: %w", err)
	}
	return nil
}

func Hash(value any) string {
	data, _ := json.Marshal(value)
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}
