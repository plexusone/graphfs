// Package format provides deterministic JSON serialization for git-friendly diffs.
package format

import (
	"bytes"
	"encoding/json"
	"sort"
)

// MarshalCanonical produces deterministic JSON output with sorted keys
// and consistent formatting for git-friendly diffs.
func MarshalCanonical(v any) ([]byte, error) {
	// First marshal to get the data
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// Unmarshal into interface{} to normalize
	var normalized any
	if err := json.Unmarshal(data, &normalized); err != nil {
		return nil, err
	}

	// Sort and re-marshal with indentation
	sorted := sortValue(normalized)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(sorted); err != nil {
		return nil, err
	}

	// Remove trailing newline for consistency
	result := buf.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result, nil
}

// UnmarshalCanonical unmarshals JSON data.
func UnmarshalCanonical(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// sortValue recursively sorts maps by key for deterministic output.
func sortValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		return sortMap(val)
	case []any:
		result := make([]any, len(val))
		for i, item := range val {
			result[i] = sortValue(item)
		}
		return result
	default:
		return v
	}
}

// sortedMap is a map wrapper that marshals with sorted keys.
type sortedMap struct {
	keys   []string
	values map[string]any
}

func sortMap(m map[string]any) *sortedMap {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	values := make(map[string]any, len(m))
	for k, v := range m {
		values[k] = sortValue(v)
	}

	return &sortedMap{keys: keys, values: values}
}

func (s *sortedMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, k := range s.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		// Marshal key
		keyBytes, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteByte(':')
		// Marshal value
		valBytes, err := json.Marshal(s.values[k])
		if err != nil {
			return nil, err
		}
		buf.Write(valBytes)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}
