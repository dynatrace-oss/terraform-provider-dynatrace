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

package documents

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

var ValidateIsUUID = func(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

var ValidateNotUUID = func(i any, p cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	if ValidateIsUUID(i.(string)) {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "The value must not be a UUID",
			Detail:   fmt.Sprintf("%v is a UUID", i),
		}
		diags = append(diags, diag)
	}
	return diags
}
