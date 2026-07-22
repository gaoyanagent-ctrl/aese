package worldcontract

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

func TestSchemasCompileAndValidateFixtures(t *testing.T) {
	root := filepath.Join("..", "..", "world-contracts")
	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat()
	entries, err := os.ReadDir(filepath.Join(root, "schemas"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, "schemas", entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		var document any
		if err := json.Unmarshal(data, &document); err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		if err := compiler.AddResource("https://aese.local/world/v1/"+entry.Name(), document); err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
	}
	fixtures, err := os.ReadDir(filepath.Join(root, "fixtures"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range fixtures {
		name := entry.Name()
		t.Run(name, func(t *testing.T) {
			schemaName := name[:len(name)-len(filepath.Ext(name))] + ".schema.json"
			schema, err := compiler.Compile("https://aese.local/world/v1/" + schemaName)
			if err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(root, "fixtures", name))
			if err != nil {
				t.Fatal(err)
			}
			decoder := json.NewDecoder(bytes.NewReader(data))
			decoder.UseNumber()
			var value any
			if err := decoder.Decode(&value); err != nil {
				t.Fatal(err)
			}
			if err := schema.Validate(value); err != nil {
				t.Fatal(err)
			}
		})
	}
}
