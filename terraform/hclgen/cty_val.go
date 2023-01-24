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

package hclgen

import (
	"github.com/zclconf/go-cty/cty"
)

func ctyVal(v any, indent string) cty.Value {
	switch rv := v.(type) {
	case string:
		return cty.StringVal(rv)
	case *string:
		return cty.StringVal(*rv)
	case int:
		return cty.NumberIntVal(int64(rv))
	case int8:
		return cty.NumberIntVal(int64(rv))
	case int16:
		return cty.NumberIntVal(int64(rv))
	case int32:
		return cty.NumberIntVal(int64(rv))
	case int64:
		return cty.NumberIntVal(int64(rv))
	case *int:
		return cty.NumberIntVal(int64(*rv))
	case *int8:
		return cty.NumberIntVal(int64(*rv))
	case *int16:
		return cty.NumberIntVal(int64(*rv))
	case *int32:
		return cty.NumberIntVal(int64(*rv))
	case *int64:
		return cty.NumberIntVal(int64(*rv))
	case uint:
		return cty.NumberUIntVal(uint64(rv))
	case uint8:
		return cty.NumberUIntVal(uint64(rv))
	case uint16:
		return cty.NumberUIntVal(uint64(rv))
	case uint32:
		return cty.NumberUIntVal(uint64(rv))
	case uint64:
		return cty.NumberUIntVal(uint64(rv))
	case *uint:
		return cty.NumberUIntVal(uint64(*rv))
	case *uint8:
		return cty.NumberUIntVal(uint64(*rv))
	case *uint16:
		return cty.NumberUIntVal(uint64(*rv))
	case *uint32:
		return cty.NumberUIntVal(uint64(*rv))
	case *uint64:
		return cty.NumberUIntVal(uint64(*rv))
	case bool:
		return cty.BoolVal(rv)
	case *bool:
		return cty.BoolVal(*rv)
	case float32:
		return cty.NumberFloatVal(float64(rv))
	case float64:
		return cty.NumberFloatVal(rv)
	case *float32:
		return cty.NumberFloatVal(float64(*rv))
	case *float64:
		return cty.NumberFloatVal(*rv)
	case []float64:
		vals := []cty.Value{}
		for _, v := range rv {
			vals = append(vals, cty.NumberFloatVal(v))
		}
		return cty.ListVal(vals)
	case []int:
		vals := []cty.Value{}
		for _, v := range rv {
			vals = append(vals, cty.NumberIntVal(int64(v)))
		}
		return cty.ListVal(vals)
	case []int8:
		vals := []cty.Value{}
		for _, v := range rv {
			vals = append(vals, cty.NumberIntVal(int64(v)))
		}
		return cty.ListVal(vals)
	case []int16:
		vals := []cty.Value{}
		for _, v := range rv {
			vals = append(vals, cty.NumberIntVal(int64(v)))
		}
		return cty.ListVal(vals)
	case []int32:
		vals := []cty.Value{}
		for _, v := range rv {
			vals = append(vals, cty.NumberIntVal(int64(v)))
		}
		return cty.ListVal(vals)
	case []int64:
		vals := []cty.Value{}
		for _, v := range rv {
			vals = append(vals, cty.NumberIntVal(int64(v)))
		}
		return cty.ListVal(vals)
	case []uint:
		vals := []cty.Value{}
		for _, v := range rv {
			vals = append(vals, cty.NumberUIntVal(uint64(v)))
		}
		return cty.ListVal(vals)
	case []uint8:
		vals := []cty.Value{}
		for _, v := range rv {
			vals = append(vals, cty.NumberUIntVal(uint64(v)))
		}
		return cty.ListVal(vals)
	case []uint16:
		vals := []cty.Value{}
		for _, v := range rv {
			vals = append(vals, cty.NumberUIntVal(uint64(v)))
		}
		return cty.ListVal(vals)
	case []uint32:
		vals := []cty.Value{}
		for _, v := range rv {
			vals = append(vals, cty.NumberUIntVal(uint64(v)))
		}
		return cty.ListVal(vals)
	case []uint64:
		vals := []cty.Value{}
		for _, v := range rv {
			vals = append(vals, cty.NumberUIntVal(uint64(v)))
		}
		return cty.ListVal(vals)
	case []string:
		vals := []cty.Value{}
		for _, v := range rv {
			vals = append(vals, cty.StringVal(v))
		}
		return cty.ListVal(vals)
	case []*string:
		vals := []cty.Value{}
		for _, v := range rv {
			vals = append(vals, cty.StringVal(*v))
		}
		return cty.ListVal(vals)
	default:
		return cty.NilVal
	}
}
