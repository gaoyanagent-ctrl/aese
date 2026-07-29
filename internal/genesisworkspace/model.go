package genesisworkspace

import (
	"fmt"
	"strings"
	"time"
)

type Status string

const (
	StatusProvisioning Status = "provisioning"
	StatusActive       Status = "active"
	StatusFailed       Status = "failed"
)

type Workspace struct {
	WorkspaceID   string    `json:"workspace_id"`
	OwnerPlayerID string    `json:"owner_player_id"`
	DisplayName   string    `json:"display_name"`
	TenantID      string    `json:"tenant_id"`
	WorldRunID    string    `json:"world_run_id"`
	CaseCode      string    `json:"case_code"`
	Status        Status    `json:"status"`
	CurrentStep   string    `json:"current_step"`
	LastError     string    `json:"last_error,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateRequest struct {
	OwnerPlayerID  string `json:"owner_player_id"`
	DisplayName    string `json:"display_name"`
	IdempotencyKey string `json:"idempotency_key"`
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
	return nil
}

type Result struct {
	Workspace
	TenantToken string `json:"tenant_token,omitempty"`
}
