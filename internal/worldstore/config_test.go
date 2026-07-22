package worldstore

import "testing"

func TestFromEnvRequiresDedicatedBoundary(t *testing.T) {
	t.Setenv(DatabaseURLEnv, "postgres://aese_world_app:dev@127.0.0.1:55432/aese_world?sslmode=disable")
	if _, err := FromEnv(); err != nil {
		t.Fatal(err)
	}
	t.Setenv(DatabaseURLEnv, "postgres://iaos_app:dev@127.0.0.1:5433/iaos")
	if _, err := FromEnv(); err == nil {
		t.Fatal("IAOS database accepted")
	}
}
