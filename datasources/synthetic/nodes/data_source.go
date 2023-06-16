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

package nodes

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	nodes "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/nodes"
	nodessettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/nodes/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: logging.EnableDS(DataSourceRead),
		Schema: map[string]*schema.Schema{
			"nodes": {
				Type:     schema.TypeList,
				Elem:     &schema.Resource{Schema: new(nodessettings.Settings).Schema()},
				Computed: true,
			},
		},
	}
}

func DataSourceRead(d *schema.ResourceData, m any) (err error) {
	var stubs api.Stubs
	if stubs, err = nodes.Service(config.Credentials(m)).List(); err != nil {
		return err
	}
	nodes := []*nodessettings.Settings{}
	for _, stub := range stubs {
		value := stub.Value.(nodessettings.Settings)
		nodes = append(nodes, &value)
	}
	marshalled := hcl.Properties{}
	if marshalled.EncodeSlice("nodes", nodes); err != nil {
		return err
	}
	d.Set("nodes", marshalled["nodes"])
	d.SetId("DYNATRACE_SYNTHETIC_NODES")
	return nil
}
