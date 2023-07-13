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

package vault

import (
	"fmt"
	"regexp"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"

	vault "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/credentials/vault/settings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: logging.EnableDS(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the credential. Possible values are `CERTIFICATE`, `PUBLIC_CERTIFICATE`, `TOKEN`, `USERNAME_PASSWORD` and `UNKNOWN`. If not specified all credential types will match",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the credential as shown within the Dynatrace WebUI. If not specified all names will match",
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The scope of the credential. Possible values are `ALL`, `EXTENSION` and `SYNTHETIC`. If not specified all scopes will match.",
			},
		},
	}
}

var notFoundRegexp = regexp.MustCompile(`Setting with id '.*' not found \(offline mode\)`)

func DataSourceRead(d *schema.ResourceData, m any) (err error) {
	name := ""
	typ := ""
	scope := ""
	if value, ok := d.GetOk("name"); ok {
		name = value.(string)
	}
	if value, ok := d.GetOk("type"); ok {
		typ = value.(string)
	}
	if value, ok := d.GetOk("scope"); ok {
		scope = value.(string)
	}
	if name == "" && typ == "" && scope == "" {
		return fmt.Errorf("at least one of `name`, `type` or `scope` needs to be specified as a non empty string")
	}

	service := export.Service(config.Credentials(m), export.ResourceTypes.Credentials)
	var stubs api.Stubs
	if stubs, err = service.List(); err != nil {
		return err
	}
	if len(stubs) == 0 {
		d.SetId("")
	}
	for _, stub := range stubs {
		if name != "" && stub.Name != name {
			continue
		}
		var credentials vault.Credentials
		if err = service.Get(stub.ID, &credentials); err != nil {
			/*
				Identically configured credentials are allowed to be configured via REST and WebUI.
				Therefore the block
				  data "dynatrace_credentials" "..." {
				    	name = "..."
				  }
				is not guaranteed to deliver the same ID consistently.
				During normal usage that's a limitation the user needs to be aware of - therefore not that big of an issue.

				But `terraform-provider-dynatrace -export -import-state` will perform the import in offline mode (based on local cache).
				Here the error message `Setting with id '######' not found (offline mode)` hints, that LIST delivered an ID that may exist
				on the remote side, but doesn't locally.

				In that case we just skip to the next stub. There must exist another setting with the same name.
			*/
			if notFoundRegexp.MatchString(err.Error()) {
				continue
			}
			return err
		}
		if scope != "" && string(credentials.Scope) != scope {
			continue
		}
		if typ != "" && string(credentials.Type) != typ {
			continue
		}
		d.Set("scope", string(credentials.Scope))
		d.Set("type", string(credentials.Type))
		d.Set("name", stub.Name)
		d.SetId(stub.ID)
		return nil
	}

	d.SetId("")
	return nil
}
