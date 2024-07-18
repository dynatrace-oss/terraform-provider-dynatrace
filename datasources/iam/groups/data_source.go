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

package groups

import (
	"context"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/groups"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

var listMutex sync.Mutex
var listStubs api.Stubs
var lastGroupsServiceRevision = "-"

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	// we're working with a cached list of stubs/IDs here
	// therefore we're restricting access to that method
	// to just one at a time
	// ... in case we need to refresh that cache
	listMutex.Lock()
	defer listMutex.Unlock()

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	creds, err := config.Credentials(m, config.CredValIAM)
	if err != nil {
		return diag.FromErr(err)
	}

	var stubs api.Stubs
	// Groups Service updates the revision string every time
	// a CREATE or DELETE happens
	// In that case we need to fetch the currently available groups IDs anew
	if groups.GetRevision() != lastGroupsServiceRevision {
		lastGroupsServiceRevision = groups.GetRevision()
		var err error
		service := groups.Service(creds)
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
				return diag.Diagnostics{}
			}
		}
	}
	d.SetId("")
	return diag.Diagnostics{}
}
