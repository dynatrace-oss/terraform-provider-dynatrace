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

package environments

import (
	"context"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Attribute keys
const (
	// this one is different from attrID inside environments list because "id" is a reserved keyword on root level.
	attrEnvID        = "env_id"
	attrID           = "id"
	attrName         = "name"
	attrActive       = "active"
	attrEnvironments = "environments"
	attrURL          = "url"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Description: "Returns the environments matching the provided filter criteria",
		Schema: map[string]*schema.Schema{
			attrEnvID: {
				Type:        schema.TypeString,
				Description: "Filter by ID of the environment",
				Optional:    true,
			},
			attrName: {
				Type:        schema.TypeString,
				Description: "Filter by friendly name of the environment",
				Optional:    true,
			},
			attrActive: {
				Type:        schema.TypeBool,
				Description: "Filter by active status of the environment",
				Optional:    true,
			},
			attrEnvironments: {
				Type:        schema.TypeList,
				Description: "List of environments matching the filter criteria",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						attrID: {
							Type:        schema.TypeString,
							Description: "The ID of the environment",
							Computed:    true,
						},
						attrName: {
							Type:        schema.TypeString,
							Description: "Friendly name of the environment",
							Computed:    true,
						},
						attrActive: {
							Type:        schema.TypeBool,
							Description: "Property to determine if environment is active",
							Computed:    true,
						},
						attrURL: {
							Type:        schema.TypeString,
							Description: "The URL of the environment",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValIAM)
	if err != nil {
		return diag.FromErr(err)
	}

	service := newEnvironmentService(creds)
	envs, err := service.Get(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get filter values
	filterID, hasFilterID := d.GetOk(attrEnvID)
	filterName, hasFilterName := d.GetOk(attrName)
	filterActive, hasFilterActive := d.GetOk(attrActive)
	// Create a unique ID for the data source based on filter criteria
	dataSourceID := fmt.Sprintf("%#v.%#v.%#v.%#v", creds.IAM.AccountID, filterID, filterName, filterActive)

	// Filter environments based on provided criteria
	var filteredEnvs []map[string]any
	for _, env := range envs {
		// Apply filters
		if hasFilterID && env.ID != filterID.(string) {
			continue
		}
		if hasFilterName && env.Name != filterName.(string) {
			continue
		}
		if hasFilterActive && env.Active != filterActive.(bool) {
			continue
		}

		// Add matching environment to result
		filteredEnvs = append(filteredEnvs, map[string]any{
			attrID:     env.ID,
			attrName:   env.Name,
			attrActive: env.Active,
			attrURL:    env.URL,
		})
	}

	d.SetId(dataSourceID)

	if err := d.Set(attrEnvironments, filteredEnvs); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
