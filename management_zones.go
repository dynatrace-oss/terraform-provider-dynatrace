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

package main

import (
	"context"
	"log"
	"reflect"
	"strings"

	api "github.com/dtcookie/dynatrace/api/config"
	"github.com/dtcookie/dynatrace/api/config/managementzones"
	"github.com/dtcookie/dynatrace/rest"
	"github.com/dtcookie/dynatrace/terraform"
	"github.com/dtcookie/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type managementZones struct {
	diffs map[string]diff
}

func (mzs *managementZones) AttachDiffSuppressFunc(sch *schema.Schema) {
	sch.DiffSuppressFunc = logging.EnableSchemaDiff(func(k, old, new string, d *schema.ResourceData) bool {
		if strings.HasPrefix(k, "metadata.") {
			return true
		}

		if mzs.diffs == nil {
			mzs.diffs = map[string]diff{}
		}
		if strings.HasSuffix(k, ".#") {
			if old == "0" && new == "1" {
				prefix := k[0 : len(k)-1]
				found := false
				for st := range mzs.diffs {
					if strings.HasPrefix(st, prefix) {
						found = true
						break
					}
				}
				if !found {
					return true
				}
			}
		}
		mzs.diffs[k] = diff{old, new}
		return false
	})
	if sch.Elem != nil {
		switch typedSchema := sch.Elem.(type) {
		case *schema.Schema:
			mzs.AttachDiffSuppressFunc(typedSchema)
		case *schema.Resource:
			mzs.AttachDiffSuppressFuncs(typedSchema.Schema)
		}
	}
}

func (mzs *managementZones) AttachDiffSuppressFuncs(schemas map[string]*schema.Schema) {
	if schemas == nil {
		return
	}
	for _, sch := range schemas {
		mzs.AttachDiffSuppressFunc(sch)
	}
}

func (mzs *managementZones) wrap(fn func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		mzs.diffs = map[string]diff{}
		result := fn(ctx, d, m)
		mzs.diffs = map[string]diff{}
		return result
	}
}

// Resource produces terraform resource definition for Management Zones
func (mzs *managementZones) Resource() *schema.Resource {
	resource := terraform.ResourceFor(new(managementzones.ManagementZone))
	resource.CreateContext = logging.Enable(mzs.wrap(mzs.Create))
	resource.UpdateContext = logging.Enable(mzs.wrap(mzs.Update))
	resource.ReadContext = logging.Enable(mzs.wrap(mzs.Read))
	resource.DeleteContext = logging.Enable(mzs.wrap(mzs.Delete))
	mzs.AttachDiffSuppressFuncs(resource.Schema)
	resource.Importer = &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext}
	return resource
}

// Create expects the configuration of a Management Zone within the given ResourceData
// and send them to the Dynatrace Server in order to create that resource
func (mzs *managementZones) Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("ManagementZones.Create(..)")
	}
	var err error
	var resolver terraform.Resolver
	if resolver, err = terraform.NewResolver(d); err != nil {
		return diag.FromErr(err)
	}
	var untypedMgmz interface{}
	if untypedMgmz, err = resolver.Resolve(reflect.TypeOf(managementzones.ManagementZone{})); err != nil {
		return diag.FromErr(err)
	}

	managementZone := untypedMgmz.(managementzones.ManagementZone)
	managementZone.ID = nil
	conf := m.(*config.ProviderConfiguration)
	managementZoneService := managementzones.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	var managementZoneStub *api.EntityShortRepresentation
	if managementZoneStub, err = managementZoneService.Create(&managementZone); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(managementZoneStub.ID)

	return mzs.Read(ctx, d, m)
}

// Update expects the configuration of a Management Zone within the given ResourceData
// and send them to the Dynatrace Server in order to update that resource
func (mzs *managementZones) Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("ManagementZones.Update(..)")
	}
	var err error
	var resolver terraform.Resolver
	if resolver, err = terraform.NewResolver(d); err != nil {
		return diag.FromErr(err)
	}
	var untypedMgmz interface{}
	if untypedMgmz, err = resolver.Resolve(reflect.TypeOf(managementzones.ManagementZone{})); err != nil {
		return diag.FromErr(err)
	}

	managementZone := untypedMgmz.(managementzones.ManagementZone)
	managementZone.ID = opt.NewString(d.Id())
	conf := m.(*config.ProviderConfiguration)
	managementZoneService := managementzones.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	if err = managementZoneService.Update(&managementZone); err != nil {
		return diag.FromErr(err)
	}
	return mzs.Read(ctx, d, m)
}

// Read queries the Dynatrace Server for the configuration of a Management Zone
// identified by the ID within the given ResourceData
func (mzs *managementZones) Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("ManagementZones.Read(..)")
	}
	var err error
	conf := m.(*config.ProviderConfiguration)
	managementZoneService := managementzones.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	managementZone := new(managementzones.ManagementZone)
	if managementZone, err = managementZoneService.Get(d.Id(), false); err != nil {
		return diag.FromErr(err)
	}
	if err = terraform.ToTerraform(managementZone, d); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

// Delete a Management Zone on the Dynatrace Server
func (mzs *managementZones) Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("ManagementZones.Delete(..)")
	}
	var err error
	conf := m.(*config.ProviderConfiguration)
	managementZoneService := managementzones.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	if err = managementZoneService.Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
