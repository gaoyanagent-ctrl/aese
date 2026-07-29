package creative

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestCreativeJobStoreIsIdempotentAndFailClosedWhileRunning(t *testing.T) {
	store := NewJobStore(filepath.Join(t.TempDir(), "jobs.json"))
	job := CreativeJob{JobID: "creative-1", TenantID: "tenant-a", CaseCode: "INC-A", Status: "running"}
	if _, replay, err := store.Begin(job); err != nil || replay {
		t.Fatalf("first begin = replay %v err %v", replay, err)
	}
	if _, _, err := store.Begin(job); !errors.Is(err, ErrJobRunning) {
		t.Fatalf("concurrent begin error = %v", err)
	}
	job.Status = "completed"
	job.Parameters = map[string]any{"proposals": []any{map[string]any{"short_name": "测试"}}}
	if err := store.Save(job); err != nil {
		t.Fatal(err)
	}
	existing, replay, err := store.Begin(job)
	if err != nil || !replay || existing.Status != "completed" {
		t.Fatalf("completed replay = %#v, %v, %v", existing, replay, err)
	}
}

func TestCreativeJobStoreFiltersTenantAndCase(t *testing.T) {
	store := NewJobStore(filepath.Join(t.TempDir(), "jobs.json"))
	for _, job := range []CreativeJob{
		{JobID: "a1", TenantID: "tenant-a", CaseCode: "INC-1", Status: "completed"},
		{JobID: "a2", TenantID: "tenant-a", CaseCode: "INC-2", Status: "completed"},
		{JobID: "b1", TenantID: "tenant-b", CaseCode: "INC-1", Status: "completed"},
	} {
		if err := store.Save(job); err != nil {
			t.Fatal(err)
		}
	}
	items, err := store.List("tenant-a", "INC-1")
	if err != nil || len(items) != 1 || items[0].JobID != "a1" {
		t.Fatalf("unexpected filtered jobs: %#v %v", items, err)
	}
}
