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

package workloaddetection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	CloudFoundry *CloudFoundryCAWD        `json:"cloudFoundry"` // Enable this setting to get \n * Processes of Cloud Foundry application instances merged into process groups by Cloud Foundry application. \n *  Container resource metrics (Container group instance entities) and [related screens](https://www.dynatrace.com/support/help/shortlink/container-groups).
	Docker       *DockerCAWD              `json:"docker"`       // Enable this setting for plain Docker environments to get \n * Container resource metrics (Container group instance entities) and [related screens](https://www.dynatrace.com/support/help/shortlink/container-groups).
	Kubernetes   *KubernetesOpenShiftCAWD `json:"kubernetes"`   // Enable this setting to get \n * Insights into your Kubernetes namespaces, workloads and pods (cloud application namespace, cloud application and cloud application instance and entities). \n * Container resource metrics (container group instance entities) and [related screens](https://www.dynatrace.com/support/help/shortlink/container-groups). \n * Similar workloads merged into process groups based on defined rules (see below). \n * Version detection for services that run in Kubernetes workloads.
}

func (me *Settings) Name() string {
	return "cloud_app_workload_detection"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cloud_foundry": {
			Type:        schema.TypeList,
			Description: "Enable this setting to get \n * Processes of Cloud Foundry application instances merged into process groups by Cloud Foundry application. \n *  Container resource metrics (Container group instance entities) and [related screens](https://www.dynatrace.com/support/help/shortlink/container-groups).",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(CloudFoundryCAWD).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"docker": {
			Type:        schema.TypeList,
			Description: "Enable this setting for plain Docker environments to get \n * Container resource metrics (Container group instance entities) and [related screens](https://www.dynatrace.com/support/help/shortlink/container-groups).",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(DockerCAWD).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"kubernetes": {
			Type:        schema.TypeList,
			Description: "Enable this setting to get \n * Insights into your Kubernetes namespaces, workloads and pods (cloud application namespace, cloud application and cloud application instance and entities). \n * Container resource metrics (container group instance entities) and [related screens](https://www.dynatrace.com/support/help/shortlink/container-groups). \n * Similar workloads merged into process groups based on defined rules (see below). \n * Version detection for services that run in Kubernetes workloads.",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(KubernetesOpenShiftCAWD).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cloud_foundry": me.CloudFoundry,
		"docker":        me.Docker,
		"kubernetes":    me.Kubernetes,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cloud_foundry": &me.CloudFoundry,
		"docker":        &me.Docker,
		"kubernetes":    &me.Kubernetes,
	})
}
