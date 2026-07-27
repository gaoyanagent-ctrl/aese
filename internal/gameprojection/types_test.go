package gameprojection

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
)

func TestFixtureParsesStrictlyAndValidates(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "world-contracts", "fixtures", "game-projection.json"))
	if err != nil {
		t.Fatal(err)
	}
	projection, err := worldcontract.ParseStrict[Projection](data)
	if err != nil {
		t.Fatal(err)
	}
	if projection.Chapter != "founder_office" || len(projection.WorkItems) != 1 {
		t.Fatalf("unexpected projection: %#v", projection)
	}
}

func TestProjectionRejectsPresentationThatCannotBeTraced(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "world-contracts", "fixtures", "game-projection.json"))
	if err != nil {
		t.Fatal(err)
	}
	var projection Projection
	if err := json.Unmarshal(data, &projection); err != nil {
		t.Fatal(err)
	}
	projection.WorkItems[0].EvidenceRef = ""
	if err := projection.Validate(); err == nil {
		t.Fatal("expected untraceable work item to fail")
	}
}

func TestProjectionRejectsUnsupportedTimeScale(t *testing.T) {
	projection := Projection{
		SchemaVersion: SchemaVersion,
		ProjectionID:  "projection-1",
		TenantID:      "tenant-hctm-genesis",
		CaseCode:      "INC-1",
		WorldRunID:    "world-1",
		SimTime:       "2026-07-27T08:00:00+08:00",
		TimeScale:     3,
		Scene:         Scene{Mode: "2.5d"},
	}
	if err := projection.Validate(); err == nil {
		t.Fatal("expected unsupported time scale to fail")
	}
}
