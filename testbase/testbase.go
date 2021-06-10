package testbase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var DisabledTests = map[string]bool{
	"dashboards":         false,
	"notifications":      false,
	"custom_services":    false,
	"aws_credentials":    false,
	"k8s_credentials":    false,
	"azure_credentials":  false,
	"auto_tags":          false,
	"alerting_profiles":  false,
	"management_zones":   false,
	"request_attributes": false,
}

type ResourceTest interface {
	ResourceKey() string
	CreateTestCase(string, string, *testing.T) (*resource.TestCase, error)
	Anonymize(m map[string]interface{})
	URL(id string) string
}

func DeepEqual(a interface{}, b interface{}) bool {
	return deepEqual(a, b, "", nil)
}

func deepEqual(a interface{}, b interface{}, addr string, t *testing.T) bool {
	// if t != nil {
	// 	t.Logf("deepEqual(%v)", addr)
	// }
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
		return deepEqualMap(ta, b.(map[string]interface{}), addr, t)
	case bool:
		return ta == b.(bool)
	case string:
		return ta == b.(string)
	case float64:
		return ta == b.(float64)
	case []interface{}:
		return deepEqualSlice(ta, b.([]interface{}), addr, t)
	default:
		panic(fmt.Errorf("unsupported type %T", ta))
	}
}

func deepEqualSlice(a []interface{}, b []interface{}, addr string, t *testing.T) bool {
	if len(a) != len(b) {
		return false
	}
	for idx, va := range a {
		found := false
		for _, vb := range b {
			if deepEqual(va, vb, fmt.Sprintf("%v[%v]", addr, idx), t) {
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

func deepEqualMap(a map[string]interface{}, b map[string]interface{}, addr string, t *testing.T) bool {
	for k, va := range a {
		vb, found := b[k]
		if !found {
			return false
		}
		if !deepEqual(va, vb, addr+"."+k, t) {
			return false
		}
	}
	return true
}

func LoadHTTP(url string, token string) (map[string]interface{}, error) {
	var err error
	var request *http.Request
	var response *http.Response
	var data []byte

	if request, err = http.NewRequest("GET", url, nil); err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Api-Token "+token)

	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(nil)}}
	if response, err = client.Do(request); err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if data, err = ioutil.ReadAll(response.Body); err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func LoadLocal(file string) (map[string]interface{}, error) {
	var err error
	var data []byte
	if data, err = ioutil.ReadFile(file); err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func CompareLocalRemote(test ResourceTest, n string, localJSONFile string, t *testing.T) resource.TestCheckFunc {
	return CompareLocalRemoteExt(test, n, localJSONFile, t, false)
}

func CompareLocalRemoteExt(test ResourceTest, n string, localJSONFile string, t *testing.T, loadHTTPOnly bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var err error
		var localMap map[string]interface{}
		var remoteMap map[string]interface{}

		if rs, ok := s.RootModule().Resources[n]; ok {
			token := TestAccProvider.Meta().(*config.ProviderConfiguration).APIToken
			url := test.URL(rs.Primary.ID)
			if remoteMap, err = LoadHTTP(url, token); err != nil {
				return err
			}
			if !loadHTTPOnly {
				if localMap, err = LoadLocal(localJSONFile); err != nil {
					return err
				}
				test.Anonymize(localMap)
				test.Anonymize(remoteMap)
				if !deepEqual(localMap, remoteMap, "", t) {
					sLocalMap, _ := json.Marshal(localMap)
					sRemoteMap, _ := json.Marshal(remoteMap)
					return fmt.Errorf("--LOCAL--\n%v\n\n\n--REMOTE--\n%v", string(sLocalMap), string(sRemoteMap))
				}
			}
			return nil
		}

		return fmt.Errorf("not found: %s", n)
	}
}
