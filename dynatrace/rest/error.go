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
