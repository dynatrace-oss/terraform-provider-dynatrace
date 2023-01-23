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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Device struct {
	Name         *string      `json:"deviceName,omitempty"`   // The name of the preconfigured device—when editing in the browser, press `Crtl+Spacebar` to see the list of available devices
	Orientation  *Orientation `json:"orientation,omitempty"`  // The orientation of the device. Possible values are `portrait` or `landscape`. Desktop and laptop devices are not allowed to use the `portrait` orientation
	Mobile       *bool        `json:"mobile,omitempty"`       // The flag of the mobile device.\nSet to `true` for mobile devices or `false` for a desktop or laptop. Required if `touchEnabled` is specified.
	TouchEnabled *bool        `json:"touchEnabled,omitempty"` // The flag of the touchscreen.\nSet to `true` if the device uses touchscreen. In that case, use can set interaction event as `tap`. Required if `mobile` is specified.
	Width        *int         `json:"width,omitempty"`        // The width of the screen in pixels.\nThe maximum allowed width is `1920`
	Height       *int         `json:"height,omitempty"`       // The height of the screen in pixels.\nThe maximum allowed width is `1080`
	ScaleFactor  *float64     `json:"scaleFactor,omitempty"`  // The pixel ratio of the device
}

func (me *Device) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the preconfigured device—when editing in the browser, press `Crtl+Spacebar` to see the list of available devices",
			Optional:    true,
		},
		"orientation": {
			Type:        schema.TypeString,
			Description: "The orientation of the device. Possible values are `portrait` or `landscape`. Desktop and laptop devices are not allowed to use the `portrait` orientation",
			Optional:    true,
		},
		"mobile": {
			Type:        schema.TypeBool,
			Description: "The flag of the mobile device.\nSet to `true` for mobile devices or `false` for a desktop or laptop.",
			Optional:    true,
		},
		"touch_enabled": {
			Type:        schema.TypeBool,
			Description: "The flag of the touchscreen.\nSet to `true` if the device uses touchscreen. In that case, use can set interaction event as `tap`.",
			Optional:    true,
		},
		"width": {
			Type:        schema.TypeInt,
			Description: "The width of the screen in pixels.\nThe maximum allowed width is `1920`.",
			Optional:    true,
		},
		"height": {
			Type:        schema.TypeInt,
			Description: "The height of the screen in pixels.\nThe maximum allowed width is `1080`.",
			Optional:    true,
		},
		"scale_factor": {
			Type:        schema.TypeFloat,
			Description: "The pixel ratio of the device.",
			Optional:    true,
		},
	}
}

func (me *Device) MarshalHCL(properties hcl.Properties) error {
	if me.Name != nil && len(*me.Name) > 0 {
		if err := properties.Encode("name", me.Name); err != nil {
			return err
		}
	}
	if err := properties.Encode("orientation", me.Orientation); err != nil {
		return err
	}
	if err := properties.Encode("mobile", me.Mobile); err != nil {
		return err
	}
	if err := properties.Encode("touch_enabled", me.TouchEnabled); err != nil {
		return err
	}
	if err := properties.Encode("width", me.Width); err != nil {
		return err
	}
	if err := properties.Encode("height", me.Height); err != nil {
		return err
	}
	if err := properties.Encode("scale_factor", me.ScaleFactor); err != nil {
		return err
	}
	return nil
}

func (me *Device) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("name", &me.Name); err != nil {
		return err
	}
	if err := decoder.Decode("orientation", &me.Orientation); err != nil {
		return err
	}
	if err := decoder.Decode("mobile", &me.Mobile); err != nil {
		return err
	}
	if err := decoder.Decode("touch_enabled", &me.TouchEnabled); err != nil {
		return err
	}
	if me.Mobile == nil && me.Name == nil {
		me.Mobile = opt.NewBool(false)
	}
	if me.TouchEnabled == nil && me.Name == nil {
		me.TouchEnabled = opt.NewBool(false)
	}
	if err := decoder.Decode("width", &me.Width); err != nil {
		return err
	}
	if err := decoder.Decode("height", &me.Height); err != nil {
		return err
	}
	if err := decoder.Decode("scale_factor", &me.ScaleFactor); err != nil {
		return err
	}
	return nil
}
