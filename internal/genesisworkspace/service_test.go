package genesisworkspace

import (
	"context"
	"errors"
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
