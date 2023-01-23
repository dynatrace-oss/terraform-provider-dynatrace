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

// JavaScriptFrameworkSupport configures support of various JavaScript frameworks
type JavaScriptFrameworkSupport struct {
	Angular       bool `json:"angular"`       // AngularJS and Angular support enabled/disabled
	Dojo          bool `json:"dojo"`          // Dojo support enabled/disabled
	ExtJS         bool `json:"extJS"`         // ExtJS, Sencha Touch support enabled/disabled
	ICEfaces      bool `json:"icefaces"`      // ICEfaces support enabled/disabled
	JQuery        bool `json:"jQuery"`        // jQuery, Backbone.js support enabled/disabled
	MooTools      bool `json:"mooTools"`      // MooTools support enabled/disabled
	Prototype     bool `json:"prototype"`     // Prototype support enabled/disabled
	ActiveXObject bool `json:"activeXObject"` // ActiveXObject support enabled/disabled
}

func (me *JavaScriptFrameworkSupport) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"angular": {
			Type:        schema.TypeBool,
			Description: "AngularJS and Angular support enabled/disabled",
			Optional:    true,
		},
		"dojo": {
			Type:        schema.TypeBool,
			Description: "Dojo support enabled/disabled",
			Optional:    true,
		},
		"extjs": {
			Type:        schema.TypeBool,
			Description: "ExtJS, Sencha Touch support enabled/disabled",
			Optional:    true,
		},
		"icefaces": {
			Type:        schema.TypeBool,
			Description: "ICEfaces support enabled/disabled",
			Optional:    true,
		},
		"jquery": {
			Type:        schema.TypeBool,
			Description: "jQuery, Backbone.js support enabled/disabled",
			Optional:    true,
		},
		"moo_tools": {
			Type:        schema.TypeBool,
			Description: "MooTools support enabled/disabled",
			Optional:    true,
		},
		"prototype": {
			Type:        schema.TypeBool,
			Description: "Prototype support enabled/disabled",
			Optional:    true,
		},
		"active_x_object": {
			Type:        schema.TypeBool,
			Description: "ActiveXObject support enabled/disabled",
			Optional:    true,
		},
	}
}

func (me *JavaScriptFrameworkSupport) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"angular":         me.Angular,
		"dojo":            me.Dojo,
		"extjs":           me.ExtJS,
		"icefaces":        me.ICEfaces,
		"jquery":          me.JQuery,
		"moo_tools":       me.MooTools,
		"prototype":       me.Prototype,
		"active_x_object": me.ActiveXObject,
	})
}

func (me *JavaScriptFrameworkSupport) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"angular":         &me.Angular,
		"dojo":            &me.Dojo,
		"extjs":           &me.ExtJS,
		"icefaces":        &me.ICEfaces,
		"jquery":          &me.JQuery,
		"moo_tools":       &me.MooTools,
		"prototype":       &me.Prototype,
		"active_x_object": &me.ActiveXObject,
	})
}
