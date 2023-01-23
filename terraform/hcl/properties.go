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

package hcl

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Properties map[string]any

func (me Properties) Unknowns(unknowns map[string]json.RawMessage) error {
	if len(unknowns) == 0 {
		return nil
	}
	data, err := json.Marshal(unknowns)
	if err != nil {
		return err
	}
	me["unknowns"] = string(data)
	return nil
}

func (me Properties) EncodeSlice(key string, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Slice {
		return fmt.Errorf("type %T is not a slice", v)
	}
	if rv.Len() == 0 {
		return nil
	}
	entries := []any{}
	for idx := 0; idx < rv.Len(); idx++ {
		vElem := rv.Index(idx)
		elem := vElem.Interface()
		if marshaler, ok := elem.(Marshaler); ok {
			marshalled := Properties{}
			if err := marshaler.MarshalHCL(marshalled); err == nil {
				entries = append(entries, marshalled)
			} else {
				return err
			}
		} else {
			return fmt.Errorf("slice entries of type %T are expected to implement hcl.Marshaler but don't", elem)
		}
	}
	me[key] = entries
	return nil
}

func (me Properties) EncodeAll(items map[string]any) error {
	if items == nil {
		return nil
	}
	for k, v := range items {
		if err := me.Encode(k, v); err != nil {
			return err
		}
	}
	return nil
}

type StringSet []string

func (me Properties) setPrimitiveSlice(key string, v any) {
	if reflect.ValueOf(v).Len() == 0 {
		me[key] = nil
	} else {
		me[key] = v
	}
}

func (me Properties) Encode(key string, v any) error {
	if v == nil {
		return nil
	}
	switch t := v.(type) {
	case *string:
		if t == nil {
			me[key] = nil
			return nil
		}
		return me.Encode(key, *t)
	case *bool:
		if t == nil {
			me[key] = nil
			return nil
		}
		return me.Encode(key, *t)
	case *int:
		if t == nil {
			me[key] = nil
			return nil
		}
		return me.Encode(key, *t)
	case *int8:
		if t == nil {
			me[key] = nil
			return nil
		}
		return me.Encode(key, *t)
	case *int16:
		if t == nil {
			me[key] = nil
			return nil
		}
		return me.Encode(key, *t)
	case *int32:
		if t == nil {
			me[key] = nil
			return nil
		}
		return me.Encode(key, *t)
	case *int64:
		if t == nil {
			me[key] = nil
			return nil
		}
		return me.Encode(key, *t)
	case *uint:
		if t == nil {
			me[key] = nil
			return nil
		}
		return me.Encode(key, *t)
	case *uint16:
		if t == nil {
			me[key] = nil
			return nil
		}
		return me.Encode(key, *t)
	case *uint8:
		if t == nil {
			me[key] = nil
			return nil
		}
		return me.Encode(key, *t)
	case *uint32:
		if t == nil {
			me[key] = nil
			return nil
		}
		return me.Encode(key, *t)
	case *uint64:
		if t == nil {
			me[key] = nil
			return nil
		}
		return me.Encode(key, *t)
	case *float32:
		if t == nil {
			me[key] = nil
			return nil
		}
		return me.Encode(key, *t)
	case *float64:
		if t == nil {
			me[key] = nil
			return nil
		}
		return me.Encode(key, *t)
	case StringSet:
		if len(t) > 0 {
			me[key] = t
		} else {
			me[key] = nil
		}
	case []string, []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []float32, []float64, []bool:
		me.setPrimitiveSlice(key, t)
	case string:
		me[key] = t
	case int:
		me[key] = t
	case bool:
		me[key] = t
	case int8:
		me[key] = int(t)
	case int16:
		me[key] = int(t)
	case int32:
		me[key] = int(t)
	case int64:
		me[key] = int(t)
	case uint:
		me[key] = int(t)
	case uint8:
		me[key] = int(t)
	case uint16:
		me[key] = int(t)
	case uint32:
		me[key] = int(t)
	case uint64:
		me[key] = int(t)
	case float32:
		me[key] = float64(t)
	case float64:
		me[key] = float64(t)
	case map[string]json.RawMessage:
		if len(t) == 0 {
			me["unknowns"] = nil
			return nil
		}
		data, err := json.Marshal(t)
		if err != nil {
			return err
		}
		me["unknowns"] = string(data)
	default:
		if marshaller, ok := v.(Marshaler); ok {
			if reflect.ValueOf(v).IsNil() {
				me[key] = nil
				return nil
			}
			marshalled := Properties{}
			if err := marshaller.MarshalHCL(marshalled); err == nil {
				if len(marshalled) > 0 {
					me[key] = []any{marshalled}
				} else {
					me[key] = nil
				}
				return nil
			} else {
				return err
			}
		}
		if reflect.TypeOf(v).Kind() == reflect.Slice {
			if reflect.ValueOf(v).Len() == 0 {
				me[key] = nil
				return nil
			}
			if reflect.TypeOf(v).Elem().Kind() == reflect.String {
				entries := []string{}
				vValue := reflect.ValueOf(v)
				for i := 0; i < vValue.Len(); i++ {
					entries = append(entries, fmt.Sprintf("%v", vValue.Index(i).Interface()))
				}
				if len(entries) > 0 {
					me[key] = entries
				} else {
					me[key] = nil
				}
				return nil
			} else if reflect.TypeOf(v).Elem().Kind() == reflect.Float64 {
				entries := []float64{}
				vValue := reflect.ValueOf(v)
				for i := 0; i < vValue.Len(); i++ {
					entries = append(entries, vValue.Index(i).Interface().(float64))
				}
				if len(entries) > 0 {
					me[key] = entries
				} else {
					me[key] = nil
				}
				return nil
			} else if reflect.TypeOf(v).Elem().Implements(marshalerType) {
				entries := []any{}
				vValue := reflect.ValueOf(v)
				for i := 0; i < vValue.Len(); i++ {
					mars := vValue.Index(i).Interface().(Marshaler)
					marshalled := Properties{}
					if err := mars.MarshalHCL(marshalled); err != nil {
						return err
					}
					entries = append(entries, marshalled)
				}
				if len(entries) > 0 {
					me[key] = entries
				} else {
					me[key] = nil
				}
				return nil
			}

		}
		if reflect.TypeOf(v).Kind() == reflect.String {
			me[key] = fmt.Sprintf("%v", v)
			return nil
		} else if marshaller, ok := v.(Marshaler); ok {
			if reflect.ValueOf(v).IsNil() {
				me[key] = nil
				return nil
			}
			marshalled := Properties{}
			if err := marshaller.MarshalHCL(marshalled); err == nil {
				if len(marshalled) > 0 {
					me[key] = []any{marshalled}
				} else {
					me[key] = nil
				}
				return nil
			} else {
				return err
			}

		} else if reflect.TypeOf(v).Kind() == reflect.Ptr {
			switch reflect.TypeOf(v).Elem().Kind() {
			case reflect.String:
				if reflect.ValueOf(v).IsNil() {
					me[key] = nil
					return nil
				}
				if reflect.ValueOf(v).IsZero() {
					me[key] = nil
					return nil
				}
				if !reflect.ValueOf(v).Elem().IsValid() {
					me[key] = nil
					return nil
				}
				return me.Encode(key, reflect.ValueOf(v).Elem().Interface())
			}
		}
		panic(fmt.Sprintf("unsupported type %T", v))
	}
	return nil
}
