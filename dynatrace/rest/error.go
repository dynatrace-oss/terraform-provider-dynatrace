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
	"errors"
	"fmt"
	"net/http"
	"strings"

	api2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api"
)

type errorEnvelope struct {
	Error *Error `json:"error"`
}

type Warning struct {
	Message string
}

func (w Warning) Error() string {
	return w.Message
}

type Error struct {
	Code                 int                   `json:"code"`
	Message              string                `json:"message"`
	ConstraintViolations []ConstraintViolation `json:"constraintViolations,omitempty"`
	URL                  string                `json:"-"`
	Method               string                `json:"-"`
}

// IsNotFoundError checks if an error is a 404 error
// api_token_client client, hybrid_client, etc. are using Error while core clients like automation use api2.APIError
func IsNotFoundError(err error) bool {
	var restErr Error
	return errors.As(err, &restErr) && restErr.Code == http.StatusNotFound || api2.IsNotFoundError(err)
}

func ConstraintViolations(err error) []ConstraintViolation {
	if err == nil {
		return []ConstraintViolation{}
	}
	if restErr, ok := err.(Error); ok {
		return restErr.ConstraintViolations
	} else if restErr, ok := err.(*Error); ok {
		return restErr.ConstraintViolations
	}
	return []ConstraintViolation{}
}

func Is404(err error) bool {
	if err == nil {
		return false
	}
	if restErr, ok := err.(Error); ok {
		return restErr.Code == 404
	}
	if restErr, ok := err.(*Error); ok {
		return restErr.Code == 404
	}
	return false
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

// ViolationMessage returns an error message containing all constraint violations of the given error.
// If there are no constraint violations, the message of the given error will be used
func (me Error) ViolationMessage() string {
	if len(me.ConstraintViolations) == 0 {
		return me.Message
	}
	messages := make([]string, 0, len(me.ConstraintViolations))
	for _, violation := range me.ConstraintViolations {
		messages = append(messages, removeSchemaAndObjectIndexFromConstraintViolationPath(violation.Path)+": "+violation.Message)
	}
	return strings.Join(messages, "\n")
}

func removeSchemaAndObjectIndexFromConstraintViolationPath(path string) string {
	pathSubstrings := strings.SplitN(path, "/", 3)
	if len(pathSubstrings) < 3 {
		return path
	}
	return strings.Join(pathSubstrings[2:], "/")
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

func IsRequiresOAuthError(err error) bool {
	var restErr Error
	if !errors.As(err, &restErr) {
		return false
	}

	if restErr.Code != http.StatusBadRequest {
		return false
	}

	return restErr.ContainsViolation("Could not do validation as request was not done using oAuth.") || restErr.Message == "No OAuth token provided"
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

var NoAPITokenError = errors.New(`No API Token has been specified.

Specifying an API Token:
- environment variable 'DYNATRACE_API_TOKEN'
- provider configuration attribute 'dt_api_token'`)

var NoOAuthCredentialsError = errors.New(`Neither OAuth Credentials nor Platform Token have been specified.

Specifying OAuth credentials:
- environment variables 'DYNATRACE_CLIENT_ID' and 'DYNATRACE_CLIENT_SECRET'
- provider configuration attributes 'client_id' and 'client_secret'
Specifying a Platform Token:

- environment variable 'DYNATRACE_PLATFORM_TOKEN'
- provider configuration attribute 'platform_token'`)

var NoClusterURLError = errors.New(`No Cluster URL has been specified.

Specifying an Cluster URL:
- environment variable 'DYNATRACE_CLUSTER_URL'
- provider configuration attribute 'dt_cluster_url'`)

var NoClusterTokenError = errors.New(`No Cluster Token has been specified.

Specifying an Cluster Token:
- environment variable 'DYNATRACE_CLUSTER_API_TOKEN'
- provider configuration attribute 'dt_cluster_api_token'`)
