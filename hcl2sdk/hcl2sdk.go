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

package hcl2sdk

import (
	"fmt"

	"github.com/dtcookie/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Convert(m map[string]*hcl.Schema) map[string]*schema.Schema {
	result := map[string]*schema.Schema{}
	for k, v := range m {
		result[k] = convertSchema(v)
	}
	return result
}

func convertElem(e interface{}) interface{} {
	if e == nil {
		return nil
	}
	switch re := e.(type) {
	case *hcl.Schema:
		return convertSchema(re)
	case *hcl.Resource:
		return convertResource(re)
	default:
		panic(fmt.Sprintf("unsupported elem type %T", re))
	}
}

func convertType(t hcl.ValueType) schema.ValueType {
	switch t {
	case hcl.TypeBool:
		return schema.TypeBool
	case hcl.TypeFloat:
		return schema.TypeFloat
	case hcl.TypeInt:
		return schema.TypeInt
	case hcl.TypeInvalid:
		return schema.TypeInvalid
	case hcl.TypeList:
		return schema.TypeList
	case hcl.TypeMap:
		return schema.TypeMap
	case hcl.TypeSet:
		return schema.TypeSet
	case hcl.TypeString:
		return schema.TypeString
	default:
		panic(fmt.Sprintf("unsupported type %v", t))
	}
}

func convertResource(r *hcl.Resource) *schema.Resource {
	result := new(schema.Resource)
	result.Schema = Convert(r.Schema)
	return result
}

func convertSchema(s *hcl.Schema) *schema.Schema {
	result := new(schema.Schema)
	result.Description = s.Description
	result.Deprecated = s.Deprecated
	result.Optional = s.Optional
	result.Required = s.Required
	result.MaxItems = s.MaxItems
	result.MinItems = s.MinItems
	result.Sensitive = s.Sensitive
	result.Default = s.Default
	result.Computed = s.Computed
	result.ConflictsWith = s.ConflictsWith
	result.ExactlyOneOf = s.ExactlyOneOf
	result.AtLeastOneOf = s.AtLeastOneOf
	result.RequiredWith = s.RequiredWith
	result.ForceNew = s.ForceNew

	result.Type = convertType(s.Type)
	result.Elem = convertElem(s.Elem)

	return result
}
