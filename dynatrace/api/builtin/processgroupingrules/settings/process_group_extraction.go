/**
* @license
* Copyright 2026 Dynatrace LLC
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

package processgroupingrules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProcessGroupExtractions []*ProcessGroupExtraction

func (me *ProcessGroupExtractions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"process_group_extraction": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(ProcessGroupExtraction).Schema()},
		},
	}
}

func (me ProcessGroupExtractions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("process_group_extraction", me)
}

func (me *ProcessGroupExtractions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("process_group_extraction", me)
}

type ProcessGroupExtraction struct {
	Detection   DetectionConditions `json:"detection"`             // ### 2. Define detection rules\n  Define process detection rules to select processes on which this rule will apply to. **At least one rule must be defined.**
	Name        *string             `json:"name,omitempty"`        // When this field is empty, OneAgent will automatically assign the process group name based on process type and properties like executable name. If you expect that multiple processes will be matched by the rule, it is highly recommended that you fill this field because it is unspecified which process will be used as the group name source.
	PgIdSource  *GroupIdSource      `json:"pgIdSource"`            // ### 3. Define grouping rules\n  **3.1. Process group id source**
	PgiIdSource *InstanceIdSource   `json:"pgiIdSource,omitempty"` // **3.2. Process id source (optional)**\n\n  Define a property that should be used to identify your process.
	ProcessType *ProcessType        `json:"processType,omitempty"` // Note: Not all types can be detected at startup.. Restrict this rule to specific process types to avoid mixing deep monitored properties leading to confusing results. Possible values: `PROCESS_TYPE_APACHE_HTTPD`, `PROCESS_TYPE_GLASSFISH`, `PROCESS_TYPE_GO`, `PROCESS_TYPE_IBM_CICS_REGION`, `PROCESS_TYPE_IBM_IMS_CONTROL`, `PROCESS_TYPE_IBM_IMS_MPR`, `PROCESS_TYPE_IIS_APP_POOL`, `PROCESS_TYPE_JAVA`, `PROCESS_TYPE_JBOSS`, `PROCESS_TYPE_NGINX`, `PROCESS_TYPE_NODE_JS`, `PROCESS_TYPE_PHP`, `PROCESS_TYPE_RUBY`, `PROCESS_TYPE_TOMCAT`, `PROCESS_TYPE_WEBLOGIC`, `PROCESS_TYPE_WEBSPHERE`
	Report      ReportItem          `json:"report"`                // Auto reports only processes which are important - meaning deep monitored or with high resource usage. Possible values: `always`, `auto`, `never`
}

func (me *ProcessGroupExtraction) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"detection": {
			Type:        schema.TypeList,
			Description: "### 2. Define detection rules\n  Define process detection rules to select processes on which this rule will apply to. **At least one rule must be defined.**",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(DetectionConditions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "When this field is empty, OneAgent will automatically assign the process group name based on process type and properties like executable name. If you expect that multiple processes will be matched by the rule, it is highly recommended that you fill this field because it is unspecified which process will be used as the group name source.",
			Optional:    true, // nullable
		},
		"pg_id_source": {
			Type:        schema.TypeList,
			Description: "### 3. Define grouping rules\n  **3.1. Process group id source**",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(GroupIdSource).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"pgi_id_source": {
			Type:        schema.TypeList,
			Description: "**3.2. Process id source (optional)**\n\n  Define a property that should be used to identify your process.",
			Optional:    true, // nullable
			Elem:        &schema.Resource{Schema: new(InstanceIdSource).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"process_type": {
			Type:        schema.TypeString,
			Description: "Note: Not all types can be detected at startup.. Restrict this rule to specific process types to avoid mixing deep monitored properties leading to confusing results. Possible values: `PROCESS_TYPE_APACHE_HTTPD`, `PROCESS_TYPE_GLASSFISH`, `PROCESS_TYPE_GO`, `PROCESS_TYPE_IBM_CICS_REGION`, `PROCESS_TYPE_IBM_IMS_CONTROL`, `PROCESS_TYPE_IBM_IMS_MPR`, `PROCESS_TYPE_IIS_APP_POOL`, `PROCESS_TYPE_JAVA`, `PROCESS_TYPE_JBOSS`, `PROCESS_TYPE_NGINX`, `PROCESS_TYPE_NODE_JS`, `PROCESS_TYPE_PHP`, `PROCESS_TYPE_RUBY`, `PROCESS_TYPE_TOMCAT`, `PROCESS_TYPE_WEBLOGIC`, `PROCESS_TYPE_WEBSPHERE`",
			Optional:    true, // nullable
		},
		"report": {
			Type:        schema.TypeString,
			Description: "Auto reports only processes which are important - meaning deep monitored or with high resource usage. Possible values: `always`, `auto`, `never`",
			Required:    true,
		},
	}
}

func (me *ProcessGroupExtraction) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"detection":     me.Detection,
		"name":          me.Name,
		"pg_id_source":  me.PgIdSource,
		"pgi_id_source": me.PgiIdSource,
		"process_type":  me.ProcessType,
		"report":        me.Report,
	})
}

func (me *ProcessGroupExtraction) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"detection":     &me.Detection,
		"name":          &me.Name,
		"pg_id_source":  &me.PgIdSource,
		"pgi_id_source": &me.PgiIdSource,
		"process_type":  &me.ProcessType,
		"report":        &me.Report,
	})
}
