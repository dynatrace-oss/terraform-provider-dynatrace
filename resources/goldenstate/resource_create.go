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

package goldenstate

import (
	"context"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	cfg "github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/confighcl"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Create(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	// logging.File.Println("-- CREATE --")
	d.SetId(uuid.NewString())
	return update(ctx, d, m, "  ")
}

func Update(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	// logging.File.Println("-- UPDATE --")
	return update(ctx, d, m, "  ")
}

func update(ctx context.Context, d *schema.ResourceData, m any, indent string) diag.Diagnostics {
	allDiags := diag.Diagnostics{}
	creds, _ := cfg.Credentials(m, cfg.CredValNone)
	for key := range serviceMap {
		diags, err := CommonUpdate(ctx, d, key, creds, "  ")
		if len(diags) > 0 {
			details := ""
			for _, diagElem := range diags {
				details = fmt.Sprintf("%s\n%s", details, diagElem.Summary)
			}
			allDiags = append(allDiags, diag.Diagnostic{Severity: diag.Warning, Summary: fmt.Sprintf("There exist resources of type `%s` not managed by terraform", key), Detail: strings.TrimSpace(details)})
		}
		if err != nil {
			return append(allDiags, diag.FromErr(err)...)
		}
	}
	return append(allDiags, read(ctx, d, m, indent+"  ")...)
}

func CommonUpdate(ctx context.Context, d *schema.ResourceData, key export.ResourceType, creds *settings.Credentials, indent string) (diag.Diagnostics, error) {
	mode := "QUIET"
	utMode := d.Get("mode")
	if tMode, ok := utMode.(string); ok {
		if tMode == "WARN" || tMode == "DELETE" {
			mode = tMode
		} else {
			mode = "QUIET"
		}
	}
	// logging.File.Println(indent, "CommonUpdate", key)
	cd := confighcl.DecoderFrom(d, Resource())
	untypedIDs, ok := cd.GetOk(string(key))
	// logging.File.Println("  cd.untypedIDs:", untypedIDs)
	// logging.File.Println("  cd.ok:", ok)
	if !ok {
		return diag.Diagnostics{}, nil
	}
	service := serviceMap[key](creds)

	stubs, err := service.List(ctx)

	if err != nil {
		return diag.Diagnostics{}, err
	}
	existingIDs := map[string]string{}
	for _, stub := range stubs {
		existingIDs[stub.ID] = stub.Name
	}
	if idSet, ok := untypedIDs.(*schema.Set); ok {
		for _, untypedID := range idSet.List() {
			if id, ok := untypedID.(string); ok {
				if matches := regexpNameId.FindStringSubmatch(id); len(matches) == 3 {
					id = matches[2]
				}
				delete(existingIDs, id)
			}
		}
	}

	diags := diag.Diagnostics{}

	for id, name := range existingIDs {
		switch mode {
		case "QUIET":
		case "WARN":
			diags = append(diags, diag.Diagnostic{Summary: fmt.Sprintf("[ %-24s ] %s", trimName(name), id)})
		case "DELETE":
			if Debug {
				logging.File.Printf("[%s] DELETING (%s) %s", key, name, id)
			}
			// if err := service.Delete(ctx, id); err != nil {
			// 	return err
			// }
		}
	}
	return diags, nil
}
