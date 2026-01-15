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

package serviceusers

import (
	"context"
	"errors"
	"fmt"
	"slices"
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
	id, serviceUser, err := findServiceUser(ctx, d, service)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(id)
	if err := d.Set("name", serviceUser.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email", serviceUser.Email); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", serviceUser.Description); err != nil {
		return diag.FromErr(err)
	}
	groups := slices.Clone(serviceUser.Groups)
	sort.Strings(groups)
	if err := d.Set("groups", groups); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func findServiceUser(ctx context.Context, d *schema.ResourceData, service settings.CRUDService[*su.ServiceUser]) (string, *su.ServiceUser, error) {
	if id, ok := tryGetStringByKey(d, "id"); ok {
		return findServiceUserByID(ctx, id, service)
	} else if name, ok := tryGetStringByKey(d, "name"); ok {
		return findServiceUserByName(ctx, name, service)
	} else if email, ok := tryGetStringByKey(d, "email"); ok {
		return findServiceUserByEmail(ctx, email, service)
	}

	return "", nil, errors.New("either id, name, or email must be provided for lookup")
}

func findServiceUserByID(ctx context.Context, id string, service settings.CRUDService[*su.ServiceUser]) (string, *su.ServiceUser, error) {
	var serviceUser su.ServiceUser
	if err := service.Get(ctx, id, &serviceUser); err != nil {
		return "", nil, err
	}

	return id, &serviceUser, nil
}

func findServiceUserByName(ctx context.Context, name string, service settings.CRUDService[*su.ServiceUser]) (string, *su.ServiceUser, error) {
	stubs, err := service.List(ctx)
	if err != nil {
		return "", nil, err
	}

	var foundID string
	var foundServiceUser su.ServiceUser
	for _, stub := range stubs {
		var serviceUser su.ServiceUser
		if err := service.Get(ctx, stub.ID, &serviceUser); err != nil {
			return "", nil, err
		}

		if serviceUser.Name == name {
			if foundID != "" {
				return "", nil, fmt.Errorf("multiple service users found with name '%s'", name)
			}
			foundID = stub.ID
			foundServiceUser = serviceUser
		}
	}

	if foundID == "" {
		return "", nil, fmt.Errorf("no service user found with name '%s'", name)
	}

	return foundID, &foundServiceUser, nil
}

func findServiceUserByEmail(ctx context.Context, email string, service settings.CRUDService[*su.ServiceUser]) (string, *su.ServiceUser, error) {
	stubs, err := service.List(ctx)
	if err != nil {
		return "", nil, err
	}

	for _, stub := range stubs {
		var serviceUser su.ServiceUser
		if err := service.Get(ctx, stub.ID, &serviceUser); err != nil {
			return "", nil, err
		}

		if serviceUser.Email == email {
			return stub.ID, &serviceUser, nil
		}
	}

	return "", nil, fmt.Errorf("no service user found with email '%s'", email)
}

func tryGetStringByKey(d *schema.ResourceData, key string) (string, bool) {
	v, ok := d.GetOk(key)
	if !ok {
		return "", false
	}

	vStr, ok := v.(string)
	return vStr, ok
}
