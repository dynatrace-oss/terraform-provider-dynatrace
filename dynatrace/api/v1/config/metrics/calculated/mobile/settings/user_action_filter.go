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

package mobile

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// User action filter of a calculated web metric.
type UserActionFilter struct {
	ActionDurationFromMilliseconds *int    `json:"actionDurationFromMilliseconds,omitempty"` // Only actions with a duration more than or equal to this value (in milliseconds) are included in the metric calculation.
	ActionDurationToMilliseconds   *int    `json:"actionDurationToMilliseconds,omitempty"`   // Only actions with a duration less than or equal to this value (in milliseconds) are included in the metric calculation.
	Apdex                          *string `json:"apdex,omitempty"`                          // Only actions with the specified Apdex score are included in the metric calculation. Possible values: [ Frustrated, Satisfied, Tolerating, Unknown ]
	UserActionName                 *string `json:"userActionName,omitempty"`                 // Only actions with this name are included in the metric calculation.
	HasReportedError               *bool   `json:"hasReportedError,omitempty"`               // The error status of the actions to be included in the metric calculation: `true` or `false`
	HasHttpError                   *bool   `json:"hasHttpError,omitempty"`                   // The request error status of the actions to be included in the metric calculation: `true` or `false`
	City                           *string `json:"city,omitempty"`                           // Only actions of users from this city are included in the metric calculation. Specify geolocation ID here.
	Continent                      *string `json:"continent,omitempty"`                      // Only actions of users from this continent are included in the metric calculation. Specify geolocation ID here.
	Country                        *string `json:"country,omitempty"`                        // Only actions of users from this country are included in the metric calculation. Specify geolocation ID here.
	Region                         *string `json:"region,omitempty"`                         // Only actions of users from this region are included in the metric calculation. Specify geolocation ID here.
	OSFamily                       *string `json:"osFamily,omitempty"`                       // Only actions coming from this OS family are included in the metric calculation.
	OSVersion                      *string `json:"osVersion,omitempty"`                      // Only actions coming from this OS version are included in the metric calculation.
	AppVersion                     *string `json:"appVersion,omitempty"`                     // Only actions coming from this app version are included in the metric calculation.
	Device                         *string `json:"device,omitempty"`                         // Only actions coming from this app version are included in the metric calculation.
	Manufacturer                   *string `json:"manufacturer,omitempty"`                   // Only actions coming from devices of this manufacturer are included in the metric calculation.
	Carrier                        *string `json:"carrier,omitempty"`                        // Only actions coming from this carrier type are included in the metric calculation.
	ConnectionType                 *string `json:"connectionType,omitempty"`                 // Only actions coming from this connection type are included in the metric calculation. Possible values: [ LAN, MOBILE, OFFLINE, UNKNOWN, WIFI ]
	NetworkTechnology              *string `json:"networkTechnology,omitempty"`              // Filter by network technology
	ISP                            *string `json:"isp,omitempty"`                            // Only actions coming from this internet service provider are included in the metric calculation.
	Orientation                    *string `json:"orientation,omitempty"`                    // Only actions coming from devices with this display orientation are included in the metric calculation. Possible values: [ LANDSCAPE, PORTRAIT, UNKNOWN ]
	Resolution                     *string `json:"resolution,omitempty"`                     // Only actions coming from devices with this display resolution are included in the metric calculation. Possible values: [ CGA, DCI2K, DCI4K, DVGA, FHD, FWVGA, FWXGA, GHDPlus, HD, HQVGA, HQVGA2, HSXGA, HUXGA, HVGA, HXGA, NTSC, PAL, QHD, QQVGA, QSXGA, QUXGA, QVGA, QWXGA, QXGA, SVGA, SXGA, SXGAMinus, SXGAPlus, UGA, UHD16K, UHD4K, UHD8K, UHDPlus, UNKNOWN, UWQHD, UXGA, VGA, WHSXGA, WHUXGA, WHXGA, WQSXGA, WQUXGA, WQVGA, WQVGA2, WQVGA3, WQXGA, WQXGA2, WSVGA, WSVGA2, WSXGA, WSXGAPlus, WUXGA, WVGA, WVGA2, WXGA, WXGA2, WXGA3, WXGAPlus, XGA, XGAPLUS, _1280x854, nHD, qHD ]
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
		"apdex": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions with the specified Apdex score are included in the metric calculation. Possible values: [ Frustrated, Satisfied, Tolerating, Unknown ]",
		},
		"user_action_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions with this name are included in the metric calculation.",
		},
		"has_reported_error": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The error status of the actions to be included in the metric calculation: `true` or `false`",
		},
		"has_http_error": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The request error status of the actions to be included in the metric calculation: `true` or `false`",
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
		"app_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions coming from this app version are included in the metric calculation.",
		},
		"device": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions coming from this app version are included in the metric calculation.",
		},
		"manufacturer": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions coming from devices of this manufacturer are included in the metric calculation.",
		},
		"carrier": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions coming from this carrier type are included in the metric calculation.",
		},
		"connection_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions coming from this connection type are included in the metric calculation. Possible values: [ LAN, MOBILE, OFFLINE, UNKNOWN, WIFI ]",
		},
		"network_technology": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Filter by network technology",
		},
		"isp": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions coming from this internet service provider are included in the metric calculation.",
		},
		"orientation": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions coming from devices with this display orientation are included in the metric calculation. Possible values: [ LANDSCAPE, PORTRAIT, UNKNOWN ]",
		},
		"resolution": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Only actions coming from devices with this display resolution are included in the metric calculation. Possible values: [ CGA, DCI2K, DCI4K, DVGA, FHD, FWVGA, FWXGA, GHDPlus, HD, HQVGA, HQVGA2, HSXGA, HUXGA, HVGA, HXGA, NTSC, PAL, QHD, QQVGA, QSXGA, QUXGA, QVGA, QWXGA, QXGA, SVGA, SXGA, SXGAMinus, SXGAPlus, UGA, UHD16K, UHD4K, UHD8K, UHDPlus, UNKNOWN, UWQHD, UXGA, VGA, WHSXGA, WHUXGA, WHXGA, WQSXGA, WQUXGA, WQVGA, WQVGA2, WQVGA3, WQXGA, WQXGA2, WSVGA, WSVGA2, WSXGA, WSXGAPlus, WUXGA, WVGA, WVGA2, WXGA, WXGA2, WXGA3, WXGAPlus, XGA, XGAPLUS, _1280x854, nHD, qHD ]",
		},
	}
}

