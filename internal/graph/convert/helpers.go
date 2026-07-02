// Package convert maps internal/uigraphapi REST DTOs onto internal/graph/model
// GraphQL models. Every function here is pure — no I/O, no context — which is
// what makes this package unit-testable without a running server.
package convert

import "encoding/json"

func StrFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func BoolFromMap(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func OptStrFromMap(m map[string]interface{}, key string) *string {
	if v, ok := m[key].(string); ok && v != "" {
		return &v
	}
	return nil
}

func UnmarshalJSONString(s string, out interface{}) error {
	return json.Unmarshal([]byte(s), out)
}

// APIEndpointInputMap converts endpoint create/update input to a REST body map.
// JSON string fields (requestBody, responses, parameters) are embedded as raw JSON
// so the REST API stores objects/arrays in JSONB instead of double-encoded strings.
func APIEndpointInputMap(input interface{}) map[string]interface{} {
	m := ToMap(input)
	for _, key := range []string{
		"requestBody", "responses", "parameters", "exampleRequests", "exampleResponses",
	} {
		s, ok := m[key].(string)
		if !ok || s == "" {
			continue
		}
		if json.Valid([]byte(s)) {
			m[key] = json.RawMessage(s)
		}
	}
	return m
}

// ToMap JSON-round-trips a struct into map[string]interface{}.
// This correctly handles optional fields: nil pointer fields are omitted
// from the resulting map (because of omitempty in the input struct tags).
func ToMap(v interface{}) map[string]interface{} {
	b, _ := json.Marshal(v)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m
}

// RawStr returns the JSON string of a raw message, defaulting to "{}".
func RawStr(b json.RawMessage) string {
	if len(b) == 0 {
		return "{}"
	}
	return string(b)
}

func RawArrStr(b json.RawMessage) string {
	if len(b) == 0 {
		return "[]"
	}
	return string(b)
}
