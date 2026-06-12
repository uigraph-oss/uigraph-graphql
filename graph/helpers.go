package graph

import "encoding/json"

// Shared resolver helpers. Kept in a non-resolver file so `gqlgen generate`
// never moves or comments them out.

func strFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func boolFromMap(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func optStrFromMap(m map[string]interface{}, key string) *string {
	if v, ok := m[key].(string); ok && v != "" {
		return &v
	}
	return nil
}

func unmarshalJSONString(s string, out interface{}) error {
	return json.Unmarshal([]byte(s), out)
}
