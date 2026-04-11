package format

import (
	"encoding/json"
	"testing"
)

func TestMarshalCanonical_SortedKeys(t *testing.T) {
	input := map[string]string{
		"zebra": "last",
		"apple": "first",
		"mango": "middle",
	}

	data, err := MarshalCanonical(input)
	if err != nil {
		t.Fatalf("MarshalCanonical failed: %v", err)
	}

	// Check that keys appear in sorted order
	expected := `{
  "apple": "first",
  "mango": "middle",
  "zebra": "last"
}`
	if string(data) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(data))
	}
}

func TestMarshalCanonical_NoTrailingNewline(t *testing.T) {
	input := map[string]string{"key": "value"}

	data, err := MarshalCanonical(input)
	if err != nil {
		t.Fatalf("MarshalCanonical failed: %v", err)
	}

	if len(data) > 0 && data[len(data)-1] == '\n' {
		t.Error("output has trailing newline")
	}
}

func TestMarshalCanonical_NestedObjects(t *testing.T) {
	input := map[string]any{
		"outer": map[string]any{
			"z_key": "z",
			"a_key": "a",
		},
	}

	data, err := MarshalCanonical(input)
	if err != nil {
		t.Fatalf("MarshalCanonical failed: %v", err)
	}

	// Verify nested keys are also sorted
	str := string(data)
	aPos := indexOf(str, "a_key")
	zPos := indexOf(str, "z_key")

	if aPos > zPos {
		t.Error("nested keys not sorted: a_key should appear before z_key")
	}
}

func TestMarshalCanonical_Deterministic(t *testing.T) {
	input := map[string]string{
		"c": "3",
		"a": "1",
		"b": "2",
	}

	// Marshal multiple times and verify same output
	var results []string
	for i := 0; i < 5; i++ {
		data, err := MarshalCanonical(input)
		if err != nil {
			t.Fatalf("MarshalCanonical failed: %v", err)
		}
		results = append(results, string(data))
	}

	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Errorf("non-deterministic output: run %d differs from run 0", i)
		}
	}
}

func TestMarshalCanonical_Arrays(t *testing.T) {
	input := []string{"first", "second", "third"}

	data, err := MarshalCanonical(input)
	if err != nil {
		t.Fatalf("MarshalCanonical failed: %v", err)
	}

	expected := `[
  "first",
  "second",
  "third"
]`
	if string(data) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(data))
	}
}

func TestUnmarshalCanonical(t *testing.T) {
	input := `{"key": "value", "num": 42}`

	var result map[string]any
	err := UnmarshalCanonical([]byte(input), &result)
	if err != nil {
		t.Fatalf("UnmarshalCanonical failed: %v", err)
	}

	if result["key"] != "value" {
		t.Errorf("expected key='value', got '%v'", result["key"])
	}
	if result["num"] != float64(42) {
		t.Errorf("expected num=42, got '%v'", result["num"])
	}
}

func TestMarshalCanonical_Struct(t *testing.T) {
	type TestStruct struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Label string `json:"label,omitempty"`
	}

	input := TestStruct{
		ID:   "test_id",
		Type: "function",
	}

	data, err := MarshalCanonical(input)
	if err != nil {
		t.Fatalf("MarshalCanonical failed: %v", err)
	}

	// Verify it can be unmarshaled back
	var result TestStruct
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if result.ID != input.ID {
		t.Errorf("expected ID '%s', got '%s'", input.ID, result.ID)
	}
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
