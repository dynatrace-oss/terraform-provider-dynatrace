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

package web

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// User action filter of a calculated web metric.
type UserActionFilter struct {
	ActionDurationFromMilliseconds *int                  `json:"actionDurationFromMilliseconds,omitempty"` // Only actions with a duration more than or equal to this value (in milliseconds) are included in the metric calculation.
	ActionDurationToMilliseconds   *int                  `json:"actionDurationToMilliseconds,omitempty"`   // Only actions with a duration less than or equal to this value (in milliseconds) are included in the metric calculation.
	LoadAction                     *bool                 `json:"loadAction,omitempty"`                     // The status of load actions in the metric calculation: `true` or `false`
	XHRAction                      *bool                 `json:"xhrAction,omitempty"`                      // The status of xhr actions in the metric calculation: `true` or `false`
	XHRRouteChangeAction           *bool                 `json:"xhrRouteChangeAction,omitempty"`           // The status of route actions in the metric calculation: `true` or `false`
	CustomAction                   *bool                 `json:"customAction,omitempty"`                   // The status of custom actions in the metric calculation: `true` or `false`
	Apdex                          *string               `json:"apdex,omitempty"`                          // Only actions with the specified Apdex score are included in the metric calculation. Possible values: [ Frustrated, Satisfied, Tolerating, Unknown ]
	Domain                         *string               `json:"domain,omitempty"`                         // Only user actions coming from the specified domain are included in the metric calculation.
	UserActionName                 *string               `json:"userActionName,omitempty"`                 // Only actions with this name are included in the metric calculation.
	RealUser                       *bool                 `json:"realUser,omitempty"`                       // The status of actions coming from real users in the metric calculation: `true` or `false`
	Robot                          *bool                 `json:"robot,omitempty"`                          // The status of actions coming from robots in the metric calculation: `true` or `false`
	Synthetic                      *bool                 `json:"synthetic,omitempty"`                      // The status of actions coming from synthetic monitors in the metric calculation: `true` or `false`
	BrowserFamily                  *string               `json:"browserFamily,omitempty"`                  // Only user actions coming from the specified browser family are included in the metric calculation.
	BrowserType                    *string               `json:"browserType,omitempty"`                    // Only user actions coming from the specified browser type are included in the metric calculation.
	BrowserVersion                 *string               `json:"browserVersion,omitempty"`                 // Only user actions coming from the specified browser version are included in the metric calculation.
	HasCustomErrors                *bool                 `json:"hasCustomErrors,omitempty"`                // The custom error status of the actions to be included in the metric calculation: `true` or `false`
	HasAnyError                    *bool                 `json:"hasAnyError,omitempty"`                    // The error status of the actions to be included in the metric calculation: `true` or `false`
	HasHttpErrors                  *bool                 `json:"hasHttpErrors,omitempty"`                  // The request error status of the actions to be included in the metric calculation: `true` or `false`
	HasJavascriptErrors            *bool                 `json:"hasJavascriptErrors,omitempty"`            // The JavaScript error status of the actions to be included in the metric calculation: `true` or `false`
	City                           *string               `json:"city,omitempty"`                           // Only actions of users from this city are included in the metric calculation. Specify geolocation ID here.
	Continent                      *string               `json:"continent,omitempty"`                      // Only actions of users from this continent are included in the metric calculation. Specify geolocation ID here.
	Country                        *string               `json:"country,omitempty"`                        // Only actions of users from this country are included in the metric calculation. Specify geolocation ID here.
	Region                         *string               `json:"region,omitempty"`                         // Only actions of users from this region are included in the metric calculation. Specify geolocation ID here.
	IP                             *string               `json:"ip,omitempty"`                             // Only actions coming from this IP address are included in the metric calculation.
	IPV6Traffic                    *bool                 `json:"ipV6Traffic,omitempty"`                    // The IPv6 status of the actions to be included in the metric calculation: `true` or `false`
	OSFamily                       *string               `json:"osFamily,omitempty"`                       // Only actions coming from this OS family are included in the metric calculation.
	OSVersion                      *string               `json:"osVersion,omitempty"`                      // Only actions coming from this OS version are included in the metric calculation.
	HTTPErrorCode                  *int                  `json:"httpErrorCode,omitempty"`                  // The HTTP error status code of the actions to be included in the metric calculation.
	HTTPErrorCodeTo                *int                  `json:"httpErrorCodeTo,omitempty"`                // Can be used in combination with httpErrorCode to define a range of error codes that will be included in the metric calculation.
	HTTPPath                       *string               `json:"httpPath,omitempty"`                       // The request path that has been determined to be the origin of an HTTP error of the actions to be included in the metric calculation.
	CustomErrorType                *string               `json:"customErrorType,omitempty"`                // The custom error type of the actions to be included in the metric calculation.
	CustomErrorName                *string               `json:"customErrorName,omitempty"`                // The custom error name of the actions to be included in the metric calculation.
	UserActionProperties           *UserActionProperties `json:"userActionProperties,omitempty"`           // Only actions with the specified properties are included in the metric calculation.
	TargetViewName                 *string               `json:"targetViewName,omitempty"`                 // Only actions on the specified view are included in the metric calculation.
	TargetViewNameMatchType        *string               `json:"targetViewNameMatchType,omitempty"`        // Specifies the match type of the view name filter, e.g. using Contains or Equals. Defaults to Equals.
	TargetViewGroup                *string               `json:"targetViewGroup,omitempty"`                // Only actions on the specified group of views are included in the metric calculation.
	TargetViewGroupNameMatchType   *string               `json:"targetViewGroupNameMatchType,omitempty"`   // Specifies the match type of the view group filter, e.g. using Contains or Equals. Defaults to Equals.
}

