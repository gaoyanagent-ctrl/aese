package agenttrace

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/industrial-ai/iaos-aese/internal/iaosclient"
)

type SetupSummary struct {
	Mode            string   `json:"mode"`
	PackKey         string   `json:"pack_key"`
	MetadataSchemas int      `json:"metadata_schemas"`
	ToolsCreated    int      `json:"tools_created"`
	ToolsUpdated    int      `json:"tools_updated"`
	ToolKeys        []string `json:"tool_keys"`
}

func Setup(ctx context.Context, client *iaosclient.Client, bundle Bundle, apply bool) (SetupSummary, error) {
	summary := SetupSummary{Mode: "dry-run", PackKey: bundle.PackKey, MetadataSchemas: len(bundle.MetadataSchemas)}
	for _, tool := range bundle.Tools {
		summary.ToolKeys = append(summary.ToolKeys, tool.ToolKey)
	}
	if !apply {
		summary.ToolsCreated = len(bundle.Tools)
		return summary, nil
	}
	summary.Mode = "apply"
	for _, schema := range bundle.MetadataSchemas {
		fields, err := json.Marshal(schema.Fields)
		if err != nil {
			return summary, err
		}
		if err := client.UpsertMetadataSchema(ctx, schema.EntityCode, iaosclient.MetadataSchemaRequest{DisplayName: schema.DisplayName, Fields: fields}); err != nil {
			return summary, fmt.Errorf("upsert metadata %s: %w", schema.EntityCode, err)
		}
	}
	existing, err := client.ListAITools(ctx, true)
	if err != nil {
		return summary, fmt.Errorf("list AI tools: %w", err)
	}
	known := make(map[string]bool, len(existing))
	for _, tool := range existing {
		known[tool.ToolKey] = true
	}
	for _, tool := range bundle.Tools {
		if known[tool.ToolKey] {
			if err := client.UpdateAITool(ctx, tool); err != nil {
				return summary, fmt.Errorf("update AI tool %s: %w", tool.ToolKey, err)
			}
			summary.ToolsUpdated++
		} else {
			if err := client.CreateAITool(ctx, tool); err != nil {
				return summary, fmt.Errorf("create AI tool %s: %w", tool.ToolKey, err)
			}
			summary.ToolsCreated++
		}
		if err := client.EnableAITool(ctx, tool.ToolKey); err != nil {
			return summary, fmt.Errorf("enable AI tool %s: %w", tool.ToolKey, err)
		}
	}
	return summary, nil
}
