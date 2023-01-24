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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Bandwidth struct {
	NetworkType *string `json:"networkType,omitempty"` // The type of the preconfigured network—when editing in the browser, press `Crtl+Spacebar` to see the list of available networks
	Latency     *int    `json:"latency,omitempty"`     // The latency of the network, in milliseconds
	Download    *int    `json:"download,omitempty"`    // The download speed of the network, in bytes per second
	Upload      *int    `json:"upload,omitempty"`      // The upload speed of the network, in bytes per second
}

func (me *Bandwidth) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"network_type": {
			Type:        schema.TypeString,
			Description: "The type of the preconfigured network—when editing in the browser, press `Crtl+Spacebar` to see the list of available networks",
			Optional:    true,
		},
		"latency": {
			Type:        schema.TypeInt,
			Description: "The latency of the network, in milliseconds",
			Optional:    true,
		},
		"download": {
			Type:        schema.TypeInt,
			Description: "The download speed of the network, in bytes per second",
			Optional:    true,
		},
		"upload": {
			Type:        schema.TypeInt,
			Description: "The upload speed of the network, in bytes per second",
			Optional:    true,
		},
	}
}

func (me *Bandwidth) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("network_type", me.NetworkType); err != nil {
		return err
	}
	if err := properties.Encode("latency", me.Latency); err != nil {
		return err
	}
	if err := properties.Encode("download", me.Download); err != nil {
		return err
	}
	if err := properties.Encode("upload", me.Upload); err != nil {
		return err
	}
	return nil
}

func (me *Bandwidth) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("network_type", &me.NetworkType); err != nil {
		return err
	}
	if err := decoder.Decode("latency", &me.Latency); err != nil {
		return err
	}
	if err := decoder.Decode("download", &me.Download); err != nil {
		return err
	}
	if err := decoder.Decode("upload", &me.Upload); err != nil {
		return err
	}
	return nil
}
