package application

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/industrial-ai/iaos-aese/internal/scenariopack"
)

// StageID marks a deterministic execution stage that is reused by run planning,
// runbook rendering and recovery logic.
type StageID string

const (
	StagePreflight  StageID = "preflight"
	StageInitialize StageID = "initialize"
	StageAct1       StageID = "act-1"
	StageAct2       StageID = "act-2"
	StageAct3       StageID = "act-3"
	StageAct4       StageID = "act-4"
	StageAct5       StageID = "act-5"
	StageAct6       StageID = "act-6"
	StageAct7       StageID = "act-7"
	StageAnalyze    StageID = "analyze"
	StageVerify     StageID = "verify"
	StageReset      StageID = "reset"
)

type StagePlan struct {
	StageID     StageID  `json:"stage"`
	EventIDs    []string `json:"event_ids"`
	EventTypes  []string `json:"event_types"`
	EventCount  int      `json:"event_count"`
	ActionHints []string `json:"action_hints"`
}

type Plan struct {
	PackKey      string      `json:"pack_key"`
	PackVersion  string      `json:"pack_version"`
	ScenarioKey  string      `json:"scenario_key"`
	Correlation  string      `json:"correlation_id"`
	TotalEvents  int         `json:"total_events"`
	PlanHash     string      `json:"plan_hash"`
	Stages       []StagePlan `json:"stages"`
	ActCount     int         `json:"act_count"`
	AllowableRun []string    `json:"allowable_run_actions"`
}

// PlanConfig binds the first-party stage contract for a fixed known story.
//
// The first version of M7 only ships `hctm/order-expedite-01`, so this map is
// intentionally narrow and explicit:
// - act-1: sequence 1..3
// - act-2: sequence 4..7
// - act-3: sequence 8
// - act-4: sequence 9
// - act-5: sequence 10..16
// - act-6: sequence 17..20
// - act-7: sequence 21..22
var planEventRanges = []struct {
	stage StageID
	start int
	end   int
}{
	{StageAct1, 0, 2},
	{StageAct2, 3, 6},
	{StageAct3, 7, 7},
	{StageAct4, 8, 8},
	{StageAct5, 9, 15},
	{StageAct6, 16, 19},
	{StageAct7, 20, 21},
}

// FindStory returns the matching story by key and validates record-level invariants.
func FindStory(pack *scenariopack.Pack, key string) (scenariopack.Story, error) {
	if pack == nil {
		return scenariopack.Story{}, fmt.Errorf("pack is required")
	}
	if key == "" {
		return scenariopack.Story{}, fmt.Errorf("story key is required")
	}
	for _, story := range pack.Stories {
		if story.Ref.Key == key || story.Initial.StoryKey == key {
			if len(story.Events.Events) == 0 {
				return scenariopack.Story{}, fmt.Errorf("story %s has no events", key)
			}
			return story, nil
		}
	}
	return scenariopack.Story{}, fmt.Errorf("story %q not found", key)
}

// eventStage maps a zero-based sequence index to deterministic act stage.
func eventStage(index int) (StageID, bool) {
	for _, bucket := range planEventRanges {
		if index >= bucket.start && index <= bucket.end {
			return bucket.stage, true
		}
	}
	return "", false
}

// CompilePlan builds a deterministic run plan from scenario events and hashes it.
func CompilePlan(pack *scenariopack.Pack, storyKey string) (Plan, error) {
	story, err := FindStory(pack, storyKey)
	if err != nil {
		return Plan{}, err
	}
	planByStage := make(map[StageID]*StagePlan, len(planEventRanges))
	for _, bucket := range planEventRanges {
		planByStage[bucket.stage] = &StagePlan{StageID: bucket.stage, ActionHints: []string{"advance", "run-to-end"}}
	}

	for index, event := range story.Events.Events {
		stage, ok := eventStage(index)
		if !ok {
			return Plan{}, fmt.Errorf("event %s has unsupported sequence position %d", event.EventID, index+1)
		}
		entry := planByStage[stage]
		entry.EventIDs = append(entry.EventIDs, event.EventID)
		entry.EventTypes = append(entry.EventTypes, event.EventType)
		entry.EventCount++
	}

	stages := make([]StagePlan, 0, 1+len(planByStage)+2)
	stages = append(stages,
		StagePlan{StageID: StagePreflight, ActionHints: []string{"preflight"}, EventCount: 0},
		StagePlan{StageID: StageInitialize, ActionHints: []string{"initialize"}, EventCount: 0},
	)
	for _, bucket := range planEventRanges {
		stages = append(stages, *planByStage[bucket.stage])
	}
	stages = append(stages,
		StagePlan{StageID: StageAnalyze, ActionHints: []string{"analyze"}, EventCount: 0},
		StagePlan{StageID: StageVerify, ActionHints: []string{"verify"}, EventCount: 0},
		StagePlan{StageID: StageReset, ActionHints: []string{"reset"}, EventCount: 0},
	)

	correlation := story.Events.CorrelationID
	output := Plan{
		PackKey:      pack.Manifest.PackKey,
		PackVersion:  pack.Manifest.PackVersion,
		ScenarioKey:  story.Ref.Key,
		Correlation:  correlation,
		TotalEvents:  len(story.Events.Events),
		Stages:       stages,
		ActCount:     len(planEventRanges),
		AllowableRun: []string{"preflight", "initialize", "advance", "run-to-end", "analyze", "verify", "reset"},
	}
	planHash, err := hashPlan(output)
	if err != nil {
		return Plan{}, fmt.Errorf("compute plan hash: %w", err)
	}
	output.PlanHash = planHash
	if correlation == "" {
		return output, fmt.Errorf("story %s missing correlation_id", storyKey)
	}
	if output.TotalEvents == 0 {
		return output, fmt.Errorf("story %s has zero events", storyKey)
	}
	return output, nil
}

func hashPlan(plan Plan) (string, error) {
	encoded, err := json.Marshal(plan)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}
