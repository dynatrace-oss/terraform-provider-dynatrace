/**
* @license
* Copyright 2020 Dynatrace LLC
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

package resources

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
)

type Attributes map[string]string

func NewAttributes(m map[string]string) Attributes {
	attributes := Attributes{}
	for k, v := range m {
		attributes[k] = v
	}
	return attributes
}

func (attributes Attributes) Keys() []string {
	result := []string{}
	if len(attributes) == 0 {
		return result
	}
	for k := range attributes {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}

func (attributes Attributes) MatchingKeys(key string) []string {
	matches := []string{}
	if len(attributes) == 0 {
		return matches
	}
	for k := range attributes {
		if k == key {
			return []string{k}
		}
		if address(key).match(address(k)) {
			matches = append(matches, k)
		}
	}
	return matches
}

func (attributes Attributes) collect(addr address, v any) {
	switch value := v.(type) {
	case string:
		attributes[string(addr)] = value
	case int, int8, int16, int32, int64, float32, float64, uint, uint8, uint16, uint32, uint64, bool:
		attributes[string(addr)] = fmt.Sprintf("%v", value)
	case map[string]any:
		for k := range value {
			attributes.collect(addr.dot(k), value[k])
		}
	case hcl.Properties:
		for k := range value {
			attributes.collect(addr.dot(k), value[k])
		}
	case []any:
		attributes[fmt.Sprintf("%s.#", addr)] = fmt.Sprintf("%v", len(value))
		for idx := range value {
			attributes.collect(addr.dot(idx), value[idx])
		}
	default:
		// data, err := json.Marshal(value)
		// if err != nil {
		// 	panic(err)
		// }
		// attributes[string(addr)] = string(data)
	}
}

func (attributes Attributes) Siblings(key string) Siblings {
	result := Siblings{}
	idx := strings.LastIndex(key, ".")
	prefix := ""
	if idx >= 0 {
		prefix = key[:idx+1]
	}
	for k, v := range attributes {
		if k != key && strings.HasPrefix(k, prefix) {
			result = append(result, Sibling{Key: strings.TrimPrefix(k, prefix), Value: v})
		}
	}
	return result
}

type Siblings []Sibling

func (siblings Siblings) String() string {
	data, _ := json.Marshal(siblings)
	return string(data)
}

func (siblings Siblings) Contains(others ...Sibling) bool {
	for _, other := range others {
		found := false
		for _, elem := range siblings {
			if elem.Equals(other) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

type Sibling struct {
	Key   string
	Value string
}

func (sibling Sibling) String() string {
	data, _ := json.Marshal(sibling)
	return string(data)
}

func (sibling Sibling) Equals(other Sibling) bool {
	return other.Key == sibling.Key && other.Value == sibling.Value
}

type address string

func (a address) split() []string {
	return strings.Split(string(a), ".")
}

func (a address) match(other address) bool {
	if strings.HasSuffix(string(other), ".#") {
		return false
	}
	if !strings.Contains(string(other), ".") {
		return false
	}
	parts := a.split()
	oparts := other.split()
	if len(parts) != len(oparts) {
		return false
	}
	for idx := range parts {
		part := parts[idx]
		opart := oparts[idx]
		if part == opart {
			continue
		}
		_, err := strconv.Atoi(part)
		isnum := err == nil
		_, err = strconv.Atoi(opart)
		oisnum := err == nil
		if isnum != oisnum {
			return false
		}
		if !isnum && (part != opart) {
			return false
		}
	}
	return true
}

func (a address) dot(key any) address {
	switch typedKey := key.(type) {
	case string:
		if len(a) == 0 {
			return address(typedKey)
		}
		return address(fmt.Sprintf("%s.%s", a, typedKey))
	case int:
		if len(a) == 0 {
			return address(fmt.Sprintf("%d", typedKey))
		}
		return address(fmt.Sprintf("%s.%d", a, typedKey))
	default:
		panic("here be dragons")
	}
}

func store(m map[string]any, key string, value string) {
	idx := strings.Index(key, ".")
	if idx < 0 {
		m[key] = value
		return
	}
	prefix := key[:idx]
	remainder := key[idx+1:]
	stored := m[prefix]
	switch container := stored.(type) {
	case map[string]any:
		store(container, remainder, value)
		return
	case []any:
		idx = strings.Index(remainder, ".")
		prefix = remainder[:idx]
		remainder = remainder[idx+1:]
		cidx, _ := strconv.Atoi(prefix)
		cont := container[cidx]
		switch tcont := cont.(type) {
		case map[string]any:
			store(tcont, remainder, value)
		case hcl.Properties:
			store(map[string]any(tcont), remainder, value)
		}
	}
}

func remove(m map[string]any, key string) {
	idx := strings.Index(key, ".")
	if idx < 0 {
		delete(m, key)
		return
	}
	prefix := key[:idx]
	remainder := key[idx+1:]
	stored := m[prefix]
	switch container := stored.(type) {
	case map[string]any:
		remove(container, remainder)
		return
	case []any:
		idx = strings.Index(remainder, ".")
		prefix = remainder[:idx]
		remainder = remainder[idx+1:]
		cidx, _ := strconv.Atoi(prefix)
		cont := container[cidx]
		switch tcont := cont.(type) {
		case map[string]any:
			remove(tcont, remainder)
		case hcl.Properties:
			remove(map[string]any(tcont), remainder)
		}
	}
}
