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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/logging"

	api "github.com/dtcookie/dynatrace/api/config"
	"github.com/dtcookie/dynatrace/api/config/alertingprofiles"
	"github.com/dtcookie/dynatrace/rest"
	"github.com/dtcookie/dynatrace/terraform"
	"github.com/dtcookie/opt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type alertingProfiles struct {
	diffs map[string]diff
}

func (aps *alertingProfiles) AttachDiffSuppressFunc(sch *schema.Schema) {
	sch.DiffSuppressFunc = logging.EnableSchemaDiff(func(k, old, new string, d *schema.ResourceData) bool {
		if strings.HasPrefix(k, "metadata.") {
			return true
		}
		if aps.diffs == nil {
			aps.diffs = map[string]diff{}
		}
		if strings.HasSuffix(k, ".#") {
			if old == "0" && new == "1" {
				prefix := k[0 : len(k)-1]
				found := false
				for st := range aps.diffs {
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
		aps.diffs[k] = diff{old, new}
		return false
	})
	if sch.Elem != nil {
		switch typedSchema := sch.Elem.(type) {
		case *schema.Schema:
			aps.AttachDiffSuppressFunc(typedSchema)
		case *schema.Resource:
			aps.AttachDiffSuppressFuncs(typedSchema.Schema)
		}
	}
}

func (aps *alertingProfiles) AttachDiffSuppressFuncs(schemas map[string]*schema.Schema) {
	if schemas == nil {
		return
	}
	for _, sch := range schemas {
		aps.AttachDiffSuppressFunc(sch)
	}
}

func (aps *alertingProfiles) wrap(fn func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		aps.diffs = map[string]diff{}
		result := fn(ctx, d, m)
		aps.diffs = map[string]diff{}
		return result
	}
}

// Resource produces terraform resource definition for Management Zones
func (aps *alertingProfiles) Resource() *schema.Resource {
	resource := terraform.ResourceFor(new(alertingprofiles.AlertingProfile))
	resource.CreateContext = logging.Enable(aps.wrap(aps.Create))
	resource.UpdateContext = logging.Enable(aps.wrap(aps.Update))
	resource.ReadContext = logging.Enable(aps.wrap(aps.Read))
	resource.DeleteContext = logging.Enable(aps.wrap(aps.Delete))
	aps.AttachDiffSuppressFuncs(resource.Schema)
	return resource
}

// Create expects the configuration of a Management Zone within the given ResourceData
// and send them to the Dynatrace Server in order to create that resource
func (aps *alertingProfiles) Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("AlertingProfiles.Create(..)")
	}
	var err error
	var resolver terraform.Resolver
	if resolver, err = terraform.NewResolver(d); err != nil {
		return diag.FromErr(err)
	}
	var untypedProfile interface{}
	if untypedProfile, err = resolver.Resolve(reflect.TypeOf(alertingprofiles.AlertingProfile{})); err != nil {
		return diag.FromErr(err)
	}

	alertingprofile := untypedProfile.(alertingprofiles.AlertingProfile)
	alertingprofile.ID = nil
	conf := m.(*config.ProviderConfiguration)
	apiService := alertingprofiles.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	var objStub *api.EntityShortRepresentation
	if objStub, err = apiService.Create(&alertingprofile); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(objStub.ID)

	return aps.Read(ctx, d, m)
}

// Update expects the configuration of a Management Zone within the given ResourceData
// and send them to the Dynatrace Server in order to update that resource
func (aps *alertingProfiles) Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("AlertingProfiles.Update(..)")
	}
	var err error
	var resolver terraform.Resolver
	if resolver, err = terraform.NewResolver(d); err != nil {
		return diag.FromErr(err)
	}
	var untypedProfile interface{}
	if untypedProfile, err = resolver.Resolve(reflect.TypeOf(alertingprofiles.AlertingProfile{})); err != nil {
		return diag.FromErr(err)
	}

	alertingprofile := untypedProfile.(alertingprofiles.AlertingProfile)
	alertingprofile.ID = opt.NewString(d.Id())
	conf := m.(*config.ProviderConfiguration)
	apiService := alertingprofiles.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	if err = apiService.Update(&alertingprofile); err != nil {
		return diag.FromErr(err)
	}
	return aps.Read(ctx, d, m)
}

// Read queries the Dynatrace Server for the configuration of a Management Zone
// identified by the ID within the given ResourceData
func (aps *alertingProfiles) Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("AlertingProfiles.Read(..)")
	}
	var err error
	conf := m.(*config.ProviderConfiguration)
	apiService := alertingprofiles.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	alertingprofile := new(alertingprofiles.AlertingProfile)
	if alertingprofile, err = apiService.Get(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	if err = terraform.ToTerraform(alertingprofile, d); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

// Delete a Management Zone on the Dynatrace Server
func (aps *alertingProfiles) Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("AlertingProfiles.Delete(..)")
	}
	var err error
	conf := m.(*config.ProviderConfiguration)
	apiService := alertingprofiles.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	if err = apiService.Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
