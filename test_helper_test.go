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

// type ResourceTest interface {
// 	ResourceKey() string
// 	CreateTestCase(string, string, *testing.T) (*resource.TestCase, error)
// 	Anonymize(m map[string]any)
// 	URL(id string) string
// }

// func compareLocalRemote(test ResourceTest, n string, localJSONFile string, t *testing.T) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		var err error
// 		var localMap map[string]any
// 		var remoteMap map[string]any

// 		if rs, ok := s.RootModule().Resources[n]; ok {
// 			token := testAccProvider.Meta().(*config.ProviderConfiguration).APIToken
// 			url := test.URL(rs.Primary.ID)
// 			if remoteMap, err = loadHTTP(url, token); err != nil {
// 				return err
// 			}
// 			if localMap, err = loadLocal(localJSONFile); err != nil {
// 				return err
// 			}
// 			test.Anonymize(localMap)
// 			test.Anonymize(remoteMap)
// 			if !deepEqual(localMap, remoteMap) {
// 				sLocalMap, _ := json.Marshal(localMap)
// 				sRemoteMap, _ := json.Marshal(remoteMap)
// 				return fmt.Errorf("--LOCAL--\n%v\n\n\n--REMOTE--\n%v", string(sLocalMap), string(sRemoteMap))
// 			}
// 			return nil
// 		}

// 		return fmt.Errorf("not found: %s", n)
// 	}
// }

// func deepEqual(a any, b any) bool {
// 	if a == nil && b == nil {
// 		return true
// 	}
// 	if a == nil && b != nil {
// 		return false
// 	}
// 	if a != nil && b == nil {
// 		return false
// 	}
// 	if reflect.TypeOf(a) != reflect.TypeOf(b) {
// 		return false
// 	}
// 	switch ta := a.(type) {
// 	case map[string]any:
// 		return deepEqualMap(ta, b.(map[string]any))
// 	case bool:
// 		return ta == b.(bool)
// 	case string:
// 		return ta == b.(string)
// 	case float64:
// 		return ta == b.(float64)
// 	case []any:
// 		return deepEqualSlice(ta, b.([]any))
// 	default:
// 		panic(fmt.Errorf("unsupported type %T", ta))
// 	}
// }

// func deepEqualSlice(a []any, b []any) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}
// 	for _, va := range a {
// 		found := false
// 		for _, vb := range b {
// 			if deepEqual(va, vb) {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			return false
// 		}
// 	}
// 	return true
// }

// func deepEqualMap(a map[string]any, b map[string]any) bool {
// 	for k, va := range a {
// 		vb, found := b[k]
// 		if !found {
// 			return false
// 		}
// 		if !deepEqual(va, vb) {
// 			return false
// 		}
// 	}
// 	return true
// }

// func loadHTTP(url string, token string) (map[string]any, error) {
// 	var err error
// 	var request *http.Request
// 	var response *http.Response
// 	var data []byte

// 	if request, err = http.NewRequest("GET", url, nil); err != nil {
// 		return nil, err
// 	}
// 	request.Header.Set("Authorization", "Api-Token "+token)

// 	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(nil)}}
// 	if response, err = client.Do(request); err != nil {
// 		return nil, err
// 	}
// 	defer response.Body.Close()

// 	if data, err = ioutil.ReadAll(response.Body); err != nil {
// 		return nil, err
// 	}

// 	m := map[string]any{}
// 	if err = json.Unmarshal(data, &m); err != nil {
// 		return nil, err
// 	}
// 	return m, nil
// }

// func loadLocal(file string) (map[string]any, error) {
// 	var err error
// 	var data []byte
// 	if data, err = ioutil.ReadFile(file); err != nil {
// 		return nil, err
// 	}
// 	m := map[string]any{}
// 	if err = json.Unmarshal(data, &m); err != nil {
// 		return nil, err
// 	}
// 	return m, nil
// }
