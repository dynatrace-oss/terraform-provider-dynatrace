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

package autotagrules

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/tags/autotagging"
	auto_tag_rule_settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/tags/autotagging/rules/settings"
	auto_tag_settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/tags/autotagging/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/confighcl"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resource produces terraform resource definition for Management Zones
func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        new(auto_tag_rule_settings.Settings).Schema(),
		CreateContext: logging.Enable(Create),
		UpdateContext: logging.Enable(Update),
		ReadContext:   logging.Enable(Read),
		DeleteContext: logging.Enable(Delete),
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

// Create expects the configuration within the given ResourceData and sends it to the Dynatrace Server in order to create that resource
func Create(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve auto tag rules from TF file
	tfConfig := new(auto_tag_rule_settings.Settings)
	if err := tfConfig.UnmarshalHCL(confighcl.DecoderFrom(d, Resource())); err != nil {
		return diag.FromErr(err)
	}
	// Retrieve auto tag rules from API
	apiConfig := new(auto_tag_settings.Settings)
	if err := autotagging.Service(creds).Get(tfConfig.AutoTagId, apiConfig); err != nil {
		return diag.FromErr(err)
	}

	// Concat API and TF rules
	// Concatenated rules may contain duplicates but API does not care
	apiConfig.Rules = append(apiConfig.Rules, tfConfig.Rules...)

	if err := autotagging.Service(creds).Update(tfConfig.AutoTagId, apiConfig); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(uuid.New().String())

	marshalled := hcl.Properties{}
	if err := tfConfig.MarshalHCL(marshalled); err != nil {
		return diag.FromErr(err)
	}
	for k, v := range marshalled {
		d.Set(k, v)
	}

	bytes, err := json.Marshal(tfConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("current_state", string(bytes))

	return diag.Diagnostics{}
}

// Update expects the configuration within the given ResourceData and send them to the Dynatrace Server in order to update that resource
func Update(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve auto tag rules from state
	stateDecoder := confighcl.StateDecoderFrom(d, Resource())
	stateConfig := new(auto_tag_rule_settings.Settings)
	if val, ok := stateDecoder.GetOk("current_state"); ok {
		state := val.(string)
		if len(state) > 0 {
			if err := json.Unmarshal([]byte(state), stateConfig); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	// Retrieve auto tag rules from TF file
	tfConfig := new(auto_tag_rule_settings.Settings)
	if err := tfConfig.UnmarshalHCL(confighcl.DecoderFrom(d, Resource())); err != nil {
		return diag.FromErr(err)
	}
	// Retrieve auto tag rules from API
	apiConfig := new(auto_tag_settings.Settings)
	if err := autotagging.Service(creds).Get(tfConfig.AutoTagId, apiConfig); err != nil {
		return diag.FromErr(err)
	}

	// Find rules only in state based off of match above to mark for deletion
	deleteRules := auto_tag_settings.Rules{}
	for _, stateRule := range stateConfig.Rules {
		found := false
		for _, tfRule := range tfConfig.Rules {
			if reflect.DeepEqual(*stateRule, *tfRule) {
				found = true
				break
			}
		}
		if !found {
			deleteRules = append(deleteRules, stateRule)
		}
	}
	// Concat API and TF rules
	apiConfig.Rules = append(apiConfig.Rules, tfConfig.Rules...)
	// Take rules of API and TF, remove rules marked for deletion
	finalRules := auto_tag_settings.Rules{}
	for _, apiRule := range apiConfig.Rules {
		found := false
		for _, deleteRule := range deleteRules {
			if reflect.DeepEqual(*apiRule, *deleteRule) {
				found = true
				break
			}
		}
		if !found {
			finalRules = append(finalRules, apiRule)
		}
	}
	apiConfig.Rules = finalRules

	if err := autotagging.Service(creds).Update(tfConfig.AutoTagId, apiConfig); err != nil {
		return diag.FromErr(err)
	}

	bytes, err := json.Marshal(tfConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("current_state", string(bytes))

	return diag.Diagnostics{}
}

// Read queries the Dynatrace Server for the configuration
func Read(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve auto tag rules from state
	stateDecoder := confighcl.StateDecoderFrom(d, Resource())
	stateConfig := new(auto_tag_rule_settings.Settings)
	if val, ok := stateDecoder.GetOk("auto_tag_id"); ok {
		stateConfig.AutoTagId = val.(string)
	}
	if val, ok := stateDecoder.GetOk("current_state"); ok {
		state := val.(string)
		if len(state) > 0 {
			if err := json.Unmarshal([]byte(state), stateConfig); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	// Retrieve auto tag rules from API
	apiConfig := new(auto_tag_settings.Settings)
	if err := autotagging.Service(creds).Get(stateConfig.AutoTagId, apiConfig); err != nil {
		return diag.FromErr(err)
	}

	// Find matching rules with state and API
	matchingRules := auto_tag_settings.Rules{}
	for _, stateRule := range stateConfig.Rules {
		for _, apiRule := range apiConfig.Rules {
			if reflect.DeepEqual(*stateRule, *apiRule) {
				matchingRules = append(matchingRules, stateRule)
			}
		}
	}
	stateConfig.Rules = matchingRules

	marshalled := hcl.Properties{}
	if err := stateConfig.MarshalHCL(marshalled); err != nil {
		return diag.FromErr(err)
	}
	for k, v := range marshalled {
		d.Set(k, v)
	}

	return diag.Diagnostics{}
}

// Delete the configuration
func Delete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	// Retrieve auto tag rules from state
	stateDecoder := confighcl.StateDecoderFrom(d, Resource())
	stateConfig := new(auto_tag_rule_settings.Settings)
	if val, ok := stateDecoder.GetOk("auto_tag_id"); ok {
		stateConfig.AutoTagId = val.(string)
	}
	if val, ok := stateDecoder.GetOk("current_state"); ok {
		state := val.(string)
		if len(state) > 0 {
			if err := json.Unmarshal([]byte(state), stateConfig); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	// Retrieve auto tag rules from API
	apiConfig := new(auto_tag_settings.Settings)
	if err := autotagging.Service(creds).Get(stateConfig.AutoTagId, apiConfig); err != nil {
		return diag.FromErr(err)
	}

	// If matching rule is found from the state, remove from API
	for _, stateRule := range stateConfig.Rules {
		for i, apiRule := range apiConfig.Rules {
			if reflect.DeepEqual(*stateRule, *apiRule) {
				apiConfig.Rules = append(apiConfig.Rules[:i], apiConfig.Rules[i+1:]...)
			}
		}
	}

	if err := autotagging.Service(creds).Update(stateConfig.AutoTagId, apiConfig); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
