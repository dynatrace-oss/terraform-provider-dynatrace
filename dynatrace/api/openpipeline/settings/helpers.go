/**
* @license
* Copyright 2025 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

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
