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

package dataminingblocklist

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	BlocklistEntries            DataminingBlocklistEntries `json:"blocklistEntries,omitempty"`            // You can exclude specific data buckets and tables from the semantic index. Learn more about [configuring data access](https://dt-url.net/lc62i1q \"Dynatrace Generative AI data access\").
	EnableCopilot               bool                       `json:"enableCopilot"`                         // Please note that once enabled, you still need to [assign permissions](https://dt-url.net/rh22idn \"Dynatrace Generative AI permissions\") to the relevant user groups.
	EnableDocumentSuggestion    *bool                      `json:"enableDocumentSuggestion,omitempty"`    // By enabling document suggestions, Dynatrace AI can find similarities between Problems and existing Notebooks and Dashboards in order to suggest relevant troubleshooting guides. Learn more about [document suggestions](https://dt-url.net/xy02gpo \"Dynatrace AI document suggestions\").
	EnableTenantAwareDataMining *bool                      `json:"enableTenantAwareDataMining,omitempty"` // You can enrich Dynatrace Generative AI with your environment data. This lets you generate more accurate queries that identify and reference relevant entities, events, spans, logs, and metrics from your environment. Once enabled, Dynatrace AI periodically scans your Grail data to create its own semantic index. Please note, it can take up to 24 hours to reflect changes. Learn more about [environment-aware queries](https://dt-url.net/4g42iu7 \"Dynatrace Generative AI environment aware queries\").
}

func (me *Settings) Name() string {
	return "environment"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"blocklist_entries": {
			Type:        schema.TypeList,
			Description: "You can exclude specific data buckets and tables from the semantic index. Learn more about [configuring data access](https://dt-url.net/lc62i1q \"Dynatrace Generative AI data access\").",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Resource{Schema: new(DataminingBlocklistEntries).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"enable_copilot": {
			Type:        schema.TypeBool,
			Description: "Please note that once enabled, you still need to [assign permissions](https://dt-url.net/rh22idn \"Dynatrace Generative AI permissions\") to the relevant user groups.",
			Required:    true,
		},
		"enable_document_suggestion": {
			Type:        schema.TypeBool,
			Description: "By enabling document suggestions, Dynatrace AI can find similarities between Problems and existing Notebooks and Dashboards in order to suggest relevant troubleshooting guides. Learn more about [document suggestions](https://dt-url.net/xy02gpo \"Dynatrace AI document suggestions\").",
			Optional:    true, // precondition
		},
		"enable_tenant_aware_data_mining": {
			Type:        schema.TypeBool,
			Description: "You can enrich Dynatrace Generative AI with your environment data. This lets you generate more accurate queries that identify and reference relevant entities, events, spans, logs, and metrics from your environment. Once enabled, Dynatrace AI periodically scans your Grail data to create its own semantic index. Please note, it can take up to 24 hours to reflect changes. Learn more about [environment-aware queries](https://dt-url.net/4g42iu7 \"Dynatrace Generative AI environment aware queries\").",
			Optional:    true, // precondition
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"blocklist_entries":               me.BlocklistEntries,
		"enable_copilot":                  me.EnableCopilot,
		"enable_document_suggestion":      me.EnableDocumentSuggestion,
		"enable_tenant_aware_data_mining": me.EnableTenantAwareDataMining,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.EnableDocumentSuggestion == nil) && (me.EnableCopilot) {
		me.EnableDocumentSuggestion = opt.NewBool(false)
	}
	if (me.EnableTenantAwareDataMining == nil) && (me.EnableCopilot) {
		me.EnableTenantAwareDataMining = opt.NewBool(false)
	}
	// ---- BlocklistEntries DataminingBlocklistEntries -> {"expectedValue":true,"property":"enableTenantAwareDataMining","type":"EQUALS"}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"blocklist_entries":               &me.BlocklistEntries,
		"enable_copilot":                  &me.EnableCopilot,
		"enable_document_suggestion":      &me.EnableDocumentSuggestion,
		"enable_tenant_aware_data_mining": &me.EnableTenantAwareDataMining,
	})
}
