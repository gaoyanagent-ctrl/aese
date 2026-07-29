package genesisworkspace

import (
	"fmt"
	"strings"
	"time"
)

type Status string

const (
	StatusProvisioning  Status = "provisioning"
	StatusAwaitingWorld Status = "awaiting_world"
	StatusActive        Status = "active"
	StatusFailed        Status = "failed"
)

type Workspace struct {
	WorkspaceID   string             `json:"workspace_id"`
	OwnerPlayerID string             `json:"owner_player_id"`
	DisplayName   string             `json:"display_name"`
	TenantID      string             `json:"tenant_id"`
	WorldRunID    string             `json:"world_run_id"`
	CaseCode      string             `json:"case_code"`
	TemplateKey   string             `json:"template_key,omitempty"`
	Region        string             `json:"region,omitempty"`
	Timezone      string             `json:"timezone,omitempty"`
	RealismLevel  string             `json:"realism_level,omitempty"`
	Status        Status             `json:"status"`
	CurrentStep   string             `json:"current_step"`
	CorrelationID string             `json:"correlation_id,omitempty"`
	Attempt       int                `json:"attempt,omitempty"`
	EvidenceRefs  map[string]string  `json:"evidence_refs,omitempty"`
	Steps         []ProvisioningStep `json:"steps,omitempty"`
	LastError     string             `json:"last_error,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

type ProvisioningStep struct {
	StepKey     string    `json:"step_key"`
	Status      string    `json:"status"`
	Attempt     int       `json:"attempt"`
	EvidenceRef string    `json:"evidence_ref,omitempty"`
	LastError   string    `json:"last_error,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateRequest struct {
	OwnerPlayerID          string `json:"owner_player_id"`
	DisplayName            string `json:"display_name"`
	IdempotencyKey         string `json:"idempotency_key"`
	TemplateKey            string `json:"template_key,omitempty"`
	Region                 string `json:"region,omitempty"`
	Timezone               string `json:"timezone,omitempty"`
	RealismLevel           string `json:"realism_level,omitempty"`
	DataRetentionConfirmed bool   `json:"data_retention_confirmed"`
}

func (r CreateRequest) Validate() error {
	if len(strings.TrimSpace(r.OwnerPlayerID)) < 8 {
		return fmt.Errorf("owner_player_id must contain at least 8 characters")
	}
	if len([]rune(strings.TrimSpace(r.DisplayName))) < 2 || len([]rune(strings.TrimSpace(r.DisplayName))) > 80 {
		return fmt.Errorf("display_name must contain 2 to 80 characters")
	}
	if len(strings.TrimSpace(r.IdempotencyKey)) < 12 || len(r.IdempotencyKey) > 160 {
		return fmt.Errorf("idempotency_key must contain 12 to 160 characters")
	}
	if r.TemplateKey != "" && r.TemplateKey != "manufacturing-enterprise" {
		return fmt.Errorf("unsupported template_key")
	}
	if r.Timezone != "" && r.Timezone != "Asia/Shanghai" {
		return fmt.Errorf("M9 manufacturing template requires Asia/Shanghai")
	}
	if r.RealismLevel != "" && r.RealismLevel != "standard" && r.RealismLevel != "strict" {
		return fmt.Errorf("realism_level must be standard or strict")
	}
	return nil
}

type Result struct {
	Workspace
	TenantToken string `json:"tenant_token,omitempty"`
}
