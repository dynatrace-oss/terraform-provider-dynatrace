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

package settings

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings interface {
	MarshalHCL(hcl.Properties) error
	UnmarshalHCL(hcl.Decoder) error
	Schema() map[string]*schema.Schema
}

func SetID(settings Settings, id string) {
	rv := reflect.ValueOf(settings).Elem()
	if field := getIDField(rv); field.IsValid() {
		field.Set(reflect.ValueOf(&id))
	}
}

func getIDField(v reflect.Value) reflect.Value {
	field := v.FieldByName("ID")
	if !field.IsValid() {
		return reflect.Value{}
	}
	if field.Type() != stringPointerType {
		return reflect.Value{}
	}
	return field
}

func getLegacyField(v reflect.Value) reflect.Value {
	if tField, found := v.Type().FieldByName("LegacyID"); found {
		if tField.Tag.Get("json") != "-" {
			return reflect.Value{}
		}
	}
	field := v.FieldByName("LegacyID")
	if !field.IsValid() {
		return reflect.Value{}
	}
	if field.Type() != stringPointerType {
		return reflect.Value{}
	}
	return field
}

func GetFlawedReasons(settings any) []string {
	v := reflect.ValueOf(settings)
	if tField, found := v.Type().Elem().FieldByName("FlawedReasons"); found {
		if tField.Tag.Get("json") != "-" {
			return []string{}
		}
	}
	field := v.Elem().FieldByName("FlawedReasons")
	if !field.IsValid() {
		return []string{}
	}
	if field.Type() != stringSliceType {
		return []string{}
	}
	return field.Interface().([]string)
}

func SupportsFlawedReasons(settings any) bool {
	v := reflect.ValueOf(settings)
	if tField, found := v.Type().Elem().FieldByName("FlawedReasons"); found {
		if tField.Tag.Get("json") != "-" {
			return false
		}
	}
	return true
}

func LegacyID(id string) string {
	objID := &ObjectID{ID: id}
	if e := objID.Decode(); e == nil && len(objID.Key) > 0 {
		return objID.Key
	}
	return ""
}

func UUIDToLong(uid uuid.UUID) (int64, error) {
	var err error
	var data []byte

	if data, err = uid.MarshalBinary(); err != nil {
		return 0, err
	}
	if uid.Variant() == uuid.RFC4122 && uid.Version() == 4 {
		_, msb := new(big.Int).SetBytes(data[:8]).Int64(), data[:8]
		_, lsb := new(big.Int).SetBytes(data[8:]).Int64(), data[8:]
		buf := make([]byte, 10)
		for idx := 0; idx < 6; idx++ {
			buf[idx] = msb[idx]
		}
		for idx := 6; idx < 10; idx++ {
			buf[idx] = lsb[idx-2]
		}
		result, _ := binary.Varint(buf)
		return result, nil
	}
	result, _ := new(big.Int).SetBytes(data[:8]).Int64(), data[:8]
	return result, nil

}

var LegacyObjIDDecode = func(id string) string {
	objID := &ObjectID{ID: id}
	if e := objID.Decode(); e == nil && len(objID.Key) > 0 {
		return objID.Key
	}
	return ""
}

var LegacyLongDecode = func(id string) string {
	key := LegacyObjIDDecode(id)
	if len(key) == 0 {
		return key
	}
	uid, err := uuid.Parse(key)
	if err != nil {
		return err.Error()
	}
	result, err := UUIDToLong(uid)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%v", result)
}

func SetLegacyID(id string, converter func(string) string, v any) {
	if GetLegacyID(v) != nil {
		return
	}
	if converter == nil {
		return
	}
	legacyID := converter(id)
	if len(legacyID) == 0 {
		return
	}
	rv := reflect.ValueOf(v).Elem()
	if field := getLegacyField(rv); field.IsValid() {
		field.Set(reflect.ValueOf(&legacyID))
	}
}

func GetLegacyID(v any) *string {
	rv := reflect.ValueOf(v).Elem()
	if field := getLegacyField(rv); field.IsValid() {
		return field.Interface().(*string)
	}
	return nil
}

func ClearLegacyID(v any) *string {
	rv := reflect.ValueOf(v).Elem()
	if field := getLegacyField(rv); field.IsValid() {
		legacyID := field.Interface().(*string)
		var n *string
		field.Set(reflect.ValueOf(n))
		return legacyID
	}
	return nil
}

type ErrorSettings struct {
	Error error
}

func (me *ErrorSettings) MarshalHCL(properties hcl.Properties) error {
	return nil
}

func (me *ErrorSettings) UnmarshalHCL(decoder hcl.Decoder) error {
	return nil
}

func (me *ErrorSettings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (me *ErrorSettings) FillDemoValues() []string {
	return []string{me.Error.Error()}
}
