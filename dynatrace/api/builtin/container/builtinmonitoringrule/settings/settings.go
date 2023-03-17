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

package builtinmonitoringrule

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	IgnoreDockerPauseContainer     bool `json:"ignoreDockerPauseContainer"`     // Disable monitoring of platform internal pause containers in Kubernetes and OpenShift.
	IgnoreKubernetesPauseContainer bool `json:"ignoreKubernetesPauseContainer"` // Disable monitoring of platform internal pause containers in Kubernetes and OpenShift.
	IgnoreOpenShiftBuildPodName    bool `json:"ignoreOpenShiftBuildPodName"`    // Disable monitoring of intermediate containers created during image build.
	IgnoreOpenShiftSdnNamespace    bool `json:"ignoreOpenShiftSdnNamespace"`    // Disable monitoring of platform internal containers in the openshift-sdn namespace.
}

func (me *Settings) Name() string {
	return "container_builtin_rule"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ignore_docker_pause_container": {
			Type:        schema.TypeBool,
			Description: "Disable monitoring of platform internal pause containers in Kubernetes and OpenShift.",
			Required:    true,
		},
		"ignore_kubernetes_pause_container": {
			Type:        schema.TypeBool,
			Description: "Disable monitoring of platform internal pause containers in Kubernetes and OpenShift.",
			Required:    true,
		},
		"ignore_open_shift_build_pod_name": {
			Type:        schema.TypeBool,
			Description: "Disable monitoring of intermediate containers created during image build.",
			Required:    true,
		},
		"ignore_open_shift_sdn_namespace": {
			Type:        schema.TypeBool,
			Description: "Disable monitoring of platform internal containers in the openshift-sdn namespace.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"ignore_docker_pause_container":     me.IgnoreDockerPauseContainer,
		"ignore_kubernetes_pause_container": me.IgnoreKubernetesPauseContainer,
		"ignore_open_shift_build_pod_name":  me.IgnoreOpenShiftBuildPodName,
		"ignore_open_shift_sdn_namespace":   me.IgnoreOpenShiftSdnNamespace,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"ignore_docker_pause_container":     &me.IgnoreDockerPauseContainer,
		"ignore_kubernetes_pause_container": &me.IgnoreKubernetesPauseContainer,
		"ignore_open_shift_build_pod_name":  &me.IgnoreOpenShiftBuildPodName,
		"ignore_open_shift_sdn_namespace":   &me.IgnoreOpenShiftSdnNamespace,
	})
}
