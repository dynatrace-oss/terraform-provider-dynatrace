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

package browser

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/request"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ScriptConfig contains the setup of the monitor
type ScriptConfig struct {
	UserAgent          *string                 `json:"userAgent,omitempty"`          // The user agent of the request
	Device             *Device                 `json:"device,omitempty"`             // The emulated device of the monitor—holds either the parameters of the custom device or the name and orientation of the preconfigured device.\n\nIf not set, then the Desktop preconfigured device is used
	Bandwidth          *Bandwidth              `json:"bandwidth,omitempty"`          // The emulated network conditions of the monitor.\n\nIf not set, then the full available bandwidth is used
	RequestHeaders     *request.HeadersSection `json:"requestHeaders,omitempty"`     // The list of HTTP headers to be sent with requests of the monitor
	Cookies            request.Cookies         `json:"cookies,omitempty"`            // These cookies are added before execution of the first step
	BlockRequests      []string                `json:"blockRequests,omitempty"`      // Block these URLs
	BypassCSP          bool                    `json:"bypassCSP"`                    // Bypass Content Security Policy of monitored pages
	MonitorFrames      *MonitorFrames          `json:"monitorFrames,omitempty"`      // Capture performance metrics for pages loaded in frames
	JavascriptSettings *JavascriptSettings     `json:"javaScriptSettings,omitempty"` // Custom JavaScript Agent settings
	DisableWebSecurity bool                    `json:"disableWebSecurity"`           // No documentation available
	IgnoredErrorCodes  *IgnoredErrorCodes      `json:"ignoredErrorCodes"`            // Ignore specific status codes
}

func (me *ScriptConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"disable_web_security": {
			Type:        schema.TypeBool,
			Description: "No documentation available",
			Optional:    true,
		},
		"bypass_csp": {
			Type:        schema.TypeBool,
			Description: "Bypass Content Security Policy of monitored pages",
			Optional:    true,
		},
		"monitor_frames": {
			Type:        schema.TypeBool,
			Description: "Capture performance metrics for pages loaded in frames",
			Optional:    true,
		},
		"user_agent": {
			Type:        schema.TypeString,
			Description: "The user agent of the request",
			Optional:    true,
		},
		"ignored_error_codes": {
			Type:        schema.TypeList,
			Description: "Ignore specific status codes",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(IgnoredErrorCodes).Schema()},
		},
		"javascript_setttings": {
			Type:        schema.TypeList,
			Description: "Custom JavaScript Agent settings",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(JavascriptSettings).Schema()},
		},
		"device": {
			Type:        schema.TypeList,
			Description: "The emulated device of the monitor—holds either the parameters of the custom device or the name and orientation of the preconfigured device.\n\nIf not set, then the Desktop preconfigured device is used",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Device).Schema()},
		},
		"bandwidth": {
			Type:        schema.TypeList,
			Description: "The emulated device of the monitor—holds either the parameters of the custom device or the name and orientation of the preconfigured device.\n\nIf not set, then the Desktop preconfigured device is used",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Bandwidth).Schema()},
		},
		"headers": {
			Type:        schema.TypeList,
			Description: "The list of HTTP headers to be sent with requests of the monitor",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(request.HeadersSection).Schema()},
		},
		"block": {
			Type:        schema.TypeSet,
			Description: "Block these URLs",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"cookies": {
			Type:        schema.TypeList,
			Description: "These cookies are added before execution of the first step",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(request.Cookies).Schema()},
		},
	}
}

func (me *ScriptConfig) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("disable_web_security", me.DisableWebSecurity); err != nil {
		return err
	}
	if err := properties.Encode("bypass_csp", me.BypassCSP); err != nil {
		return err
	}
	if me.MonitorFrames != nil && me.MonitorFrames.Enabled {
		if err := properties.Encode("monitor_frames", me.MonitorFrames.Enabled); err != nil {
			return err
		}
	} else {
		if err := properties.Encode("monitor_frames", false); err != nil {
			return err
		}
	}
	if err := properties.Encode("user_agent", me.UserAgent); err != nil {
		return err
	}
	if err := properties.Encode("ignored_error_codes", me.IgnoredErrorCodes); err != nil {
		return err
	}
	if err := properties.Encode("javascript_setttings", me.JavascriptSettings); err != nil {
		return err
	}
	if err := properties.Encode("device", me.Device); err != nil {
		return err
	}
	if err := properties.Encode("bandwidth", me.Bandwidth); err != nil {
		return err
	}
	if err := properties.Encode("headers", me.RequestHeaders); err != nil {
		return err
	}
	if err := properties.Encode("block", me.BlockRequests); err != nil {
		return err
	}
	if err := properties.Encode("cookies", me.Cookies); err != nil {
		return err
	}
	return nil
}

