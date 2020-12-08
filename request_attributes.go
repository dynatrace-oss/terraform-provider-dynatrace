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

	api "github.com/dtcookie/dynatrace/api/config"
	"github.com/dtcookie/dynatrace/api/config/requestattributes"
	"github.com/dtcookie/dynatrace/rest"
	"github.com/dtcookie/dynatrace/terraform"
	"github.com/dtcookie/opt"
	"github.com/dtcookie/terraform-provider-dynatrace/config"
	"github.com/dtcookie/terraform-provider-dynatrace/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type reqAtts struct{}

// ResourceRequestAttributes has no documentation
func (ra *reqAtts) Resource() *schema.Resource {
	resource := terraform.ResourceFor(new(requestattributes.RequestAttribute))
	resource.CreateContext = logging.Enable(ra.Create)
	resource.UpdateContext = logging.Enable(ra.Update)
	resource.ReadContext = logging.Enable(ra.Read)
	resource.DeleteContext = logging.Enable(ra.Delete)

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
	instance := new(requestattributes.RequestAttribute)
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
