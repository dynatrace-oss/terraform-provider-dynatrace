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

package segments

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ValidateUUID = func(i any, p cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	if _, err := uuid.Parse(i.(string)); err != nil {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "The value is expected to be a UUID",
			Detail:   fmt.Sprintf("%v is not a UUID", i),
		}
		diags = append(diags, diag)
	}
	return diags
}

var ValidateUUIDOrEmpty = func(i any, p cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	if len(i.(string)) == 0 {
		return diags
	}
	return ValidateUUID(i, p)
}

var ValidateMaxLength = func(maxLength int) schema.SchemaValidateDiagFunc {
	return func(i any, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		if len(i.(string)) > maxLength {
			diag := diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("%s exceeds the maximum length of %d characters", i.(string), maxLength),
			}
			diags = append(diags, diag)
		}
		return diags
	}
}

func isInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

var ValidateTypePossibleValues = func(valuesList []string) schema.SchemaValidateDiagFunc {
	return func(i any, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		if !isInList(i.(string), valuesList) {
			diag := diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("%s is not on the list of possible values: %s", i.(string), valuesList),
			}
			diags = append(diags, diag)
		}
		return diags
	}
}

var ValidateMinLength = func(min int) schema.SchemaValidateDiagFunc {
	return func(i any, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		if len(i.(string)) < min {
			diag := diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("%s must be at least %d characters long", i.(string), min),
			}
			diags = append(diags, diag)
		}
		return diags
	}
}

var ValidateMax = func(max int) schema.SchemaValidateDiagFunc {
	return func(i any, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		var iv int
		switch v := i.(type) {
		case int:
			iv = v
		case string:
			var err error
			v = strings.TrimSpace(v)
			if strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}") {
				return diags
			}
			iv, err = strconv.Atoi(v)
			if err != nil {
				diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: fmt.Sprintf("%s is not a number", i.(string))})
				return diags
			}
		}
		if iv > max {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: fmt.Sprintf("%v must not be greater than %d", i, max)})
		}
		return diags
	}
}

var ValidateRange = func(min int, max int) schema.SchemaValidateDiagFunc {
	return func(i any, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		var iv int
		switch v := i.(type) {
		case int:
			iv = v
		case string:
			var err error
			v = strings.TrimSpace(v)
			if strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}") {
				return diags
			}
			iv, err = strconv.Atoi(v)
			if err != nil {
				diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: fmt.Sprintf("%s is not a number", i.(string))})
				return diags
			}
		}
		if iv < min || iv > max {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: fmt.Sprintf("%v is not in range %d and %d", i, min, max)})
		}
		return diags
	}
}

var ValidateMin = func(min int) schema.SchemaValidateDiagFunc {
	return func(i any, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		var iv int
		switch v := i.(type) {
		case int:
			iv = v
		case string:
			var err error
			v = strings.TrimSpace(v)
			if strings.HasPrefix(v, "{{") && strings.HasSuffix(v, "}}") {
				return diags
			}
			iv, err = strconv.Atoi(v)
			if err != nil {
				diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: fmt.Sprintf("%s is not a number", i.(string))})
				return diags
			}
		}
		if iv < min {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Summary: fmt.Sprintf("%v must not be less than %d", i, min)})
		}
		return diags
	}
}

var ValidateRegex = func(r *regexp.Regexp, errorMessage string) schema.SchemaValidateDiagFunc {
	return func(i any, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		if !r.MatchString(i.(string)) {
			diag := diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("%s contains invalid characters", i.(string)),
				Detail:   errorMessage,
			}
			diags = append(diags, diag)
		}
		return diags
	}
}

var Validate = func(fn ...schema.SchemaValidateDiagFunc) schema.SchemaValidateDiagFunc {
	return func(i any, p cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		for _, validateFunc := range fn {
			diags = append(diags, validateFunc(i, p)...)
		}
		return diags
	}
}