func (me *ScriptConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("disable_web_security", &me.DisableWebSecurity); err != nil {
		return err
	}
	if err := decoder.Decode("block", &me.BlockRequests); err != nil {
		return err
	}
	if err := decoder.Decode("bypass_csp", &me.BypassCSP); err != nil {
		return err
	}
	monitorFrames := false
	if err := decoder.Decode("monitor_frames", &monitorFrames); err != nil {
		return err
	}
	if monitorFrames {
		me.MonitorFrames = &MonitorFrames{Enabled: true}
	}
	if err := decoder.Decode("user_agent", &me.UserAgent); err != nil {
		return err
	}
	if err := decoder.Decode("ignored_error_codes", &me.IgnoredErrorCodes); err != nil {
		return err
	}
	if err := decoder.Decode("javascript_setttings", &me.JavascriptSettings); err != nil {
		return err
	}
	if err := decoder.Decode("device", &me.Device); err != nil {
		return err
	}
	if err := decoder.Decode("bandwidth", &me.Bandwidth); err != nil {
		return err
	}
	if _, ok := decoder.GetOk("headers.#"); ok {
		me.RequestHeaders = new(request.HeadersSection)
		if err := me.RequestHeaders.UnmarshalHCL(hcl.NewDecoder(decoder, "headers", 0)); err != nil {
			return err
		}
	}
	if err := decoder.Decode("headers", &me.RequestHeaders); err != nil {
		return err
	}
	if err := decoder.Decode("cookies", &me.Cookies); err != nil {
		return err
	}
	return nil
}

type IgnoredErrorCodes struct {
	StatusCodes              string  `json:"statusCodes"`                        // You can use exact number, range or status class mask. Multiple values can be separated by comma, i.e. 404, 405-410, 5xx
	MatchingDocumentRequests *string `json:"matchingDocumentRequests,omitempty"` // Only apply to document request matching this regex
}

func (me *IgnoredErrorCodes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"status_codes": {
			Type:        schema.TypeString,
			Description: "You can use exact number, range or status class mask. Multiple values can be separated by comma, i.e. 404, 405-410, 5xx",
			Required:    true,
		},
		"matching_document_requests": {
			Type:        schema.TypeString,
			Description: "Only apply to document request matching this regex",
			Optional:    true,
		},
	}
}

func (me *IgnoredErrorCodes) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("status_codes", me.StatusCodes); err != nil {
		return err
	}
	if err := properties.Encode("matching_document_requests", me.MatchingDocumentRequests); err != nil {
		return err
	}
	return nil
}

func (me *IgnoredErrorCodes) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("status_codes", &me.StatusCodes); err != nil {
		return err
	}
	if err := decoder.Decode("matching_document_requests", &me.MatchingDocumentRequests); err != nil {
		return err
	}
	return nil
}

type JavascriptSettings struct {
	TimeoutSettings         *TimeoutSettings         `json:"timeoutSettings"`
	CustomProperties        *string                  `json:"customProperties"`
	VisuallyCompleteOptions *VisuallyCompleteOptions `json:"visuallyCompleteOptions"`
}

func (me *JavascriptSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"timeout_settings": {
			Type:        schema.TypeList,
			Description: "Custom JavaScript Agent settings",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(TimeoutSettings).Schema()},
		},
		"custom_properties": {
			Type:        schema.TypeString,
			Description: "Additional Javascript Agent Properties",
			Optional:    true,
		},
		"visually_complete_options": {
			Type:        schema.TypeList,
			Description: "Parameters for Visually complete and Speed index calculation",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(VisuallyCompleteOptions).Schema()},
		},
	}
}

func (me *JavascriptSettings) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("timeout_settings", me.TimeoutSettings); err != nil {
		return err
	}
	if err := properties.Encode("visually_complete_options", me.VisuallyCompleteOptions); err != nil {
		return err
	}
	if me.CustomProperties != nil && len(*me.CustomProperties) > 0 {
		if err := properties.Encode("custom_properties", me.CustomProperties); err != nil {
			return err
		}
	}
	return nil
}

func (me *JavascriptSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("timeout_settings", &me.TimeoutSettings); err != nil {
		return err
	}
	if err := decoder.Decode("visually_complete_options", &me.VisuallyCompleteOptions); err != nil {
		return err
	}
	if err := decoder.Decode("custom_properties", &me.CustomProperties); err != nil {
		return err
	}
	return nil
}

