package genesisworkspace

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Provisioner interface {
	Provision(context.Context, Workspace) (string, error)
}

type Service struct {
	Store        *Store
	Provisioner  Provisioner
	ControlPlane *ControlPlaneClient
	Now          func() time.Time
}

func (s *Service) Create(ctx context.Context, request CreateRequest) (Result, error) {
	if err := request.Validate(); err != nil {
		return Result{}, err
	}
	if s.Store == nil || (s.Provisioner == nil && (s.ControlPlane == nil || !s.ControlPlane.enabled(ctx))) {
		return Result{}, fmt.Errorf("genesis workspace service is not configured")
	}
	request.OwnerPlayerID = strings.TrimSpace(request.OwnerPlayerID)
	request.DisplayName = strings.TrimSpace(request.DisplayName)
	request.IdempotencyKey = strings.TrimSpace(request.IdempotencyKey)
	request.TemplateKey = strings.TrimSpace(request.TemplateKey)
	request.Region = strings.TrimSpace(request.Region)
	request.Timezone = strings.TrimSpace(request.Timezone)
	request.RealismLevel = strings.TrimSpace(request.RealismLevel)
	if request.TemplateKey == "" {
		request.TemplateKey = "manufacturing-enterprise"
	}
	if request.Region == "" {
		request.Region = "CN-JS"
	}
	if request.Timezone == "" {
		request.Timezone = "Asia/Shanghai"
	}
	if request.RealismLevel == "" {
		request.RealismLevel = "standard"
	}
	if s.ControlPlane != nil && s.ControlPlane.enabled(ctx) {
		result, err := s.ControlPlane.Create(ctx, request.OwnerPlayerID, request)
		if err != nil {
			return Result{}, err
		}
		if err := s.Store.Save(request.OwnerPlayerID, request.IdempotencyKey, result.Workspace); err != nil {
			return Result{}, err
		}
		return result, nil
	}
	if existing, ok, err := s.Store.Find(request.OwnerPlayerID, request.IdempotencyKey); err != nil {
		return Result{}, err
	} else if ok {
		token, err := s.Provisioner.Provision(ctx, existing)
		if err != nil {
			existing.Status = StatusFailed
			existing.LastError = err.Error()
			existing.UpdatedAt = s.now()
			_ = s.Store.Save(request.OwnerPlayerID, request.IdempotencyKey, existing)
			return Result{}, err
		}
		existing.Status = StatusActive
		existing.CurrentStep = "identity_studio_ready"
		existing.LastError = ""
		existing.UpdatedAt = s.now()
		if err := s.Store.Save(request.OwnerPlayerID, request.IdempotencyKey, existing); err != nil {
			return Result{}, err
		}
		return Result{Workspace: existing, TenantToken: token}, nil
	}
	now := s.now()
	suffix, err := randomHex(8)
	if err != nil {
		return Result{}, err
	}
	workspace := Workspace{
		WorkspaceID: "gxw-" + suffix, OwnerPlayerID: request.OwnerPlayerID,
		DisplayName: request.DisplayName, TenantID: "tenant-gx-" + suffix,
		WorldRunID: "world-gx-" + suffix, CaseCode: "INC-GX-" + strings.ToUpper(suffix),
		TemplateKey: request.TemplateKey, Region: request.Region, Timezone: request.Timezone,
		RealismLevel: request.RealismLevel,
		Status:       StatusProvisioning, CurrentStep: "tenant_create",
		CreatedAt: now, UpdatedAt: now,
	}
	if err := s.Store.Save(request.OwnerPlayerID, request.IdempotencyKey, workspace); err != nil {
		return Result{}, err
	}
	token, err := s.Provisioner.Provision(ctx, workspace)
	if err != nil {
		workspace.Status = StatusFailed
		workspace.LastError = err.Error()
		workspace.UpdatedAt = now
		_ = s.Store.Save(request.OwnerPlayerID, request.IdempotencyKey, workspace)
		return Result{}, err
	}
	workspace.Status = StatusActive
	workspace.CurrentStep = "identity_studio_ready"
	workspace.LastError = ""
	workspace.UpdatedAt = now
	if err := s.Store.Save(request.OwnerPlayerID, request.IdempotencyKey, workspace); err != nil {
		return Result{}, err
	}
	return Result{Workspace: workspace, TenantToken: token}, nil
}

func (s *Service) now() time.Time {
	if s.Now != nil {
		return s.Now().UTC()
	}
	return time.Now().UTC()
}

func (s *Service) List(ctx context.Context, owner string) ([]Workspace, error) {
	if s.Store == nil {
		return nil, fmt.Errorf("genesis workspace store is not configured")
	}
	if s.ControlPlane != nil && s.ControlPlane.enabled(ctx) {
		return s.ControlPlane.List(ctx, strings.TrimSpace(owner))
	}
	items, err := s.Store.List(strings.TrimSpace(owner))
	if err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	return items, nil
}

func (s *Service) RefreshSession(ctx context.Context, owner, workspaceID string) (Result, error) {
	owner = strings.TrimSpace(owner)
	workspaceID = strings.TrimSpace(workspaceID)
	if owner == "" || workspaceID == "" {
		return Result{}, fmt.Errorf("owner and workspace_id are required")
	}
	if s.ControlPlane != nil && s.ControlPlane.enabled(ctx) {
		return s.ControlPlane.Session(ctx, owner, workspaceID)
	}
	items, err := s.List(ctx, owner)
	if err != nil {
		return Result{}, err
	}
	for _, workspace := range items {
		if workspace.WorkspaceID != workspaceID {
			continue
		}
		if workspace.Status != StatusActive {
			return Result{}, fmt.Errorf("workspace %s is not active", workspaceID)
		}
		token, err := s.Provisioner.Provision(ctx, workspace)
		if err != nil {
			return Result{}, err
		}
		return Result{Workspace: workspace, TenantToken: token}, nil
	}
	return Result{}, fmt.Errorf("workspace %s not found for owner", workspaceID)
}

func randomHex(bytes int) (string, error) {
	value := make([]byte, bytes)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}
