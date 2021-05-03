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

package main_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type ResourceTest interface {
	ResourceKey() string
	CreateTestCase(string, string, *testing.T) (*resource.TestCase, error)
	Anonymize(m map[string]interface{})
	URL(id string) string
}

func compareLocalRemote(test ResourceTest, n string, localJSONFile string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var err error
		var localMap map[string]interface{}
		var remoteMap map[string]interface{}

		if rs, ok := s.RootModule().Resources[n]; ok {
			token := testAccProvider.Meta().(*config.ProviderConfiguration).APIToken
			url := test.URL(rs.Primary.ID)
			if remoteMap, err = loadHTTP(url, token); err != nil {
				return err
			}
			if localMap, err = loadLocal(localJSONFile); err != nil {
				return err
			}
			test.Anonymize(localMap)
			test.Anonymize(remoteMap)
			if !deepEqual(localMap, remoteMap) {
				sLocalMap, _ := json.Marshal(localMap)
				sRemoteMap, _ := json.Marshal(remoteMap)
				return fmt.Errorf("--LOCAL--\n%v\n\n\n--REMOTE--\n%v", string(sLocalMap), string(sRemoteMap))
			}
			return nil
		}

		return fmt.Errorf("not found: %s", n)
	}
}

func deepEqual(a interface{}, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil && b != nil {
		return false
	}
	if a != nil && b == nil {
		return false
	}
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	switch ta := a.(type) {
	case map[string]interface{}:
		return deepEqualMap(ta, b.(map[string]interface{}))
	case bool:
		return ta == b.(bool)
	case string:
		return ta == b.(string)
	case float64:
		return ta == b.(float64)
	case []interface{}:
		return deepEqualSlice(ta, b.([]interface{}))
	default:
		panic(fmt.Errorf("unsupported type %T", ta))
	}
}

func deepEqualSlice(a []interface{}, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for _, va := range a {
		found := false
		for _, vb := range b {
			if deepEqual(va, vb) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func deepEqualMap(a map[string]interface{}, b map[string]interface{}) bool {
	for k, va := range a {
		vb, found := b[k]
		if !found {
			return false
		}
		if !deepEqual(va, vb) {
			return false
		}
	}
	return true
}
