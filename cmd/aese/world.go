package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/world"
	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
)

func worldCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "world requires validate, inspect, run, or replay")
		return 2
	}
	switch args[0] {
	case "validate":
		return worldValidate(args[1:], stdout, stderr)
	case "inspect":
		return worldInspect(args[1:], stdout, stderr)
	case "run":
		return worldRun(args[1:], stdout, stderr, false)
	case "replay":
		return worldRun(args[1:], stdout, stderr, true)
	default:
		fmt.Fprintf(stderr, "unknown world command %q\n", args[0])
		return 2
	}
}
func oneWorldDir(name string, args []string, stderr io.Writer) (world.Bundle, int) {
	if len(args) != 1 {
		fmt.Fprintf(stderr, "world %s requires exactly one world directory\n", name)
		return world.Bundle{}, 2
	}
	bundle, err := world.LoadBundle(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return world.Bundle{}, 1
	}
	if _, err := world.New(bundle.Run, bundle.Initial, bundle.Events); err != nil {
		fmt.Fprintln(stderr, err)
		return world.Bundle{}, 1
	}
	return bundle, 0
}
func worldValidate(args []string, stdout, stderr io.Writer) int {
	bundle, code := oneWorldDir("validate", args, stderr)
	if code != 0 {
		return code
	}
	fmt.Fprintf(stdout, "valid: %s@%s run=%s events=%d\n", bundle.Run.WorldPackKey, bundle.Run.WorldPackVersion, bundle.Run.WorldRunID, len(bundle.Events))
	return 0
}
func worldInspect(args []string, stdout, stderr io.Writer) int {
	bundle, code := oneWorldDir("inspect", args, stderr)
	if code != 0 {
		return code
	}
	hash, err := worldcontract.CanonicalHash(bundle.Initial)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	_ = writeJSON(stdout, map[string]any{"world_run": bundle.Run, "scheduled_events": len(bundle.Events), "initial_state_hash": hash, "mode": "read-only"})
	return 0
}

type worldFlags struct {
	apply              bool
	output, until, log string
}

func worldRun(args []string, stdout, stderr io.Writer, replay bool) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "world run/replay requires a world directory")
		return 2
	}
	fs := flag.NewFlagSet("world", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var f worldFlags
	fs.BoolVar(&f.apply, "apply", false, "write event log and snapshot; default is dry-run")
	fs.StringVar(&f.output, "output", "", "output directory, required with --apply")
	fs.StringVar(&f.until, "until", "", "stop at RFC3339 virtual time")
	fs.StringVar(&f.log, "log", "", "event log to verify during replay")
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(stderr, "unexpected positional arguments")
		return 2
	}
	if f.apply && f.output == "" {
		fmt.Fprintln(stderr, "--output is required with --apply")
		return 2
	}
	if !f.apply && f.output != "" {
		fmt.Fprintln(stderr, "--output requires --apply")
		return 2
	}
	bundle, err := world.LoadBundle(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	var engine *world.Engine
	if replay {
		if f.log == "" {
			fmt.Fprintln(stderr, "--log is required for replay")
			return 2
		}
		expected, err := readWorldLog(f.log)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		engine, err = world.Replay(bundle.Run, bundle.Initial, bundle.Events, expected)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	} else {
		engine, err = world.New(bundle.Run, bundle.Initial, bundle.Events)
		if err == nil {
			if f.until != "" {
				var until time.Time
				until, err = time.Parse(time.RFC3339, f.until)
				if err == nil {
					err = engine.RunUntil(until)
				}
			} else {
				err = engine.RunAll()
			}
		}
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	snapshot, err := engine.Snapshot()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	log := engine.Log()
	mode := "dry-run"
	if f.apply {
		mode = "apply"
		if err := writeWorldArtifacts(f.output, log, snapshot); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	_ = writeJSON(stdout, map[string]any{"mode": mode, "replay": replay, "world_run_id": bundle.Run.WorldRunID, "event_count": len(log), "state_hash": snapshot.StateHash, "through_sequence": snapshot.ThroughSequence, "artifacts_written": f.apply})
	return 0
}
func readWorldLog(path string) ([]worldcontract.WorldEvent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var wrapper struct {
		SchemaVersion string                     `json:"schema_version"`
		Events        []worldcontract.WorldEvent `json:"events"`
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&wrapper); err != nil {
		return nil, err
	}
	if wrapper.SchemaVersion != worldcontract.SchemaVersion {
		return nil, fmt.Errorf("unsupported event log schema_version %q", wrapper.SchemaVersion)
	}
	for i := range wrapper.Events {
		if err := wrapper.Events[i].Validate(); err != nil {
			return nil, fmt.Errorf("event log item %d: %w", i, err)
		}
	}
	return wrapper.Events, nil
}
func writeWorldArtifacts(dir string, log []worldcontract.WorldEvent, snapshot worldcontract.Snapshot) error {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	if err := writeArtifact(filepath.Join(dir, "event-log.json"), map[string]any{"schema_version": worldcontract.SchemaVersion, "events": log}); err != nil {
		return err
	}
	return writeArtifact(filepath.Join(dir, "snapshot.json"), snapshot)
}
func writeArtifact(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o640); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
