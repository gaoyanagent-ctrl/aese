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

type IAOSClient struct {
	BaseURL       string
	PlatformToken string
	HTTPClient    *http.Client
}

func (c IAOSClient) Provision(ctx context.Context, workspace Workspace) (string, error) {
	base := strings.TrimRight(c.BaseURL, "/")
	if base == "" {
		return "", fmt.Errorf("IAOS base URL is required")
	}
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 90 * time.Second}
	}
	platformToken := strings.TrimSpace(c.PlatformToken)
	if platformToken == "" {
		parsed, err := url.Parse(base)
		if err != nil || (parsed.Hostname() != "127.0.0.1" && parsed.Hostname() != "localhost") {
			return "", fmt.Errorf("GENESIS_PLATFORM_TOKEN is required outside loopback development")
		}
		var tokenResponse struct {
			Token string `json:"token"`
		}
		if err := c.doJSON(ctx, client, http.MethodGet, base+"/api/v1/dev/token?tenant_id=tenant-001", "", nil, &tokenResponse); err != nil {
			return "", fmt.Errorf("obtain local platform token: %w", err)
		}
		platformToken = tokenResponse.Token
	}
	metadata := map[string]any{
		"source": "aese-enterprise-genesis", "workspace_id": workspace.WorkspaceID,
		"owner_player_id": workspace.OwnerPlayerID, "world_run_id": workspace.WorldRunID,
	}
	var tenantState struct {
		Status string `json:"status"`
	}
	tenantExists := c.doJSON(ctx, client, http.MethodGet, base+"/api/v1/platform/tenants/"+url.PathEscape(workspace.TenantID), platformToken, nil, &tenantState) == nil
	if !tenantExists {
		if err := c.doJSON(ctx, client, http.MethodPost, base+"/api/v1/platform/tenants", platformToken, map[string]any{
			"tenant_id": workspace.TenantID, "display_name": workspace.DisplayName,
			"region": "local", "metadata": metadata,
		}, &tenantState); err != nil {
			return "", fmt.Errorf("create IAOS tenant: %w", err)
		}
	}
	password := workspace.WorkspaceID + "-Founder!9"
	if err := c.doJSON(ctx, client, http.MethodPost, base+"/api/v1/platform-identities/founder/bootstrap", platformToken, map[string]any{
		"tenant_id": workspace.TenantID, "tenant_name": workspace.DisplayName,
		"password": password, "apply": true,
	}, nil); err != nil {
		return "", fmt.Errorf("bootstrap founder identity: %w", err)
	}
	var founderSession struct {
		Token string `json:"token"`
	}
	if err := c.doJSON(ctx, client, http.MethodPost, base+"/api/v1/auth/login", "", map[string]any{
		"tenant_id": workspace.TenantID,
		"username":  "founder-principal",
		"password":  password,
	}, &founderSession); err != nil {
		return "", fmt.Errorf("issue founder session: %w", err)
	}
	if strings.TrimSpace(founderSession.Token) == "" {
		return "", fmt.Errorf("issue founder session: IAOS returned an empty token")
	}
	if err := c.doJSON(ctx, client, http.MethodPost, base+"/api/v1/incorporations/runtime/install", founderSession.Token, map[string]any{
		"apply": true,
	}, nil); err != nil {
		return "", fmt.Errorf("install M9 runtime: %w", err)
	}
	if tenantState.Status != "active" {
		if err := c.doJSON(ctx, client, http.MethodPost, base+"/api/v1/platform/tenants/"+url.PathEscape(workspace.TenantID)+"/activate", platformToken, map[string]any{
			"reason": "Enterprise Genesis provisioning completed",
		}, nil); err != nil {
			return "", fmt.Errorf("activate IAOS tenant: %w", err)
		}
	}
	return founderSession.Token, nil
}

func (c IAOSClient) doJSON(ctx context.Context, client *http.Client, method, endpoint, token string, input, output any) error {
	var body io.Reader
	if input != nil {
		data, err := json.Marshal(input)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return err
	}
	if input != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s %s returned %d: %s", method, endpoint, resp.StatusCode, strings.TrimSpace(string(data)))
	}
	if output != nil && len(data) > 0 {
		if err := json.Unmarshal(data, output); err != nil {
			return err
		}
	}
	return nil
}
