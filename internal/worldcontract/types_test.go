package worldcontract

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "world-contracts", "fixtures", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestStrictFixturesAndStableHashes(t *testing.T) {
	tests := []struct {
		name  string
		parse func([]byte) (any, error)
	}{
		{"world-run.json", func(b []byte) (any, error) { return ParseStrict[WorldRun](b) }},
		{"world-event.json", func(b []byte) (any, error) { return ParseStrict[WorldEvent](b) }},
		{"snapshot.json", func(b []byte) (any, error) { return ParseStrict[Snapshot](b) }},
		{"knowledge.json", func(b []byte) (any, error) { return ParseStrict[Knowledge](b) }},
		{"discrepancy.json", func(b []byte) (any, error) { return ParseStrict[Discrepancy](b) }},
		{"observation.json", func(b []byte) (any, error) { return ParseStrict[Observation](b) }},
		{"intent.json", func(b []byte) (any, error) { return ParseStrict[Intent](b) }},
		{"committed-outcome.json", func(b []byte) (any, error) { return ParseStrict[CommittedOutcome](b) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := tt.parse(fixture(t, tt.name))
			if err != nil {
				t.Fatal(err)
			}
			a, _ := CanonicalHash(v)
			b, _ := CanonicalHash(v)
			if a != b || !strings.HasPrefix(a, "sha256:") {
				t.Fatalf("unstable hash %q %q", a, b)
			}
		})
	}
}

func TestStrictParsingRejectsUnknownAndTrailing(t *testing.T) {
	if _, err := ParseStrict[WorldRun]([]byte(`{"schema_version":"1.0","unknown":true}`)); err == nil {
		t.Fatal("unknown field accepted")
	}
	if _, err := ParseStrict[WorldRun]([]byte(`{} {}`)); err == nil {
		t.Fatal("trailing document accepted")
	}
	if _, err := ParseStrict[WorldRun]([]byte(`{}`)); err == nil {
		t.Fatal("missing required fields accepted")
	}
}

func TestCanonicalHashIgnoresObjectKeyOrder(t *testing.T) {
	a, _ := CanonicalHash(map[string]any{"b": "2", "a": "1"})
	b, _ := CanonicalHash(map[string]any{"a": "1", "b": "2"})
	if a != b {
		t.Fatalf("hash differs: %s %s", a, b)
	}
}
