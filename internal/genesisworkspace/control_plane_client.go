package genesisworkspace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type iaosTokenContextKey struct{}

func WithIAOSToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, iaosTokenContextKey{}, strings.TrimSpace(token))
}

func IAOSToken(ctx context.Context) string {
	value, _ := ctx.Value(iaosTokenContextKey{}).(string)
	return strings.TrimSpace(value)
}

type ControlPlaneClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

type ControlPlaneHTTPError struct {
	Path       string
	StatusCode int
	Body       string
}

func (e *ControlPlaneHTTPError) Error() string {
	return fmt.Sprintf("IAOS Genesis control plane %s returned %d: %s", e.Path, e.StatusCode, e.Body)
}

func (e *ControlPlaneHTTPError) UpstreamStatusCode() int { return e.StatusCode }

type controlPlaneWorkspace struct {
	WorkspaceID       string             `json:"workspace_id"`
	TenantID          string             `json:"tenant_id"`
	WorldRunID        string             `json:"world_run_id"`
	CaseCode          string             `json:"case_code"`
	DisplayName       string             `json:"display_name"`
	Status            Status             `json:"status"`
	CurrentCheckpoint string             `json:"current_checkpoint"`
	CorrelationID     string             `json:"correlation_id"`
	Attempt           int                `json:"attempt"`
	LastError         string             `json:"last_error"`
	EvidenceRefs      map[string]string  `json:"evidence_refs"`
	Steps             []ProvisioningStep `json:"steps"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

func (c ControlPlaneClient) enabled(ctx context.Context) bool {
	return strings.TrimSpace(c.BaseURL) != "" && IAOSToken(ctx) != ""
}

func (c ControlPlaneClient) Create(ctx context.Context, owner string, request CreateRequest) (Result, error) {
	var workspace controlPlaneWorkspace
	err := c.doJSON(ctx, http.MethodPost, "/api/v1/genesis/workspaces", map[string]any{
		"display_name": request.DisplayName, "idempotency_key": request.IdempotencyKey,
		"template_key": request.TemplateKey, "template_version": "m9/1.8.0",
		"region": request.Region, "timezone": request.Timezone,
	}, &workspace)
	if err != nil {
		return Result{}, err
	}
	if workspace.Status == StatusAwaitingWorld || workspace.CurrentCheckpoint == "runtime_installed" {
		evidenceRef := "aese://genesis-workspaces/" + workspace.WorkspaceID + "/world-runs/" + workspace.WorldRunID
		if err = c.doJSON(ctx, http.MethodPost, "/api/v1/genesis/workspaces/"+url.PathEscape(workspace.WorkspaceID)+"/world-ready", map[string]string{
			"world_run_id": workspace.WorldRunID, "evidence_ref": evidenceRef,
		}, &workspace); err != nil {
			return Result{}, err
		}
	}
	return c.SessionFromWorkspace(ctx, owner, workspace)
}

func (c ControlPlaneClient) List(ctx context.Context, owner string) ([]Workspace, error) {
	var response struct {
		Items []controlPlaneWorkspace `json:"items"`
	}
	if err := c.doJSON(ctx, http.MethodGet, "/api/v1/genesis/workspaces", nil, &response); err != nil {
		return nil, err
	}
	out := make([]Workspace, 0, len(response.Items))
	for _, item := range response.Items {
		out = append(out, mapControlPlaneWorkspace(owner, item))
	}
	return out, nil
}

func (c ControlPlaneClient) Session(ctx context.Context, owner, workspaceID string) (Result, error) {
	var response struct {
		Workspace controlPlaneWorkspace `json:"workspace"`
		Token     string                `json:"token"`
	}
	if err := c.doJSON(ctx, http.MethodPost, "/api/v1/genesis/workspaces/"+url.PathEscape(workspaceID)+"/session", map[string]any{}, &response); err != nil {
		return Result{}, err
	}
	return Result{Workspace: mapControlPlaneWorkspace(owner, response.Workspace), TenantToken: response.Token}, nil
}

func (c ControlPlaneClient) AdoptLegacy(ctx context.Context, owner string, workspace Workspace) (Result, error) {
	var adopted controlPlaneWorkspace
	err := c.doJSON(ctx, http.MethodPost, "/api/v1/genesis/workspaces/legacy-adoptions", map[string]any{
		"workspace_id": workspace.WorkspaceID,
		"world_run_id": workspace.WorldRunID,
		"case_code":    workspace.CaseCode,
		"display_name": workspace.DisplayName,
		"evidence_ref": "aese-local-store:" + workspace.WorkspaceID,
		"template_key": workspace.TemplateKey,
		"region":       workspace.Region,
		"timezone":     workspace.Timezone,
	}, &adopted)
	if err != nil {
		return Result{}, err
	}
	return c.SessionFromWorkspace(ctx, owner, adopted)
}

func (c ControlPlaneClient) SessionFromWorkspace(ctx context.Context, owner string, workspace controlPlaneWorkspace) (Result, error) {
	return c.Session(ctx, owner, workspace.WorkspaceID)
}

func mapControlPlaneWorkspace(owner string, item controlPlaneWorkspace) Workspace {
	return Workspace{
		WorkspaceID: item.WorkspaceID, OwnerPlayerID: owner, DisplayName: item.DisplayName,
		TenantID: item.TenantID, WorldRunID: item.WorldRunID, CaseCode: item.CaseCode,
		Status: item.Status, CurrentStep: item.CurrentCheckpoint, CorrelationID: item.CorrelationID,
		Attempt: item.Attempt, LastError: item.LastError, EvidenceRefs: item.EvidenceRefs, Steps: item.Steps,
		CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
	}
}

func (c ControlPlaneClient) doJSON(ctx context.Context, method, path string, input, output any) error {
	var body io.Reader
	if input != nil {
		raw, err := json.Marshal(input)
		if err != nil {
			return err
		}
		body = bytes.NewReader(raw)
	}
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 90 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(c.BaseURL, "/")+path, body)
	if err != nil {
		return err
	}
	if input != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+IAOSToken(ctx))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &ControlPlaneHTTPError{
			Path: path, StatusCode: resp.StatusCode, Body: strings.TrimSpace(string(raw)),
		}
	}
	if output != nil && len(raw) > 0 {
		return json.Unmarshal(raw, output)
	}
	return nil
}
