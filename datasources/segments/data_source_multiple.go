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

	segmentservice "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/grail/segments"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceMulti() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve a list of all segments.",
		ReadContext: logging.EnableDSCtx(DataSourceMultipleRead),
		Schema: map[string]*schema.Schema{
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of the segment.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the segment.",
						},
					},
				},
			},
		},
	}
}

func DataSourceMultipleRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientSet, err := config.ClientSet(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	service, err := segmentservice.Service(clientSet)
	if err != nil {
		return diag.FromErr(err)
	}

	stubs, err := service.List(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	values := make([]map[string]any, 0, len(stubs))
	for _, stub := range stubs {
		values = append(values, map[string]any{"id": stub.ID, "name": stub.Name})
	}

	if err = d.Set("values", values); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("segments")
	return diag.Diagnostics{}
}
