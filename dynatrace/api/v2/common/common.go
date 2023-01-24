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

package common

type SettingsObjectList struct {
	Items []*SettingsObjectListItem `json:"items"`
}

type SettingsObjectListItem struct {
	ObjectID string `json:"objectId"`
}

type SettingsObjectUpdate struct {
	SchemaVersion string `json:"schemaVersion"`
	Value         any    `json:"value"`
}

type SettingsObjectCreate struct {
	SchemaVersion string `json:"schemaVersion"`
	SchemaID      string `json:"schemaId"`
	Scope         string `json:"scope"`
	Value         any    `json:"value"`
}

type SettingsObjectResponse struct {
	ObjectID string `json:"objectId"` // The ID of the settings object
	Code     int32  // The HTTP status code for the object
}

type SettingsObjectErrorResponse struct {
	InvalidValue map[string]any `json:"invalidValue,omitempty"` // The value of the setting. \n\n It defines the actual values of settings' parameters. \n\nThe actual content depends on the object's schema.
	Error        *Error         `json:"error,omitempty"`        // Error details
	Code         *int32         `json:"code,omitempty"`         // The HTTP status code for the object
}

type Error struct {
	ConstraintViolations []*ConstraintViolation `json:"constraintViolations,omitempty"` // A list of constraint violations
	Message              string                 `json:"message,omitempty"`              // The error message
	Code                 int32                  `json:"code,omitempty"`                 // The HTTP status code
}

type ConstraintViolation struct {
	ParmeterLocation *ParameterLocation `json:"parameterLocation,omitempty"`
	Location         *string            `json:"location,omitempty"`
	Message          *string            `json:"message,omitempty"`
	Path             *string            `json:"path,omitempty"`
}

type ParameterLocation string

var ParameterLocations = struct {
	Path        ParameterLocation
	PayloadBody ParameterLocation
	Query       ParameterLocation
}{
	Path:        ParameterLocation("PATH"),
	PayloadBody: ParameterLocation("PAYLOAD_BODY"),
	Query:       ParameterLocation("QUERY"),
}
