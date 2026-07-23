package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/industrial-ai/iaos-aese/internal/application"
	bridgeiaos "github.com/industrial-ai/iaos-aese/internal/bridge/iaos"
	"github.com/industrial-ai/iaos-aese/internal/scenariopack"
	"github.com/industrial-ai/iaos-aese/internal/validate"
	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
)

const usage = `Usage:
  aese validate <pack-dir> [--json]
  aese inspect <pack-dir> [--json]
  aese world validate <world-dir>
  aese world inspect <world-dir>
  aese world run <world-dir> [--until <RFC3339>] [--apply --output <dir>]
  aese world replay <world-dir> --log <event-log.json> [--apply --output <dir>]
  aese reconcile <bridge-journal.json>
  aese experiment validate|inspect|expand|run|compare|evidence|replay [--definition <json>] [--apply]

Online commands are available after an IAOS target is configured:
  aese apply <pack-dir> --target <url> [--apply]
  aese replay <pack-dir> --story <key> --target <url> [--apply]
  aese verify <pack-dir> --story <key> --target <url>
  aese reset <pack-dir> --story <key> --target <url> [--apply]
  aese agent-setup <pack-dir> --target <url> [--apply]
  aese agent-run <pack-dir> --story <key> --target <url> [--apply]
`

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	switch args[0] {
	case "validate":
		return validateCommand(args[1:], stdout, stderr)
	case "inspect":
		return inspectCommand(args[1:], stdout, stderr)
	case "apply":
		return applyCommand(args[1:], stdout, stderr)
	case "replay":
		return replayCommand(args[1:], stdout, stderr)
	case "verify":
		return verifyCommand(args[1:], stdout, stderr)
	case "reset":
		return resetCommand(args[1:], stdout, stderr)
	case "agent-setup":
		return agentSetupCommand(args[1:], stdout, stderr)
	case "agent-run":
		return agentRunCommand(args[1:], stdout, stderr)
	case "world":
		return worldCommand(args[1:], stdout, stderr)
	case "experiment":
		return experimentCommand(args[1:], stdout, stderr)
	case "reconcile":
		return reconcileCommand(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n%s", args[0], usage)
		return 2
	}
}

func reconcileCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "reconcile requires exactly one bridge journal JSON file")
		return 2
	}
	file, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	var entries []worldcontract.Envelope
	if err := decoder.Decode(&entries); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	report := bridgeiaos.Reconcile(entries)
	_ = writeJSON(stdout, report)
	if !report.Converged {
		return 1
	}
	return 0
}

func validateCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "validate requires exactly one pack directory")
		return 2
	}
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	jsonOutput := fs.Bool("json", false, "emit JSON")
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(stderr, "validate received unexpected positional arguments")
		return 2
	}
	pack, err := scenariopack.Load(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	result := validate.Pack(pack)
	if *jsonOutput {
		_ = writeJSON(stdout, result)
	} else if result.Valid() {
		fmt.Fprintf(stdout, "valid: %s@%s (%d record sets, %d stories)\n", pack.Manifest.PackKey, pack.Manifest.PackVersion, len(pack.RecordSets), len(pack.Stories))
	} else {
		for _, issue := range result.Issues {
			fmt.Fprintln(stderr, issue.Error())
		}
	}
	if !result.Valid() {
		return 1
	}
	return 0
}

func inspectCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "inspect requires exactly one pack directory")
		return 2
	}
	fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
	fs.SetOutput(stderr)
	jsonOutput := fs.Bool("json", true, "emit JSON")
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(stderr, "inspect received unexpected positional arguments")
		return 2
	}
	pack, err := scenariopack.Load(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	result := validate.Pack(pack)
	if !result.Valid() {
		for _, issue := range result.Issues {
			fmt.Fprintln(stderr, issue.Error())
		}
		return 1
	}
	summary := scenariopack.Inspect(pack)
	if *jsonOutput {
		_ = writeJSON(stdout, summary)
	} else {
		fmt.Fprintf(stdout, "%s@%s: %d master records, %d initial records, %d stories\n", summary.PackKey, summary.PackVersion, summary.MasterRecords, summary.InitialRecords, len(summary.Stories))
	}
	return 0
}

