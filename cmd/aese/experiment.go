package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/industrial-ai/iaos-aese/internal/experiment"
)

func experimentCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "experiment requires validate|inspect|expand|run|compare|evidence|replay")
		return 2
	}
	name := args[0]
	fs := flag.NewFlagSet("experiment "+name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	file := fs.String("definition", "world-packs/hctm-genesis/experiments/m14/experiment.json", "strict experiment definition")
	apply := fs.Bool("apply", false, "execute isolated matrix; default preflight only")
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		return 2
	}
	d, err := loadExperiment(*file)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err = experiment.Validate(d); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	switch name {
	case "validate":
		return emit(stdout, map[string]any{"valid": true, "definition": d.Code, "production_writes": 0})
	case "inspect", "expand":
		return emit(stdout, map[string]any{"definition": d, "run_count": len(d.Profiles) * len(d.Policies) * len(d.Seeds), "estimated_artifact_bytes": len(d.Profiles) * len(d.Policies) * len(d.Seeds) * 4096, "dry_run": true})
	case "run":
		if !*apply {
			return emit(stdout, map[string]any{"dry_run": true, "apply_required": true, "run_count": len(d.Profiles) * len(d.Policies) * len(d.Seeds), "production_writes": 0})
		}
	case "compare", "evidence", "replay":
	default:
		fmt.Fprintln(stderr, "unknown experiment command")
		return 2
	}
	e, err := experiment.BuildEvidence(d)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return emit(stdout, e)
}

func loadExperiment(path string) (experiment.Definition, error) {
	f, err := os.Open(path)
	if err != nil {
		return experiment.Definition{}, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	var d experiment.Definition
	if err = dec.Decode(&d); err != nil {
		return d, err
	}
	if dec.Decode(&struct{}{}) != io.EOF {
		return d, fmt.Errorf("definition must contain one JSON value")
	}
	return d, nil
}
func emit(w io.Writer, v any) int {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	if err := e.Encode(v); err != nil {
		return 1
	}
	return 0
}

var _ = strings.TrimSpace