func (me *UserActionFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action_duration_from_milliseconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Only actions with a duration more than or equal to this value (in milliseconds) are included in the metric calculation.",
		},
		"action_duration_to_milliseconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Only actions with a duration less than or equal to this value (in milliseconds) are included in the metric calculation.",
		},
		"load_action": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The status of load actions in the metric calculation: `true` or `false`",
		},
		"xhr_action": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The status of xhr actions in the metric calculation: `true` or `false`",
		},
		"xhr_route_change_action": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The status of route actions in the metric calculation: `true` or `false`",
		},
		"custom_action": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The status of custom actions in the metric calculation: `true` or `false`",
		},
		"apdex": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions with the specified Apdex score are included in the metric calculation. Possible values: [ Frustrated, Satisfied, Tolerating, Unknown ]",
		},
		"domain": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only user actions coming from the specified domain are included in the metric calculation.",
		},
		"user_action_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions with this name are included in the metric calculation.",
		},
		"real_user": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The status of actions coming from real users in the metric calculation: `true` or `false`",
		},
		"robot": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The status of actions coming from robots in the metric calculation: `true` or `false`",
		},
		"synthetic": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The status of actions coming from synthetic monitors in the metric calculation: `true` or `false`",
		},
		"browser_family": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only user actions coming from the specified browser family are included in the metric calculation.",
		},
		"browser_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only user actions coming from the specified browser type are included in the metric calculation.",
		},
		"browser_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only user actions coming from the specified browser version are included in the metric calculation.",
		},
		"has_custom_errors": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The custom error status of the actions to be included in the metric calculation: `true` or `false`",
		},
		"has_any_error": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The error status of the actions to be included in the metric calculation: `true` or `false`",
		},
		"has_http_errors": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The request error status of the actions to be included in the metric calculation: `true` or `false`",
		},
		"has_javascript_errors": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The JavaScript error status of the actions to be included in the metric calculation: `true` or `false`",
		},
		"city": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions of users from this city are included in the metric calculation. Specify geolocation ID here.",
		},
		"continent": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions of users from this continent are included in the metric calculation. Specify geolocation ID here.",
		},
		"country": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions of users from this country are included in the metric calculation. Specify geolocation ID here.",
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions of users from this region are included in the metric calculation. Specify geolocation ID here.",
		},
		"ip": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions coming from this IP address are included in the metric calculation.",
		},
		"ip_v6_traffic": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The IPv6 status of the actions to be included in the metric calculation: `true` or `false`",
		},
		"os_family": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions coming from this OS family are included in the metric calculation.",
		},
		"os_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions coming from this OS version are included in the metric calculation.",
		},
		"http_error_code": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The HTTP error status code of the actions to be included in the metric calculation.",
		},
		"http_error_code_to": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Can be used in combination with httpErrorCode to define a range of error codes that will be included in the metric calculation.",
		},
		"http_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The request path that has been determined to be the origin of an HTTP error of the actions to be included in the metric calculation.",
		},
		"custom_error_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The custom error type of the actions to be included in the metric calculation.",
		},
		"custom_error_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The custom error name of the actions to be included in the metric calculation.",
		},
		"user_action_properties": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "The definition of a calculated web metric.",
			Elem:        &schema.Resource{Schema: new(UserActionProperties).Schema()},
		},
		"target_view_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions on the specified view are included in the metric calculation.",
		},
		"target_view_name_match_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies the match type of the view name filter, e.g. using Contains or Equals. Defaults to Equals.",
		},
		"target_view_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions on the specified group of views are included in the metric calculation.",
		},
		"target_view_group_name_match_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies the match type of the view group filter, e.g. using Contains or Equals. Defaults to Equals.",
		},
	}
}

