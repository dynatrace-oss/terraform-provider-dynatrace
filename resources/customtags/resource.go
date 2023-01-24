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

package customtags

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/common"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/customtags"
	settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/customtags/settings"
	cfg "github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
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
		Schema:        new(settings.Settings).Schema(),
		CreateContext: logging.Enable(Create),
		UpdateContext: logging.Enable(Update),
		ReadContext:   logging.Enable(Read),
		DeleteContext: logging.Enable(Delete),
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

// Create expects the configuration within the given ResourceData and sends it to the Dynatrace Server in order to create that resource
func Create(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	config := new(settings.Settings)
	if err := config.UnmarshalHCL(confighcl.DecoderFrom(d, Resource())); err != nil {
		return diag.FromErr(err)
	}

	if _, err := customtags.Service(cfg.Credentials(m)).Create(config); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(uuid.New().String())

	marshalled := hcl.Properties{}
	if err := config.MarshalHCL(marshalled); err != nil {
		return diag.FromErr(err)
	}
	for k, v := range marshalled {
		d.Set(k, v)
	}

	d.Set("matched_entities", config.MatchedEntities)
	bytes, err := json.Marshal(config)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("current_state", string(bytes))

	return diag.Diagnostics{}
}

// Update expects the configuration within the given ResourceData and send them to the Dynatrace Server in order to update that resource
func Update(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	config := new(settings.Settings)
	if err := config.UnmarshalHCL(confighcl.DecoderFrom(d, Resource())); err != nil {
		return diag.FromErr(err)
	}

	stateDecoder := confighcl.StateDecoderFrom(d, Resource())
	stateConfig := new(settings.Settings)
	if val, ok := stateDecoder.GetOk("current_state"); ok {
		state := val.(string)
		if len(state) > 0 {
			if err := json.Unmarshal([]byte(state), stateConfig); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	var deleteTags []*common.TagFilter
	for _, stateTag := range stateConfig.Tags {
		if stateTag.Value != nil && len(*stateTag.Value) == 0 {
			stateTag.Value = nil
		}
		found := false
		for _, tfTag := range config.Tags {
			if stateTag.Equals(tfTag) {
				found = true
				break
			}
		}
		if !found {
			deleteTags = append(deleteTags, stateTag)
		}
	}

	if len(deleteTags) > 0 {
		delConfig := new(settings.Settings)
		delConfig.EntitySelector = config.EntitySelector
		delConfig.Tags = deleteTags
		srv := customtags.Service(cfg.Credentials(m))
		if fullDeleter, ok := srv.(FullDeleter); ok {
			if err := fullDeleter.DeleteValue(stateConfig); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if err := customtags.Service(cfg.Credentials(m)).Update(config.EntitySelector, config); err != nil {
		return diag.FromErr(err)
	}

	d.Set("matched_entities", config.MatchedEntities)
	bytes, err := json.Marshal(config)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("current_state", string(bytes))

	return diag.Diagnostics{}
}

// Read queries the Dynatrace Server for the configuration
func Read(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	stateDecoder := confighcl.StateDecoderFrom(d, Resource())
	stateConfig := new(settings.Settings)

	var selector string
	if val, ok := stateDecoder.GetOk("entity_selector"); ok {
		selector = val.(string)
	}
	if val, ok := stateDecoder.GetOk("current_state"); ok {
		state := val.(string)
		if len(state) > 0 {
			if err := json.Unmarshal([]byte(state), stateConfig); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	apiConfig := new(settings.Settings)
	if err := customtags.Service(cfg.Credentials(m)).Get(selector, apiConfig); err != nil {
		return diag.FromErr(err)
	}

	filteredTags := common.TagFilters{}
	for _, tag := range apiConfig.Tags {
		for _, stateTag := range stateConfig.Tags {
			if tag.Equals(stateTag) {
				filteredTags = append(filteredTags, stateTag)
				break
			}
		}
	}
	apiConfig.Tags = filteredTags

	marshalled := hcl.Properties{}
	if err := apiConfig.MarshalHCL(marshalled); err != nil {
		return diag.FromErr(err)
	}
	for k, v := range marshalled {
		d.Set(k, v)
	}

	return diag.Diagnostics{}
}

// Delete the configuration
func Delete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	stateDecoder := confighcl.StateDecoderFrom(d, Resource())
	stateConfig := new(settings.Settings)
	if val, ok := stateDecoder.GetOk("entity_selector"); ok {
		stateConfig.EntitySelector = val.(string)
	}
	if val, ok := stateDecoder.GetOk("current_state"); ok {
		state := val.(string)
		if len(state) > 0 {
			if err := json.Unmarshal([]byte(state), stateConfig); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	srv := customtags.Service(cfg.Credentials(m))
	if fullDeleter, ok := srv.(FullDeleter); ok {
		if err := fullDeleter.DeleteValue(stateConfig); err != nil {
			return diag.FromErr(err)
		}
	}

	responseConfig := new(settings.Settings)
	srv.Get(stateConfig.EntitySelector, responseConfig)
	dd, _ := json.Marshal(responseConfig)
	log.Println(string(dd))

	return diag.Diagnostics{}
}

type FullDeleter interface {
	DeleteValue(v *settings.Settings) error
}
