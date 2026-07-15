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

func ToMap(v interface{}) map[string]interface{} {
	b, _ := json.Marshal(v)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m
}

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
