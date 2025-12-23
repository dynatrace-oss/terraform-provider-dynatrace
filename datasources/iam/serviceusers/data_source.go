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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/serviceusers"
	su "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/serviceusers/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Description: "Fetches IAM service user details by name, ID, or email",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The name of the service user",
				ExactlyOneOf: []string{"name", "id", "email"},
			},
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The UUID of the service user",
				ExactlyOneOf: []string{"name", "id", "email"},
			},
			"email": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The email of the service user",
				ExactlyOneOf: []string{"name", "id", "email"},
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
				Description: "The UUIDs of the groups the service user belongs to",
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValIAM)
	if err != nil {
		return diag.FromErr(err)
	}

	service := serviceusers.Service(creds)
	return dataSourceReadWithService(ctx, d, service)
}

func dataSourceReadWithService(ctx context.Context, d *schema.ResourceData, service settings.CRUDService[*su.ServiceUser]) diag.Diagnostics {
	// Determine lookup method
	var lookupByID string
	var lookupByName string
	var lookupByEmail string

	if v, ok := d.GetOk("id"); ok {
		lookupByID = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		lookupByName = v.(string)
	}
	if v, ok := d.GetOk("email"); ok {
		lookupByEmail = v.(string)
	}

	// If looking up by ID, get directly
	if len(lookupByID) > 0 {
		var serviceUser su.ServiceUser
		if err := service.Get(ctx, lookupByID, &serviceUser); err != nil {
			return diag.FromErr(err)
		}
		return setServiceUserData(d, &serviceUser)
	}

	// For name or email lookup, we need to list and find
	stubs, err := service.List(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, stub := range stubs {
		var serviceUser su.ServiceUser
		if err := service.Get(ctx, stub.ID, &serviceUser); err != nil {
			return diag.FromErr(err)
		}

		if len(lookupByName) > 0 && serviceUser.Name == lookupByName {
			return setServiceUserData(d, &serviceUser)
		}
		if len(lookupByEmail) > 0 && serviceUser.Email == lookupByEmail {
			return setServiceUserData(d, &serviceUser)
		}
	}

	// Not found
	d.SetId("")
	return diag.Errorf("Service user not found")
}

func setServiceUserData(d *schema.ResourceData, serviceUser *su.ServiceUser) diag.Diagnostics {
	d.SetId(serviceUser.UID)
	d.Set("id", serviceUser.UID)
	d.Set("name", serviceUser.Name)
	d.Set("email", serviceUser.Email)
	d.Set("description", serviceUser.Description)

	if len(serviceUser.Groups) > 0 {
		sort.Strings(serviceUser.Groups)
		d.Set("groups", serviceUser.Groups)
	} else {
		d.Set("groups", []string{})
	}

	return diag.Diagnostics{}
}
