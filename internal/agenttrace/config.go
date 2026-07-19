package agenttrace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/industrial-ai/iaos-aese/internal/iaosclient"
)

const BundleSchemaVersion = "1.0.0"

type MetadataField struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

type MetadataSchema struct {
	EntityCode  string          `json:"entity_code"`
	DisplayName string          `json:"display_name"`
	Fields      []MetadataField `json:"fields"`
}

type Bundle struct {
	SchemaVersion   string                      `json:"schema_version"`
	PackKey         string                      `json:"pack_key"`
	MetadataSchemas []MetadataSchema            `json:"metadata_schemas"`
	Tools           []iaosclient.AIToolManifest `json:"tools"`
}

func LoadBundle(packRoot string) (Bundle, error) {
	path := filepath.Join(packRoot, "agent-tools.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return Bundle{}, fmt.Errorf("read agent tool bundle: %w", err)
	}
	var bundle Bundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return Bundle{}, fmt.Errorf("decode agent tool bundle: %w", err)
	}
	if bundle.SchemaVersion != BundleSchemaVersion || bundle.PackKey == "" || len(bundle.Tools) == 0 {
		return Bundle{}, fmt.Errorf("invalid agent tool bundle header")
	}
	seen := map[string]bool{}
	for _, tool := range bundle.Tools {
		if tool.ToolKey == "" || tool.ToolType != "query" || tool.SourceRef != "entity.records" || seen[tool.ToolKey] {
			return Bundle{}, fmt.Errorf("invalid or duplicate agent tool %q", tool.ToolKey)
		}
		seen[tool.ToolKey] = true
	}
	for _, toolKey := range requiredToolKeys() {
		if !seen[toolKey] {
			return Bundle{}, fmt.Errorf("agent tool bundle is missing required tool %q", toolKey)
		}
	}
	return bundle, nil
}

func (b Bundle) ValidatePack(packKey string) error {
	if b.PackKey != packKey {
		return fmt.Errorf("agent tool bundle pack %q does not match scenario pack %q", b.PackKey, packKey)
	}
	return nil
}
