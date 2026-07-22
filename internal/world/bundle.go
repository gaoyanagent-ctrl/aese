package world

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/industrial-ai/iaos-aese/internal/rules"
	"github.com/industrial-ai/iaos-aese/internal/simevent"
	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
)

type Bundle struct {
	Root    string
	Run     worldcontract.WorldRun
	Initial rules.State
	Events  []simevent.Scheduled
}

func LoadBundle(root string) (Bundle, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return Bundle{}, err
	}
	read := func(name string) ([]byte, error) {
		data, err := os.ReadFile(filepath.Join(abs, name))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return data, nil
	}
	runData, err := read("world-run.json")
	if err != nil {
		return Bundle{}, err
	}
	run, err := worldcontract.ParseStrict[worldcontract.WorldRun](runData)
	if err != nil {
		return Bundle{}, fmt.Errorf("world-run.json: %w", err)
	}
	stateData, err := read("initial-state.json")
	if err != nil {
		return Bundle{}, err
	}
	var state rules.State
	if err := strictDecode(stateData, &state); err != nil {
		return Bundle{}, fmt.Errorf("initial-state.json: %w", err)
	}
	if state == nil {
		return Bundle{}, fmt.Errorf("initial-state.json must be an object")
	}
	eventData, err := read("scheduled-events.json")
	if err != nil {
		return Bundle{}, err
	}
	var wrapper struct {
		SchemaVersion string               `json:"schema_version"`
		Events        []simevent.Scheduled `json:"events"`
	}
	if err := strictDecode(eventData, &wrapper); err != nil {
		return Bundle{}, fmt.Errorf("scheduled-events.json: %w", err)
	}
	if wrapper.SchemaVersion != worldcontract.SchemaVersion {
		return Bundle{}, fmt.Errorf("scheduled-events.json: unsupported schema_version %q", wrapper.SchemaVersion)
	}
	if len(wrapper.Events) == 0 {
		return Bundle{}, fmt.Errorf("scheduled-events.json: events are required")
	}
	return Bundle{Root: abs, Run: run, Initial: state, Events: wrapper.Events}, nil
}
func strictDecode(data []byte, dst any) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	d.UseNumber()
	if err := d.Decode(dst); err != nil {
		return err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("multiple JSON values")
		}
		return err
	}
	return nil
}
