/*
 * @license
 * Copyright 2026 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ingestsources

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func genericDataSource[T settings.Settings](serviceFn export.ServiceFn[T]) *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(func(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
			return dataSourceRead[T](ctx, d, m, serviceFn)
		}),
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Description: "The display name of the ingest source",
				Optional:    true,
			},
			"path_segment": {
				Type:        schema.TypeString,
				Description: "The path segment of the ingest source",
				Optional:    true,
			},
			"source": {
				Type:        schema.TypeString,
				Description: "The ingest source",
				Optional:    true,
			},
			"source_type": {
				Type:        schema.TypeString,
				Description: "The source type of the ingest source",
				Optional:    true,
			},
			"items": {
				Type:        schema.TypeList,
				Description: "An array of the found configurations",
				Computed:    true,
				Elem:        &schema.Resource{Schema: itemSchema[T]()},
			},
		},
	}
}

func itemSchema[T settings.Settings]() map[string]*schema.Schema {
	scm := hcl.SetComputedSchema(reflect.New(reflect.TypeFor[T]().Elem()).Interface().(settings.Settings).Schema())
	scm["id"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The ID of the ingest source",
	}
	return scm
}

func dataSourceRead[T settings.Settings](ctx context.Context, d *schema.ResourceData, m any, serviceFn export.ServiceFn[T]) diag.Diagnostics {
	clientSet, err := config.ClientSet(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}
	svc, err := serviceFn(clientSet)
	if err != nil {
		return diag.FromErr(err)
	}

	s, ok := svc.(settings.ListFilterService)
	if !ok {
		return diag.FromErr(settings.ErrListFilterNotImplemented(svc.SchemaID()))
	}

	filter := getFilter(d)

	stubs, err := s.ListWithFilter(ctx, filter)
	if err != nil {
		return diag.FromErr(err)
	}
	items := make([]any, 0, len(stubs))

	for _, stub := range stubs {
		val := stub.Value.(T)
		item := hcl.Properties{}
		if err = val.MarshalHCL(item); err != nil {
			return diag.FromErr(err)
		}

		item["id"] = stub.ID
		items = append(items, item)
	}

	if err = d.Set("items", items); err != nil {
		return diag.FromErr(err)
	}

	if filter == "" {
		d.SetId(svc.SchemaID())
	} else {
		d.SetId(filter)
	}

	return diag.Diagnostics{}
}

func getFilter(d *schema.ResourceData) string {
	filters := make([]string, 0, 4)

	if val, ok := d.GetOk("display_name"); ok {
		filters = append(filters, fmt.Sprintf("value.displayName='%s'", settings.EscapeFilterString(val.(string))))
	}

	if val, ok := d.GetOk("path_segment"); ok {
		filters = append(filters, fmt.Sprintf("value.pathSegment='%s'", settings.EscapeFilterString(val.(string))))
	}

	if val, ok := d.GetOk("source"); ok {
		filters = append(filters, fmt.Sprintf("value.source='%s'", settings.EscapeFilterString(val.(string))))
	}

	if val, ok := d.GetOk("source_type"); ok {
		filters = append(filters, fmt.Sprintf("value.sourceType='%s'", settings.EscapeFilterString(val.(string))))
	}

	return strings.Join(filters, " and ")
}
