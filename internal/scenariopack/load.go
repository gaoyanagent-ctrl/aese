package scenariopack

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Load(root string) (*Pack, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve pack path: %w", err)
	}
	manifestPath := filepath.Join(absRoot, "manifest.json")
	var manifest Manifest
	if err := readJSON(manifestPath, &manifest); err != nil {
		return nil, err
	}
	if manifest.SchemaVersion != SupportedSchemaVersion {
		return nil, fmt.Errorf("manifest.json: schema_version: unsupported %q (want %q)", manifest.SchemaVersion, SupportedSchemaVersion)
	}
	pack := &Pack{Root: absRoot, Manifest: manifest}
	for _, ref := range manifest.MasterData {
		path, err := safeJoin(absRoot, ref.Path)
		if err != nil {
			return nil, fmt.Errorf("manifest.json: master_data path %q: %w", ref.Path, err)
		}
		var sets []RecordSet
		if err := readRecordSets(path, &sets); err != nil {
			return nil, err
		}
		for i := range sets {
			sets[i].Source = rel(absRoot, path)
		}
		pack.RecordSets = append(pack.RecordSets, sets...)
	}
	for _, ref := range manifest.Stories {
		story := Story{Ref: ref}
		if err := readAt(absRoot, ref.InitialState, &story.Initial); err != nil {
			return nil, err
		}
		story.Initial.Source = ref.InitialState
		for i := range story.Initial.RecordSets {
			story.Initial.RecordSets[i].Source = ref.InitialState
		}
		if err := readAt(absRoot, ref.Events, &story.Events); err != nil {
			return nil, err
		}
		story.Events.Source = ref.Events
		if err := readAt(absRoot, ref.ExpectedOutcomes, &story.Expected); err != nil {
			return nil, err
		}
		story.Expected.Source = ref.ExpectedOutcomes
		pack.Stories = append(pack.Stories, story)
	}
	return pack, nil
}

func readRecordSets(path string, sets *[]RecordSet) error {
	var one RecordSet
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	if err := json.Unmarshal(data, &one); err == nil && one.Entity != "" {
		*sets = []RecordSet{one}
		return nil
	}
	var wrapper struct {
		RecordSets []RecordSet `json:"record_sets"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return fmt.Errorf("%s: invalid JSON: %w", path, err)
	}
	if len(wrapper.RecordSets) == 0 {
		return fmt.Errorf("%s: record_sets: must contain at least one set", path)
	}
	*sets = wrapper.RecordSets
	return nil
}

func readAt(root, name string, dst any) error {
	path, err := safeJoin(root, name)
	if err != nil {
		return fmt.Errorf("path %q: %w", name, err)
	}
	return readJSON(path, dst)
}

func readJSON(path string, dst any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.UseNumber()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("%s: invalid JSON: %w", path, err)
	}
	return nil
}

func safeJoin(root, name string) (string, error) {
	if name == "" || filepath.IsAbs(name) {
		return "", fmt.Errorf("path must be a non-empty relative path")
	}
	path := filepath.Clean(filepath.Join(root, filepath.FromSlash(name)))
	relPath, err := filepath.Rel(root, path)
	if err != nil || relPath == ".." || strings.HasPrefix(relPath, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes pack root")
	}
	return path, nil
}

func rel(root, path string) string {
	value, _ := filepath.Rel(root, path)
	return filepath.ToSlash(value)
}