func (me *UserActionFilter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"action_duration_from_milliseconds": me.ActionDurationFromMilliseconds,
		"action_duration_to_milliseconds":   me.ActionDurationToMilliseconds,
		"apdex":                             me.Apdex,
		"user_action_name":                  me.UserActionName,
		"has_reported_error":                me.HasReportedError,
		"has_http_errors":                   me.HasHttpError,
		"city":                              me.City,
		"continent":                         me.Continent,
		"country":                           me.Country,
		"region":                            me.Region,
		"os_family":                         me.OSFamily,
		"os_version":                        me.OSVersion,
		"app_version":                       me.AppVersion,
		"device":                            me.Device,
		"manufacturer":                      me.Manufacturer,
		"carrier":                           me.Carrier,
		"connection_type":                   me.ConnectionType,
		"network_technology":                me.NetworkTechnology,
		"isp":                               me.ISP,
		"orientation":                       me.Orientation,
		"resolution":                        me.Resolution,
	})
}

func (me *UserActionFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"action_duration_from_milliseconds": &me.ActionDurationFromMilliseconds,
		"action_duration_to_milliseconds":   &me.ActionDurationToMilliseconds,
		"apdex":                             &me.Apdex,
		"user_action_name":                  &me.UserActionName,
		"has_reported_error":                &me.HasReportedError,
		"has_http_errors":                   &me.HasHttpError,
		"city":                              &me.City,
		"continent":                         &me.Continent,
		"country":                           &me.Country,
		"region":                            &me.Region,
		"os_family":                         &me.OSFamily,
		"os_version":                        &me.OSVersion,
		"app_version":                       &me.AppVersion,
		"device":                            &me.Device,
		"manufacturer":                      &me.Manufacturer,
		"carrier":                           &me.Carrier,
		"connection_type":                   &me.ConnectionType,
		"network_technology":                &me.NetworkTechnology,
		"isp":                               &me.ISP,
		"orientation":                       &me.Orientation,
		"resolution":                        &me.Resolution,
	})
}
