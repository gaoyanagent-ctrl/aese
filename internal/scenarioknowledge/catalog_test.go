package scenarioknowledge

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

func TestM9KnowledgeManifest(t *testing.T) {
	m, err := Load("../../scenario-packs/hctm/knowledge/m9-incorporation.json")
	if err != nil {
		t.Fatal(err)
	}
	if errs := m.Validate(); len(errs) > 0 {
		t.Fatalf("invalid manifest: %v", errs)
	}
}

func TestM9KnowledgeManifestSchema(t *testing.T) {
	schemaData, err := os.ReadFile("../../scenario-packs/hctm/knowledge/schemas/knowledge-edition.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var schemaDocument any
	if err := json.Unmarshal(schemaData, &schemaDocument); err != nil {
		t.Fatal(err)
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("https://iaos.local/aese/knowledge-edition.schema.json", schemaDocument); err != nil {
		t.Fatal(err)
	}
	schema, err := compiler.Compile("https://iaos.local/aese/knowledge-edition.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	manifestData, err := os.ReadFile("../../scenario-packs/hctm/knowledge/m9-incorporation.json")
	if err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(bytes.NewReader(manifestData))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		t.Fatal(err)
	}
	if err := schema.Validate(value); err != nil {
		t.Fatal(err)
	}
}
