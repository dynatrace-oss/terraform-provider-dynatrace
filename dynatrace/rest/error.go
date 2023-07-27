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

package rest

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type errorEnvelope struct {
	Error *Error `json:"error"`
}

type Error struct {
	Code                 int                   `json:"code"`
	Message              string                `json:"message"`
	ConstraintViolations []ConstraintViolation `json:"constraintViolations,omitempty"`
	URL                  string                `json:"-"`
	Method               string                `json:"-"`
}

// CorrectPayload investigates a failed request and the
// error that got returned.
// If it is an error message originating from Settings 2.0
// it may contain constraint violations complaining about
// unknown properties within the payload.
//
// If it was possible to modify the payload of the request
// based on these constraint violations, CorrectPayload will
// return `true`, signalling that it makes sense to re-send
// the given request (with corrected payload).
//
// Note: The `payload` field of the request won't be the
// original type after this.
func CorrectPayload(err error, request *request) bool {
	if request.payload == nil {
		return false
	}
	if resterr, ok := err.(Error); ok {
		if violations := resterr.PropertyViolations(); len(violations) > 0 {
			return violations.Remove(request)
		}
	}
	return false
}

// PropertyViolations represents ALL constraint violations the
// Settings 2.0 API returned complaining about properties that
// were not expected
type PropertyViolations []PropertyViolation

// PropertyViolation represents the "address" of a property
// the Settings 2.0 API wasn't expecting at this position.
// Example: [ "0", "eventFilters", "0", "unsupportedProperty" ]
//
// These entries originate from the `path` of a constraint
// violation, e.g. `builtin:alerting.profile/0/eventFilters/0/unsupportedProperty`
// The first part of the path is always the schemaID and is irrelevant
//
// Entries that are convertible to an int are considered to be indizes
// of arrays within the payload
//
// The first entry ALWAYS denotes an array index, because at least one
// setting has been sent for create or update
//
// The last entry ALWAYS denotes the name of the property
type PropertyViolation []string

// If an error returned by the Settings 2.0 API contains constraint violations
// based on missing properties, this function tries to remove them from the
// payload of the given request.
// If it was possible to "correct" the payload according to the constraint
// violations, `true` will be returned, otherwise `false`
func (me PropertyViolations) Remove(request *request) bool {
	var data []byte
	var err error
	if data, err = json.Marshal(request.payload); err != nil {
		return false
	}
	var untypedPayload any
	if err := json.Unmarshal(data, &untypedPayload); err != nil {
		return false
	}
	if me.remove(untypedPayload) {
		request.payload = untypedPayload
		return true
	}
	return false
}

// The given value may either be an `[]any` (POST) or `map[string]any` (PUT)
// [{"schemaId":"...","scope":"...","value":{ ... }}]
// {"schemaId":"...","scope":"...","value":{ ... }}
// If any of the constraint violations were matching up with the given payload
// and it was indeed possible to make corrections, `true` will be returned,
// otherwise `false`
func (me PropertyViolations) remove(v any) bool {
	if len(me) == 0 {
		return false
	}
	switch typedPayload := v.(type) {
	case []any:
		removed := false
		for _, elem := range typedPayload {
			if m, ok := elem.(map[string]any); ok {
				if value, found := m["value"]; found {
					for _, violation := range me {
						res := violation.Remove([]any{value})
						removed = removed || res
					}
				}
			}
		}
		return removed
	case map[string]any:
		removed := false
		if value, found := typedPayload["value"]; found {
			for _, violation := range me {
				res := violation.Remove([]any{value})
				removed = removed || res
			}

		}
		return removed
	}
	return false
}

func (me PropertyViolation) Remove(v any) bool {
	// Sanity checks. No payload or no address pieces left?
	if v == nil || len(me) == 0 {
		return false
	}
	addr := me[0]
	switch typedV := v.(type) {
	case []any:
		// For `[]any` we expect to address an index
		if index, err := strconv.Atoi(addr); err == nil {
			// check if index is in range and if there
			// is a remainder address left to remove
			// a property
			if len(typedV) > index && len(me) > 1 {
				// remove the index from the "address"
				// and walk down one level
				return me[1:].Remove(typedV[index])
			}
		}
		return false
	case map[string]any:
		if len(me) == 1 {
			// Last entry in an "address" is the property name
			// Let's try to remove from the given map
			if _, found := typedV[addr]; found {
				delete(typedV, addr)
				return true
			}
		} else {
			// Further levels to step down available
			if elem, found := typedV[addr]; found {
				// Look up the current property name
				// remove that name from the "address"
				// and walk down one level
				return me[1:].Remove(elem)
			}
		}
		return false
	}
	return false
}

// The message of a constraint violation needs to match that regex.
// Otherwise it doesn't qualify for a complaint about an unknown property.
var unsupportedPropertyRegex = regexp.MustCompile("Given property '(.*)' with value: '.*' was not found in the schema")

func (me Error) PropertyViolations() PropertyViolations {
	violations := PropertyViolations{}
	if len(me.ConstraintViolations) == 0 {
		return violations
	}
	violationMap := map[string]PropertyViolation{}
	for _, violation := range me.ConstraintViolations {
		if len(violation.Message) == 0 {
			continue
		}
		if matches := unsupportedPropertyRegex.FindStringSubmatch(violation.Message); len(matches) > 1 {
			if property := matches[1]; len(property) > 0 {
				var address []string
				path := violation.Path
				parts := strings.Split(path, "/")
				if len(parts) > 1 {
					// builtin:alerting.profile/0/unsupportedAttribute
					// the first part is always the schemaID, we don't need it
					address = parts[1:]
				}
				if len(property) > 0 && len(address) > 0 {
					violationMap[violation.Path] = PropertyViolation(address)
				}
			}
		}
	}
	for _, violation := range violationMap {
		violations = append(violations, violation)
	}
	return violations
}

func (me Error) Error() string {
	if len(me.ConstraintViolations) == 0 {
		return me.Message
	}
	msg := fmt.Sprintf("%s %s: %s", me.Method, me.URL, me.Message)
	for _, violation := range me.ConstraintViolations {
		d, _ := json.Marshal(violation)
		msg = msg + "\n" + string(d)
	}
	return msg
}

func (me Error) ViolationMessage() string {
	if len(me.ConstraintViolations) == 0 {
		return ""
	}
	result := ""
	for _, violation := range me.ConstraintViolations {
		if len(result) > 0 {
			result = result + "\n"
		}
		if strings.Contains(violation.Message, "must not be null") || strings.Contains(violation.Message, "Element may not be null on creation") {
			result = result + violation.Path + " " + violation.Message
		} else {
			result = result + violation.Message
		}
	}
	return result
}

type ConstraintViolation struct {
	Description       string `json:"description,omitempty"`
	ParameterLocation string `json:"parameterLocation,omitempty"`
	Message           string `json:"message,omitempty"`
	Path              string `json:"path,omitempty"`
}

func ContainsViolation(err error, message string) bool {
	if restError, ok := err.(Error); ok {
		return restError.ContainsViolation(message)
	}
	return false
}

func (me Error) ContainsViolation(message string) bool {
	if len(me.ConstraintViolations) == 0 {
		return false
	}
	for _, violation := range me.ConstraintViolations {
		if violation.Message == message {
			return true
		}
	}
	return false
}
