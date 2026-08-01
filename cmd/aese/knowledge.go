package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/scenarioknowledge"
)

func knowledgeInstallCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(stderr, "knowledge install requires <bundle.json> --target <IAOS URL>")
		return 2
	}
	fs := flag.NewFlagSet("knowledge install", flag.ContinueOnError)
	fs.SetOutput(stderr)
	target := fs.String("target", "", "IAOS base URL")
	token := fs.String("token", os.Getenv("IAOS_TOKEN"), "IAOS bearer token (defaults to IAOS_TOKEN)")
	tenantID := fs.String("tenant", "", "optional target tenant header")
	apply := fs.Bool("apply", false, "perform installation; default validates only")
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}
	if *target == "" || *token == "" {
		fmt.Fprintln(stderr, "--target and --token/IAOS_TOKEN are required")
		return 2
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	var bundle scenarioknowledge.Bundle
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&bundle); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.Marshal(struct {
		Apply  bool                     `json:"apply"`
		Bundle scenarioknowledge.Bundle `json:"bundle"`
	}{*apply, bundle})
	request, err := http.NewRequest(http.MethodPost, strings.TrimRight(*target, "/")+"/api/v1/platform/packages/knowledge-editions/install", bytes.NewReader(payload))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	request.Header.Set("Authorization", "Bearer "+*token)
	request.Header.Set("Content-Type", "application/json")
	if *tenantID != "" {
		request.Header.Set("X-Tenant-ID", *tenantID)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fmt.Fprintf(stderr, "IAOS %d: %s\n", response.StatusCode, strings.TrimSpace(string(body)))
		return 1
	}
	fmt.Fprintln(stdout, strings.TrimSpace(string(body)))
	return 0
}
