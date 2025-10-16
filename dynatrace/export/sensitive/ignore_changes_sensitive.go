/**
* @license
* Copyright 2023 Dynatrace LLC
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

package sensitive

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/envutil"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/meta"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Avoid overwriting password that the user has changed manually with known bad data
var IGNORE_CHANGES_REQUIRES_ATTENTION = envutil.GetBoolEnv(envutil.EnvIgnoreChangesRequiresAttention, false)

const SecretValueExact = "${state.secret_value.exact}"
const SecretValue = "${state.secret_value}"

// Replaces all values within the given string slice with `${state.secret_value.exact}`,
// signalling that before handing over attributes to Terraform, the state should get
// searched for the exact same key and that state value get taken for it
func SecretifyExact(values []string) []string {
	for idx := range values {
		values[idx] = SecretValueExact
	}
	return values
}

func ConditionalIgnoreChangesMap(schema map[string]*schema.Schema, itemsToEncode map[string]any) map[string]any {
	return ConditionalIgnoreChangesMapPlus(schema, itemsToEncode, []string{})
}
func ConditionalIgnoreChangesMapPlus(schema map[string]*schema.Schema, itemsToEncode map[string]any, additionalIgnoreFields []string) map[string]any {
	if IGNORE_CHANGES_REQUIRES_ATTENTION {

		lifeCycleIgnoreChanges := buildIgnoreSensitiveFromSchema(schema, additionalIgnoreFields)

		currentlifeCycle, lifeCycleFound := itemsToEncode["lifecycle"]
		if lifeCycleFound {
			fmt.Printf("ERROR: Ignore Changes: lifecycle already exists: %v, overwriting with: %v", currentlifeCycle, lifeCycleIgnoreChanges)
		}

		if len(lifeCycleIgnoreChanges.IgnoreChanges) > 0 {
			itemsToEncode["lifecycle"] = &lifeCycleIgnoreChanges
		}
	}

	return itemsToEncode
}

func ConditionalIgnoreChangesSingle(schema map[string]*schema.Schema, properties *hcl.Properties) error {
	return ConditionalIgnoreChangesSinglePlus(schema, properties, []string{})
}

func ConditionalIgnoreChangesSinglePlus(schema map[string]*schema.Schema, properties *hcl.Properties, additionalIgnoreFields []string) error {
	lifeCycleIgnoreChanges := buildIgnoreSensitiveFromSchema(schema, additionalIgnoreFields)

	if IGNORE_CHANGES_REQUIRES_ATTENTION {
		if err := properties.Encode("lifecycle", &lifeCycleIgnoreChanges); err != nil {
			return err
		}
	}

	return nil
}

func buildIgnoreSensitiveFromSchema(schema map[string]*schema.Schema, additionalIgnoreFields []string) meta.LifeCycle {
	sensitiveFields := []string{}

	sensitiveFields = append(sensitiveFields, additionalIgnoreFields...)

	for schemaKey, schemaField := range schema {
		if schemaField.Sensitive {
			sensitiveFields = appendIfNotExists(sensitiveFields, schemaKey)
		}
	}

	lifeCycleIgnoreChanges := meta.LifeCycle{
		IgnoreChanges: sensitiveFields,
	}
	return lifeCycleIgnoreChanges
}

func appendIfNotExists(slice []string, element string) []string {
	for _, existingElement := range slice {
		if existingElement == element {
			return slice
		}
	}

	return append(slice, element)
}
