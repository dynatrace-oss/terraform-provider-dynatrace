package jsonmodel

import (
	"encoding/json"
)

func ExtractType(message json.RawMessage) (string, error) {
	s := struct {
		Type string `json:"type"`
	}{}
	if err := json.Unmarshal(message, &s); err != nil {
		return "", err
	}

	return s.Type, nil
}

func ExtractEndpointProcessorType(ep EndpointProcessor) (string, error) {
	return ExtractType(ep.union)
}

func ClonePtr[T any](original *T) *T {
	newT := new(T)
	*newT = *original
	return newT
}
