package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/industrial-ai/iaos-aese/internal/agenttrace"
	"github.com/industrial-ai/iaos-aese/internal/iaosclient"
	"github.com/industrial-ai/iaos-aese/internal/legacyprojection"
	"github.com/industrial-ai/iaos-aese/internal/replay"
	"github.com/industrial-ai/iaos-aese/internal/scenariopack"
)

type ClientConfig struct {
	BaseURL  string
	Token    string
	TenantID string
}

type ReplayOptions struct {
	Apply    bool
	Target   string
	Tenant   string
	Actor    string
	PackKey  string
	OrderID  string
	Entities map[string]bool
}

type VerifyOptions struct {
	Target string
	Tenant string
	Actor  string
}

func NewIAOSClient(config ClientConfig) (*iaosclient.Client, error) {
	return iaosclient.New(iaosclient.Config{BaseURL: config.BaseURL, Token: config.Token, TenantID: config.TenantID})
}

func NewReplayOptions(values ReplayOptions) replay.Options {
	return replay.Options{
		Apply:    values.Apply,
		Target:   values.Target,
		Tenant:   values.Tenant,
		Actor:    values.Actor,
		PackKey:  values.PackKey,
		OrderID:  values.OrderID,
		Entities: values.Entities,
	}
}

func EffectiveRunID(value, prefix string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return fmt.Sprintf("hctm-%s-%d", prefix, time.Now().UTC().UnixNano())
}

func ParseEntityAllowlist(raw string) map[string]bool {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	out := map[string]bool{}
	for _, part := range strings.Split(raw, ",") {
		if value := strings.TrimSpace(part); value != "" {
			out[value] = true
		}
	}
	return out
}

func ApplyScenario(ctx context.Context, client *iaosclient.Client, pack *scenariopack.Pack, storyKey, runID string, apply bool) (iaosclient.ScenarioSummary, []legacyprojection.Warning, error) {
	story, err := FindStory(pack, storyKey)
	if err != nil {
		return iaosclient.ScenarioSummary{}, nil, err
	}
	projection, err := legacyprojection.Project(pack, legacyprojection.Options{StoryKey: story.Ref.Key, RunID: runID, DryRun: !apply})
	if err != nil {
		return iaosclient.ScenarioSummary{}, nil, err
	}
	summary, err := client.ApplyScenario(ctx, projection.Request, story.Events.CorrelationID)
	return summary, projection.Warnings, err
}

func ReplayScenario(ctx context.Context, client *iaosclient.Client, pack *scenariopack.Pack, storyKey string, opts ReplayOptions) (replay.ReplaySummary, error) {
	story, err := FindStory(pack, storyKey)
	if err != nil {
		return replay.ReplaySummary{}, err
	}
	runner, err := replay.New(client)
	if err != nil {
		return replay.ReplaySummary{}, err
	}
	input := NewReplayOptions(opts)
	if input.PackKey == "" && pack != nil {
		input.PackKey = pack.Manifest.PackKey
	}
	summary, err := runner.Replay(ctx, story, input)
	return summary, err
}

func VerifyScenario(ctx context.Context, client *iaosclient.Client, pack *scenariopack.Pack, storyKey string, opts VerifyOptions) (replay.VerifySummary, error) {
	story, err := FindStory(pack, storyKey)
	if err != nil {
		return replay.VerifySummary{}, err
	}
	assertions := make([]replay.Assertion, 0, len(story.Expected.IAOSAssertions))
	for _, assertion := range story.Expected.IAOSAssertions {
		operator := assertion.Operator
		if operator == "" {
			operator = assertion.Type
		}
		assertions = append(assertions, replay.Assertion{
			Key:      assertion.Key,
			Entity:   assertion.Entity,
			Match:    assertion.Match,
			Field:    assertion.Field,
			Operator: operator,
			Expected: assertion.Expected,
		})
	}
	runner, err := replay.New(client)
	if err != nil {
		return replay.VerifySummary{}, err
	}
	summary, err := runner.Verify(ctx, assertions, replay.Options{Target: opts.Target, Tenant: opts.Tenant, Actor: opts.Actor})
	return summary, err
}

func ResetScenario(ctx context.Context, client *iaosclient.Client, pack *scenariopack.Pack, storyKey, runID string, apply bool) (iaosclient.ScenarioSummary, error) {
	story, err := FindStory(pack, storyKey)
	if err != nil {
		return iaosclient.ScenarioSummary{}, err
	}
	req := iaosclient.ScenarioResetRequest{
		PackKey: pack.Manifest.PackKey, PackVersion: pack.Manifest.PackVersion,
		ScenarioKey: story.Ref.Key, RunID: runID, DryRun: !apply,
	}
	return client.ResetScenario(ctx, req, story.Events.CorrelationID)
}

func SetupAgents(ctx context.Context, client *iaosclient.Client, pack *scenariopack.Pack, apply bool) (agenttrace.SetupSummary, error) {
	if pack == nil {
		return agenttrace.SetupSummary{}, fmt.Errorf("scenario pack is required")
	}
	bundle, err := agenttrace.LoadBundle(pack.Root)
	if err != nil {
		return agenttrace.SetupSummary{}, err
	}
	if err := bundle.ValidatePack(pack.Manifest.PackKey); err != nil {
		return agenttrace.SetupSummary{}, err
	}
	return agenttrace.Setup(ctx, client, bundle, apply)
}

func RunAgents(ctx context.Context, client *iaosclient.Client, pack *scenariopack.Pack, storyKey, runID string, apply bool) (agenttrace.RunSummary, error) {
	if pack == nil {
		return agenttrace.RunSummary{}, fmt.Errorf("scenario pack is required")
	}
	story, err := FindStory(pack, storyKey)
	if err != nil {
		return agenttrace.RunSummary{}, err
	}
	bundle, err := agenttrace.LoadBundle(pack.Root)
	if err != nil {
		return agenttrace.RunSummary{}, err
	}
	if err := bundle.ValidatePack(pack.Manifest.PackKey); err != nil {
		return agenttrace.RunSummary{}, err
	}
	return agenttrace.Run(ctx, client, pack.Manifest.PackKey, story.Ref.Key, story.Events.CorrelationID, runID, apply)
}
