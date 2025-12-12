/**
* @license
* Copyright 2025 Dynatrace LLC
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

package serviceusers

import (
	"context"
	"sort"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/serviceusers"
	svcusr "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/serviceusers/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Description: "Fetches IAM service user details by name or ID",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
				Description:  "The name of the service user",
			},
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
				Description:  "The ID of the service user",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the service user",
			},
			"groups": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "The groups assigned to the service user",
			},
		},
	}
}

var listMutex sync.Mutex
var listStubs api.Stubs
var lastServiceUsersServiceRevision = "-"

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	// we're working with a cached list of stubs/IDs here
	// therefore we're restricting access to that method
	// to just one at a time
	// ... in case we need to refresh that cache
	listMutex.Lock()
	defer listMutex.Unlock()

	var name string
	var id string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if v, ok := d.GetOk("id"); ok {
		id = v.(string)
	}

	creds, err := config.Credentials(m, config.CredValIAM)
	if err != nil {
		return diag.FromErr(err)
	}

	service := serviceusers.Service(creds)

	// If we have an ID, fetch directly
	if len(id) > 0 {
		var serviceUser svcusr.ServiceUser
		if err := service.Get(ctx, id, &serviceUser); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(id)
		d.Set("name", serviceUser.Name)
		if len(serviceUser.Description) > 0 {
			d.Set("description", serviceUser.Description)
		}
		if len(serviceUser.Groups) > 0 {
			sort.Strings(serviceUser.Groups)
			d.Set("groups", serviceUser.Groups)
		} else {
			d.Set("groups", []string{})
		}
		return diag.Diagnostics{}
	}

	// If we have a name, fetch list and search
	if len(name) > 0 {
		var stubs api.Stubs
		// Service users service updates the revision string every time
		// a CREATE or DELETE happens
		// In that case we need to fetch the currently available service users IDs anew
		if serviceusers.GetRevision() != lastServiceUsersServiceRevision {
			lastServiceUsersServiceRevision = serviceusers.GetRevision()
			var err error
			stubs, err = service.List(ctx)
			if err != nil {
				return diag.FromErr(err)
			}
			listStubs = stubs
		} else {
			stubs = listStubs
		}

		if len(stubs) > 0 {
			for _, stub := range stubs {
				if name == stub.Name {
					d.SetId(stub.ID)
					// Fetch full details
					var serviceUser svcusr.ServiceUser
					if err := service.Get(ctx, stub.ID, &serviceUser); err != nil {
						return diag.FromErr(err)
					}
					d.Set("name", serviceUser.Name)
					if len(serviceUser.Description) > 0 {
						d.Set("description", serviceUser.Description)
					}
					if len(serviceUser.Groups) > 0 {
						sort.Strings(serviceUser.Groups)
						d.Set("groups", serviceUser.Groups)
					} else {
						d.Set("groups", []string{})
					}
					return diag.Diagnostics{}
				}
			}
		}
	}

	d.SetId("")
	return diag.Diagnostics{}
}
