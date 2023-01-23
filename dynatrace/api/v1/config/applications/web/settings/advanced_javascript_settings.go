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

// AdvancedJavaScriptTagSettings Advanced JavaScript tag settings
type AdvancedJavaScriptTagSettings struct {
	SyncBeaconFirefox                   bool                        `json:"syncBeaconFirefox"`                   // Send the beacon signal as a synchronous XMLHttpRequest using Firefox enabled/disabled
	SyncBeaconInternetExplorer          bool                        `json:"syncBeaconInternetExplorer"`          // Send the beacon signal as a synchronous XMLHttpRequest using Internet Explorer enabled/disabled
	InstrumentUnsupportedAjaxFrameworks bool                        `json:"instrumentUnsupportedAjaxFrameworks"` // Instrumentation of unsupported Ajax frameworks enabled/disabled
	SpecialCharactersToEscape           string                      `json:"specialCharactersToEscape"`           // Additional special characters that are to be escaped using non-alphanumeric characters in HTML escape format. Maximum length 30 character. Allowed characters are `^`, `\`, `<` and `>`.
	MaxActionNameLength                 int32                       `json:"maxActionNameLength"`                 // Maximum character length for action names. Valid values range from 5 to 10000.
	MaxErrorsToCapture                  int32                       `json:"maxErrorsToCapture"`                  // Maximum number of errors to be captured per page. Valid values range from 0 to 50.
	AdditionalEventHandlers             *AdditionalEventHandlers    `json:"additionalEventHandlers"`             // Additional event handlers and wrappers
	EventWrapperSettings                EventWrapperSettings        `json:"eventWrapperSettings"`                // In addition to the event handlers, events called using `addEventListener` or `attachEvent` can be captured. Be careful with this option! Event wrappers can conflict with the JavaScript code on a web page
	GlobalEventCaptureSettings          *GlobalEventCaptureSettings `json:"globalEventCaptureSettings"`          // Global event capture settings
}

func (me *AdvancedJavaScriptTagSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sync_beacon_firefox": {
			Type:        schema.TypeBool,
			Description: "Send the beacon signal as a synchronous XMLHttpRequest using Firefox enabled/disabled",
			Optional:    true,
		},
		"sync_beacon_internet_explorer": {
			Type:        schema.TypeBool,
			Description: "Send the beacon signal as a synchronous XMLHttpRequest using Internet Explorer enabled/disabled",
			Optional:    true,
		},
		"instrument_unsupported_ajax_frameworks": {
			Type:        schema.TypeBool,
			Description: "Instrumentation of unsupported Ajax frameworks enabled/disabled",
			Optional:    true,
		},
		"special_characters_to_escape": {
			Type:        schema.TypeString,
			Description: "Additional special characters that are to be escaped using non-alphanumeric characters in HTML escape format. Maximum length 30 character. Allowed characters are `^`, `\\`, `<` and `>`.",
			Optional:    true,
		},
		"max_action_name_length": {
			Type:        schema.TypeInt,
			Description: "Maximum character length for action names. Valid values range from 5 to 10000.",
			Required:    true,
		},
		"max_errors_to_capture": {
			Type:        schema.TypeInt,
			Description: "Maximum number of errors to be captured per page. Valid values range from 0 to 50.",
			Required:    true,
		},
		"additional_event_handlers": {
			Type:        schema.TypeList,
			Description: "Additional event handlers and wrappers",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(AdditionalEventHandlers).Schema()},
		},
		"event_wrapper_settings": {
			Type:        schema.TypeList,
			Description: "In addition to the event handlers, events called using `addEventListener` or `attachEvent` can be captured. Be careful with this option! Event wrappers can conflict with the JavaScript code on a web page",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(EventWrapperSettings).Schema()},
		},
		"global_event_capture_settings": {
			Type:        schema.TypeList,
			Description: "Global event capture settings",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(GlobalEventCaptureSettings).Schema()},
		},
	}
}

func (me *AdvancedJavaScriptTagSettings) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeAll(map[string]any{
		"sync_beacon_firefox":                    me.SyncBeaconFirefox,
		"sync_beacon_internet_explorer":          me.SyncBeaconInternetExplorer,
		"instrument_unsupported_ajax_frameworks": me.InstrumentUnsupportedAjaxFrameworks,
		"special_characters_to_escape":           me.SpecialCharactersToEscape,
		"max_action_name_length":                 me.MaxActionNameLength,
		"max_errors_to_capture":                  me.MaxErrorsToCapture,
		"additional_event_handlers":              me.AdditionalEventHandlers,
		"event_wrapper_settings":                 &me.EventWrapperSettings,
		"global_event_capture_settings":          me.GlobalEventCaptureSettings,
	}); err != nil {
		return err
	}
	if me.EventWrapperSettings.IsDefault() {
		properties["event_wrapper_settings"] = nil
	}
	return nil
}

func (me *AdvancedJavaScriptTagSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"sync_beacon_firefox":                    &me.SyncBeaconFirefox,
		"sync_beacon_internet_explorer":          &me.SyncBeaconInternetExplorer,
		"instrument_unsupported_ajax_frameworks": &me.InstrumentUnsupportedAjaxFrameworks,
		"special_characters_to_escape":           &me.SpecialCharactersToEscape,
		"max_action_name_length":                 &me.MaxActionNameLength,
		"max_errors_to_capture":                  &me.MaxErrorsToCapture,
		"additional_event_handlers":              &me.AdditionalEventHandlers,
		"event_wrapper_settings":                 &me.EventWrapperSettings,
		"global_event_capture_settings":          &me.GlobalEventCaptureSettings,
	})
}
