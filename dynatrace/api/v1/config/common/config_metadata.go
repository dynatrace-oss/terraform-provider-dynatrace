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

package common

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ConfigMetadata Metadata useful for debugging
type ConfigMetadata struct {
	ClusterVersion               *string  `json:"clusterVersion,omitempty"`               // Dynatrace server version.
	ConfigurationVersions        []int64  `json:"configurationVersions,omitempty"`        // A Sorted list of the version numbers of the configuration.
	CurrentConfigurationVersions []string `json:"currentConfigurationVersions,omitempty"` // A Sorted list of string version numbers of the configuration.
}

func (me *ConfigMetadata) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_version": {
			Type:        schema.TypeString,
			Description: "Dynatrace server version",
			Optional:    true,
		},
		"configuration_versions": {
			Type:        schema.TypeList,
			Description: "A Sorted list of the version numbers of the configuration",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
		},
		"current_configuration_versions": {
			Type:        schema.TypeList,
			Description: "A Sorted list of the version numbers of the configuration",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *ConfigMetadata) MarshalHCL(properties hcl.Properties) error {
	if me.ClusterVersion != nil && len(*me.ClusterVersion) > 0 {
		if err := properties.Encode("cluster_version", me.ClusterVersion); err != nil {
			return err
		}
	}
	if err := properties.Encode("configuration_versions", me.ConfigurationVersions); err != nil {
		return err
	}
	if err := properties.Encode("current_configuration_versions", me.CurrentConfigurationVersions); err != nil {
		return err
	}
	return nil
}

func (me *ConfigMetadata) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("cluster_version"); ok {
		me.ClusterVersion = opt.NewString(value.(string))
	}
	if _, ok := decoder.GetOk("configuration_versions.#"); ok {
		me.ConfigurationVersions = []int64{}
		if entries, ok := decoder.GetOk("configuration_versions"); ok {
			for _, entry := range entries.([]any) {
				me.ConfigurationVersions = append(me.ConfigurationVersions, int64(entry.(int)))
			}
		}
	}
	if _, ok := decoder.GetOk("current_configuration_versions.#"); ok {
		me.CurrentConfigurationVersions = []string{}
		if entries, ok := decoder.GetOk("current_configuration_versions"); ok {
			for _, entry := range entries.([]any) {
				me.CurrentConfigurationVersions = append(me.CurrentConfigurationVersions, entry.(string))
			}
		}
	}
	return nil
}
