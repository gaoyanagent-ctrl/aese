package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestReconcileCommandExitCodes(t *testing.T) {
	dir := t.TempDir()
	ok := filepath.Join(dir, "ok.json")
	if err := os.WriteFile(ok, []byte(`[
{"schema_version":"1.0","message_id":"i","kind":"intent","tenant_id":"t","world_pack_key":"p","world_pack_version":"1","world_run_id":"r","branch_id":"main","sim_occurred_at":"2026-07-23T00:00:00Z","correlation_id":"c","idempotency_key":"ii","producer":{"system":"iaos","component":"runtime","version":"1"},"subject_ref":{"namespace":"hctm","type":"case","code":"x"},"payload_type":"intent.v1","payload":{}},
{"schema_version":"1.0","message_id":"o","kind":"observation","tenant_id":"t","world_pack_key":"p","world_pack_version":"1","world_run_id":"r","branch_id":"main","sim_occurred_at":"2026-07-23T00:00:00Z","correlation_id":"c","idempotency_key":"oo","producer":{"system":"aese","component":"world","version":"1"},"subject_ref":{"namespace":"hctm","type":"case","code":"x"},"payload_type":"observation.v1","payload":{}},
{"schema_version":"1.0","message_id":"r","kind":"committed_outcome","tenant_id":"t","world_pack_key":"p","world_pack_version":"1","world_run_id":"r","branch_id":"main","sim_occurred_at":"2026-07-23T00:00:00Z","correlation_id":"c","idempotency_key":"rr","producer":{"system":"iaos","component":"runtime","version":"1"},"subject_ref":{"namespace":"hctm","type":"case","code":"x"},"payload_type":"outcome.v1","payload":{}}
]`), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := reconcileCommand([]string{ok}, &stdout, &stderr); code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr.String())
	}
	if !bytes.Contains(stdout.Bytes(), []byte(`"converged": true`)) {
		t.Fatalf("output=%s", stdout.String())
	}
}
