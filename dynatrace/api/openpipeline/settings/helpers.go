package openpipeline

import "encoding/json"

// ExtractType extracts the value of type field as a string from the specified raw message.
func ExtractType(message json.RawMessage) (string, error) {
	s := struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(message, &s); err != nil {
		return "", err
	}

	return s.Type, nil
}

// MarshalAsJSONWithType converts the specified value to JSON with an additional type field.
func MarshalAsJSONWithType(v any, ttype string) (json.RawMessage, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	m["type"] = ttype

	return json.Marshal(m)
}
