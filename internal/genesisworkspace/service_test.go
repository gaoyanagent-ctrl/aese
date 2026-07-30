package genesisworkspace

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"
)

type fakeProvisioner struct {
	calls int
	fail  int
}

func (f *fakeProvisioner) Provision(_ context.Context, workspace Workspace) (string, error) {
	f.calls++
	if f.calls <= f.fail {
		return "", errors.New("injected provisioning failure")
	}
	return "token-for-" + workspace.TenantID, nil
}

func TestServiceCreatesIndependentTenantAndReplaysIdempotently(t *testing.T) {
	provisioner := &fakeProvisioner{}
	service := &Service{
		Store:       NewStore(filepath.Join(t.TempDir(), "workspaces.json")),
		Provisioner: provisioner,
		Now:         func() time.Time { return time.Date(2026, 7, 28, 3, 0, 0, 0, time.UTC) },
	}
	request := CreateRequest{
		OwnerPlayerID: "player-local-001", DisplayName: "新能源汽车热管理创业项目",
		IdempotencyKey: "create-company-001",
	}
	first, err := service.Create(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Create(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if first.TenantID == "" || first.TenantID == "tenant-hctm-genesis" {
		t.Fatalf("tenant was not independently allocated: %#v", first)
	}
	if first.TenantID != second.TenantID || first.WorkspaceID != second.WorkspaceID {
		t.Fatalf("idempotent replay allocated a second workspace: %#v %#v", first, second)
	}
	if first.Status != StatusActive || first.CurrentStep != "identity_studio_ready" {
		t.Fatalf("workspace is not ready: %#v", first)
	}
}

func TestTwoCreateKeysAllocateDifferentTenants(t *testing.T) {
	service := &Service{
		Store:       NewStore(filepath.Join(t.TempDir(), "workspaces.json")),
		Provisioner: &fakeProvisioner{},
	}
	one, err := service.Create(context.Background(), CreateRequest{
		OwnerPlayerID: "player-local-001", DisplayName: "企业一", IdempotencyKey: "create-company-001",
	})
	if err != nil {
		t.Fatal(err)
	}
	two, err := service.Create(context.Background(), CreateRequest{
		OwnerPlayerID: "player-local-001", DisplayName: "企业二", IdempotencyKey: "create-company-002",
	})
	if err != nil {
		t.Fatal(err)
	}
	if one.TenantID == two.TenantID {
		t.Fatal("different workspaces share a tenant")
	}
}

func TestRefreshSessionReissuesTokenOnlyForWorkspaceOwner(t *testing.T) {
	provisioner := &fakeProvisioner{}
	service := &Service{
		Store:       NewStore(filepath.Join(t.TempDir(), "workspaces.json")),
		Provisioner: provisioner,
	}
	created, err := service.Create(context.Background(), CreateRequest{
		OwnerPlayerID: "player-local-001", DisplayName: "企业一", IdempotencyKey: "create-company-001",
	})
	if err != nil {
		t.Fatal(err)
	}
	refreshed, err := service.RefreshSession(context.Background(), "player-local-001", created.WorkspaceID)
	if err != nil {
		t.Fatal(err)
	}
	if refreshed.TenantToken != "token-for-"+created.TenantID {
		t.Fatalf("unexpected refreshed token: %q", refreshed.TenantToken)
	}
	if _, err := service.RefreshSession(context.Background(), "player-local-other", created.WorkspaceID); err == nil {
		t.Fatal("another player refreshed the workspace session")
	}
}

func TestRefreshSessionAdoptsOwnedLegacyWorkspaceAfterControlPlane404(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(filepath.Join(dir, "workspaces.json"))
	legacy := Workspace{
		WorkspaceID: "gxw-legacy1234", OwnerPlayerID: "player-local-001",
		DisplayName: "旧企业", TenantID: "tenant-legacy",
		WorldRunID: "world-legacy1234", CaseCode: "INC-LEGACY-1",
		TemplateKey: "manufacturing-enterprise", Region: "CN-JS",
		Timezone: "Asia/Shanghai", Status: StatusActive,
	}
	if err := store.Save(legacy.OwnerPlayerID, "legacy-workspace-001", legacy); err != nil {
		t.Fatal(err)
	}
	var adopted bool
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/genesis/workspaces/gxw-legacy1234/session":
			if !adopted {
				http.NotFound(w, r)
				return
			}
			fmt.Fprint(w, `{"workspace":{"workspace_id":"gxw-legacy1234","tenant_id":"tenant-legacy","world_run_id":"world-legacy1234","case_code":"INC-LEGACY-1","display_name":"旧企业","status":"active","current_checkpoint":"session_issued"},"token":"tenant-token"}`)
		case "/api/v1/genesis/workspaces/legacy-adoptions":
			adopted = true
			fmt.Fprint(w, `{"workspace_id":"gxw-legacy1234","tenant_id":"tenant-legacy","world_run_id":"world-legacy1234","case_code":"INC-LEGACY-1","display_name":"旧企业","status":"active","current_checkpoint":"tenant_active"}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer upstream.Close()
	control := &ControlPlaneClient{BaseURL: upstream.URL}
	service := &Service{Store: store, ControlPlane: control}
	ctx := WithIAOSToken(context.Background(), "player-token")
	result, err := service.RefreshSession(ctx, legacy.OwnerPlayerID, legacy.WorkspaceID)
	if err != nil {
		t.Fatal(err)
	}
	if !adopted || result.TenantToken != "tenant-token" {
		t.Fatalf("legacy workspace was not adopted: %#v", result)
	}
}

func TestProvisioningFailureRetriesSameWorkspaceAndTenant(t *testing.T) {
	provisioner := &fakeProvisioner{fail: 1}
	service := &Service{
		Store:       NewStore(filepath.Join(t.TempDir(), "workspaces.json")),
		Provisioner: provisioner,
	}
	request := CreateRequest{
		OwnerPlayerID: "player-recovery-001", DisplayName: "恢复测试企业",
		IdempotencyKey: "create-company-recovery-001",
	}
	if _, err := service.Create(context.Background(), request); err == nil {
		t.Fatal("injected failure was not returned")
	}
	failed, ok, err := service.Store.Find(request.OwnerPlayerID, request.IdempotencyKey)
	if err != nil || !ok || failed.Status != StatusFailed {
		t.Fatalf("failed checkpoint was not persisted: %#v %v", failed, err)
	}
	recovered, err := service.Create(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if recovered.WorkspaceID != failed.WorkspaceID || recovered.TenantID != failed.TenantID {
		t.Fatalf("retry allocated a new identity: failed=%#v recovered=%#v", failed, recovered)
	}
	if recovered.Status != StatusActive || recovered.LastError != "" {
		t.Fatalf("workspace did not recover: %#v", recovered)
	}
}
