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

package segments

import (
	"context"
	"encoding/json"

	segmentsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/grail/segments/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	segmentsclient "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/segments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceMulti() *schema.Resource {
	itemSchema := hcl.SetComputedSchema(new(segmentsettings.Segment).Schema())
	itemSchema["id"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The unique identifier of the segment.",
	}

	return &schema.Resource{
		Description: "Retrieve a list of all segments.",
		ReadContext: logging.EnableDSCtx(DataSourceMultipleRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Only return segments with exactly this name",
				Optional:    true,
			},
			"include_ready_made": {
				Type:        schema.TypeBool,
				Description: "If `true`, ready-made segments are also returned. Defaults to `false`",
				Optional:    true,
			},
			"items": {
				Type:        schema.TypeList,
				Description: "An array of the found segments",
				Computed:    true,
				Elem:        &schema.Resource{Schema: itemSchema},
			},
		},
	}
}

func DataSourceMultipleRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := config.ClientSet(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	platformClient, err := clientSet.PlatformClient()
	if err != nil {
		return diag.FromErr(err)
	}

	responses, err := segmentsclient.NewClient(platformClient).GetAll(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	name, filterByName := d.GetOk("name")
	includeReadyMade := d.Get("include_ready_made").(bool)

	items := make([]any, 0, len(responses))
	for _, response := range responses {
		var meta struct {
			UID         string `json:"uid"`
			IsReadyMade bool   `json:"isReadyMade"`
		}
		if err = json.Unmarshal(response.Data, &meta); err != nil {
			return diag.FromErr(err)
		}
		if meta.IsReadyMade && !includeReadyMade {
			continue
		}

		var segment segmentsettings.Segment
		if err = json.Unmarshal(response.Data, &segment); err != nil {
			return diag.FromErr(err)
		}
		if filterByName && segment.Name != name.(string) {
			continue
		}

		item := hcl.Properties{}
		if err = segment.MarshalHCL(item); err != nil {
			return diag.FromErr(err)
		}
		item["id"] = meta.UID
		items = append(items, item)
	}

	if err = d.Set("items", items); err != nil {
		return diag.FromErr(err)
	}

	if filterByName {
		d.SetId(name.(string))
	} else {
		d.SetId("segments")
	}
	return diag.Diagnostics{}
}
