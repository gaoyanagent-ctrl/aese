package scenariopack

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsEscapingPath(t *testing.T) {
	dir := t.TempDir()
	manifest := `{"schema_version":"1.0.0","pack_key":"test","pack_version":"1","display_name":"test","timezone":"UTC","tenant_template":"tenant-test","master_data":[{"path":"../secret.json"}],"stories":[]}`
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(manifest), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := Load(dir)
	if err == nil || !strings.Contains(err.Error(), "escapes pack root") {
		t.Fatalf("expected path escape error, got %v", err)
	}
}

func TestLoadRejectsUnsupportedVersion(t *testing.T) {
	dir := t.TempDir()
	manifest := `{"schema_version":"2.0.0","pack_key":"test","pack_version":"1","tenant_template":"tenant-test"}`
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(manifest), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := Load(dir)
	if err == nil || !strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("expected unsupported version, got %v", err)
	}
}
