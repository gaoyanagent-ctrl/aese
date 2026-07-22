package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestWorldCommandsDryRunApplyReplay(t *testing.T) {
	root := "../../world-contracts/runtime-example"
	var out, errOut bytes.Buffer
	if code := run([]string{"world", "validate", root}, &out, &errOut); code != 0 {
		t.Fatalf("validate=%d %s", code, errOut.String())
	}
	out.Reset()
	errOut.Reset()
	if code := run([]string{"world", "run", root}, &out, &errOut); code != 0 {
		t.Fatalf("dry=%d %s", code, errOut.String())
	}
	if !bytes.Contains(out.Bytes(), []byte(`"artifacts_written": false`)) {
		t.Fatalf("unexpected %s", out.String())
	}
	dir := t.TempDir()
	out.Reset()
	if code := run([]string{"world", "run", root, "--apply", "--output", dir}, &out, &errOut); code != 0 {
		t.Fatalf("apply=%d %s", code, errOut.String())
	}
	log := filepath.Join(dir, "event-log.json")
	if _, err := os.Stat(log); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if code := run([]string{"world", "replay", root, "--log", log}, &out, &errOut); code != 0 {
		t.Fatalf("replay=%d %s", code, errOut.String())
	}
}
func TestWorldApplyRequiresOutput(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"world", "run", "../../world-contracts/runtime-example", "--apply"}, &out, &errOut); code != 2 {
		t.Fatalf("code=%d", code)
	}
}
