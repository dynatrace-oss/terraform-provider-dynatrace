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


package customservice

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getBool(data *schema.ResourceData, key string) bool {
	if data == nil {
		return false
	}
	value := data.Get(key)
	if value == nil {
		return false
	}
	if bVal, ok := value.(bool); ok {
		return bVal
	}
	log.Printf("Warning: Value stored for '%s' expected to be of type 'bool' (actual type: %T)\n", key, value)
	return false
}

func getString(data *schema.ResourceData, key string) string {
	if data == nil {
		return ""
	}
	value := data.Get(key)
	if value == nil {
		return ""
	}
	if bVal, ok := value.(string); ok {
		return bVal
	}
	log.Printf("Warning: Value stored for '%s' expected to be of type 'string' (actual type: %T)\n", key, value)
	return ""
}

// func extractSet(value interface{}) []string {
// 	if value == nil {
// 		return nil
// 	}
// 	if valueSet, ok := value.(*schema.Set); ok {
// 		valueList := valueSet.List()
// 		if valueList != nil {
// 			values := []string{}
// 			for _, v := range valueList {
// 				values = append(values, v.(string))
// 			}
// 			return values
// 		}
// 	}
// 	return nil
// }

func extractMap(value interface{}) map[string]interface{} {
	if value == nil {
		return nil
	}
	if valueSet, ok := value.(*schema.Set); ok {
		valueList := valueSet.List()
		if valueList != nil {
			for _, v := range valueList {
				var mapValue map[string]interface{}
				if mapValue, ok = v.(map[string]interface{}); ok {
					return mapValue
				}
			}
		}
		return nil
	}
	return nil
}