func (me *UserActionFilter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"action_duration_from_milliseconds": me.ActionDurationFromMilliseconds,
		"action_duration_to_milliseconds":   me.ActionDurationToMilliseconds,
		"load_action":                       me.LoadAction,
		"xhr_action":                        me.XHRAction,
		"xhr_route_change_action":           me.XHRRouteChangeAction,
		"custom_action":                     me.CustomAction,
		"apdex":                             me.Apdex,
		"domain":                            me.Domain,
		"user_action_name":                  me.UserActionName,
		"real_user":                         me.RealUser,
		"robot":                             me.Robot,
		"synthetic":                         me.Synthetic,
		"browser_family":                    me.BrowserFamily,
		"browser_type":                      me.BrowserType,
		"browser_version":                   me.BrowserVersion,
		"has_custom_errors":                 me.HasCustomErrors,
		"has_any_error":                     me.HasAnyError,
		"has_http_errors":                   me.HasHttpErrors,
		"has_javascript_errors":             me.HasJavascriptErrors,
		"city":                              me.City,
		"continent":                         me.Continent,
		"country":                           me.Country,
		"region":                            me.Region,
		"ip":                                me.IP,
		"ip_v6_traffic":                     me.IPV6Traffic,
		"os_family":                         me.OSFamily,
		"os_version":                        me.OSVersion,
		"http_error_code":                   me.HTTPErrorCode,
		"http_error_code_to":                me.HTTPErrorCodeTo,
		"http_path":                         me.HTTPPath,
		"custom_error_type":                 me.CustomErrorType,
		"custom_error_name":                 me.CustomErrorName,
		"user_action_properties":            me.UserActionProperties,
		"target_view_name":                  me.TargetViewName,
		"target_view_name_match_type":       me.TargetViewNameMatchType,
		"target_view_group":                 me.TargetViewGroup,
		"target_view_group_name_match_type": me.TargetViewGroupNameMatchType,
	})
}

func (me *UserActionFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"action_duration_from_milliseconds": &me.ActionDurationFromMilliseconds,
		"action_duration_to_milliseconds":   &me.ActionDurationToMilliseconds,
		"load_action":                       &me.LoadAction,
		"xhr_action":                        &me.XHRAction,
		"xhr_route_change_action":           &me.XHRRouteChangeAction,
		"custom_action":                     &me.CustomAction,
		"apdex":                             &me.Apdex,
		"domain":                            &me.Domain,
		"user_action_name":                  &me.UserActionName,
		"real_user":                         &me.RealUser,
		"robot":                             &me.Robot,
		"synthetic":                         &me.Synthetic,
		"browser_family":                    &me.BrowserFamily,
		"browser_type":                      &me.BrowserType,
		"browser_version":                   &me.BrowserVersion,
		"has_custom_errors":                 &me.HasCustomErrors,
		"has_any_error":                     &me.HasAnyError,
		"has_http_errors":                   &me.HasHttpErrors,
		"has_javascript_errors":             &me.HasJavascriptErrors,
		"city":                              &me.City,
		"continent":                         &me.Continent,
		"country":                           &me.Country,
		"region":                            &me.Region,
		"ip":                                &me.IP,
		"ip_v6_traffic":                     &me.IPV6Traffic,
		"os_family":                         &me.OSFamily,
		"os_version":                        &me.OSVersion,
		"http_error_code":                   &me.HTTPErrorCode,
		"http_error_code_to":                &me.HTTPErrorCodeTo,
		"http_path":                         &me.HTTPPath,
		"custom_error_type":                 &me.CustomErrorType,
		"custom_error_name":                 &me.CustomErrorName,
		"user_action_properties":            &me.UserActionProperties,
		"target_view_name":                  &me.TargetViewName,
		"target_view_name_match_type":       &me.TargetViewNameMatchType,
		"target_view_group":                 &me.TargetViewGroup,
		"target_view_group_name_match_type": &me.TargetViewGroupNameMatchType,
	})
}
