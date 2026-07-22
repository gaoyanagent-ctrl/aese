package legacyprojection

import (
	"encoding/json"
	"github.com/industrial-ai/iaos-aese/internal/scenariopack"
	"github.com/industrial-ai/iaos-aese/internal/simevent"
	"github.com/industrial-ai/iaos-aese/internal/worldcontract"
)

func ToWorldEvents(story scenariopack.Story) []simevent.Scheduled {
	out := make([]simevent.Scheduled, 0, len(story.Events.Events))
	for _, e := range story.Events.Events {
		payload, _ := json.Marshal(map[string]any{"key": "legacy:" + e.EventID, "value": map[string]any{"event_type": e.EventType, "payload": e.Payload}})
		out = append(out, simevent.Scheduled{EventID: "legacy-" + e.EventID, EventType: "legacy.scenario.event.v1", SimOccurredAt: e.Timestamp, Priority: 500, CorrelationID: e.Correlation(), CausationID: func() string {
			if e.Causation() == "" {
				return ""
			}
			return "legacy-" + e.Causation()
		}(), IdempotencyKey: "legacy:" + e.Idempotency(), SubjectRef: worldcontract.StableRef{Namespace: "hctm", Type: "legacy_event", Code: e.EventID}, PayloadType: "state.set.v1", Payload: payload})
	}
	return out
}
