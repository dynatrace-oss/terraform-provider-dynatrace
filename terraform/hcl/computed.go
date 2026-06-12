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

package hcl

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// setComputedItem modifies a provided HCL schema item by setting it to computed and adjusting incompatible properties (MinItems, MaxItems, Optional and Required)
// to ensure the schema remains valid when items are set to computed.
func setComputedItem(item *schema.Schema) *schema.Schema {
	item.MinItems = 0
	item.MaxItems = 0
	item.Optional = false
	item.Required = false
	item.Computed = true
	return item
}

// SetComputedSchema modifies and returns a provided HCL schema by setting all of its items to computed.
// It also recursively modifies nested schemas and takes care of incompatible properties (MinItems, MaxItems, Optional and Required)
// to ensure the schema remains valid when items are set to computed.
// use case: data sources that require Computed to be set
func SetComputedSchema(item map[string]*schema.Schema) map[string]*schema.Schema {
	for k, v := range item {
		item[k] = setComputedItem(v)
		if item[k].Elem != nil {
			if nested, ok := item[k].Elem.(*schema.Resource); ok {
				nested.Schema = SetComputedSchema(nested.Schema)
			}
		}
	}
	return item
}