type VisuallyCompleteOptions struct {
	ImageSizeThreshold int      `json:"imageSizeThreshold"` // Use this setting to define the minimum visible area per element (in pixels) for an element to be counted towards Visually complete and Speed index
	InactivityTimeout  int      `json:"inactivityTimeout"`  // The time the Visually complete module waits for inactivity and no further mutations on the page after the load action
	MutationTimeout    int      `json:"mutationTimeout"`    // The time the Visually complete module waits after an XHR or custom action closes to start the calculation
	ExcludedURLs       []string `json:"excludedUrls"`       // Use regular expressions to define URLs for images and iFrames to exclude from detection by the Visually complete module
	ExcludedElements   []string `json:"excludedElements"`   // Query CSS selectors to specify mutation nodes (elements that change) to ignore in Visually complete and Speed index calculation
}

func (me *VisuallyCompleteOptions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"image_size_threshold": {
			Type:        schema.TypeInt,
			Description: "Use this setting to define the minimum visible area per element (in pixels) for an element to be counted towards Visually complete and Speed index",
			Required:    true,
		},
		"inactivity_timeout": {
			Type:        schema.TypeInt,
			Description: "The time the Visually complete module waits for inactivity and no further mutations on the page after the load action",
			Required:    true,
		},
		"mutation_timeout": {
			Type:        schema.TypeInt,
			Description: "The time the Visually complete module waits after an XHR or custom action closes to start the calculation",
			Required:    true,
		},
		"excluded_urls": {
			Type:        schema.TypeList,
			Description: "Parameters for Visually complete and Speed index calculation",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"excluded_elements": {
			Type:        schema.TypeList,
			Description: "Query CSS selectors to specify mutation nodes (elements that change) to ignore in Visually complete and Speed index calculation",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *VisuallyCompleteOptions) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("image_size_threshold", me.ImageSizeThreshold); err != nil {
		return err
	}
	if err := properties.Encode("inactivity_timeout", me.InactivityTimeout); err != nil {
		return err
	}
	if err := properties.Encode("mutation_timeout", me.MutationTimeout); err != nil {
		return err
	}
	if err := properties.Encode("excluded_urls", me.ExcludedURLs); err != nil {
		return err
	}
	if err := properties.Encode("excluded_elements", me.ExcludedElements); err != nil {
		return err
	}
	return nil
}

func (me *VisuallyCompleteOptions) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("image_size_threshold", &me.ImageSizeThreshold); err != nil {
		return err
	}
	if err := decoder.Decode("inactivity_timeout", &me.InactivityTimeout); err != nil {
		return err
	}
	if err := decoder.Decode("mutation_timeout", &me.MutationTimeout); err != nil {
		return err
	}
	if value, ok := decoder.GetOk("excluded_urls"); ok {
		me.ExcludedURLs = []string{}
		for _, e := range value.([]any) {
			me.ExcludedURLs = append(me.ExcludedURLs, e.(string))
		}
	}
	if value, ok := decoder.GetOk("excluded_elements"); ok {
		me.ExcludedElements = []string{}
		for _, e := range value.([]any) {
			me.ExcludedElements = append(me.ExcludedElements, e.(string))
		}
	}
	return nil
}

type TimeoutSettings struct {
	TemporaryActionLimit        int `json:"temporaryActionLimit"`        // Track up to n cascading setTimeout calls
	TemporaryActionTotalTimeout int `json:"temporaryActionTotalTimeout"` // Limit cascading timeouts cumulatively to n ms
}

func (me *TimeoutSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action_limit": {
			Type:        schema.TypeInt,
			Description: "Track up to n cascading setTimeout calls",
			Required:    true,
		},
		"total_timeout": {
			Type:        schema.TypeInt,
			Description: "Limit cascading timeouts cumulatively to n ms",
			Required:    true,
		},
	}
}

func (me *TimeoutSettings) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("action_limit", int(me.TemporaryActionLimit)); err != nil {
		return err
	}
	if err := properties.Encode("total_timeout", int(me.TemporaryActionTotalTimeout)); err != nil {
		return err
	}
	return nil
}

func (me *TimeoutSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("action_limit", &me.TemporaryActionLimit); err != nil {
		return err
	}
	if err := decoder.Decode("total_timeout", &me.TemporaryActionTotalTimeout); err != nil {
		return err
	}
	return nil
}

type MonitorFrames struct {
	Enabled bool `json:"enabled"`
}
