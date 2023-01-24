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

package locations

import (
	locations "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/locations"
	locsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/locations/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func UniqueDataSource() *schema.Resource {
	return &schema.Resource{
		Read:   UniqueDataSourceRead,
		Schema: new(locsettings.SyntheticLocation).Schema(),
	}
}

func UniqueDataSourceRead(d *schema.ResourceData, m any) (err error) {
	var id *string
	var name *string
	var typeLoc *string
	var status *string
	var stage *string
	var cloudPlatform *string
	var ips []string

	if v, ok := d.GetOk("entity_id"); ok {
		d.SetId(v.(string))
		id = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		name = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("type"); ok {
		typeLoc = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("status"); ok {
		status = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("stage"); ok {
		stage = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("cloud_platform"); ok {
		cloudPlatform = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("ips"); ok {
		if vt, ok := v.([]string); ok {
			ips = vt
		}
	}
	var stubs settings.Stubs
	if stubs, err = locations.Service(config.Credentials(m)).List(); err != nil {
		return err
	}

	for _, stub := range stubs {
		if id != nil {
			if *id != stub.ID {
				continue
			}
		}
		if name != nil {
			if *name != stub.Name {
				continue
			}
		}
		value := stub.Value.(*locsettings.SyntheticLocation)
		if typeLoc != nil {
			if *typeLoc != string(value.Type) {
				continue
			}
		}
		if status != nil {
			if *status != string(*value.Status) {
				continue
			}
		}
		if stage != nil {
			if *stage != string(*value.Stage) {
				continue
			}
		}
		if cloudPlatform != nil {
			if *cloudPlatform != string(*value.CloudPlatform) {
				continue
			}
		}
		if len(ips) > 0 {
			if !subsetCheck(value.IPs, ips) {
				continue
			}
		}

		marshalled := hcl.Properties{}
		if err := value.MarshalHCL(marshalled); err != nil {
			return err
		}
		for k, v := range marshalled {
			if k == "entity_id" {
				d.SetId(v.(string))
			}
			d.Set(k, v)
		}
		return nil
	}

	d.SetId("")
	return nil
}

// subsetCheck verifies that the input strings are a subset of source strings
// Arguments: source slice of strings, input slice of strings
// Return: true if subset, false if not
func subsetCheck(source []string, input []string) bool {
	for _, inputString := range input {
		found := false
		for _, sourceString := range source {
			if inputString == sourceString {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
