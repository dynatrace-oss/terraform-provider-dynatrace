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

package technology

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	BoshProcessManager bool    `json:"boshProcessManager"` // Platform: Cloud Foundry\n\nStatus: Released\n\nOperating system: Linux\n\nMin agent version: 1.159
	Containerd         bool    `json:"containerd"`         // Platform: Kubernetes\n\nStatus: Released\n\nOperating system: Linux\n\nMin agent version: 1.169
	Crio               bool    `json:"crio"`               // Platform: Kubernetes\n\nStatus: Released\n\nOperating system: Linux\n\nMin agent version: 1.163
	Docker             bool    `json:"docker"`             // Platform: Docker and Kubernetes\n\nStatus: Released\n\nOperating system: Linux
	DockerWindows      bool    `json:"dockerWindows"`      // Platform: Docker\n\nStatus: Early adopter\n\nOperating system: Windows\n\nMin agent version: 1.149
	Garden             bool    `json:"garden"`             // Platform: Cloud Foundry\n\nStatus: Released\n\nOperating system: Linux\n\nMin agent version: 1.133
	Podman             bool    `json:"podman"`             // Platform: Podman\n\nStatus: Released\n\nOperating system: Linux\n\nMin agent version: 1.267
	Scope              *string `json:"-" scope:"scope"`    // The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.
	Winc               bool    `json:"winc"`               // Platform: Cloud Foundry\n\nStatus: Early adopter\n\nOperating system: Windows\n\nMin agent version: 1.175
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bosh_process_manager": {
			Type:        schema.TypeBool,
			Description: "Platform: Cloud Foundry\n\nStatus: Released\n\nOperating system: Linux\n\nMin agent version: 1.159",
			Required:    true,
		},
		"containerd": {
			Type:        schema.TypeBool,
			Description: "Platform: Kubernetes\n\nStatus: Released\n\nOperating system: Linux\n\nMin agent version: 1.169",
			Required:    true,
		},
		"crio": {
			Type:        schema.TypeBool,
			Description: "Platform: Kubernetes\n\nStatus: Released\n\nOperating system: Linux\n\nMin agent version: 1.163",
			Required:    true,
		},
		"docker": {
			Type:        schema.TypeBool,
			Description: "Platform: Docker and Kubernetes\n\nStatus: Released\n\nOperating system: Linux",
			Required:    true,
		},
		"docker_windows": {
			Type:        schema.TypeBool,
			Description: "Platform: Docker\n\nStatus: Early adopter\n\nOperating system: Windows\n\nMin agent version: 1.149",
			Required:    true,
		},
		"garden": {
			Type:        schema.TypeBool,
			Description: "Platform: Cloud Foundry\n\nStatus: Released\n\nOperating system: Linux\n\nMin agent version: 1.133",
			Required:    true,
		},
		"podman": {
			Type:        schema.TypeBool,
			Description: "Platform: Podman\n\nStatus: Released\n\nOperating system: Linux\n\nMin agent version: 1.267",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
		"winc": {
			Type:        schema.TypeBool,
			Description: "Platform: Cloud Foundry\n\nStatus: Early adopter\n\nOperating system: Windows\n\nMin agent version: 1.175",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"bosh_process_manager": me.BoshProcessManager,
		"containerd":           me.Containerd,
		"crio":                 me.Crio,
		"docker":               me.Docker,
		"docker_windows":       me.DockerWindows,
		"garden":               me.Garden,
		"podman":               me.Podman,
		"scope":                me.Scope,
		"winc":                 me.Winc,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"bosh_process_manager": &me.BoshProcessManager,
		"containerd":           &me.Containerd,
		"crio":                 &me.Crio,
		"docker":               &me.Docker,
		"docker_windows":       &me.DockerWindows,
		"garden":               &me.Garden,
		"podman":               &me.Podman,
		"scope":                &me.Scope,
		"winc":                 &me.Winc,
	})
}
