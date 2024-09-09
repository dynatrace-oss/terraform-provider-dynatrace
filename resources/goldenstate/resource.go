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
	"os"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const Debug = true

var Enabled = os.Getenv("DYNATRACE_GOLDEN_STATE_ENABLED") == "true"

func Resource() *schema.Resource {
	schemaMap := map[string]*schema.Schema{
		"mode": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "WARN",
			Description:  "Possible values are:\n* `DELETE` if you want resources to automatally get deleted`n* `WARN` if you want to get notified about resources that aren't managed by Terraform via a warning message from this resource`\nDefault is `WARN`.",
			ValidateFunc: validation.StringInSlice([]string{"DELETE", "WARN"}, false),
		},
	}
	for resource := range serviceMap {
		schemaMap[string(resource)] = &schema.Schema{
			Type:        schema.TypeSet,
			Description: fmt.Sprintf("The IDs for resource of type `%s` this `dynatrace_golden_state` should ignore (and therefore neither warn about their existence nor attempt to delete them). Specify `[]` if you expect no such resources to exist in Dynatrace. Omit this attribute if you don't care about these kinds of resources regarding the golden state of the environment.", resource),
			Optional:    true,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		}
	}
	return &schema.Resource{
		Schema:        schemaMap,
		CreateContext: logging.Enable(Create),
		UpdateContext: logging.Enable(Update),
		ReadContext:   logging.Enable(Read),
		DeleteContext: logging.Enable(Delete),
	}
}

func Delete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	d.SetId("")
	if !Enabled {
		return diag.Diagnostics{diag.Diagnostic{Severity: diag.Warning, Summary: DisabledMessage}}
	}
	return diag.Diagnostics{}
}