func writeJSON(w io.Writer, value any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

type onlineFlags struct {
	target, tenant, token, actor, story, entities, runID, orderID string
	apply                                                         bool
}

func addOnlineFlags(fs *flag.FlagSet, values *onlineFlags, withStory, withApply bool) {
	fs.StringVar(&values.target, "target", "", "IAOS base URL")
	fs.StringVar(&values.tenant, "tenant", "", "target tenant (defaults to pack tenant_template)")
	fs.StringVar(&values.token, "token", "", "IAOS bearer token (defaults to IAOS_TOKEN)")
	fs.StringVar(&values.actor, "actor", "aese-cli", "actor recorded in run summary")
	fs.StringVar(&values.entities, "entities", "", "comma-separated entity allowlist")
	fs.StringVar(&values.runID, "run-id", "", "stable run ID for apply/reset idempotency")
	fs.StringVar(&values.orderID, "order-id", "", "sales order UUID returned by scenario apply")
	if withStory {
		fs.StringVar(&values.story, "story", "", "story key")
	}
	if withApply {
		fs.BoolVar(&values.apply, "apply", false, "perform writes; default is dry-run")
	}
}

func parseOnline(name string, args []string, stderr io.Writer, withStory, withApply bool) (*scenariopack.Pack, onlineFlags, int) {
	var values onlineFlags
	if len(args) == 0 {
		fmt.Fprintf(stderr, "%s requires a pack directory\n", name)
		return nil, values, 2
	}
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	addOnlineFlags(fs, &values, withStory, withApply)
	if err := fs.Parse(args[1:]); err != nil {
		return nil, values, 2
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(stderr, "%s received unexpected positional arguments\n", name)
		return nil, values, 2
	}
	if values.target == "" {
		fmt.Fprintln(stderr, "--target is required")
		return nil, values, 2
	}
	if values.token == "" {
		values.token = os.Getenv("IAOS_TOKEN")
	}
	if values.token == "" {
		fmt.Fprintln(stderr, "IAOS bearer token is required via --token or IAOS_TOKEN")
		return nil, values, 2
	}
	pack, err := scenariopack.Load(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return nil, values, 1
	}
	validation := validate.Pack(pack)
	if !validation.Valid() {
		for _, issue := range validation.Issues {
			fmt.Fprintln(stderr, issue.Error())
		}
		return nil, values, 1
	}
	if values.tenant == "" {
		values.tenant = pack.Manifest.TenantTemplate
	}
	return pack, values, 0
}

func applyCommand(args []string, stdout, stderr io.Writer) int {
	pack, values, code := parseOnline("apply", args, stderr, true, true)
	if code != 0 {
		return code
	}
	if values.story == "" {
		fmt.Fprintln(stderr, "--story is required")
		return 2
	}
	client, err := application.NewIAOSClient(application.ClientConfig{BaseURL: values.target, Token: values.token, TenantID: values.tenant})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	summary, warnings, runErr := application.ApplyScenario(
		context.Background(),
		client,
		pack,
		values.story,
		application.EffectiveRunID(values.runID, "apply"),
		values.apply,
	)
	_ = writeJSON(stdout, map[string]any{"summary": summary, "mapping_warnings": warnings})
	if runErr != nil {
		fmt.Fprintln(stderr, runErr)
		return 1
	}
	return 0
}

func replayCommand(args []string, stdout, stderr io.Writer) int {
	pack, values, code := parseOnline("replay", args, stderr, true, true)
	if code != 0 {
		return code
	}
	client, err := application.NewIAOSClient(application.ClientConfig{BaseURL: values.target, Token: values.token, TenantID: values.tenant})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	summary, runErr := application.ReplayScenario(
		context.Background(),
		client,
		pack,
		values.story,
		application.ReplayOptions{
			Apply:    values.apply,
			Target:   values.target,
			Tenant:   values.tenant,
			Actor:    values.actor,
			PackKey:  pack.Manifest.PackKey,
			OrderID:  values.orderID,
			Entities: application.ParseEntityAllowlist(values.entities),
		},
	)
	_ = writeJSON(stdout, summary)
	if runErr != nil {
		fmt.Fprintln(stderr, runErr)
		return 1
	}
	return 0
}

func verifyCommand(args []string, stdout, stderr io.Writer) int {
	pack, values, code := parseOnline("verify", args, stderr, true, false)
	if code != 0 {
		return code
	}
	story, err := application.FindStory(pack, values.story)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if len(story.Expected.IAOSAssertions) == 0 {
		fmt.Fprintln(stderr, "story has no iaos_assertions for online verification")
		return 1
	}
	client, err := application.NewIAOSClient(application.ClientConfig{BaseURL: values.target, Token: values.token, TenantID: values.tenant})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	summary, runErr := application.VerifyScenario(
		context.Background(),
		client,
		pack,
		values.story,
		application.VerifyOptions{
			Target: values.target,
			Tenant: values.tenant,
			Actor:  values.actor,
		},
	)
	_ = writeJSON(stdout, summary)
	if runErr != nil {
		fmt.Fprintln(stderr, runErr)
		return 1
	}
	return 0
}

func resetCommand(args []string, stdout, stderr io.Writer) int {
	pack, values, code := parseOnline("reset", args, stderr, true, true)
	if code != 0 {
		return code
	}
	client, err := application.NewIAOSClient(application.ClientConfig{BaseURL: values.target, Token: values.token, TenantID: values.tenant})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	summary, runErr := application.ResetScenario(
		context.Background(),
		client,
		pack,
		values.story,
		application.EffectiveRunID(values.runID, "reset"),
		values.apply,
	)
	_ = writeJSON(stdout, summary)
	if runErr != nil {
		fmt.Fprintln(stderr, runErr)
		return 1
	}
	return 0
}

func agentSetupCommand(args []string, stdout, stderr io.Writer) int {
	pack, values, code := parseOnline("agent-setup", args, stderr, false, true)
	if code != 0 {
		return code
	}
	client, err := application.NewIAOSClient(application.ClientConfig{BaseURL: values.target, Token: values.token, TenantID: values.tenant})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	summary, err := application.SetupAgents(context.Background(), client, pack, values.apply)
	_ = writeJSON(stdout, summary)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func agentRunCommand(args []string, stdout, stderr io.Writer) int {
	pack, values, code := parseOnline("agent-run", args, stderr, true, true)
	if code != 0 {
		return code
	}
	client, err := application.NewIAOSClient(application.ClientConfig{BaseURL: values.target, Token: values.token, TenantID: values.tenant})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	summary, err := application.RunAgents(context.Background(), client, pack, values.story, application.EffectiveRunID(values.runID, "agent"), values.apply)
	_ = writeJSON(stdout, summary)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}
