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

package provider

import (
	"context"
	"log"
	"reflect"
	"strings"

	api "github.com/dtcookie/dynatrace/api/config"
	"github.com/dtcookie/dynatrace/api/config/dashboards"
	"github.com/dtcookie/dynatrace/rest"
	"github.com/dtcookie/dynatrace/terraform"
	"github.com/dtcookie/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type diff struct {
	old string
	new string
}

type resDashboards struct {
	diffs map[string]diff
}

func (db *resDashboards) AttachDiffSuppressFunc(sch *schema.Schema) {
	sch.DiffSuppressFunc = logging.EnableSchemaDiff(func(k, old, new string, d *schema.ResourceData) bool {
		if strings.Contains(k, ".sharing_details.") {
			return true
		}
		if strings.HasPrefix(k, "metadata.") {
			return true
		}
		if db.diffs == nil {
			db.diffs = map[string]diff{}
		}
		if strings.HasSuffix(k, ".#") {
			if old == "0" && new == "1" {
				prefix := k[0 : len(k)-1]
				found := false
				for st := range db.diffs {
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
		db.diffs[k] = diff{old, new}
		return false
	})
	if sch.Elem != nil {
		switch typedSchema := sch.Elem.(type) {
		case *schema.Schema:
			db.AttachDiffSuppressFunc(typedSchema)
		case *schema.Resource:
			db.AttachDiffSuppressFuncs(typedSchema.Schema)
		}
	}
}

func (db *resDashboards) AttachDiffSuppressFuncs(schemas map[string]*schema.Schema) {
	if schemas == nil {
		return
	}
	for _, sch := range schemas {
		db.AttachDiffSuppressFunc(sch)
	}
}

func (db *resDashboards) wrap(fn func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		db.diffs = map[string]diff{}
		result := fn(ctx, d, m)
		db.diffs = map[string]diff{}
		return result
	}
}

// ResourceDashboards produces resource definitions for Dashboards
func (db *resDashboards) Resource() *schema.Resource {
	resource := terraform.ResourceFor(new(dashboards.Dashboard))
	resource.CreateContext = logging.Enable(db.wrap(db.Create))
	resource.UpdateContext = logging.Enable(db.wrap(db.Update))
	resource.ReadContext = logging.Enable(db.wrap(db.Read))
	resource.DeleteContext = logging.Enable(db.wrap(db.Delete))
	resource.Importer = &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext}
	db.AttachDiffSuppressFuncs(resource.Schema)
	return resource
}

// Create expects configuration for a Dashboard within the given ResourceData
// and sends it via POST request to the Dynatrace Server in order to create that resource.
func (db *resDashboards) Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("Dashboards.Create(...)")
	}

	var err error
	var resolver terraform.Resolver
	if resolver, err = terraform.NewResolver(d); err != nil {
		return diag.FromErr(err)
	}
	var untypedDashboard interface{}
	if untypedDashboard, err = resolver.Resolve(reflect.TypeOf(dashboards.Dashboard{})); err != nil {
		return diag.FromErr(err)
	}

	dashboard := untypedDashboard.(dashboards.Dashboard)
	dashboard.ID = nil
	dashboard.ConfigurationMetadata = nil
	conf := m.(*config.ProviderConfiguration)
	dashboardService := dashboards.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	var dashboardStub *api.EntityShortRepresentation
	if dashboardStub, err = dashboardService.Create(&dashboard); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dashboardStub.ID)

	return db.Read(ctx, d, m)
}

// Update expects configuration for a Dashboard within the given ResourceData
// and sends it via POST request to the Dynatrace Server in order to update that resource
func (db *resDashboards) Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("Dashboards.Update(...)")
	}
	var err error
	var resolver terraform.Resolver
	if resolver, err = terraform.NewResolver(d); err != nil {
		return diag.FromErr(err)
	}
	var untypedDashboard interface{}
	if untypedDashboard, err = resolver.Resolve(reflect.TypeOf(dashboards.Dashboard{})); err != nil {
		return diag.FromErr(err)
	}

	dashboard := untypedDashboard.(dashboards.Dashboard)
	dashboard.ConfigurationMetadata = nil
	dashboard.ID = opt.NewString(d.Id())
	conf := m.(*config.ProviderConfiguration)
	dashboardService := dashboards.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	if err = dashboardService.Update(&dashboard); err != nil {
		return diag.FromErr(err)
	}
	return db.Read(ctx, d, m)
}

// Read queries the Dynatrace REST API for the configuration of a Dashboard
func (db *resDashboards) Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("Dashboards.Read(...)")
	}
	var err error
	conf := m.(*config.ProviderConfiguration)
	dashboardService := dashboards.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	var dashboard *dashboards.Dashboard
	if dashboard, err = dashboardService.Get(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	dashboard.ConfigurationMetadata = nil
	if err = terraform.ToTerraform(dashboard, d); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

// Delete deletes the dashboard with the given ID
func (db *resDashboards) Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("Dashboards.Delete(...)")
	}

	var err error
	conf := m.(*config.ProviderConfiguration)
	dashboardService := dashboards.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	if err = dashboardService.Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
