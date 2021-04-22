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
	"github.com/dtcookie/dynatrace/api/config/requestattributes"
	"github.com/dtcookie/dynatrace/rest"
	"github.com/dtcookie/dynatrace/terraform"
	"github.com/dtcookie/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type reqAtts struct {
	diffs map[string]diff
}

func (ra *reqAtts) AttachDiffSuppressFunc(sch *schema.Schema) {
	sch.DiffSuppressFunc = logging.EnableSchemaDiff(func(k, old, new string, d *schema.ResourceData) bool {
		if strings.HasPrefix(k, "metadata.") {
			return true
		}

		if ra.diffs == nil {
			ra.diffs = map[string]diff{}
		}
		if strings.HasSuffix(k, ".#") {
			if old == "0" && new == "1" {
				prefix := k[0 : len(k)-1]
				found := false
				for st := range ra.diffs {
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
		ra.diffs[k] = diff{old, new}
		return false
	})
	if sch.Elem != nil {
		switch typedSchema := sch.Elem.(type) {
		case *schema.Schema:
			ra.AttachDiffSuppressFunc(typedSchema)
		case *schema.Resource:
			ra.AttachDiffSuppressFuncs(typedSchema.Schema)
		}
	}
}

func (ra *reqAtts) AttachDiffSuppressFuncs(schemas map[string]*schema.Schema) {
	if schemas == nil {
		return
	}
	for _, sch := range schemas {
		ra.AttachDiffSuppressFunc(sch)
	}
}

func (ra *reqAtts) wrap(fn func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		ra.diffs = map[string]diff{}
		result := fn(ctx, d, m)
		ra.diffs = map[string]diff{}
		return result
	}
}

// ResourceRequestAttributes has no documentation
func (ra *reqAtts) Resource() *schema.Resource {
	resource := terraform.ResourceFor(new(requestattributes.RequestAttribute))
	resource.CreateContext = logging.Enable(ra.wrap(ra.Create))
	resource.UpdateContext = logging.Enable(ra.wrap(ra.Update))
	resource.ReadContext = logging.Enable(ra.wrap(ra.Read))
	resource.DeleteContext = logging.Enable(ra.wrap(ra.Delete))
	ra.AttachDiffSuppressFuncs(resource.Schema)
	resource.Importer = &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext}
	return resource
}

// FooError has no documentation
type FooError struct {
	Message string
}

func (fe *FooError) Error() string {
	return fe.Message
}

// CreateRequestAttribute has no documentation
func (ra *reqAtts) Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("RequestAttributes.Create(..)")
	}
	var err error
	var resolver terraform.Resolver
	if resolver, err = terraform.NewResolver(d); err != nil {
		return diag.FromErr(err)
	}
	var untypedInstance interface{}
	if untypedInstance, err = resolver.Resolve(reflect.TypeOf(requestattributes.RequestAttribute{})); err != nil {
		return []diag.Diagnostic{
			{
				Severity: diag.Error,
				Summary:  "Unable to resolve resource data to a well formed Request Attribute",
				Detail:   err.Error(),
			},
		}
	}

	instance := untypedInstance.(requestattributes.RequestAttribute)
	instance.ID = nil
	conf := m.(*config.ProviderConfiguration)
	apiService := requestattributes.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	var stub *api.EntityShortRepresentation
	if stub, err = apiService.Create(&instance); err != nil {
		return []diag.Diagnostic{
			{
				Severity: diag.Error,
				Summary:  "Dynatrace REST API rejected payload",
				Detail:   err.Error(),
			},
		}
	}
	d.SetId(stub.ID)

	return ra.Read(ctx, d, m)
}

// UpdateRequestAttribute has no documentation
func (ra *reqAtts) Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("RequestAttributes.Update(..)")
	}
	var err error
	var resolver terraform.Resolver
	if resolver, err = terraform.NewResolver(d); err != nil {
		return diag.FromErr(err)
	}
	var untypedInstance interface{}
	if untypedInstance, err = resolver.Resolve(reflect.TypeOf(requestattributes.RequestAttribute{})); err != nil {
		return diag.FromErr(err)
	}

	instance := untypedInstance.(requestattributes.RequestAttribute)
	instance.ID = opt.NewString(d.Id())
	conf := m.(*config.ProviderConfiguration)
	apiService := requestattributes.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	if err = apiService.Update(&instance); err != nil {
		return diag.FromErr(err)
	}
	return ra.Read(ctx, d, m)
}

// Read has no documentation
func (ra *reqAtts) Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("RequestAttributes.Read(..)")
	}
	var err error
	conf := m.(*config.ProviderConfiguration)
	apiService := requestattributes.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	var instance *requestattributes.RequestAttribute
	if instance, err = apiService.Get(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	if err = terraform.ToTerraform(instance, d); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

// Delete has no documentation
func (ra *reqAtts) Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if config.HTTPVerbose {
		log.Println("RequestAttributes.Delete(..)")
	}
	var err error
	conf := m.(*config.ProviderConfiguration)
	apiService := requestattributes.NewService(conf.DTenvURL, conf.APIToken)
	rest.Verbose = config.HTTPVerbose
	if err = apiService.Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
